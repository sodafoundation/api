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

package controller

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller/volume"
	"github.com/opensds/opensds/pkg/db"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/testutils/collection"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
)

type fakeSelector struct {
	res *model.StoragePoolSpec
	err error
}

func (s *fakeSelector) SelectSupportedPool(tags map[string]interface{}) (*model.StoragePoolSpec, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.res, nil
}

func NewFakeVolumeController() volume.Controller {
	return &fakeVolumeController{}
}

type fakeVolumeController struct {
}

func (fvc *fakeVolumeController) CreateVolume(*pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	return &SampleVolumes[0], nil
}

func (fvc *fakeVolumeController) DeleteVolume(*pb.DeleteVolumeOpts) error {
	return nil
}

func (fvc *fakeVolumeController) ExtendVolume(*pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	return &SampleVolumes[0], nil
}

func (fvc *fakeVolumeController) CreateVolumeAttachment(*pb.CreateAttachmentOpts) (*model.VolumeAttachmentSpec, error) {
	return &SampleAttachments[0], nil
}

func (fvc *fakeVolumeController) DeleteVolumeAttachment(*pb.DeleteAttachmentOpts) error {
	return nil
}

func (fvc *fakeVolumeController) CreateVolumeSnapshot(*pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	return &SampleSnapshots[0], nil
}

func (fvc *fakeVolumeController) DeleteVolumeSnapshot(*pb.DeleteVolumeSnapshotOpts) error {
	return nil
}

func (fvc *fakeVolumeController) SetDock(dockInfo *model.DockSpec) { return }

func TestCreateVolume(t *testing.T) {
	mockClient := new(dbtest.MockClient)
	mockClient.On("GetDock", "b7602e18-771e-11e7-8f38-dbd6d291f4e0").Return(&SampleDocks[0], nil)
	db.C = mockClient

	var req = &model.VolumeSpec{
		BaseModel:   &model.BaseModel{},
		Name:        "sample-volume",
		Description: "This is a sample volume for testing",
		Size:        int64(1),
		ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
	}
	var c = &Controller{
		volumeController: NewFakeVolumeController(),
	}

	var pol = &model.StoragePoolSpec{
		BaseModel: &model.BaseModel{},
		DockId:    "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
	}

	var errchan = make(chan error, 1)
	c.CreateVolume(context.NewAdminContext(), req, pol, &SampleProfiles[0], errchan)
	if err := <-errchan; err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}
}

func TestDeleteVolume(t *testing.T) {
	mockClient := new(dbtest.MockClient)
	mockClient.On("GetProfile", "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&SampleProfiles[0], nil)
	mockClient.On("GetDockByPoolId", "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&SampleDocks[0], nil)
	db.C = mockClient

	var req = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
		ProfileId: "1106b972-66ef-11e7-b172-db03f3689c9c",
		PoolId:    "084bf71e-a102-11e7-88a8-e31fe6d52248",
	}
	var c = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{BaseModel: &model.BaseModel{}, DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0"},
			err: nil,
		},
		volumeController: NewFakeVolumeController(),
	}
	var errchan = make(chan error, 1)
	c.DeleteVolume(context.NewAdminContext(), req, errchan)

	if err := <-errchan; err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}
}

func TestExtendVolume(t *testing.T) {
	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolume", "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
	mockClient.On("GetPool", "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&SamplePools[0], nil)
	mockClient.On("GetDefaultProfile").Return(&SampleProfiles[0], nil)
	mockClient.On("GetProfile", "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&SampleProfiles[0], nil)
	mockClient.On("GetDockByPoolId", "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&SampleDocks[0], nil)
	mockClient.On("GetDock", "b7602e18-771e-11e7-8f38-dbd6d291f4e0").Return(&SampleDocks[0], nil)
	db.C = mockClient

	var c = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{BaseModel: &model.BaseModel{}, DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0"},
			err: nil,
		},
		volumeController: NewFakeVolumeController(),
	}

	newSize := int64(1)
	var errchan = make(chan error, 1)
	c.ExtendVolume(context.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8", newSize, errchan)
	expectedError := "new size(1) <= old size(1)"

	if err := <-errchan; err == nil {
		t.Errorf("Expected Non-%v, got %v\n", nil, err)
	} else {
		if expectedError != err.Error() {
			t.Errorf("Expected Non-%v, got %v\n", expectedError, err.Error())
		}
	}

	newSize = int64(92)
	var errchan2 = make(chan error, 1)
	c.ExtendVolume(context.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8", newSize, errchan2)
	expectedError = "pool free capacity(90) < new size(92) - old size(1)"

	if err := <-errchan2; err == nil {

		t.Errorf("Expected Non-%v, got %v\n", nil, err)
	} else {
		if expectedError != err.Error() {
			t.Errorf("Expected Non-%v, got %v\n", expectedError, err.Error())
		}
	}

	newSize = int64(2)
	var errchan3 = make(chan error, 1)
	c.ExtendVolume(context.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8", newSize, errchan3)

	if err := <-errchan3; err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}
}

func TestCreateVolumeAttachment(t *testing.T) {
	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolume", "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
	mockClient.On("GetDockByPoolId", "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&SampleDocks[0], nil)
	db.C = mockClient
	var req = &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{},
		VolumeId:  "bd5b12a8-a101-11e7-941e-d77981b584d8",
		HostInfo:  model.HostInfo{},
		Status:    "creating",
	}
	var c = &Controller{
		volumeController: NewFakeVolumeController(),
	}

	var errchan = make(chan error, 1)

	c.CreateVolumeAttachment(context.NewAdminContext(), req, errchan)
	if err := <-errchan; err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}
}

func TestDeleteVolumeAttachment(t *testing.T) {
	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolume", "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
	mockClient.On("GetDockByPoolId", "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&SampleDocks[0], nil)
	db.C = mockClient

	var req = &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id: "f2dda3d2-bf79-11e7-8665-f750b088f63e",
		},
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		HostInfo: model.HostInfo{},
	}
	var c = &Controller{
		volumeController: NewFakeVolumeController(),
	}
	var errchan = make(chan error, 1)

	c.DeleteVolumeAttachment(context.NewAdminContext(), req, errchan)

	if err := <-errchan; err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}
}

func TestCreateVolumeSnapshot(t *testing.T) {
	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolume", "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
	mockClient.On("GetDockByPoolId", "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&SampleDocks[0], nil)
	db.C = mockClient

	var req = &model.VolumeSnapshotSpec{
		BaseModel:   &model.BaseModel{},
		VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:        "sample-snapshot-01",
		Description: "This is the first sample snapshot for testing",
		Size:        int64(1),
	}
	var c = &Controller{
		volumeController: NewFakeVolumeController(),
	}

	var errchan = make(chan error, 1)

	c.CreateVolumeSnapshot(context.NewAdminContext(), req, errchan)
	if err := <-errchan; err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}
}

func TestDeleteVolumeSnapshot(t *testing.T) {
	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolume", "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
	mockClient.On("GetDockByPoolId", "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&SampleDocks[0], nil)
	db.C = mockClient

	var req = &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "3769855c-a102-11e7-b772-17b880d2f537",
		},
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
	}
	var c = &Controller{
		volumeController: NewFakeVolumeController(),
	}
	var errchan = make(chan error, 1)

	c.DeleteVolumeSnapshot(context.NewAdminContext(), req, errchan)
	if err := <-errchan; err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}
}

