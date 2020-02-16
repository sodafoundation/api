// Copyright 2019 The OpenSDS Authors.
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

var fzn = &AvailabilityZoneMgr{
	Receiver: NewFakeZoneReceiver(),
}

func TestCreateAvailabilityZone(t *testing.T) {
	expected := &model.AvailabilityZoneSpec{
		BaseModel: &model.BaseModel{
			Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		Name:        "default",
		Description: "default zone",
	}

	zn, err := fzn.CreateAvailabilityZone(&model.AvailabilityZoneSpec{})
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(zn, expected) {
		t.Errorf("Expected %v, got %v", expected, zn)
		return
	}
}

func TestGetAvailabilityZone(t *testing.T) {
	var znID = "1106b972-66ef-11e7-b172-db03f3689c9c"
	expected := &model.AvailabilityZoneSpec{
		BaseModel: &model.BaseModel{
			Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		Name:        "default",
		Description: "default zone",
	}

	zn, err := fzn.GetAvailabilityZone(znID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(zn, expected) {
		t.Errorf("Expected %v, got %v", expected, zn)
		return
	}
}

func TestListAvailabilityZone(t *testing.T) {
	expected := []*model.AvailabilityZoneSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
			},
			Name:        "default",
			Description: "default zone",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "2f9c0a04-66ef-11e7-ade2-43158893e017",
			},
			Name:        "test",
			Description: "test zone",
		},
	}

	zns, err := fzn.ListAvailabilityZones()
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(zns, expected) {
		t.Errorf("Expected %v, got %v", expected[1], zns[1])
		return
	}
}

func TestUpdateAvailabilityZone(t *testing.T) {
	expected := &model.AvailabilityZoneSpec{
		BaseModel: &model.BaseModel{
			Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		Name:        "default",
		Description: "default zone",
	}

	zn, err := fzn.UpdateAvailabilityZone("1106b972-66ef-11e7-b172-db03f3689c9c", &model.AvailabilityZoneSpec{})
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(zn, expected) {
		t.Errorf("Expected %v, got %v", expected, zn)
		return
	}
}

func TestDeleteAvailabilityZone(t *testing.T) {
	var znID = "1106b972-66ef-11e7-b172-db03f3689c9c"

	if err := fzn.DeleteAvailabilityZone(znID); err != nil {
		t.Error(err)
		return
	}
}
