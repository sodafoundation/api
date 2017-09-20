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

	log.Println("[Info] Calling volume plugin to create volume...")

	//Call function of StoragePlugins configured by storage plugins.
	return d.VolumePlugin.CreateVolume(name, size)
}

func (d *DockHub) GetVolume(volID string) (*api.VolumeSpec, error) {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	log.Println("[Info] Calling volume plugin to get volume...")

	//Call function of StoragePlugins configured by storage plugins.
	return d.VolumePlugin.GetVolume(volID)
}

func (d *DockHub) DeleteVolume(volID string) error {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	log.Println("[Info] Calling volume plugin to delete volume...")

	//Call function of StoragePlugins configured by storage plugins.
	return d.VolumePlugin.DeleteVolume(volID)
}

func (d *DockHub) CreateVolumeAttachment(volID string, doLocalAttach, multiPath bool, hostInfo *api.HostInfo) (*api.VolumeAttachmentSpec, error) {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	log.Println("[Info] Calling volume plugin to initialize connection...")

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

	log.Println("[Info] Calling volume plugin to attach volume...")

	//Call function of StoragePlugins configured by storage plugins.
	return d.VolumePlugin.AttachVolume(volID, host, mountpoint)
}

func (d *DockHub) DeleteVolumeAttachment(volID string) error {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	log.Println("[Info] Calling volume plugin to detach volume...")

	//Call function of StoragePlugins configured by storage plugins.
	return d.VolumePlugin.DetachVolume(volID)
}

func (d *DockHub) CreateSnapshot(name, volID, description string) (*api.VolumeSnapshotSpec, error) {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	log.Println("[Info] Calling volume plugin to create snapshot...")

	//Call function of StoragePlugins configured by storage plugins.
	return d.VolumePlugin.CreateSnapshot(name, volID, description)
}

func (d *DockHub) GetSnapshot(snapID string) (*api.VolumeSnapshotSpec, error) {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	log.Println("[Info] Calling volume plugin to get snapshot...")

	//Call function of StoragePlugins configured by storage plugins.
	return d.VolumePlugin.GetSnapshot(snapID)
}

func (d *DockHub) DeleteSnapshot(snapID string) error {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	log.Println("[Info] Calling volume plugin to delete snapshot...")

	//Call function of StoragePlugins configured by storage plugins.
	return d.VolumePlugin.DeleteSnapshot(snapID)
}

func (d *DockHub) ListPools() (*[]api.StoragePoolSpec, error) {
	//Get the storage plugins and do some initializations.
	d.VolumePlugin = plugins.InitVP(d.ResourceType)

	log.Println("[Info] Calling volume plugin to list pools...")

	//Call function of StoragePlugins configured by storage plugins.
	return d.VolumePlugin.ListPools()
}
