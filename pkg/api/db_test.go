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
		BaseModel:   &model.BaseModel{},
		Name:        "volume sample",
		Description: "This is a sample volume for testing",
		Size:        int64(1),
		Status:      "creating",
	}

	mockClient := new(dbtest.Client)
	mockClient.On("CreateVolume", context.NewAdminContext(), req).Return(&SampleVolumes[0], nil)
	db.C = mockClient

	var expected = &SampleVolumes[0]
	result, err := CreateVolumeDBEntry(context.NewAdminContext(), req)

	if err != nil {
		t.Errorf("Failed to create volume asynchronously, err is %v\n", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}

	var req2 = &model.VolumeSpec{
		BaseModel:   &model.BaseModel{},
		Name:        "volume sample",
		Description: "This is a sample volume for testing",
		Size:        int64(-2),
		Status:      "creating",
	}
	result, err = CreateVolumeDBEntry(context.NewAdminContext(), req2)

	mockClient = new(dbtest.Client)
	mockClient.On("CreateVolume", context.NewAdminContext(), req).Return(nil, errors.New("not find the volume"))
	db.C = mockClient
	result, err = CreateVolumeDBEntry(context.NewAdminContext(), req)
}

func TestCreateVolumeFromSnapshotDBEntry(t *testing.T) {
	var req = &model.VolumeSpec{
		BaseModel:   &model.BaseModel{},
		Name:        "volume sample",
		Description: "This is a sample volume for testing",
		Size:        int64(1),
		Status:      "creating",
		SnapshotId:  "3769855c-a102-11e7-b772-17b880d2f537",
	}
	var snap = &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "3769855c-a102-11e7-b772-17b880d2f537",
		},
		Size:   int64(1),
		Status: "available",
	}

	mockClient := new(dbtest.Client)
	mockClient.On("CreateVolume", context.NewAdminContext(), req).Return(&SampleVolumes[1], nil)
	mockClient.On("GetVolumeSnapshot", context.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f537").Return(snap, nil)
	db.C = mockClient

	var expected = &SampleVolumes[1]
	result, err := CreateVolumeDBEntry(context.NewAdminContext(), req)

	if err != nil {
		t.Errorf("Failed to create volume with snapshot, err is %v\n", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteVolumeDBEntry(t *testing.T) {
	var vol = &model.VolumeSpec{
		BaseModel: &model.BaseModel{},
		Status:    "available",
		ProfileId: "3769855c-a102-11e7-b772-17b880d2f537",
		PoolId:    "3762355c-a102-11e7-b772-17b880d2f537",
	}

	mockClient := new(dbtest.Client)
	mockClient.On("UpdateVolume", context.NewAdminContext(), vol).Return(nil, nil)
	mockClient.On("DeleteVolume", context.NewAdminContext(), vol.Id).Return(nil)
	mockClient.On("ListSnapshotsByVolumeId", context.NewAdminContext(), vol.Id).Return(nil, nil)
	mockClient.On("ListVolumeAttachments", context.NewAdminContext(), vol.Id).Return(nil, nil)
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

	mockClient := new(dbtest.Client)
	mockClient.On("ExtendVolume", context.NewAdminContext(), vol).Return(nil, nil)
	mockClient.On("GetVolume", context.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(vol, nil)
	db.C = mockClient

	_, err := ExtendVolumeDBEntry(context.NewAdminContext(), vol.Id)
	if err != nil {
		t.Errorf("Failed to delete volume, err is %v\n", err)
	}

	mockClient = new(dbtest.Client)
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
	mockClient = new(dbtest.Client)
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
	mockClient := new(dbtest.Client)
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

	mockClient := new(dbtest.Client)
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

	mockClient := new(dbtest.Client)
	mockClient.On("UpdateVolumeSnapshot", context.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f537", req).Return(nil, nil)
	mockClient.On("GetVolume", context.NewAdminContext(), req.VolumeId).Return(nil, nil)
	db.C = mockClient

	err := DeleteVolumeSnapshotDBEntry(context.NewAdminContext(), req)

	if err != nil {
		t.Errorf("Failed to delete volume snapshot, err is %v\n", err)
	}
}
