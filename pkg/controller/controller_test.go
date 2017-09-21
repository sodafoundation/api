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
	"github.com/opensds/opensds/pkg/controller/volume"
	"github.com/opensds/opensds/pkg/db"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	api "github.com/opensds/opensds/pkg/model"
)

func NewFakeVolumeController(req *pb.DockRequest) volume.Controller {
	return &fakeVolumeController{
		Request: req,
	}
}

type fakeVolumeController struct {
	Request *pb.DockRequest
}

func (fvc *fakeVolumeController) CreateVolume() (*api.VolumeSpec, error) {
	return &sampleVolume, nil
}

func (fvc *fakeVolumeController) DeleteVolume() *api.Response {
	return &api.Response{Status: "Success"}
}

func (fvc *fakeVolumeController) CreateVolumeAttachment() (*api.VolumeAttachmentSpec, error) {
	return &sampleAttachment, nil
}

func (fvc *fakeVolumeController) UpdateVolumeAttachment() (*api.VolumeAttachmentSpec, error) {
	return &sampleModifiedAttachment, nil
}

func (fvc *fakeVolumeController) DeleteVolumeAttachment() *api.Response {
	return &api.Response{Status: "Success"}
}

func (fvc *fakeVolumeController) CreateVolumeSnapshot() (*api.VolumeSnapshotSpec, error) {
	return &sampleSnapshot, nil
}

func (fvc *fakeVolumeController) DeleteVolumeSnapshot() *api.Response {
	return &api.Response{Status: "Success"}
}

