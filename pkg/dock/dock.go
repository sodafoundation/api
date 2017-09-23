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

/*
This module implements a standard SouthBound interface of resources to
storage plugins.

*/

package dock

import (
	log "github.com/golang/glog"

	"github.com/opensds/opensds/contrib/drivers"
	api "github.com/opensds/opensds/pkg/model"
)

// A reference to DockHub structure with fields that represent some required
// parameters for initializing and controlling the volume driver.
type DockHub struct {
	// ResourceType represents the type of backend resources. This field is used
	// for initializing the specified volume driver.
	ResourceType string

	Driver drivers.VolumeDriver
}

func NewDockHub(resourceType string) *DockHub {
	return &DockHub{
		ResourceType: resourceType,
	}
}

func (d *DockHub) CreateVolume(name string, size int64) (*api.VolumeSpec, error) {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(d.ResourceType)

	log.Info("Calling volume driver to create volume...")

	//Call function of StorageDrivers configured by storage drivers.
	return d.Driver.CreateVolume(name, size)
}

func (d *DockHub) GetVolume(volID string) (*api.VolumeSpec, error) {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(d.ResourceType)

	log.Info("Calling volume driver to get volume...")

	//Call function of StorageDrivers configured by storage drivers.
	return d.Driver.GetVolume(volID)
}

func (d *DockHub) DeleteVolume(volID string) error {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(d.ResourceType)

	log.Info("Calling volume driver to delete volume...")

	//Call function of StorageDrivers configured by storage drivers.
	return d.Driver.DeleteVolume(volID)
}

func (d *DockHub) CreateVolumeAttachment(volID string, doLocalAttach, multiPath bool, hostInfo *api.HostInfo) (*api.VolumeAttachmentSpec, error) {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(d.ResourceType)

	log.Info("Calling volume driver to initialize volume connection...")

	//Call function of StorageDrivers configured by storage drivers.
	connInfo, err := d.Driver.InitializeConnection(volID, doLocalAttach, multiPath, hostInfo)
	if err != nil {
		log.Error("Call driver to initialize volume connection failed:", err)
		return &api.VolumeAttachmentSpec{}, err
	}

	return &api.VolumeAttachmentSpec{
		HostInfo:       hostInfo,
		ConnectionInfo: connInfo,
	}, nil
}

func (d *DockHub) UpdateVolumeAttachment(volID, host, mountpoint string) error {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(d.ResourceType)

	log.Info("Calling volume driver to attach volume...")

	//Call function of StorageDrivers configured by storage drivers.
	return d.Driver.AttachVolume(volID, host, mountpoint)
}

func (d *DockHub) DeleteVolumeAttachment(volID string) error {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(d.ResourceType)

	log.Info("Calling volume driver to detach volume...")

	//Call function of StorageDrivers configured by storage drivers.
	return d.Driver.DetachVolume(volID)
}

func (d *DockHub) CreateSnapshot(name, volID, description string) (*api.VolumeSnapshotSpec, error) {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(d.ResourceType)

	log.Info("Calling volume driver to create snapshot...")

	//Call function of StorageDrivers configured by storage drivers.
	return d.Driver.CreateSnapshot(name, volID, description)
}

func (d *DockHub) GetSnapshot(snapID string) (*api.VolumeSnapshotSpec, error) {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(d.ResourceType)

	log.Info("Calling volume driver to get snapshot...")

	//Call function of StorageDrivers configured by storage drivers.
	return d.Driver.GetSnapshot(snapID)
}

func (d *DockHub) DeleteSnapshot(snapID string) error {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(d.ResourceType)

	log.Info("Calling volume driver to delete snapshot...")

	//Call function of StorageDrivers configured by storage drivers.
	return d.Driver.DeleteSnapshot(snapID)
}

func (d *DockHub) ListPools() (*[]api.StoragePoolSpec, error) {
	//Get the storage drivers and do some initializations.
	d.Driver = drivers.Init(d.ResourceType)

	log.Info("Calling volume driver to list pools...")

	//Call function of StorageDrivers configured by storage drivers.
	return d.Driver.ListPools()
}
