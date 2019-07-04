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

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	. "github.com/opensds/opensds/testutils/collection"
	ctrtest "github.com/opensds/opensds/testutils/controller/testing"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
)

////////////////////////////////////////////////////////////////////////////////
//                      Prepare for mock server                               //
////////////////////////////////////////////////////////////////////////////////
func init() {
	beego.Router("/v1beta/file/shares", NewFakeFileSharePortal(),
		"post:CreateFileShare;get:ListFileShares")
	beego.Router("/v1beta/file/shares/:fileshareId", NewFakeFileSharePortal(),
		"get:GetFileShare;put:UpdateFileShare;delete:DeleteFileShare")

	beego.Router("/v1beta/file/snapshots", &FileShareSnapshotPortal{},
		"post:CreateFileShareSnapshot;get:ListFileShareSnapshots")
	beego.Router("/v1beta/file/snapshots/:snapshotId", &FileShareSnapshotPortal{},
		"get:GetFileShareSnapshot;put:UpdateFileShareSnapshot;delete:DeleteFileShareSnapshot")

}

func NewFakeFileSharePortal() *FileSharePortal {
	mockClient := new(ctrtest.Client)

	mockClient.On("Connect", "localhost:50049").Return(nil)
	mockClient.On("Close").Return(nil)
	mockClient.On("CreateFileShare", ctx.Background(), &pb.CreateFileShareOpts{
		Context: c.NewAdminContext().ToJson(),
	}).Return(&pb.GenericResponse{}, nil)
	mockClient.On("DeleteFileShare", ctx.Background(), &pb.DeleteFileShareOpts{
		Context: c.NewAdminContext().ToJson(),
	}).Return(&pb.GenericResponse{}, nil)

	return &FileSharePortal{
		CtrClient: mockClient,
	}
}

////////////////////////////////////////////////////////////////////////////////
//                            Tests for FileShare                             //
////////////////////////////////////////////////////////////////////////////////

var (
	fakeFileShare = &model.FileShareSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			CreatedAt: "2017-10-24T16:21:32",
		},
		Name:             "fake FileShare",
		Description:      "fake FileShare",
		Size:             99,
		AvailabilityZone: "unknown",
		Status:           "available",
		PoolId:           "831fa5fb-17cf-4410-bec6-1f4b06208eef",
		ProfileId:        "d3a109ff-3e51-4625-9054-32604c79fa90",
	}
	fakeFileShares = []*model.FileShareSpec{fakeFileShare}
)

func TestListFileShares(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		var sampleFileShares = []*model.FileShareSpec{&SampleFileShares[0], &SampleFileShares[1]}
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListFileSharesWithFilter", c.NewAdminContext(), m).Return(sampleFileShares, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/file/shares?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		var output []*model.FileShareSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, sampleFileShares)
	})

	t.Run("Should return 500 if list file share with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListFileSharesWithFilter", c.NewAdminContext(), m).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/file/shares?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestGetFileShare(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetFileShare", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleFileShares[0], nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/file/shares/bd5b12a8-a101-11e7-941e-d77981b584d8", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		var output model.FileShareSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, &output, &SampleFileShares[0])
	})

	t.Run("Should return 404 if get file share replication with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetFileShare", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/file/shares/bd5b12a8-a101-11e7-941e-d77981b584d8", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}

func TestUpdateFileShare(t *testing.T) {
	var jsonStr = []byte(`{
		"id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
		"name":"fake FileShare",
		"description":"fake Fileshare"
	}`)
	var expectedJson = []byte(`{
		"id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
		"name": "fake FileShare",
		"description": "fake FileShare",
		"size": 1,
		"status": "available",
		"poolId": "084bf71e-a102-11e7-88a8-e31fe6d52248",
		"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c"
	}`)
	var expected model.FileShareSpec
	json.Unmarshal(expectedJson, &expected)

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		fileshare := model.FileShareSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&fileshare)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateFileShare", c.NewAdminContext(), &fileshare).Return(&expected, nil)
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/file/shares/bd5b12a8-a101-11e7-941e-d77981b584d8", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.FileShareSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &expected)
	})

	t.Run("Should return 500 if update file share with bad request", func(t *testing.T) {
		fileshare := model.FileShareSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&fileshare)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateFileShare", c.NewAdminContext(), &fileshare).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/file/shares/bd5b12a8-a101-11e7-941e-d77981b584d8", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

////////////////////////////////////////////////////////////////////////////////
//                      Tests for fileshare snapshot                          //
////////////////////////////////////////////////////////////////////////////////

func TestListFileShareSnapshots(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		var sampleSnapshots = []*model.FileShareSnapshotSpec{&SampleFileShareSnapshots[0], &SampleFileShareSnapshots[1]}
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListFileShareSnapshotsWithFilter", c.NewAdminContext(), m).Return(sampleSnapshots, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/file/snapshots?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		var output []*model.FileShareSnapshotSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, sampleSnapshots)
	})

	t.Run("Should return 500 if list fileshare snapshots with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListFileShareSnapshotsWithFilter", c.NewAdminContext(), m).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/file/snapshots?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestGetFileShareSnapshot(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetFileShareSnapshot", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f537").Return(&SampleFileShareSnapshots[0], nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/file/snapshots/3769855c-a102-11e7-b772-17b880d2f537", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.FileShareSnapshotSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SampleFileShareSnapshots[0])
	})

	t.Run("Should return 404 if get fileshare group with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetFileShareSnapshot", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f537").Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/file/snapshots/3769855c-a102-11e7-b772-17b880d2f537", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}

func TestUpdateFileShareSnapshot(t *testing.T) {
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
		"fileshareId": "bd5b12a8-a101-11e7-941e-d77981b584d8"
	}`)
	var expected model.FileShareSnapshotSpec
	json.Unmarshal(expectedJson, &expected)

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		snapshot := model.FileShareSnapshotSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&snapshot)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateFileShareSnapshot", c.NewAdminContext(), snapshot.Id, &snapshot).
			Return(&expected, nil)
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/file/snapshots/3769855c-a102-11e7-b772-17b880d2f537", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.FileShareSnapshotSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &expected)
	})

	t.Run("Should return 500 if update fileshare snapshot with bad request", func(t *testing.T) {
		snapshot := model.FileShareSnapshotSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&snapshot)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateFileShareSnapshot", c.NewAdminContext(), snapshot.Id, &snapshot).
			Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/file/snapshots/3769855c-a102-11e7-b772-17b880d2f537", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}
