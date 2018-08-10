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
This module implements the database operation of data structure
defined in api module.

*/

package db

import (
	"fmt"
	"strings"

	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db/drivers/etcd"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/pkg/utils/config"
	fakedb "github.com/opensds/opensds/testutils/db"
)

// C is a global variable that controls database module.
var C Client

// Init function can perform some initialization work of different databases.
func Init(db *Database) {
	switch db.Driver {
	case "mysql":
		// C = mysql.Init(db.Driver, db.Crendential)
		fmt.Printf("mysql is not implemented right now!")
		return
	case "etcd":
		C = etcd.NewClient(strings.Split(db.Endpoint, ","))
		return
	case "fake":
		C = fakedb.NewFakeDbClient()
		return
	default:
		fmt.Printf("Can't find database driver %s!\n", db.Driver)
	}
}

// Client is an interface for exposing some operations of managing database
// client.
type Client interface {
	CreateDock(ctx *c.Context, dck *model.DockSpec) (*model.DockSpec, error)

	GetDock(ctx *c.Context, dckID string) (*model.DockSpec, error)

	ListDocks(ctx *c.Context) ([]*model.DockSpec, error)

	ListDocksWithFilter(ctx *c.Context, m map[string][]string) ([]*model.DockSpec, error)

	UpdateDock(ctx *c.Context, dckID, name, desp string) (*model.DockSpec, error)

	DeleteDock(ctx *c.Context, dckID string) error

	GetDockByPoolId(ctx *c.Context, poolId string) (*model.DockSpec, error)

	CreatePool(ctx *c.Context, pol *model.StoragePoolSpec) (*model.StoragePoolSpec, error)

	GetPool(ctx *c.Context, polID string) (*model.StoragePoolSpec, error)

	ListAvailabilityZones(ctx *c.Context) ([]string, error)

	ListPools(ctx *c.Context) ([]*model.StoragePoolSpec, error)

	ListPoolsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.StoragePoolSpec, error)

	UpdatePool(ctx *c.Context, polID, name, desp string, usedCapacity int64, used bool) (*model.StoragePoolSpec, error)

	DeletePool(ctx *c.Context, polID string) error

	CreateProfile(ctx *c.Context, prf *model.ProfileSpec) (*model.ProfileSpec, error)

	GetProfile(ctx *c.Context, prfID string) (*model.ProfileSpec, error)

	GetDefaultProfile(ctx *c.Context) (*model.ProfileSpec, error)

	ListProfiles(ctx *c.Context) ([]*model.ProfileSpec, error)

	ListProfilesWithFilter(ctx *c.Context, m map[string][]string) ([]*model.ProfileSpec, error)

	UpdateProfile(ctx *c.Context, prfID string, input *model.ProfileSpec) (*model.ProfileSpec, error)

	DeleteProfile(ctx *c.Context, prfID string) error

	AddExtraProperty(ctx *c.Context, prfID string, ext model.ExtraSpec) (*model.ExtraSpec, error)

	ListExtraProperties(ctx *c.Context, prfID string) (*model.ExtraSpec, error)

	RemoveExtraProperty(ctx *c.Context, prfID, extraKey string) error

	CreateVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error)

	GetVolume(ctx *c.Context, volID string) (*model.VolumeSpec, error)

	ListVolumes(ctx *c.Context) ([]*model.VolumeSpec, error)

	ListVolumesWithFilter(ctx *c.Context, m map[string][]string) ([]*model.VolumeSpec, error)

	UpdateVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error)

	DeleteVolume(ctx *c.Context, volID string) error

	ExtendVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error)

	CreateVolumeAttachment(ctx *c.Context, attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error)

	GetVolumeAttachment(ctx *c.Context, attachmentId string) (*model.VolumeAttachmentSpec, error)

	ListVolumeAttachments(ctx *c.Context, volumeId string) ([]*model.VolumeAttachmentSpec, error)

	ListVolumeAttachmentsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.VolumeAttachmentSpec, error)

	UpdateVolumeAttachment(ctx *c.Context, attachmentId string, attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error)

	DeleteVolumeAttachment(ctx *c.Context, attachmentId string) error

	CreateVolumeSnapshot(ctx *c.Context, vs *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error)

	GetVolumeSnapshot(ctx *c.Context, snapshotID string) (*model.VolumeSnapshotSpec, error)

	ListVolumeSnapshots(ctx *c.Context) ([]*model.VolumeSnapshotSpec, error)

	ListVolumeSnapshotsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.VolumeSnapshotSpec, error)

	UpdateVolumeSnapshot(ctx *c.Context, snapshotID string, vs *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error)

	DeleteVolumeSnapshot(ctx *c.Context, snapshotID string) error

	CreateReplication(ctx *c.Context, replication *model.ReplicationSpec) (*model.ReplicationSpec, error)

	GetReplication(ctx *c.Context, replicationId string) (*model.ReplicationSpec, error)

	GetReplicationByVolumeId(ctx *c.Context, volumeId string) (*model.ReplicationSpec, error)

	ListReplication(ctx *c.Context) ([]*model.ReplicationSpec, error)

	ListReplicationWithFilter(ctx *c.Context, m map[string][]string) ([]*model.ReplicationSpec, error)

	DeleteReplication(ctx *c.Context, replicationId string) error

	UpdateReplication(ctx *c.Context, replicationId string, input *model.ReplicationSpec) (*model.ReplicationSpec, error)

	CreateVolumeGroup(ctx *c.Context, vg *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error)

	GetVolumeGroup(ctx *c.Context, vgId string) (*model.VolumeGroupSpec, error)

	UpdateVolumeGroup(ctx *c.Context, vg *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error)

	UpdateStatus(ctx *c.Context, object interface{}, status string) error

	ListVolumesByGroupId(ctx *c.Context, vgId string) ([]*model.VolumeSpec, error)

	ListSnapshotsByVolumeId(ctx *c.Context, volId string) ([]*model.VolumeSnapshotSpec, error)

	DeleteVolumeGroup(ctx *c.Context, vgId string) error

	ListVolumeGroups(ctx *c.Context) ([]*model.VolumeGroupSpec, error)

	VolumesToUpdate(ctx *c.Context, volumeList []*model.VolumeSpec) ([]*model.VolumeSpec, error)

	ListVolumeGroupsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.VolumeGroupSpec, error)
}
