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

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	c "github.com/sodafoundation/api/pkg/context"
	"github.com/sodafoundation/api/pkg/db"
	"github.com/sodafoundation/api/pkg/model"
	pb "github.com/sodafoundation/api/pkg/model/proto"
	"github.com/sodafoundation/api/pkg/utils/constants"
	. "github.com/sodafoundation/api/testutils/collection"
	ctrtest "github.com/sodafoundation/api/testutils/controller/testing"
	dbtest "github.com/sodafoundation/api/testutils/db/testing"
)

////////////////////////////////////////////////////////////////////////////////
//                      Prepare for mock server                               //
////////////////////////////////////////////////////////////////////////////////

func init() {
	beego.Router("/v1beta/block/volumes", NewFakeVolumePortal(),
		"post:CreateVolume;get:ListVolumes")
	beego.Router("/v1beta/block/volumes/:volumeId", NewFakeVolumePortal(),
		"get:GetVolume;put:UpdateVolume;delete:DeleteVolume")
	beego.Router("/v1beta/block/volumes/:volumeId/resize", NewFakeVolumePortal(),
		"post:ExtendVolume")

	beego.Router("/v1beta/block/snapshots", &VolumeSnapshotPortal{},
		"post:CreateVolumeSnapshot;get:ListVolumeSnapshots")
	beego.Router("/v1beta/block/snapshots/:snapshotId", &VolumeSnapshotPortal{},
		"get:GetVolumeSnapshot;put:UpdateVolumeSnapshot;delete:DeleteVolumeSnapshot")
}

func NewFakeVolumePortal() *VolumePortal {
	mockClient := new(ctrtest.Client)

	mockClient.On("Connect", "localhost:50049").Return(nil)
	mockClient.On("Close").Return(nil)
	mockClient.On("CreateVolume", ctx.Background(), &pb.CreateVolumeOpts{
		Id:          "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:        "sample-volume",
		Size:        1,
		Description: "This is a sample volume for testing",
		//SnapshotId:            "",
		AvailabilityZone: "default",
		ProfileId:        "1106b972-66ef-11e7-b172-db03f3689c9c",
		PoolId:           "084bf71e-a102-11e7-88a8-e31fe6d52248",
		Context:          c.NewAdminContext().ToJson(),
		Profile:          SampleProfiles[0].ToJson(),
	}).Return(&pb.GenericResponse{}, nil)
	mockClient.On("ExtendVolume", ctx.Background(), &pb.ExtendVolumeOpts{
		Id:      "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Size:    int64(20),
		Context: c.NewAdminContext().ToJson(),
		Profile: SampleProfiles[0].ToJson(),
	}).Return(&pb.GenericResponse{}, nil)
	mockClient.On("DeleteVolume", ctx.Background(), &pb.DeleteVolumeOpts{
		Id:        "bd5b12a8-a101-11e7-941e-d77981b584d8",
		ProfileId: "1106b972-66ef-11e7-b172-db03f3689c9c",
		PoolId:    "084bf71e-a102-11e7-88a8-e31fe6d52248",
		Context:   c.NewAdminContext().ToJson(),
		Profile:   SampleProfiles[0].ToJson(),
	}).Return(&pb.GenericResponse{}, nil)

	return &VolumePortal{
		CtrClient: mockClient,
	}
}

////////////////////////////////////////////////////////////////////////////////
//                            Tests for volume                                //
////////////////////////////////////////////////////////////////////////////////

var (
	fakeVolume = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			CreatedAt: "2017-10-24T16:21:32",
		},
		Name:             "fake Vol",
		Description:      "fake Vol",
		Size:             99,
		AvailabilityZone: "unknown",
		Status:           "available",
		PoolId:           "831fa5fb-17cf-4410-bec6-1f4b06208eef",
		ProfileId:        "d3a109ff-3e51-4625-9054-32604c79fa90",
	}
	fakeVolumes = []*model.VolumeSpec{fakeVolume}
)

func TestListVolumes(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		var sampleVolumes = []*model.VolumeSpec{&SampleVolumes[0], &SampleVolumes[1]}
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListVolumesWithFilter", c.NewAdminContext(), m).Return(sampleVolumes, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/volumes?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		var output []*model.VolumeSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, sampleVolumes)
	})

	t.Run("Should return 500 if list volume with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListVolumesWithFilter", c.NewAdminContext(), m).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/volumes?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestGetVolume(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolume", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		var output model.VolumeSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, &output, &SampleVolumes[0])
	})

	t.Run("Should return 404 if get volume replication with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolume", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}

