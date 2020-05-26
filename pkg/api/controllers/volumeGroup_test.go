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

package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	c "github.com/sodafoundation/api/pkg/context"
	"github.com/sodafoundation/api/pkg/db"
	"github.com/sodafoundation/api/pkg/model"
	. "github.com/sodafoundation/api/testutils/collection"
	dbtest "github.com/sodafoundation/api/testutils/db/testing"
)

func init() {
	beego.Router("/v1beta/block/volumeGroups", &VolumeGroupPortal{}, "post:CreateVolumeGroup;get:ListVolumeGroups")
	beego.Router("/v1beta/block/volumeGroups/:groupId", &VolumeGroupPortal{}, "put:UpdateVolumeGroup;get:GetVolumeGroup;delete:DeleteVolumeGroup")
}

func TestCreateVolumeGroup (t *testing.T) {
	var fakeBody = `{
		"name": "volumeGroup-demo",
		"description": "volume group test",
		"profiles": ["993c87dc-1928-498b-9767-9da8f901d6ce",
    		"90d667f0-e9a9-427c-8a7f-cc714217c7bd"]
	}`

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("CreateVolumeGroup", c.NewAdminContext(), &model.VolumeGroupSpec{
			BaseModel:   &model.BaseModel{},
			Name:        "volumeGroup-demo",
			Description: "volume group test",
			Profiles: []string {"993c87dc-1928-498b-9767-9da8f901d6ce",
			"90d667f0-e9a9-427c-8a7f-cc714217c7bd"},
			Status: "creating",
			AvailabilityZone: "default",
			}).Return(&SampleVolumeGroups[1], nil)
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/block/volumeGroups", strings.NewReader(fakeBody))
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.VolumeGroupSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SampleVolumeGroups[1])
	})
}

func TestUpdateVolumeGroup(t *testing.T) {
	var jsonStr = []byte(`{
		"id": "3769855c-a102-11e7-b772-17b880d2f555",
		"name": "volumeGroup-demo",
		"description": "volumeGroup test"
	}`)
	var expectedJson = []byte(`{
		"id": "084bf71e-a102-11e7-88a8-e31fe6d52248",
		"createdAt": "2017-07-10T14:36:58.014Z",
		"updatedAt": "2017-07-10T14:36:58.014Z",
		"tenantId": "string",
		"userId": "string",
		"name": "volumeGroup-demo",
		"description": "volume group test",
		"availabilityZone": "default",
		"status": "string",
		"poolId": "string",
		"profiles": [
			993c87dc-1928-498b-9767-9da8f901d6ce",
			"90d667f0-e9a9-427c-8a7f-cc714217c7bd"
  		]
	}`)
	var expected model.VolumeGroupSpec
	json.Unmarshal(expectedJson, &expected)

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		volumegroup := model.VolumeGroupSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&volumegroup)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateVolumeGroup", c.NewAdminContext(), volumegroup.Id, &volumegroup).
			Return(&expected, nil)
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/block/volumeGroups/3769855c-a102-11e7-b772-17b880d2f555", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.VolumeGroupSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &expected)
	})

	t.Run("Should return 500 if update volume group with bad request", func(t *testing.T) {
		volumegroup := model.VolumeGroupSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&volumegroup)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateVolumeGroup", c.NewAdminContext(), volumegroup.Id, &volumegroup).
			Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/block/volumeGroups/3769855c-a102-11e7-b772-17b880d2f555", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})

}

func TestListVolumeGroups(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		var sampleVGs = []*model.VolumeGroupSpec{&SampleVolumeGroups[0]}
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListVolumeGroupsWithFilter", c.NewAdminContext(), m).Return(sampleVGs, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/volumeGroups?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output []*model.VolumeGroupSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, sampleVGs)
	})

	t.Run("Should return 500 if list volume groups with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListVolumeGroupsWithFilter", c.NewAdminContext(), m).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/volumeGroups?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestGetVolumeGroup(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolumeGroup", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f555").Return(&SampleVolumeGroups[0], nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/volumeGroups/3769855c-a102-11e7-b772-17b880d2f555", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.VolumeGroupSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SampleVolumeGroups[0])
	})

	t.Run("Should return 404 if get volume group with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolumeGroup", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f555").Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/volumeGroups/3769855c-a102-11e7-b772-17b880d2f555", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}

func TestDeleteVolumeGroup(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolumeGroup", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f555").Return(
			&SampleVolumeGroups[0], nil)
		mockClient.On("DeleteVolumeGroup", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f555").Return(nil)
		db.C = mockClient

		r, _ := http.NewRequest("DELETE",
			"/v1beta/block/volumeGroups/3769855c-a102-11e7-b772-17b880d2f555", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 200)
	})

	t.Run("Should return 404 if delete volume group with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolumeGroup", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f555").Return(
			nil, errors.New("Invalid resource uuid"))
		db.C = mockClient

		r, _ := http.NewRequest("DELETE",
			"/v1beta/block/volumeGroups/3769855c-a102-11e7-b772-17b880d2f555", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}


