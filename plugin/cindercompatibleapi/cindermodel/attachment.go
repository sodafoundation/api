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
This module implements the common data structure.

*/

package cindermodel

// *******************show*******************

// ShowAttachmentRespSpec ...
type ShowAttachmentRespSpec struct {
	VolumeAttachment AttachmentOfShowResp `json:"attachment,omitempty"`
}

// AttachmentOfShowResp ...
type AttachmentOfShowResp struct {
	Status         string `json:"status"`
	DetachedAt     string `json:"detached_at,omitempty"`
	ConnectionInfo `json:"connection_info"`
	AttachedAt     string `json:"attached_at,omitempty"`
	AttachMode     string `json:"attach_mode,omitempty"`
	Instance       string `json:"instance,omitempty"`
	VolumeID       string `json:"volume_id"`
	ID             string `json:"id"`
}

// *******************List Detail*******************

// ListAttachmentDetailRespSpec ...
type ListAttachmentDetailRespSpec struct {
	Attachments []AttachmentOfListDetailResp `json:"attachments"`
}

// AttachmentOfListDetailResp ...
type AttachmentOfListDetailResp struct {
	Status         string         `json:"status"`
	DetachedAt     string         `json:"detached_at,omitempty"`
	ConnectionInfo ConnectionInfo `json:"connection_info"`
	AttachedAt     string         `json:"attached_at,omitempty"`
	AttachMode     string         `json:"attach_mode,omitempty"`
	Instance       string         `json:"instance,omitempty"`
	VolumeID       string         `json:"volume_id"`
	ID             string         `json:"id"`
}

// *******************List*******************

// ListAttachmentRespSpec ...
type ListAttachmentRespSpec struct {
	Attachments []AttachmentOfListResp `json:"attachments"`
}

// AttachmentOfListResp ...
type AttachmentOfListResp struct {
	Status   string `json:"status"`
	Instance string `json:"instance,omitempty"`
	VolumeID string `json:"volume_id"`
	ID       string `json:"id"`
}

// *******************Create*******************

// CreateAttachmentReqSpec ...
type CreateAttachmentReqSpec struct {
	Attachment AttachmentOfCreateReq `json:"attachment"`
}

// AttachmentOfCreateReq ...
type AttachmentOfCreateReq struct {
	InstanceUuID string    `json:"instance_uuid"`
	Connector    Connector `json:"connector,omitempty"`
	VolumeUuID   string    `json:"volume_uuid"`
}

// CreateAttachmentRespSpec ...
type CreateAttachmentRespSpec struct {
	Attachment AttachmentOfCreateResp `json:"attachment"`
}

// AttachmentOfCreateResp ...
type AttachmentOfCreateResp struct {
	Status         string `json:"status"`
	DetachedAt     string `json:"detached_at,omitempty"`
	ConnectionInfo `json:"connection_info"`
	AttachedAt     string `json:"attached_at,omitempty"`
	AttachMode     string `json:"attach_mode,omitempty"`
	Instance       string `json:"instance,omitempty"`
	VolumeID       string `json:"volume_id"`
	ID             string `json:"id"`
}

// *******************Update*******************

// UpdateAttachmentReqSpec ...
type UpdateAttachmentReqSpec struct {
	Attachment AttachmentOfUpdateReq `json:"attachment"`
}

// AttachmentOfUpdateReq ...
type AttachmentOfUpdateReq struct {
	Connector Connector `json:"connector"`
}

// UpdateAttachmentRespSpec ...
type UpdateAttachmentRespSpec struct {
	Attachment AttachmentOfUpdateResp `json:"attachment"`
}

// AttachmentOfUpdateResp ...
type AttachmentOfUpdateResp struct {
	Status         string `json:"status"`
	DetachedAt     string `json:"detached_at,omitempty"`
	ConnectionInfo `json:"connection_info"`
	AttachedAt     string `json:"attached_at,omitempty"`
	AttachMode     string `json:"attach_mode,omitempty"`
	Instance       string `json:"instance,omitempty"`
	VolumeID       string `json:"volume_id"`
	ID             string `json:"id"`
}

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
