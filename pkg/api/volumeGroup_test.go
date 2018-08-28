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
	//"github.com/opensds/opensds/pkg/controller"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	//. "github.com/opensds/opensds/testutils/collection"
	c "github.com/opensds/opensds/pkg/context"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
)

func init() {
	beego.Router("/v1beta/block/groups", &VolumeGroupPortal{}, "post:CreateVolumeGroup")
	beego.Router("/v1beta/block/groups/:groupId", &VolumeGroupPortal{}, "put:UpdateVolumeGroup;get:GetVolumeGroup;delete:DeleteVolumeGroup")
}

var (
	fakeVolumeGroup = &model.VolumeGroupSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			CreatedAt: "2017-10-24T16:21:32",
		},
		Name:             "fakeGroup",
		Description:      "fakeGroup",
		AvailabilityZone: "unknown",
		Status:           "available",
		PoolId:           "831fa5fb-17cf-4410-bec6-1f4b06208eef",
	}
	fakeVolumeGroups = []*model.VolumeGroupSpec{fakeVolumeGroup}
)

//func TestListVolumeGroups(t *testing.T) {

//	mockClient := new(dbtest.Client)
//	mockClient.On("ListVolumeGroups").Return(fakeVolumeGroups, nil)
//	db.C = mockClient

//	r, _ := http.NewRequest("GET", "/v1beta/block/groups?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
//	w := httptest.NewRecorder()
//	beego.BeeApp.Handlers.ServeHTTP(w, r)

//	var output []model.VolumeGroupSpec
//	json.Unmarshal(w.Body.Bytes(), &output)

//	expectedJson := `[{
//		    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
//			"createdAt": "2017-10-24T16:21:32",
//			"name": "fakeGroup",
//			"description": "fakeGroup",
//			"availabilityZone": "unknown",
//			"status": "available",
//			"poolId": "831fa5fb-17cf-4410-bec6-1f4b06208eef"
//		}]`

//	var expected []model.VolumeGroupSpec
//	json.Unmarshal([]byte(expectedJson), &expected)

//	if w.Code != 200 {
//		t.Errorf("Expected 200, actual %v", w.Code)
//	}

//	if !reflect.DeepEqual(expected, output) {
//		t.Errorf("Expected %v, actual %v", expected, output)
//	}
//}

//func TestListVolumeGroupsWithBadRequest(t *testing.T) {

//	mockClient := new(dbtest.Client)
//	mockClient.On("ListVolumeGroups").Return(nil, errors.New("db error"))
//	db.C = mockClient

//	r, _ := http.NewRequest("GET", "/v1beta/block/groups?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
//	w := httptest.NewRecorder()
//	beego.BeeApp.Handlers.ServeHTTP(w, r)

//	if w.Code != 400 {
//		t.Errorf("Expected 400, actual %v", w.Code)
//	}
//}

