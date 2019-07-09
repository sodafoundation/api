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

# Temporary directory
OPT_DIR=/opt/opensds
mkdir -p $OPT_DIR

# Golang version
MINIMUM_GO_VERSION=${MINIMUM_GO_VERSION:-go1.11.1}
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

# if not found, install it.
if [[ -z "$(which go)" ]]; then
    log "Golang is not installed, downloading..."
    wget https://storage.googleapis.com/golang/${MINIMUM_GO_VERSION}.linux-amd64.tar.gz -O $OPT_DIR/${MINIMUM_GO_VERSION}.linux-amd64.tar.gz > /dev/null
    log "tar xzf $OPT_DIR/${MINIMUM_GO_VERSION}.linux-amd64.tar.gz -C /usr/local/"
    tar xzf $OPT_DIR/${MINIMUM_GO_VERSION}.linux-amd64.tar.gz -C /usr/local/
    echo 'export GOROOT=/usr/local/go' > $GOENV_PROFILE
    echo 'export GOPATH=$HOME/go' >> $GOENV_PROFILE
    echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' >> $GOENV_PROFILE
    source $GOENV_PROFILE
fi

# verify go version
IFS=" " read -ra go_version <<< "$(go version)"
if [[ "${MINIMUM_GO_VERSION}" != $(echo -e "${MINIMUM_GO_VERSION}\n${go_version[2]}" | sort -s -t. -k 1,1 -k 2,2n -k 3,3n | head -n1) && "${go_version[2]}" != "devel" ]]; then
    log_error "Detected go version: ${go_version[*]}, OpenSDS requires ${MINIMUM_GO_VERSION} or greater."
    log_error "Please remove golang old version ${go_version[2]}, bootstrap will install ${MINIMUM_GO_VERSION} automatically"
    exit 2
fi

GOPATH=${GOPATH:-$HOME/go}
OPENSDS_ROOT=${GOPATH}/src/github.com/opensds
OPENSDS_DIR=${GOPATH}/src/github.com/opensds/opensds
mkdir -p ${OPENSDS_ROOT}

cd ${OPENSDS_ROOT}
if [ ! -d ${OPENSDS_DIR} ]; then
    log "Downloading the OpenSDS source code..."
    git clone https://github.com/opensds/opensds.git -b master
fi

# make sure 'make' has been installed.
if [[ -z "$(which make)" ]]; then
    log "Installing make ..."
    sudo apt-get install make -y
fi

cd ${OPENSDS_DIR}
if [ ! -d ${OPENSDS_DIR}/build ]; then
    log "Building OpenSDS ..."
    make ubuntu-dev-setup
    make
fi

log OpenSDS bootstrapped successfully. you can execute 'source /etc/profile' to load golang ENV.
