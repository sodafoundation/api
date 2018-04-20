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
	"errors"

	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/plugin/cindercompatibleapi/cindermodel"
)

var (
	// APIVersion ...
	APIVersion = "v3"
	// EndPoint ...
	EndPoint = "http://127.0.0.1:8777/v3"
)

// *******************List Detail*******************

// ListVolumeDetailResp ...
func ListVolumeDetailResp(volumes []*model.VolumeSpec) *cindermodel.ListVolumeDetailRespSpec {
	var resp cindermodel.ListVolumeDetailRespSpec
	var volumeForList cindermodel.VolumeForListDetailResp

	if 0 == len(volumes) {
		resp.Volumes = make([]cindermodel.VolumeForListDetailResp, 0, 0)
	} else {
		for _, volume := range volumes {

			volumeForList.Attachments = make([]cindermodel.AttachmentOfVolumeResp, 0, 0)
			volumeForList.AvailabilityZone = volume.AvailabilityZone
			volumeForList.UpdatedAt = volume.BaseModel.UpdatedAt
			volumeForList.ID = volume.BaseModel.Id
			volumeForList.Size = volume.Size
			volumeForList.UserID = volume.UserId
			volumeForList.Metadata = make(map[string]string)
			//volumeForList.TenantID = volume.TenantId
			volumeForList.Status = volume.Status
			volumeForList.Description = volume.Description
			volumeForList.Name = volume.Name
			volumeForList.CreatedAt = volume.BaseModel.CreatedAt

			resp.Volumes = append(resp.Volumes, volumeForList)
		}
	}

	return &resp
}

// *******************Create*******************

// CreateVolumeReq ...
func CreateVolumeReq(cinderReq *cindermodel.CreateVolumeReqSpec) (*model.VolumeSpec, error) {
	volume := model.VolumeSpec{}
	volume.BaseModel = &model.BaseModel{}
	volume.Name = cinderReq.Volume.Name
	volume.Description = cinderReq.Volume.Description
	volume.Size = cinderReq.Volume.Size
	volume.AvailabilityZone = cinderReq.Volume.AvailabilityZone

	if ("" != cinderReq.Volume.SourceVolID) || (false != cinderReq.Volume.Multiattach) ||
		("" != cinderReq.Volume.SnapshotID) || ("" != cinderReq.Volume.BackupID) ||
		("" != cinderReq.Volume.ImageRef) || ("" != cinderReq.Volume.VolumeType) ||
		(0 != len(cinderReq.Volume.Metadata)) || ("" != cinderReq.Volume.ConsistencygroupID) {

		return nil, errors.New("When creating a volume, opensds does not support " +
			"id/source_volid/multiattach/snapshot_id/backup_id/imageRef/volume_type/metadata/consistencygroup_id in body")
	}

	return &volume, nil
}

// CreateVolumeResp ...
func CreateVolumeResp(volume *model.VolumeSpec) *cindermodel.CreateVolumeRespSpec {
	resp := cindermodel.CreateVolumeRespSpec{}

	resp.Volume.Attachments = make([]cindermodel.AttachmentOfVolumeResp, 0, 0)
	resp.Volume.AvailabilityZone = volume.AvailabilityZone
	resp.Volume.UpdatedAt = volume.BaseModel.UpdatedAt
	resp.Volume.ID = volume.BaseModel.Id
	resp.Volume.Size = volume.Size
	resp.Volume.UserID = volume.UserId
	resp.Volume.Metadata = make(map[string]string)
	resp.Volume.Status = volume.Status
	resp.Volume.Description = volume.Description
	resp.Volume.Name = volume.Name
	resp.Volume.CreatedAt = volume.BaseModel.CreatedAt

	return &resp
}

// *******************List*******************

// ListVolumeResp ...
func ListVolumeResp(volumes []*model.VolumeSpec) *cindermodel.ListVolumeRespSpec {
	var resp cindermodel.ListVolumeRespSpec
	var volumeForList cindermodel.VolumeForListResp

	if 0 == len(volumes) {
		resp.Volumes = make([]cindermodel.VolumeForListResp, 0, 0)
	} else {
		for _, volume := range volumes {
			volumeForList.ID = volume.Id
			volumeForList.Name = volume.Name

			resp.Volumes = append(resp.Volumes, volumeForList)
		}
	}

	return &resp
}

// *******************Show*******************

// ShowVolumeResp ...
func ShowVolumeResp(volume *model.VolumeSpec) *cindermodel.ShowVolumeRespSpec {
	resp := cindermodel.ShowVolumeRespSpec{}

	resp.Volume.Attachments = make([]cindermodel.AttachmentOfVolumeResp, 0, 0)
	resp.Volume.AvailabilityZone = volume.AvailabilityZone
	resp.Volume.UpdatedAt = volume.BaseModel.UpdatedAt
	resp.Volume.ID = volume.BaseModel.Id
	resp.Volume.Size = volume.Size
	resp.Volume.UserID = volume.UserId
	resp.Volume.Metadata = make(map[string]string)
	resp.Volume.Status = volume.Status
	resp.Volume.Description = volume.Description
	resp.Volume.Name = volume.Name
	resp.Volume.CreatedAt = volume.BaseModel.CreatedAt
	//resp.Volume.TenantID = volume.TenantId

	return &resp
}

// *******************Update*******************

// UpdateVolumeReq ...
func UpdateVolumeReq(cinderReq *cindermodel.UpdateVolumeReqSpec) (*model.VolumeSpec, error) {
	volume := model.VolumeSpec{}
	volume.BaseModel = &model.BaseModel{}
	volume.Description = cinderReq.Volume.Description
	volume.Name = cinderReq.Volume.Name

	if 0 != len(cinderReq.Volume.Metadata) {

		return nil, errors.New("When updating a volume, opensds does not support metadata")
	}

	return &volume, nil
}

// UpdateVolumeResp ...
func UpdateVolumeResp(volume *model.VolumeSpec) *cindermodel.UpdateVolumeRespSpec {
	resp := cindermodel.UpdateVolumeRespSpec{}
	resp.Volume.Attachments = make([]cindermodel.AttachmentOfVolumeResp, 0, 0)
	resp.Volume.AvailabilityZone = volume.AvailabilityZone
	resp.Volume.UpdatedAt = volume.BaseModel.UpdatedAt
	resp.Volume.ID = volume.BaseModel.Id
	resp.Volume.Size = volume.Size
	resp.Volume.UserID = volume.UserId
	resp.Volume.Metadata = make(map[string]string)
	resp.Volume.Status = volume.Status
	resp.Volume.Description = volume.Description
	resp.Volume.Name = volume.Name
	resp.Volume.CreatedAt = volume.BaseModel.CreatedAt

	return &resp
}
