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
	"errors"
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/testutils/collection"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
)

func TestCreateVolumeDBEntry(t *testing.T) {
	var req = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id:        "9f1514a7-ca31-4de9-8a27-54c9e6fdba9b",
			CreatedAt: "2018-06-06T21:48:04",
		},
		Name:             "volume sample",
		Description:      "This is a sample volume for testing",
		Size:             int64(1),
		Status:           "creating",
		AvailabilityZone: "az1",
	}
	var newReq = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id:        "9f1514a7-ca31-4de9-8a27-54c9e6fdba9b",
			CreatedAt: "2018-06-06T21:48:04",
		},
		Name:             "volume sample",
		Description:      "This is a sample volume for testing",
		Size:             int64(1),
		Status:           "creating",
		AvailabilityZone: "az1",
		PoolId:           "f4486139-78d5-462d-a7b9-fdaf6c797e11",
		ProfileId:        "1106b972-66ef-11e7-b172-db03f3689c9c",
	}

	PoolA := model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4486139-78d5-462d-a7b9-fdaf6c797e11",
			CreatedAt: "2017-10-24T15:04:05",
		},
		FreeCapacity:     int64(50),
		AvailabilityZone: "az1",
		Extras: model.StoragePoolExtraSpec{
			Advanced: model.ExtraSpec{
				"thin":     true,
				"dedupe":   true,
				"compress": true,
				"diskType": "SSD",
			},
		},
	}

	pools := []*model.StoragePoolSpec{
		&PoolA,
	}
	mockClient := new(dbtest.MockClient)
	mockClient.On("CreateVolume", context.NewAdminContext(), newReq).Return(nil, errors.New("db error"))
	mockClient.On("GetDock", context.NewAdminContext(), "").Return(&SampleDocks[0], nil)
	mockClient.On("GetDefaultProfile", context.NewAdminContext()).Return(&SampleProfiles[0], nil)
	mockClient.On("GetProfile", context.NewAdminContext(), "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&SampleProfiles[0], nil)
	mockClient.On("ListPools", context.NewAdminContext()).Return(pools, nil)

	db.C = mockClient

	CreateVolumeDBEntry(context.NewAdminContext(), req)
}

