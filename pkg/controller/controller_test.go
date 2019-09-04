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
	"context"
	"errors"
	"fmt"
	"testing"

	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller/dr"
	"github.com/opensds/opensds/pkg/controller/volume"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	. "github.com/opensds/opensds/testutils/collection"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
	"github.com/opensds/opensds/pkg/controller/fileshare"
)

type fakeSelector struct {
	res *model.StoragePoolSpec
	err error
}

func (s *fakeSelector) SelectSupportedPoolForVolume(vol *model.VolumeSpec) (*model.StoragePoolSpec, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.res, nil
}

func (s *fakeSelector) SelectSupportedPoolForFileShare(vol *model.FileShareSpec) (*model.StoragePoolSpec, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.res, nil
}

func (s *fakeSelector) SelectSupportedPoolForVG(vg *model.VolumeGroupSpec) (*model.StoragePoolSpec, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.res, nil
}

// NewController method creates a controller structure and expose its pointer.
func NewFakeDrController() dr.Controller {
	return &fakeDrController{}
}

type fakeDrController struct {
}

func (d *fakeDrController) CreateReplication(ctx *c.Context, replica *model.ReplicationSpec, primaryVol,
	secondaryVol *model.VolumeSpec) (*model.ReplicationSpec, error) {
	return &SampleReplications[0], nil
}

func (d *fakeDrController) DeleteReplication(ctx *c.Context, replica *model.ReplicationSpec, primaryVol,
	secondaryVol *model.VolumeSpec) error {
	return nil
}

func (d *fakeDrController) EnableReplication(ctx *c.Context, replica *model.ReplicationSpec, primaryVol,
	secondaryVol *model.VolumeSpec) error {
	return nil
}

func (d *fakeDrController) DisableReplication(ctx *c.Context, replica *model.ReplicationSpec, primaryVol,
	secondaryVol *model.VolumeSpec) error {
	return nil
}

func (d *fakeDrController) FailoverReplication(ctx *c.Context, replica *model.ReplicationSpec,
	failover *model.FailoverReplicationSpec, primaryVol, secondaryVol *model.VolumeSpec) error {
	return nil
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

func (fvc *fakeVolumeController) CreateVolumeAttachment(*pb.CreateVolumeAttachmentOpts) (*model.VolumeAttachmentSpec, error) {
	return &SampleAttachments[0], nil
}

func (fvc *fakeVolumeController) DeleteVolumeAttachment(*pb.DeleteVolumeAttachmentOpts) error {
	return nil
}

func (fvc *fakeVolumeController) CreateVolumeSnapshot(*pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	return &SampleSnapshots[0], nil
}

func (fvc *fakeVolumeController) DeleteVolumeSnapshot(*pb.DeleteVolumeSnapshotOpts) error {
	return nil
}

func (fvc *fakeVolumeController) AttachVolume(*pb.AttachVolumeOpts) (string, error) {
	return "", nil
}

func (fvc *fakeVolumeController) DetachVolume(*pb.DetachVolumeOpts) error {
	return nil
}

func (fvc *fakeVolumeController) CreateReplication(opts *pb.CreateReplicationOpts) (*model.ReplicationSpec, error) {
	return &SampleReplications[0], nil
}

func (fvc *fakeVolumeController) DeleteReplication(opt *pb.DeleteReplicationOpts) error {
	return nil
}

func (fvc *fakeVolumeController) EnableReplication(opt *pb.EnableReplicationOpts) error {
	return nil
}

func (fvc *fakeVolumeController) DisableReplication(opt *pb.DisableReplicationOpts) error {
	return nil
}

func (fvc *fakeVolumeController) FailoverReplication(opt *pb.FailoverReplicationOpts) error {
	return nil
}

func (fvc *fakeVolumeController) CreateVolumeGroup(*pb.CreateVolumeGroupOpts) (*model.VolumeGroupSpec, error) {
	return &SampleVolumeGroups[0], nil
}

func (fvc *fakeVolumeController) UpdateVolumeGroup(*pb.UpdateVolumeGroupOpts) (*model.VolumeGroupSpec, error) {
	return &SampleVolumeGroups[0], nil
}

func (fvc *fakeVolumeController) DeleteVolumeGroup(*pb.DeleteVolumeGroupOpts) error {
	return nil
}
func (fvc *fakeVolumeController) SetDock(dockInfo *model.DockSpec) { return }

func TestCreateVolume(t *testing.T) {
	var req = &pb.CreateVolumeOpts{
		Id:          "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:        "sample-volume",
		Description: "This is a sample volume for testing",
		Size:        int64(1),
		ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
		Context:     c.NewAdminContext().ToJson(),
	}
	var vol = &SampleVolumes[0]
	mockClient := new(dbtest.Client)
	mockClient.On("GetVolume", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
	mockClient.On("GetDock", c.NewAdminContext(), "b7602e18-771e-11e7-8f38-dbd6d291f4e0").Return(&SampleDocks[0], nil)
	mockClient.On("GetDefaultProfile", c.NewAdminContext()).Return(&SampleProfiles[0], nil)
	mockClient.On("GetProfile", c.NewAdminContext(), "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&SampleProfiles[0], nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), vol, vol.Status).Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{
				BaseModel: &model.BaseModel{
					Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
				},
				DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			err: nil,
		},
		volumeController: NewFakeVolumeController(),
	}

	if _, err := ctrl.CreateVolume(context.Background(), req); err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}
}

