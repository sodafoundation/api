#!/usr/bin/env bash

# Copyright 2019 The OpenSDS Authors.
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

# Default host ip.
HOST_IP=0.0.0.0
# OpenSDS version configuration.
OPENSDS_VERSION=${OPENSDS_VERSION:-v1beta}
# OpenSDS service name in keystone.
OPENSDS_SERVER_NAME=${OPENSDS_SERVER_NAME:-opensds}

# devstack keystone configuration
STACK_GIT_BASE=${STACK_GIT_BASE:-https://git.openstack.org}
STACK_USER_NAME=${STACK_USER_NAME:-stack}
STACK_PASSWORD=${STACK_PASSWORD:-opensds@123}
STACK_HOME=${STACK_HOME:-/opt/stack}
STACK_BRANCH=${STACK_BRANCH:-stable/queens}
DEV_STACK_DIR=$STACK_HOME/devstack

# Multi-Cloud service name in keystone
MULTICLOUD_SERVER_NAME=${MULTICLOUD_SERVER_NAME:-multicloud}
# Multi-cloud 
MULTICLOUD_VERSION=${MULTICLOUD_VERSION:-v1}

osds::keystone::create_user(){
    if id ${STACK_USER_NAME} &> /dev/null; then
        return
    fi
    sudo useradd -s /bin/bash -d ${STACK_HOME} -m ${STACK_USER_NAME}
    echo "stack ALL=(ALL) NOPASSWD: ALL" | sudo tee /etc/sudoers.d/stack
}

osds::keystone::devstack_local_conf(){
DEV_STACK_LOCAL_CONF=${DEV_STACK_DIR}/local.conf
cat > $DEV_STACK_LOCAL_CONF << DEV_STACK_LOCAL_CONF_DOCK
[[local|localrc]]
# use TryStack git mirror
GIT_BASE=$STACK_GIT_BASE
disable_service mysql
enable_service postgresql
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
ENABLED_SERVICES=postgresql,key
# Using stable/queens branches
# ---------------------------------
KEYSTONE_BRANCH=$STACK_BRANCH
KEYSTONECLIENT_BRANCH=$STACK_BRANCH
DEV_STACK_LOCAL_CONF_DOCK
chown stack:stack $DEV_STACK_LOCAL_CONF
}

osds::keystone::create_user_and_endpoint(){
    . $DEV_STACK_DIR/openrc admin admin
    
    # for_hotpot
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
    
    # for_gelato
    openstack user create --domain default --password "$STACK_PASSWORD" "$MULTICLOUD_SERVER_NAME"
    openstack role add --project service --user "$MULTICLOUD_SERVER_NAME" admin
    openstack group add user service "$MULTICLOUD_SERVER_NAME"
    openstack service create --name "multicloud$MULTICLOUD_VERSION" --description "Multi-cloud Block Storage" "multicloud$MULTICLOUD_VERSION"
    openstack endpoint create --region RegionOne "multicloud$MULTICLOUD_VERSION" public "http://$HOST_IP:8089/$MULTICLOUD_VERSION/%(tenant_id)s"
    openstack endpoint create --region RegionOne "multicloud$MULTICLOUD_VERSION" internal "http://$HOST_IP:8089/$MULTICLOUD_VERSION/%(tenant_id)s"
    openstack endpoint create --region RegionOne "multicloud$MULTICLOUD_VERSION" admin "http://$HOST_IP:8089/$MULTICLOUD_VERSION/%(tenant_id)s"
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
        git clone ${STACK_GIT_BASE}/openstack-dev/devstack -b ${STACK_BRANCH} ${DEV_STACK_DIR}
        chown stack:stack -R ${DEV_STACK_DIR}
    fi

}

osds::keystone::install(){
	KEYSTONE_IP=$HOST_IP
	osds::keystone::create_user
	osds::keystone::download_code

	osds::keystone::devstack_local_conf
	cd ${DEV_STACK_DIR}
	su $STACK_USER_NAME -c ${DEV_STACK_DIR}/stack.sh
	osds::keystone::create_user_and_endpoint
	osds::keystone::delete_redundancy_data
}

osds::keystone::install
# set entrypoint.sh as init command
sed -i '14i\/entrypoint\.sh' /etc/rc.local
