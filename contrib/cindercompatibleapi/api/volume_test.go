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
	"time"

	"github.com/astaxie/beego"
	c "github.com/opensds/opensds/client"
	"github.com/opensds/opensds/contrib/cindercompatibleapi/converter"
)

func init() {
	beego.Router("/v3/volumes/:volumeId/action", &VolumePortal{},
		"post:VolumeAction")
	beego.Router("/v3/volumes/:volumeId", &VolumePortal{},
		"get:GetVolume;delete:DeleteVolume;put:UpdateVolume")
	beego.Router("/v3/volumes/detail", &VolumePortal{},
		"get:ListVolumesDetails")
	beego.Router("/v3/volumes", &VolumePortal{},
		"post:CreateVolume;get:ListVolumes")
	if false == IsFakeClient {
		client = NewFakeClient(&c.Config{Endpoint: TestEp})
	}
}

////////////////////////////////////////////////////////////////////////////////
//                            Tests for Volume                              //
////////////////////////////////////////////////////////////////////////////////
func TestGetVolume(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v3/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.ShowVolumeRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJSON := `
    {
        "volume": {
            "id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
            "name": "sample-volume",
            "description": "This is a sample volume for testing",
            "metadata": {
                
            },
            "size": 1
        }
    }`

	var expected converter.ShowVolumeRespSpec
	json.Unmarshal([]byte(expectedJSON), &expected)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	expected.Volume.Attachments = make([]converter.RespAttachment, 0, 0)
	expected.Volume.Status = "available"
	expected.Volume.VolumeType = "1106b972-66ef-11e7-b172-db03f3689c9c"

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListVolumes(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v3/volumes", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.ListVolumesRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJSON := `
    {
        "volumes": [{
            "id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
            "metadata": {
                
            },
            "name": "sample-volume"
        }]
    }`

	var expected converter.ListVolumesRespSpec
	json.Unmarshal([]byte(expectedJSON), &expected)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListVolumesDetails(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v3/volumes/detail", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.ListVolumesDetailsRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJSON := `
    {
        "volumes": [{
            "id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
            "size": 1,
            "status": "available",
            "description": "This is a sample volume for testing",
            "metadata": {
                
            },
			"volume_type": "1106b972-66ef-11e7-b172-db03f3689c9c",
            "name": "sample-volume"
        }]
    }`

	var expected converter.ListVolumesDetailsRespSpec
	json.Unmarshal([]byte(expectedJSON), &expected)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	expected.Volumes[0].Attachments = make([]converter.RespAttachment, 0, 0)
	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestCreateVolume(t *testing.T) {
	RequestBodyStr := `
    {
        "volume": {
            "name": "sample-volume",
            "description": "This is a sample volume for testing",
            "size": 1
        }
    }`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("POST", "/v3/volumes", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.CreateVolumeRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	var expected converter.CreateVolumeRespSpec
	json.Unmarshal([]byte(RequestBodyStr), &expected)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected %v, actual %v", http.StatusAccepted, w.Code)
	}

	expected.Volume.Attachments = make([]converter.RespAttachment, 0, 0)
	expected.Volume.Status = "available"
	expected.Volume.ID = "bd5b12a8-a101-11e7-941e-d77981b584d8"
	expected.Volume.Metadata = make(map[string]string)
	expected.Volume.VolumeType = "1106b972-66ef-11e7-b172-db03f3689c9c"
	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestCreateVolumeWithBadRequest(t *testing.T) {
	RequestBodyStr := `
    {
        "volume": {
            "name": "sample-volume",
            "description": "This is a sample volume for testing",
            "size": 1,
        }
    }`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("POST", "/v3/volumes", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %v, actual %v", http.StatusBadRequest, w.Code)
	}

	var output ErrorSpec
	json.Unmarshal(w.Body.Bytes(), &output)
	expected := "Create a volume, parse request body failed: invalid character '}' looking for beginning of object key string"

	if expected != output.Message {
		t.Errorf("Expected %v, actual %v", expected, output.Message)
	}

	RequestBodyStr = `
    {
        "volume": {
            "name": "sample-volume",
            "description": "This is a sample volume for testing",
            "size": 1,
			"multiattach": true
        }
    }`

	jsonStr = []byte(RequestBodyStr)
	r, _ = http.NewRequest("POST", "/v3/volumes", bytes.NewBuffer(jsonStr))
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %v, actual %v", http.StatusBadRequest, w.Code)
	}

	json.Unmarshal(w.Body.Bytes(), &output)
	expected = "Create a volume failed: OpenSDS does not support the parameter: id/source_volid/multiattach/snapshot_id/backup_id/imageRef/metadata/consistencygroup_id"

	if expected != output.Message {
		t.Errorf("Expected %v, actual %v", expected, output.Message)
	}

}

func TestDeleteVolume(t *testing.T) {
	r, _ := http.NewRequest("DELETE", "/v3/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected %v, actual %v", http.StatusAccepted, w.Code)
	}
}

func TestUpdateVolume(t *testing.T) {
	RequestBodyStr := `
    {
        "volume": {
            "name": "sample-volume",
            "multiattach": false,
            "description": "This is a sample volume for testing"
        }
    }`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("PUT", "/v3/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.UpdateVolumeRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	var expected converter.UpdateVolumeRespSpec
	json.Unmarshal([]byte(RequestBodyStr), &expected)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	expected.Volume.Attachments = make([]converter.RespAttachment, 0, 0)
	expected.Volume.Status = "available"
	expected.Volume.ID = "bd5b12a8-a101-11e7-941e-d77981b584d8"
	expected.Volume.Size = 1
	expected.Volume.Metadata = make(map[string]string)
	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestUpdateVolumeWithBadRequest(t *testing.T) {
	RequestBodyStr := `
    {
        "volume": {
            "name": "sample-volume",
            "multiattach": false,
            "description": "This is a sample volume for testing",
        }
    }`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("PUT", "/v3/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %v, actual %v", http.StatusBadRequest, w.Code)
	}

	var output ErrorSpec
	json.Unmarshal(w.Body.Bytes(), &output)
	expected := "Update a volume, parse request body failed: invalid character '}' looking for beginning of object key string"

	if expected != output.Message {
		t.Errorf("Expected %v, actual %v", expected, output.Message)
	}

	RequestBodyStr = `
    {
        "volume": {
            "name": "sample-volume",
            "multiattach": false,
            "description": "This is a sample volume for testing",
            "metadata": {
                "key1": "value1"
            }
        }
    }`

	jsonStr = []byte(RequestBodyStr)
	r, _ = http.NewRequest("PUT", "/v3/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8", bytes.NewBuffer(jsonStr))

	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %v, actual %v", http.StatusBadRequest, w.Code)
	}

	json.Unmarshal(w.Body.Bytes(), &output)
	expected = "Update a volume failed: OpenSDS does not support the parameter: metadata"

	if expected != output.Message {
		t.Errorf("Expected %v, actual %v", expected, output.Message)
	}
}

func TestVolumeAction(t *testing.T) {
	RequestBodyStr := `{"os-reserve": null}`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("POST", "/v3/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8/action", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected %v, actual %v", http.StatusAccepted, w.Code)
	}

	RequestBodyStr = `
    {
        "os-extend": {
            "new_size": 3
        }
    }`

	jsonStr = []byte(RequestBodyStr)
	r, _ = http.NewRequest("POST", "/v3/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8/action", bytes.NewBuffer(jsonStr))
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected %v, actual %v", http.StatusNotFound, w.Code)
	}
}

func TestVolumeActionInitializeConnectionWithError(t *testing.T) {
	SleepDuration = time.Nanosecond
	Req := converter.InitializeConnectionReqSpec{}

	Req.InitializeConnection.Connector.Platform = "x86_64"
	Req.InitializeConnection.Connector.Host = "ubuntu"
	Req.InitializeConnection.Connector.DoLocalAttach = false
	Req.InitializeConnection.Connector.IP = "10.10.3.173"
	Req.InitializeConnection.Connector.OsType = "linux2"
	Req.InitializeConnection.Connector.Multipath = false
	Req.InitializeConnection.Connector.Initiator = "iqn.1993-08.org.debian:01:6acaf7eab14"
	body, _ := json.Marshal(Req)

	RequestBodyStr := string(body)
	fmt.Println("397")
	fmt.Println(string(body))
	fmt.Println("399")
	fmt.Println(RequestBodyStr)

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("POST", "/v3/volumes/bd5b12a8-a101-11e7-941e-d77981b584d8/action", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected %v, actual %v", http.StatusInternalServerError, w.Code)
	}

	var output ErrorSpec
	json.Unmarshal(w.Body.Bytes(), &output)
	expected := "Initialize connection, attachment is not available or connectionInfo is incorrect"

	if expected != output.Message {
		t.Errorf("Expected %v, actual %v", expected, output.Message)
	}
}
