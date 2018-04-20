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

// *******************Create*******************

// CreateSnapshotReq ...
func CreateSnapshotReq(cinderReq *cindermodel.CreateSnapshotReqSpec) (*model.VolumeSnapshotSpec, error) {
	req := model.VolumeSnapshotSpec{}
	req.VolumeId = cinderReq.Snapshot.VolumeID
	req.Name = cinderReq.Snapshot.Name
	req.Description = cinderReq.Snapshot.Description

	if false != cinderReq.Snapshot.Force {
		return nil, errors.New("Opensds does not support the parameter: force")
	}

	if 0 != len(cinderReq.Snapshot.Metadata) {
		return nil, errors.New("Opensds does not support the parameter: metadata")
	}

	return &req, nil
}

// CreateSnapshotResp ...
func CreateSnapshotResp(snapshot *model.VolumeSnapshotSpec) *cindermodel.CreateSnapshotRespSpec {
	resp := cindermodel.CreateSnapshotRespSpec{}
	resp.Snapshot.Status = snapshot.Status
	resp.Snapshot.Description = snapshot.Description
	resp.Snapshot.CreatedAt = snapshot.BaseModel.CreatedAt
	resp.Snapshot.Name = snapshot.Name
	resp.Snapshot.UserID = snapshot.UserId
	resp.Snapshot.VolumeID = snapshot.VolumeId
	//resp.Snapshot.Metadata = snapshot.Metadata
	resp.Snapshot.ID = snapshot.BaseModel.Id
	resp.Snapshot.Size = snapshot.Size
	resp.Snapshot.UpdatedAt = snapshot.BaseModel.UpdatedAt

	return &resp
}

// *******************Update*******************

// UpdateSnapshotReq ...
func UpdateSnapshotReq(cinderSnapshot *cindermodel.UpdateSnapshotReqSpec) *model.VolumeSnapshotSpec {
	req := model.VolumeSnapshotSpec{}
	req.Name = cinderSnapshot.Snapshot.Name
	req.Description = cinderSnapshot.Snapshot.Description

	return &req
}

// UpdateSnapshotResp ...
func UpdateSnapshotResp(snapshot *model.VolumeSnapshotSpec) *cindermodel.UpdateSnapshotRespSpec {
	resp := cindermodel.UpdateSnapshotRespSpec{}
	resp.Snapshot.Status = snapshot.Status
	resp.Snapshot.Description = snapshot.Description
	resp.Snapshot.CreatedAt = snapshot.BaseModel.CreatedAt
	resp.Snapshot.Name = snapshot.Name
	resp.Snapshot.ID = snapshot.BaseModel.Id
	resp.Snapshot.Size = snapshot.Size
	resp.Snapshot.VolumeID = snapshot.VolumeId
	resp.Snapshot.UserID = snapshot.UserId

	return &resp
}

// *******************Show details*******************

// ShowSnapshotDetailsResp ...
func ShowSnapshotDetailsResp(snapshot *model.VolumeSnapshotSpec) *cindermodel.ShowSnapshotDetailsRespSpec {
	resp := cindermodel.ShowSnapshotDetailsRespSpec{}
	resp.Snapshot.Status = snapshot.Status
	resp.Snapshot.Description = snapshot.Description
	resp.Snapshot.CreatedAt = snapshot.BaseModel.CreatedAt
	resp.Snapshot.Name = snapshot.Name
	resp.Snapshot.UserID = snapshot.UserId
	resp.Snapshot.VolumeID = snapshot.VolumeId
	resp.Snapshot.Size = snapshot.Size
	resp.Snapshot.ID = snapshot.BaseModel.Id

	return &resp
}

// *******************List*******************

// ListSnapshotResp ...
func ListSnapshotResp(snapshots []*model.VolumeSnapshotSpec) *cindermodel.ListSnapshotRespSpec {
	var resp cindermodel.ListSnapshotRespSpec
	var cinderSnapshot cindermodel.ListSnapshotResp

	if 0 == len(snapshots) {
		resp.Snapshots = make([]cindermodel.ListSnapshotResp, 0, 0)
	} else {
		for _, snapshot := range snapshots {
			cinderSnapshot.Status = snapshot.Status
			cinderSnapshot.Description = snapshot.Description
			cinderSnapshot.CreatedAt = snapshot.BaseModel.CreatedAt
			cinderSnapshot.Name = snapshot.Name
			cinderSnapshot.UserID = snapshot.UserId
			cinderSnapshot.VolumeID = snapshot.VolumeId
			//cinderSnapshot.Metadata = snapshot.Metadata
			cinderSnapshot.ID = snapshot.BaseModel.Id
			cinderSnapshot.Size = snapshot.Size

			resp.Snapshots = append(resp.Snapshots, cinderSnapshot)
		}
	}

	return &resp
}

// ListSnapshotDetailResp ...
func ListSnapshotDetailResp(snapshots []*model.VolumeSnapshotSpec) *cindermodel.ListSnapshotDetailRespSpec {
	var resp cindermodel.ListSnapshotDetailRespSpec
	var cinderSnapshot cindermodel.ListSnapshotDetailResp

	if 0 == len(snapshots) {
		resp.Snapshots = make([]cindermodel.ListSnapshotDetailResp, 0, 0)
	} else {
		for _, snapshot := range snapshots {
			cinderSnapshot.Status = snapshot.Status
			cinderSnapshot.Description = snapshot.Description
			cinderSnapshot.CreatedAt = snapshot.BaseModel.CreatedAt
			cinderSnapshot.Name = snapshot.Name
			cinderSnapshot.UserID = snapshot.UserId
			cinderSnapshot.VolumeID = snapshot.VolumeId
			//cinderSnapshot.Metadata = snapshot.Metadata
			cinderSnapshot.ID = snapshot.BaseModel.Id
			cinderSnapshot.Size = snapshot.Size

			resp.Snapshots = append(resp.Snapshots, cinderSnapshot)
		}
	}

	return &resp
}
