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
This module implements the database operation of data structure
defined in api module.

*/

package db

import (
	"fmt"

	"github.com/opensds/opensds/pkg/db/drivers/etcd"
	_ "github.com/opensds/opensds/pkg/db/drivers/mysql"
	api "github.com/opensds/opensds/pkg/model"
)

var C Client

type DBConfig struct {
	DriverName string
	Endpoints  []string
	Credential string
}

func Init(conf *DBConfig) {
	switch conf.DriverName {
	case "mysql":
		// C = mysql.Init(conf.DriverName, conf.Crendential)
		fmt.Errorf("mysql is not implemented right now!")
	case "etcd":
		C = etcd.Init(conf.Endpoints)
	default:
		fmt.Errorf("Can't find database driver %s!\n", conf.DriverName)
	}
}

type Client interface {
	CreateDock(dck *api.DockSpec) (*api.DockSpec, error)

	GetDock(dckID string) (*api.DockSpec, error)

	ListDocks() (*[]api.DockSpec, error)

	UpdateDock(dckID, name, desp string) (*api.DockSpec, error)

	DeleteDock(dckID string) error

	CreatePool(pol *api.StoragePoolSpec) (*api.StoragePoolSpec, error)

	GetPool(polID string) (*api.StoragePoolSpec, error)

	ListPools() (*[]api.StoragePoolSpec, error)

	UpdatePool(polID, name, desp string, usedCapacity int64, used bool) (*api.StoragePoolSpec, error)

	DeletePool(polID string) error

	CreateProfile(prf *api.ProfileSpec) (*api.ProfileSpec, error)

	GetProfile(prfID string) (*api.ProfileSpec, error)

	ListProfiles() (*[]api.ProfileSpec, error)

	UpdateProfile(prfID string, input *api.ProfileSpec) (*api.ProfileSpec, error)

	DeleteProfile(prfID string) error

	AddExtraProperty(prfID string, ext api.ExtraSpec) (*api.ExtraSpec, error)

	ListExtraProperties(prfID string) (*api.ExtraSpec, error)

	RemoveExtraProperty(prfID, extraKey string) error

	CreateVolume(vol *api.VolumeSpec) (*api.VolumeSpec, error)

	GetVolume(volID string) (*api.VolumeSpec, error)

	ListVolumes() (*[]api.VolumeSpec, error)

	DeleteVolume(volID string) error

	CreateVolumeAttachment(volID string, atc *api.VolumeAttachmentSpec) (*api.VolumeAttachmentSpec, error)

	GetVolumeAttachment(volID, attachmentID string) (*api.VolumeAttachmentSpec, error)

	ListVolumeAttachments(volID string) (*[]api.VolumeAttachmentSpec, error)

	UpdateVolumeAttachment(volID, attachmentID, mountpoint string, hostInfo *api.HostInfo) (*api.VolumeAttachmentSpec, error)

	DeleteVolumeAttachment(volID, attachmentID string) error

	CreateVolumeSnapshot(vs *api.VolumeSnapshotSpec) (*api.VolumeSnapshotSpec, error)

	GetVolumeSnapshot(snapshotID string) (*api.VolumeSnapshotSpec, error)

	ListVolumeSnapshots() (*[]api.VolumeSnapshotSpec, error)

	DeleteVolumeSnapshot(snapshotID string) error
}
