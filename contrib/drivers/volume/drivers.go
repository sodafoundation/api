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
This module defines an standard table of storage driver. The default storage
driver is sample driver used for testing. If you want to use other storage
plugin, just modify Init() and Clean() method.

*/

package volume

import (
	_ "github.com/opensds/opensds/contrib/backup/multicloud"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
)

// VolumeDriver is an interface for exposing some operations of different volume
// drivers, currently support sample, lvm, ceph, cinder and so forth.
type VolumeDriver interface {
	//Any initialization the volume driver does while starting.
	Setup() error
	//Any operation the volume driver does while stopping.
	Teardown() error

	CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error)

	PullVolume(volIdentifier string) (*model.VolumeSpec, error)

	DeleteVolume(opt *pb.DeleteVolumeOpts) error

	ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error)

	InitializeConnection(opt *pb.CreateVolumeAttachmentOpts) (*model.ConnectionInfo, error)

	TerminateConnection(opt *pb.DeleteVolumeAttachmentOpts) error

	CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error)

	PullSnapshot(snapIdentifier string) (*model.VolumeSnapshotSpec, error)

	DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error

	InitializeSnapshotConnection(opt *pb.CreateSnapshotAttachmentOpts) (*model.ConnectionInfo, error)

	TerminateSnapshotConnection(opt *pb.DeleteSnapshotAttachmentOpts) error

	// NOTE Parameter vg means complete volume group information, because driver
	// may use it to do something and return volume group status.
	CreateVolumeGroup(opt *pb.CreateVolumeGroupOpts) (*model.VolumeGroupSpec, error)

	// NOTE Parameter addVolumesRef or removeVolumesRef means complete volume
	// information that will be added or removed from group. Driver may use
	// them to do some related operations and return their status.
	UpdateVolumeGroup(opt *pb.UpdateVolumeGroupOpts) (*model.VolumeGroupSpec, error)

	// NOTE Parameter volumes means volumes deleted from group, driver may use
	// their compelete information to do some related operations and return
	// their status.
	DeleteVolumeGroup(opt *pb.DeleteVolumeGroupOpts) error

	ListPools() ([]*model.StoragePoolSpec, error)
}

type MetricDriver interface {
	//Any initialization the metric driver does while starting.
	Setup() error
	//Any operation the metric driver does while stopping.
	Teardown() error
	// Collect metrics for all supported resources
	CollectMetrics() ([]*model.MetricSpec, error)
}