func TestCreateVolumeFromSnapshot(t *testing.T) {
	var req = &pb.CreateVolumeOpts{
		Id:          "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:        "sample-volume",
		Description: "This is a sample volume for testing",
		Size:        int64(1),
		ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
		SnapshotId:  "3769855c-a102-11e7-b772-17b880d2f537",
		Context:     c.NewAdminContext().ToJson(),
	}
	var vol = &SampleVolumes[0]
	mockClient := new(dbtest.Client)
	mockClient.On("GetDock", c.NewAdminContext(), "b7602e18-771e-11e7-8f38-dbd6d291f4e0").Return(&SampleDocks[0], nil)
	mockClient.On("GetVolumeSnapshot", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f537").Return(&SampleSnapshots[0], nil)
	mockClient.On("GetVolume", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
	mockClient.On("GetProfile", c.NewAdminContext(), "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&SampleProfiles[0], nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), vol, vol.Status).Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{
				BaseModel: &model.BaseModel{
					Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
				},
				DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			err: nil,
		},
		volumeController: NewFakeVolumeController(),
	}

	if _, err := ctrl.CreateVolume(context.Background(), req); err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}
}

func TestDeleteVolume(t *testing.T) {
	var req = &pb.DeleteVolumeOpts{
		Id:        "bd5b12a8-a101-11e7-941e-d77981b584d8",
		ProfileId: "1106b972-66ef-11e7-b172-db03f3689c9c",
		PoolId:    "084bf71e-a102-11e7-88a8-e31fe6d52248",
		Context:   c.NewAdminContext().ToJson(),
	}

	mockClient := new(dbtest.Client)
	mockClient.On("GetProfile", c.NewAdminContext(), req.ProfileId).Return(&SampleProfiles[0], nil)
	mockClient.On("GetDockByPoolId", c.NewAdminContext(), req.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("DeleteVolume", c.NewAdminContext(), req.Id).Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{
				BaseModel: &model.BaseModel{
					Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
				},
				DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			err: nil,
		},
		volumeController: NewFakeVolumeController(),
	}

	if _, err := ctrl.DeleteVolume(context.Background(), req); err != nil {
		t.Errorf("Failed to delete volume, err is %v\n", err)
	}
}

func TestExtendVolume(t *testing.T) {
	var req = &pb.ExtendVolumeOpts{
		Id:        "bd5b12a8-a101-11e7-941e-d77981b584d8",
		PoolId:    "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId: "1106b972-66ef-11e7-b172-db03f3689c9c",
		Size:      int64(1),
		Context:   c.NewAdminContext().ToJson(),
	}
	var vol2 = &SampleVolumes[0]
	mockClient := new(dbtest.Client)
	mockClient.On("GetVolume", c.NewAdminContext(), req.Id).Return(vol2, nil)
	mockClient.On("GetPool", c.NewAdminContext(), req.PoolId).Return(&SamplePools[0], nil)
	mockClient.On("GetDefaultProfile", c.NewAdminContext()).Return(&SampleProfiles[0], nil)
	mockClient.On("GetProfile", c.NewAdminContext(), req.ProfileId).Return(&SampleProfiles[0], nil)
	mockClient.On("GetDockByPoolId", c.NewAdminContext(), req.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), vol2, vol2.Status).Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{
				BaseModel: &model.BaseModel{
					Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
				},
				DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			err: nil,
		},
		volumeController: NewFakeVolumeController(),
	}

	req.Size = int64(92)
	_, err := ctrl.ExtendVolume(context.Background(), req)
	expectedError := "pool free capacity(90) < new size(92) - old size(1)"
	if err == nil {
		t.Errorf("Expected Non-%v, got %v\n", nil, err)
	} else {
		if expectedError != err.Error() {
			t.Errorf("Expected Non-%v, got %v\n", expectedError, err.Error())
		}
	}

	req.Size = int64(2)
	if _, err = ctrl.ExtendVolume(context.Background(), req); err != nil {
		t.Errorf("Failed to extend volume: %v\n", err)
	}
}

func TestCreateVolumeAttachment(t *testing.T) {
	var req = &pb.CreateVolumeAttachmentOpts{
		Id:       "f2dda3d2-bf79-11e7-8665-f750b088f63e",
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		HostInfo: &pb.HostInfo{},
		Context:  c.NewAdminContext().ToJson(),
	}
	var vol, volatm = &SampleVolumes[0], &SampleAttachments[0]
	mockClient := new(dbtest.Client)
	mockClient.On("GetVolume", c.NewAdminContext(), req.VolumeId).Return(vol, nil)
	mockClient.On("GetPool", c.NewAdminContext(), vol.PoolId).Return(&SamplePools[0], nil)
	mockClient.On("GetDock", c.NewAdminContext(), "b7602e18-771e-11e7-8f38-dbd6d291f4e0").Return(&SampleDocks[0], nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), volatm, volatm.Status).Return(nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), vol, model.VolumeInUse).Return(nil)

	db.C = mockClient

	var ctrl = &Controller{
		volumeController: NewFakeVolumeController(),
	}

	if _, err := ctrl.CreateVolumeAttachment(context.Background(), req); err != nil {
		t.Errorf("Failed to create volume attachment: %v\n", err)
	}
}

