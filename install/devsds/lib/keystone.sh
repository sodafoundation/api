#!/usr/bin/env bash

# Copyright 2018 The OpenSDS Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


_XTRACE_KEYSTONE=$(set +o | grep xtrace)
set +o xtrace

# 'stack' user is just for install keystone through devstack
osds::keystone::create_user(){
    if id ${STACK_USER_NAME} &> /dev/null; then
        return
    fi
    sudo useradd -s /bin/bash -d ${STACK_HOME} -m ${STACK_USER_NAME}
    echo "stack ALL=(ALL) NOPASSWD: ALL" | sudo tee /etc/sudoers.d/stack
}


osds::keystone::remove_user(){
    userdel ${STACK_USER_NAME} -f -r
    rm /etc/sudoers.d/stack
}

osds::keystone::devstack_local_conf(){
DEV_STACK_LOCAL_CONF=${DEV_STACK_DIR}/local.conf
cat > $DEV_STACK_LOCAL_CONF << DEV_STACK_LOCAL_CONF_DOCK
[[local|localrc]]
# use TryStack git mirror
GIT_BASE=$STACK_GIT_BASE

# If the ``*_PASSWORD`` variables are not set here you will be prompted to enter
# values for them by ``stack.sh``and they will be added to ``local.conf``.
ADMIN_PASSWORD=$STACK_PASSWORD
DATABASE_PASSWORD=$STACK_PASSWORD
RABBIT_PASSWORD=$STACK_PASSWORD
SERVICE_PASSWORD=$STACK_PASSWORD

# Neither is set by default.
HOST_IP=$HOST_IP

# path of the destination log file.  A timestamp will be appended to the given name.
LOGFILE=\$DEST/logs/stack.sh.log

# Old log files are automatically removed after 7 days to keep things neat.  Change
# the number of days by setting ``LOGDAYS``.
LOGDAYS=2

ENABLED_SERVICES=mysql,key
# Using stable/queens branches
# ---------------------------------
KEYSTONE_BRANCH=$STACK_BRANCH
KEYSTONECLIENT_BRANCH=$STACK_BRANCH
DEV_STACK_LOCAL_CONF_DOCK
chown stack:stack $DEV_STACK_LOCAL_CONF
}

osds::keystone::opensds_conf() {
cat >> $OPENSDS_CONFIG_DIR/opensds.conf << OPENSDS_GLOBAL_CONFIG_DOC
[keystone_authtoken]
memcached_servers = $KEYSTONE_IP:11211
signing_dir = /var/cache/opensds
cafile = /opt/stack/data/ca-bundle.pem
auth_uri = http://$KEYSTONE_IP/identity
project_domain_name = Default
project_name = service
user_domain_name = Default
password = $STACK_PASSWORD
# Whether to encrypt the password. If enabled, the value of the password must be ciphertext.
enable_encrypted = False
# Encryption and decryption tool. Default value is aes. The decryption tool can only decrypt the corresponding ciphertext.
pwd_encrypter = aes
username = $OPENSDS_SERVER_NAME
auth_url = http://$KEYSTONE_IP/identity
auth_type = password

OPENSDS_GLOBAL_CONFIG_DOC

cp $OPENSDS_DIR/examples/policy.json $OPENSDS_CONFIG_DIR
}

osds::keystone::create_user_and_endpoint(){
    . $DEV_STACK_DIR/openrc admin admin
    openstack user create --domain default --password $STACK_PASSWORD $OPENSDS_SERVER_NAME
    openstack role add --project service --user opensds admin
    openstack group create service
    openstack group add user service opensds
    openstack role add service --project service --group service
    openstack group add user admins admin
    openstack service create --name opensds$OPENSDS_VERSION --description "OpenSDS Block Storage" opensds$OPENSDS_VERSION
    openstack endpoint create --region RegionOne opensds$OPENSDS_VERSION public http://$HOST_IP:50040/$OPENSDS_VERSION/%\(tenant_id\)s
    openstack endpoint create --region RegionOne opensds$OPENSDS_VERSION internal http://$HOST_IP:50040/$OPENSDS_VERSION/%\(tenant_id\)s
    openstack endpoint create --region RegionOne opensds$OPENSDS_VERSION admin http://$HOST_IP:50040/$OPENSDS_VERSION/%\(tenant_id\)s
}

osds::keystone::delete_user(){
    . $DEV_STACK_DIR/openrc admin admin
    openstack service delete opensds$OPENSDS_VERSION
    openstack role remove service --project service --group service
    openstack group remove user service opensds
    openstack group delete service    
    openstack user delete $OPENSDS_SERVER_NAME --domain default
}

osds::keystone::delete_redundancy_data() {
    . $DEV_STACK_DIR/openrc admin admin
    openstack project delete demo
    openstack project delete alt_demo
    openstack project delete invisible_to_admin
    openstack user delete demo
    openstack user delete alt_demo
}

osds::keystone::download_code(){
    if [ ! -d ${DEV_STACK_DIR} ];then
        git clone ${STACK_GIT_BASE}/openstack-dev/devstack.git -b ${STACK_BRANCH} ${DEV_STACK_DIR}
        chown stack:stack -R ${DEV_STACK_DIR}
    fi

}

osds::keystone::install(){
    if [ "true" == $USE_CONTAINER_KEYSTONE ] 
    then
        KEYSTONE_IP=$HOST_IP
        docker pull opensdsio/opensds-authchecker:latest
        docker run -d --privileged=true --net=host --name=opensds-authchecker opensdsio/opensds-authchecker:latest
        osds::keystone::opensds_conf
        docker cp $TOP_DIR/lib/keystone.policy.json opensds-authchecker:/etc/keystone/policy.json
    else
        if [ "true" != $USE_EXISTING_KEYSTONE ] 
        then
            KEYSTONE_IP=$HOST_IP
            osds::keystone::create_user
            osds::keystone::download_code
            osds::keystone::opensds_conf

            # If keystone is ready to start, there is no need continue next step.
            if osds::util::wait_for_url http://$HOST_IP/identity "keystone" 0.25 4; then
                return
            fi
            osds::keystone::devstack_local_conf
            cd ${DEV_STACK_DIR}
            su $STACK_USER_NAME -c ${DEV_STACK_DIR}/stack.sh
            osds::keystone::create_user_and_endpoint
            osds::keystone::delete_redundancy_data
            # add opensds customize policy.json for keystone
            cp $TOP_DIR/lib/keystone.policy.json /etc/keystone/policy.json
        else
            osds::keystone::opensds_conf
            cd ${DEV_STACK_DIR}
            osds::keystone::create_user_and_endpoint
        fi    
    fi
}

osds::keystone::cleanup() {
    : #do nothing
}

osds::keystone::uninstall(){
    if [ "true" == $USE_CONTAINER_KEYSTONE ] 
    then
        docker stop opensds-authchecker
        docker rm opensds-authchecker
    else
        if [ "true" != $USE_EXISTING_KEYSTONE ] 
        then
            su $STACK_USER_NAME -c ${DEV_STACK_DIR}/unstack.sh
        else
            osds::keystone::delete_user
        fi
    fi
}

osds::keystone::uninstall_purge(){
    rm $STACK_HOME/* -rf
    osds::keystone::remove_user
}

## Restore xtrace
$_XTRACE_KEYSTONE
