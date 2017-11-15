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

func TestSelectProfile(t *testing.T) {
	s := NewFakeSelector()

	// Test if the method would return default profile when no profile id
	// assigned.
	var prfID = ""
	var expectedDefaultProfile = &model.ProfileSpec{
		BaseModel: &model.BaseModel{
			Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		Name:        "default",
		Description: "default policy",
		Extra:       model.ExtraSpec{},
	}

	prf, err := s.SelectProfile(prfID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedDefaultProfile, prf) {
		t.Fatalf("Expected %v, get %v", expectedDefaultProfile, prf)
	}

	// Test if the method would return specified profile when profile id
	// assigned.
	prfID = "2f9c0a04-66ef-11e7-ade2-43158893e017"
	var expectedAssignedProfile = &model.ProfileSpec{
		BaseModel: &model.BaseModel{
			Id: "2f9c0a04-66ef-11e7-ade2-43158893e017",
		},
		Name:        "silver",
		Description: "silver policy",
		Extra: model.ExtraSpec{
			"diskType":  "SAS",
			"iops":      300,
			"bandwidth": 500,
		},
	}

	prf, err = s.SelectProfile(prfID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedAssignedProfile, prf) {
		t.Fatalf("Expected %v, get %v", expectedAssignedProfile, prf)
	}
}

func TestSelectSupportedPool(t *testing.T) {
	s := NewFakeSelector()

	var expectedPool = &model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
		},
		Name:          "sample-pool-01",
		Description:   "This is the first sample storage pool for testing",
		TotalCapacity: int64(100),
		FreeCapacity:  int64(90),
		DockId:        "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		Parameters: map[string]interface{}{
			"diskType":  "SSD",
			"iops":      1000,
			"bandwidth": 1000,
		},
	}
	var inputTag = map[string]interface{}{
		"diskType":  "SSD",
		"iops":      1000,
		"bandwidth": 1000,
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
			Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
		},
		Name:          "sample-pool-01",
		Description:   "This is the first sample storage pool for testing",
		TotalCapacity: int64(100),
		FreeCapacity:  int64(90),
		DockId:        "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		Parameters: map[string]interface{}{
			"diskType":  "SSD",
			"iops":      1000,
			"bandwidth": 1000,
		},
	}
	var expectedDock = &model.DockSpec{
		BaseModel: &model.BaseModel{
			Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		},
		Name:        "sample",
		Description: "sample backend service",
		Endpoint:    "localhost:50050",
		DriverName:  "sample",
	}

	// Test if the method would return correct dock when storage pool assigned.
	dck, err := s.SelectDock(inputPool)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedDock, dck) {
		t.Fatalf("Expected %v, got %v", expectedDock, dck)
	}

	var inputVolID = "bd5b12a8-a101-11e7-941e-d77981b584d8"

	// Test if the method would return correct dock when volume id assigned.
	dck, err = s.SelectDock(inputVolID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedDock, dck) {
		t.Fatalf("Expected %v, got %v", expectedDock, dck)
	}
}
