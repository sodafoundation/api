#!/bin/bash

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

#set -e errexit
# Keep track of the script directory
TOP_DIR=$(cd $(dirname "$0") && pwd)
# OpenSDS Root directory
OPENSDS_DIR=$(cd $TOP_DIR/../.. && pwd)

split_line(){
    echo "================================================================================================"
    echo $*
    echo "================================================================================================"
}

# Start e2e test
split_line "Start e2e test"
sudo $OPENSDS_DIR/install/devsds/install.sh
ps -ef|grep osds
go test -v github.com/opensds/opensds/test/e2e/... -tags e2e

# Start e2e flow test
split_line "Start e2e flow test"
go build -o ./test/e2e/volume-connector github.com/opensds/opensds/test/e2e/connector/ 
go test -v github.com/opensds/opensds/test/e2e/... -tags e2ef
sudo $OPENSDS_DIR/install/devsds/uninstall.sh
