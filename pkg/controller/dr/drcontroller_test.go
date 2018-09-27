// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package dr

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
	"github.com/stretchr/testify/mock"
)

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

func (fvc *fakeVolumeController) AttachVolume(*pb.AttachVolumeOpts) (string, error) {
	return "/dev/disk/by-path/ip-192.168.56.100:3260-iscsi-iqn.2017-10.io.opensds:baec258b-8f79-4bbc-bf97-28addfa903d3-lun-1", nil
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
	return nil, nil
}

func (fvc *fakeVolumeController) UpdateVolumeGroup(*pb.UpdateVolumeGroupOpts) error {
	return nil
}

func (fvc *fakeVolumeController) DeleteVolumeGroup(*pb.DeleteVolumeGroupOpts) error {
	return nil
}
func (fvc *fakeVolumeController) SetDock(dockInfo *model.DockSpec) { return }

var (
	pool = model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
		},
		Name:             "sample-pool-01",
		Description:      "This is the first sample storage pool for testing",
		TotalCapacity:    int64(100),
		FreeCapacity:     int64(90),
		DockId:           "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		AvailabilityZone: "default",
		ReplicationType:  model.ReplicationTypeArray,
	}
	volumes = []model.VolumeSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			},
			Name:        "sample-volume",
			Description: "This is a sample volume for testing",
			Size:        int64(1),
			Status:      "available",
			PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
			Metadata: map[string]string{
				"lvPath": "/dev/opensds-volumes-default/volume-ab14d4ea-edd4-41bd-b37b-391f66115e8b",
			},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "e000bf78-7cf7-4fd2-a085-e94bd61daf31",
			},
			Name:        "sample-volume",
			Description: "This is a sample volume for testing",
			Size:        int64(1),
			Status:      "available",
			PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
			Metadata: map[string]string{
				"lvPath": "/dev/opensds-volumes-default/volume-7bce5fb6-a229-4584-bad4-15f1a6a6aadd",
			},
		},
	}
)

