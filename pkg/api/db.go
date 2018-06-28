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
	"strings"
	"time"

	log "github.com/golang/glog"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller"
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
	if in.SnapshotId != "" {
		snap, err := db.C.GetVolumeSnapshot(ctx, in.SnapshotId)
		if err != nil {
			log.Error("Get snapshot failed in create volume method: ", err)
			return nil, err
		}
		if snap.Status != model.VolumeSnapAvailable {
			var errMsg = "Only if the snapshot is available, the volume can be created"
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		if snap.Size > in.Size {
			var errMsg = "Size of volume must be equal to or bigger than size of the snapshot"
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
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
		UserId:           ctx.UserId,
		Name:             in.Name,
		Description:      in.Description,
		ProfileId:        in.ProfileId,
		Size:             in.Size,
		AvailabilityZone: in.AvailabilityZone,
		Status:           model.VolumeCreating,
		SnapshotId:       in.SnapshotId,
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
	validStatus := []string{model.VolumeSnapAvailable, model.VolumeSnapError,
		model.VolumeSnapErrorDeleting}
	if !utils.Contained(in.Status, validStatus) {
		errMsg := fmt.Sprintf("Only the volume snapshot with the status available, error, error_deleting can be deleted, the volume status is %s", in.Status)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	// If volume id is invalid, it would mean that volume snapshot creation failed before the create method
	// in storage driver was called, and delete its db entry directly.
	_, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		if err := db.C.DeleteVolumeSnapshot(ctx, in.Id); err != nil {
			log.Error("when delete volume snapshot in db:", err)
			return err
		}
		return nil
	}

	in.Status = model.VolumeSnapDeleting
	_, err = db.C.UpdateVolumeSnapshot(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

//Just modify the state of the volume to be deleted in the DB, the real deletion in another thread
func DeleteVolumeDBEntry(ctx *c.Context, in *model.VolumeSpec) error {
	validStatus := []string{model.VolumeAvailable, model.VolumeError,
		model.VolumeErrorDeleting, model.VolumeErrorExtending}
	if !utils.Contained(in.Status, validStatus) {
		errMsg := fmt.Sprintf("Only the volume with the status available, error, error_deleting, error_extending can be deleted, the volume status is %s", in.Status)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	// If profileId or poolId of the volume doesn't exist, it would mean that the volume provisioning operation failed before the create method
	// in storage driver was called, therefore the volume entry should be deleted from db direclty.
	if in.ProfileId == "" || in.PoolId == "" {
		if err := db.C.DeleteVolume(ctx, in.Id); err != nil {
			log.Error("when delete volume in db:", err)
			return err
		}
		return nil
	}

	snaps, err := db.C.ListSnapshotsByVolumeId(ctx, in.Id)
	if err != nil {
		return err
	}

	if len(snaps) > 0 {
		return fmt.Errorf("Volume %s can not be deleted, because it still has snapshots", in.Id)
	}

	volAttachments, err := db.C.ListVolumeAttachments(ctx, in.Id)
	if err != nil {
		return err
	}

	if len(volAttachments) > 0 {
		return fmt.Errorf("Volume %s can not be deleted, because it's in use", in.Id)
	}

	in.Status = model.VolumeDeleting
	_, err = db.C.UpdateVolume(ctx, in)
	if err != nil {
		return err
	}
	return nil
}

func DeleteReplicationDBEntry(ctx *c.Context, in *model.ReplicationSpec) error {
	invalidStatus := []string{model.ReplicationCreating, model.ReplicationDeleting, model.ReplicationEnabling,
		model.ReplicationDisabling, model.ReplicationFailingOver, model.ReplicationFailingBack}

	if utils.Contained(in.ReplicationStatus, invalidStatus) {
		errMsg := fmt.Sprintf("Can't delete the replication in %s", in.ReplicationStatus)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	in.ReplicationStatus = model.ReplicationDeleting
	_, err := db.C.UpdateReplication(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

func EnableReplicationDBEntry(ctx *c.Context, in *model.ReplicationSpec) error {
	invalidStatus := []string{model.ReplicationCreating, model.ReplicationDeleting, model.ReplicationEnabling,
		model.ReplicationDisabling, model.ReplicationFailingOver, model.ReplicationFailingBack}
	if utils.Contained(in.ReplicationStatus, invalidStatus) {
		errMsg := fmt.Sprintf("Can't enable the replication in %s", in.ReplicationStatus)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	in.ReplicationStatus = model.ReplicationEnabling
	_, err := db.C.UpdateReplication(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

func DisableReplicationDBEntry(ctx *c.Context, in *model.ReplicationSpec) error {
	invalidStatus := []string{model.ReplicationCreating, model.ReplicationDeleting, model.ReplicationEnabling,
		model.ReplicationDisabling, model.ReplicationFailingOver, model.ReplicationFailingBack}
	if utils.Contained(in.ReplicationStatus, invalidStatus) {
		errMsg := fmt.Sprintf("Can't disable the replication in %s", in.ReplicationStatus)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	in.ReplicationStatus = model.ReplicationDisabling
	_, err := db.C.UpdateReplication(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

func FailoverReplicationDBEntry(ctx *c.Context, in *model.ReplicationSpec, secondaryBackendId string) error {
	invalidStatus := []string{model.ReplicationCreating, model.ReplicationDeleting, model.ReplicationEnabling,
		model.ReplicationDisabling, model.ReplicationFailingOver, model.ReplicationFailingBack}
	if utils.Contained(in.ReplicationStatus, invalidStatus) {
		errMsg := fmt.Sprintf("Can't fail over/back the replication in %s", in.ReplicationStatus)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	if secondaryBackendId == model.ReplicationDefaultBackendId {
		in.ReplicationStatus = model.ReplicationFailingOver
	} else {
		in.ReplicationStatus = model.ReplicationFailingBack
	}
	_, err := db.C.UpdateReplication(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

func CreateVolumeGroupDBEntry(ctx *c.Context, in *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
	if len(in.Profiles) == 0 {
		msg := fmt.Sprintf("Profiles must be provided to create volume group.")
		log.Error(msg)
		return nil, errors.New(msg)
	}

	if in.Id == "" {
		in.Id = uuid.NewV4().String()
	}
	if in.AvailabilityZone == "" {
		log.Warning("Use default availability zone when user doesn't specify availabilityZone.")
		in.AvailabilityZone = "default"
	}

	vg := &model.VolumeGroupSpec{
		BaseModel: &model.BaseModel{
			Id: in.Id,
		},
		UserId:           ctx.UserId,
		Name:             in.Name,
		Description:      in.Description,
		AvailabilityZone: in.AvailabilityZone,
		Status:           model.VolumeGroupCreating,
		Profiles:         in.Profiles,
	}
	result, err := db.C.CreateVolumeGroup(ctx, vg)
	if err != nil {
		log.Error("When add volume to db:", err)
		return nil, err
	}
	// TODO:Rpc call to create group.
	// Create volume group request is sent to the Dock. Dock will update volume status to "available"
	// after volume group creation is completed.
	controller.Brain.CreateVolumeGroup(ctx, vg)
	return result, nil
}

func UpdateVolumeGroupDBEntry(ctx *c.Context, vgUpdate *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
	vg, err := db.C.GetVolumeGroup(ctx, vgUpdate.Id)
	if err != nil {
		return nil, err
	}

	var name string
	if vg.Name == vgUpdate.Name {
		name = ""
	} else {
		name = vgUpdate.Name
	}
	var description string
	if vg.Description == vgUpdate.Description {
		description = ""
	} else {
		description = vgUpdate.Description
	}
	vgUpdate.Profiles = vg.Profiles
	vgUpdate.PoolId = vg.PoolId

	var invalid_uuids []string
	for _, uuidAdd := range vgUpdate.AddVolumes {
		for _, uuidRemove := range vgUpdate.RemoveVolumes {
			if uuidAdd == uuidRemove {
				invalid_uuids = append(invalid_uuids, uuidAdd)
			}
		}
	}
	if len(invalid_uuids) > 0 {
		msg := fmt.Sprintf("UUID %s is in both add and remove volume list", strings.Join(invalid_uuids, ","))
		log.Error(msg)
		return nil, errors.New(msg)
	}

	volumes, err := db.C.ListVolumesByGroupId(ctx, vgUpdate.Id)
	if err != nil {
		return nil, err
	}

	var addVolumesNew, removeVolumeNew []string
	// Validate volumes in AddVolumes and RemoveVolumes.
	if len(vgUpdate.AddVolumes) > 0 {
		if addVolumesNew, err = ValidateAddVolumes(ctx, volumes, vgUpdate.AddVolumes, vgUpdate); err != nil {
			return nil, err
		}
	}
	if len(vgUpdate.RemoveVolumes) > 0 {
		if removeVolumeNew, err = ValidateRemoveVolumes(ctx, volumes, vgUpdate.RemoveVolumes, vgUpdate); err != nil {
			return nil, err
		}
	}

	if name == "" && description == "" && len(addVolumesNew) == 0 && len(removeVolumeNew) == 0 {
		msg := fmt.Sprintf("Update group %s faild, because no valid name, description, addvolumes or removevolumes were provided", vgUpdate.Id)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	vgNew := &model.VolumeGroupSpec{
		BaseModel: &model.BaseModel{
			Id: vg.Id,
		},
	}

	vgNew.UpdatedAt = time.Now().Format(constants.TimeFormat)
	// Only update name or description. No need to send them over through an RPC call and set status to available.
	if name != "" {
		vgNew.Name = name
	}
	if description != "" {
		vgNew.Description = description
	}
	if len(addVolumesNew) == 0 && len(removeVolumeNew) == 0 {
		vgNew.Status = model.VolumeGroupAvailable
	} else {
		vgNew.Status = model.VolumeGroupUpdating
	}

	result, err := db.C.UpdateVolumeGroup(ctx, vgNew)
	if err != nil {
		log.Error("When update volume group in db:", err.Error())
		return nil, err
	}

	//TODO: Do an RPC call only if addVolumesNew or removeVolumeNew is not nil.
	if len(addVolumesNew) > 0 || len(removeVolumeNew) > 0 {
		controller.Brain.UpdateVolumeGroup(ctx, vg, addVolumesNew, removeVolumeNew)
	}

	return result, nil
}

func ValidateAddVolumes(ctx *c.Context, volumes []*model.VolumeSpec, addVolumes []string, vg *model.VolumeGroupSpec) ([]string, error) {
	var addVolumeRef []string
	var flag bool
	for _, volumeId := range addVolumes {
		flag = true
		for _, volume := range volumes {
			if volumeId == volume.Id {
				// Volume already in group. Remove it from addVolumes.
				flag = false
				break
			}
		}
		if flag {
			addVolumeRef = append(addVolumeRef, volumeId)
		}
	}

	var addVolumesNew []string
	for _, addVol := range addVolumeRef {
		addVolRef, err := db.C.GetVolume(ctx, addVol)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot add volume %s to group %s, volume cannot be found.", addVol, vg.Id))
			return nil, err
		}
		if addVolRef.GroupId != "" {
			return nil, fmt.Errorf("Cannot add volume %s to group %s beacuse it is already in group %s", addVolRef.Id, vg.Id, addVolRef.GroupId)
		}
		if addVolRef.ProfileId == "" {
			return nil, fmt.Errorf("Cannot add volume %s to group %s , volume has no profile.", addVolRef.Id, vg.Id)
		}
		if !utils.Contained(addVolRef.ProfileId, vg.Profiles) {
			return nil, fmt.Errorf("Cannot add volume %s to group %s , volume profile is not supported by the group.", addVolRef.Id, vg.Id)
		}
		if addVolRef.Status != model.VolumeAvailable && addVolRef.Status != model.VolumeInUse {
			return nil, fmt.Errorf("Cannot add volume %s to group %s beacuse volume is in invalid status %s", addVolRef.Id, vg.Id, addVolRef.Status)
		}
		if addVolRef.PoolId != vg.PoolId {
			return nil, fmt.Errorf("Cannot add volume %s to group %s , volume is not local to the pool of group.", addVolRef.Id, vg.Id)
		}

		addVolumesNew = append(addVolumesNew, addVolRef.Id)
	}

	return addVolumesNew, nil
}

func ValidateRemoveVolumes(ctx *c.Context, volumes []*model.VolumeSpec, removeVolumes []string, vg *model.VolumeGroupSpec) ([]string, error) {

	for _, v := range removeVolumes {
		for _, volume := range volumes {
			if v == volume.Id {
				if volume.Status != model.VolumeAvailable && volume.Status != model.VolumeInUse && volume.Status != model.VolumeError && volume.Status != model.VolumeErrorDeleting {
					return nil, fmt.Errorf("Cannot remove volume %s from group %s, volume is in invalid status %s", volume.Id, vg.Id, volume.Status)
				}
				break
			}

		}
	}
	for _, v := range removeVolumes {
		var available = false
		for _, volume := range volumes {
			if v == volume.Id {
				available = true
				break
			}
		}
		if available == false {
			return nil, fmt.Errorf("Cannot remove volume %s from group %s, volume is not in group ", v, vg.Id)
		}
	}

	return removeVolumes, nil
}

func DeleteVolumeGroupDBEntry(ctx *c.Context, volumeGroupId string) error {
	vg, err := db.C.GetVolumeGroup(ctx, volumeGroupId)
	if err != nil {
		return err
	}

	// If pool id is invalid, it would mean that volume group creation failed before the create method
	// in storage driver was called, and delete its db entry directly.
	_, err = db.C.GetDockByPoolId(ctx, vg.PoolId)
	if err != nil {
		if err := db.C.DeleteVolumeGroup(ctx, vg.Id); err != nil {
			log.Error("when delete volume group in db:", err)
			return err
		}
		return nil
	}

	//TODO DeleteVolumes tag is set by policy.
	deleteVolumes := true

	if deleteVolumes == false && vg.Status != model.VolumeGroupAvailable && vg.Status != model.VolumeGroupError {
		msg := fmt.Sprintf("The status of the Group must be available or error , group can be deleted. But current status is %s", vg.Status)
		log.Error(msg)
		return errors.New(msg)
	}

	if vg.GroupSnapshots != nil {
		msg := fmt.Sprintf("Group can not be deleted, because group has existing snapshots")
		log.Error(msg)
		return errors.New(msg)
	}

	volumes, err := db.C.ListVolumesByGroupId(ctx, vg.Id)
	if err != nil {
		return err
	}

	if len(volumes) > 0 && deleteVolumes == false {
		msg := fmt.Sprintf("Group %s still contains volumes. The deleteVolumes flag is required to delete it.", vg.Id)
		log.Error(msg)
		return errors.New(msg)
	}

	var volumesUpdate []*model.VolumeSpec
	for _, value := range volumes {
		if value.AttachStatus == model.VolumeAttached {
			msg := fmt.Sprintf("Volume %s in group %s is attached. Need to deach first.", value.Id, vg.Id)
			log.Error(msg)
			return errors.New(msg)
		}

		snapshots, err := db.C.ListSnapshotsByVolumeId(ctx, value.Id)
		if err != nil {
			return err
		}
		if len(snapshots) > 0 {
			msg := fmt.Sprintf("Volume %s in group still has snapshots", value.Id)
			log.Error(msg)
			return errors.New(msg)
		}

		volumesUpdate = append(volumesUpdate, &model.VolumeSpec{
			BaseModel: &model.BaseModel{
				Id: value.Id,
			},
			Status:  model.VolumeDeleting,
			GroupId: volumeGroupId,
		})
	}

	db.C.UpdateStatus(ctx, volumesUpdate, "")

	db.C.UpdateStatus(ctx, vg, model.VolumeGroupDeleting)
	//TODO Rpc call to delete group.
	controller.Brain.DeleteVolumeGroup(ctx, vg)

	return nil
}
