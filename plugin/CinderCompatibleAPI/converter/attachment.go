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
	"github.com/opensds/opensds/plugin/CinderCompatibleAPI/CinderModel"
)

// *******************show*******************

// ShowAttachmentResp ...
func ShowAttachmentResp(attachment *model.VolumeAttachmentSpec) *CinderModel.ShowAttachmentRespSpec {
	resp := CinderModel.ShowAttachmentRespSpec{}
	resp.VolumeAttachment.Status = attachment.Status
	resp.VolumeAttachment.ConnectionInfo.DriverVolumeType = attachment.ConnectionInfo.DriverVolumeType
	resp.VolumeAttachment.ConnectionInfo.ConnectionData = attachment.ConnectionInfo.ConnectionData
	//resp.VolumeAttachment.AttachedAt = attachment.Mountpoint
	resp.VolumeAttachment.Instance = attachment.Metadata["instance_uuid"]
	resp.VolumeAttachment.VolumeID = attachment.VolumeId
	resp.VolumeAttachment.ID = attachment.BaseModel.Id

	return &resp
}

// *******************List Detail*******************

// ListAttachmentDetailResp ...
func ListAttachmentDetailResp(attachments []*model.VolumeAttachmentSpec) *CinderModel.ListAttachmentDetailRespSpec {
	var resp CinderModel.ListAttachmentDetailRespSpec
	attachmentOfListDetail := CinderModel.AttachmentOfListDetailResp{}

	if 0 == len(attachments) {
		// Even if the number is 0, it must return {"attachments":[]}
		resp.Attachments = make([]CinderModel.AttachmentOfListDetailResp, 0, 0)
	} else {
		for _, attachment := range attachments {
			attachmentOfListDetail.Status = attachment.Status
			attachmentOfListDetail.ConnectionInfo.DriverVolumeType = attachment.ConnectionInfo.DriverVolumeType
			attachmentOfListDetail.ConnectionInfo.ConnectionData = attachment.ConnectionInfo.ConnectionData
			//attachmentOfListDetail.AttachedAt = attachment.Mountpoint
			attachmentOfListDetail.Instance = attachment.Metadata["instance_uuid"]
			attachmentOfListDetail.VolumeID = attachment.VolumeId
			attachmentOfListDetail.ID = attachment.Id

			resp.Attachments = append(resp.Attachments, attachmentOfListDetail)
		}
	}

	return &resp
}

// *******************List*******************

// ListAttachmentResp ...
func ListAttachmentResp(attachments []*model.VolumeAttachmentSpec) *CinderModel.ListAttachmentRespSpec {
	var resp CinderModel.ListAttachmentRespSpec
	var attachmentOfList CinderModel.AttachmentOfListResp

	if 0 == len(attachments) {
		// Even if the number is 0, it must return {"attachments":[]}
		resp.Attachments = make([]CinderModel.AttachmentOfListResp, 0, 0)
	} else {
		for _, attachment := range attachments {
			attachmentOfList.Status = attachment.Status
			attachmentOfList.Instance = attachment.Metadata["instance_uuid"]
			attachmentOfList.VolumeID = attachment.VolumeId
			attachmentOfList.ID = attachment.Id

			resp.Attachments = append(resp.Attachments, attachmentOfList)
		}
	}

	return &resp
}

// *******************Create*******************

// CreateAttachmentReq ...
func CreateAttachmentReq(cinderReq *CinderModel.CreateAttachmentReqSpec) *model.VolumeAttachmentSpec {
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
func CreateAttachmentResp(attachment *model.VolumeAttachmentSpec) *CinderModel.CreateAttachmentRespSpec {
	resp := CinderModel.CreateAttachmentRespSpec{}
	resp.Attachment.Status = attachment.Status
	resp.Attachment.ConnectionInfo.DriverVolumeType = attachment.ConnectionInfo.DriverVolumeType
	resp.Attachment.ConnectionInfo.ConnectionData = attachment.ConnectionInfo.ConnectionData
	//resp.Attachment.AttachedAt = attachment.Mountpoint
	resp.Attachment.Instance = attachment.Metadata["instance_uuid"]
	resp.Attachment.VolumeID = attachment.VolumeId
	resp.Attachment.ID = attachment.BaseModel.Id

	return &resp
}

// *******************Update*******************

// UpdateAttachmentReq ...
func UpdateAttachmentReq(cinderReq *CinderModel.UpdateAttachmentReqSpec) *model.VolumeAttachmentSpec {
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
func UpdateAttachmentResp(attachment *model.VolumeAttachmentSpec) *CinderModel.UpdateAttachmentRespSpec {
	resp := CinderModel.UpdateAttachmentRespSpec{}
	resp.Attachment.Status = attachment.Status
	resp.Attachment.ConnectionInfo.DriverVolumeType = attachment.ConnectionInfo.DriverVolumeType
	resp.Attachment.ConnectionInfo.ConnectionData = attachment.ConnectionInfo.ConnectionData
	//resp.Attachment.AttachedAt = attachment.Mountpoint
	resp.Attachment.Instance = attachment.Metadata["instance_uuid"]
	resp.Attachment.VolumeID = attachment.VolumeId
	resp.Attachment.ID = attachment.BaseModel.Id

	return &resp
}
