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

# default backend list
BACKEND_LIST=${BACKEND_LIST:-lvm}

osds::usage(){
    cat  << OSDS_HELP_INFO_DOC
Usage:
  $(basename $0) [-h|--help]
  $(basename $0) [-b|--backends xxx]
Flags:
  -h, --help     Print this information.
  -b, --backends Specify backend list,separated by a comma (default: "lvm").
OSDS_HELP_INFO_DOC
}

osds::backendlist_check(){
    local backendlist=$1
    for backend in $(echo $backendlist | tr "," " ");do
        case $backend in
        lvm|ceph)
        ;;
        *)
        echo "Error: backends must be one of lvm,ceph" >&2
        exit -1
        ;;
        esac
    done
}

# Parse parameter first
case "$# $1" in
    "0 ")
    echo "Not specified the backend, using default(lvm)"
    ;;
    "2 -b"|"2 --backends")
    BACKEND_LIST=$2
    osds::backendlist_check $BACKEND_LIST
    ;;
    "1 -h"|"2 --help")
    osds::usage
    exit 0
    ;;
     *)
    osds::usage
    exit 1
    ;;
esac

# Print the commands being run so that we can see the command that triggers
# an error.  It is also useful for following along as the install occurs.
set -o xtrace
set -o errexit

# Keep track of the script directory
TOP_DIR=$(cd $(dirname "$0") && pwd)
# OpenSDS source code root directory
OPENSDS_DIR=$(cd $TOP_DIR/../.. && pwd)

# OpenSDS configuration directory
OPENSDS_CONFIG_DIR=${OPENSDS_CONFIG_DIR:-/etc/opensds}
OPENSDS_DRIVER_CONFIG_DIR=${OPENSDS_CONFIG_DIR}/driver

mkdir -p $OPENSDS_DRIVER_CONFIG_DIR

# Temporary directory for testing
OPT_DIR=/opt/opensds
OPT_BIN=$OPT_DIR/bin
mkdir -p $OPT_BIN
export PATH=$OPT_BIN:$PATH

# Store backend list.
echo -n  $BACKEND_LIST > $OPT_DIR/backend.list

# Echo text to the log file, summary log file and stdout
# osds::echo_summary "something to say"
function osds::echo_summary {
    echo -e $@ >&6
}

# Echo text only to stdout, no log files
# osds::echo_nolog "something not for the logs"
function osds::echo_nolog {
    echo $@ >&3
}

# Log file
LOGFILE=/var/log/opensds/devsds.log
TIMESTAMP_FORMAT=${TIMESTAMP_FORMAT:-"%F-%H%M%S"}
LOGDAYS=${LOGDAYS:-7}
CURRENT_LOG_TIME=$(date "+$TIMESTAMP_FORMAT")

# Clean up old log files.  Append '.*' to the user-specified
# ``LOGFILE`` to match the date in the search template.
LOGFILE_DIR="${LOGFILE%/*}"           # dirname
LOGFILE_NAME="${LOGFILE##*/}"         # basename
mkdir -p $LOGFILE_DIR
find $LOGFILE_DIR -maxdepth 1 -name $LOGFILE_NAME.\* -mtime +$LOGDAYS -exec rm {} \;
LOGFILE=$LOGFILE.${CURRENT_LOG_TIME}
SUMFILE=$LOGFILE.summary.${CURRENT_LOG_TIME}

# Set fd 3 to a copy of stdout. So we can set fd 1 without losing
# stdout later.
exec 3>&1
# Set fd 1 and 2 to write the log file
exec 1> >( $TOP_DIR/tools/outfilter.py -v -o "${LOGFILE}" ) 2>&1
# Set fd 6 to summary log file
exec 6> >( $TOP_DIR/tools/outfilter.py -o "${SUMFILE}" )

osds::echo_summary "install.sh log $LOGFILE"

# Specified logfile name always links to the most recent log
ln -sf $LOGFILE $LOGFILE_DIR/$LOGFILE_NAME
ln -sf $SUMFILE $LOGFILE_DIR/$LOGFILE_NAME.summary

source $TOP_DIR/lib/util.sh
source $TOP_DIR/lib/etcd.sh
osds::etcd::start

# Set global configuration.
cat > $OPENSDS_CONFIG_DIR/opensds.conf << OPENSDS_GLOBAL_CONFIG_DOC
[osdslet]
api_endpoint = 0.0.0.0:50040
graceful = True
log_file = /var/log/opensds/osdslet.log
socket_order = inc

[osdsdock]
api_endpoint = localhost:50050
log_file = /var/log/opensds/osdsdock.log
# Specify which backends should be enabled, sample,ceph,cinder,lvm and so on.
enabled_backends = $BACKEND_LIST

[database]
endpoint = localhost:$ETCD_PORT,localhost:$ETCD_PEER_PORT
driver = etcd
OPENSDS_GLOBAL_CONFIG_DOC

for backend in $(echo $BACKEND_LIST | tr "," " "); do
    case $backend in
        "lvm")
        source $TOP_DIR/lib/lvm.sh
        osds::lvm::init
        ;;
        "ceph")
        source $TOP_DIR/lib/ceph.sh
        osds::ceph::init
        ;;
    esac
done


# Run osdsdock and osdslet daemon in background.
(
cd ${OPENSDS_DIR}
sudo build/out/bin/osdslet --daemon --alsologtostderr
sudo build/out/bin/osdsdock --daemon --alsologtostderr

osds::echo_summary "Waiting for osdslet to come up."
osds::util::wait_for_url localhost:50040 "osdslet" 0.25 80

export OPENSDS_ENDPOINT=http://localhost:50040
build/out/bin/osdsctl profile create '{"name": "default", "description": "default policy"}'
# Copy bash completion script to system.
cp ${OPENSDS_DIR}/osdsctl/completion/osdsctl.bash_completion /etc/bash_completion.d/

if [ $? == 0 ]; then
osds::echo_summary devsds installed successfully !!
fi
)

set +o xtrace
exec 1>&3
# Force all output to stdout and logs now
exec 1> >( tee -a "${LOGFILE}" ) 2>&1

echo -e "\n\n\n"
echo "Execute command blow to set up ENV OPENSDS_ENDPOINT:"
echo ""
echo "export OPENSDS_ENDPOINT=http://localhost:50040"
echo ""
echo "Enjoy it !!"

# Restore/close logging file descriptors
exec 1>&3
exec 2>&3
exec 3>&-
exec 6>&-
echo ""
