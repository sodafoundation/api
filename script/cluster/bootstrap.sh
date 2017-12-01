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

# This script helps new contributors or users set up their local workstation
# for opensds installation and development.

OPENSDS_DIR=${HOME}/gopath/src/github.com/opensds
OPENSDS_ROOT=${OPENSDS_DIR}/opensds
ETCD_URL=https://github.com/coreos/etcd/releases/download/v3.2.0
ETCD_TARBALL=etcd-v3.2.0-linux-amd64.tar.gz
ETCD_DIR=etcd-v3.2.0-linux-amd64

# Install Golang environment

if ! which go &>/dev/null; then
	wget https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz
	tar xvf go1.9.linux-amd64.tar.gz -C /usr/local/
	echo 'export GOROOT=/usr/local/go' >> /etc/profile
	echo 'export GOPATH=$HOME/gopath' >> /etc/profile
	echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' >> /etc/profile
fi
source /etc/profile

# If etcd file not exists, download it from etcd release url.
if [ ! -d ${HOME}/${ETCD_DIR} ]; then
	curl -L ${ETCD_URL}/${ETCD_TARBALL} -o ${HOME}/${ETCD_TARBALL}
	cd ${HOME}
	tar xzvf ${HOME}/${ETCD_TARBALL}
fi

# OpenSDS Download and Build
if [ ! -d $OPENSDS_DIR ]; then
	mkdir -p ${OPENSDS_DIR}
fi
cd ${OPENSDS_DIR}
if [ ! -d $OPENSDS_ROOT ]; then
	git clone https://github.com/opensds/opensds.git -b development
fi
cd ${OPENSDS_ROOT}
if [ ! -d $OPENSDS_ROOT/build ]; then
	sudo apt-get install librados-dev librbd-dev -y
	make
fi