func TestGetVolumeGroup(t *testing.T) {

	mockClient := new(dbtest.Client)
	mockClient.On("GetVolumeGroup", c.NewAdminContext(), "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(fakeVolumeGroup, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/groups/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.VolumeGroupSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `{
		    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"createdAt": "2017-10-24T16:21:32",
			"name": "fakeGroup",
			"description": "fakeGroup",
			"availabilityZone": "unknown",
			"status": "available",
			"poolId": "831fa5fb-17cf-4410-bec6-1f4b06208eef"
		}`

	var expected model.VolumeGroupSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestGetVolumeGroupWithBadRequest(t *testing.T) {

	mockClient := new(dbtest.Client)
	mockClient.On("GetVolumeGroup", c.NewAdminContext(), "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(nil, errors.New("db error"))
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/block/groups/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}

var (
	fakeGroupVolumeTest = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: "f4a5e666-c669-4c64-a2a1-8f9ecd560c70",
		},
		Name:        "sample-volume",
		Description: "This is a sample volume for testing",
		Size:        int64(1),
		Status:      "available",
		PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
	}
	fakeGroupVolumes = []*model.VolumeSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "f4a5e666-c669-4c64-a2a1-8f9ecd560c71",
			},
			Name:        "sample-volume",
			Description: "This is a sample volume for testing",
			Size:        int64(1),
			Status:      "available",
			PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
			GroupId:     "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "f4a5e666-c669-4c64-a2a1-8f9ecd560c72",
			},
			Name:        "sample-volume",
			Description: "This is a sample volume for testing",
			Size:        int64(1),
			Status:      "available",
			PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
			GroupId:     "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "f4a5e666-c669-4c64-a2a1-8f9ecd560c73",
			},
			Name:        "sample-volume",
			Description: "This is a sample volume for testing",
			Size:        int64(1),
			Status:      "available",
			PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
			GroupId:     "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "f4a5e666-c669-4c64-a2a1-8f9ecd560c74",
			},
			Name:        "sample-volume",
			Description: "This is a sample volume for testing",
			Size:        int64(1),
			Status:      "other",
			PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
			GroupId:     "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "f4a5e666-c669-4c64-a2a1-8f9ecd560c75",
			},
			Name:        "sample-volume",
			Description: "This is a sample volume for testing",
			Size:        int64(1),
			Status:      "available",
			PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
			GroupId:     "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "f4a5e666-c669-4c64-a2a1-8f9ecd560c76",
			},
			Name:        "sample-volume",
			Description: "This is a sample volume for testing",
			Size:        int64(1),
			Status:      "available",
			PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
			GroupId:     "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
		},
	}
)

func TestUpdateVolumeGroup(t *testing.T) {
	var jsonStr = []byte(
		`{
		    "id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"name": "fakeGroupUpdate",
			"description": "fakeGroupUpdate",
			"addVolumes":["f4a5e666-c669-4c64-a2a1-8f9ecd560c70","f4a5e666-c669-4c64-a2a1-8f9ecd560c72","f4a5e666-c669-4c64-a2a1-8f9ecd560c73"],
            "removeVolumes":["f4a5e666-c669-4c64-a2a1-8f9ecd560c74","f4a5e666-c669-4c64-a2a1-8f9ecd560c75","f4a5e666-c669-4c64-a2a1-8f9ecd560c76"]
		}`)
	r, _ := http.NewRequest("PUT",
		"/v1beta/block/groups/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/JSON")

	var vg = model.VolumeGroupSpec{
		BaseModel: &model.BaseModel{},
	}
	json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&vg)

	mockClient := new(dbtest.Client)

	mockClient.On("GetVolumeGroup", c.NewAdminContext(), fakeVolumeGroup.Id).Return(fakeVolumeGroup, nil)
	mockClient.On("ListVolumesByGroupId", c.NewAdminContext(), fakeVolumeGroup.Id).Return(fakeGroupVolumes, nil)
	mockClient.On("GetVolume", c.NewAdminContext(), fakeGroupVolumeTest.Id).Return(fakeGroupVolumeTest, nil)
	mockClient.On("GetDockByPoolId", c.NewAdminContext(), fakeVolumeGroup.PoolId).Return(nil, errors.New("db error"))
	var vgUpdate = &model.VolumeGroupSpec{
		BaseModel:   &model.BaseModel{},
		Name:        "fakeGroupUpdate",
		Status:      "updating",
		Description: "fakeGroupUpdate",
	}
	mockClient.On("UpdateVolumeGroup", vgUpdate).Return(nil, errors.New("db error"))
	db.C = mockClient
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}
}

func TestDeleteVolumeGroup(t *testing.T) {
	fakeVolumeGroupDelete := &model.VolumeGroupSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			CreatedAt: "2017-10-24T16:21:32",
		},
		Name:             "fakeGroup",
		Description:      "fakeGroup",
		AvailabilityZone: "unknown",
		Status:           "available",
		PoolId:           "831fa5fb-17cf-4410-bec6-1f4b06208eef",
		//GroupSnapshots:   []string{"feafef"},
	}
	Snapshots := []*model.VolumeSnapshotSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "3769855c-a102-11e7-b772-17b880d2f537",
			},
			Name:        "sample-snapshot-01",
			Description: "This is the first sample snapshot for testing",
			Size:        int64(1),
			Status:      "created",
			VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
	}
	mockClient := new(dbtest.Client)
	mockClient.On("GetVolumeGroup", c.NewAdminContext(), "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(fakeVolumeGroupDelete, nil)
	mockClient.On("ListVolumesByGroupId", c.NewAdminContext(), fakeVolumeGroup.Id).Return(fakeGroupVolumes, nil)
	mockClient.On("ListSnapshotsByVolumeId", c.NewAdminContext(), fakeGroupVolumes[0].Id).Return(Snapshots, nil)
	mockClient.On("GetDockByPoolId", c.NewAdminContext(), fakeVolumeGroupDelete.PoolId).Return(nil, nil)
	db.C = mockClient

	r, _ := http.NewRequest("DELETE", "/v1beta/block/groups/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}
}
