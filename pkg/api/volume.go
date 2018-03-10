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
This module implements a entry into the OpenSDS northbound service.

*/

package api

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/policy"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
)

type VolumePortal struct {
	beego.Controller
}

func (this *VolumePortal) CreateVolume() {
	if !policy.Authorize(this.Ctx, "volume:create") {
		return
	}
	var volume = model.VolumeSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&volume); err != nil {
		reason := fmt.Sprintf("Parse volume request body failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Call global controller variable to handle create volume request.
	result, err := controller.Brain.CreateVolume(c.GetContext(this.Ctx), &volume)
	if err != nil {
		reason := fmt.Sprintf("Create volume failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal volume created result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)
	this.Ctx.Output.Body(body)
	return
}

func (this *VolumePortal) ListVolumes() {
	if !policy.Authorize(this.Ctx, "volume:list") {
		return
	}
	// Call db api module to handle list volumes request.
	result, err := db.C.ListVolumes(c.GetContext(this.Ctx))
	if err != nil {
		reason := fmt.Sprintf("List volumes failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal volumes listed result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *VolumePortal) GetVolume() {
	if !policy.Authorize(this.Ctx, "volume:get") {
		return
	}
	id := this.Ctx.Input.Param(":volumeId")

	// Call db api module to handle get volume request.
	result, err := db.C.GetVolume(c.GetContext(this.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("Get volume failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal volume showed result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *VolumePortal) UpdateVolume() {
	if !policy.Authorize(this.Ctx, "volume:update") {
		return
	}
	var volume = model.VolumeSpec{
		BaseModel: &model.BaseModel{},
	}

	id := this.Ctx.Input.Param(":volumeId")
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&volume); err != nil {
		reason := fmt.Sprintf("Parse volume request body failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	volume.Id = id
	result, err := db.C.UpdateVolume(c.GetContext(this.Ctx), &volume)

	if err != nil {
		reason := fmt.Sprintf("Update volume failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal volume updated result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)

	return
}

// ExtendVolume ...
func (this *VolumePortal) ExtendVolume() {
	if !policy.Authorize(this.Ctx, "volume:extend") {
		return
	}
	var extendRequestBody = model.ExtendVolumeSpec{}

	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&extendRequestBody); err != nil {
		reason := fmt.Sprintf("Parse volume request body failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	id := this.Ctx.Input.Param(":volumeId")
	volume, err := db.C.GetVolume(c.GetContext(this.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("Get volume failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	if extendRequestBody.Extend.NewSize > volume.Size {
		volume.Size = extendRequestBody.Extend.NewSize
	} else {
		reason := fmt.Sprintf("Extend volume failed: new size(%d) <= old size(%d)",
			extendRequestBody.Extend.NewSize, volume.Size)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Call global controller variable to handle extend volume request.
	result, err := controller.Brain.ExtendVolume(c.GetContext(this.Ctx), volume)
	if err != nil {
		reason := fmt.Sprintf("Extend volume failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal volume extended result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)

	return
}

func (this *VolumePortal) DeleteVolume() {
	if !policy.Authorize(this.Ctx, "volume:delete") {
		return
	}
	id := this.Ctx.Input.Param(":volumeId")
	volume, err := db.C.GetVolume(c.GetContext(this.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("Get volume failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Call global controller variable to handle delete volume request.
	err = controller.Brain.DeleteVolume(c.GetContext(this.Ctx), volume)
	if err != nil {
		reason := fmt.Sprintf("Delete volume failed: %v", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)
	return
}

type VolumeAttachmentPortal struct {
	beego.Controller
}

func (this *VolumeAttachmentPortal) CreateVolumeAttachment() {
	if !policy.Authorize(this.Ctx, "volume:create_attachment") {
		return
	}
	var attachment = model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{},
	}

	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&attachment); err != nil {
		reason := fmt.Sprintf("Parse volume attachment request body failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Call global controller variable to handle create volume attachment request.
	result, err := controller.Brain.CreateVolumeAttachment(c.GetContext(this.Ctx), &attachment)
	if err != nil {
		reason := fmt.Sprintf("Create volume attachment failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal volume attachment created result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)
	this.Ctx.Output.Body(body)
	return
}

func (this *VolumeAttachmentPortal) ListVolumeAttachments() {
	if !policy.Authorize(this.Ctx, "volume:list_attachments") {
		return
	}
	volId := this.GetString("volumeId")

	result, err := db.C.ListVolumeAttachments(c.GetContext(this.Ctx), volId)
	if err != nil {
		reason := fmt.Sprintf("List volume attachments failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal volume attachments listed result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *VolumeAttachmentPortal) GetVolumeAttachment() {
	if !policy.Authorize(this.Ctx, "volume:get_attachment") {
		return
	}
	id := this.Ctx.Input.Param(":attachmentId")

	result, err := db.C.GetVolumeAttachment(c.GetContext(this.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("Get volume attachment failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal volume attachment showed result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *VolumeAttachmentPortal) UpdateVolumeAttachment() {
	if !policy.Authorize(this.Ctx, "volume:update_attachment") {
		return
	}
	var attachment = model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{},
	}
	id := this.Ctx.Input.Param(":attachmentId")

	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&attachment); err != nil {
		reason := fmt.Sprintf("Parse volume attachment request body failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}
	attachment.Id = id

	result, err := db.C.UpdateVolumeAttachment(c.GetContext(this.Ctx), id, &attachment)
	if err != nil {
		reason := fmt.Sprintf("Update volume attachment failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal volume attachment updated result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *VolumeAttachmentPortal) DeleteVolumeAttachment() {
	if !policy.Authorize(this.Ctx, "volume:delete_attachment") {
		return
	}
	id := this.Ctx.Input.Param(":attachmentId")
	attachment, err := db.C.GetVolumeAttachment(c.GetContext(this.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("Get volume attachment failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Call global controller variable to handle delete volume attachment request.
	err = controller.Brain.DeleteVolumeAttachment(c.GetContext(this.Ctx), attachment)
	if err != nil {
		reason := fmt.Sprintf("Delete volume attachment failed: %v", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)
	return
}

type VolumeSnapshotPortal struct {
	beego.Controller
}

func (this *VolumeSnapshotPortal) CreateVolumeSnapshot() {
	if !policy.Authorize(this.Ctx, "snapshot:create") {
		return
	}
	var snapshot = model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{},
	}

	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&snapshot); err != nil {
		reason := fmt.Sprintf("Parse volume snapshot request body failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Call global controller variable to handle create volume snapshot request.
	result, err := controller.Brain.CreateVolumeSnapshot(c.GetContext(this.Ctx), &snapshot)
	if err != nil {
		reason := fmt.Sprintf("Create volume snapshot failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal volume snapshot created result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)
	this.Ctx.Output.Body(body)
	return
}

func (this *VolumeSnapshotPortal) ListVolumeSnapshots() {
	if !policy.Authorize(this.Ctx, "snapshot:list") {
		return
	}
	result, err := db.C.ListVolumeSnapshots(c.GetContext(this.Ctx))
	if err != nil {
		reason := fmt.Sprintf("List volume snapshots failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal volume snapshots listed result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *VolumeSnapshotPortal) GetVolumeSnapshot() {
	if !policy.Authorize(this.Ctx, "snapshot:get") {
		return
	}
	id := this.Ctx.Input.Param(":snapshotId")

	result, err := db.C.GetVolumeSnapshot(c.GetContext(this.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("Get volume snapshot failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal volume snapshot showed result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *VolumeSnapshotPortal) UpdateVolumeSnapshot() {
	if !policy.Authorize(this.Ctx, "snapshot:update") {
		return
	}
	var snapshot = model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{},
	}

	id := this.Ctx.Input.Param(":snapshotId")

	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&snapshot); err != nil {
		reason := fmt.Sprintf("Parse volume snapshot request body failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}
	snapshot.Id = id

	result, err := db.C.UpdateVolumeSnapshot(c.GetContext(this.Ctx), id, &snapshot)
	if err != nil {
		reason := fmt.Sprintf("Update volume snapshot failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal volume snapshot updated result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *VolumeSnapshotPortal) DeleteVolumeSnapshot() {
	if !policy.Authorize(this.Ctx, "snapshot:delete") {
		return
	}
	id := this.Ctx.Input.Param(":snapshotId")

	snapshot, err := db.C.GetVolumeSnapshot(c.GetContext(this.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("Get volume snapshot failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Call global controller variable to handle delete volume snapshot request.
	err = controller.Brain.DeleteVolumeSnapshot(c.GetContext(this.Ctx), snapshot)
	if err != nil {
		reason := fmt.Sprintf("Delete volume snapshot failed: %v", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)
	return
}
