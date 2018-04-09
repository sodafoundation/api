// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package dorado

import (
	"fmt"
	"github.com/opensds/opensds/pkg/utils/config"
	"testing"
)

func TestLoadConf(t *testing.T) {
	r := ReplicationDriver{}
	config.CONF.OsdsDock.Backends.HuaweiDorado.ConfigPath = "./testdata/dorado.yaml"
	r.Setup()
	fmt.Println(r.conf)
}
