// Copyright 2017 The OpenSDS Authors.
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
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
)

func init() {
	beego.Router("/v1beta/block/volumes", &VolumePortal{},
		"post:CreateVolume;get:ListVolumes")
	beego.Router("/v1beta/block/volumes/:volumeId", &VolumePortal{},
		"get:GetVolume;put:UpdateVolume;delete:DeleteVolume")

	beego.Router("/v1beta/block/attachments", &VolumeAttachmentPortal{},
		"post:CreateVolumeAttachment;get:ListVolumeAttachments")
	beego.Router("/v1beta/block/attachments/:attachmentId", &VolumeAttachmentPortal{},
		"get:GetVolumeAttachment;put:UpdateVolumeAttachment;delete:DeleteVolumeAttachment")

	beego.Router("/v1beta/block/snapshots", &VolumeSnapshotPortal{},
		"post:CreateVolumeSnapshot;get:ListVolumeSnapshots")
	beego.Router("/v1beta/block/snapshots/:snapshotId", &VolumeSnapshotPortal{},
		"get:GetVolumeSnapshot;put:UpdateVolumeSnapshot;delete:DeleteVolumeSnapshot")
}

////////////////////////////////////////////////////////////////////////////////
//                            Tests for volume                               //
////////////////////////////////////////////////////////////////////////////////

var (
	fakeVolume = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			CreatedAt: "2017-10-24T16:21:32",
		},
		Name:             "fake Vol",
		Description:      "fake Vol",
		Size:             99,
		AvailabilityZone: "unknown",
		Status:           "available",
		PoolId:           "831fa5fb-17cf-4410-bec6-1f4b06208eef",
		ProfileId:        "d3a109ff-3e51-4625-9054-32604c79fa90",
	}
	fakeVolumes = []*model.VolumeSpec{fakeVolume}
)