func TestDeleteVolumeAttachment(t *testing.T) {
	var req = &pb.DeleteVolumeAttachmentOpts{
		Id:       "f2dda3d2-bf79-11e7-8665-f750b088f63e",
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		HostInfo: &pb.HostInfo{},
		Context:  c.NewAdminContext().ToJson(),
	}
	var vol = &SampleVolumes[0]
	mockClient := new(dbtest.Client)
	mockClient.On("GetVolume", c.NewAdminContext(), req.VolumeId).Return(vol, nil)
	mockClient.On("GetDockByPoolId", c.NewAdminContext(), vol.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("DeleteVolumeAttachment", c.NewAdminContext(), req.Id).Return(nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), vol, model.VolumeAvailable).Return(nil)

	db.C = mockClient

	var ctrl = &Controller{
		volumeController: NewFakeVolumeController(),
	}

	if _, err := ctrl.DeleteVolumeAttachment(context.Background(), req); err != nil {
		t.Errorf("Failed to delete volume attachment: %v\n", err)
	}
}

func TestCreateVolumeSnapshot(t *testing.T) {
	var req = &pb.CreateVolumeSnapshotOpts{
		Id:          "3769855c-a102-11e7-b772-17b880d2f537",
		VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:        "sample-snapshot-01",
		Description: "This is the first sample snapshot for testing",
		Size:        int64(1),
		Context:     c.NewAdminContext().ToJson(),
	}
	var vol = &SampleVolumes[0]
	var snp = &SampleSnapshots[0]
	mockClient := new(dbtest.Client)
	mockClient.On("GetVolume", c.NewAdminContext(), req.VolumeId).Return(vol, nil)
	mockClient.On("GetDockByPoolId", c.NewAdminContext(), vol.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), snp, "available").Return(nil)

	db.C = mockClient

	var ctrl = &Controller{
		volumeController: NewFakeVolumeController(),
	}

	if _, err := ctrl.CreateVolumeSnapshot(context.Background(), req); err != nil {
		t.Errorf("Failed to create volume snapshot: %v\n", err)
	}
}

func TestDeleteVolumeSnapshot(t *testing.T) {
	var req = &pb.DeleteVolumeSnapshotOpts{
		Id:       "3769855c-a102-11e7-b772-17b880d2f537",
		VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Context:  c.NewAdminContext().ToJson(),
	}
	var vol = &SampleVolumes[0]
	mockClient := new(dbtest.Client)
	mockClient.On("GetVolume", c.NewAdminContext(), req.VolumeId).Return(vol, nil)
	mockClient.On("GetDockByPoolId", c.NewAdminContext(), vol.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("DeleteVolumeSnapshot", c.NewAdminContext(), req.Id).Return(nil)

	db.C = mockClient

	var ctrl = &Controller{
		volumeController: NewFakeVolumeController(),
	}

	if _, err := ctrl.DeleteVolumeSnapshot(context.Background(), req); err != nil {
		t.Errorf("Failed to delete volume snapshot: %v\n", err)
	}
}

func TestCreateReplication(t *testing.T) {
	var req = &pb.CreateReplicationOpts{
		Id:              "c299a978-4f3e-11e8-8a5c-977218a83359",
		PrimaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		// Just adapt the mock method,the volume must be different in real scenario.
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
		Context:           c.NewAdminContext().ToJson(),
	}

	var replica = &SampleReplications[0]
	mockClient := new(dbtest.Client)
	mockClient.On("GetReplication", c.NewAdminContext(), req.Id).Return(replica, nil)
	mockClient.On("GetVolume", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), replica, model.ReplicationAvailable).Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		drController: NewFakeDrController(),
	}

	if _, err := ctrl.CreateReplication(context.Background(), req); err != nil {
		t.Errorf("Failed to create volume replication: %v\n", err)
	}
}

func TestDeleteReplication(t *testing.T) {
	var req = &pb.DeleteReplicationOpts{
		Id:              "c299a978-4f3e-11e8-8a5c-977218a83359",
		PrimaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		// Just adapt the mock method,the volume must be different in real scenario.
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
		Context:           c.NewAdminContext().ToJson(),
	}

	mockClient := new(dbtest.Client)
	mockClient.On("GetReplication", c.NewAdminContext(), req.Id).Return(&SampleReplications[0], nil)
	mockClient.On("GetVolume", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
	mockClient.On("DeleteReplication", c.NewAdminContext(), req.Id).Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		drController: NewFakeDrController(),
	}

	if _, err := ctrl.DeleteReplication(context.Background(), req); err != nil {
		t.Errorf("Failed to delete volume replication: %v\n", err)
	}
}

func TestEnableReplication(t *testing.T) {
	var req = &pb.EnableReplicationOpts{
		Id:              "c299a978-4f3e-11e8-8a5c-977218a83359",
		PrimaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		// Just adapt the mock method,the volume must be different in real scenario.
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
		Context:           c.NewAdminContext().ToJson(),
	}

	mockClient := new(dbtest.Client)
	mockClient.On("GetReplication", c.NewAdminContext(), req.Id).Return(&SampleReplications[0], nil)
	mockClient.On("GetVolume", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), &SampleReplications[0], model.ReplicationEnabled).Return(nil)
	db.C = mockClient
	var ctrl = &Controller{
		drController: NewFakeDrController(),
	}

	if _, err := ctrl.EnableReplication(context.Background(), req); err != nil {
		t.Errorf("Failed to enable volume replication: %v\n", err)
	}
}

func TestDisableReplication(t *testing.T) {
	var req = &pb.DisableReplicationOpts{
		Id:              "c299a978-4f3e-11e8-8a5c-977218a83359",
		PrimaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		// Just adapt the mock method,the volume must be different in real scenario.
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
		Context:           c.NewAdminContext().ToJson(),
	}

	mockClient := new(dbtest.Client)
	mockClient.On("GetReplication", c.NewAdminContext(), req.Id).Return(&SampleReplications[0], nil)
	mockClient.On("GetVolume", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), &SampleReplications[0], model.ReplicationDisabled).Return(nil)
	db.C = mockClient
	var ctrl = &Controller{
		drController: NewFakeDrController(),
	}

	if _, err := ctrl.DisableReplication(context.Background(), req); err != nil {
		t.Errorf("Failed to disable volume replication: %v\n", err)
	}
}

