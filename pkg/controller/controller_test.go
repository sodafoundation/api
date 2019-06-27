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
	"testing"

	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller/dr"
	"github.com/opensds/opensds/pkg/controller/volume"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	. "github.com/opensds/opensds/testutils/collection"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
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
