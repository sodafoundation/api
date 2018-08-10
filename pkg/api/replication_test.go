// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
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

var (
	fakeReplication = &model.ReplicationSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			CreatedAt: "2017-10-24T16:21:32",
		},
		Name:              "fake replication",
		Description:       "fake replication",
		AvailabilityZone:  "default",
		PrimaryVolumeId:   "d3a109ff-3e51-4625-9054-32604c79fa90",
		SecondaryVolumeId: "d3a109ff-3e51-4625-9054-32604c79fa92",
		ReplicationMode:   "sync",
		ReplicationPeriod: 0,
		ReplicationStatus: model.ReplicationEnabled,
	}
	fakeReplications = []*model.ReplicationSpec{fakeReplication}
)

func TestListReplicationsDetail(t *testing.T) {

	mockClient := new(dbtest.Client)
	m := map[string][]string{
		"offset":  []string{"0"},
		"limit":   []string{"1"},
		"sortDir": []string{"asc"},
		"sortKey": []string{"name"},
	}
	mockClient.On("ListReplicationWithFilter", c.NewAdminContext(), m).Return(fakeReplications, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/replications/detail?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output []model.ReplicationSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `[{
		    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"createdAt": "2017-10-24T16:21:32",
			"name": "fake replication",
			"description": "fake replication",
			"availabilityZone": "default",
			"PrimaryVolumeId":   "d3a109ff-3e51-4625-9054-32604c79fa90",
			"SecondaryVolumeId": "d3a109ff-3e51-4625-9054-32604c79fa92",
			"ReplicationMode": "sync",
			"ReplicationPeriod": 0,
			"ReplicationStatus": "enabled"
		}]`

	var expected []model.ReplicationSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListReplications(t *testing.T) {

	mockClient := new(dbtest.Client)
	m := map[string][]string{
		"offset":  []string{"0"},
		"limit":   []string{"1"},
		"sortDir": []string{"asc"},
		"sortKey": []string{"name"},
	}
	mockClient.On("ListReplicationWithFilter", c.NewAdminContext(), m).Return(fakeReplications, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/replications?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output []model.ReplicationSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `[{
		    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"name": "fake replication",
			"ReplicationStatus": "enabled"
		}]`

	var expected []model.ReplicationSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListReplicationsWithBadRequest(t *testing.T) {

	mockClient := new(dbtest.Client)
	m := map[string][]string{
		"offset":  []string{"0"},
		"limit":   []string{"1"},
		"sortDir": []string{"asc"},
		"sortKey": []string{"name"},
	}
	mockClient.On("ListReplicationWithFilter", c.NewAdminContext(), m).Return(nil, errors.New("db error"))
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/replications?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}

func TestGetReplication(t *testing.T) {

	mockClient := new(dbtest.Client)
	mockClient.On("GetReplication", c.NewAdminContext(), "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(fakeReplication, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/replications/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.ReplicationSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `{
		    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"createdAt": "2017-10-24T16:21:32",
			"name": "fake replication",
			"description": "fake replication",
			"availabilityZone": "default",
			"PrimaryVolumeId":   "d3a109ff-3e51-4625-9054-32604c79fa90",
			"SecondaryVolumeId": "d3a109ff-3e51-4625-9054-32604c79fa92",
			"ReplicationMode": "sync",
			"ReplicationPeriod": 0,
			"ReplicationStatus": "enabled"
		}`

	var expected model.ReplicationSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestGetReplicationWithBadRequest(t *testing.T) {

	mockClient := new(dbtest.Client)
	mockClient.On("GetReplication", c.NewAdminContext(), "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(nil, errors.New("db error"))
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/replications/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}

func TestUpdateReplication(t *testing.T) {
	var jsonStr = []byte(`{
		    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"name":"fake replication",
			"description":"fake replication"}`)
	r, _ := http.NewRequest("PUT",
		"/v1beta/block/replications/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/json")

	var replication = model.ReplicationSpec{
		BaseModel: &model.BaseModel{},
	}
	json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&replication)
	mockClient := new(dbtest.Client)
	mockClient.On("UpdateReplication", c.NewAdminContext(), replication.Id, &replication).Return(fakeReplication, nil)
	db.C = mockClient
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.ReplicationSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `{
		    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"createdAt": "2017-10-24T16:21:32",
			"name": "fake replication",
			"description": "fake replication",
			"availabilityZone": "default",
			"PrimaryVolumeId":   "d3a109ff-3e51-4625-9054-32604c79fa90",
			"SecondaryVolumeId": "d3a109ff-3e51-4625-9054-32604c79fa92",
			"ReplicationMode": "sync",
			"ReplicationPeriod": 0,
			"ReplicationStatus": "enabled"
		}`

	var expected model.ReplicationSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestUpdateReplicationWithBadRequest(t *testing.T) {
	var jsonStr = []byte(``)
	r, _ := http.NewRequest("PUT",
		"/v1beta/block/replications/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/json")
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}

	jsonStr = []byte(`{
		    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"name":"fake replication",
			"description":"fake replication"}`)
	r, _ = http.NewRequest("PUT",
		"/v1beta/block/replications/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", bytes.NewBuffer(jsonStr))
	w = httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/json")

	var replication = model.ReplicationSpec{
		BaseModel: &model.BaseModel{},
	}
	json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&replication)

	mockClient := new(dbtest.Client)
	mockClient.On("UpdateReplication", c.NewAdminContext(), replication.Id,
		&replication).Return(nil, errors.New("db error"))
	db.C = mockClient
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}
