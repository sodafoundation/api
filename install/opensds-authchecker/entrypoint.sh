#!/usr/bin/env bash

# Copyright 2019 The OpenSDS Authors.
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

get_default_host_ip() {
    local host_ip=$1
    local af=$2
    # Search for an IP unless an explicit is set by ``HOST_IP`` environment variable
    if [ -z "$host_ip" ]; then
        host_ip=""
        # Find the interface used for the default route
        host_ip_iface=${host_ip_iface:-$(ip -f $af route | awk '/default/ {print $5}' | head -1)}
        local host_ips
        host_ips=$(LC_ALL=C ip -f $af addr show ${host_ip_iface} | sed /temporary/d |awk /$af'/ {split($2,parts,"/");  print parts[1]}')
        local ip
        for ip in $host_ips; do
            host_ip=$ip
            break;
        done
    fi
    echo $host_ip
}

HOST_IP=$(get_default_host_ip "$HOST_IP" "inet")

sed -i "s,^admin_endpoint.*$,admin_endpoint = http://$HOST_IP/identity,g" /etc/keystone/keystone.conf
sed -i "s,^public_endpoint.*$,public_endpoint = http://$HOST_IP/identity,g" /etc/keystone/keystone.conf

systemctl restart devstack@keystone.service
