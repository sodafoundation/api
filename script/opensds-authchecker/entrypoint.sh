#!/bin/bash

# Copyright (c) 2019 Huawei Technologies Co., Ltd. All Rights Reserved.
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

HOST_IP=`/sbin/ifconfig eth0 | grep 'inet addr' | cut -d: -f2 | awk '{print $1}'`
#. /opt/stack/devstack/openrc admin admin

#openstack endpoint set 04198085606f489bbeb03e3de534774d  --url 'http://172.17.0.19/identity'
sed -i "s,^admin_endpoint.*$,admin_endpoint = http://$HOST_IP/identity,g" /etc/keystone/keystone.conf
sed -i "s,^public_endpoint.*$,public_endpoint = http://$HOST_IP/identity,g" /etc/keystone/keystone.conf

systemctl restart devstack@keystone.service

/sbin/init
