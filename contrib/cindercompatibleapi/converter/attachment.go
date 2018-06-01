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

package converter

import (
	"github.com/opensds/opensds/pkg/model"
)

// ExtraSpec ...
type ExtraSpec map[string]interface{}

// ConnectionInfo ...
type ConnectionInfo struct {
	DriverVolumeType string                 `json:"driver_volume_type,omitempty"`
	ConnectionData   map[string]interface{} `json:"data,omitempty"`
}

// Connector ...
type Connector struct {
	Initiator  string `json:"initiator,omitempty"`
	IP         string `json:"ip,omitempty"`
	Platform   string `json:"platform,omitempty"`
	Host       string `json:"host,omitempty"`
	OsType     string `json:"os_type,omitempty"`
	Multipath  bool   `json:"multipath,omitempty"`
	Mountpoint string `json:"mountpoint,omitempty"`
	Mode       string `json:"mode,omitempty"`
}

// *******************Show attachment details*******************

// ShowAttachmentRespSpec ...
type ShowAttachmentRespSpec struct {
	VolumeAttachment ShowRespAttachment `json:"attachment,omitempty"`
}

// ShowRespAttachment ...
type ShowRespAttachment struct {
	Status         string `json:"status"`
	DetachedAt     string `json:"detached_at,omitempty"`
	ConnectionInfo `json:"connection_info"`
	AttachedAt     string `json:"attached_at,omitempty"`
	AttachMode     string `json:"attach_mode,omitempty"`
	Instance       string `json:"instance,omitempty"`
	VolumeID       string `json:"volume_id"`
	ID             string `json:"id"`
}

// ShowAttachmentResp ...
func ShowAttachmentResp(attachment *model.VolumeAttachmentSpec) *ShowAttachmentRespSpec {
	resp := ShowAttachmentRespSpec{}
	resp.VolumeAttachment.Status = attachment.Status
	resp.VolumeAttachment.ConnectionInfo.DriverVolumeType = attachment.ConnectionInfo.DriverVolumeType
	resp.VolumeAttachment.ConnectionInfo.ConnectionData = attachment.ConnectionInfo.ConnectionData
	//resp.VolumeAttachment.AttachedAt = attachment.Mountpoint
	resp.VolumeAttachment.Instance = attachment.Metadata["instance_uuid"]
	resp.VolumeAttachment.VolumeID = attachment.VolumeId
	resp.VolumeAttachment.ID = attachment.BaseModel.Id

	return &resp
}

// *******************List attachments with details*******************

// ListAttachmentsDetailsRespSpec ...
type ListAttachmentsDetailsRespSpec struct {
	Attachments []ListRespAttachmentDetails `json:"attachments"`
}

// ListRespAttachmentDetails ...
type ListRespAttachmentDetails struct {
	Status         string         `json:"status"`
	DetachedAt     string         `json:"detached_at,omitempty"`
	ConnectionInfo ConnectionInfo `json:"connection_info"`
	AttachedAt     string         `json:"attached_at,omitempty"`
	AttachMode     string         `json:"attach_mode,omitempty"`
	Instance       string         `json:"instance,omitempty"`
	VolumeID       string         `json:"volume_id"`
	ID             string         `json:"id"`
}

// ListAttachmentsDetailsResp ...
func ListAttachmentsDetailsResp(attachments []*model.VolumeAttachmentSpec) *ListAttachmentsDetailsRespSpec {
	var resp ListAttachmentsDetailsRespSpec
	cinderAttachment := ListRespAttachmentDetails{}

	if 0 == len(attachments) {
		// Even if the number is 0, it must return {"attachments":[]}
		resp.Attachments = make([]ListRespAttachmentDetails, 0, 0)
	} else {
		for _, attachment := range attachments {
			cinderAttachment.Status = attachment.Status
			cinderAttachment.ConnectionInfo.DriverVolumeType = attachment.ConnectionInfo.DriverVolumeType
			cinderAttachment.ConnectionInfo.ConnectionData = attachment.ConnectionInfo.ConnectionData
			//cinderAttachment.AttachedAt = attachment.Mountpoint
			cinderAttachment.Instance = attachment.Metadata["instance_uuid"]
			cinderAttachment.VolumeID = attachment.VolumeId
			cinderAttachment.ID = attachment.Id

			resp.Attachments = append(resp.Attachments, cinderAttachment)
		}
	}

	return &resp
}

// *******************List attachments*******************

// ListAttachmentsRespSpec ...
type ListAttachmentsRespSpec struct {
	Attachments []ListRespAttachment `json:"attachments"`
}

// ListRespAttachment ...
type ListRespAttachment struct {
	Status   string `json:"status"`
	Instance string `json:"instance,omitempty"`
	VolumeID string `json:"volume_id"`
	ID       string `json:"id"`
}

// ListAttachmentsResp ...
func ListAttachmentsResp(attachments []*model.VolumeAttachmentSpec) *ListAttachmentsRespSpec {
	var resp ListAttachmentsRespSpec
	var cinderAttachment ListRespAttachment

	if 0 == len(attachments) {
		// Even if the number is 0, it must return {"attachments":[]}
		resp.Attachments = make([]ListRespAttachment, 0, 0)
	} else {
		for _, attachment := range attachments {
			cinderAttachment.Status = attachment.Status
			cinderAttachment.Instance = attachment.Metadata["instance_uuid"]
			cinderAttachment.VolumeID = attachment.VolumeId
			cinderAttachment.ID = attachment.Id

			resp.Attachments = append(resp.Attachments, cinderAttachment)
		}
	}

	return &resp
}

