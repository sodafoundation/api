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

// *******************Create*******************

// CreateSnapshotReqSpec ...
type CreateSnapshotReqSpec struct {
	Snapshot CreateSnapshotReq `json:"snapshot"`
}

// CreateSnapshotReq ...
type CreateSnapshotReq struct {
	VolumeID    string            `json:"volume_id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Force       bool              `json:"force,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// CreateSnapshotRespSpec ...
type CreateSnapshotRespSpec struct {
	Snapshot CreateSnapshotResp `json:"snapshot,omitempty"`
}

// CreateSnapshotResp ...
type CreateSnapshotResp struct {
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

// *******************Update*******************

// UpdateSnapshotReqSpec ...
type UpdateSnapshotReqSpec struct {
	Snapshot UpdateSnapshotReq `json:"snapshot"`
}

// UpdateSnapshotReq ...
type UpdateSnapshotReq struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
}

// UpdateSnapshotRespSpec ...
type UpdateSnapshotRespSpec struct {
	Snapshot UpdateSnapshotResp `json:"snapshot"`
}

// UpdateSnapshotResp ...
type UpdateSnapshotResp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Size        int64  `json:"size"`
	VolumeID    string `json:"volume_id"`
	UserID      string `json:"user_id"`
}

// *******************Show details*******************

// ShowSnapshotDetailsRespSpec ...
type ShowSnapshotDetailsRespSpec struct {
	Snapshot ShowSnapshotDetailsResp `json:"snapshot"`
}

// ShowSnapshotDetailsResp ...
type ShowSnapshotDetailsResp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	Name        string `json:"name"`
	UserID      string `json:"user_id"`
	VolumeID    string `json:"volume_id"`
	Size        int64  `json:"size"`
	ID          string `json:"id"`
}

// *******************List*******************

// ListSnapshotRespSpec ...
type ListSnapshotRespSpec struct {
	Snapshots []ListSnapshotResp `json:"snapshots"`
	Count     int64              `json:"count,omitempty"`
}

// ListSnapshotResp ...
type ListSnapshotResp struct {
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

// *******************List detail*******************

// ListSnapshotDetailRespSpec ...
type ListSnapshotDetailRespSpec struct {
	Snapshots []ListSnapshotDetailResp `json:"snapshots"`
	Count     int64                    `json:"count,omitempty"`
}

// ListSnapshotDetailResp ...
type ListSnapshotDetailResp struct {
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
