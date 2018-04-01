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
	var req = &model.VolumeSpec{
		BaseModel:   &model.BaseModel{},
		Name:        "sample-volume",
		Description: "This is a sample volume for testing",
		Size:        int64(1),
		ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
	}
	var vol = &SampleVolumes[0]
	mockClient := new(dbtest.MockClient)
	mockClient.On("GetDock", "b7602e18-771e-11e7-8f38-dbd6d291f4e0").Return(&SampleDocks[0], nil)
	mockClient.On("GetDefaultProfile").Return(&SampleProfiles[0], nil)
	mockClient.On("GetProfile", "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&SampleProfiles[0], nil)
	mockClient.On("UpdateVolume", vol.Id, vol).Return(req, nil)
	db.C = mockClient

	var c = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{BaseModel: &model.BaseModel{}, DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0"},
			err: nil,
		},
		volumeController: NewFakeVolumeController(),
	}

	var errchan = make(chan error, 1)
	c.CreateVolume(context.NewAdminContext(), req, errchan)
	if err := <-errchan; err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}
}

func TestDeleteVolume(t *testing.T) {
	var req = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
		ProfileId: "1106b972-66ef-11e7-b172-db03f3689c9c",
		PoolId:    "084bf71e-a102-11e7-88a8-e31fe6d52248",
	}

	mockClient := new(dbtest.MockClient)
	mockClient.On("GetProfile", req.ProfileId).Return(&SampleProfiles[0], nil)
	mockClient.On("GetDockByPoolId", req.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("DeleteVolume", req.Id).Return(nil)
	db.C = mockClient

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
	var vol = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
		PoolId:    "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId: "1106b972-66ef-11e7-b172-db03f3689c9c",
		Size:      int64(1),
	}
	var vol2 = &SampleVolumes[0]
	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolume", vol.Id).Return(vol, nil)
	mockClient.On("GetPool", vol.PoolId).Return(&SamplePools[0], nil)
	mockClient.On("GetDefaultProfile").Return(&SampleProfiles[0], nil)
	mockClient.On("GetProfile", vol.ProfileId).Return(&SampleProfiles[0], nil)
	mockClient.On("GetDockByPoolId", vol.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("UpdateVolume", vol2.Id, vol2).Return(vol, nil)
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
	c.ExtendVolume(context.NewAdminContext(), vol.Id, newSize, errchan)
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
	c.ExtendVolume(context.NewAdminContext(), vol.Id, newSize, errchan2)
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
	c.ExtendVolume(context.NewAdminContext(), vol.Id, newSize, errchan3)

	if err := <-errchan3; err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}
}

func TestCreateVolumeAttachment(t *testing.T) {
	var req = &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{},
		VolumeId:  "bd5b12a8-a101-11e7-941e-d77981b584d8",
		HostInfo:  model.HostInfo{},
		Status:    "creating",
	}
	var vol = &SampleVolumes[0]
	var volattm = &SampleAttachments[0]
	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolume", req.VolumeId).Return(vol, nil)
	mockClient.On("GetDockByPoolId", vol.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("UpdateVolumeAttachment", volattm.Id, volattm).Return(volattm, nil)

	db.C = mockClient

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
	var req = &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id: "f2dda3d2-bf79-11e7-8665-f750b088f63e",
		},
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		HostInfo: model.HostInfo{},
	}
	var vol = &SampleVolumes[0]
	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolume", req.VolumeId).Return(vol, nil)
	mockClient.On("GetDockByPoolId", vol.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("DeleteVolumeAttachment", req.Id).Return(nil)

	db.C = mockClient

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
	var req = &model.VolumeSnapshotSpec{
		BaseModel:   &model.BaseModel{},
		VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:        "sample-snapshot-01",
		Description: "This is the first sample snapshot for testing",
		Size:        int64(1),
	}
	var vol = &SampleVolumes[0]
	var snp = &SampleSnapshots[0]
	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolume", req.VolumeId).Return(vol, nil)
	mockClient.On("GetDockByPoolId", vol.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("UpdateVolumeSnapshot", snp.Id, snp).Return(snp, nil)

	db.C = mockClient

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
	var req = &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "3769855c-a102-11e7-b772-17b880d2f537",
		},
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
	}
	var vol = &SampleVolumes[0]
	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolume", req.VolumeId).Return(vol, nil)
	mockClient.On("GetDockByPoolId", vol.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("DeleteVolumeSnapshot", req.Id).Return(nil)

	db.C = mockClient

	var c = &Controller{
		volumeController: NewFakeVolumeController(),
	}
	var errchan = make(chan error, 1)

	c.DeleteVolumeSnapshot(context.NewAdminContext(), req, errchan)
	if err := <-errchan; err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}
}