func TestFailoverReplication(t *testing.T) {
	var req = &pb.FailoverReplicationOpts{
		Id:              "c299a978-4f3e-11e8-8a5c-977218a83359",
		PrimaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		// Just adapt the mock method,the volume must be different in real scenario.
		SecondaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:                "sample-replication-01",
		Description:         "This is a sample replication for testing",
		PoolId:              "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:           "1106b972-66ef-11e7-b172-db03f3689c9c",
		AllowAttachedVolume: true,
		SecondaryBackendId:  model.ReplicationDefaultBackendId,
		Context:             c.NewAdminContext().ToJson(),
	}

	mockClient := new(dbtest.Client)
	mockClient.On("GetReplication", c.NewAdminContext(), req.Id).Return(&SampleReplications[0], nil)
	mockClient.On("GetVolume", c.NewAdminContext(), "bd5b12a8-a101-11e7-941e-d77981b584d8").Return(&SampleVolumes[0], nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), &SampleReplications[0], model.ReplicationFailover).Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		volumeController: NewFakeVolumeController(),
		drController:     NewFakeDrController(),
	}

	if _, err := ctrl.FailoverReplication(context.Background(), req); err != nil {
		t.Errorf("Failed to failover volume replication: %v\n", err)
	}
}

func TestCreateVolumeGroup(t *testing.T) {
	var req = &pb.CreateVolumeGroupOpts{
		Id:          "3769855c-a102-11e7-b772-17b880d2f555",
		Name:        "sample-group-01",
		Description: "This is the first sample group for testing",
		AddVolumes:  []string{"bd5b12a8-a101-11e7-941e-d77981b584d8"},
		Context:     c.NewAdminContext().ToJson(),
	}
	var vg = &SampleVolumeGroups[0]

	mockClient := new(dbtest.Client)
	mockClient.On("GetVolumeGroup", c.NewAdminContext(), req.Id).Return(vg, nil)
	mockClient.On("GetDock", c.NewAdminContext(), "b7602e18-771e-11e7-8f38-dbd6d291f4e0").Return(&SampleDocks[0], nil)
	mockClient.On("UpdateVolume", c.NewAdminContext(), &model.VolumeSpec{
		BaseModel: &model.BaseModel{Id: req.AddVolumes[0]},
		GroupId:   req.Id,
	}).Return(&SampleVolumes[0], nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), vg, model.VolumeGroupAvailable).Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{
				BaseModel: &model.BaseModel{
					Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
				},
				DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			err: nil,
		},
		volumeController: NewFakeVolumeController(),
	}

	if _, err := ctrl.CreateVolumeGroup(context.Background(), req); err != nil {
		t.Errorf("Failed to create volume group: %v\n", err)
	}
}

func TestUpdateVolumeGroup(t *testing.T) {
	var req = &pb.UpdateVolumeGroupOpts{
		Id:         "3769855c-a102-11e7-b772-17b880d2f555",
		AddVolumes: []string{"bd5b12a8-a101-11e7-941e-d77981b584d8"},
		PoolId:     "084bf71e-a102-11e7-88a8-e31fe6d52248",
		Context:    c.NewAdminContext().ToJson(),
	}

	mockClient := new(dbtest.Client)
	mockClient.On("GetDockByPoolId", c.NewAdminContext(), req.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("UpdateVolume", c.NewAdminContext(), &model.VolumeSpec{
		BaseModel: &model.BaseModel{Id: req.AddVolumes[0]},
		GroupId:   req.Id,
	}).Return(&SampleVolumes[0], nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), &SampleVolumeGroups[0], model.VolumeGroupAvailable).Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		volumeController: NewFakeVolumeController(),
	}

	if _, err := ctrl.UpdateVolumeGroup(context.Background(), req); err != nil {
		t.Errorf("Failed to update volume group: %v\n", err)
	}
}

func TestDeleteVolumeGroup(t *testing.T) {
	var req = &pb.DeleteVolumeGroupOpts{
		Id:      "3769855c-a102-11e7-b772-17b880d2f555",
		PoolId:  "084bf71e-a102-11e7-88a8-e31fe6d52248",
		Context: c.NewAdminContext().ToJson(),
	}

	mockClient := new(dbtest.Client)
	mockClient.On("GetDockByPoolId", c.NewAdminContext(), req.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("DeleteVolumeGroup", c.NewAdminContext(), req.Id).Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		volumeController: NewFakeVolumeController(),
	}

	if _, err := ctrl.DeleteVolumeGroup(context.Background(), req); err != nil {
		t.Errorf("Failed to delete volume group: %v\n", err)
	}
}

func NewFakeFileShareController() fileshare.Controller {
	return &fakeFileShareController{}
}

type fakeFileShareController struct {}

func (fakeFileShareController) SetDock(dockInfo *model.DockSpec) { return }

func (fakeFileShareController) CreateFileShare(opt *pb.CreateFileShareOpts) (*model.FileShareSpec, error) {
	return &SampleFileShares[0], nil
}

func (fakeFileShareController) CreateFileShareAcl(opt *pb.CreateFileShareAclOpts) (*model.FileShareAclSpec, error) {
	return &SampleFileSharesAcl[2], nil
}

func (fakeFileShareController) DeleteFileShareAcl(opt *pb.DeleteFileShareAclOpts) error {return nil}

