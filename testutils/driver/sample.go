// Copyright 2019 The OpenSDS Authors.
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

/*
This module implements a sample driver for OpenSDS. This driver will handle all
operations of volume and return a fake value.

*/

package sample

import (
	"errors"

	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	. "github.com/opensds/opensds/testutils/collection"
)

// Driver
type Driver struct{}

// Setup
func (*Driver) Setup() error { return nil }

// Unset
func (*Driver) Unset() error { return nil }

// CreateVolume
func (*Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	return &SampleVolumes[0], nil
}

// PullVolume
func (*Driver) PullVolume(volIdentifier string) (*model.VolumeSpec, error) {
	for _, volume := range SampleVolumes {
		if volIdentifier == volume.Id {
			return &volume, nil
		}
	}

	return nil, errors.New("Can't find volume " + volIdentifier)
}

// DeleteVolume
func (*Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	return nil
}

// ExtendVolume ...
func (*Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	return &SampleVolumes[0], nil
}

// InitializeConnection
func (*Driver) InitializeConnection(opt *pb.CreateVolumeAttachmentOpts) (*model.ConnectionInfo, error) {
	return &SampleConnection, nil
}

// TerminateConnection
func (*Driver) TerminateConnection(opt *pb.DeleteVolumeAttachmentOpts) error { return nil }

// CreateSnapshot
func (*Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	return &SampleSnapshots[0], nil
}

// PullSnapshot
func (*Driver) PullSnapshot(snapIdentifier string) (*model.VolumeSnapshotSpec, error) {
	for _, snapshot := range SampleSnapshots {
		if snapIdentifier == snapshot.Id {
			return &snapshot, nil
		}
	}

	return nil, errors.New("Can't find snapshot " + snapIdentifier)
}

// DeleteSnapshot
func (*Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	return nil
}

// ListPools
func (*Driver) ListPools() ([]*model.StoragePoolSpec, error) {
	var pols []*model.StoragePoolSpec

	for i := range SamplePools {
		pols = append(pols, &SamplePools[i])
	}
	return pols, nil
}

func (d *Driver) InitializeSnapshotConnection(opt *pb.CreateSnapshotAttachmentOpts) (*model.ConnectionInfo, error) {
	return nil, &model.NotImplementError{S: "method InitializeSnapshotConnection has not been implemented yet"}
}

func (d *Driver) TerminateSnapshotConnection(opt *pb.DeleteSnapshotAttachmentOpts) error {
	return &model.NotImplementError{S: "method TerminateSnapshotConnection has not been implemented yet"}
}

func (d *Driver) CreateVolumeGroup(opt *pb.CreateVolumeGroupOpts) (*model.VolumeGroupSpec, error) {
	return nil, &model.NotImplementError{"method CreateVolumeGroup has not been implemented yet"}
}

func (d *Driver) UpdateVolumeGroup(opt *pb.UpdateVolumeGroupOpts) (*model.VolumeGroupSpec, error) {
	return nil, &model.NotImplementError{"method UpdateVolumeGroup has not been implemented yet"}
}

func (d *Driver) DeleteVolumeGroup(opt *pb.DeleteVolumeGroupOpts) error {
	return &model.NotImplementError{"method DeleteVolumeGroup has not been implemented yet"}
}

func (d *Driver) CreateFileShare(opt *pb.CreateFileShareOpts) (*model.FileShareSpec, error) {
	return &SampleFileShares[0], nil
}

func (d *Driver) DeleteFileShare(opt *pb.DeleteFileShareOpts) error {
	return nil
}

func (d *Driver) CreateFileShareAcl(opt *pb.CreateFileShareAclOpts) (*model.FileShareAclSpec, error) {
	return &SampleFileSharesAcl[0], nil
}

func (d *Driver) DeleteFileShareAcl(opt *pb.DeleteFileShareAclOpts) error {
	return nil
}

// CreateFileShareSnapshot
func (d *Driver) CreateFileShareSnapshot(opt *pb.CreateFileShareSnapshotOpts) (*model.FileShareSnapshotSpec, error) {
	return &SampleFileShareSnapshots[0], nil
}

// DeleteFileShareSnapshot
func (d *Driver) DeleteFileShareSnapshot(opt *pb.DeleteFileShareSnapshotOpts) error {
	return nil
}
