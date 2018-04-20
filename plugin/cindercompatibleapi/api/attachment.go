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
This module implements a entry into the OpenSDS northbound service.

*/

package api

import (
	"encoding/json"
	"fmt"

	log "github.com/golang/glog"
	OpenSDSAPI "github.com/opensds/opensds/pkg/api"
	"github.com/opensds/opensds/pkg/api/policy"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/plugin/cindercompatibleapi/cindermodel"
	"github.com/opensds/opensds/plugin/cindercompatibleapi/converter"
)

// AttachmentPortal ...
type AttachmentPortal struct {
	OpenSDSAPI.BasePortal
}

// DeleteAttachment ...
func (portal *AttachmentPortal) DeleteAttachment() {
	if !policy.Authorize(portal.Ctx, "attachment:delete") {
		return
	}

	id := portal.Ctx.Input.Param(":attachmentId")
	attachment := model.VolumeAttachmentSpec{}
	err := client.DeleteVolumeAttachment(id, &attachment)

	if err != nil {
		reason := fmt.Sprintf("Delete attachment failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(StatusOK)
	return
}

// GetAttachment ...
func (portal *AttachmentPortal) GetAttachment() {
	if !policy.Authorize(portal.Ctx, "attachment:get") {
		return
	}

	id := portal.Ctx.Input.Param(":attachmentId")
	attachment, err := client.GetVolumeAttachment(id)

	if err != nil {
		reason := fmt.Sprintf("Show attachment details failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.ShowAttachmentResp(attachment)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Show attachment details, marshal result failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// ListAttachmentsDetail ...
func (portal *AttachmentPortal) ListAttachmentsDetail() {
	if !policy.Authorize(portal.Ctx, "attachment:list_detail") {
		return
	}

	attachments, err := client.ListVolumeAttachments()
	if err != nil {
		reason := fmt.Sprintf("List attachments with details failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.ListAttachmentDetailResp(attachments)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("List attachments with details, marshal result failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// ListAttachment ...
func (portal *AttachmentPortal) ListAttachment() {
	if !policy.Authorize(portal.Ctx, "attachment:list") {
		return
	}

	attachments, err := client.ListVolumeAttachments()
	if err != nil {
		reason := fmt.Sprintf("List attachments failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.ListAttachmentResp(attachments)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("List attachments, marshal result failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// CreateAttachment ...
func (portal *AttachmentPortal) CreateAttachment() {
	if !policy.Authorize(portal.Ctx, "attachment:create") {
		return
	}
	var cinderReq = cindermodel.CreateAttachmentReqSpec{}

	if err := json.NewDecoder(portal.Ctx.Request.Body).Decode(&cinderReq); err != nil {
		reason := fmt.Sprintf("Create attachment, parse request body failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	attachment := converter.CreateAttachmentReq(&cinderReq)
	attachment, err := client.CreateVolumeAttachment(attachment)

	if err != nil {
		reason := fmt.Sprintf("Create attachment failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.CreateAttachmentResp(attachment)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Create attachment, marshal result failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// UpdateAttachment ...
func (portal *AttachmentPortal) UpdateAttachment() {
	if !policy.Authorize(portal.Ctx, "attachment:update") {
		return
	}

	id := portal.Ctx.Input.Param(":attachmentId")
	var cinderReq = cindermodel.UpdateAttachmentReqSpec{}

	if err := json.NewDecoder(portal.Ctx.Request.Body).Decode(&cinderReq); err != nil {
		reason := fmt.Sprintf("Update an attachment, parse request body failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	attachment := converter.UpdateAttachmentReq(&cinderReq)
	attachment, err := client.UpdateVolumeAttachment(id, attachment)

	if err != nil {
		reason := fmt.Sprintf("Update an attachment failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.UpdateAttachmentResp(attachment)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Update an attachment, marshal result failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(StatusOK)
	portal.Ctx.Output.Body(body)
	return
}
