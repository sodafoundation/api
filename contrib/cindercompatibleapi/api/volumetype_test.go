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
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/astaxie/beego"
	c "github.com/opensds/opensds/client"
	"github.com/opensds/opensds/contrib/cindercompatibleapi/converter"
)

func init() {
	beego.Router("/v3/types", &TypePortal{},
		"post:CreateType;get:ListTypes")
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
	r, _ := http.NewRequest("GET", "/v3/types/1106b972-66ef-11e7-b172-db03f3689c9c", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.ShowTypeRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJSON := `
    {
        "volume_type": {
            "id": "1106b972-66ef-11e7-b172-db03f3689c9c",
            "name": "default",
            "description": "default policy",
            "extra_specs": []
        }
    }`

	var expected converter.ShowTypeRespSpec
	json.Unmarshal([]byte(expectedJSON), &expected)
	expected.VolumeType.IsPublic = true

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
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

	var output converter.ShowTypeRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	if DefaultTypeName != output.VolumeType.Name {
		t.Errorf("Expected %v, actual %v", DefaultTypeName, output.VolumeType.Description)
	}
}

func TestListTypes(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v3/types", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.ListTypesRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJSON := `
    {
        "volume_types": [{
            "id": "1106b972-66ef-11e7-b172-db03f3689c9c",
            "name": "default",
            "description": "default policy",
            "os-volume-type-access:is_public": true,
            "is_public": true
        },
        {
            "id": "2f9c0a04-66ef-11e7-ade2-43158893e017",
            "name": "silver",
            "description": "silver policy",
            "os-volume-type-access:is_public": true,
            "is_public": true,
            "extra_specs": {
                "dataStorage": {
                    "provisioningPolicy": "Thin",
                    "isSpaceEfficient": true
                },
                "ioConnectivity": {
                    "accessProtocol": "rbd",
                    "maxIOPS": 5000000,
                    "maxBWS": 500
                }
            }
        }]
    }`

	var expected converter.ListTypesRespSpec
	json.Unmarshal([]byte(expectedJSON), &expected)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestCreateType(t *testing.T) {
	RequestBodyStr := `
    {
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

	var output converter.CreateTypeRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	var expected converter.CreateTypeRespSpec
	json.Unmarshal([]byte(RequestBodyStr), &expected)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	expected.VolumeType.ID = "1106b972-66ef-11e7-b172-db03f3689c9c"
	expected.VolumeType.IsPublic = true

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestCreateTypeWithBadRequest(t *testing.T) {
	RequestBodyStr := `
    {
        "volume_type": {
            "name": "default",
            "os-volume-type-access:is_public": true,
            "description": "default policy",
        }
    }`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("POST", "/v3/types", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %v, actual %v", http.StatusBadRequest, w.Code)
	}

	var output ErrorSpec
	json.Unmarshal(w.Body.Bytes(), &output)
	expected := "Create a volume type, parse request body failed: invalid character '}' looking for beginning of object key string"

	if expected != output.Message {
		t.Errorf("Expected %v, actual %v", expected, output.Message)
	}

	RequestBodyStr = `
    {
        "volume_type": {
            "name": "default",
            "os-volume-type-access:is_public": false,
            "description": "default policy"
        }
    }`

	jsonStr = []byte(RequestBodyStr)
	r, _ = http.NewRequest("POST", "/v3/types", bytes.NewBuffer(jsonStr))

	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %v, actual %v", http.StatusBadRequest, w.Code)
	}

	json.Unmarshal(w.Body.Bytes(), &output)
	expected = "Create a volume type failed: OpenSDS does not support os-volume-type-access:is_public = false"

	if expected != output.Message {
		t.Errorf("Expected %v, actual %v", expected, output.Message)
	}
}

func TestDeleteType(t *testing.T) {
	r, _ := http.NewRequest("DELETE", "/v3/types/1106b972-66ef-11e7-b172-db03f3689c9c", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected %v, actual %v", http.StatusAccepted, w.Code)
	}
}

func TestUpdateType(t *testing.T) {
	RequestBodyStr := `
    {
        "volume_type": {
            "name": "default",
            "description": "default policy",
            "is_public": true
        }
    }`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("PUT", "/v3/types/1106b972-66ef-11e7-b172-db03f3689c9c", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.UpdateTypeRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	var expected converter.UpdateTypeRespSpec
	json.Unmarshal([]byte(RequestBodyStr), &expected)
	expected.VolumeType.IsPublic = true

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	expected.VolumeType.ID = "1106b972-66ef-11e7-b172-db03f3689c9c"
	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestUpdateTypeWithBadRequest(t *testing.T) {
	RequestBodyStr := `
    {
        "volume_type": {
            "name": "default",
            "description": "default policy",
            "is_public": true,
        }
    }`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("PUT", "/v3/types/1106b972-66ef-11e7-b172-db03f3689c9c", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %v, actual %v", http.StatusBadRequest, w.Code)
	}

	var output ErrorSpec
	json.Unmarshal(w.Body.Bytes(), &output)
	expected := "Update a volume type, parse request body failed: invalid character '}' looking for beginning of object key string"

	if expected != output.Message {
		t.Errorf("Expected %v, actual %v", expected, output.Message)
	}

	RequestBodyStr = `
    {
        "volume_type": {
            "name": "default",
            "description": "default policy",
            "is_public": false
        }
    }`

	jsonStr = []byte(RequestBodyStr)
	r, _ = http.NewRequest("PUT", "/v3/types/1106b972-66ef-11e7-b172-db03f3689c9c", bytes.NewBuffer(jsonStr))
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %v, actual %v", http.StatusBadRequest, w.Code)
	}

	json.Unmarshal(w.Body.Bytes(), &output)
	expected = "Update a volume type failed: OpenSDS does not support is_public = false"

	if expected != output.Message {
		t.Errorf("Expected %v, actual %v", expected, output.Message)
	}
}

func TestAddExtraProperty(t *testing.T) {
	RequestBodyStr := `
    {
        "extra_specs": {
            "dataStorage": {
                "provisioningPolicy": "Thin",
                "isSpaceEfficient": true
            },
            "ioConnectivity": {
                "accessProtocol": "rbd",
                "maxIOPS": 5000000,
                "maxBWS": 500
            }
        }
    }`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("POST", "/v3/types/1106b972-66ef-11e7-b172-db03f3689c9c/extra_specs", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.ExtraSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	var expected converter.ExtraSpec
	json.Unmarshal([]byte(RequestBodyStr), &expected)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestAddExtraPropertyWithBadRequest(t *testing.T) {
	RequestBodyStr := `
    {
        "extra_specs": {
            "dataStorage": {
                "provisioningPolicy": "Thin",
                "isSpaceEfficient": true
            },
            "ioConnectivity": {
                "accessProtocol": "rbd",
                "maxIOPS": 5000000,
                "maxBWS": 500,
            }
        }
    }`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("POST", "/v3/types/1106b972-66ef-11e7-b172-db03f3689c9c/extra_specs", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %v, actual %v", http.StatusBadRequest, w.Code)
	}

	var output ErrorSpec
	json.Unmarshal(w.Body.Bytes(), &output)
	expected := "Create or update extra specs for volume type, parse request body failed: invalid character '}' looking for beginning of object key string"

	if expected != output.Message {
		t.Errorf("Expected %v, actual %v", expected, output.Message)
	}
}

func TestListExtraProperties(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v3/types/1106b972-66ef-11e7-b172-db03f3689c9c/extra_specs", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.ExtraSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	var expected converter.ExtraSpec
	expectedStr := `
    {
        "extra_specs": {
            "dataStorage": {
                "provisioningPolicy": "Thin",
                "isSpaceEfficient": true
            },
            "ioConnectivity": {
                "accessProtocol": "rbd",
                "maxIOPS": 5000000,
                "maxBWS": 500
            }
        }
    }`
	json.Unmarshal([]byte(expectedStr), &expected)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestShowExtraPropertie(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v3/types/1106b972-66ef-11e7-b172-db03f3689c9c/extra_specs/dataStorage", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.ExtraSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	var expected converter.ExtraSpec
	expectedStr := `
    {
        "dataStorage": {
            "provisioningPolicy": "Thin",
            "isSpaceEfficient": true
        }
    }`
	json.Unmarshal([]byte(expectedStr), &expected)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestShowExtraPropertieWithBadRequest(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v3/types/1106b972-66ef-11e7-b172-db03f3689c9c/extra_specs/disk", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.ExtraSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected %v, actual %v", http.StatusNotFound, w.Code)
	}

	if nil != output["disk"] {
		t.Errorf("Expected %v, actual %v", nil, output["disk"])
	}
}

func TestUpdateExtraPropertiy(t *testing.T) {
	RequestBodyStr := `
    {
        "dataStorage": {
            "provisioningPolicy": "Thin",
            "isSpaceEfficient": true
        }
    }`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("PUT", "/v3/types/1106b972-66ef-11e7-b172-db03f3689c9c/extra_specs/dataStorage", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.ExtraSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedStr := `
    {
        "dataStorage": {
            "provisioningPolicy": "Thin",
            "isSpaceEfficient": true
        }
    }`

	var expected converter.ExtraSpec
	json.Unmarshal([]byte(expectedStr), &expected)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestUpdateExtraPropertiyWithBadRequest(t *testing.T) {
	RequestBodyStr := `
    {
        "dataStorage": {
            "provisioningPolicy": "Thin",
            "isSpaceEfficient": true,
        }
    }`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("PUT", "/v3/types/1106b972-66ef-11e7-b172-db03f3689c9c/extra_specs/dataStorage", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %v, actual %v", http.StatusBadRequest, w.Code)
	}

	var output ErrorSpec
	json.Unmarshal(w.Body.Bytes(), &output)
	expected := "Update extra specification for volume type, parse request body failed: invalid character '}' looking for beginning of object key string"

	if expected != output.Message {
		t.Errorf("Expected %v, actual %v", expected, output.Message)
	}

	RequestBodyStr = `
    {
        "Storage": {
            "provisioningPolicy": "Thin",
            "isSpaceEfficient": true
        }
    }`

	jsonStr = []byte(RequestBodyStr)
	r, _ = http.NewRequest("PUT", "/v3/types/1106b972-66ef-11e7-b172-db03f3689c9c/extra_specs/dataStorage", bytes.NewBuffer(jsonStr))
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %v, actual %v", http.StatusBadRequest, w.Code)
	}

	json.Unmarshal(w.Body.Bytes(), &output)
	expected = "Update extra specification for volume type failed: The body of the request is wrong"

	if expected != output.Message {
		t.Errorf("Expected %v, actual %v", expected, output.Message)
	}
}

func TestDeleteExtraPropertie(t *testing.T) {
	r, _ := http.NewRequest("DELETE", "/v3/types/1106b972-66ef-11e7-b172-db03f3689c9c/extra_specs/diskType", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected %v, actual %v", http.StatusAccepted, w.Code)
	}
}
