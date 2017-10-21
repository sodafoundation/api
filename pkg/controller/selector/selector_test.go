// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

/*
This module implements the policy-based scheduling by parsing storage
profiles configured by admin.

*/

package selector

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/model"
)

/*
func TestSelectSupportedPool(t *testing.T) {
	s := NewFakeSelector()

	var expectedPool = &model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		},
		Name:          "rbd-pool",
		Description:   "ceph pool1",
		StorageType:   "block",
		DockId:        "076454a8-65da-11e7-9a65-5f5d9b935b9f",
		TotalCapacity: 200,
		FreeCapacity:  200,
		Parameters: map[string]interface{}{
			"thinProvision":    "false",
			"highAvailability": "true",
		},
	}
	var inputTag = map[string]string{
		"highAvailability": "true",
	}

	// Test if the method would return correct pool when storage tag assigned.
	pol, err := s.SelectSupportedPool(inputTag)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedPool, pol) {
		t.Fatalf("Expected %v, get %v", expectedPool, pol)
	}
}

func TestSelectDock(t *testing.T) {
	s := NewFakeSelector()

	var inputPool = &model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		},
		Name:          "cinder-pool",
		Description:   "cinder pool1",
		StorageType:   "block",
		DockId:        "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		TotalCapacity: 100,
		FreeCapacity:  100,
		Parameters: map[string]interface{}{
			"thinProvision":    "true",
			"highAvailability": "false",
		},
	}
	var expectedDock = &model.DockSpec{
		BaseModel: &model.BaseModel{
			Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		},
		Name:        "cinder",
		Description: "cinder backend service",
		Endpoint:    "localhost:50050",
		DriverName:  "cinder",
		Parameters: map[string]interface{}{
			"thinProvision":    "true",
			"highAvailability": "false",
		},
	}

	// Test if the method would return correct dock when storage pool assigned.
	dck, err := s.SelectDock(inputPool)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedDock, dck) {
		t.Fatalf("Expected %v, got %v", expectedDock, dck)
	}

	var inputVolID = "9193c3ec-771f-11e7-8ca3-d32c0a8b2725"
	expectedDock = &model.DockSpec{
		BaseModel: &model.BaseModel{
			Id: "076454a8-65da-11e7-9a65-5f5d9b935b9f",
		},
		Name:        "ceph",
		Description: "ceph backend service",
		Endpoint:    "localhost:50050",
		DriverName:  "ceph",
		Parameters: map[string]interface{}{
			"thinProvision":    "false",
			"highAvailability": "true",
		},
	}

	// Test if the method would return correct dock when volume id assigned.
	dck, err = s.SelectDock(inputVolID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedDock, dck) {
		t.Fatalf("Expected %v, got %v", expectedDock, dck)
	}
}
*/
