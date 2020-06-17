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
	ctx "context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sodafoundation/api/pkg/utils/constants"

	pb "github.com/sodafoundation/api/pkg/model/proto"
	ctrtest "github.com/sodafoundation/api/testutils/controller/testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	c "github.com/sodafoundation/api/pkg/context"
	"github.com/sodafoundation/api/pkg/db"
	"github.com/sodafoundation/api/pkg/model"
	. "github.com/sodafoundation/api/testutils/collection"
	dbtest "github.com/sodafoundation/api/testutils/db/testing"
)

func init() {
	beego.Router("/v1beta/block/volumeGroups", NewFakeVolumeGroupPortal(), "post:CreateVolumeGroup;get:ListVolumeGroups")
	beego.Router("/v1beta/block/volumeGroups/:groupId", NewFakeVolumeGroupPortal(), "put:UpdateVolumeGroup;get:GetVolumeGroup;delete:DeleteVolumeGroup")
}

func NewFakeVolumeGroupPortal() *VolumeGroupPortal {
	mockClient := new(ctrtest.Client)

	mockClient.On("Connect", "localhost:50049").Return(nil)
	mockClient.On("Close").Return(nil)
	mockClient.On("CreateVolumeGroup", ctx.Background(), &pb.CreateVolumeGroupOpts{
		Context: c.NewAdminContext().ToJson(),
	}).Return(&pb.GenericResponse{}, nil)
	mockClient.On("DeleteVolumeGroup", ctx.Background(), &pb.DeleteVolumeGroupOpts{
		Context: c.NewAdminContext().ToJson(),
	}).Return(&pb.GenericResponse{}, nil)

	return &VolumeGroupPortal{
		CtrClient: mockClient,
	}
}

func TestCreateVolumeGroup(t *testing.T) {
	var jsonStr = []byte(`{
		"id": "3769855c-a102-11e7-b772-17b880d2f555",
		"name": "volumeGroup-demo",
  		"description": "volume group test",
  		"profiles": [
    		"993c87dc-1928-498b-9767-9da8f901d6ce",
    		"90d667f0-e9a9-427c-8a7f-cc714217c7bd"
  		]
	}`)

	t.Run("Should return 202 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		var volumeGroup = &model.VolumeGroupSpec{
			BaseModel: &model.BaseModel{
				Id:        "3769855c-a102-11e7-b772-17b880d2f555",
				CreatedAt: time.Now().Format(constants.TimeFormat),
			},
			Status:           "creating",
			AvailabilityZone: "default",
		}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&volumeGroup)
		mockClient.On("CreateVolumeGroup", c.NewAdminContext(), volumeGroup).Return(&SampleVolumeGroups[0], nil)
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/block/volumeGroups", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.VolumeGroupSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 202)
		assertTestResult(t, &output, &SampleVolumeGroups[0])
	})
	t.Run("Should return 400 if create volume group with bad request", func(t *testing.T) {
		vg := model.VolumeGroupSpec{BaseModel: &model.BaseModel{
			Id: "3769855c-a102-11e7-b772-17b880d2f555",
			CreatedAt: time.Now().Format(constants.TimeFormat),
		},
		Status: "creating",
		AvailabilityZone: "default",
		}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&vg)
		mockClient := new(dbtest.Client)
		mockClient.On("CreateVolumeGroup", c.NewAdminContext(), &vg).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/block/volumeGroups", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 400)
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

	t.Run("Should return 404 if get volume group resource does not exist", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolumeGroup", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f555").Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/volumeGroups/3769855c-a102-11e7-b772-17b880d2f555", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}

func TestUpdateVolumeGroup(t *testing.T) {
	var jsonStr = []byte(`{
		"id": "3769855c-a102-11e7-b772-17b880d2f555",
		"name": "volumeGroup-demo",
  		"description": "volumeGroup test"
	}`)

	t.Run("Should return 202 if everything works well", func(t *testing.T) {
		vg := model.VolumeGroupSpec{BaseModel: &model.BaseModel{
			Id:        "3769855c-a102-11e7-b772-17b880d2f555",
			UpdatedAt: time.Now().Format(constants.TimeFormat),
		},
			Status: "available",
		}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&vg)
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolumeGroup", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f555").Return(&SampleVolumeGroups[1], nil)
		mockClient.On("ListVolumesByGroupId", c.NewAdminContext(), SampleVolumeGroups[1].Id).Return(nil, nil)
		mockClient.On("GetVolume", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
		mockClient.On("UpdateVolumeGroup", c.NewAdminContext(), &vg).Return(&SampleVolumeGroups[1], nil)
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/block/volumeGroups/3769855c-a102-11e7-b772-17b880d2f555", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.VolumeGroupSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 202)
		assertTestResult(t, &output, &SampleVolumeGroups[1])
	})

	t.Run("Should return 400 if update volume fails with bad request", func(t *testing.T) {
		vg := model.VolumeGroupSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&vg)
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolumeGroup", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&vg, nil)
		mockClient.On("ListVolumesByGroupId", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(nil, nil)
		mockClient.On("UpdateVolumeGroup", c.NewAdminContext(), &vg).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/block/volumeGroups/bd5b12a8-a101-11e7-941e-d77981b584d8", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 400)
	})
}

func TestDeleteVolumeGroup(t *testing.T) {
	t.Run("Should return 202 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		var volumesUpdate []*model.VolumeSpec
		mockClient.On("GetVolumeGroup", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f555").Return(&SampleVolumeGroups[0], nil)
		mockClient.On("GetDockByPoolId", c.NewAdminContext(), SampleVolumeGroups[0].PoolId).Return(&SampleDocks[0], nil)
		mockClient.On("ListVolumesByGroupId", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f555").Return(nil, nil)
		mockClient.On("ListSnapshotsByVolumeId", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return( nil, nil)
		mockClient.On("UpdateStatus", c.NewAdminContext(), volumesUpdate, "").Return( nil)
		mockClient.On("UpdateStatus", c.NewAdminContext(), &SampleVolumeGroups[0], "deleting").Return( nil)
		mockClient.On("DeleteVolumeGroup", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f555").Return(nil)
		db.C = mockClient

		r, _ := http.NewRequest("DELETE", "/v1beta/block/volumeGroups/3769855c-a102-11e7-b772-17b880d2f555", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 202)
	})
}
