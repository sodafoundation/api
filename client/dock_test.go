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

var fd = &DockMgr{
	Receiver: NewFakeDockReceiver(),
}

func TestGetDock(t *testing.T) {
	var dckID = "b7602e18-771e-11e7-8f38-dbd6d291f4e0"
	expected := &model.DockSpec{
		BaseModel: &model.BaseModel{
			Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		},
		Name:        "sample",
		Description: "sample backend service",
		Endpoint:    "localhost:50050",
		DriverName:  "sample",
	}

	dck, err := fd.GetDock(dckID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(dck, expected) {
		t.Errorf("Expected %v, got %v", expected, dck)
		return
	}
}

func TestListDocks(t *testing.T) {
	expected := []*model.DockSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			Name:        "sample",
			Description: "sample backend service",
			Endpoint:    "localhost:50050",
			DriverName:  "sample",
		},
	}

	dcks, err := fd.ListDocks()
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(dcks, expected) {
		t.Errorf("Expected %v, got %v", expected, dcks)
		return
	}
}