func TestCreateVolumeDBEntry(t *testing.T) {
	var req = &model.VolumeSpec{
		BaseModel:   &model.BaseModel{},
		Name:        "volume sample",
		Description: "This is a sample volume for testing",
		Size:        int64(1),
		ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
		Status:      "creating",
	}

	mockClient := new(dbtest.MockClient)
	mockClient.On("GetDefaultProfile").Return(&SampleProfiles[0], nil)
	mockClient.On("GetProfile", "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&SampleProfiles[0], nil)
	mockClient.On("CreateVolume", req).Return(&SampleVolumes[0], nil)
	db.C = mockClient

	var c = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{BaseModel: &model.BaseModel{}, DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0"},
			err: nil,
		},
		volumeController: NewFakeVolumeController(),
	}
	var expected = &SampleVolumes[0]
	result, _, _, err := c.CreateVolumeDBEntry(context.NewAdminContext(), req)

	if err != nil {
		t.Errorf("Failed to create volume asynchronously, err is %v\n", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteVolumeDBEntry(t *testing.T) {
	var vol = &model.VolumeSpec{
		BaseModel: &model.BaseModel{},
		Status:    "available"}

	mockClient := new(dbtest.MockClient)
	mockClient.On("UpdateVolume", "", vol).Return(nil, nil)
	db.C = mockClient

	var c = &Controller{
		volumeController: NewFakeVolumeController(),
	}

	err := c.DeleteVolumeDBEntry(context.NewAdminContext(), vol)
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
	mockClient.On("ExtendVolume", vol).Return(nil, nil)
	mockClient.On("GetVolume", "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(vol, nil)
	db.C = mockClient

	var c = &Controller{
		volumeController: NewFakeVolumeController(),
	}

	_, err := c.ExtendVolumeDBEntry(context.NewAdminContext(), vol.Id)
	if err != nil {
		t.Errorf("Failed to delete volume, err is %v\n", err)
	}
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
	mockClient.On("GetVolume", "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(vol, nil)
	mockClient.On("CreateVolumeAttachment", req).Return(&SampleAttachments[0], nil)
	db.C = mockClient

	var c = &Controller{
		volumeController: NewFakeVolumeController(),
	}
	var expected = &SampleAttachments[0]

	result, err := c.CreateVolumeAttachmentDBEntry(context.NewAdminContext(), req)

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
	mockClient.On("GetVolume", "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(vol, nil)
	mockClient.On("CreateVolumeSnapshot", req).Return(&SampleSnapshots[0], nil)
	db.C = mockClient

	var c = &Controller{
		volumeController: NewFakeVolumeController(),
	}
	var expected = &SampleSnapshots[0]
	result, err := c.CreateVolumeSnapshotDBEntry(context.NewAdminContext(), req)

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
	mockClient.On("UpdateVolumeSnapshot", "3769855c-a102-11e7-b772-17b880d2f537", req).Return(nil, nil)
	db.C = mockClient

	var c = &Controller{
		volumeController: NewFakeVolumeController(),
	}
	err := c.DeleteVolumeSnapshotDBEntry(context.NewAdminContext(), req)

	if err != nil {
		t.Errorf("Failed to delete volume snapshot, err is %v\n", err)
	}
}
