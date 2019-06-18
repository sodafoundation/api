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

package client

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/model"
)

var fakeVersion = &VersionMgr{
	Receiver: NewFakeVersionReceiver(),
}

func TestGetVersion(t *testing.T) {
	expected := &model.VersionSpec{
		Name:      "v1beta",
		Status:    "SUPPORTED",
		UpdatedAt: "2017-04-10T14:36:58.014Z",
	}

	vol, err := fakeVersion.GetVersion("v1beta")
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(vol, expected) {
		t.Errorf("Expected %v, got %v", expected, vol)
		return
	}
}

func TestListVersions(t *testing.T) {
	expected := []*model.VersionSpec{
		{
			Name:      "v1beta",
			Status:    "CURRENT",
			UpdatedAt: "2017-07-10T14:36:58.014Z",
		},
	}

	vols, err := fakeVersion.ListVersions()
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(vols, expected) {
		t.Errorf("Expected %v, got %v", expected, vols)
		return
	}
}
