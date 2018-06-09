// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package drivers

import (
	"github.com/opensds/opensds/contrib/drivers/ceph"
	"github.com/opensds/opensds/contrib/drivers/huawei/dorado"
	"github.com/opensds/opensds/contrib/drivers/lvm"
	"github.com/opensds/opensds/contrib/drivers/openstack/cinder"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
	sample "github.com/opensds/opensds/testutils/driver"
)

// VolumeDriver is an interface for exposing some operations of different volume
// drivers, currently support sample, lvm, ceph, cinder and so forth.
type VolumeDriver interface {
	//Any initialization the volume driver does while starting.
	Setup() error
	//Any operation the volume driver does while stopping.
	Unset() error

	CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error)

	PullVolume(volIdentifier string) (*model.VolumeSpec, error)

	DeleteVolume(opt *pb.DeleteVolumeOpts) error

	ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error)

	InitializeConnection(opt *pb.CreateAttachmentOpts) (*model.ConnectionInfo, error)

	TerminateConnection(opt *pb.DeleteAttachmentOpts) error

	CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error)

	PullSnapshot(snapIdentifier string) (*model.VolumeSnapshotSpec, error)

	DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error

	// NOTE Parameter vg means complete volume group information, because driver
	// may use it to do something and return volume group status.
	CreateVolumeGroup(opt *pb.CreateVolumeGroupOpts, vg *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error)

	// NOTE Parameter addVolumesRef or removeVolumesRef means complete volume
	// information that will be added or removed from group. Driver may use
	// them to do some related operations and return their status.
	UpdateVolumeGroup(opt *pb.UpdateVolumeGroupOpts, vg *model.VolumeGroupSpec, addVolumesRef []*model.VolumeSpec, removeVolumesRef []*model.VolumeSpec) (*model.VolumeGroupSpec, []*model.VolumeSpec, []*model.VolumeSpec, error)

	// NOTE Parameter volumes means volumes deleted from group, driver may use
	// their compelete information to do some related operations and return
	// their status.
	DeleteVolumeGroup(opt *pb.DeleteVolumeGroupOpts, vg *model.VolumeGroupSpec, volumes []*model.VolumeSpec) (*model.VolumeGroupSpec, []*model.VolumeSpec, error)

	ListPools() ([]*model.StoragePoolSpec, error)
}

// Init
func Init(resourceType string) VolumeDriver {
	var d VolumeDriver
	switch resourceType {
	case "cinder":
		d = &cinder.Driver{}
		break
	case "ceph":
		d = &ceph.Driver{}
		break
	case "lvm":
		d = &lvm.Driver{}
		break
	case "huawei_dorado":
		d = &dorado.Driver{}
		break
	default:
		d = &sample.Driver{}
		break
	}
	d.Setup()
	return d
}

// Clean
func Clean(d VolumeDriver) VolumeDriver {
	// Execute different clean operations according to the VolumeDriver type.
	switch d.(type) {
	case *cinder.Driver:
		break
	case *ceph.Driver:
		break
	case *lvm.Driver:
		break
	case *dorado.Driver:
		break
	default:
		break
	}
	d.Unset()
	d = nil

	return d
}
