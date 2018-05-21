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
	beego.Router("/V3/attachments/:attachmentId", &AttachmentPortal{},
		"get:GetAttachment;delete:DeleteAttachment;put:UpdateAttachment")
	beego.Router("/V3/attachments/detail", &AttachmentPortal{},
		"get:ListAttachmentsDetails")

	beego.Router("/V3/attachments", &AttachmentPortal{},
		"post:CreateAttachment;get:ListAttachments")
	if false == IsFakeClient {
		client = NewFakeClient(&c.Config{Endpoint: TestEp})
	}
}

////////////////////////////////////////////////////////////////////////////////
//                            Tests for Attachment                              //
////////////////////////////////////////////////////////////////////////////////
func TestGetAttachment(t *testing.T) {
	r, _ := http.NewRequest("GET", "/V3/attachments/f2dda3d2-bf79-11e7-8665-f750b088f63e", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.ShowAttachmentRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJSON := `
	{
        "attachment": {
            "id": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
            "volume_id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
            "status": "available",
            "connection_info": {
                "driver_volume_type": "iscsi",
                "data": {
                    "discard": false,
                    "targetDiscovered": true,
                    "targetIqn": "iqn.2017-10.io.opensds:volume:00000001",
                    "targetPortal": "127.0.0.0.1:3260"
                }
            }
        }
    }`

	var expected converter.ShowAttachmentRespSpec
	json.Unmarshal([]byte(expectedJSON), &expected)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListAttachments(t *testing.T) {
	r, _ := http.NewRequest("GET", "/V3/attachments", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.ListAttachmentsRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJSON := `
	{
        "attachments": [{
            "id": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
            "volume_id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
            "status": "available"
        }]
    }`

	var expected converter.ListAttachmentsRespSpec
	json.Unmarshal([]byte(expectedJSON), &expected)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListAttachmentsDetails(t *testing.T) {
	r, _ := http.NewRequest("GET", "/V3/attachments/detail", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.ListAttachmentsDetailsRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJSON := `
	{
        "attachments": [{
            "id": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
            "volume_id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
            "status": "available",
            "connection_info": {
                "driver_volume_type": "iscsi",
                "data": {
                    "discard": false,
                    "targetDiscovered": true,
                    "targetIqn": "iqn.2017-10.io.opensds:volume:00000001",
                    "targetPortal": "127.0.0.0.1:3260"
                }
            }
        }]
    }`

	var expected converter.ListAttachmentsDetailsRespSpec
	json.Unmarshal([]byte(expectedJSON), &expected)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestCreateAttachment(t *testing.T) {
	RequestBodyStr := `
	{
        "attachment": {
            "id": "",
            "volume_uuid": "bd5b12a8-a101-11e7-941e-d77981b584d8"
        }
    }`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("POST", "/V3/attachments", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.CreateAttachmentRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJSON := `
	{
        "attachment": {
            "id": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
            "volume_id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
            "status": "available",
            "connection_info": {
                "driver_volume_type": "iscsi",
                "data": {
                    "discard": false,
                    "targetDiscovered": true,
                    "targetIqn": "iqn.2017-10.io.opensds:volume:00000001",
                    "targetPortal": "127.0.0.0.1:3260"
                }
            }
        }
    }`

	var expected converter.CreateAttachmentRespSpec
	json.Unmarshal([]byte(expectedJSON), &expected)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	expected.Attachment.ID = "f2dda3d2-bf79-11e7-8665-f750b088f63e"
	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}

}

func TestCreateAttachmentWithBadRequest(t *testing.T) {
	RequestBodyStr := `
	{
        "attachment": {
            "id": "",
            "volume_uuid": "bd5b12a8-a101-11e7-941e-d77981b584d8",
        }
    }`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("POST", "/V3/attachments", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %v, actual %v", http.StatusBadRequest, w.Code)
	}

	var output ErrorSpec
	json.Unmarshal(w.Body.Bytes(), &output)
	expected := "Create attachment, parse request body failed: invalid character '}' looking for beginning of object key string"

	if expected != output.Message {
		t.Errorf("Expected %v, actual %v", expected, output.Message)
	}
}

func TestDeleteAttachment(t *testing.T) {
	r, _ := http.NewRequest("DELETE", "/V3/attachments/f2dda3d2-bf79-11e7-8665-f750b088f63e", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}
}

func TestUpdateAttachment(t *testing.T) {
	RequestBodyStr := `
    {
        "attachment": {
            "connector": {
                "ip": "127.0.0.0.1"
            }
        }
    }`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("PUT", "/V3/attachments/f2dda3d2-bf79-11e7-8665-f750b088f63e", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.UpdateAttachmentRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJSON := `
	{
        "attachment": {
            "id": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
            "volume_id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
            "status": "available",
            "connection_info": {
                "driver_volume_type": "iscsi",
                "data": {
                    "discard": false,
                    "targetDiscovered": true,
                    "targetIqn": "iqn.2017-10.io.opensds:volume:00000001",
                    "targetPortal": "127.0.0.0.1:3260"
                }
            }
        }
    }`

	var expected converter.UpdateAttachmentRespSpec
	json.Unmarshal([]byte(expectedJSON), &expected)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, actual %v", http.StatusOK, w.Code)
	}

	expected.Attachment.ID = "f2dda3d2-bf79-11e7-8665-f750b088f63e"
	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestUpdateAttachmentWithBadRequest(t *testing.T) {
	RequestBodyStr := `
    {
        "attachment": {
            "connector": {
                "ip": "127.0.0.0.1",
            }
        }
    }`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("PUT", "/V3/attachments/f2dda3d2-bf79-11e7-8665-f750b088f63e", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %v, actual %v", http.StatusBadRequest, w.Code)
	}

	var output ErrorSpec
	json.Unmarshal(w.Body.Bytes(), &output)
	expected := "Update an attachment, parse request body failed: invalid character '}' looking for beginning of object key string"

	if expected != output.Message {
		t.Errorf("Expected %v, actual %v", expected, output.Message)
	}
}
