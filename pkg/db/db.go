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
This module implements the database operation of data structure
defined in api module.

*/

package db

import (
	"fmt"
	"strings"

	"github.com/opensds/opensds/pkg/db/drivers/etcd"
	_ "github.com/opensds/opensds/pkg/db/drivers/mysql"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/pkg/utils/config"
	fakedb "github.com/opensds/opensds/testutils/db"
)

var C Client

func Init(db *Database) {
	switch db.Driver {
	case "mysql":
		// C = mysql.Init(db.Driver, db.Crendential)
		fmt.Errorf("mysql is not implemented right now!")
		return
	case "etcd":
		C = etcd.NewClient(strings.Split(db.Endpoint, ","))
		return
	case "fake":
		C = fakedb.NewFakeDbClient()
		return
	default:
		fmt.Errorf("Can't find database driver %s!\n", db.Driver)
	}
}

type Client interface {
	CreateDock(dck *model.DockSpec) error

	GetDock(dckID string) (*model.DockSpec, error)

	ListDocks() ([]*model.DockSpec, error)

	UpdateDock(dckID, name, desp string) (*model.DockSpec, error)

	DeleteDock(dckID string) error

	GetDockByPoolId(poolId string) (*model.DockSpec, error)

	CreatePool(pol *model.StoragePoolSpec) error

	GetPool(polID string) (*model.StoragePoolSpec, error)

	ListPools() ([]*model.StoragePoolSpec, error)

	UpdatePool(polID, name, desp string, usedCapacity int64, used bool) (*model.StoragePoolSpec, error)

	DeletePool(polID string) error

	CreateProfile(prf *model.ProfileSpec) error

	GetProfile(prfID string) (*model.ProfileSpec, error)

	GetDefaultProfile() (*model.ProfileSpec, error)

	ListProfiles() ([]*model.ProfileSpec, error)

	UpdateProfile(prfID string, input *model.ProfileSpec) (*model.ProfileSpec, error)

	DeleteProfile(prfID string) error

	AddExtraProperty(prfID string, ext model.ExtraSpec) (*model.ExtraSpec, error)

	ListExtraProperties(prfID string) (*model.ExtraSpec, error)

	RemoveExtraProperty(prfID, extraKey string) error

	CreateVolume(vol *model.VolumeSpec) error

	GetVolume(volID string) (*model.VolumeSpec, error)

	ListVolumes() ([]*model.VolumeSpec, error)

	DeleteVolume(volID string) error

	CreateVolumeAttachment(attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error)

	GetVolumeAttachment(attachmentId string) (*model.VolumeAttachmentSpec, error)

	ListVolumeAttachments(volumeId string) ([]*model.VolumeAttachmentSpec, error)

	UpdateVolumeAttachment(attachmentId string, attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error)

	DeleteVolumeAttachment(attachmentId string) error

	CreateVolumeSnapshot(vs *model.VolumeSnapshotSpec) error

	GetVolumeSnapshot(snapshotID string) (*model.VolumeSnapshotSpec, error)

	ListVolumeSnapshots() ([]*model.VolumeSnapshotSpec, error)

	DeleteVolumeSnapshot(snapshotID string) error
}
