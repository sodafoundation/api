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
	"reflect"
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
		Size:             in.Size,
		AvailabilityZone: in.AvailabilityZone,
		Status:           model.VOLUME_CREATING,
		ProfileId:        in.ProfileId,
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

func CreateVolumeGroupDBEntry(ctx *c.Context, in *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
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
		Status:           model.VOLUMEGROUP_CREATING,
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
		vgNew.Status = model.VOLUMEGROUP_AVAILABLE
	} else {
		vgNew.Status = model.VOLUMEGROUP_UPDATING
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
		if addVolRef.Status != model.VOLUME_AVAILABLE && addVolRef.Status != model.VOLUME_IN_USE {
			return nil, fmt.Errorf("Cannot add volume %s to group %s beacuse volume is in invalid status %s", addVolRef.Id, vg.Id, addVolRef.Status)
		}

		addVolumesNew = append(addVolumesNew, addVolRef.Id)
	}

	return addVolumesNew, nil
}

func ValidateRemoveVolumes(ctx *c.Context, volumes []*model.VolumeSpec, removeVolumes []string, vg *model.VolumeGroupSpec) ([]string, error) {

	for _, v := range removeVolumes {
		for _, volume := range volumes {
			if v == volume.Id {
				if volume.Status != model.VOLUME_AVAILABLE && volume.Status != model.VOLUME_IN_USE && volume.Status != model.VOLUME_ERROR && volume.Status != model.VOLUEM_ERROR_DELETING {
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
	//TODO DeleteVolumes tag is set by policy.
	deleteVolumes := true

	if deleteVolumes == false && vg.Status != model.VOLUMEGROUP_AVAILABLE && vg.Status != model.VOLUMEGROUP_ERROR {
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
		if value.AttachStatus == model.VOLUME_ATTACHED {
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
			Status:  model.VOLUME_DELETING,
			GroupId: volumeGroupId,
		})
	}

	db.C.UpdateStatus(ctx, volumesUpdate, "")

	db.C.UpdateStatus(ctx, vg, model.VOLUMEGROUP_DELETING)
	//TODO Rpc call to delete group.
	controller.Brain.DeleteVolumeGroup(ctx, vg)
	return nil
}

func ListVolumeWithFilter(ctx *c.Context, filter map[string][]string) ([]*model.VolumeSpec, error) {
	volumes, err := db.C.ListVolumes(ctx)
	if err != nil {
		log.Error("List volumes failed: ", err)
		return nil, err
	}

	volumesSelected := Select(filter, volumes)

	v := reflect.ValueOf(volumesSelected)
	l := v.Len()

	var volList []*model.VolumeSpec

	for i := 0; i < l; i++ {
		volList = append(volList, v.Index(i).Interface().(*model.VolumeSpec))
	}

	var vol *model.VolumeSpec

	p := ParameterFilter(filter, len(volList), []string{"ID", "NAME", "STATUS", "AVAILABILITYZONE", "PROFILEID", "PROJECTID", "SIZE", "POOLID", "DESCRIPTION"})

	return vol.SortList(volList, p.sortKey, p.sortDir)[p.beginIdx:p.endIdx], nil
}

func ListDocksWithFilter(ctx *c.Context, filter map[string][]string) ([]*model.DockSpec, error) {
	docks, err := db.C.ListDocks(ctx)
	if err != nil {
		log.Error("List docks failed: ", err.Error())
		return nil, err
	}

	dcksSelected := Select(filter, docks)

	v := reflect.ValueOf(dcksSelected)
	l := v.Len()

	var dockList []*model.DockSpec

	for i := 0; i < l; i++ {
		dockList = append(dockList, v.Index(i).Interface().(*model.DockSpec))
	}

	var d *model.DockSpec

	p := ParameterFilter(filter, len(dockList), []string{"ID", "NAME", "ENDPOINT", "DRIVERNAME", "DESCRIPTION", "STATUS"})

	return d.SortList(dockList, p.sortKey, p.sortDir)[p.beginIdx:p.endIdx], nil
}

func ListPoolsWithFilter(ctx *c.Context, filter map[string][]string) ([]*model.StoragePoolSpec, error) {
	pools, err := db.C.ListPools(ctx)
	if err != nil {
		log.Error("List pools failed: ", err.Error())
		return nil, err
	}

	poolsSelected := Select(filter, pools)

	v := reflect.ValueOf(poolsSelected)
	l := v.Len()

	var poolList []*model.StoragePoolSpec

	for i := 0; i < l; i++ {
		poolList = append(poolList, v.Index(i).Interface().(*model.StoragePoolSpec))
	}

	var d *model.StoragePoolSpec

	p := ParameterFilter(filter, len(poolList), []string{"ID", "NAME", "STATUS", "AVAILABILITYZONE", "DOCKID", "DESCRIPTION"})

	return d.SortList(poolList, p.sortKey, p.sortDir)[p.beginIdx:p.endIdx], nil
}

func ListProfilesWithFilter(ctx *c.Context, filter map[string][]string) ([]*model.ProfileSpec, error) {
	profiles, err := db.C.ListProfiles(ctx)
	if err != nil {
		log.Error("List profiles failed: ", err)
		return nil, err
	}

	prfsSelected := Select(filter, profiles)

	v := reflect.ValueOf(prfsSelected)
	l := v.Len()

	var profList []*model.ProfileSpec

	for i := 0; i < l; i++ {
		profList = append(profList, v.Index(i).Interface().(*model.ProfileSpec))
	}

	var d *model.ProfileSpec

	p := ParameterFilter(filter, len(profList), []string{"ID", "NAME", "DESCRIPTION"})

	return d.SortList(profList, p.sortKey, p.sortDir)[p.beginIdx:p.endIdx], nil
}

func ListVolumeAttachmentsWithFilter(ctx *c.Context, filter map[string][]string) ([]*model.VolumeAttachmentSpec, error) {
	var volumeId string
	if v, ok := filter["VolumeId"]; ok {
		volumeId = v[0]
	}
	volumeAttachments, err := db.C.ListVolumeAttachments(ctx, volumeId)
	if err != nil {
		log.Error("List volumes failed: ", err)
		return nil, err
	}

	atcsSelected := Select(filter, volumeAttachments)

	v := reflect.ValueOf(atcsSelected)
	l := v.Len()

	var atcsList []*model.VolumeAttachmentSpec

	for i := 0; i < l; i++ {
		atcsList = append(atcsList, v.Index(i).Interface().(*model.VolumeAttachmentSpec))
	}

	var d *model.VolumeAttachmentSpec

	p := ParameterFilter(filter, len(atcsList), []string{"ID", "VOLUMEID", "STATUS", "USERID", "TENANTID"})

	return d.SortList(atcsList, p.sortKey, p.sortDir)[p.beginIdx:p.endIdx], nil
}

func ListVolumeSnapshotsWithFilter(ctx *c.Context, filter map[string][]string) ([]*model.VolumeSnapshotSpec, error) {
	volumeSnapshots, err := db.C.ListVolumeSnapshots(ctx)
	if err != nil {
		log.Error("List volumeSnapshots failed: ", err)
		return nil, err
	}

	snpsSelected := Select(filter, volumeSnapshots)

	v := reflect.ValueOf(snpsSelected)
	l := v.Len()

	var snpsList []*model.VolumeSnapshotSpec

	for i := 0; i < l; i++ {
		snpsList = append(snpsList, v.Index(i).Interface().(*model.VolumeSnapshotSpec))
	}

	var d *model.VolumeSnapshotSpec

	p := ParameterFilter(filter, len(snpsList), []string{"ID", "VOLUMEID", "STATUS", "USERID", "TENANTID"})

	return d.SortList(snpsList, p.sortKey, p.sortDir)[p.beginIdx:p.endIdx], nil
}

func ListVolumeGroupsWithFilter(ctx *c.Context, filter map[string][]string) ([]*model.VolumeGroupSpec, error) {
	volumeGroups, err := db.C.ListVolumeGroups(ctx)
	if err != nil {
		log.Error("List volumeGroups failed: ", err)
		return nil, err
	}

	volumeGroupsSelected := Select(filter, volumeGroups)

	v := reflect.ValueOf(volumeGroupsSelected)
	l := v.Len()

	var vgList []*model.VolumeGroupSpec

	for i := 0; i < l; i++ {
		vgList = append(vgList, v.Index(i).Interface().(*model.VolumeGroupSpec))
	}

	var d *model.VolumeGroupSpec

	p := ParameterFilter(filter, len(vgList), []string{"ID", "CREATEDAT", "NAME", "STATUS", "POOLID", "AVAILABILITYZONE", "USERID", "TENANTID"})

	return d.SortList(vgList, p.sortKey, p.sortDir)[p.beginIdx:p.endIdx], nil

}
