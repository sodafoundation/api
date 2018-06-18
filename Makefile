# Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

.PHONY: all build protoc osdsdock osdslet osdsctl docker clean


build:osdsdock osdslet osdsctl

all:package build

package:
	sudo apt-get update && sudo apt-get install -y \
	  build-essential gcc librados-dev librbd-dev

osdsdock:
	mkdir -p  ./build/out/bin/
	go build -o ./build/out/bin/osdsdock github.com/opensds/opensds/cmd/osdsdock

osdslet:
	mkdir -p  ./build/out/bin/
	go build -o ./build/out/bin/osdslet github.com/opensds/opensds/cmd/osdslet

osdsctl:
	mkdir -p  ./build/out/bin/
	go build -o ./build/out/bin/osdsctl github.com/opensds/opensds/osdsctl

docker:build
	cp ./build/out/bin/osdsdock ./cmd/osdsdock
	cp ./build/out/bin/osdslet ./cmd/osdslet
	docker build cmd/osdsdock -t opensdsio/opensds-dock:latest
	docker build cmd/osdslet -t opensdsio/opensds-controller:latest

test:build
	script/CI/test

protoc:
	cd pkg/dock/proto && protoc --go_out=plugins=grpc:. dock.proto

clean:
	rm -rf ./build ./cmd/osdslet/osdslet ./cmd/osdsdock/osdsdock
