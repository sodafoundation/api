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
		Status:           model.VOLUME_CREATING,
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

	if volume.Status != model.VOLUME_AVAILABLE {
		errMsg := "The status of the volume to be extended must be available"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	volume.Status = model.VOLUME_EXTENDING
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
	if vol.Status != model.VOLUME_AVAILABLE {
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
		Status:         model.VOLUMEATM_CREATING,
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
	if vol.Status != model.VOLUME_AVAILABLE && vol.Status != model.VOLUME_IN_USE {
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
		Status:      model.VOLUMESNAP_CREATING,
	}

	result, err := db.C.CreateVolumeSnapshot(ctx, snap)
	if err != nil {
		log.Error("Error occurred in dock module when create volume snapshot in db:", err)
		return nil, err
	}
	return result, nil
}

func DeleteVolumeSnapshotDBEntry(ctx *c.Context, in *model.VolumeSnapshotSpec) error {
	if in.Status != model.VOLUMESNAP_AVAILABLE {
		errMsg := "Only the volume snapshot with the status available can be deleted"
		log.Error(errMsg)
		return errors.New(errMsg)
	}
	in.Status = model.VOLUMESNAP_DELETING
	_, err := db.C.UpdateVolumeSnapshot(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

//Just modify the state of the volume to be deleted in the DB, the real deletion in another thread
func DeleteVolumeDBEntry(ctx *c.Context, in *model.VolumeSpec) error {
	invalidStatus := []string{model.VOLUME_AVAILABLE, model.VOLUME_ERROR,
		model.VOLUEM_ERROR_DELETING, model.VOLUME_ERROR_EXTENDING}
	if !utils.Contained(in.Status, invalidStatus) {
		errMsg := fmt.Sprintf("Can't delete the volume in %s", in.Status)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	in.Status = model.VOLUME_DELETING
	_, err := db.C.UpdateVolume(ctx, in)
	if err != nil {
		return err
	}
	return nil
}
