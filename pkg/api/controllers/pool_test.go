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
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/testutils/collection"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
)

func init() {
	var poolPortal PoolPortal
	beego.Router("/v1beta/pools", &poolPortal, "get:ListPools")
	beego.Router("/v1beta/availabilityZones", &poolPortal, "get:ListAvailabilityZones")
	beego.Router("/v1beta/pools/:poolId", &poolPortal, "get:GetPool")
}

func TestListAvailabilityZones(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("ListAvailabilityZones", c.NewAdminContext()).Return(SampleAvailabilityZones, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/availabilityZones", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output []string
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, SampleAvailabilityZones)
	})
}

func TestListPools(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		var samplePools = []*model.StoragePoolSpec{&SamplePools[0], &SamplePools[1]}
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListPoolsWithFilter", c.NewAdminContext(), m).Return(samplePools, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/pools?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output []*model.StoragePoolSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, samplePools)
	})

	t.Run("Should return 500 if list pools with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListPoolsWithFilter", c.NewAdminContext(), m).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/pools?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestGetPool(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetPool", c.NewAdminContext(), "f4486139-78d5-462d-a7b9-fdaf6c797e1b").Return(&SamplePools[0], nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/pools/f4486139-78d5-462d-a7b9-fdaf6c797e1b", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.StoragePoolSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SamplePools[0])
	})

	t.Run("Should return 404 if get docks with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetPool", c.NewAdminContext(), "f4486139-78d5-462d-a7b9-fdaf6c797e1b").Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET",
			"/v1beta/pools/f4486139-78d5-462d-a7b9-fdaf6c797e1b", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}
