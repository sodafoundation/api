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
This module implements the common data structure.

*/

package model

// FileShareSpec is an file share created by storage service, it can be mounted to physical machine or virtual machine instance.
type FileShareSpec struct {
	*BaseModel

	// The uuid of the project that the fileshare belongs to.
	TenantId string `json:"tenantId,omitempty"`

	// The uuid of the user that the file share belongs to.
	// +optional
	UserId string `json:"userId,omitempty"`

	// The name of the file share.
	Name string `json:"name,omitempty"`

	// The description of the file share.
	// +optional
	Description string `json:"description,omitempty"`

	// Creation time of fileshare.
	CreatedAt string `json:"createdAt,omitempty"`

	// Updation time of fileshare.
	UpdatedAt string `json:"updatedAt,omitempty"`

	// The protocol of the fileshare. e.g NFS, SMB etc.
	Protocols []string `json:"protocols,omitempty"`

	// The size of the file share requested by the user.
	// Default unit of file share Size is GB.
	Size int64 `json:"size,omitempty"`

	// The locality that file share belongs to.
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// The status of the file share.
	// One of: "available", "error", "in-use", etc.
	Status string `json:"status,omitempty"`

	// The uuid of the pool which the file share belongs to.
	// +readOnly
	PoolId string `json:"poolId,omitempty"`

	// The uuid of the profile which the file share belongs to.
	ProfileId string `json:"profileId,omitempty"`

	// Metadata should be kept until the scemantics between opensds file share
	// and backend storage resouce description are clear.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`

	// The uuid of the snapshot which the fileshare is created
	SnapshotId string `json:"snapshotId,omitempty"`

	// ExportLocations of the fileshare.
	ExportLocations []string `json:"exportLocations,omitempty"`
}

