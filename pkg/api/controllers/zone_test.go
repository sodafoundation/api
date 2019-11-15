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
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/testutils/collection"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
)

func init() {
	var zonePortal AvailabilityZonePortal
	beego.Router("/v1beta/availabilityZones", &zonePortal, "post:CreateAvailabilityZone;get:ListAvailabilityZones")
	beego.Router("/v1beta/availabilityZones/:zoneId", &zonePortal, "get:GetAvailabilityZone;put:UpdateAvailabilityZone;delete:DeleteAvailabilityZone")
}

////////////////////////////////////////////////////////////////////////////////
//                            Tests for zone                                  //
////////////////////////////////////////////////////////////////////////////////

func TestCreateAvailabilityZone(t *testing.T) {
	var fakeBody = `{
		"name": "default",
		"description": "default zone"
	}`

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("CreateAvailabilityZone", c.NewAdminContext(), &model.AvailabilityZoneSpec{
			BaseModel:   &model.BaseModel{},
			Name:        "default",
			Description: "default zone",
		}).Return(&SampleAvailabilityZones[1], nil)
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/availabilityZones", strings.NewReader(fakeBody))
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.AvailabilityZoneSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SampleAvailabilityZones[1])
	})
}

func TestUpdateAvailabilityZone(t *testing.T) {
	var jsonStr = []byte(`{
		"id": "2f9c0a04-66ef-11e7-ade2-43158893e017",
		"name": "test",
		"description": "test zone"
	}`)
	var expectedJson = []byte(`{
		"id": "2f9c0a04-66ef-11e7-ade2-43158893e017",
		"name": "test",
		"description": "test zone"
	}`)
	var expected model.AvailabilityZoneSpec
	json.Unmarshal(expectedJson, &expected)

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		zone := model.AvailabilityZoneSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&zone)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateAvailabilityZone", c.NewAdminContext(), zone.Id, &zone).
			Return(&expected, nil)
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/availabilityZones/2f9c0a04-66ef-11e7-ade2-43158893e017", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.AvailabilityZoneSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &expected)
	})

	t.Run("Should return 500 if update zone with bad request", func(t *testing.T) {
		zone := model.AvailabilityZoneSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&zone)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateAvailabilityZone", c.NewAdminContext(), zone.Id, &zone).
			Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/availabilityZones/2f9c0a04-66ef-11e7-ade2-43158893e017", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestListAvailabilityZone(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		var sampleZones = []*model.AvailabilityZoneSpec{&SampleAvailabilityZones[1]}
		mockClient := new(dbtest.Client)
		mockClient.On("ListAvailabilityZones", c.NewAdminContext()).Return(
			sampleZones, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/availabilityZones?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output []*model.AvailabilityZoneSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, sampleZones)
	})

	t.Run("Should return 500 if list zones with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("ListAvailabilityZones", c.NewAdminContext()).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/availabilityZones?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestGetAvailabilityZone(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetAvailabilityZone", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").
			Return(&SampleAvailabilityZones[1], nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/availabilityZones/2f9c0a04-66ef-11e7-ade2-43158893e017", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.AvailabilityZoneSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SampleAvailabilityZones[1])
	})

	t.Run("Should return 404 if get zone with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetAvailabilityZone", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(
			nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET",
			"/v1beta/availabilityZones/2f9c0a04-66ef-11e7-ade2-43158893e017", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}

func TestDeleteAvailabilityZone(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetAvailabilityZone", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(
			&SampleAvailabilityZones[1], nil)
		mockClient.On("DeleteAvailabilityZone", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(nil)
		db.C = mockClient

		r, _ := http.NewRequest("DELETE",
			"/v1beta/availabilityZones/2f9c0a04-66ef-11e7-ade2-43158893e017", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 200)
	})

	t.Run("Should return 404 if delete zone with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetAvailabilityZone", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(
			nil, errors.New("Invalid resource uuid"))
		db.C = mockClient

		r, _ := http.NewRequest("DELETE",
			"/v1beta/availabilityZones/2f9c0a04-66ef-11e7-ade2-43158893e017", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}