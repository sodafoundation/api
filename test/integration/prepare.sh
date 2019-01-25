#!/bin/bash

# Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

# Keep track of the script directory
TOP_DIR=$(cd $(dirname "$0") && pwd)

# OpenSDS Root directory
OPENSDS_DIR=$(cd $TOP_DIR/../.. && pwd)
OPENSDS_CONF=/etc/opensds/opensds.conf

# Config backend info.
mkdir -p /etc/opensds
cat > ${OPENSDS_CONF} << OPENSDS_GLOBAL_CONFIG_DOC
[osdsapiserver]
api_endpoint = 0.0.0.0:50040
log_file = /var/log/opensds/osdsapiserver.log

[osdslet]
api_endpoint = 0.0.0.0:50049
log_file = /var/log/opensds/osdslet.log

[osdsdock]
api_endpoint = 0.0.0.0:50050
log_file = /var/log/opensds/osdsdock.log
# Choose the type of dock resource, only support 'provisioner' and 'attacher'.
dock_type = provisioner
# Specify which backends should be enabled, sample,ceph,cinder,lvm and so on.
enabled_backends = sample

[sample]
name = sample
description = Sample backend for testing
driver_name = default

[database]
# Enabled database types, such as etcd, mysql, fake, etc.
driver = fake
OPENSDS_GLOBAL_CONFIG_DOC

# Create certs
export OPENSSL_CONF="${OPENSDS_DIR}"/script/devsds/lib/openssl.cnf
source "${OPENSDS_DIR}"/script/devsds/lib/certificate.sh
osds::certificate::install

# Run osdsdock and osdslet daemon in background.
cd ${OPENSDS_DIR}
sudo ${OPENSDS_DIR}/build/out/bin/osdsdock -daemon
sudo ${OPENSDS_DIR}/build/out/bin/osdslet -daemon
sudo ${OPENSDS_DIR}/build/out/bin/osdsapiserver -daemon
