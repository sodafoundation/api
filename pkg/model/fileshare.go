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
This module implements the common data structure.

*/
package model

type FileShareAclSpec struct {
	*BaseModel

	// The uuid of the project that the fileshare belongs to.
	TenantId string `json:"tenantId,omitempty"`

	// The uuid of the fileshare.
	FileShareId string `json:"fileshareId,omitempty"`

	// The type of access. Ex: IP based.
	Type string `json:"type,omitempty"`

	// The accessCapability for fileshare.
	AccessCapability []string `json:"accessCapability,omitempty"`

	// accessTo of the fileshare.
	AccessTo string `json:"accessTo,omitempty"`

	// The description of the fileshare acl.
	Description string `json:"description,omitempty"`

	// The uuid of the profile which the fileshare belongs to.
	ProfileId string `json:"profileId,omitempty"`

	// Metadata should be kept until the scemantics between opensds fileshare
	// and backend storage resouce description are clear.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`
}

// FileShareSpec is a schema for fileshare API. Fileshare will be created on some backend
// and can be shared among multiple users.

type FileShareSpec struct {
	*BaseModel
	// The uuid of the project that the fileshare belongs to.
	TenantId string `json:"tenantId,omitempty"`

	// The uuid of the user that the fileshare belongs to.
	// +optional
	UserId string `json:"userId,omitempty"`

	// The name of the fileshare.
	Name string `json:"name,omitempty"`

	// The protocol of the fileshare. e.g NFS, SMB etc.
	Protocols []string `json:"protocols,omitempty"`

	// The description of the fileshare.
	// +optional
	Description string `json:"description,omitempty"`

	// The size of the fileshare requested by the user.
	// Default unit of fileshare Size is GB.
	Size int64 `json:"size,omitempty"`

	// The locality that fileshare belongs to.
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// The status of the fileshare.
	// One of: "available", "error" etc.
	Status string `json:"status,omitempty"`

	// The uuid of the pool which the fileshare belongs to.
	// +readOnly
	PoolId string `json:"poolId,omitempty"`

	// The uuid of the profile which the fileshare belongs to.
	ProfileId string `json:"profileId,omitempty"`

	// The uuid of the snapshot which the fileshare is created
	SnapshotId string `json:"snapshotId,omitempty"`

	// ExportLocations of the fileshare.
	ExportLocations []string `json:"exportLocations,omitempty"`

	// Metadata should be kept until the scemantics between opensds fileshare
	// and backend storage resouce description are clear.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`
}

// FileShareSnapshotSpec is a description of fileshare snapshot resource.
type FileShareSnapshotSpec struct {
	*BaseModel

	// The uuid of the project that the fileshare belongs to.
	TenantId string `json:"tenantId,omitempty"`

	// The uuid of the user that the fileshare belongs to.
	// +optional
	UserId string `json:"userId,omitempty"`

	// The uuid of the fileshare.
	FileShareId string `json:"fileshareId,omitempty"`

	// The name of the fileshare snapshot.
	Name string `json:"name,omitempty"`

	// The description of the fileshare snapshot.
	// +optional
	Description string `json:"description,omitempty"`

	// The size of the fileshare which the snapshot belongs to.
	// Default unit of filesahre Size is GB.
	ShareSize int64 `json:"shareSize,omitempty"`

	// The size of the snapshot. Default unit of files snapshot Size is GB.
	SnapshotSize int64 `json:"snapshotSize,omitempty"`

	// The status of the fileshare snapshot.
	// One of: "available", "error", etc.
	Status string `json:"status,omitempty"`

	// The uuid of the profile which the fileshare belongs to.
	ProfileId string `json:"profileId,omitempty"`

	// Metadata should be kept until the scemantics between opensds fileshare
	// and backend storage resouce description are clear.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`
}
