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

package controller

import (
	"reflect"
	"testing"

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
	mockClient.On("GetDefaultProfile").Return(&SampleProfiles[0], nil)
	mockClient.On("GetProfile", "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&SampleProfiles[0], nil)
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
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{BaseModel: &model.BaseModel{}, DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0"},
			err: nil,
		},
		volumeController: NewFakeVolumeController(),
	}
	var expected = &SampleVolumes[0]

	result, err := c.CreateVolume(req)
	if err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
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

	result := c.DeleteVolume(req)
	if result != nil {
		t.Errorf("Expected %v, got %v\n", nil, result)
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
	}
	var c = &Controller{
		volumeController: NewFakeVolumeController(),
	}
	var expected = &SampleAttachments[0]

	result, err := c.CreateVolumeAttachment(req)
	if err != nil {
		t.Errorf("Failed to create volume attachment, err is %v\n", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
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

	result := c.DeleteVolumeAttachment(req)
	if result != nil {
		t.Errorf("Expected %v, got %v\n", nil, result)
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
	var expected = &SampleSnapshots[0]

	result, err := c.CreateVolumeSnapshot(req)
	if err != nil {
		t.Errorf("Failed to create volume snapshot, err is %v\n", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
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

	result := c.DeleteVolumeSnapshot(req)
	if result != nil {
		t.Errorf("Expected %v, got %v\n", nil, result)
	}
}
