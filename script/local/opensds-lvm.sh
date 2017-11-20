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

ETCD_URL=https://github.com/coreos/etcd/releases/download/v3.2.0
ETCD_TARBALL=etcd-v3.2.0-linux-amd64.tar.gz
ETCD_DIR=etcd-v3.2.0-linux-amd64
GO_PATH=${HOME}/gopath
OPENSDS_DIR=${GO_PATH}/src/github.com/opensds/opensds
VG_NAME=vg001

# Install some lvm tools.
sudo apt-get install -y lvm2

if [ -z $HOME ];then
	echo "home path not exist"
	exit
fi

if [ ! -b /dev/loop1 ]; then
	# Create a new physical volume and add it in volume group.
	dd if=/dev/zero of=${HOME}/lvm.img bs=1GB count=20
	losetup /dev/loop1 ${HOME}/lvm.img
	pvcreate /dev/loop1
	vgcreate ${VG_NAME} /dev/loop1
fi

# If etcd file not exists, download it from etcd release url.
if [ ! -d ${HOME}/${ETCD_DIR} ]; then
	curl -L ${ETCD_URL}/${ETCD_TARBALL} -o ${HOME}/${ETCD_TARBALL}
	tar xzvf ${HOME}/${ETCD_TARBALL} ${HOME}
fi

# Run etcd daemon in background.
cd ${HOME}/${ETCD_DIR}
nohup sudo ./etcd > nohup.out 2> nohup.err < /dev/null &

# Run osdsdock and osdslet daemon in background.
source /etc/profile
cd ${OPENSDS_DIR} && make
nohup sudo build/out/bin/osdsdock > nohup.out 2> nohup.err < /dev/null &
nohup sudo build/out/bin/osdslet > nohup.out 2> nohup.err < /dev/null &
