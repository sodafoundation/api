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



# Start the openapi-spec validation.
split_line "Start openapi-spec validation"

wget https://github.com/go-swagger/go-swagger/releases/download/v0.19.0/swagger_linux_amd64
chmod +x ./swagger_linux_amd64 && ./swagger_linux_amd64 validate --stop-on-error $OPENSDS_DIR/openapi-spec/swagger.yaml
rm ./swagger_linux_amd64

# Start unit test.
split_line "Start unit test"
go test -v github.com/opensds/opensds/osdsctl/... -cover
go test -v github.com/opensds/opensds/client/... -cover
go test -v github.com/opensds/opensds/pkg/... -cover
go test -v github.com/opensds/opensds/contrib/... -cover

