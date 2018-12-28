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

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/policy"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller/client"
	pb "github.com/opensds/opensds/pkg/controller/proto"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/pkg/utils/config"
	"golang.org/x/net/context"
)

type VolumePortal struct {
	BasePortal
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
	// NOTE:It will create a volume entry into the database and initialize its status
	// as "creating". It will not wait for the real volume creation to complete
	// and will return result immediately.
	result, err := CreateVolumeDBEntry(c.GetContext(this.Ctx), &volume)
	if err != nil {
		reason := fmt.Sprintf("Create volume failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	this.Ctx.Output.Body(body)
	this.Ctx.Output.SetStatus(StatusAccepted)

	// NOTE:The real volume creation process.
	// Volume creation request is sent to the Dock. Dock will update volume status to "available"
	// after volume creation is completed.
	ctrClient := client.NewClient()
	if err := ctrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer ctrClient.Close()

	opt := &pb.CreateVolumeOpts{
		Message: string(body),
	}
	if _, err = ctrClient.CreateVolume(context.Background(), opt); err != nil {
		log.Error("Create volume failed in controller service:", err)
		return
	}

	return
}

func (this *VolumePortal) ListVolumes() {
	if !policy.Authorize(this.Ctx, "volume:list") {
		return
	}
	// Call db api module to handle list volumes request.
	m, err := this.GetParameters()
	if err != nil {
		reason := fmt.Sprintf("List volumes failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result, err := db.C.ListVolumesWithFilter(c.GetContext(this.Ctx), m)
	if err != nil {
		reason := fmt.Sprintf("List volumes failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	this.Ctx.Output.Body(body)
	this.Ctx.Output.SetStatus(StatusOK)

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
	body, _ := json.Marshal(result)
	this.Ctx.Output.Body(body)
	this.Ctx.Output.SetStatus(StatusOK)

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
	body, _ := json.Marshal(result)
	this.Ctx.Output.Body(body)
	this.Ctx.Output.SetStatus(StatusOK)

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
	// NOTE:It will update the the status of the volume waiting for expansion in
	// the database to "extending" and return the result immediately.
	result, err := ExtendVolumeDBEntry(c.GetContext(this.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("Extend volume failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	this.Ctx.Output.Body(body)
	this.Ctx.Output.SetStatus(StatusAccepted)

	// NOTE:The real volume extension process.
	// Volume extension request is sent to the Dock. Dock will update volume status to "available"
	// after volume extension is completed.
	ctrClient := client.NewClient()
	if err := ctrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer ctrClient.Close()

	body, _ = json.Marshal(&extendRequestBody)
	opt := &pb.ExtendVolumeOpts{
		Id:      id,
		Message: string(body),
	}
	if _, err = ctrClient.ExtendVolume(context.Background(), opt); err != nil {
		log.Error("Extend volume failed in controller service:", err)
		return
	}

	return
}

func (this *VolumePortal) DeleteVolume() {
	if !policy.Authorize(this.Ctx, "volume:delete") {
		return
	}
	var err error
	id := this.Ctx.Input.Param(":volumeId")
	volume, err := db.C.GetVolume(c.GetContext(this.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("Get volume failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// NOTE:It will update the the status of the volume waiting for deletion in
	// the database to "deleting" and return the result immediately.
	if err = DeleteVolumeDBEntry(c.GetContext(this.Ctx), volume); err != nil {
		reason := fmt.Sprintf("Delete volume failed: %v", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}
	this.Ctx.Output.SetStatus(StatusAccepted)

	// NOTE:The real volume deletion process.
	// Volume deletion request is sent to the Dock. Dock will delete volume from driver
	// and database or update volume status to "errorDeleting" if deletion from driver faild.
	ctrClient := client.NewClient()
	if err := ctrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer ctrClient.Close()

	body, _ := json.Marshal(volume)
	opt := &pb.DeleteVolumeOpts{
		Message: string(body),
	}
	if _, err = ctrClient.DeleteVolume(context.Background(), opt); err != nil {
		log.Error("Delete volume failed in controller service:", err)
		return
	}

	return
}

type VolumeAttachmentPortal struct {
	BasePortal
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

	// NOTE:It will create a volume attachment entry into the database and initialize its status
	// as "creating". It will not wait for the real volume attachment creation to complete
	// and will return result immediately.
	result, err := CreateVolumeAttachmentDBEntry(c.GetContext(this.Ctx), &attachment)
	if err != nil {
		reason := fmt.Sprintf("Create volume attachment failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	this.Ctx.Output.Body(body)
	this.Ctx.Output.SetStatus(StatusAccepted)

	// NOTE:The real volume attachment creation process.
	// Volume attachment creation request is sent to the Dock. Dock will update volume attachment status to "available"
	// after volume attachment creation is completed.
	ctrClient := client.NewClient()
	if err := ctrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer ctrClient.Close()

	opt := &pb.CreateVolumeAttachmentOpts{
		Message: string(body),
	}
	if _, err = ctrClient.CreateVolumeAttachment(context.Background(), opt); err != nil {
		log.Error("Create volume attachment failed in controller service:", err)
		return
	}

	return
}

func (this *VolumeAttachmentPortal) ListVolumeAttachments() {
	if !policy.Authorize(this.Ctx, "volume:list_attachments") {
		return
	}

	m, err := this.GetParameters()
	if err != nil {
		reason := fmt.Sprintf("List volume attachments failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result, err := db.C.ListVolumeAttachmentsWithFilter(c.GetContext(this.Ctx), m)
	if err != nil {
		reason := fmt.Sprintf("List volume attachments failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	this.Ctx.Output.Body(body)
	this.Ctx.Output.SetStatus(StatusOK)

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
	body, _ := json.Marshal(result)
	this.Ctx.Output.Body(body)
	this.Ctx.Output.SetStatus(StatusOK)

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
	body, _ := json.Marshal(result)
	this.Ctx.Output.Body(body)
	this.Ctx.Output.SetStatus(StatusOK)

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
	// NOTE:It will not wait for the real volume attachment deletion to complete
	// and will return ok immediately.
	this.Ctx.Output.SetStatus(StatusAccepted)

	// NOTE:The real volume attachment deletion process.
	// Volume attachment deletion request is sent to the Dock. Dock will delete volume attachment from database
	// or update its status to "errorDeleting" if volume connection termination failed.
	ctrClient := client.NewClient()
	if err := ctrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer ctrClient.Close()

	body, _ := json.Marshal(attachment)
	opt := &pb.DeleteVolumeAttachmentOpts{
		Message: string(body),
	}
	if _, err = ctrClient.DeleteVolumeAttachment(context.Background(), opt); err != nil {
		log.Error("Delete volume attachment failed in controller service:", err)
		return
	}

	return
}

type VolumeSnapshotPortal struct {
	BasePortal
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

	// NOTE:It will create a volume snapshot entry into the database and initialize its status
	// as "creating". It will not wait for the real volume snapshot creation to complete
	// and will return result immediately.
	result, err := CreateVolumeSnapshotDBEntry(c.GetContext(this.Ctx), &snapshot)
	if err != nil {
		reason := fmt.Sprintf("Create volume snapshot failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	this.Ctx.Output.Body(body)
	this.Ctx.Output.SetStatus(StatusAccepted)

	// NOTE:The real volume snapshot creation process.
	// Volume snapshot creation request is sent to the Dock. Dock will update volume snapshot status to "available"
	// after volume snapshot creation complete.
	ctrClient := client.NewClient()
	if err := ctrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer ctrClient.Close()

	opt := &pb.CreateVolumeSnapshotOpts{
		Message: string(body),
	}
	if _, err = ctrClient.CreateVolumeSnapshot(context.Background(), opt); err != nil {
		log.Error("Create volume snapthot failed in controller service:", err)
		return
	}

	return
}

func (this *VolumeSnapshotPortal) ListVolumeSnapshots() {
	if !policy.Authorize(this.Ctx, "snapshot:list") {
		return
	}
	m, err := this.GetParameters()
	if err != nil {
		reason := fmt.Sprintf("List volume snapshots failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result, err := db.C.ListVolumeSnapshotsWithFilter(c.GetContext(this.Ctx), m)
	if err != nil {
		reason := fmt.Sprintf("List volume snapshots failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, _ := json.Marshal(result)
	this.Ctx.Output.Body(body)
	this.Ctx.Output.SetStatus(StatusOK)

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
	body, _ := json.Marshal(result)
	this.Ctx.Output.Body(body)
	this.Ctx.Output.SetStatus(StatusOK)

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
	body, _ := json.Marshal(result)
	this.Ctx.Output.Body(body)
	this.Ctx.Output.SetStatus(StatusOK)

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

	// NOTE:It will update the the status of the volume snapshot waiting for deletion in
	// the database to "deleting" and return the result immediately.
	err = DeleteVolumeSnapshotDBEntry(c.GetContext(this.Ctx), snapshot)
	if err != nil {
		reason := fmt.Sprintf("Delete volume snapshot failed: %v", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// NOTE:The real volume snapshot deletion process.
	// Volume snapshot deletion request is sent to the Dock. Dock will delete volume snapshot from driver and
	// database or update its status to "errorDeleting" if volume snapshot deletion from driver failed.
	ctrClient := client.NewClient()
	if err := ctrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("When connecting controller client:", err)
		return
	}
	defer ctrClient.Close()

	body, _ := json.Marshal(snapshot)
	opt := &pb.DeleteVolumeSnapshotOpts{
		Message: string(body),
	}
	if _, err = ctrClient.DeleteVolumeSnapshot(context.Background(), opt); err != nil {
		log.Error("Delete volume snapthot failed in controller service:", err)
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)
	return
}
