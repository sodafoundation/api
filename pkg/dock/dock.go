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
	"log"

	"github.com/opensds/opensds/contrib/plugins"
	api "github.com/opensds/opensds/pkg/model"
)

// A reference to DockHub structure with fields that represent some required
// parameters for initializing and controlling the volume plugin.
type DockHub struct {
	// ResourceType represents the type of backend resources. This field is used
	// for initializing the specified volume plugin.
	ResourceType string

	plugins.VolumePlugin
}

func NewDockHub(resourceType string) *DockHub {
	return &DockHub{
		ResourceType: resourceType,
	}
}

func (d *DockHub) CreateVolume(name string, size int64) (*api.VolumeSpec, error) {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	//Call function of StoragePlugins configured by storage plugins.
	result, err := d.VolumePlugin.CreateVolume(name, size)
	if err != nil {
		log.Println("Call plugin to create volume failed:", err)
		return &api.VolumeSpec{}, err
	} else {
		return result, nil
	}
}

func (d *DockHub) GetVolume(volID string) (*api.VolumeSpec, error) {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	//Call function of StoragePlugins configured by storage plugins.
	result, err := d.VolumePlugin.GetVolume(volID)
	if err != nil {
		log.Println("Call plugin to get volume failed:", err)
		return &api.VolumeSpec{}, err
	} else {
		return result, nil
	}
}

func (d *DockHub) DeleteVolume(volID string) error {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	//Call function of StoragePlugins configured by storage plugins.
	if err := d.VolumePlugin.DeleteVolume(volID); err != nil {
		log.Println("Call plugin to delete volume failed:", err)
		return err
	}
	return nil
}

func (d *DockHub) CreateVolumeAttachment(volID string, doLocalAttach, multiPath bool, hostInfo *api.HostInfo) (*api.VolumeAttachmentSpec, error) {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	//Call function of StoragePlugins configured by storage plugins.
	connInfo, err := d.VolumePlugin.InitializeConnection(volID, doLocalAttach, multiPath, hostInfo)
	if err != nil {
		log.Println("Call plugin to initialize volume connection failed:", err)
		return &api.VolumeAttachmentSpec{}, err
	}

	return &api.VolumeAttachmentSpec{
		HostInfo:       hostInfo,
		ConnectionInfo: connInfo,
	}, nil
}

func (d *DockHub) UpdateVolumeAttachment(volID, host, mountpoint string) error {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	//Call function of StoragePlugins configured by storage plugins.
	if err := d.VolumePlugin.AttachVolume(volID, host, mountpoint); err != nil {
		log.Println("Call plugin to update volume attachment failed:", err)
		return err
	}
	return nil
}

func (d *DockHub) DeleteVolumeAttachment(volID string) error {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	//Call function of StoragePlugins configured by storage plugins.
	if err := d.VolumePlugin.DetachVolume(volID); err != nil {
		log.Println("Call plugin to delete volume attachment failed:", err)
		return err
	}
	return nil
}

func (d *DockHub) CreateSnapshot(name, volID, description string) (*api.VolumeSnapshotSpec, error) {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	//Call function of StoragePlugins configured by storage plugins.
	result, err := d.VolumePlugin.CreateSnapshot(name, volID, description)
	if err != nil {
		log.Println("Call plugin to create snapshot failed:", err)
		return &api.VolumeSnapshotSpec{}, err
	} else {
		return result, nil
	}
}

func (d *DockHub) GetSnapshot(snapID string) (*api.VolumeSnapshotSpec, error) {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	//Call function of StoragePlugins configured by storage plugins.
	result, err := d.VolumePlugin.GetSnapshot(snapID)
	if err != nil {
		log.Println("Call plugin to get snapshot failed:", err)
		return &api.VolumeSnapshotSpec{}, err
	} else {
		return result, nil
	}
}

func (d *DockHub) DeleteSnapshot(snapID string) error {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	//Call function of StoragePlugins configured by storage plugins.
	if err := d.VolumePlugin.DeleteSnapshot(snapID); err != nil {
		log.Println("Call plugin to delete snapshot failed:", err)
		return err
	}
	return nil
}
