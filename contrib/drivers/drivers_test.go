// Copyright 2017 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package drivers

import (
	"reflect"
	"testing"

	_ "github.com/opensds/opensds/contrib/drivers/ceph"
	"github.com/opensds/opensds/contrib/drivers/sample"
)

func TestInit(t *testing.T) {
	var rsList = []string{"others"}
	var expectedVd = []VolumeDriver{&sample.Driver{}}

	for i, rs := range rsList {
		if vp := Init(rs); !reflect.DeepEqual(vp, expectedVd[i]) {
			t.Errorf("Expected %v, got %v\n", expectedVd, vp)
		}
	}
}