func TestCreateVolumeFromSnapshotDBEntry(t *testing.T) {
	PoolA := model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4486139-78d5-462d-a7b9-fdaf6c797e11",
			CreatedAt: "2017-10-24T15:04:05",
		},
		FreeCapacity:     int64(50),
		AvailabilityZone: "az1",
		Extras: model.StoragePoolExtraSpec{
			Advanced: model.ExtraSpec{
				"thin":     true,
				"dedupe":   true,
				"compress": true,
				"diskType": "SSD",
			},
		},
	}

	pools := []*model.StoragePoolSpec{
		&PoolA,
	}
	var req = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id:        "9f1514a7-ca31-4de9-8a27-54c9e6fdba9b",
			CreatedAt: "2018-06-06T21:48:04",
		},
		Name:             "volume sample",
		Description:      "This is a sample volume for testing",
		Size:             int64(1),
		Status:           "creating",
		SnapshotId:       "3769855c-a102-11e7-b772-17b880d2f537",
		AvailabilityZone: "az1",
		PoolId:           "f4486139-78d5-462d-a7b9-fdaf6c797e11",
	}
	var snap = &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "3769855c-a102-11e7-b772-17b880d2f537",
		},
		Size:   int64(1),
		Status: "available",
	}
	var newReq = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id:        "9f1514a7-ca31-4de9-8a27-54c9e6fdba9b",
			CreatedAt: "2018-06-06T21:48:04",
		},
		Name:             "volume sample",
		Description:      "This is a sample volume for testing",
		Size:             int64(1),
		Status:           "creating",
		AvailabilityZone: "az1",
		SnapshotId:       "3769855c-a102-11e7-b772-17b880d2f537",
		PoolId:           "f4486139-78d5-462d-a7b9-fdaf6c797e11",
		ProfileId:        "1106b972-66ef-11e7-b172-db03f3689c9c",
	}

	mockClient := new(dbtest.MockClient)
	mockClient.On("CreateVolume", context.NewAdminContext(), newReq).Return(nil, errors.New("DB error"))
	mockClient.On("GetVolumeSnapshot", context.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f537").Return(snap, nil)
	mockClient.On("GetDock", context.NewAdminContext(), "").Return(&SampleDocks[0], nil)
	mockClient.On("GetDefaultProfile", context.NewAdminContext()).Return(&SampleProfiles[0], nil)
	mockClient.On("GetProfile", context.NewAdminContext(), "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&SampleProfiles[0], nil)
	mockClient.On("GetVolume", context.NewAdminContext(), "").Return(req, nil)
	mockClient.On("ListPools", context.NewAdminContext()).Return(pools, nil)
	db.C = mockClient

	CreateVolumeDBEntry(context.NewAdminContext(), req)
}

func TestDeleteVolumeDBEntry(t *testing.T) {
	var vol = &model.VolumeSpec{
		BaseModel: &model.BaseModel{},
		Status:    "available"}

	mockClient := new(dbtest.MockClient)
	mockClient.On("UpdateVolume", context.NewAdminContext(), vol).Return(nil, nil)
	db.C = mockClient

	err := DeleteVolumeDBEntry(context.NewAdminContext(), vol)
	if err != nil {
		t.Errorf("Failed to delete volume, err is %v\n", err)
	}
}

func TestExtendVolumeDBEntry(t *testing.T) {
	var vol = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
		Status: "available",
		Size:   2,
	}

	mockClient := new(dbtest.MockClient)
	mockClient.On("ExtendVolume", context.NewAdminContext(), vol).Return(nil, nil)
	mockClient.On("GetVolume", context.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(vol, nil)
	db.C = mockClient

	_, err := ExtendVolumeDBEntry(context.NewAdminContext(), vol.Id)
	if err != nil {
		t.Errorf("Failed to delete volume, err is %v\n", err)
	}

	mockClient = new(dbtest.MockClient)
	mockClient.On("GetVolume", context.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(nil, errors.New("error occurs when get volume"))
	db.C = mockClient
	_, err = ExtendVolumeDBEntry(context.NewAdminContext(), vol.Id)

	var vol2 = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
		Status: "error",
		Size:   2,
	}
	mockClient = new(dbtest.MockClient)
	mockClient.On("GetVolume", context.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(vol2, nil)
	db.C = mockClient
	_, err = ExtendVolumeDBEntry(context.NewAdminContext(), vol.Id)
}

func TestCreateVolumeAttachmentDBEntry(t *testing.T) {
	var m = map[string]string{"a": "a"}

	var req = &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{},
		VolumeId:  "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Metadata:  m,
		Status:    "creating",
	}
	var vol = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
		Status: "available",
	}
	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolume", context.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(vol, nil)
	mockClient.On("CreateVolumeAttachment", context.NewAdminContext(), req).Return(&SampleAttachments[0], nil)
	db.C = mockClient

	var expected = &SampleAttachments[0]

	result, err := CreateVolumeAttachmentDBEntry(context.NewAdminContext(), req)

	if err != nil {
		t.Errorf("Failed to create volume attachment, err is %v\n", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestCreateVolumeSnapshotDBEntry(t *testing.T) {
	var m = map[string]string{"a": "a"}
	var vol = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
		Size:   1,
		Status: "available",
	}
	var req = &model.VolumeSnapshotSpec{
		BaseModel:   &model.BaseModel{},
		VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:        "sample-snapshot-01",
		Description: "This is the first sample snapshot for testing",
		Size:        int64(1),
		Status:      "creating",
		Metadata:    m,
	}

	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolume", context.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(vol, nil)
	mockClient.On("CreateVolumeSnapshot", context.NewAdminContext(), req).Return(&SampleSnapshots[0], nil)
	db.C = mockClient

	var expected = &SampleSnapshots[0]
	result, err := CreateVolumeSnapshotDBEntry(context.NewAdminContext(), req)

	if err != nil {
		t.Errorf("Failed to create volume snapshot, err is %v\n", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteVolumeSnapshotDBEntry(t *testing.T) {
	var req = &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "3769855c-a102-11e7-b772-17b880d2f537",
		},
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Status:   "available",
	}

	mockClient := new(dbtest.MockClient)
	mockClient.On("UpdateVolumeSnapshot", context.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f537", req).Return(nil, nil)
	db.C = mockClient

	err := DeleteVolumeSnapshotDBEntry(context.NewAdminContext(), req)

	if err != nil {
		t.Errorf("Failed to delete volume snapshot, err is %v\n", err)
	}
}
