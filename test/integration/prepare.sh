#!/bin/bash

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

[osdslet]
api_endpoint = localhost:50049

[osdsdock]
api_endpoint = localhost:50050
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
export OPENSSL_CONF="${OPENSDS_DIR}"/install/devsds/lib/openssl.cnf
source "${OPENSDS_DIR}"/install/devsds/lib/certificate.sh
osds::certificate::install

# Run osdsdock, osdslet and osdsapiserver daemon in background.
cd ${OPENSDS_DIR}
sudo ${OPENSDS_DIR}/build/out/bin/osdsdock -daemon
sudo ${OPENSDS_DIR}/build/out/bin/osdslet -daemon
sudo ${OPENSDS_DIR}/build/out/bin/osdsapiserver -daemon