func TestNewControllerWithVolumeConfig(t *testing.T) {
	db.C = &fakeDbClient{}

	/*
		CASE 1:
		Test the case that user only specifies volume request.
	*/
	var expectedController = &Controller{
		request: &pb.DockRequest{
			VolumeId:          "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
			VolumeName:        "fake-volume",
			VolumeDescription: "fake volume for testing",
			VolumeSize:        int64(1),
			ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		profile: &api.ProfileSpec{
			BaseModel: &api.BaseModel{
				Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
			},
			Name:        "default",
			Description: "default policy",
			Extra:       api.ExtraSpec{},
		},
	}
	expectedController.searcher = NewDbSearcher()
	expectedController.policyController = policy.NewController(expectedController.profile)
	expectedController.volumeController = volume.NewController(expectedController.request)

	c, err := NewControllerWithVolumeConfig(&sampleVolume, nil, nil)
	if err != nil {
		t.Errorf("Failed to create controller, err is %v\n", err)
	}

	if !reflect.DeepEqual(c, expectedController) {
		t.Errorf("Expected %v, got %v\n", expectedController, c)
	}

	/*
		CASE 2:
		Test the case that user only specifies volume attachment request.
	*/
	expectedController = &Controller{
		request: &pb.DockRequest{
			AttachmentId:          "80287bf8-66de-11e7-b031-f3b0af1675ba",
			AttachmentName:        "fake-volume-attachment",
			AttachmentDescription: "fake volume attachment for testing",
			VolumeId:              "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
		},
	}
	expectedController.searcher = NewDbSearcher()
	expectedController.volumeController = volume.NewController(expectedController.request)

	c, err = NewControllerWithVolumeConfig(nil, &sampleAttachment, nil)
	if err != nil {
		t.Errorf("Failed to create controller, err is %v\n", err)
	}

	if !reflect.DeepEqual(c, expectedController) {
		t.Errorf("Expected %v, got %v\n", expectedController, c)
	}

	/*
		CASE 3:
		Test the case that user only specifies volume snapshot request.
	*/
	expectedController = &Controller{
		request: &pb.DockRequest{
			SnapshotId:          "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			SnapshotName:        "fake-volume-snapshot",
			SnapshotDescription: "fake volume snapshot for testing",
			VolumeId:            "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
		},
	}
	expectedController.searcher = NewDbSearcher()
	expectedController.volumeController = volume.NewController(expectedController.request)

	c, err = NewControllerWithVolumeConfig(nil, nil, &sampleSnapshot)
	if err != nil {
		t.Errorf("Failed to create controller, err is %v\n", err)
	}

	if !reflect.DeepEqual(c, expectedController) {
		t.Errorf("Expected %v, got %v\n", expectedController, c)
	}
}

func TestCreateVolume(t *testing.T) {
	var req = &pb.DockRequest{
		VolumeName:        "fake-volume",
		VolumeDescription: "fake volume for testing",
		VolumeSize:        int64(1),
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
	}
	var c = &Controller{
		searcher:         NewFakeDbSearcher(),
		volumeController: NewFakeVolumeController(req),
		policyController: policy.NewController(&sampleProfile),
		profile:          &sampleProfile,
		request:          req,
	}
	var expected = &sampleVolume

	result, err := c.CreateVolume()
	if err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteVolume(t *testing.T) {
	var req = &pb.DockRequest{
		VolumeId:  "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
		ProfileId: "1106b972-66ef-11e7-b172-db03f3689c9c",
	}
	var c = &Controller{
		searcher:         NewFakeDbSearcher(),
		volumeController: NewFakeVolumeController(req),
		policyController: policy.NewController(&sampleProfile),
		profile:          &sampleProfile,
		request:          req,
	}
	var expected = &api.Response{Status: "Success"}

	result := c.DeleteVolume()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestCreateVolumeAttachment(t *testing.T) {
	var req = &pb.DockRequest{
		AttachmentName:        "fake-volume-attachment",
		AttachmentDescription: "fake volume attachment for testing",
		VolumeId:              "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}
	var c = &Controller{
		searcher:         NewFakeDbSearcher(),
		volumeController: NewFakeVolumeController(req),
		request:          req,
	}
	var expected = &sampleAttachment

	result, err := c.CreateVolumeAttachment()
	if err != nil {
		t.Errorf("Failed to create volume attachment, err is %v\n", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestUpdateVolumeAttachment(t *testing.T) {
	var req = &pb.DockRequest{
		AttachmentId:          "80287bf8-66de-11e7-b031-f3b0af1675ba",
		AttachmentName:        "modified-fake-volume-attachment",
		AttachmentDescription: "modified fake volume attachment for testing",
		VolumeId:              "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}
	var c = &Controller{
		searcher:         NewFakeDbSearcher(),
		volumeController: NewFakeVolumeController(req),
		request:          req,
	}
	var expected = &sampleModifiedAttachment

	result, err := c.UpdateVolumeAttachment()
	if err != nil {
		t.Errorf("Failed to update volume attachment, err is %v\n", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteVolumeAttachment(t *testing.T) {
	var req = &pb.DockRequest{
		AttachmentId: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		VolumeId:     "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}
	var c = &Controller{
		searcher:         NewFakeDbSearcher(),
		volumeController: NewFakeVolumeController(req),
		request:          req,
	}
	var expected = &api.Response{Status: "Success"}

	result := c.DeleteVolumeAttachment()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestCreateVolumeSnapshot(t *testing.T) {
	var req = &pb.DockRequest{
		SnapshotName:        "fake-volume-snapshot",
		SnapshotDescription: "fake volume snapshot for testing",
		VolumeId:            "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}
	var c = &Controller{
		searcher:         NewFakeDbSearcher(),
		volumeController: NewFakeVolumeController(req),
		request:          req,
	}
	var expected = &sampleSnapshot

	result, err := c.CreateVolumeSnapshot()
	if err != nil {
		t.Errorf("Failed to create volume snapshot, err is %v\n", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteVolumeSnapshot(t *testing.T) {
	var req = &pb.DockRequest{
		SnapshotId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		VolumeId:   "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}
	var c = &Controller{
		searcher:         NewFakeDbSearcher(),
		volumeController: NewFakeVolumeController(req),
		request:          req,
	}
	var expected = &api.Response{Status: "Success"}

	result := c.DeleteVolumeSnapshot()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

var (
	sampleProfile = api.ProfileSpec{
		BaseModel: &api.BaseModel{
			Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		Name:        "default",
		Description: "default policy",
		Extra:       api.ExtraSpec{},
	}

	sampleAttachment = api.VolumeAttachmentSpec{
		BaseModel: &api.BaseModel{
			Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		},
		Name:        "fake-volume-attachment",
		Description: "fake volume attachment for testing",
		VolumeId:    "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}

	sampleModifiedAttachment = api.VolumeAttachmentSpec{
		BaseModel: &api.BaseModel{
			Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		},
		Name:        "modified-fake-volume-attachment",
		Description: "modified fake volume attachment for testing",
		VolumeId:    "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}

	sampleSnapshot = api.VolumeSnapshotSpec{
		BaseModel: &api.BaseModel{
			Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		},
		Name:        "fake-volume-snapshot",
		Description: "fake volume snapshot for testing",
		VolumeId:    "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}
)
