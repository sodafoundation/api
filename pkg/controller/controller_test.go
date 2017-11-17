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

	"github.com/opensds/opensds/pkg/controller/policy"
	"github.com/opensds/opensds/pkg/controller/volume"
	"github.com/opensds/opensds/pkg/db"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
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
	return &sampleVolume, nil
}

func (fvc *fakeVolumeController) DeleteVolume(*pb.DeleteVolumeOpts) error {
	return nil
}

func (fvc *fakeVolumeController) CreateVolumeAttachment(*pb.CreateAttachmentOpts) (*model.VolumeAttachmentSpec, error) {
	return &sampleAttachment, nil
}

func (fvc *fakeVolumeController) DeleteVolumeAttachment(*pb.DeleteAttachmentOpts) error {
	return nil
}

func (fvc *fakeVolumeController) CreateVolumeSnapshot(*pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	return &sampleSnapshot, nil
}

func (fvc *fakeVolumeController) DeleteVolumeSnapshot(*pb.DeleteVolumeSnapshotOpts) error {
	return nil
}

func (fvc *fakeVolumeController) SetDock(dockInfo *model.DockSpec) { return }

func TestCreateVolume(t *testing.T) {

	mockClient := new(dbtest.MockClient)
	mockClient.On("GetDefaultProfile").Return(&sampleProfile, nil)
	mockClient.On("GetProfile", "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&sampleProfile, nil)
	mockClient.On("GetDock", "9193c3ec-771f-11e7-8ca3-d32c0a8b2725").Return(&sampleDock, nil)
	db.C = mockClient

	var req = &model.VolumeSpec{
		BaseModel:   &model.BaseModel{},
		Name:        "fake-volume",
		Description: "fake volume for testing",
		Size:        int64(1),
		ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
	}
	var c = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{BaseModel: &model.BaseModel{}, DockId: "9193c3ec-771f-11e7-8ca3-d32c0a8b2725"},
			err: nil,
		},
		volumeController: NewFakeVolumeController(),
		policyController: policy.NewController(&sampleProfile),
	}
	var expected = &sampleVolume

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
	mockClient.On("GetProfile", "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&sampleProfile, nil)
	mockClient.On("GetDockByPoolId", "71a9e23b-2b46-4331-88bc-c13fcf0cc7b1").Return(&sampleDock, nil)
	db.C = mockClient

	var req = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
		},
		ProfileId: "1106b972-66ef-11e7-b172-db03f3689c9c",
		PoolId:    "71a9e23b-2b46-4331-88bc-c13fcf0cc7b1",
	}
	var c = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{BaseModel: &model.BaseModel{}, DockId: "9193c3ec-771f-11e7-8ca3-d32c0a8b2725"},
			err: nil,
		},
		volumeController: NewFakeVolumeController(),
		policyController: policy.NewController(&sampleProfile),
	}

	result := c.DeleteVolume(req)
	if result != nil {
		t.Errorf("Expected %v, got %v\n", nil, result)
	}
}

func TestCreateVolumeAttachment(t *testing.T) {
	mockClient := new(dbtest.MockClient)
	mockClient.On("GetVolume", "9193c3ec-771f-11e7-8ca3-d32c0a8b2725").Return(&sampleVolume, nil)
	mockClient.On("GetDockByPoolId", "71a9e23b-2b46-4331-88bc-c13fcf0cc7b1").Return(&sampleDock, nil)
	db.C = mockClient
	var req = &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{},
		VolumeId:  "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
		HostInfo:  &model.HostInfo{},
	}
	var c = &Controller{
		volumeController: NewFakeVolumeController(),
	}
	var expected = &sampleAttachment

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
	mockClient.On("GetVolume", "9193c3ec-771f-11e7-8ca3-d32c0a8b2725").Return(&sampleVolume, nil)
	mockClient.On("GetDockByPoolId", "71a9e23b-2b46-4331-88bc-c13fcf0cc7b1").Return(&sampleDock, nil)
	db.C = mockClient

	var req = &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		},
		VolumeId: "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
		HostInfo: &model.HostInfo{},
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
	mockClient.On("GetVolume", "9193c3ec-771f-11e7-8ca3-d32c0a8b2725").Return(&sampleVolume, nil)
	mockClient.On("GetDockByPoolId", "71a9e23b-2b46-4331-88bc-c13fcf0cc7b1").Return(&sampleDock, nil)
	db.C = mockClient

	var req = &model.VolumeSnapshotSpec{
		BaseModel:   &model.BaseModel{},
		VolumeId:    "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
		Name:        "fake-volumesnapshot",
		Description: "fake volumesnapshot for testing",
		Size:        int64(1),
	}
	var c = &Controller{
		volumeController: NewFakeVolumeController(),
	}
	var expected = &sampleSnapshot

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
	mockClient.On("GetVolume", "9193c3ec-771f-11e7-8ca3-d32c0a8b2725").Return(&sampleVolume, nil)
	mockClient.On("GetDockByPoolId", "71a9e23b-2b46-4331-88bc-c13fcf0cc7b1").Return(&sampleDock, nil)
	db.C = mockClient

	var req = &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "8193c3ec-771f-11e7-8ca3-d32c0a8b2727",
		},
		VolumeId: "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}
	var c = &Controller{
		volumeController: NewFakeVolumeController(),
	}

	result := c.DeleteVolumeSnapshot(req)
	if result != nil {
		t.Errorf("Expected %v, got %v\n", nil, result)
	}
}

var (
	sampleDock = model.DockSpec{
		BaseModel: &model.BaseModel{
			Id:        "7ea922da-774b-417a-a226-f8717a4c3cc3",
			CreatedAt: "2017-08-02T09:17:05",
		},
		Name:        "fake-volume",
		Description: "fake volume for testing",
	}

	sampleVolume = model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id:        "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
			CreatedAt: "2017-08-02T09:17:05",
		},
		Name:        "fake-volume",
		Description: "fake volume for testing",
		Size:        1,
		PoolId:      "71a9e23b-2b46-4331-88bc-c13fcf0cc7b1",
	}

	sampleProfile = model.ProfileSpec{
		BaseModel: &model.BaseModel{
			Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		Name:        "default",
		Description: "default policy",
		Extra:       model.ExtraSpec{},
	}

	sampleAttachment = model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		},
		VolumeId: "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}

	sampleModifiedAttachment = model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		},
		VolumeId: "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}

	sampleSnapshot = model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		},
		Name:        "fake-volume-snapshot",
		Description: "fake volume snapshot for testing",
		VolumeId:    "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}
)
