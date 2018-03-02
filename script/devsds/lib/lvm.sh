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

# Save trace setting
_XTRACE_LVM=$(set +o | grep xtrace)
set +o xtrace

# Defaults
# --------
# Name of the lvm volume groups to use/create for iscsi volumes
VOLUME_GROUP_NAME=${VOLUME_GROUP_NAME:-opensds-volumes}
DEFAULT_VOLUME_GROUP_NAME=$VOLUME_GROUP_NAME-default
# Backing file name is of the form $VOLUME_GROUP$BACKING_FILE_SUFFIX
BACKING_FILE_SUFFIX=-backing-file
# Default volume size
VOLUME_BACKING_FILE_SIZE=${VOLUME_BACKING_FILE_SIZE:-20G}
LVM_DIR=$OPT_DIR/lvm
DATA_DIR=$LVM_DIR
mkdir -p $LVM_DIR

osds::lvm::pkg_install(){
    sudo apt-get install -y lvm2 tgt open-iscsi
}

osds::lvm::pkg_uninstall(){
    sudo apt-get purge -y lvm2 tgt open-iscsi
}

osds::lvm::create_volume_group(){
    local vg=$1
    local size=$2

    local backing_file=$DATA_DIR/$vg$BACKING_FILE_SUFFIX
    if ! sudo vgs $vg; then
        # Only create if the file doesn't already exists
        [[ -f $backing_file ]] || truncate -s $size $backing_file
        local vg_dev
        vg_dev=`sudo losetup -f --show $backing_file`

        # Only create volume group if it doesn't already exist
        if ! sudo vgs $vg; then
            sudo vgcreate $vg $vg_dev
        fi
    fi
}

osds::lvm::set_configuration(){
cat > $OPENSDS_DRIVER_CONFIG_DIR/lvm.yaml << OPENSDS_LVM_CONFIG_DOC
pool:
  $DEFAULT_VOLUME_GROUP_NAME:
    diskType: NL-SAS
    AZ: default
OPENSDS_LVM_CONFIG_DOC

cat >> $OPENSDS_CONFIG_DIR/opensds.conf << OPENSDS_LVM_GLOBAL_CONFIG_DOC
[lvm]
name = lvm
description = LVM Test
driver_name = lvm
config_path = /etc/opensds/driver/lvm.yaml
OPENSDS_LVM_GLOBAL_CONFIG_DOC
}

osds::lvm::init() {
    local vg=$DEFAULT_VOLUME_GROUP_NAME
    local size=$VOLUME_BACKING_FILE_SIZE

    # Install lvm relative packages.
    osds::lvm::pkg_install
    osds::lvm::create_volume_group $vg $size

    # Remove iscsi targets
    sudo tgtadm --op show --mode target | awk '/Target/ {print $3}' | sudo xargs -r -n1 tgt-admin --delete
    # Remove volumes that already exist.
    osds::lvm::remove_volumes $vg
    osds::lvm::set_configuration
    osds::lvm::set_lvm_filter
}

osds::lvm::remove_volumes() {
    local vg=$1

    # Clean out existing volumes
    sudo lvremove -f $vg
}

osds::lvm::remove_volume_group() {
    local vg=$1

    # Remove the volume group
    sudo vgremove -f $vg
}

osds::lvm::clean_backing_file() {
    local backing_file=$1
    # If the backing physical device is a loop device, it was probably setup by DevStack
    if [[ -n "$backing_file" ]] && [[ -e "$backing_file" ]]; then
        local vg_dev
        vg_dev=$(sudo losetup -j $backing_file | awk -F':' '/'$BACKING_FILE_SUFFIX'/ { print $1}')
        if [[ -n "$vg_dev" ]]; then
            sudo losetup -d $vg_dev
        fi
        rm -f $backing_file
    fi
}

osds::lvm::clean_volume_group() {
    local vg=$1
    osds::lvm::remove_volumes $vg
    osds::lvm::remove_volume_group $vg
    # if there is no logical volume left, it's safe to attempt a cleanup
    # of the backing file
    if [[ -z "$(sudo lvs --noheadings -o lv_name $vg 2>/dev/null)" ]]; then
        osds::lvm::clean_backing_file $DATA_DIR/$vg$BACKING_FILE_SUFFIX
    fi
}

osds::lvm::cleanup(){
    osds::lvm::clean_volume_group $DEFAULT_VOLUME_GROUP_NAME
    osds::lvm::clean_lvm_filter
}

# osds::lvm::clean_lvm_filter() Remove the filter rule set in set_lvm_filter()

osds::lvm::clean_lvm_filter() {
    sudo sed -i "s/^.*# from devsds$//" /etc/lvm/lvm.conf
}

# osds::lvm::set_lvm_filter() Gather all devices configured for LVM and
# use them to build a global device filter
# osds::lvm::set_lvm_filter() Create a device filter
# and add to /etc/lvm.conf.  Note this uses
# all current PV's in use by LVM on the
# system to build it's filter.
osds::lvm::set_lvm_filter() {
    local filter_suffix='"r|.*|" ]  # from devsds'
    local filter_string="global_filter = [ "
    local pv
    local vg
    local line

    for pv_info in $(sudo pvs --noheadings -o name); do
        pv=$(echo -e "${pv_info}" | sed 's/ //g' | sed 's/\/dev\///g')
        new="\"a|$pv|\", "
        filter_string=$filter_string$new
    done
    filter_string=$filter_string$filter_suffix

    osds::lvm::clean_lvm_filter
    sudo sed -i "/# global_filter = \[.*\]/a\    $global_filter$filter_string" /etc/lvm/lvm.conf
    osds::echo_summary "set lvm.conf device global_filter to: $filter_string"
}

# Restore xtrace
$_XTRACE_LVM