func TestArrayBasedCreateReplication(t *testing.T) {
	mockClient := new(dbtest.Client)
	pool.ReplicationType = model.ReplicationTypeArray
	mockClient.On("GetProfile", context.NewAdminContext(), "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&SampleProfiles[0], nil)
	mockClient.On("GetPool", context.NewAdminContext(), "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&pool, nil)
	mockClient.On("GetDock", context.NewAdminContext(), "b7602e18-771e-11e7-8f38-dbd6d291f4e0").Return(&SampleDocks[0], nil)
	mockClient.On("GetDockByPoolId", context.NewAdminContext(), "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&SampleDocks[0], nil)
	db.C = mockClient

	r := &model.ReplicationSpec{
		BaseModel: &model.BaseModel{
			Id: "c299a978-4f3e-11e8-8a5c-977218a83359",
		},
		PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
	}

	c := NewController(NewFakeVolumeController())
	result, err := c.CreateReplication(context.NewAdminContext(), r, &volumes[0], &volumes[1])
	if err != nil {
		t.Error("Test DR CreateReplication failed, ", err)
	}

	var expected = &model.ReplicationSpec{
		BaseModel: &model.BaseModel{
			Id: "c299a978-4f3e-11e8-8a5c-977218a83359",
		},
		PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
		PrimaryReplicationDriverData: map[string]string{
			"lvPath": "/dev/opensds-volumes-default/volume-ab14d4ea-edd4-41bd-b37b-391f66115e8b",
		},
		SecondaryReplicationDriverData: map[string]string{
			"lvPath": "/dev/opensds-volumes-default/volume-7bce5fb6-a229-4584-bad4-15f1a6a6aadd",
		},
		Metadata: map[string]string{},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestHostBasedCreateReplication(t *testing.T) {
	mockClient := new(dbtest.Client)
	pool.ReplicationType = model.ReplicationTypeHost
	mockClient.On("GetProfile", context.NewAdminContext(), "1106b972-66ef-11e7-b172-db03f3689c9c").Return(&SampleProfiles[0], nil)
	mockClient.On("GetPool", context.NewAdminContext(), "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&pool, nil)
	mockClient.On("GetDock", context.NewAdminContext(), mock.Anything).Return(&SampleDocks[0], nil)
	mockClient.On("GetDockByPoolId", context.NewAdminContext(), "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&SampleDocks[0], nil)
	mockClient.On("CreateVolumeAttachment", context.NewAdminContext(), &SampleAttachments[0]).Return(&SampleAttachments[0], nil)
	mockClient.On("UpdateVolumeAttachment", context.NewAdminContext(), SampleAttachments[0].Id, &SampleAttachments[0]).Return(&SampleAttachments[0], nil)
	volumeList := []*model.VolumeSpec{}
	mockClient.On("ListVolumes", context.NewAdminContext()).Return(volumeList, nil)

	db.C = mockClient

	r := &model.ReplicationSpec{
		BaseModel: &model.BaseModel{
			Id: "c299a978-4f3e-11e8-8a5c-977218a83359",
		},
		PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
	}

	c := NewController(NewFakeVolumeController())
	result, err := c.CreateReplication(context.NewAdminContext(), r, &volumes[0], &volumes[1])
	if err != nil {
		t.Error("Test DR CreateReplication failed, ", err)
	}

	var expected = &model.ReplicationSpec{
		BaseModel: &model.BaseModel{
			Id: "c299a978-4f3e-11e8-8a5c-977218a83359",
		},
		PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
		PrimaryReplicationDriverData: map[string]string{
			"lvPath":       "/dev/opensds-volumes-default/volume-ab14d4ea-edd4-41bd-b37b-391f66115e8b",
			"Mountpoint":   "/dev/disk/by-path/ip-192.168.56.100:3260-iscsi-iqn.2017-10.io.opensds:baec258b-8f79-4bbc-bf97-28addfa903d3-lun-1",
			"AttachmentId": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
			"HostName":     "",
			"HostIp":       "",
		},
		SecondaryReplicationDriverData: map[string]string{
			"lvPath":       "/dev/opensds-volumes-default/volume-7bce5fb6-a229-4584-bad4-15f1a6a6aadd",
			"Mountpoint":   "/dev/disk/by-path/ip-192.168.56.100:3260-iscsi-iqn.2017-10.io.opensds:baec258b-8f79-4bbc-bf97-28addfa903d3-lun-1",
			"AttachmentId": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
			"HostName":     "",
			"HostIp":       "",
		},
		Metadata: map[string]string{},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestArrayBasedDeleteReplication(t *testing.T) {
	pool.ReplicationType = model.ReplicationTypeArray
	mockClient := new(dbtest.Client)
	mockClient.On("GetPool", context.NewAdminContext(), "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&pool, nil)
	mockClient.On("GetDock", context.NewAdminContext(), "b7602e18-771e-11e7-8f38-dbd6d291f4e0").Return(&SampleDocks[0], nil)
	mockClient.On("GetDockByPoolId", context.NewAdminContext(), "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&SampleDocks[0], nil)
	mockClient.On("DeleteReplication", context.NewAdminContext(), "c299a978-4f3e-11e8-8a5c-977218a83359").Return(nil)
	mockClient.On("UpdateVolume", context.NewAdminContext(), mock.Anything).Return(&SampleVolumes[0], nil)
	db.C = mockClient

	r := &model.ReplicationSpec{
		BaseModel: &model.BaseModel{
			Id: "c299a978-4f3e-11e8-8a5c-977218a83359",
		},
		PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
	}

	c := NewController(NewFakeVolumeController())
	err := c.DeleteReplication(context.NewAdminContext(), r, &volumes[0], &volumes[1])
	if err != nil {
		t.Error("Test DR DeleteReplication failed, ", err)
	}
}

func TestHostBasedDeleteReplication(t *testing.T) {
	pool.ReplicationType = model.ReplicationTypeHost
	mockClient := new(dbtest.Client)
	mockClient.On("GetPool", context.NewAdminContext(), "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&pool, nil)
	mockClient.On("GetDock", context.NewAdminContext(), mock.Anything).Return(&SampleDocks[0], nil)
	mockClient.On("GetDockByPoolId", context.NewAdminContext(), "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&SampleDocks[0], nil)
	mockClient.On("DeleteReplication", context.NewAdminContext(), "c299a978-4f3e-11e8-8a5c-977218a83359").Return(nil)
	mockClient.On("GetVolumeAttachment", context.NewAdminContext(), "f2dda3d2-bf79-11e7-8665-f750b088f63e").Return(&SampleAttachments[0], nil)
	mockClient.On("DeleteVolumeAttachment", context.NewAdminContext(), "f2dda3d2-bf79-11e7-8665-f750b088f63e").Return(nil)
	mockClient.On("UpdateVolume", context.NewAdminContext(), mock.Anything).Return(&SampleVolumes[0], nil)

	db.C = mockClient

	r := &model.ReplicationSpec{
		BaseModel: &model.BaseModel{
			Id: "c299a978-4f3e-11e8-8a5c-977218a83359",
		},
		PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
		PrimaryReplicationDriverData: map[string]string{
			"lvPath":       "/dev/opensds-volumes-default/volume-ab14d4ea-edd4-41bd-b37b-391f66115e8b",
			"Mountpoint":   "/dev/disk/by-path/ip-192.168.56.100:3260-iscsi-iqn.2017-10.io.opensds:baec258b-8f79-4bbc-bf97-28addfa903d3-lun-1",
			"AttachmentId": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
			"HostName":     "",
			"HostIp":       "",
		},
		SecondaryReplicationDriverData: map[string]string{
			"lvPath":       "/dev/opensds-volumes-default/volume-7bce5fb6-a229-4584-bad4-15f1a6a6aadd",
			"Mountpoint":   "/dev/disk/by-path/ip-192.168.56.100:3260-iscsi-iqn.2017-10.io.opensds:baec258b-8f79-4bbc-bf97-28addfa903d3-lun-1",
			"AttachmentId": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
			"HostName":     "",
			"HostIp":       "",
		},
	}

	c := NewController(NewFakeVolumeController())
	err := c.DeleteReplication(context.NewAdminContext(), r, &volumes[0], &volumes[1])
	if err != nil {
		t.Error("Test DR DeleteReplication failed, ", err)
	}
}

func TestEnableReplication(t *testing.T) {
	pool.ReplicationType = model.ReplicationTypeArray
	mockClient := new(dbtest.Client)
	mockClient.On("GetPool", context.NewAdminContext(), "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&pool, nil)
	mockClient.On("GetDock", context.NewAdminContext(), "b7602e18-771e-11e7-8f38-dbd6d291f4e0").Return(&SampleDocks[0], nil)
	mockClient.On("GetDockByPoolId", context.NewAdminContext(), "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&SampleDocks[0], nil)
	db.C = mockClient

	r := &model.ReplicationSpec{
		BaseModel: &model.BaseModel{
			Id: "c299a978-4f3e-11e8-8a5c-977218a83359",
		},
		PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
	}

	c := NewController(NewFakeVolumeController())
	err := c.EnableReplication(context.NewAdminContext(), r, &volumes[0], &volumes[1])
	if err != nil {
		t.Error("Test DR EnableReplication failed, ", err)
	}
}

func TestDisableReplication(t *testing.T) {
	pool.ReplicationType = model.ReplicationTypeArray
	mockClient := new(dbtest.Client)
	mockClient.On("GetPool", context.NewAdminContext(), "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&pool, nil)
	mockClient.On("GetDock", context.NewAdminContext(), "b7602e18-771e-11e7-8f38-dbd6d291f4e0").Return(&SampleDocks[0], nil)
	mockClient.On("GetDockByPoolId", context.NewAdminContext(), "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&SampleDocks[0], nil)
	db.C = mockClient

	r := &model.ReplicationSpec{
		BaseModel: &model.BaseModel{
			Id: "c299a978-4f3e-11e8-8a5c-977218a83359",
		},
		PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
	}

	c := NewController(NewFakeVolumeController())
	err := c.DisableReplication(context.NewAdminContext(), r, &volumes[0], &volumes[1])
	if err != nil {
		t.Error("Test DR DisableReplication failed, ", err)
	}
}

func TestFailoverReplication(t *testing.T) {
	pool.ReplicationType = model.ReplicationTypeArray
	mockClient := new(dbtest.Client)
	mockClient.On("GetPool", context.NewAdminContext(), "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&pool, nil)
	mockClient.On("GetDock", context.NewAdminContext(), "b7602e18-771e-11e7-8f38-dbd6d291f4e0").Return(&SampleDocks[0], nil)
	mockClient.On("GetDockByPoolId", context.NewAdminContext(), "084bf71e-a102-11e7-88a8-e31fe6d52248").Return(&SampleDocks[0], nil)
	db.C = mockClient

	r := &model.ReplicationSpec{
		BaseModel: &model.BaseModel{
			Id: "c299a978-4f3e-11e8-8a5c-977218a83359",
		},
		PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
		SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:              "sample-replication-01",
		Description:       "This is a sample replication for testing",
		PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
		ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
	}
	f := &model.FailoverReplicationSpec{
		AllowAttachedVolume: true,
		SecondaryBackendId:  model.ReplicationDefaultBackendId,
	}
	c := NewController(NewFakeVolumeController())
	err := c.FailoverReplication(context.NewAdminContext(), r, f, &volumes[0], &volumes[1])
	if err != nil {
		t.Error("Test DR FailoverReplication failed, ", err)
	}
}
