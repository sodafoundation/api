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

ETCD_DIR=etcd-v3.2.0-linux-amd64
OPENSDS_DIR=${HOME}/gopath/src/github.com/opensds/opensds
IMAGE_PATH=${HOME}/lvm.img
DEIVCE_PATH=/dev/loop1
VG_NAME=vg001

# Install some lvm tools.
sudo apt-get install -y lvm2

if [ -z $HOME ];then
	echo "home path not exist"
	exit
fi

pvoutput=`pvdisplay`
vgoutput=`vgdisplay`

if [ ! -f ${IMAGE_PATH} ]; then
	dd if=/dev/zero of=${IMAGE_PATH} bs=1GB count=20
fi
if [ ! -n "$pvoutput" ]; then
	# Create a new physical volume.
	losetup ${DEIVCE_PATH} ${IMAGE_PATH}
	pvcreate ${DEIVCE_PATH}
fi
if [ ! -n "$vgoutput" ]; then
	# Add pv in volume group.
	vgcreate ${VG_NAME} ${DEIVCE_PATH}
fi

# Run etcd daemon in background.
cd ${HOME}/${ETCD_DIR}
nohup sudo ./etcd > nohup.out 2> nohup.err < /dev/null &

# Config opensds backend info.
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
	enabled_backends = lvm

	[lvm]
	name = lvm
	description = LVM Test
	driver_name = lvm
	config_path = /etc/opensds/driver/lvm.yaml

	[database]
	endpoint = localhost:2379,localhost:2380
	driver = etcd
	' >> /etc/opensds/opensds.conf
fi
if [ ! -f /etc/opensds/driver/lvm.yaml ]; then
	echo '
	pool:
	  vg001:
	    diskType: SSD
	    iops: 1000
	    bandwidth: 1000
	    AZ: default
	' >> /etc/opensds/driver/lvm.yaml
fi

# Run osdsdock and osdslet daemon in background.
cd ${OPENSDS_DIR}
nohup sudo build/out/bin/osdsdock > nohup.out 2> nohup.err < /dev/null &
nohup sudo build/out/bin/osdslet > nohup.out 2> nohup.err < /dev/null &
