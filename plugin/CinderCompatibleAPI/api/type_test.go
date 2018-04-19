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
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/astaxie/beego"
	c "github.com/opensds/opensds/client"
	"github.com/opensds/opensds/plugin/CinderCompatibleAPI/CinderModel"
)

func init() {
	beego.Router("/v3/types", &TypePortal{},
		"post:CreateType;get:ListType")
	beego.Router("/v3/types/:volumeTypeId", &TypePortal{},
		"get:GetType;delete:DeleteType;put:UpdateType")
	beego.Router("/v3/types/:volumeTypeId/extra_specs", &TypePortal{},
		"post:AddExtraProperty;get:ListExtraProperties")
	beego.Router("/v3/types/:volumeTypeId/extra_specs/:key", &TypePortal{},
		"get:ShowExtraProperty;put:UpdateExtraProperty;delete:DeleteExtraProperty")

	if false == IsFakeClient {
		client = NewFakeClient(&c.Config{Endpoint: TestEp})
	}
}

////////////////////////////////////////////////////////////////////////////////
//                            Tests for volume types                               //
////////////////////////////////////////////////////////////////////////////////
func TestGetType(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v3/types/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output CinderModel.ShowTypeRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJSON := `{
    	"volume_type": {
		"id": "1106b972-66ef-11e7-b172-db03f3689c9c",
        "name": "default",
        "description": "default policy",
		"extra_specs": []
    	}
	}`

	var expected CinderModel.ShowTypeRespSpec
	json.Unmarshal([]byte(expectedJSON), &expected)
	expected.VolumeType.IsPublic = true

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestGetDefaultType(t *testing.T) {
	DefaultTypeName = "default"

	r, _ := http.NewRequest("GET", "/v3/types/default", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output CinderModel.ShowTypeRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if DefaultTypeName != output.VolumeType.Name {
		t.Errorf("Expected %v, actual %v", DefaultTypeName, output.VolumeType.Description)
	}
}

func TestListType(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v3/types", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output CinderModel.ListTypeRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJSON := `{
	"volume_types":	
	[{
		"id": "1106b972-66ef-11e7-b172-db03f3689c9c",
        "name": "default",
        "description": "default policy",
		"os-volume-type-access:is_public": true,
		"is_public": true
    	}
		,
		{
		"id": "2f9c0a04-66ef-11e7-ade2-43158893e017",
        "name": "silver",
        "description": "silver policy",
		"os-volume-type-access:is_public": true,
		"is_public": true,
		"extra_specs": {
            "dataStorage": {
					"provisioningPolicy": "Thin",
					"isSpaceEfficient":   true
				},
				"ioConnectivity": {
					"accessProtocol": "rbd",
					"maxIOPS":        5000000,
					"maxBWS":         500
				}
        }
    	}
		]
		}`

	var expected CinderModel.ListTypeRespSpec
	json.Unmarshal([]byte(expectedJSON), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestCreateType(t *testing.T) {
	RequestBodyStr := `{
    	"volume_type": {
        "name": "default",
		"os-volume-type-access:is_public": true,
        "description": "default policy"
    	}
	}`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("POST", "/v3/types", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output CinderModel.CreateTypeRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	var expected CinderModel.CreateTypeRespSpec
	json.Unmarshal([]byte(RequestBodyStr), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}
	fmt.Println(expected)
	fmt.Println(output)
	expected.VolumeType.ID = "1106b972-66ef-11e7-b172-db03f3689c9c"
	expected.VolumeType.IsPublic = true

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestDeleteType(t *testing.T) {
	r, _ := http.NewRequest("DELETE", "/v3/types/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 202 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}
}

func TestUpdateType(t *testing.T) {
	RequestBodyStr := `{
    	"volume_type": {
        "name": "default",
        "description": "default policy",
		"is_public": true
    	}
	}`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("PUT", "/v3/types/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output CinderModel.UpdateTypeRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	var expected CinderModel.UpdateTypeRespSpec
	json.Unmarshal([]byte(RequestBodyStr), &expected)
	expected.VolumeType.IsPublic = true

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}
	expected.VolumeType.ID = "1106b972-66ef-11e7-b172-db03f3689c9c"
	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestAddExtraProperty(t *testing.T) {
	RequestBodyStr := `{
		"extra_specs":{
			"dataStorage": {
					"provisioningPolicy": "Thin",
					"isSpaceEfficient":   true
				},
				"ioConnectivity": {
					"accessProtocol": "rbd",
					"maxIOPS":        5000000,
					"maxBWS":         500
				}
		}
	}`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("POST", "/v3/types/f4a5e666-c669-4c64-a2a1-8f9ecd560c78/extra_specs", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output CinderModel.ExtraSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	var expected CinderModel.ExtraSpec
	json.Unmarshal([]byte(RequestBodyStr), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListExtraProperties(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v3/types/f4a5e666-c669-4c64-a2a1-8f9ecd560c78/extra_specs", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output CinderModel.ExtraSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	var expected CinderModel.ExtraSpec
	expectedStr := `{
		"extra_specs":{
			"dataStorage": {
					"provisioningPolicy": "Thin",
					"isSpaceEfficient":   true
				},
				"ioConnectivity": {
					"accessProtocol": "rbd",
					"maxIOPS":        5000000,
					"maxBWS":         500
				}
		}
	}`
	json.Unmarshal([]byte(expectedStr), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestShowExtraPropertie(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v3/types/f4a5e666-c669-4c64-a2a1-8f9ecd560c78/extra_specs/dataStorage", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output CinderModel.ExtraSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	var expected CinderModel.ExtraSpec
	expectedStr := `{
		"dataStorage": {
					"provisioningPolicy": "Thin",
					"isSpaceEfficient":   true
				}
	}`
	json.Unmarshal([]byte(expectedStr), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestShowExtraPropertieWithBadRequest(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v3/types/f4a5e666-c669-4c64-a2a1-8f9ecd560c78/extra_specs/disk", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output CinderModel.ExtraSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if nil != output["disk"] {
		t.Errorf("Expected %v, actual %v", nil, output["disk"])
	}
}

func TestUpdateExtraPropertie(t *testing.T) {
	RequestBodyStr := `{
		"dataStorage": {
					"provisioningPolicy": "Thin",
					"isSpaceEfficient":   true
				}
	}`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("PUT", "/v3/types/f4a5e666-c669-4c64-a2a1-8f9ecd560c78/extra_specs/dataStorage", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output CinderModel.ExtraSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedStr := `{
		"dataStorage": {
					"provisioningPolicy": "Thin",
					"isSpaceEfficient":   true
				}
	}`

	var expected CinderModel.ExtraSpec
	json.Unmarshal([]byte(expectedStr), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestDeleteExtraPropertie(t *testing.T) {
	r, _ := http.NewRequest("DELETE", "/v3/types/f4a5e666-c669-4c64-a2a1-8f9ecd560c78/extra_specs/diskType", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != StatusAccepted {
		t.Errorf("Expected %v, actual %v", StatusAccepted, w.Code)
	}
}