func TestUpdateVolume(t *testing.T) {
	var jsonStr = []byte(`{
		"id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
		"name":"fake Vol",
		"description":"fake Vol"
	}`)
	var expectedJson = []byte(`{
		"id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
		"name": "fake Vol",
		"description": "fake Vol",
		"size": 1,
		"status": "available",
		"poolId": "084bf71e-a102-11e7-88a8-e31fe6d52248",
		"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c"
	}`)
	var expected model.VolumeSpec
	json.Unmarshal(expectedJson, &expected)

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		volume := model.VolumeSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&volume)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateVolume", c.NewAdminContext(), &volume).Return(&expected, nil)
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/block/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.VolumeSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &expected)
	})

	t.Run("Should return 500 if update volume with bad request", func(t *testing.T) {
		volume := model.VolumeSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&volume)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateVolume", c.NewAdminContext(), &volume).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/block/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestExtendVolume(t *testing.T) {
	var jsonStr = []byte(`{
		"newSize":20
	}`)
	var expectedJson = []byte(`{
		"id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
		"name": "sample-volume",
		"description": "This is a sample volume for testing",
		"size": 1,
        "availabilityZone": "default",
		"status": "extending",
		"poolId": "084bf71e-a102-11e7-88a8-e31fe6d52248",
		"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c"
	}`)
	var expected model.VolumeSpec
	json.Unmarshal(expectedJson, &expected)

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolume", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
		mockClient.On("ExtendVolume", c.NewAdminContext(), &expected).Return(&expected, nil)
		mockClient.On("GetProfile", c.NewAdminContext(), SampleReplications[0].ProfileId).Return(&SampleProfiles[0], nil)
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/block/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8/resize", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.VolumeSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 202)
		assertTestResult(t, &output, &expected)
	})

	t.Run("Should return 400 if extend volume with bad request", func(t *testing.T) {
		jsonStr = []byte(`{
			"newSize": 1
		}`)
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolume", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
		mockClient.On("ExtendVolume", c.NewAdminContext(), &expected).Return(&expected, nil)
		mockClient.On("GetProfile", c.NewAdminContext(), SampleReplications[0].ProfileId).Return(&SampleProfiles[0], nil)
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/block/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8/resize", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 400)
	})
}

func TestCreateVolume(t *testing.T) {
	var jsonStr = []byte(`{
		"id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
		"name":"fake Vol",
		"size": 1,
		"description":"fake Vol"
	}`)

	t.Run("Should return 202 if everything works well", func(t *testing.T) {
		volume := model.VolumeSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&volume)
		volume.CreatedAt = time.Now().Format(constants.TimeFormat)
		volume.AvailabilityZone = "default"
		volume.Status = "creating"
		volume.ProfileId = "1106b972-66ef-11e7-b172-db03f3689c9c"
		mockClient := new(dbtest.Client)
		mockClient.On("GetDefaultProfile", c.NewAdminContext()).Return(&SampleProfiles[0], nil)
		mockClient.On("CreateVolume", c.NewAdminContext(), &volume).Return(&SampleVolumes[0], nil)
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/block/volumes", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.VolumeSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 202)
		assertTestResult(t, &output, &SampleVolumes[0])
	})

	t.Run("Should return 400 - get default profile failed: db error", func(t *testing.T) {
		volume := model.VolumeSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&volume)
		volume.CreatedAt = time.Now().Format(constants.TimeFormat)
		volume.AvailabilityZone = "default"
		volume.Status = "creating"
		mockClient := new(dbtest.Client)
		mockClient.On("GetDefaultProfile", c.NewAdminContext()).Return(nil, errors.New("get default profile failed:"))
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/block/volumes", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 400)
	})

	t.Run("Should return 400 - get profile failed: db error", func(t *testing.T) {
		var testjsonStr = []byte(`{
			"id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"name":"fake Vol",
			"size": 1,
			"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c",
			"description":"fake Vol"
		}`)
		volume := model.VolumeSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&volume)
		volume.CreatedAt = time.Now().Format(constants.TimeFormat)
		volume.AvailabilityZone = "default"
		volume.Status = "creating"
		volume.ProfileId = "1106b972-66ef-11e7-b172-db03f3689c9c"
		mockClient := new(dbtest.Client)
		mockClient.On("GetProfile", c.NewAdminContext(), volume.ProfileId).Return(nil, errors.New("get profile failed:"))
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/block/volumes", bytes.NewBuffer(testjsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 400)
	})

	t.Run("Should return 400 - get default profile failed: db error", func(t *testing.T) {
		volume := model.VolumeSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&volume)
		volume.CreatedAt = time.Now().Format(constants.TimeFormat)
		volume.AvailabilityZone = "default"
		volume.Status = "creating"
		volume.ProfileId = "1106b972-66ef-11e7-b172-db03f3689c9c"
		mockClient := new(dbtest.Client)
		mockClient.On("GetDefaultProfile", c.NewAdminContext()).Return(&SampleProfiles[0], nil)
		mockClient.On("CreateVolume", c.NewAdminContext(), &volume).Return(nil, errors.New("create volume failed:"))
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/block/volumes", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 400)
	})
}

