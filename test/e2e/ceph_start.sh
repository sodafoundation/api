#!/bin/bash

# Copyright 2017 The OpenSDS Authors.
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

OPENSDS_DIR=${GOPATH}/src/github.com/opensds
OPENSDS_ROOT=${OPENSDS_DIR}/opensds
OPENSDS_LOG_DIR=/var/log/opensds
OPENSDS_CONFIG_DIR=/etc/opensds/driver
ETCD_DIR=etcd-v3.2.0-linux-amd64

function log() {
DATE=`date "+%Y-%m-%d %H:%M:%S"`
USER=$(whoami)
echo "${DATE} ${USER} execute $0 [INFO] $@"
}

function log_error ()
{
DATE=`date "+%Y-%m-%d %H:%M:%S"`
USER=$(whoami)
echo "${DATE} ${USER} execute $0 [ERROR] $@" 2>&1
}

function cleanup(){
    rm ${HOME}/${ETCD_DIR}/default.etcd -rf
    killall osdslet osdsdock etcd &>/dev/null
}

# OpenSDS cluster installation.
cd ${OPENSDS_ROOT} && script/cluster/bootstrap.sh

# Import some pre-defined environment variables.
source /etc/profile

[ ! -d $OPENSDS_CONFIG_DIR ] && mkdir -p ${OPENSDS_CONFIG_DIR}
[ ! -d $OPENSDS_LOG_DIR ] && mkdir -p ${OPENSDS_LOG_DIR}

# Config backend info.
cat > /etc/opensds/opensds.conf << OPENSDS_GLOABL_CONFIG_DOC
[osdslet]
api_endpoint = localhost:50040
graceful = True
log_file = $OPENSDS_LOG_DIR/osdslet.log
socket_order = inc

[osdsdock]
api_endpoint = localhost:50050
log_file = $OPENSDS_LOG_DIR/osdsdock.log
# Specify which backends should be enabled, sample,ceph,cinder,lvm and so on.
enabled_backends = ceph

[ceph]
name = ceph
description = Ceph E2E Test
driver_name = ceph
config_path = $OPENSDS_CONFIG_DIR/ceph.yaml

[database]
endpoint = localhost:2379,localhost:2380
driver = etcd
OPENSDS_GLOABL_CONFIG_DOC

cat > ${OPENSDS_CONFIG_DIR}/ceph.yaml <<OPENSDS_CEPH_DIRVER_CONFIG_DOC
configFile: /etc/ceph/ceph.conf
pool:
  "rbd":
    diskType: SSD
    AZ: default
OPENSDS_CEPH_DIRVER_CONFIG_DOC

# Run etcd daemon in background.
cd ${HOME}/${ETCD_DIR}
./etcd &>>${OPENSDS_LOG_DIR}/etcd.log &
# Waiting for the etcd up.
n=1
export ETCDCTL_API=3
while ! ./etcdctl endpoint status &>/dev/null
do
    echo try $n times
    let n++
    if [ $n -ge 10 ]; then
        log_error "The etcd is not up exited"
        cleanup
        exit 1
    fi
    sleep 1
done


# Run osdsdock and osdslet daemon in background.
cd ${OPENSDS_ROOT}
sudo build/out/bin/osdsdock -daemon
sudo build/out/bin/osdslet -daemon

# Start e2e test.
go test -v github.com/opensds/opensds/test/e2e/... -tags e2e

cleanup
exit 0
