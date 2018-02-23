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

osds:usage()
{
    echo "Usage: $(basename $0) [--help|--purge]"
}

# Parse parameter first
case "$# $*" in
    "0 "|"1 --purge")
    ;;
    "1 --help")
    osds:usage
    exit 0
    ;;
     *)
    osds:usage
    exit 1
    ;;
esac

set -o xtrace

# Keep track of the script directory
TOP_DIR=$(cd $(dirname "$0") && pwd)
# Temporary dir for testing
OPT_DIR=/opt/opensds
OPT_BIN=$OPT_DIR/bin

source $TOP_DIR/lib/util.sh
source $TOP_DIR/lib/etcd.sh
source $TOP_DIR/lib/lvm.sh
source $TOP_DIR/lib/ceph.sh

osds::stop()
{
    OSDSLET_PID=$(pgrep osdslet)
    OSDSDOCK_PID=$(pgrep osdsdock)
    if [ ! -z "$OSDSLET_PID" ]; then
        kill $OSDSLET_PID
    fi
    if [ ! -z "$OSDSDOCK_PID" ]; then
        kill $OSDSDOCK_PID
    fi
}

osds::lvm_enabled(){
    cat $OPT_DIR/backend.list | grep lvm
    return $?
}

osds::ceph_enabled(){
    cat $OPT_DIR/backend.list | grep ceph
    return $?
}

osds::cleanup(){
    osds::etcd::cleanup
    osds::lvm_enabled && osds::lvm::cleanup
    osds::ceph_enabled && osds::ceph::cleanup
    osds::stop
}

osds::purge_cleanup(){
    osds::cleanup
    osds::lvm::pkg_uninstall
    rm /opt/opensds -rf
    rm /etc/opensds -rf
    rm /var/log/opensds -rf
    rm /etc/bash_completion.d/osdsctl.bash_completion
}

case "$# $*" in
    "0 ")
    osds::cleanup
    ;;
    "1 --purge")
    osds::purge_cleanup
    ;;
     *)
    osds:usage
    exit 1
    ;;
esac