func TestDeleteVolume(t *testing.T) {

	t.Run("Should return 202 if everything works well", func(t *testing.T) {
		SampleVolumes[0].Status = "available"
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolume", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
		mockClient.On("GetProfile", c.NewAdminContext(), "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&SampleProfiles[0], nil)
		mockClient.On("ListSnapshotsByVolumeId", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(nil, nil)
		mockClient.On("ListAttachmentsByVolumeId", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(nil, nil)
		mockClient.On("UpdateVolume", c.NewAdminContext(), &SampleVolumes[0]).Return(nil, nil)
		mockClient.On("DeleteVolume", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(nil)
		db.C = mockClient

		r, _ := http.NewRequest("DELETE", "/v1beta/block/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		assertTestResult(t, w.Code, 202)
	})

	t.Run("Should return 404 - Get volume id failed: db error ", func(t *testing.T) {
		SampleVolumes[0].Status = "available"
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolume", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(nil, errors.New("specified volume(bd5b12a8-a101-11e7-941e-d77981b584d8) can't find"))
		db.C = mockClient

		r, _ := http.NewRequest("DELETE", "/v1beta/block/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		assertTestResult(t, w.Code, 404)
	})
}

////////////////////////////////////////////////////////////////////////////////
//                         Tests for volume snapshot                          //
////////////////////////////////////////////////////////////////////////////////

func TestListVolumeSnapshots(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		var sampleSnapshots = []*model.VolumeSnapshotSpec{&SampleSnapshots[0], &SampleSnapshots[1]}
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListVolumeSnapshotsWithFilter", c.NewAdminContext(), m).Return(sampleSnapshots, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/snapshots?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		var output []*model.VolumeSnapshotSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, sampleSnapshots)
	})

	t.Run("Should return 500 if list volume snapshots with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListVolumeSnapshotsWithFilter", c.NewAdminContext(), m).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/snapshots?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestGetVolumeSnapshot(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolumeSnapshot", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f537").Return(&SampleSnapshots[0], nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/snapshots/3769855c-a102-11e7-b772-17b880d2f537", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.VolumeSnapshotSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SampleSnapshots[0])
	})

	t.Run("Should return 404 if get volume group with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolumeSnapshot", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f537").Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/snapshots/3769855c-a102-11e7-b772-17b880d2f537", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}

func TestUpdateVolumeSnapshot(t *testing.T) {
	var jsonStr = []byte(`{
		"id": "3769855c-a102-11e7-b772-17b880d2f537",
		"name":"fake snapshot",
		"description":"fake snapshot"
	}`)
	var expectedJson = []byte(`{
		"id": "3769855c-a102-11e7-b772-17b880d2f537",
		"name": "fake snapshot",
		"description": "fake snapshot",
		"size": 1,
		"status": "available",
		"volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8"
	}`)
	var expected model.VolumeSnapshotSpec
	json.Unmarshal(expectedJson, &expected)

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		snapshot := model.VolumeSnapshotSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&snapshot)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateVolumeSnapshot", c.NewAdminContext(), snapshot.Id, &snapshot).
			Return(&expected, nil)
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/block/snapshots/3769855c-a102-11e7-b772-17b880d2f537", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.VolumeSnapshotSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &expected)
	})

	t.Run("Should return 500 if update volume snapshot with bad request", func(t *testing.T) {
		snapshot := model.VolumeSnapshotSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&snapshot)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateVolumeSnapshot", c.NewAdminContext(), snapshot.Id, &snapshot).
			Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/block/snapshots/3769855c-a102-11e7-b772-17b880d2f537", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}
