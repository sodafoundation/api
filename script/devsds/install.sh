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

set -o xtrace
set -o errexit


# Keep track of the script directory
TOP_DIR=$(cd $(dirname "$0") && pwd)

OPENSDS_DIR=$(cd $TOP_DIR/../.. && pwd)

# Temp dir for testing
OPT_DIR=/opt/opensds/
OPT_BIN=/opt/opensds/bin
mkdir -p $OPT_BIN
export PATH=$OPT_BIN:$PATH
VERBOSE=True


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

source $TOP_DIR/lib/lvm.sh
osds::lvm::init_default

# Create opensds config dir.
mkdir -p /etc/opensds
mkdir -p /etc/opensds/driver

# Config opensds backend info.

cat > /etc/opensds/opensds.conf << OPENSDS_GLOABL_CONFIG_DOC
[osdslet]
api_endpoint = 0.0.0.0:50040
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
endpoint = localhost:$ETCD_PORT,localhost:$ETCD_PEER_PORT
driver = etcd
OPENSDS_GLOABL_CONFIG_DOC

cat > /etc/opensds/driver/lvm.yaml << OPENSDS_LVM_CONFIG_DOC
pool:
  $DEFAULT_VOLUME_GROUP_NAME:
    diskType: NL-SAS
    AZ: default
OPENSDS_LVM_CONFIG_DOC

# Run osdsdock and osdslet daemon in background.

cd ${OPENSDS_DIR}
sudo build/out/bin/osdslet -daemon
sudo build/out/bin/osdsdock -daemon

osds::echo_summary "Waiting for osdslet to come up."
osds::util::wait_for_url localhost:50040 "osdslet" 0.25 80

export OPENSDS_ENDPOINT=http://localhost:50040
build/out/bin/osdsctl profile create '{"name": "default", "description": "default policy"}'
if [ $? == 0 ]; then
osds::echo_summary devsds installed successfully !!
fi

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
