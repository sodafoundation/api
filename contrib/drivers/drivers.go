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

package drivers

import (
	_ "github.com/opensds/opensds/contrib/backup/multicloud"
	"github.com/opensds/opensds/contrib/drivers/ceph"
	"github.com/opensds/opensds/contrib/drivers/hpe/nimble"
	"github.com/opensds/opensds/contrib/drivers/huawei/dorado"
	"github.com/opensds/opensds/contrib/drivers/huawei/fusionstorage"
	"github.com/opensds/opensds/contrib/drivers/lvm"
	"github.com/opensds/opensds/contrib/drivers/openstack/cinder"
	"github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
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

// Init
func Init(resourceType string) VolumeDriver {
	var d VolumeDriver
	switch resourceType {
	case config.CinderDriverType:
		d = &cinder.Driver{}
		break
	case config.CephDriverType:
		d = &ceph.Driver{}
		break
	case config.LVMDriverType:
		d = &lvm.Driver{}
		break
	case config.HuaweiDoradoDriverType:
		d = &dorado.Driver{}
		break
	case config.HuaweiFusionStorageDriverType:
		d = &fusionstorage.Driver{}
	case config.HpeNimbleDriverType:
		d = &nimble.Driver{}
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
	case *fusionstorage.Driver:
		break
	case *nimble.Driver:
		break
	default:
		break
	}
	d.Unset()
	d = nil

	return d
}

func CleanMetricDriver(d MetricDriver) MetricDriver {
	// Execute different clean operations according to the MetricDriver type.
	switch d.(type) {
	case *lvm.MetricDriver:
		break
	default:
		break
	}
	_ = d.Teardown()
	d = nil

	return d
}

type MetricDriver interface {
	//Any initialization the metric driver does while starting.
	Setup() error
	//Any operation the metric driver does while stopping.
	Teardown() error
	// Collect metrics for all supported resources
	CollectMetrics() ([]*model.MetricSpec, error)
}

// Init
func InitMetricDriver(resourceType string) MetricDriver {
	var d MetricDriver
	switch resourceType {
	case config.LVMDriverType:
		d = &lvm.MetricDriver{}
		break
	case config.CephDriverType:
		d = &ceph.MetricDriver{}
		break
	case config.HuaweiDoradoDriverType:
		d = &dorado.MetricDriver{}
		break
	default:
		//d = &sample.Driver{}
		break
	}
	d.Setup()
	return d
}
