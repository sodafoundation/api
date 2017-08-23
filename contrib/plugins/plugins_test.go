// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package plugins

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/contrib/plugins/ceph"
	"github.com/opensds/opensds/contrib/plugins/upsplugin"
)

func TestInitVP(t *testing.T) {
	var rsList = []string{"ceph", "others"}
	var expectedVp = []VolumePlugin{&ceph.CephPlugin{}, &upsplugin.Plugin{}}

	for i, rs := range rsList {
		if vp := InitVP(rs); !reflect.DeepEqual(vp, expectedVp[i]) {
			t.Errorf("Expected %v, got %v\n", expectedVp, vp)
		}
	}
}
