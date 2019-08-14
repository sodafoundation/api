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

.PHONY: clean
clean:
	bash -x hack/clean.sh
