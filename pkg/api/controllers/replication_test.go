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
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/testutils/collection"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
)

func init() {
	beego.Router("/v1beta/block/replications", NewReplicationPortal(),
		"post:CreateReplication;get:ListReplications")
	beego.Router("/v1beta/block/replications/detail", NewReplicationPortal(),
		"get:ListReplicationsDetail")
	beego.Router("/v1beta/block/replications/:replicationId", NewReplicationPortal(),
		"get:GetReplication;put:UpdateReplication;delete:DeleteReplication")
}

func TestListReplicationsDetail(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		var sampleReplications = []*model.ReplicationSpec{&SampleReplications[0]}
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListReplicationWithFilter", c.NewAdminContext(), m).Return(sampleReplications, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/replications/detail?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output []*model.ReplicationSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, sampleReplications)
	})
}

func TestListReplications(t *testing.T) {
	var expectedJson = []byte(`[
		{
		    "id": "c299a978-4f3e-11e8-8a5c-977218a83359",
			"name": "sample-replication-01"
		}
	]`)
	var expected []*model.ReplicationSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		var sampleReplications = []*model.ReplicationSpec{&SampleReplications[0]}
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListReplicationWithFilter", c.NewAdminContext(), m).Return(sampleReplications, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/replications?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output []*model.ReplicationSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, expected)
	})

	t.Run("Should return 500 if list volume replications with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListReplicationWithFilter", c.NewAdminContext(), m).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/replications?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestGetReplication(t *testing.T) {
	var expectedJson = []byte(`{
			"id": "c299a978-4f3e-11e8-8a5c-977218a83359",
			"primaryVolumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"secondaryVolumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"name": "sample-replication-01",
			"description": "This is a sample replication for testing",
			"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c"
	}`)
	var expected model.ReplicationSpec
	json.Unmarshal(expectedJson, &expected)

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetReplication", c.NewAdminContext(), "c299a978-4f3e-11e8-8a5c-977218a83359").Return(&SampleReplications[0], nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/replications/c299a978-4f3e-11e8-8a5c-977218a83359", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.ReplicationSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &expected)
	})

	t.Run("Should return 404 if get volume replication with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetReplication", c.NewAdminContext(), "c299a978-4f3e-11e8-8a5c-977218a83359").Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/replications/c299a978-4f3e-11e8-8a5c-977218a83359", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}

func TestUpdateReplication(t *testing.T) {
	var jsonStr = []byte(`{
		    "id": "c299a978-4f3e-11e8-8a5c-977218a83359",
			"name":"fake replication",
			"description":"fake replication"
	}`)
	var expectedJson = []byte(`{
			"id": "c299a978-4f3e-11e8-8a5c-977218a83359",
			"primaryVolumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"secondaryVolumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"name": "fake replication",
			"description": "fake replication",
			"poolId": "084bf71e-a102-11e7-88a8-e31fe6d52248",
			"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c"
	}`)
	var expected model.ReplicationSpec
	json.Unmarshal(expectedJson, &expected)

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		replication := model.ReplicationSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&replication)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateReplication", c.NewAdminContext(), replication.Id, &replication).Return(&expected, nil)
		mockClient.On("GetProfile", c.NewAdminContext(), SampleReplications[0].ProfileId).Return(&SampleProfiles[0], nil)
		db.C = mockClient

		r, _ := http.NewRequest("PUT",
			"/v1beta/block/replications/c299a978-4f3e-11e8-8a5c-977218a83359", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.ReplicationSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &expected)
	})

	t.Run("Should return 500 if update volume replication with bad request", func(t *testing.T) {
		replication := model.ReplicationSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&replication)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateReplication", c.NewAdminContext(), replication.Id,
			&replication).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("PUT",
			"/v1beta/block/replications/c299a978-4f3e-11e8-8a5c-977218a83359", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}
