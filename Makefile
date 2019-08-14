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
BUILD_TGT := opensds-hotpot-$(VERSION)-linux-amd64

DOCKER_TAG ?= latest

PROTOC_VERSION ?= 3.8.0

all: build

ubuntu-dev-setup:
	sudo apt-get update && sudo apt-get install -y \
	  build-essential gcc librados-dev librbd-dev

build: prebuild osdsdock osdslet osdsapiserver osdsctl metricexporter

prebuild:
	mkdir -p $(BUILD_DIR)

.PHONY: osdsdock osdslet osdsapiserver osdsctl docker test protoc goimports

osdsdock:
	go build -ldflags '-w -s' -o $(BUILD_DIR)/bin/osdsdock github.com/opensds/opensds/cmd/osdsdock

osdslet:
	go build -ldflags '-w -s' -o $(BUILD_DIR)/bin/osdslet github.com/opensds/opensds/cmd/osdslet

osdsapiserver:
	go build -ldflags '-w -s' -o $(BUILD_DIR)/bin/osdsapiserver github.com/opensds/opensds/cmd/osdsapiserver

osdsctl:
	go build -ldflags '-w -s' -o $(BUILD_DIR)/bin/osdsctl github.com/opensds/opensds/osdsctl

metricexporter:
	go build -ldflags '-w -s' -o $(BUILD_DIR)/bin/lvm_exporter github.com/opensds/opensds/contrib/exporters/lvm_exporter

docker: build
	cp $(BUILD_DIR)/bin/osdsdock ./cmd/osdsdock
	cp $(BUILD_DIR)/bin/osdslet ./cmd/osdslet
	cp $(BUILD_DIR)/bin/osdsapiserver ./cmd/osdsapiserver
	docker build cmd/osdsdock -t opensdsio/opensds-dock:$(DOCKER_TAG)
	docker build cmd/osdslet -t opensdsio/opensds-controller:$(DOCKER_TAG)
	docker build cmd/osdsapiserver -t opensdsio/opensds-apiserver:$(DOCKER_TAG)

test: build #osds_verify osds_unit_test osds_integration_test osds_e2eflowtest_build osds_e2etest_build
	install/CI/test
# make osds_core
.PHONY: osds_core
osds_core:
	cd osds && $(MAKE)

# unit tests
.PHONY: osds_unit_test
osds_unit_test:
	cd osds && $(MAKE) unit_test

# verify
.PHONY: osds_verify
osds_verify:
	cd osds && $(MAKE) verify

.PHONY: osds_integration_test
osds_integration_test:
	cd osds && $(MAKE) integration_test

.PHONY: osds_e2etest_build
osds_e2etest_build:
	cd osds && $(MAKE) e2etest_build

protoc_precheck:
	@if ! which protoc >/dev/null; then\
		echo "No protoc in $(PATH), consider visiting https://github.com/protocolbuffers/protobuf/releases to get the protoc(version $(PROTOC_VERSION))";\
		exit 1;\
	fi;
	@if [ ! "libprotoc $(PROTOC_VERSION)" = "$(shell protoc --version)" ]; then\
		echo "protoc version should be $(PROTOC_VERSION)";\
		exit 1;\
	fi

protoc: protoc_precheck
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
	    zip -r $(DIST_DIR)/$(BUILD_TGT).zip $(BUILD_TGT) && \
	    tar zcvf $(DIST_DIR)/$(BUILD_TGT).tar.gz $(BUILD_TGT) && \
	    tree \
	)
