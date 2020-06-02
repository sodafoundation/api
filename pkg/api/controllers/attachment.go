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
	"github.com/sodafoundation/api/pkg/api/policy"
	c "github.com/sodafoundation/api/pkg/context"
	"github.com/sodafoundation/api/pkg/api/controllerclient"
	"github.com/sodafoundation/api/pkg/db"
	"github.com/sodafoundation/api/pkg/model"
	pb "github.com/sodafoundation/api/pkg/model/proto"
	"github.com/sodafoundation/api/pkg/utils"
	apiconfig "github.com/sodafoundation/api/pkg/utils/config"
	"github.com/sodafoundation/api/pkg/utils/constants"
)

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

	// Check if host exists
	host, err := db.C.GetHost(ctx, attachment.HostId)
	if err != nil {
		errMsg := fmt.Sprintf("get host failed in create volume attachment method: %v", err)
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// Check if volume exists and volume status is normal
	vol, err := db.C.GetVolume(ctx, attachment.VolumeId)
	if err != nil {
		errMsg := fmt.Sprintf("get volume failed in create volume attachment method: %v", err)
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	if !utils.Contains(host.AvailabilityZones, vol.AvailabilityZone) {
		errMsg := fmt.Sprintf("availability zone of volume: %s is not in the host availability zones: %v",
			vol.AvailabilityZone,
			host.AvailabilityZones)
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	if vol.Status == model.VolumeAvailable {
		db.UpdateVolumeStatus(ctx, db.C, vol.Id, model.VolumeAttaching)
	} else if vol.Status == model.VolumeInUse {
		if vol.MultiAttach {
			db.UpdateVolumeStatus(ctx, db.C, vol.Id, model.VolumeAttaching)
		} else {
			errMsg := "volume is already attached to one of the host. If you want to attach to multiple host, volume multiattach must be true"
			v.ErrorHandle(model.ErrorBadRequest, errMsg)
			return
		}
	} else {
		errMsg := "status of volume is available. It can be attached to host"
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// Set AccessProtocol
	pol, err := db.C.GetPool(ctx, vol.PoolId)
	if err != nil {
		msg := fmt.Sprintf("get pool failed in create volume attachment method: %v", err)
		log.Error(msg)
		return
	}
	var protocol = pol.Extras.IOConnectivity.AccessProtocol
	if protocol == "" {
		// Default protocol is iscsi
		protocol = constants.ISCSIProtocol
	}
	attachment.AccessProtocol = protocol

	// Set AttachMode, rw is a default setting
	if attachment.AttachMode != "ro" && attachment.AttachMode != "rw" {
		attachment.AttachMode = "rw"
	}
	attachment.Status = model.VolumeAttachCreating

	// NOTE:It will create a volume attachment entry into the database and initialize its status
	// as "creating". It will not wait for the real volume attachment creation to complete
	// and will return result immediately.
	result, err := db.C.CreateVolumeAttachment(ctx, &attachment)
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
	if err := v.CtrClient.Connect(apiconfig.CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}

	// // Note: In some protocols, there is no related initiator
	// var initiatorPort = ""
	// for _, e := range host.Initiators {
	// 	if e.Protocol == protocol {
	// 		initiatorPort = e.PortName
	// 		break
	// 	}
	// }
	var initiators []*pb.Initiator
	for _, e := range host.Initiators {
		initiator := pb.Initiator{
			PortName: e.PortName,
			Protocol: e.Protocol,
		}
		initiators = append(initiators, &initiator)
	}

	opt := &pb.CreateVolumeAttachmentOpts{
		Id:             result.Id,
		VolumeId:       result.VolumeId,
		PoolId:         vol.PoolId,
		AccessProtocol: protocol,
		HostInfo: &pb.HostInfo{
			OsType:     host.OsType,
			Ip:         host.IP,
			Host:       host.HostName,
			Initiators: initiators,
		},
		Metadata: vol.Metadata,
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

	// Check if attachment can be deleted
	validStatus := []string{model.VolumeAttachAvailable, model.VolumeAttachError,
		model.VolumeAttachErrorDeleting}
	if !utils.Contained(attachment.Status, validStatus) {
		errMsg := fmt.Sprintf("only the volume attachment with the status available, error, error_deleting can be deleted, the volume status is %s", attachment.Status)
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// If volume id is invalid, it would mean that volume attachment creation failed before the create method
	// in storage driver was called, and delete its db entry directly.
	vol, err := db.C.GetVolume(ctx, attachment.VolumeId)
	if err != nil {
		if err := db.C.DeleteVolumeAttachment(ctx, attachment.Id); err != nil {
			errMsg := fmt.Sprintf("failed to delete volume attachment: %s", err.Error())
			v.ErrorHandle(model.ErrorBadRequest, errMsg)
			return
		}
		v.SuccessHandle(StatusAccepted, nil)
		return
	}

	host, err := db.C.GetHost(ctx, attachment.HostId)
	if err != nil {
		errMsg := fmt.Sprintf("get host failed in delete volume attachment method: %v", err)
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	attachment.Status = model.VolumeAttachDeleting
	_, err = db.C.UpdateVolumeAttachment(ctx, attachment.Id, attachment)
	if err != nil {
		errMsg := fmt.Sprintf("failed to update volume attachment: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	v.SuccessHandle(StatusAccepted, nil)

	// NOTE:The real volume attachment deletion process.
	// Volume attachment deletion request is sent to the Dock. Dock will delete volume attachment from database
	// or update its status to "errorDeleting" if volume connection termination failed.
	if err := v.CtrClient.Connect(apiconfig.CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}

	var initiators []*pb.Initiator
	for _, e := range host.Initiators {
		initiator := pb.Initiator{
			PortName: e.PortName,
			Protocol: e.Protocol,
		}
		initiators = append(initiators, &initiator)
	}
	opt := &pb.DeleteVolumeAttachmentOpts{
		Id:             attachment.Id,
		VolumeId:       attachment.VolumeId,
		PoolId:         vol.PoolId,
		AccessProtocol: attachment.AccessProtocol,
		HostInfo: &pb.HostInfo{
			OsType:     host.OsType,
			Ip:         host.IP,
			Host:       host.HostName,
			Initiators: initiators,
		},
		Metadata: vol.Metadata,
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
