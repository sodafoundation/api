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

osds:usage()
{
    echo "Usage: $(basename $0) [--help|--cleanup|--purge]"
cat  << OSDS_HELP_UNINSTALL_INFO_DOC
Usage:
    $(basename $0) [-h|--help]
    $(basename $0) [-c|--cleanup]
    $(basename $0) [-p|--purge]
Flags:
    -h, --help     Print this information.
    -c, --cleanup  Stop service and clean up some application data.
    -p, --purge    Remove package, config file, log file.
OSDS_HELP_UNINSTALL_INFO_DOC
}

# Parse parameter first
case "$# $*" in
    "0 "|"1 --purge"|"1 -p"|"1 --cleanup"|"1 -c")
    ;;
    "1 -h"|"1 --help")
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
source $TOP_DIR/sdsrc

osds::cleanup() {
    osds::util::serice_operation cleanup
}

osds::uninstall(){
    osds::cleanup
    osds::util::serice_operation uninstall
}

osds::uninstall_purge(){
    osds::uninstall
    osds::util::serice_operation uninstall_purge

    rm /opt/opensds -rf
    rm /etc/opensds -rf
    rm /var/log/opensds -rf
    rm /etc/bash_completion.d/osdsctl.bash_completion -rf
    rm /opt/opensds-security -rf
}

case "$# $*" in
    "1 -c"|"1 --cleanup")
    osds::cleanup
    ;;
    "0 ")
    osds::uninstall
    ;;
    "1 -p"|"1 --purge")
    osds::uninstall_purge
    ;;
     *)
    osds:usage
    exit 1
    ;;
esac