func (fakeFileShareController) DeleteFileShare(opt *pb.DeleteFileShareOpts) error {return nil}

func (fakeFileShareController) CreateFileShareSnapshot(opt *pb.CreateFileShareSnapshotOpts) (*model.FileShareSnapshotSpec, error) {
	return &SampleFileShareSnapshots[0], nil
}

func (fakeFileShareController) DeleteFileShareSnapshot(opts *pb.DeleteFileShareSnapshotOpts) error {return nil}

func TestCreateFileShare(t *testing.T) {
	prf := &SampleFileShareProfiles[0]
	var req = &pb.CreateFileShareOpts{
		Id:          "d2975ebe-d82c-430f-b28e-f373746a71ca",
		Name:        "sample-fileshare-01",
		Description: "This is a sample fileshare for testing",
		Size:        int64(1),
		Profile:     prf.ToJson(),
		Context:     c.NewAdminContext().ToJson(),
	}
	var fileshare = &SampleFileShares[0]
	mockClient := new(dbtest.Client)
	mockClient.On("GetFileShare", c.NewAdminContext(), req.Id).Return(&SampleFileShares[0], nil)
	mockClient.On("GetDock", c.NewAdminContext(), "b7602e18-771e-11e7-8f38-dbd6d291f4e0").Return(&SampleDocks[0], nil)
	mockClient.On("GetFileShareDefaultProfile", c.NewAdminContext()).Return(&SampleFileShareProfiles[0], nil)
	mockClient.On("GetProfile", c.NewAdminContext(), "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&SampleFileShareProfiles[0], nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), fileshare, fileshare.Status).Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{
				BaseModel: &model.BaseModel{
					Id: "bdd44c8e-b8a9-488a-89c0-d1e5beb902dg",
				},
				DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			err: nil,
		},
		fileshareController: NewFakeFileShareController(),
	}
	if _, err := ctrl.CreateFileShare(context.Background(), req); err != nil {
		t.Errorf("failed to create fileshare, err is %v\n", err)
	}
	mockClient1 := new(dbtest.Client)
	mockClient1.On("GetFileShare", c.NewAdminContext(), req.Id).Return(&SampleFileShares[0],fmt.Errorf("specified fileshare(%s) can't find", req.Id))
	mockClient1.On("UpdateStatus", c.NewAdminContext(), fileshare, "error").Return(nil)
	db.C = mockClient1

	var ctrl1  = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{
				BaseModel: &model.BaseModel{
					Id: "bdd44c8e-b8a9-488a-89c0-d1e5beb902dg",
				},
				DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			err: nil,
		},
		fileshareController: NewFakeFileShareController(),
	}

	t1, _:= ctrl1.CreateFileShare(context.Background(), req)
	err_desc := t1.GetError().Description
	expectedError := fmt.Sprintf("specified fileshare(%s) can't find", fileshare.Id)
	if err_desc != expectedError{
		t.Errorf("specified fileshare(%s) can't find\n",fileshare.Id )
	}

	mockClient2 := new(dbtest.Client)
	mockClient2.On("GetFileShare", c.NewAdminContext(), req.Id).Return(&SampleFileShares[0], nil)
	mockClient2.On("UpdateStatus", c.NewAdminContext(), fileshare, "error").Return(nil)
	db.C = mockClient2

	var ctrl2  = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{
				BaseModel: &model.BaseModel{
					Id: "bdd44c8e-b8a9-488a-89c0-d1e5beb902dg",
				},
				DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			err: errors.New("filter supported pools failed: no available pool to meet user's requirement"),
		},
		fileshareController: NewFakeFileShareController(),
	}

	t2, _:= ctrl2.CreateFileShare(context.Background(), req)
	err_desc2 := t2.GetError().Description
	expectedError2 := fmt.Sprintf("filter supported pools failed: no available pool to meet user's requirement")
	if err_desc2 != expectedError2{
		t.Errorf("filter supported pools failed: no available pool to meet user's requirement\n" )
	}

	mockClientd := new(dbtest.Client)
	mockClientd.On("GetFileShare", c.NewAdminContext(), req.Id).Return(&SampleFileShares[0], nil)
	mockClientd.On("GetDock", c.NewAdminContext(), "b7602e18-771e-11e7-8f38-dbd6d291f4e0").Return(nil, fmt.Errorf("when search supported dock resource:when get dock in db:"))
	mockClientd.On("GetFileShareDefaultProfile", c.NewAdminContext()).Return(&SampleFileShareProfiles[0], nil)
	mockClientd.On("GetProfile", c.NewAdminContext(), "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&SampleFileShareProfiles[0], nil)
	mockClientd.On("UpdateStatus", c.NewAdminContext(), fileshare, "error").Return(nil)
	db.C = mockClientd

	var ctrld = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{
				BaseModel: &model.BaseModel{
					Id: "bdd44c8e-b8a9-488a-89c0-d1e5beb902dg",
				},
				DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			err: nil,
		},
		fileshareController: NewFakeFileShareController(),
	}
	td, _:= ctrld.CreateFileShare(context.Background(), req)
	err_descd := td.GetError().Description
	expectedErrord := fmt.Sprintf("when search supported dock resource:when get dock in db:")
	if err_descd != expectedErrord{
		t.Errorf("when search supported dock resource:when get dock in db\n" )
	}

}

