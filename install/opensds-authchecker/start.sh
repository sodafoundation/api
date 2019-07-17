#!/bin/bash
exec > >(tee -i /start.log)
exec 2>&1

set -x

get_default_host_ip() {
    local host_ip=$1
    local af=$2
    # Search for an IP unless an explicit is set by ``HOST_IP`` environment variable
    if [ -z "$host_ip" ]; then
        host_ip=""
        # Find the interface used for the default route
        host_ip_iface=${host_ip_iface:-$(ip -f "${af}" route | awk '/default/ {print $5}' | head -1)}
        local host_ips
        host_ips=$(LC_ALL=C ip -f "${af}" addr show "${host_ip_iface}" | sed /temporary/d |awk /"${af}"'/ {split($2,parts,"/");  print parts[1]}')
        local ip
        for ip in $host_ips; do
            host_ip=$ip
            break;
        done
    fi
    echo "${host_ip}"
}

HOST_IP=$(get_default_host_ip "$HOST_IP" "inet")

sed -i "s,^admin_endpoint.*$,admin_endpoint = http://$HOST_IP/identity,g" /etc/keystone/keystone.conf
sed -i "s,^public_endpoint.*$,public_endpoint = http://$HOST_IP/identity,g" /etc/keystone/keystone.conf
sed -i "2a <Directory /usr/bin>\n    Require all granted\n</Directory>\n" /usr/share/keystone/wsgi-keystone.conf

openstack-config --set /etc/keystone/keystone.conf database connection 'sqlite:////var/lib/keystone/keystone.db'
keystone-manage credential_setup --keystone-user keystone --keystone-group keystone
keystone-manage fernet_setup --keystone-user keystone --keystone-group keystone
keystone-manage db_sync
keystone-manage bootstrap \
  --bootstrap-project-name admin \
  --bootstrap-username admin \
  --bootstrap-password opensds@123 \
  --bootstrap-role-name admin \
  --bootstrap-service-name keystone \
  --bootstrap-region-id RegionOne \
  --bootstrap-admin-url http://"${HOST_IP}"/identity \
  --bootstrap-public-url http://"${HOST_IP}"/identity \
  --bootstrap-internal-url http://"${HOST_IP}"/identity

ln -s /usr/share/keystone/wsgi-keystone.conf /etc/httpd/conf.d/
systemctl enable httpd.service
systemctl start httpd.service

export OS_IDENTITY_API_VERSION="3"
export OS_AUTH_URL="http://$HOST_IP/identity"
export OS_USER_DOMAIN_ID="default"
export OS_PROJECT_DOMAIN_ID="default"
export OS_PROJECT_NAME="admin"
export OS_USERNAME="admin"
export OS_PASSWORD="admin"

OPENSDS_VERSION=${OPENSDS_VERSION:-v1beta}
OPENSDS_SERVER_NAME=${OPENSDS_SERVER_NAME:-opensds}
STACK_PASSWORD=${STACK_PASSWORD:-opensds@123}
MULTICLOUD_SERVER_NAME=${MULTICLOUD_SERVER_NAME:-multicloud}
MULTICLOUD_VERSION=${MULTICLOUD_VERSION:-v1}
chmod 666 /var/lib/keystone/keystone.db

# for_hotpot
openstack user create --domain default --password "${STACK_PASSWORD}" "${OPENSDS_SERVER_NAME}"
openstack project create service
openstack role add --project service --user opensds admin
openstack group create service
openstack group add user service opensds
openstack role add service --project service --group service
openstack group create admins
openstack group add user admins admin
openstack service create --name "opensds${OPENSDS_VERSION}" --description "OpenSDS Block Storage" "opensds${OPENSDS_VERSION}"
openstack endpoint create --region RegionOne "opensds${OPENSDS_VERSION}" public "http://${HOST_IP}:50040/${OPENSDS_VERSION}/%(tenant_id)s"
openstack endpoint create --region RegionOne "opensds${OPENSDS_VERSION}" internal "http://${HOST_IP}:50040/${OPENSDS_VERSION}/%(tenant_id)s"
openstack endpoint create --region RegionOne "opensds${OPENSDS_VERSION}" admin "http://${HOST_IP}:50040/${OPENSDS_VERSION}/%(tenant_id)s"

# for_gelato
openstack user create --domain default --password "${STACK_PASSWORD}" "${MULTICLOUD_SERVER_NAME}"
openstack role add --project service --user "${MULTICLOUD_SERVER_NAME}" admin
openstack group add user service "${MULTICLOUD_SERVER_NAME}"
openstack service create --name "multicloud${MULTICLOUD_VERSION}" --description "Multi-cloud Block Storage" "multicloud${MULTICLOUD_VERSION}"
openstack endpoint create --region RegionOne "multicloud${MULTICLOUD_VERSION}" public "http://${HOST_IP}:8089/${MULTICLOUD_VERSION}/%(tenant_id)s"
openstack endpoint create --region RegionOne "multicloud${MULTICLOUD_VERSION}" internal "http://${HOST_IP}:8089/${MULTICLOUD_VERSION}/%(tenant_id)s"
openstack endpoint create --region RegionOne "multicloud${MULTICLOUD_VERSION}" admin "http://${HOST_IP}:8089/${MULTICLOUD_VERSION}/%(tenant_id)s"

