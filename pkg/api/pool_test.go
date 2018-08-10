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

package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/astaxie/beego"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
)

func init() {
	var poolPortal PoolPortal
	beego.Router("/v1beta/pools", &poolPortal, "get:ListPools")
	beego.Router("/v1beta/availabilityZones", &poolPortal, "get:ListAvailabilityZones")
	beego.Router("/v1beta/pools/:poolId", &poolPortal, "get:GetPool")
}

var (
	fakePool = &model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4486139-78d5-462d-a7b9-fdaf6c797e1b",
			CreatedAt: "2017-10-24T15:04:05",
		},
		Name:             "fakePool",
		Description:      "fake pool for testing",
		Status:           "available",
		AvailabilityZone: "unknown",
		TotalCapacity:    99999,
		FreeCapacity:     6999,
		DockId:           "ccac4f33-e603-425a-8813-371bbe10566e",
		Extras: model.StoragePoolExtraSpec{
			DataStorage: model.DataStorageLoS{
				ProvisioningPolicy: "Thin",
				IsSpaceEfficient:   true,
			},
			IOConnectivity: model.IOConnectivityLoS{
				AccessProtocol: "rbd",
				MaxIOPS:        1000,
			},
			Advanced: map[string]interface{}{
				"diskType":   "SSD",
				"throughput": float64(1000),
			},
		},
	}
	fakePools = []*model.StoragePoolSpec{fakePool}
)

func TestListAvailabilityZones(t *testing.T) {
	mockClient := new(dbtest.Client)
	mockClient.On("ListAvailabilityZones", c.NewAdminContext()).Return(fakePools, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/availabilityZones", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	expectedZones := "unknow"
	t.Log(w)
	if !strings.Contains(string(w.Body.Bytes()), expectedZones) {
		t.Errorf("Expected %v, actual %v", expectedZones, w.Body.Bytes())
	}
}

func TestListPools(t *testing.T) {

	mockClient := new(dbtest.Client)
	m := map[string][]string{
		"offset":  []string{"0"},
		"limit":   []string{"1"},
		"sortDir": []string{"asc"},
		"sortKey": []string{"name"},
	}
	mockClient.On("ListPoolsWithFilter", c.NewAdminContext(), m).Return(fakePools, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/pools?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output []model.StoragePoolSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `[
		{
			"id": "f4486139-78d5-462d-a7b9-fdaf6c797e1b",
			"name": "fakePool",
			"description": "fake pool for testing",
			"createdAt": "2017-10-24T15:04:05",
			"updatedAt": "",
			"status": "available",
			"availabilityZone": "unknown",
			"totalCapacity": 99999,
			"freeCapacity": 6999,
			"dockId": "ccac4f33-e603-425a-8813-371bbe10566e",
			"extras": {
				"dataStorage": {
					"provisioningPolicy": "Thin",
					"isSpaceEfficient":   true
				},
				"ioConnectivity": {
					"accessProtocol": "rbd",
					"maxIOPS":        1000
				},
				"advanced": {
					"diskType":   "SSD",
					"throughput": 1000
				}
			}	
		}		
	]`

	var expected []model.StoragePoolSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListPoolsWithBadRequest(t *testing.T) {

	mockClient := new(dbtest.Client)
	m := map[string][]string{
		"offset":  []string{"0"},
		"limit":   []string{"1"},
		"sortDir": []string{"asc"},
		"sortKey": []string{"name"},
	}
	mockClient.On("ListPoolsWithFilter", c.NewAdminContext(), m).Return(nil, errors.New("db error"))
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/pools?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}

func TestGetPool(t *testing.T) {

	mockClient := new(dbtest.Client)
	mockClient.On("GetPool", c.NewAdminContext(), "f4486139-78d5-462d-a7b9-fdaf6c797e1b").Return(fakePool, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/pools/f4486139-78d5-462d-a7b9-fdaf6c797e1b", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.StoragePoolSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `
		{
			"id": "f4486139-78d5-462d-a7b9-fdaf6c797e1b",
			"name": "fakePool",
			"description": "fake pool for testing",
			"createdAt": "2017-10-24T15:04:05",
			"updatedAt": "",
			"status": "available",
			"availabilityZone": "unknown",
			"totalCapacity": 99999,
			"freeCapacity": 6999,
			"dockId": "ccac4f33-e603-425a-8813-371bbe10566e",
			"extras": {
				"dataStorage": {
					"provisioningPolicy": "Thin",
					"isSpaceEfficient":   true
				},
				"ioConnectivity": {
					"accessProtocol": "rbd",
					"maxIOPS":        1000
				},
				"advanced": {
					"diskType":   "SSD",
					"throughput": 1000
				}
			}	
		}`

	var expected model.StoragePoolSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestGetPoolWithBadRequest(t *testing.T) {

	mockClient := new(dbtest.Client)
	mockClient.On("GetPool", c.NewAdminContext(), "f4486139-78d5-462d-a7b9-fdaf6c797e1b").Return(
		nil, errors.New("db error"))
	db.C = mockClient

	r, _ := http.NewRequest("GET",
		"/v1beta/pools/f4486139-78d5-462d-a7b9-fdaf6c797e1b", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}
