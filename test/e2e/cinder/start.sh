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

cd ${OPENSDS_ROOT}

# OpenSDS cluster installation.
script/cluster/bootstrap.sh

# Config backend info.
if [ ! -f /etc/opensds/opensds.conf ]; then
	echo '
	[osdslet]
	api_endpoint = localhost:50040
	graceful = True
	log_file = /var/log/opensds/osdslet.log
	socket_order = inc

	[osdsdock]
	api_endpoint = localhost:50050
	log_file = /var/log/opensds/osdsdock.log
	# Specify which backends should be enabled, sample,ceph,cinder,lvm and so on.
	enabled_backends = cinder

	[cinder]
	name = cinder
	description = Cinder Test
	driver_name = cinder
	config_path = /etc/opensds/driver/cinder.yaml

	[database]
	endpoint = localhost:2379,localhost:2380
	driver = etcd
	' >> /etc/opensds/opensds.conf
fi
if [ ! -f /etc/opensds/driver/cinder.yaml ]; then
	echo '
	authOptions:
	  endpoint: "http://192.168.56.104/identity"
	  domainId: "Default"
	  domainName: "Default"
	  username: "admin"
	  password: "admin"
	  tenantId: "04154b841eb644a3947506c54fa73c76"
	  tenantName: "admin"
	pool:
	  pool1:
	    diskType: SSD
	    iops: 1000
	    bandwidth: 1000
	    AZ: nova-01
	  pool2:
	    diskType: SAS
	    iops: 800
	    bandwidth: 800
	    AZ: nova-01
	' >> /etc/opensds/driver/cinder.yaml
fi

# Run etcd daemon in background.
cd ${HOME}/${ETCD_DIR}
nohup sudo ./etcd > nohup.out 2> nohup.err < /dev/null &

# Run osdsdock and osdslet daemon in background.
cd ${OPENSDS_ROOT}
nohup sudo build/out/bin/osdsdock > nohup.out 2> nohup.err < /dev/null &
nohup sudo build/out/bin/osdslet > nohup.out 2> nohup.err < /dev/null &

# Start e2e test.
go test -v github.com/opensds/opensds/test/e2e/... -tags e2e:cinder
