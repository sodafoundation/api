// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

// +build integration

package integration

import (
	"encoding/json"
	"testing"

	"github.com/opensds/opensds/client"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/constants"
)

var c = client.NewClient(&client.Config{
	Endpoint:    "http://localhost:50040",
	AuthOptions: client.NewNoauthOptions(constants.DefaultTenantId)})

//define a variable to store profile ID
var proId = "e4278ae6-01a4-4f5b-a92e-a4cfbdd8d645"

func TestCreateProfile(t *testing.T) string {
	var body = &model.ProfileSpec{
		Name:        "flow",
		Description: "flow policy",
		Extras: model.ExtraSpec{
			"diskType": "SAS",
		},
	}

	prf, err := c.CreateProfile(body)
	if err != nil {
		t.Error("create profile in client failed:", err)
		return
	}

	prfBody, _ := json.MarshalIndent(prf, "", "	")
	&proId = prfBody.id
	t.Log("Create Profile Success")
}

func TestGetProfileDetail(t *testing.T) {
	prfID = proId

	prf, err := c.GetProfile(prfID)
	if err != nil {
		t.Error("get profile in client failed:", err)
		return
	}

	prfBody, _ := json.MarshalIndent(prf, "", "	")
	t.Log(string(prfBody))
	t.Log("Get Profile detail Success")
}

func TestGetProfileList(t *testing.T) {
	proList, err := c.ListProfiles()
	if err != nil {
		t.Error("list profiles in client failed:", err)
		return
	}
	prfsBody, _ := json.MarshalIndent(prfs, "", "	")
	t.Log(string(prfsBody))
	t.Log("Get Profile List Success")
}

func TestDeleteProfile(t testing.T) {
	var prfID = proId
	del, err := c.DeleteProfile(prfID)
	if err != nil {
		t.Error("delete profile in client failed:", err)
		return
	}
	t.Log("Delete Profile Success")
}
