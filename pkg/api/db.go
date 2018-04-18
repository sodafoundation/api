// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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
This module implements API-related database operations.
*/

package api

import (
	"errors"
	"fmt"

	"time"

	log "github.com/golang/glog"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/constants"
	"github.com/satori/go.uuid"
)

func CreateVolumeDBEntry(ctx *c.Context, in *model.VolumeSpec) (*model.VolumeSpec, error) {
	if in.Id == "" {
		in.Id = uuid.NewV4().String()
	}
	if in.Size <= 0 {
		errMsg := fmt.Sprintf("Invalid volume size: %d", in.Size)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if in.AvailabilityZone == "" {
		log.Warning("Use default availability zone when user doesn't specify availabilityZone.")
		in.AvailabilityZone = "default"
	}
	if in.CreatedAt == "" {
		in.CreatedAt = time.Now().Format(constants.TimeFormat)
	}
	vol := &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id:        in.Id,
			CreatedAt: in.CreatedAt,
		},
		UserId:           in.UserId,
		TenantId:         in.TenantId,
		Name:             in.Name,
		Description:      in.Description,
		Size:             in.Size,
		AvailabilityZone: in.AvailabilityZone,
		Status:           model.VolumeCreating,
	}
	result, err := db.C.CreateVolume(ctx, vol)
	if err != nil {
		log.Error("When add volume to db:", err)
		return nil, err
	}

	return result, nil
}

func ExtendVolumeDBEntry(ctx *c.Context, volID string) (*model.VolumeSpec, error) {
	volume, err := db.C.GetVolume(ctx, volID)
	if err != nil {
		log.Error("Get volume failed in extend volume method: ", err)
		return nil, err
	}

	if volume.Status != model.VolumeAvailable {
		errMsg := "The status of the volume to be extended must be available"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	volume.Status = model.VolumeExtending
	// Store the volume data into database.
	result, err := db.C.ExtendVolume(ctx, volume)
	if err != nil {
		log.Error("When extend volume in db module:", err)
		return nil, err
	}
	return result, nil
}

func CreateVolumeAttachmentDBEntry(ctx *c.Context, in *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	vol, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in create volume attachment method: ", err)
		return nil, err
	}
	if vol.Status != model.VolumeAvailable {
		errMsg := "Only the status of volume is available, attachment can be created"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if in.Id == "" {
		in.Id = uuid.NewV4().String()
	}
	if in.CreatedAt == "" {
		in.CreatedAt = time.Now().Format(constants.TimeFormat)
	}
	if len(in.AdditionalProperties) == 0 {
		in.AdditionalProperties = map[string]interface{}{"attachment": "attachment"}
	}
	if len(in.ConnectionData) == 0 {
		in.ConnectionData = map[string]interface{}{"attachment": "attachment"}
	}

	var atc = &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id:        in.Id,
			CreatedAt: in.CreatedAt,
		},
		VolumeId: in.VolumeId,
		HostInfo: model.HostInfo{
			Platform:  in.Platform,
			OsType:    in.OsType,
			Ip:        in.Ip,
			Host:      in.Host,
			Initiator: in.Initiator,
		},
		Status:         model.VolumeAttachCreating,
		Metadata:       utils.MergeStringMaps(in.Metadata, vol.Metadata),
		ConnectionInfo: in.ConnectionInfo,
	}

	result, err := db.C.CreateVolumeAttachment(ctx, atc)
	if err != nil {
		log.Error("Error occurred in dock module when create volume attachment in db:", err)
		return nil, err
	}
	return result, nil
}

func CreateVolumeSnapshotDBEntry(ctx *c.Context, in *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	vol, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in create volume snapshot method: ", err)
		return nil, err
	}
	if vol.Status != model.VolumeAvailable && vol.Status != model.VolumeInUse {
		var errMsg = "Only the status of volume is available or in-use, the snapshot can be created"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	if in.Id == "" {
		in.Id = uuid.NewV4().String()
	}

	if in.CreatedAt == "" {
		in.CreatedAt = time.Now().Format(constants.TimeFormat)
	}

	var snap = &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id:        in.Id,
			CreatedAt: in.CreatedAt,
		},
		Name:        in.Name,
		Description: in.Description,
		VolumeId:    in.VolumeId,
		Size:        vol.Size,
		Metadata:    utils.MergeStringMaps(in.Metadata, vol.Metadata),
		Status:      model.VolumeSnapCreating,
	}

	result, err := db.C.CreateVolumeSnapshot(ctx, snap)
	if err != nil {
		log.Error("Error occurred in dock module when create volume snapshot in db:", err)
		return nil, err
	}
	return result, nil
}

func DeleteVolumeSnapshotDBEntry(ctx *c.Context, in *model.VolumeSnapshotSpec) error {
	if in.Status != model.VolumeSnapAvailable {
		errMsg := "Only the volume snapshot with the status available can be deleted"
		log.Error(errMsg)
		return errors.New(errMsg)
	}
	in.Status = model.VolumeSnapDeleting
	_, err := db.C.UpdateVolumeSnapshot(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

//Just modify the state of the volume to be deleted in the DB, the real deletion in another thread
func DeleteVolumeDBEntry(ctx *c.Context, in *model.VolumeSpec) error {
	if in.Status != model.VolumeAvailable {
		errMsg := "Only the volume with the status available can be deleted"
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	in.Status = model.VolumeDeleting
	_, err := db.C.UpdateVolume(ctx, in)
	if err != nil {
		return err
	}
	return nil
}

func DeleteReplicationDBEntry(ctx *c.Context, in *model.ReplicationSpec) error {
	invalidStatus := []string{model.ReplicationCreating, model.ReplicationDeleting,
		model.ReplicationEnabling, model.ReplicationDisabling, model.ReplicationFailovering}

	if utils.Contained(in.Status, invalidStatus) {
		errMsg := fmt.Sprintf("can delete the replication in %s", in.Status)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	in.Status = model.ReplicationDeleting
	_, err := db.C.UpdateReplication(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

func EnableReplicationDBEntry(ctx *c.Context, in *model.ReplicationSpec) error {
	if in.Status != model.ReplicationAvailable {
		errMsg := "Only the replication with the status available can be enbaled"
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	in.Status = model.ReplicationEnabling
	_, err := db.C.UpdateReplication(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

func DisableReplicationDBEntry(ctx *c.Context, in *model.ReplicationSpec) error {
	if in.Status != model.ReplicationAvailable {
		errMsg := "Only the replication with the status available can be enbaled"
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	in.Status = model.ReplicationDisabling
	_, err := db.C.UpdateReplication(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

func FailoverReplicationDBEntry(ctx *c.Context, in *model.ReplicationSpec) error {
	if in.Status != model.ReplicationAvailable {
		errMsg := "Only the replication with the status available can be enbaled"
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	in.Status = model.ReplicationFailovering
	_, err := db.C.UpdateReplication(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}
