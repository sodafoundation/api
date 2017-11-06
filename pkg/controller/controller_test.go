// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package controller

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/controller/policy"
	"github.com/opensds/opensds/pkg/controller/selector"
	"github.com/opensds/opensds/pkg/controller/volume"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
)

func NewFakeVolumeController() volume.Controller {
	return &fakeVolumeController{}
}

type fakeVolumeController struct {
}

func (fvc *fakeVolumeController) CreateVolume(*pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	return &sampleVolume, nil
}

func (fvc *fakeVolumeController) DeleteVolume(*pb.DeleteVolumeOpts) *model.Response {
	return &model.Response{Status: "Success"}
}

func (fvc *fakeVolumeController) CreateVolumeAttachment(*pb.CreateAttachmentOpts) (*model.VolumeAttachmentSpec, error) {
	return &sampleAttachment, nil
}

func (fvc *fakeVolumeController) DeleteVolumeAttachment(*pb.DeleteAttachmentOpts) *model.Response {
	return &model.Response{Status: "Success"}
}

func (fvc *fakeVolumeController) CreateVolumeSnapshot(*pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	return &sampleSnapshot, nil
}

func (fvc *fakeVolumeController) DeleteVolumeSnapshot(*pb.DeleteVolumeSnapshotOpts) *model.Response {
	return &model.Response{Status: "Success"}
}

func (fvc *fakeVolumeController) SetDock(dockInfo *model.DockSpec) { return }

func TestCreateVolume(t *testing.T) {
	var req = &model.VolumeSpec{
		BaseModel:   &model.BaseModel{},
		Name:        "fake-volume",
		Description: "fake volume for testing",
		Size:        int64(1),
		ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
	}
	var c = &Controller{
		Selector:         selector.NewFakeSelector(),
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
	var req = &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
		},
	}
	var c = &Controller{
		Selector:         selector.NewFakeSelector(),
		volumeController: NewFakeVolumeController(),
		policyController: policy.NewController(&sampleProfile),
	}
	var expected = &model.Response{Status: "Success"}

	result := c.DeleteVolume(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestCreateVolumeAttachment(t *testing.T) {
	var req = &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{},
		VolumeId:  "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
		HostInfo:  &model.HostInfo{},
	}
	var c = &Controller{
		Selector:         selector.NewFakeSelector(),
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
	var req = &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		},
		HostInfo: &model.HostInfo{},
	}
	var c = &Controller{
		Selector:         selector.NewFakeSelector(),
		volumeController: NewFakeVolumeController(),
	}
	var expected = &model.Response{Status: "Success"}

	result := c.DeleteVolumeAttachment(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestCreateVolumeSnapshot(t *testing.T) {
	var req = &model.VolumeSnapshotSpec{
		BaseModel:   &model.BaseModel{},
		VolumeId:    "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
		Name:        "fake-volumesnapshot",
		Description: "fake volumesnapshot for testing",
		Size:        int64(1),
	}
	var c = &Controller{
		Selector:         selector.NewFakeSelector(),
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
	var req = &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "8193c3ec-771f-11e7-8ca3-d32c0a8b2727",
		},
	}
	var c = &Controller{
		Selector:         selector.NewFakeSelector(),
		volumeController: NewFakeVolumeController(),
	}
	var expected = &model.Response{Status: "Success"}

	result := c.DeleteVolumeSnapshot(req)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

var (
	sampleVolume = model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id:        "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
			CreatedAt: "2017-08-02T09:17:05",
		},
		Name:        "fake-volume",
		Description: "fake volume for testing",
		Size:        1,
		PoolId:      "80287bf8-66de-11e7-b031-f3b0af1675ba",
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
		Name:        "fake-volume-attachment",
		Description: "fake volume attachment for testing",
		VolumeId:    "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}

	sampleModifiedAttachment = model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		},
		Name:        "modified-fake-volume-attachment",
		Description: "modified fake volume attachment for testing",
		VolumeId:    "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
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
