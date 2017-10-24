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

/*
This module implements the policy-based scheduling by parsing storage
profiles configured by admin.

*/

package controller

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
)

func TestSearchProfile(t *testing.T) {
	fd := db.NewFakeDbClient()

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

	prf, err := SearchProfile(prfID, fd)
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

	prf, err = SearchProfile(prfID, fd)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedAssignedProfile, prf) {
		t.Fatalf("Expected %v, get %v", expectedAssignedProfile, prf)
	}
}
