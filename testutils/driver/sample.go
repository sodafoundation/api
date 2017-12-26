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

/*
This module implements a sample driver for OpenSDS. This driver will handle all
operations of volume and return a fake value.

*/

package sample

import (
	"errors"

	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
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

// InitializeConnection
func (*Driver) InitializeConnection(opt *pb.CreateAttachmentOpts) (*model.ConnectionInfo, error) {
	return &SampleConnection, nil
}

// TerminateConnection
func (*Driver) TerminateConnection(opt *pb.DeleteAttachmentOpts) error { return nil }

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
