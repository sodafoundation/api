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
This module defines an standard table of storage plugin. The default storage
plugin is ups plugin. If you want to use other storage plugin, just modify
Init() method.

*/

package plugins

import (
	"github.com/opensds/opensds/contrib/plugins/ceph"
	"github.com/opensds/opensds/contrib/plugins/upsplugin"
	api "github.com/opensds/opensds/pkg/model"
)

type VolumePlugin interface {
	//Any initialization the volume driver does while starting.
	Setup()
	//Any operation the volume driver does while stoping.
	Unset()

	CreateVolume(name string, size int64) (*api.VolumeSpec, error)

	GetVolume(volID string) (*api.VolumeSpec, error)

	DeleteVolume(volID string) error

	InitializeConnection(volID string, doLocalAttach, multiPath bool, hostInfo *api.HostInfo) (*api.ConnectionInfo, error)

	AttachVolume(volID, host, mountpoint string) error

	DetachVolume(volID string) error

	CreateSnapshot(name, volID, description string) (*api.VolumeSnapshotSpec, error)

	GetSnapshot(snapID string) (*api.VolumeSnapshotSpec, error)

	DeleteSnapshot(snapID string) error

	ListPools() (*[]api.StoragePoolSpec, error)
}

func InitVP(resourceType string) VolumePlugin {
	switch resourceType {
	case "ceph":
		return &ceph.CephPlugin{}
	default:
		return &upsplugin.Plugin{}
	}
}
