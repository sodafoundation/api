#!/bin/bash

# Copyright 2018 The OpenSDS Authors.
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

# Download and install mockery tool.
go get -v github.com/vektra/mockery/.../

# Auto-generate some fake objects in db module for mocking work.
mockery -name Client -dir ../pkg/db -output ./db/testing -case underscore
# Auto-generate some fake objects in controller and dock module for mocking work.
mockery -name Client -dir ../pkg/controller/client -output ./controller/testing -case underscore
mockery -name Client -dir ../pkg/dock/client -output ./dock/testing -case underscore
# Auto-generate some fake objects in driver module for mocking work.
mockery -name VolumeDriver -dir ../contrib/drivers -output ./driver/testing -case underscore
mockery -name ReplicationDriver -dir ../contrib/drivers -output ./driver/testing -case underscore