// *******************Create attachment*******************

// CreateAttachmentReqSpec ...
type CreateAttachmentReqSpec struct {
	Attachment CreateReqAttachment `json:"attachment"`
}

// CreateReqAttachment ...
type CreateReqAttachment struct {
	InstanceUuID string    `json:"instance_uuid"`
	Connector    Connector `json:"connector,omitempty"`
	VolumeUuID   string    `json:"volume_uuid"`
}

// CreateAttachmentRespSpec ...
type CreateAttachmentRespSpec struct {
	Attachment CreateRespAttachment `json:"attachment"`
}

// CreateRespAttachment ...
type CreateRespAttachment struct {
	Status         string `json:"status"`
	DetachedAt     string `json:"detached_at,omitempty"`
	ConnectionInfo `json:"connection_info"`
	AttachedAt     string `json:"attached_at,omitempty"`
	AttachMode     string `json:"attach_mode,omitempty"`
	Instance       string `json:"instance,omitempty"`
	VolumeID       string `json:"volume_id"`
	ID             string `json:"id"`
}

// CreateAttachmentReq ...
func CreateAttachmentReq(cinderReq *CreateAttachmentReqSpec) *model.VolumeAttachmentSpec {
	attachment := model.VolumeAttachmentSpec{}
	attachment.Metadata = make(map[string]string)
	attachment.Metadata["instance_uuid"] = cinderReq.Attachment.InstanceUuID
	attachment.HostInfo.Initiator = cinderReq.Attachment.Connector.Initiator
	attachment.HostInfo.Ip = cinderReq.Attachment.Connector.IP
	attachment.HostInfo.Platform = cinderReq.Attachment.Connector.Platform
	attachment.HostInfo.Host = cinderReq.Attachment.Connector.Host
	attachment.HostInfo.OsType = cinderReq.Attachment.Connector.OsType
	attachment.Mountpoint = cinderReq.Attachment.Connector.Mountpoint
	attachment.VolumeId = cinderReq.Attachment.VolumeUuID

	return &attachment
}

// CreateAttachmentResp ...
func CreateAttachmentResp(attachment *model.VolumeAttachmentSpec) *CreateAttachmentRespSpec {
	resp := CreateAttachmentRespSpec{}
	resp.Attachment.Status = attachment.Status
	resp.Attachment.ConnectionInfo.DriverVolumeType = attachment.ConnectionInfo.DriverVolumeType
	resp.Attachment.ConnectionInfo.ConnectionData = attachment.ConnectionInfo.ConnectionData
	//resp.Attachment.AttachedAt = attachment.Mountpoint
	resp.Attachment.Instance = attachment.Metadata["instance_uuid"]
	resp.Attachment.VolumeID = attachment.VolumeId
	resp.Attachment.ID = attachment.BaseModel.Id

	return &resp
}

// *******************Update an attachment*******************

// UpdateAttachmentReqSpec ...
type UpdateAttachmentReqSpec struct {
	Attachment UpdateReqAttachment `json:"attachment"`
}

// UpdateReqAttachment ...
type UpdateReqAttachment struct {
	Connector Connector `json:"connector"`
}

// UpdateAttachmentRespSpec ...
type UpdateAttachmentRespSpec struct {
	Attachment UpdateRespAttachment `json:"attachment"`
}

// UpdateRespAttachment ...
type UpdateRespAttachment struct {
	Status         string `json:"status"`
	DetachedAt     string `json:"detached_at,omitempty"`
	ConnectionInfo `json:"connection_info"`
	AttachedAt     string `json:"attached_at,omitempty"`
	AttachMode     string `json:"attach_mode,omitempty"`
	Instance       string `json:"instance,omitempty"`
	VolumeID       string `json:"volume_id"`
	ID             string `json:"id"`
}

// UpdateAttachmentReq ...
func UpdateAttachmentReq(cinderReq *UpdateAttachmentReqSpec) *model.VolumeAttachmentSpec {
	attachment := model.VolumeAttachmentSpec{}
	attachment.HostInfo.Initiator = cinderReq.Attachment.Connector.Initiator
	attachment.HostInfo.Ip = cinderReq.Attachment.Connector.IP
	attachment.HostInfo.Platform = cinderReq.Attachment.Connector.Platform
	attachment.HostInfo.Host = cinderReq.Attachment.Connector.Host
	attachment.HostInfo.OsType = cinderReq.Attachment.Connector.OsType
	attachment.Mountpoint = cinderReq.Attachment.Connector.Mountpoint

	return &attachment
}

// UpdateAttachmentResp ...
func UpdateAttachmentResp(attachment *model.VolumeAttachmentSpec) *UpdateAttachmentRespSpec {
	resp := UpdateAttachmentRespSpec{}
	resp.Attachment.Status = attachment.Status
	resp.Attachment.ConnectionInfo.DriverVolumeType = attachment.ConnectionInfo.DriverVolumeType
	resp.Attachment.ConnectionInfo.ConnectionData = attachment.ConnectionInfo.ConnectionData
	//resp.Attachment.AttachedAt = attachment.Mountpoint
	resp.Attachment.Instance = attachment.Metadata["instance_uuid"]
	resp.Attachment.VolumeID = attachment.VolumeId
	resp.Attachment.ID = attachment.BaseModel.Id

	return &resp
}
