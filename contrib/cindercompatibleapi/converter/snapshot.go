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
)

// *******************Create a snapshot*******************

// CreateSnapshotReqSpec ...
type CreateSnapshotReqSpec struct {
	Snapshot CreateReqSnapshot `json:"snapshot"`
}

// CreateReqSnapshot ...
type CreateReqSnapshot struct {
	VolumeID    string            `json:"volume_id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Force       bool              `json:"force,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// CreateSnapshotRespSpec ...
type CreateSnapshotRespSpec struct {
	Snapshot CreateRespSnapshot `json:"snapshot,omitempty"`
}

// CreateRespSnapshot ...
type CreateRespSnapshot struct {
	Status      string            `json:"status"`
	Description string            `json:"description"`
	CreatedAt   string            `json:"created_at"`
	Name        string            `json:"name"`
	UserID      string            `json:"user_id"`
	VolumeID    string            `json:"volume_id"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	ID          string            `json:"id"`
	Size        int64             `json:"size"`
	UpdatedAt   string            `json:"updated_at"`
}

// CreateSnapshotReq ...
func CreateSnapshotReq(cinderReq *CreateSnapshotReqSpec) (*model.VolumeSnapshotSpec, error) {
	req := model.VolumeSnapshotSpec{}
	req.VolumeId = cinderReq.Snapshot.VolumeID
	req.Name = cinderReq.Snapshot.Name
	req.Description = cinderReq.Snapshot.Description

	if false != cinderReq.Snapshot.Force {
		return nil, errors.New("OpenSDS does not support the parameter: force")
	}

	if 0 != len(cinderReq.Snapshot.Metadata) {
		return nil, errors.New("OpenSDS does not support the parameter: metadata")
	}

	return &req, nil
}

// CreateSnapshotResp ...
func CreateSnapshotResp(snapshot *model.VolumeSnapshotSpec) *CreateSnapshotRespSpec {
	resp := CreateSnapshotRespSpec{}
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

// *******************Update a snapshot*******************

// UpdateSnapshotReqSpec ...
type UpdateSnapshotReqSpec struct {
	Snapshot UpdateReqSnapshot `json:"snapshot"`
}

// UpdateReqSnapshot ...
type UpdateReqSnapshot struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
}

// UpdateSnapshotRespSpec ...
type UpdateSnapshotRespSpec struct {
	Snapshot UpdateRespSnapshot `json:"snapshot"`
}

// UpdateRespSnapshot ...
type UpdateRespSnapshot struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Size        int64  `json:"size"`
	VolumeID    string `json:"volume_id"`
	UserID      string `json:"user_id"`
}

// UpdateSnapshotReq ...
func UpdateSnapshotReq(cinderSnapshot *UpdateSnapshotReqSpec) *model.VolumeSnapshotSpec {
	req := model.VolumeSnapshotSpec{}
	req.Name = cinderSnapshot.Snapshot.Name
	req.Description = cinderSnapshot.Snapshot.Description

	return &req
}

// UpdateSnapshotResp ...
func UpdateSnapshotResp(snapshot *model.VolumeSnapshotSpec) *UpdateSnapshotRespSpec {
	resp := UpdateSnapshotRespSpec{}
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

// *******************Show a snapshot's details*******************

// ShowSnapshotDetailsRespSpec ...
type ShowSnapshotDetailsRespSpec struct {
	Snapshot ShowRespSnapshotDetails `json:"snapshot"`
}

// ShowRespSnapshotDetails ...
type ShowRespSnapshotDetails struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	Name        string `json:"name"`
	UserID      string `json:"user_id"`
	VolumeID    string `json:"volume_id"`
	Size        int64  `json:"size"`
	ID          string `json:"id"`
}

// ShowSnapshotDetailsResp ...
func ShowSnapshotDetailsResp(snapshot *model.VolumeSnapshotSpec) *ShowSnapshotDetailsRespSpec {
	resp := ShowSnapshotDetailsRespSpec{}
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

// *******************List accessible snapshots*******************

// ListSnapshotsRespSpec ...
type ListSnapshotsRespSpec struct {
	Snapshots []ListRespSnapshot `json:"snapshots"`
	Count     int64              `json:"count,omitempty"`
}

// ListRespSnapshot ...
type ListRespSnapshot struct {
	Status      string            `json:"status"`
	Description string            `json:"description"`
	CreatedAt   string            `json:"created_at"`
	Name        string            `json:"name"`
	UserID      string            `json:"user_id"`
	VolumeID    string            `json:"volume_id"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	ID          string            `json:"id"`
	Size        int64             `json:"size"`
}

// ListSnapshotsResp ...
func ListSnapshotsResp(snapshots []*model.VolumeSnapshotSpec) *ListSnapshotsRespSpec {
	var resp ListSnapshotsRespSpec
	var cinderSnapshot ListRespSnapshot

	if 0 == len(snapshots) {
		resp.Snapshots = make([]ListRespSnapshot, 0, 0)
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

// *******************List snapshots and details*******************

// ListSnapshotsDetailsRespSpec ...
type ListSnapshotsDetailsRespSpec struct {
	Snapshots []ListRespSnapshotDetails `json:"snapshots"`
	Count     int64                     `json:"count,omitempty"`
}

// ListRespSnapshotDetails ...
type ListRespSnapshotDetails struct {
	Status      string            `json:"status"`
	Progress    int64             `json:"os-extended-snapshot-attributes:progress,omitempty"`
	Description string            `json:"description"`
	CreatedAt   string            `json:"created_at"`
	Name        string            `json:"name"`
	UserID      string            `json:"user_id"`
	VolumeID    string            `json:"volume_id"`
	ProjectID   string            `json:"os-extended-snapshot-attributes:project_id,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	ID          string            `json:"id"`
	Size        int64             `json:"size"`
}

// ListSnapshotsDetailsResp ...
func ListSnapshotsDetailsResp(snapshots []*model.VolumeSnapshotSpec) *ListSnapshotsDetailsRespSpec {
	var resp ListSnapshotsDetailsRespSpec
	var cinderSnapshot ListRespSnapshotDetails

	if 0 == len(snapshots) {
		resp.Snapshots = make([]ListRespSnapshotDetails, 0, 0)
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
