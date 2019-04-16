#unit tests
.PHONY: unit_test
ifeq ($(WHAT),)
       TEST_DIR="./pkg/"
else
       TEST_DIR=${WHAT}	
endif

export GOARCHAIUS_CONFIG_PATH=$(CURDIR)
test:
	find ${TEST_DIR} -name "*_test.go"|xargs -i dirname {}|uniq|xargs -i go test ${T} {}



# verify
.PHONY: verify
verify:
	bash -x ./verify.sh

.PHONY: unit_test
unit_test:
	bash -x ./../test/unittest/unittest.sh

.PHONY: integration_test
integration_test:
	#ash test/integration/scripts/execute.sh
	bash -x ./../test/integration/integrationtest.sh


.PHONY: e2etest_build
e2etest_build:
	bash -x ./../test/e2e/e2etest.sh


.PHONY: e2eflowtest_build
e2eflowtest_build:
	bash -x ./../test//e2e/e2eflowtest.sh
	
.PHONY: clean
clean:
	bash -x hack/clean.sh

.PHONY: cross_build
cross_build:
	export GOARCH=arm
	export GOOS="linux"
	export GOARM=6
	export CGO_ENABLED=1
	export CC=arm-linux-gnueabi-gcc
	go build cmd/edge_core.go

.PHONY: small_build
small_build:
	go build -ldflags="-w -s -extldflags -static" cmd/edge_core.go
	upx-ucl --best edge_core