func TestDeleteFileShare(t *testing.T) {
	prf := &SampleFileShareProfiles[0]
	var req = &pb.DeleteFileShareOpts{
		Id:      "d2975ebe-d82c-430f-b28e-f373746a71ca",
		Profile: prf.ToJson(),
		PoolId:  "084bf71e-a102-11e7-88a8-e31fe6d52248",
		Context: c.NewAdminContext().ToJson(),
	}
	profile_out := model.NewProfileFromJson(prf.ToJson())
	mockClient := new(dbtest.Client)
	mockClient.On("GetProfile", c.NewAdminContext(), profile_out.Id).Return(&SampleFileShareProfiles[0], nil)
	mockClient.On("GetDockByPoolId", c.NewAdminContext(), req.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("DeleteFileShare", c.NewAdminContext(), req.Id).Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{
				BaseModel: &model.BaseModel{
					Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
				},
				DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			err: nil,
		},
		fileshareController: NewFakeFileShareController(),
	}
	if _, err := ctrl.DeleteFileShare(context.Background(), req); err != nil {
		t.Errorf("failed to delete volume, err is %v\n", err)
	}

	mockClientd := new(dbtest.Client)
	mockClientd.On("GetFileShare", c.NewAdminContext(), req.Id).Return(&SampleFileShares[0], nil)
	mockClientd.On("GetDockByPoolId", c.NewAdminContext(), req.PoolId).Return(nil, fmt.Errorf("when search dock in db by pool id: Get dock failed by pool id:"))
	mockClientd.On("GetFileShare", c.NewAdminContext(),req.Id).Return(&SampleFileShares[0],nil)
	mockClientd.On("UpdateStatus", c.NewAdminContext(), &SampleFileShares[0], "errorDeleting").Return(nil)
	db.C = mockClientd

	var ctrld = &Controller{
		selector: &fakeSelector{
			res: &model.StoragePoolSpec{
				BaseModel: &model.BaseModel{
					Id: "bdd44c8e-b8a9-488a-89c0-d1e5beb902dg",
				},
				DockId: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			err: nil,
		},
		fileshareController: NewFakeFileShareController(),
	}
	td, _:= ctrld.DeleteFileShare(context.Background(), req)
	err_descd := td.GetError().Description
	expectedErrord := fmt.Sprintf("when search dock in db by pool id: Get dock failed by pool id:")
	if err_descd != expectedErrord{
		t.Errorf("when search supported dock dock in db by pool id: Get dock failed by pool id: failed\n" )
	}

}

func TestCreateFileShareSnapshot(t *testing.T) {
	var req = &pb.CreateFileShareSnapshotOpts{
		Id:          "3769855c-a102-11e7-b772-17b880d2f537",
		FileshareId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:        "sample-snapshot-01",
		Description: "This is the first sample snapshot for testing",
		Size:        int64(1),
		Context:     c.NewAdminContext().ToJson(),
	}
	var fileshare = &SampleFileShares[0]
	var snp = &SampleFileShareSnapshots[0]
	mockClient := new(dbtest.Client)
	mockClient.On("GetFileShare", c.NewAdminContext(), req.FileshareId).Return(fileshare, nil)
	mockClient.On("GetDockByPoolId", c.NewAdminContext(), fileshare.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("GetProfile", c.NewAdminContext(), "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&SampleFileShareProfiles[0], nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), snp, "available").Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		fileshareController: NewFakeFileShareController(),
	}

	if _, err := ctrl.CreateFileShareSnapshot(context.Background(), req); err != nil {
		t.Errorf("failed to create file share snapshot: %v\n", err)
	}

	mockClient1 := new(dbtest.Client)
	mockClient1.On("GetFileShare", c.NewAdminContext(), req.FileshareId).Return(nil, errors.New("get file share failed in create file share snapshot method:"))
	mockClient1.On("UpdateFileShareSnapshotStatus", c.NewAdminContext(), req.Id, "error").Return(&SampleFileShareSnapshots[0], nil)
	mockClient1.On("GetFileShareSnapshot", c.NewAdminContext(), req.Id).Return(&SampleFileShareSnapshots[0], nil)
	mockClient1.On("UpdateStatus", c.NewAdminContext(), snp, "error").Return(nil)
	db.C = mockClient1

	var ctrl1 = &Controller{
		fileshareController: NewFakeFileShareController(),
	}
	t1, _:= ctrl1.CreateFileShareSnapshot(context.Background(), req)
	err_desc := t1.GetError().Description
	expectedError := fmt.Sprintf("get file share failed in create file share snapshot method:")
	if err_desc != expectedError{
		t.Errorf("test of create file share snapshot failed, didn't get %v instead got %v\n", expectedError, err_desc)
	}

	mockClient2 := new(dbtest.Client)
	mockClient2.On("GetFileShare", c.NewAdminContext(), req.FileshareId).Return(&SampleFileShares[0], nil)
	mockClient2.On("GetDockByPoolId", c.NewAdminContext(), fileshare.PoolId).Return(nil, errors.New("when search supported dock resource: "))
	mockClient2.On("UpdateFileShareSnapshotStatus", c.NewAdminContext(), req.Id, "error").Return(&SampleFileShareSnapshots[0], nil)
	mockClient2.On("GetFileShareSnapshot", c.NewAdminContext(), req.Id).Return(&SampleFileShareSnapshots[0], nil)
	mockClient2.On("UpdateStatus", c.NewAdminContext(), snp, "error").Return(nil)
	db.C = mockClient2

	var ctrl2 = &Controller{
		fileshareController: NewFakeFileShareController(),
	}
	t2, _:= ctrl2.CreateFileShareSnapshot(context.Background(), req)
	err_desc2 := t2.GetError().Description
	expectedError2 := fmt.Sprintf("when search supported dock resource: ")
	if err_desc2 != expectedError2{
		t.Errorf("test of create file share snapshot failed, didn't get %v instead got %v\n", expectedError2, err_desc2)
	}
}

func TestDeleteFileShareSnapshot(t *testing.T) {
	var req = &pb.DeleteFileShareSnapshotOpts{
		Id:          "3769855c-a102-11e7-b772-17b880d2f537",
		FileshareId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Context:     c.NewAdminContext().ToJson(),
	}
	var fileshare = &SampleShares[0]
	mockClient := new(dbtest.Client)
	mockClient.On("GetFileShare", c.NewAdminContext(), req.FileshareId).Return(fileshare, nil)
	mockClient.On("GetDockByPoolId", c.NewAdminContext(), fileshare.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("DeleteFileShareSnapshot", c.NewAdminContext(), req.Id).Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		fileshareController: NewFakeFileShareController(),
	}
	if _, err := ctrl.DeleteFileShareSnapshot(context.Background(), req); err != nil {
		t.Errorf("failed to delete file share snapshot: %v\n", err)
	}

	mockClient1 := new(dbtest.Client)
	mockClient1.On("GetFileShare", c.NewAdminContext(), req.FileshareId).Return(nil, errors.New("get file share failed in delete file share snapshot method:"))
	mockClient1.On("UpdateFileShareSnapshotStatus", c.NewAdminContext(), req.Id, "errorDeleting").Return(&SampleFileShareSnapshots[0], nil)
	mockClient1.On("GetFileShareSnapshot", c.NewAdminContext(), req.Id).Return(&SampleFileShareSnapshots[0], nil)
	mockClient1.On("UpdateStatus", c.NewAdminContext(), &SampleFileShareSnapshots[0], "errorDeleting").Return(nil)
	db.C = mockClient1

	var ctrl1 = &Controller{
		fileshareController: NewFakeFileShareController(),
	}
	t1, _:= ctrl1.DeleteFileShareSnapshot(context.Background(), req)
	err_desc := t1.GetError().Description
	expectedError := fmt.Sprintf("get file share failed in delete file share snapshot method:")
	if err_desc != expectedError{
		t.Errorf("test of create file share snapshot failed, didn't get %v instead got %v\n", expectedError, err_desc)
	}

	mockClient2 := new(dbtest.Client)
	mockClient2.On("GetFileShare", c.NewAdminContext(), req.FileshareId).Return(&SampleFileShares[0], nil)
	mockClient2.On("GetDockByPoolId", c.NewAdminContext(), "bdd44c8e-b8a9-488a-89c0-d1e5beb902dg").Return(nil, errors.New("when search supported dock resource: "))
	mockClient2.On("UpdateFileShareSnapshotStatus", c.NewAdminContext(), req.Id, "errorDeleting").Return(&SampleFileShareSnapshots[0], nil)
	mockClient2.On("GetFileShareSnapshot", c.NewAdminContext(), req.Id).Return(&SampleFileShareSnapshots[0], nil)
	mockClient2.On("UpdateStatus", c.NewAdminContext(), &SampleFileShareSnapshots[0], "errorDeleting").Return(nil)
	db.C = mockClient2

	var ctrl2 = &Controller{
		fileshareController: NewFakeFileShareController(),
	}
	t2, _:= ctrl2.DeleteFileShareSnapshot(context.Background(), req)
	err_desc2 := t2.GetError().Description
	expectedError2 := fmt.Sprintf("when search supported dock resource: ")
	if err_desc2 != expectedError2{
		t.Errorf("test of create file share snapshot failed, didn't get %v instead got %v\n", expectedError2, err_desc2)
	}
}

func TestCreateFileShareAcl(t *testing.T) {
	var req = &pb.CreateFileShareAclOpts{
		Id: "d2975ebe-d82c-430f-b28e-f373746a71ca",
		Description: "This is a sample Acl for testing",
		Context:     c.NewAdminContext().ToJson(),
		Type: "ip",
		AccessTo: "10.21.23.10",
		AccessCapability:[]string{"Read", "Write"},
	}
	var fileshare= &SampleFileShares[0]
	mockClient := new(dbtest.Client)
	mockClient.On("GetFileShare", c.NewAdminContext(), "").Return(&SampleFileShares[0], nil)
	mockClient.On("GetDockByPoolId", c.NewAdminContext(), fileshare.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("CreateFileShareAcl", c.NewAdminContext(), req).Return(&SampleFileSharesAcl[0], nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), &SampleFileSharesAcl[2], "available").Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		fileshareController: NewFakeFileShareController(),
	}
	if _, err := ctrl.CreateFileShareAcl(context.Background(), req); err != nil {
		t.Errorf("failed to create file share acl: %v\n", err)
	}

	mockClient1 := new(dbtest.Client)
	mockClient1.On("GetFileShare", c.NewAdminContext(), "").Return(nil,fmt.Errorf("specified fileshare(%s) can't find", req.Id))
	mockClient1.On("GetFileShareAcl", c.NewAdminContext(), "d2975ebe-d82c-430f-b28e-f373746a71ca").Return(&SampleFileSharesAcl[2], nil)
	mockClient1.On("UpdateFileShareAclStatus", c.NewAdminContext(), SampleFileSharesAcl[2].Id, "error").Return(nil)
	mockClient1.On("UpdateStatus", c.NewAdminContext(), &SampleFileSharesAcl[2], "error").Return(nil)
	db.C = mockClient1

	var ctrl1 = &Controller{
		fileshareController: NewFakeFileShareController(),
	}

	t1, _:= ctrl1.CreateFileShareAcl(context.Background(), req)
	err_desc := t1.GetError().Description
	expectedError := fmt.Sprintf("specified fileshare(%s) can't find", fileshare.Id)
	if err_desc != expectedError{
		t.Errorf("test of create file share acl failed, didn't get %v instead got %v\n", expectedError, err_desc)
	}

	mockClient2 := new(dbtest.Client)
	mockClient2.On("GetFileShare", c.NewAdminContext(), "").Return(&SampleFileShares[0],nil)
	mockClient2.On("GetDockByPoolId", c.NewAdminContext(), fileshare.PoolId).Return(nil, errors.New("when search supported dock resource:Get dock failed by pool id: "+fileshare.PoolId))
	mockClient2.On("GetFileShareAcl", c.NewAdminContext(), "d2975ebe-d82c-430f-b28e-f373746a71ca").Return(&SampleFileSharesAcl[2], nil)
	mockClient2.On("UpdateFileShareAclStatus", c.NewAdminContext(), SampleFileSharesAcl[2].Id, "error").Return(nil)
	mockClient2.On("UpdateStatus", c.NewAdminContext(), &SampleFileSharesAcl[2], "error").Return(nil)
	db.C = mockClient2

	var ctrl2 = &Controller{
		fileshareController: NewFakeFileShareController(),
	}

	t2, _:= ctrl2.CreateFileShareAcl(context.Background(), req)
	err_desc2 := t2.GetError().Description
	expectedError2 := fmt.Sprintf("when search supported dock resource:Get dock failed by pool id: %v",fileshare.PoolId)
	if err_desc2 != expectedError2{
		t.Errorf("test of create file share acl failed, didn't get %v instead got %v\n", expectedError2, err_desc2)
	}

}

func TestDeleteFileShareAcl(t *testing.T) {
	var req = &pb.DeleteFileShareAclOpts{
		Id:          "d2975ebe-d82c-430f-b28e-f373746a71ca",
		Description: "This is a sample Acl for testing",
		Context:     c.NewAdminContext().ToJson(),
	}
	var fileshare = &SampleFileShares[0]
	mockClient := new(dbtest.Client)
	mockClient.On("GetFileShare", c.NewAdminContext(), "").Return(&SampleFileShares[0], nil)
	mockClient.On("GetDockByPoolId", c.NewAdminContext(), fileshare.PoolId).Return(&SampleDocks[0], nil)
	mockClient.On("DeleteFileShareAcl", c.NewAdminContext(), req.Id).Return(nil)
	mockClient.On("UpdateStatus", c.NewAdminContext(), &SampleFileSharesAcl[2], "available").Return(nil)
	db.C = mockClient

	var ctrl = &Controller{
		fileshareController: NewFakeFileShareController(),
	}
	if _, err := ctrl.DeleteFileShareAcl(context.Background(), req); err != nil {
		t.Errorf("failed to delete file share snapshot: %v\n", err)
	}

	mockClient1 := new(dbtest.Client)
	mockClient1.On("GetFileShare", c.NewAdminContext(), "").Return(nil, fmt.Errorf(" when delete file share acl:specified fileshare(%s) can't find", req.Id))
	mockClient1.On("GetFileShareAcl", c.NewAdminContext(), "d2975ebe-d82c-430f-b28e-f373746a71ca").Return(&SampleFileSharesAcl[2], nil)
	mockClient1.On("UpdateFileShareAclStatus", c.NewAdminContext(), SampleFileSharesAcl[2].Id, "error").Return(nil)
	mockClient1.On("UpdateStatus", c.NewAdminContext(), &SampleFileSharesAcl[2], "errorDeleting").Return(nil)
	db.C = mockClient1

	var ctrl1 = &Controller{
		fileshareController: NewFakeFileShareController(),
	}

	t1, _ := ctrl1.DeleteFileShareAcl(context.Background(), req)
	err_desc := t1.GetError().Description
	expectedError := fmt.Sprintf(" when delete file share acl:specified fileshare(%s) can't find", fileshare.Id)
	if err_desc != expectedError {
		t.Errorf("test of delete file share acl failed, didn't get %v instead got %v\n", expectedError, err_desc)
	}

	mockClient2 := new(dbtest.Client)
	mockClient2.On("GetFileShare", c.NewAdminContext(), "").Return(&SampleFileShares[0], nil)
	mockClient2.On("GetDockByPoolId", c.NewAdminContext(), fileshare.PoolId).Return(nil, errors.New("when search supported dock resource:Get dock failed by pool id: "+fileshare.PoolId))
	mockClient2.On("GetFileShareAcl", c.NewAdminContext(), "d2975ebe-d82c-430f-b28e-f373746a71ca").Return(&SampleFileSharesAcl[2], nil)
	mockClient2.On("UpdateFileShareAclStatus", c.NewAdminContext(), SampleFileSharesAcl[2].Id, "error").Return(nil)
	mockClient2.On("UpdateStatus", c.NewAdminContext(), &SampleFileSharesAcl[2], "errorDeleting").Return(nil)
	db.C = mockClient2

	var ctrl2 = &Controller{
		fileshareController: NewFakeFileShareController(),
	}

	t2, _ := ctrl2.DeleteFileShareAcl(context.Background(), req)
	err_desc2 := t2.GetError().Description
	expectedError2 := fmt.Sprintf("when search supported dock resource:Get dock failed by pool id: %v", fileshare.PoolId)
	if err_desc2 != expectedError2 {
		t.Errorf("test of delete file share acl failed, didn't get %v instead got %v\n", expectedError2, err_desc2)
	}

}