func TestListVolumes(t *testing.T) {

	mockClient := new(dbtest.MockClient)
	mockClient.On("ListVolumes").Return(fakeVolumes, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/volumes", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output []model.VolumeSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `[{
		    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"createdAt": "2017-10-24T16:21:32",
			"name": "fake Vol",
			"description": "fake Vol",
			"size": 99,
			"availabilityZone": "unknown",
			"profileId": "d3a109ff-3e51-4625-9054-32604c79fa90",
			"status": "available",
			"poolId": "831fa5fb-17cf-4410-bec6-1f4b06208eef"
		}]`

	var expected []model.VolumeSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListVolumesWithBadRequest(t *testing.T) {

	mockClient := new(dbtest.MockClient)
	mockClient.On("ListVolumes").Return(nil, errors.New("db error"))
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/volumes", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}

func TestGetVolume(t *testing.T) {

	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolume", "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(fakeVolume, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/volumes/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.VolumeSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `{
		    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"createdAt": "2017-10-24T16:21:32",
			"name": "fake Vol",
			"description": "fake Vol",
			"size": 99,
			"availabilityZone": "unknown",
			"profileId": "d3a109ff-3e51-4625-9054-32604c79fa90",
			"status": "available",
			"poolId": "831fa5fb-17cf-4410-bec6-1f4b06208eef"
		}`

	var expected model.VolumeSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestGetVolumeWithBadRequest(t *testing.T) {

	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolume", "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(nil, errors.New("db error"))
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/volumes/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}

func TestUpdateVolume(t *testing.T) {
	var jsonStr = []byte(`{"name":"fake Vol","description":"fake Vol"}`)
	r, _ := http.NewRequest("PUT",
		"/v1beta/block/volumes/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/JSON")

	var volume = model.VolumeSpec{
		BaseModel: &model.BaseModel{},
	}
	json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&volume)
	volume.Id = "f4a5e666-c669-4c64-a2a1-8f9ecd560c78"

	mockClient := new(dbtest.MockClient)
	mockClient.On("UpdateVolume", "f4a5e666-c669-4c64-a2a1-8f9ecd560c78", &volume).Return(fakeVolume, nil)
	db.C = mockClient
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.VolumeSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `{
		    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"createdAt": "2017-10-24T16:21:32",
			"name": "fake Vol",
			"description": "fake Vol",
			"size": 99,
			"availabilityZone": "unknown",
			"profileId": "d3a109ff-3e51-4625-9054-32604c79fa90",
			"status": "available",
			"poolId": "831fa5fb-17cf-4410-bec6-1f4b06208eef"
		}`

	var expected model.VolumeSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestUpdateVolumeWithBadRequest(t *testing.T) {
	var jsonStr = []byte(``)
	r, _ := http.NewRequest("PUT",
		"/v1beta/block/volumes/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/JSON")
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}

	jsonStr = []byte(`{"name":"fake Vol","description":"fake Vol"}`)
	r, _ = http.NewRequest("PUT",
		"/v1beta/block/volumes/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", bytes.NewBuffer(jsonStr))
	w = httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/JSON")

	var volume = model.VolumeSpec{
		BaseModel: &model.BaseModel{},
	}
	json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&volume)
	volume.Id = "f4a5e666-c669-4c64-a2a1-8f9ecd560c78"

	mockClient := new(dbtest.MockClient)
	mockClient.On("UpdateVolume", "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
		&volume).Return(nil, errors.New("db error"))
	db.C = mockClient
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}

////////////////////////////////////////////////////////////////////////////////
//                         Tests for volume snapshot                          //
////////////////////////////////////////////////////////////////////////////////

var (
	fakeSnapshot = &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			CreatedAt: "2017-10-24T16:21:32",
		},
		Name:        "fake snapshot",
		Description: "fake snapshot",
		Size:        99,
		Status:      "available",
		VolumeId:    "d3a109ff-3e51-4625-9054-32604c79fa90",
	}
	fakeSnapshots = []*model.VolumeSnapshotSpec{fakeSnapshot}
)

func TestListVolumeSnapshots(t *testing.T) {

	mockClient := new(dbtest.MockClient)
	mockClient.On("ListVolumeSnapshots").Return(fakeSnapshots, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/snapshots", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output []model.VolumeSnapshotSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `[{
		    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"createdAt": "2017-10-24T16:21:32",
			"name": "fake snapshot",
			"description": "fake snapshot",
			"size": 99,
			"volumeId": "d3a109ff-3e51-4625-9054-32604c79fa90",
			"status": "available"
		}]`

	var expected []model.VolumeSnapshotSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListVolumeSnapshotsWithBadRequest(t *testing.T) {

	mockClient := new(dbtest.MockClient)
	mockClient.On("ListVolumeSnapshots").Return(nil, errors.New("db error"))
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/snapshots", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}

func TestGetVolumeSnapshot(t *testing.T) {

	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolumeSnapshot", "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(fakeSnapshot, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/snapshots/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.VolumeSnapshotSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `{
		    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"createdAt": "2017-10-24T16:21:32",
			"name": "fake snapshot",
			"description": "fake snapshot",
			"size": 99,
			"volumeId": "d3a109ff-3e51-4625-9054-32604c79fa90",
			"status": "available"
		}`

	var expected model.VolumeSnapshotSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestGetVolumeSnapshotWithBadRequest(t *testing.T) {

	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolumeSnapshot", "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(nil, errors.New("db error"))
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/snapshots/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}

func TestUpdateVolumeSnapshot(t *testing.T) {
	var jsonStr = []byte(`{"name":"fake snapshot","description":"fake snapshot"}`)
	r, _ := http.NewRequest("PUT",
		"/v1beta/block/snapshots/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/JSON")

	var snapshot = model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{},
	}
	json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&snapshot)
	snapshot.Id = "f4a5e666-c669-4c64-a2a1-8f9ecd560c78"

	mockClient := new(dbtest.MockClient)
	mockClient.On("UpdateVolumeSnapshot", "f4a5e666-c669-4c64-a2a1-8f9ecd560c78", &snapshot).Return(fakeSnapshot, nil)
	db.C = mockClient
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.VolumeSnapshotSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `{
		    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"createdAt": "2017-10-24T16:21:32",
			"name": "fake snapshot",
			"description": "fake snapshot",
			"size": 99,
			"volumeId": "d3a109ff-3e51-4625-9054-32604c79fa90",
			"status": "available"
		}`

	var expected model.VolumeSnapshotSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestUpdateVolumeSnapshotWithBadRequest(t *testing.T) {
	var jsonStr = []byte(``)
	r, _ := http.NewRequest("PUT",
		"/v1beta/block/snapshots/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/JSON")
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}

	jsonStr = []byte(`{"name":"fake snapshot","description":"fake snapshot"}`)
	r, _ = http.NewRequest("PUT",
		"/v1beta/block/snapshots/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", bytes.NewBuffer(jsonStr))
	w = httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/JSON")

	var snapshot = model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{},
	}
	json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&snapshot)
	snapshot.Id = "f4a5e666-c669-4c64-a2a1-8f9ecd560c78"

	mockClient := new(dbtest.MockClient)
	mockClient.On("UpdateVolumeSnapshot", "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
		&snapshot).Return(nil, errors.New("db error"))
	db.C = mockClient
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}

////////////////////////////////////////////////////////////////////////////////
//                         Tests for volume attachment                          //
////////////////////////////////////////////////////////////////////////////////

var (
	fakeAttachment = &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			CreatedAt: "2017-10-24T16:21:32",
		},
		Status:   "available",
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		ConnectionInfo: model.ConnectionInfo{
			DriverVolumeType: "iscsi",
			ConnectionData: map[string]interface{}{
				"targetDiscovered": true,
				"targetIqn":        "iqn.2017-10.io.opensds:volume:00000001",
				"targetPortal":     "127.0.0.0.1:3260",
				"discard":          false,
			},
		},
	}
	fakeAttachments = []*model.VolumeAttachmentSpec{fakeAttachment}
)

func TestListVolumeAttachments(t *testing.T) {

	mockClient := new(dbtest.MockClient)
	mockClient.On("ListVolumeAttachments", "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(fakeAttachments, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/attachments?volumeId=bd5b12a8-a101-11e7-941e-d77981b584d8", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output []model.VolumeAttachmentSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `[
	  {
	    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
	    "createdAt": "2017-10-24T16:21:32",
	    "volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
	    "status": "available",
	    "connectionInfo": {
	      "driverVolumeType": "iscsi",
	      "data": {
	        "discard": false,
	        "targetDiscovered": true,
	        "targetIqn": "iqn.2017-10.io.opensds:volume:00000001",
	        "targetPortal": "127.0.0.0.1:3260"
	      }
	    }
	  }
	]`

	var expected []model.VolumeAttachmentSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}

}

func TestGetVolumeAttachment(t *testing.T) {

	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolumeAttachment", "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(fakeAttachment, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/attachments/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.VolumeAttachmentSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `{
	    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
	    "createdAt": "2017-10-24T16:21:32",
	    "volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
	    "status": "available",
	    "connectionInfo": {
	      "driverVolumeType": "iscsi",
	      "data": {
	        "discard": false,
	        "targetDiscovered": true,
	        "targetIqn": "iqn.2017-10.io.opensds:volume:00000001",
	        "targetPortal": "127.0.0.0.1:3260"
	      }
	    }
	  }`

	var expected model.VolumeAttachmentSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestUpdateVolumeAttachment(t *testing.T) {
	var jsonStr = []byte(`{"status": "available"}`)
	r, _ := http.NewRequest("PUT",
		"/v1beta/block/attachments/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/JSON")

	var attachment = model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{},
	}
	json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&attachment)
	attachment.Id = "f4a5e666-c669-4c64-a2a1-8f9ecd560c78"

	mockClient := new(dbtest.MockClient)
	mockClient.On("UpdateVolumeAttachment", "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
		&attachment).Return(fakeAttachment, nil)
	db.C = mockClient
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.VolumeAttachmentSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `{
	    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
	    "createdAt": "2017-10-24T16:21:32",
	    "volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
	    "status": "available",
	    "connectionInfo": {
	      "driverVolumeType": "iscsi",
	      "data": {
	        "discard": false,
	        "targetDiscovered": true,
	        "targetIqn": "iqn.2017-10.io.opensds:volume:00000001",
	        "targetPortal": "127.0.0.0.1:3260"
	      }
	    }
	  }`

	var expected model.VolumeAttachmentSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestUpdateVolumeAttachmentWithBadRequest(t *testing.T) {
	var jsonStr = []byte(``)
	r, _ := http.NewRequest("PUT",
		"/v1beta/block/attachments/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/JSON")

	beego.BeeApp.Handlers.ServeHTTP(w, r)
	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}

	jsonStr = []byte(`{"status": "available"}`)
	r, _ = http.NewRequest("PUT",
		"/v1beta/block/attachments/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", bytes.NewBuffer(jsonStr))
	w = httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/JSON")

	var attachment = model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{},
	}
	json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&attachment)
	attachment.Id = "f4a5e666-c669-4c64-a2a1-8f9ecd560c78"

	mockClient := new(dbtest.MockClient)
	mockClient.On("UpdateVolumeAttachment", "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
		&attachment).Return(nil, errors.New("db error"))
	db.C = mockClient
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}
