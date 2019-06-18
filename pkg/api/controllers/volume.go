// Copyright 2019 The OpenSDS Authors.
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
This module implements a entry into the OpenSDS northbound service.

*/

package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/policy"
	"github.com/opensds/opensds/pkg/api/util"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller/client"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	. "github.com/opensds/opensds/pkg/utils/config"
)

func NewVolumePortal() *VolumePortal {
	return &VolumePortal{
		CtrClient: client.NewClient(),
	}
}

type VolumePortal struct {
	BasePortal

	CtrClient client.Client
}

func (v *VolumePortal) CreateVolume() {
	if !policy.Authorize(v.Ctx, "volume:create") {
		return
	}
	ctx := c.GetContext(v.Ctx)
	var volume = model.VolumeSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(v.Ctx.Request.Body).Decode(&volume); err != nil {
		errMsg := fmt.Sprintf("parse volume request body failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// get profile
	var prf *model.ProfileSpec
	var err error
	if volume.ProfileId == "" {
		log.Warning("Use default profile when user doesn't specify profile.")
		prf, err = db.C.GetDefaultProfile(ctx)
		volume.ProfileId = prf.Id
	} else {
		prf, err = db.C.GetProfile(ctx, volume.ProfileId)
	}
	if err != nil {
		errMsg := fmt.Sprintf("get profile failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// NOTE:It will create a volume entry into the database and initialize its status
	// as "creating". It will not wait for the real volume creation to complete
	// and will return result immediately.
	result, err := util.CreateVolumeDBEntry(ctx, &volume)
	if err != nil {
		errMsg := fmt.Sprintf("create volume failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	log.V(8).Infof("create volume DB entry success %+v", result)

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusAccepted, body)

	// NOTE:The real volume creation process.
	// Volume creation request is sent to the Dock. Dock will update volume status to "available"
	// after volume creation is completed.
	if err := v.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer v.CtrClient.Close()

	opt := &pb.CreateVolumeOpts{
		Id:               result.Id,
		Name:             result.Name,
		Description:      result.Description,
		Size:             result.Size,
		AvailabilityZone: result.AvailabilityZone,
		// TODO: ProfileId will be removed later.
		ProfileId:         result.ProfileId,
		Profile:           prf.ToJson(),
		PoolId:            result.PoolId,
		SnapshotId:        result.SnapshotId,
		Metadata:          result.Metadata,
		SnapshotFromCloud: result.SnapshotFromCloud,
		Context:           ctx.ToJson(),
	}
	response, err := v.CtrClient.CreateVolume(context.Background(), opt)
	if err != nil {
		log.Error("create volume failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to create volume in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}

func (v *VolumePortal) ListVolumes() {
	if !policy.Authorize(v.Ctx, "volume:list") {
		return
	}
	// Call db api module to handle list volumes request.
	m, err := v.GetParameters()
	if err != nil {
		errMsg := fmt.Sprintf("list volumes failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	result, err := db.C.ListVolumesWithFilter(c.GetContext(v.Ctx), m)
	if err != nil {
		errMsg := fmt.Sprintf("list volumes failed: %s", err.Error())
		v.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusOK, body)

	return
}

func (v *VolumePortal) GetVolume() {
	if !policy.Authorize(v.Ctx, "volume:get") {
		return
	}
	id := v.Ctx.Input.Param(":volumeId")

	// Call db api module to handle get volume request.
	result, err := db.C.GetVolume(c.GetContext(v.Ctx), id)
	if err != nil {
		errMsg := fmt.Sprintf("volume %s not found: %s", id, err.Error())
		v.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusOK, body)

	return
}

func (v *VolumePortal) UpdateVolume() {
	if !policy.Authorize(v.Ctx, "volume:update") {
		return
	}
	var volume = model.VolumeSpec{
		BaseModel: &model.BaseModel{},
	}

	id := v.Ctx.Input.Param(":volumeId")
	if err := json.NewDecoder(v.Ctx.Request.Body).Decode(&volume); err != nil {
		errMsg := fmt.Sprintf("parse volume request body failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	volume.Id = id
	result, err := db.C.UpdateVolume(c.GetContext(v.Ctx), &volume)
	if err != nil {
		errMsg := fmt.Sprintf("update volume failed: %s", err.Error())
		v.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusOK, body)

	return
}

// ExtendVolume ...
func (v *VolumePortal) ExtendVolume() {
	if !policy.Authorize(v.Ctx, "volume:extend") {
		return
	}
	ctx := c.GetContext(v.Ctx)
	var extendRequestBody = model.ExtendVolumeSpec{}

	if err := json.NewDecoder(v.Ctx.Request.Body).Decode(&extendRequestBody); err != nil {
		errMsg := fmt.Sprintf("parse volume request body failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	id := v.Ctx.Input.Param(":volumeId")
	volume, err := db.C.GetVolume(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("volume %s not found: %s", id, err.Error())
		v.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	prf, err := db.C.GetProfile(ctx, volume.ProfileId)
	if err != nil {
		errMsg := fmt.Sprintf("extend volume failed: %v", err.Error())
		v.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// NOTE:It will update the the status of the volume waiting for expansion in
	// the database to "extending" and return the result immediately.
	result, err := util.ExtendVolumeDBEntry(ctx, id, &extendRequestBody)
	if err != nil {
		errMsg := fmt.Sprintf("extend volume failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusAccepted, body)

	// NOTE:The real volume extension process.
	// Volume extension request is sent to the Dock. Dock will update volume status to "available"
	// after volume extension is completed.
	if err = v.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer v.CtrClient.Close()

	opt := &pb.ExtendVolumeOpts{
		Id:       id,
		Size:     extendRequestBody.NewSize,
		Metadata: result.Metadata,
		Context:  ctx.ToJson(),
		Profile:  prf.ToJson(),
	}
	response, err := v.CtrClient.ExtendVolume(context.Background(), opt)
	if err != nil {
		log.Error("extend volume failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to extend volume in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}

func (v *VolumePortal) DeleteVolume() {
	if !policy.Authorize(v.Ctx, "volume:delete") {
		return
	}
	ctx := c.GetContext(v.Ctx)

	var err error
	id := v.Ctx.Input.Param(":volumeId")
	volume, err := db.C.GetVolume(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("volume %s not found: %s", id, err.Error())
		v.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// If profileId or poolId of the volume doesn't exist, it would mean that
	// the volume provisioning operation failed before the create method in
	// storage driver was called, therefore the volume entry should be deleted
	// from db directly.
	if volume.ProfileId == "" || volume.PoolId == "" {
		if err := db.C.DeleteVolume(ctx, volume.Id); err != nil {
			errMsg := fmt.Sprintf("delete volume failed: %v", err.Error())
			v.ErrorHandle(model.ErrorInternalServer, errMsg)
			return
		}
		v.SuccessHandle(StatusAccepted, nil)
		return
	}

	prf, err := db.C.GetProfile(ctx, volume.ProfileId)
	if err != nil {
		errMsg := fmt.Sprintf("delete volume failed: %v", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// NOTE:It will update the the status of the volume waiting for deletion in
	// the database to "deleting" and return the result immediately.
	if err = util.DeleteVolumeDBEntry(ctx, volume); err != nil {
		errMsg := fmt.Sprintf("delete volume failed: %v", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	v.SuccessHandle(StatusAccepted, nil)

	// NOTE:The real volume deletion process.
	// Volume deletion request is sent to the Dock. Dock will delete volume from driver
	// and database or update volume status to "errorDeleting" if deletion from driver faild.
	if err := v.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer v.CtrClient.Close()

	opt := &pb.DeleteVolumeOpts{
		Id:        volume.Id,
		ProfileId: volume.ProfileId,
		PoolId:    volume.PoolId,
		Metadata:  volume.Metadata,
		Context:   ctx.ToJson(),
		Profile:   prf.ToJson(),
	}
	response, err := v.CtrClient.DeleteVolume(context.Background(), opt)
	if err != nil {
		log.Error("delete volume failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to delete volume in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}

func NewVolumeAttachmentPortal() *VolumeAttachmentPortal {
	return &VolumeAttachmentPortal{
		CtrClient: client.NewClient(),
	}
}

type VolumeAttachmentPortal struct {
	BasePortal

	CtrClient client.Client
}

func (v *VolumeAttachmentPortal) CreateVolumeAttachment() {
	if !policy.Authorize(v.Ctx, "volume:create_attachment") {
		return
	}
	ctx := c.GetContext(v.Ctx)
	var attachment = model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{},
	}

	if err := json.NewDecoder(v.Ctx.Request.Body).Decode(&attachment); err != nil {
		errMsg := fmt.Sprintf("parse volume attachment request body failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// NOTE:It will create a volume attachment entry into the database and initialize its status
	// as "creating". It will not wait for the real volume attachment creation to complete
	// and will return result immediately.
	result, err := util.CreateVolumeAttachmentDBEntry(ctx, &attachment)
	if err != nil {
		errMsg := fmt.Sprintf("create volume attachment failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusAccepted, body)

	// NOTE:The real volume attachment creation process.
	// Volume attachment creation request is sent to the Dock. Dock will update volume attachment status to "available"
	// after volume attachment creation is completed.
	if err := v.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer v.CtrClient.Close()

	opt := &pb.CreateVolumeAttachmentOpts{
		Id:       result.Id,
		VolumeId: result.VolumeId,
		HostInfo: &pb.HostInfo{
			Platform:  result.Platform,
			OsType:    result.OsType,
			Ip:        result.Ip,
			Host:      result.Host,
			Initiator: result.Initiator,
		},
		Metadata: result.Metadata,
		Context:  ctx.ToJson(),
	}
	response, err := v.CtrClient.CreateVolumeAttachment(context.Background(), opt)
	if err != nil {
		log.Error("create volume attachment failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to create volume attachment in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}

func (v *VolumeAttachmentPortal) ListVolumeAttachments() {
	if !policy.Authorize(v.Ctx, "volume:list_attachments") {
		return
	}

	m, err := v.GetParameters()
	if err != nil {
		errMsg := fmt.Sprintf("list volume attachments failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	result, err := db.C.ListVolumeAttachmentsWithFilter(c.GetContext(v.Ctx), m)
	if err != nil {
		errMsg := fmt.Sprintf("list volume attachments failed: %s", err.Error())
		v.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusOK, body)

	return
}

func (v *VolumeAttachmentPortal) GetVolumeAttachment() {
	if !policy.Authorize(v.Ctx, "volume:get_attachment") {
		return
	}
	id := v.Ctx.Input.Param(":attachmentId")

	result, err := db.C.GetVolumeAttachment(c.GetContext(v.Ctx), id)
	if err != nil {
		errMsg := fmt.Sprintf("volume attachment %s not found: %s", id, err.Error())
		v.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusOK, body)

	return
}

func (v *VolumeAttachmentPortal) UpdateVolumeAttachment() {
	if !policy.Authorize(v.Ctx, "volume:update_attachment") {
		return
	}
	var attachment = model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{},
	}
	id := v.Ctx.Input.Param(":attachmentId")

	if err := json.NewDecoder(v.Ctx.Request.Body).Decode(&attachment); err != nil {
		errMsg := fmt.Sprintf("parse volume attachment request body failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}
	attachment.Id = id

	result, err := db.C.UpdateVolumeAttachment(c.GetContext(v.Ctx), id, &attachment)
	if err != nil {
		errMsg := fmt.Sprintf("update volume attachment failed: %s", err.Error())
		v.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusOK, body)

	return
}

func (v *VolumeAttachmentPortal) DeleteVolumeAttachment() {
	if !policy.Authorize(v.Ctx, "volume:delete_attachment") {
		return
	}

	ctx := c.GetContext(v.Ctx)

	id := v.Ctx.Input.Param(":attachmentId")
	attachment, err := db.C.GetVolumeAttachment(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("volume attachment %s not found: %s", id, err.Error())
		v.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}
	// NOTE:It will not wait for the real volume attachment deletion to complete
	// and will return ok immediately.
	if err = util.DeleteVolumeAttachmentDBEntry(ctx, attachment); err != nil {
		errMsg := fmt.Sprintf("delete volume attachment failed: %v", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	v.SuccessHandle(StatusAccepted, nil)

	// NOTE:The real volume attachment deletion process.
	// Volume attachment deletion request is sent to the Dock. Dock will delete volume attachment from database
	// or update its status to "errorDeleting" if volume connection termination failed.
	if err := v.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer v.CtrClient.Close()

	opt := &pb.DeleteVolumeAttachmentOpts{
		Id:             attachment.Id,
		VolumeId:       attachment.VolumeId,
		AccessProtocol: attachment.AccessProtocol,
		HostInfo: &pb.HostInfo{
			Platform:  attachment.Platform,
			OsType:    attachment.OsType,
			Ip:        attachment.Ip,
			Host:      attachment.Host,
			Initiator: attachment.Initiator,
		},
		Metadata: attachment.Metadata,
		Context:  ctx.ToJson(),
	}
	response, err := v.CtrClient.DeleteVolumeAttachment(context.Background(), opt)
	if err != nil {
		log.Error("delete volume attachment failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to delete volume attachment in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}

func NewVolumeSnapshotPortal() *VolumeSnapshotPortal {
	return &VolumeSnapshotPortal{
		CtrClient: client.NewClient(),
	}
}

type VolumeSnapshotPortal struct {
	BasePortal

	CtrClient client.Client
}

func (v *VolumeSnapshotPortal) CreateVolumeSnapshot() {
	if !policy.Authorize(v.Ctx, "snapshot:create") {
		return
	}
	ctx := c.GetContext(v.Ctx)
	var snapshot = model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{},
	}

	if err := json.NewDecoder(v.Ctx.Request.Body).Decode(&snapshot); err != nil {
		errMsg := fmt.Sprintf("parse volume snapshot request body failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// get profile
	var prf *model.ProfileSpec
	var err error

	// If user doesn't specified profile, using profile derived form volume
	if len(snapshot.ProfileId) == 0 {
		log.Warning("User doesn't specified profile id, using profile derived form volume")
		vol, err := db.C.GetVolume(ctx, snapshot.VolumeId)
		if err != nil {
			v.ErrorHandle(model.ErrorBadRequest, err.Error())
			return
		}
		snapshot.ProfileId = vol.ProfileId
	}
	prf, err = db.C.GetProfile(ctx, snapshot.ProfileId)

	if err != nil {
		errMsg := fmt.Sprintf("get profile failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// NOTE:It will create a volume snapshot entry into the database and initialize its status
	// as "creating". It will not wait for the real volume snapshot creation to complete
	// and will return result immediately.
	result, err := util.CreateVolumeSnapshotDBEntry(ctx, &snapshot)
	if err != nil {
		errMsg := fmt.Sprintf("create volume snapshot failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusAccepted, body)

	// NOTE:The real volume snapshot creation process.
	// Volume snapshot creation request is sent to the Dock. Dock will update volume snapshot status to "available"
	// after volume snapshot creation complete.
	if err := v.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer v.CtrClient.Close()

	opt := &pb.CreateVolumeSnapshotOpts{
		Id:          result.Id,
		Name:        result.Name,
		Description: result.Description,
		VolumeId:    result.VolumeId,
		Size:        result.Size,
		Metadata:    result.Metadata,
		Context:     ctx.ToJson(),
		Profile:     prf.ToJson(),
	}
	response, err := v.CtrClient.CreateVolumeSnapshot(context.Background(), opt)
	if err != nil {
		log.Error("create volume snapthot failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to create volume snapshot in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}

func (v *VolumeSnapshotPortal) ListVolumeSnapshots() {
	if !policy.Authorize(v.Ctx, "snapshot:list") {
		return
	}
	m, err := v.GetParameters()
	if err != nil {
		errMsg := fmt.Sprintf("list volume snapshots failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	result, err := db.C.ListVolumeSnapshotsWithFilter(c.GetContext(v.Ctx), m)
	if err != nil {
		errMsg := fmt.Sprintf("list volume snapshots failed: %s", err.Error())
		v.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusOK, body)

	return
}

func (v *VolumeSnapshotPortal) GetVolumeSnapshot() {
	if !policy.Authorize(v.Ctx, "snapshot:get") {
		return
	}
	id := v.Ctx.Input.Param(":snapshotId")

	result, err := db.C.GetVolumeSnapshot(c.GetContext(v.Ctx), id)
	if err != nil {
		errMsg := fmt.Sprintf("volume snapshot %s not found: %s", id, err.Error())
		v.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusOK, body)

	return
}

func (v *VolumeSnapshotPortal) UpdateVolumeSnapshot() {
	if !policy.Authorize(v.Ctx, "snapshot:update") {
		return
	}
	var snapshot = model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{},
	}

	id := v.Ctx.Input.Param(":snapshotId")

	if err := json.NewDecoder(v.Ctx.Request.Body).Decode(&snapshot); err != nil {
		errMsg := fmt.Sprintf("parse volume snapshot request body failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}
	snapshot.Id = id

	result, err := db.C.UpdateVolumeSnapshot(c.GetContext(v.Ctx), id, &snapshot)
	if err != nil {
		errMsg := fmt.Sprintf("update volume snapshot failed: %s", err.Error())
		v.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	v.SuccessHandle(StatusOK, body)

	return
}

func (v *VolumeSnapshotPortal) DeleteVolumeSnapshot() {
	if !policy.Authorize(v.Ctx, "snapshot:delete") {
		return
	}
	ctx := c.GetContext(v.Ctx)
	id := v.Ctx.Input.Param(":snapshotId")

	snapshot, err := db.C.GetVolumeSnapshot(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("volume snapshot %s not found: %s", id, err.Error())
		v.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	prf, err := db.C.GetProfile(ctx, snapshot.ProfileId)
	if err != nil {
		errMsg := fmt.Sprintf("delete snapshot failed: %v", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// NOTE:It will update the the status of the volume snapshot waiting for deletion in
	// the database to "deleting" and return the result immediately.
	err = util.DeleteVolumeSnapshotDBEntry(ctx, snapshot)
	if err != nil {
		errMsg := fmt.Sprintf("delete volume snapshot failed: %v", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	v.SuccessHandle(StatusAccepted, nil)

	// NOTE:The real volume snapshot deletion process.
	// Volume snapshot deletion request is sent to the Dock. Dock will delete volume snapshot from driver and
	// database or update its status to "errorDeleting" if volume snapshot deletion from driver failed.
	if err := v.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer v.CtrClient.Close()

	opt := &pb.DeleteVolumeSnapshotOpts{
		Id:       snapshot.Id,
		VolumeId: snapshot.VolumeId,
		Metadata: snapshot.Metadata,
		Context:  ctx.ToJson(),
		Profile:  prf.ToJson(),
	}
	response, err := v.CtrClient.DeleteVolumeSnapshot(context.Background(), opt)
	if err != nil {
		log.Error("delete volume snapthot failed in controller service:", err)
		return
	}
	if errorMsg := response.GetError(); errorMsg != nil {
		log.Errorf("failed to delete volume snapshot in controller, code: %v, message: %v",
			errorMsg.GetCode(), errorMsg.GetDescription())
		return
	}

	return
}
