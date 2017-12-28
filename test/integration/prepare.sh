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

OPENSDS_DIR=${HOME}/gopath/src/github.com/opensds
OPENSDS_ROOT=${OPENSDS_DIR}/opensds
OPENSDS_CONF=/etc/opensds/opensds.conf

# Config backend info.
if [ ! -f ${OPENSDS_CONF} ]; then
    mkdir -p /etc/opensds
	echo '
	[osdslet]
	api_endpoint = localhost:50040
	graceful = True
	log_file = /var/log/opensds/osdslet.log
	socket_order = inc

	[osdsdock]
	api_endpoint = localhost:50050
	log_file = /var/log/opensds/osdsdock.log
	# Enabled backend types, such as sample, ceph, cinder, etc.
	enabled_backends = sample

	[sample]
	name = sample
	description = Sample backend for testing
	driver_name = default

	[database]
	# Enabled database types, such as etcd, mysql, fake, etc.
	driver = fake
	' >> ${OPENSDS_CONF}
fi

# Run osdsdock and osdslet daemon in background.
cd ${OPENSDS_ROOT}
sudo build/out/bin/osdsdock -daemon
sudo build/out/bin/osdslet -daemon
