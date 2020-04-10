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
	c "github.com/sodafoundation/api/pkg/context"
	"github.com/sodafoundation/api/pkg/db"
	"github.com/sodafoundation/api/pkg/model"
	. "github.com/sodafoundation/api/testutils/collection"
	dbtest "github.com/sodafoundation/api/testutils/db/testing"
)

////////////////////////////////////////////////////////////////////////////////
//                      Prepare for mock server                               //
////////////////////////////////////////////////////////////////////////////////

func init() {
	beego.Router("/v1beta/block/attachments", &VolumeAttachmentPortal{},
		"post:CreateVolumeAttachment;get:ListVolumeAttachments")
	beego.Router("/v1beta/block/attachments/:attachmentId", &VolumeAttachmentPortal{},
		"get:GetVolumeAttachment;put:UpdateVolumeAttachment;delete:DeleteVolumeAttachment")
}

////////////////////////////////////////////////////////////////////////////////
//                         Tests for volume attachment                          //
////////////////////////////////////////////////////////////////////////////////

func TestListVolumeAttachments(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		var sampleAttachments = []*model.VolumeAttachmentSpec{&SampleAttachments[0]}
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"volumeId": {"bd5b12a8-a101-11e7-941e-d77981b584d8"},
			"offset":   {"0"},
			"limit":    {"1"},
			"sortDir":  {"asc"},
			"sortKey":  {"name"},
		}
		mockClient.On("ListVolumeAttachmentsWithFilter", c.NewAdminContext(), m).
			Return(sampleAttachments, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/attachments?volumeId=bd5b12a8-a101-11e7-941e-d77981b584d8&offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		var output []*model.VolumeAttachmentSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, sampleAttachments)
	})

	t.Run("Should return 500 if list volume attachments with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"volumeId": {"bd5b12a8-a101-11e7-941e-d77981b584d8"},
			"offset":   {"0"},
			"limit":    {"1"},
			"sortDir":  {"asc"},
			"sortKey":  {"name"},
		}
		mockClient.On("ListVolumeAttachmentsWithFilter", c.NewAdminContext(), m).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET",
			"/v1beta/block/attachments?volumeId=bd5b12a8-a101-11e7-941e-d77981b584d8&offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestGetVolumeAttachment(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolumeAttachment", c.NewAdminContext(), "f2dda3d2-bf79-11e7-8665-f750b088f63e").
			Return(&SampleAttachments[0], nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/attachments/f2dda3d2-bf79-11e7-8665-f750b088f63e", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		var output model.VolumeAttachmentSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SampleAttachments[0])
	})

	t.Run("Should return 404 if get volume attachment with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolumeAttachment", c.NewAdminContext(), "f2dda3d2-bf79-11e7-8665-f750b088f63e").
			Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/attachments/f2dda3d2-bf79-11e7-8665-f750b088f63e", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}

func TestUpdateVolumeAttachment(t *testing.T) {
	var jsonStr = []byte(`{
		"id": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
		"name": "fake volume attachment",
		"description": "fake volume attachment"
	}`)
	var expectedJson = []byte(`{
		"id": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
		"name": "fake volume attachment",
		"description": "fake volume attachment",
		"status": "available",
		"volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
		"hostId": "202964b5-8e73-46fd-b41b-a8e403f3c30b",
		"connectionInfo": {
			"driverVolumeType": "iscsi",
			"data": {
				"targetDiscovered": true,
				"targetIqn": "iqn.2017-10.io.opensds:volume:00000001",
				"targetPortal": "127.0.0.0.1:3260",
				"discard": false
			}
		}
	}`)
	var expected model.VolumeAttachmentSpec
	json.Unmarshal(expectedJson, &expected)

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		attachment := model.VolumeAttachmentSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&attachment)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateVolumeAttachment", c.NewAdminContext(), attachment.Id, &attachment).
			Return(&expected, nil)
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/block/attachments/f2dda3d2-bf79-11e7-8665-f750b088f63e", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.VolumeAttachmentSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &expected)
	})

	t.Run("Should return 500 if update volume attachment with bad request", func(t *testing.T) {
		attachment := model.VolumeAttachmentSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&attachment)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateVolumeAttachment", c.NewAdminContext(), attachment.Id, &attachment).
			Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/block/attachments/f2dda3d2-bf79-11e7-8665-f750b088f63e", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}
