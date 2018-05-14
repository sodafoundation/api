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

# This script helps new contributors or users set up their local workstation
# for opensds installation and development.

# Temporary directory
OPT_DIR=/opt/opensds
mkdir -p $OPT_DIR

# Golang version
GOLANG_VERSION=${GOLANG_VERSION:-1.9.2}
GOENV_PROFILE=${GOENV_PROFILE:-/etc/profile.d/goenv.sh}

# Log file
LOG_DIR=/var/log/opensds
LOGFILE=${LOGFILE:-/var/log/opensds/bootstrap.log}
mkdir -p $LOG_DIR

# Log function
log() {
    DATE=`date "+%Y-%m-%d %H:%M:%S"`
    USER=$(whoami)
    echo "${DATE} [INFO] $@"
    echo "${DATE} ${USER} execute $0 [INFO] $@" > $LOGFILE
}

log_error ()
{
    DATE=`date "+%Y-%m-%d %H:%M:%S"`
    USER=$(whoami)
    echo "${DATE} [ERROR] $@" 2>&1
    echo "${DATE} ${USER} execute $0 [ERROR] $@" > $LOGFILE
}
log OpenSDS bootstrap starting ...

# load profile
source /etc/profile
# Install Golang environment
if ! which go &>/dev/null; then
    log "Golang is not exist, downloading..."
	wget https://storage.googleapis.com/golang/go${GOLANG_VERSION}.linux-amd64.tar.gz -O $OPT_DIR/go${GOLANG_VERSION}.linux-amd64.tar.gz > /dev/null
	log "tar xzf $OPT_DIR/go${GOLANG_VERSION}.linux-amd64.tar.gz -C /usr/local/"
	tar xzf $OPT_DIR/go${GOLANG_VERSION}.linux-amd64.tar.gz -C /usr/local/
	echo 'export GOROOT=/usr/local/go' > $GOENV_PROFILE
	echo 'export GOPATH=$HOME/gopath' >> $GOENV_PROFILE
	echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' >> $GOENV_PROFILE
	source $GOENV_PROFILE
fi

GOPATH=${GOPATH:-$HOME/gopath}
OPENSDS_ROOT=${GOPATH}/src/github.com/opensds
OPENSDS_DIR=${GOPATH}/src/github.com/opensds/opensds
mkdir -p ${OPENSDS_ROOT}

cd ${OPENSDS_ROOT}
if [ ! -d ${OPENSDS_DIR} ]; then
    log "Download the OpenSDS source code."
	git clone https://github.com/opensds/opensds.git -b master
fi

cd ${OPENSDS_DIR}
if [ ! -d ${OPENSDS_DIR}/build ]; then
    sudo apt-get update > /dev/null
	sudo apt-get install librados-dev librbd-dev -y > /dev/null
	log "Build OpenSDS ..."
	make
fi

log OpenSDS bootstrapped successfully. you can execute 'source /etc/profile' to load golang ENV.
