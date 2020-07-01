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

BASE_DIR := $(shell pwd)
BUILD_DIR := $(BASE_DIR)/build/out
DIST_DIR := $(BASE_DIR)/build/dist
VERSION ?= $(shell git describe --exact-match 2> /dev/null || \
                 git describe --match=$(git rev-parse --short=8 HEAD) \
		 --always --dirty --abbrev=8)
BUILD_TGT := soda-api-$(VERSION)-linux-amd64

all: build

ubuntu-dev-setup:
	sudo apt-get update && sudo apt-get install -y \
	  build-essential gcc librados-dev librbd-dev

build: prebuild osdsapiserver osdsctl

prebuild:
	mkdir -p $(BUILD_DIR)

.PHONY: osdsdock osdslet osdsapiserver osdsctl docker test protoc goimports

osdsdock:
	go build -ldflags '-w -s' -o $(BUILD_DIR)/bin/osdsdock github.com/sodafoundation/dock/cmd/osdsdock

osdslet:
	go build -ldflags '-w -s' -o $(BUILD_DIR)/bin/osdslet github.com/sodafoundation/controller/cmd/osdslet

osdsapiserver:
	go build -ldflags '-w -s' -o $(BUILD_DIR)/bin/osdsapiserver github.com/sodafoundation/api/cmd/osdsapiserver

osdsctl:
	go build -ldflags '-w -s' -o $(BUILD_DIR)/bin/osdsctl github.com/sodafoundation/api/osdsctl

metricexporter:
	go build -ldflags '-w -s' -o $(BUILD_DIR)/bin/lvm_exporter github.com/sodafoundation/api/contrib/exporters/lvm_exporter

docker: build
	cp $(BUILD_DIR)/bin/osdsapiserver ./cmd/osdsapiserver
	docker build cmd/osdsapiserver -t sodafoundation/apiserver:latest

test: build
	install/CI/test

protoc:
	cd pkg/model/proto && protoc --go_out=plugins=grpc:. model.proto

goimports:
	goimports -w $(shell go list -f {{.Dir}} ./... |grep -v /vendor/)

clean:
	rm -rf $(BUILD_DIR) ./cmd/osdsapiserver/osdsapiserver ./cmd/osdslet/osdslet ./cmd/osdsdock/osdsdock

version:
	@echo ${VERSION}

.PHONY: dist
dist: build
	( \
	    rm -fr $(DIST_DIR) && mkdir $(DIST_DIR) && \
	    cd $(DIST_DIR) && \
	    mkdir $(BUILD_TGT) && \
	    cp -r $(BUILD_DIR)/bin $(BUILD_TGT)/ && \
	    cp $(BASE_DIR)/LICENSE $(BUILD_TGT)/ && \
	    cp $(BASE_DIR)/openapi-spec/swagger.yaml $(BUILD_TGT)/ && \
	    zip -r $(DIST_DIR)/$(BUILD_TGT).zip $(BUILD_TGT) && \
	    tar zcvf $(DIST_DIR)/$(BUILD_TGT).tar.gz $(BUILD_TGT) && \
	    tree \
	)
