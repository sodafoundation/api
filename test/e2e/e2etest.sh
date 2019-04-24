#copyright (c) 2019 Huawei Technologies Co., Ltd. All Rights Reserved.
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

set -o errexit
# Keep track of the script directory
TOP_DIR=$(cd $(dirname "$0") && pwd)
# OpenSDS Root directory
OPENSDS_DIR=$(cd $TOP_DIR/../.. && pwd)

split_line(){
    echo "================================================================================================"
    echo $*
    echo "================================================================================================"
}

# Start lvm e2e test
split_line "Start lvm e2e test"
sudo $OPENSDS_DIR/script/devsds/install.sh
ps -ef|grep osds
go test -v github.com/opensds/opensds/test/e2e/... -tags e2e

