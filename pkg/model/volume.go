// Copyright 2017 The OpenSDS Authors.
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

import (
	"encoding/json"
)

// Volume is an block device created by storage service, it can be attached
// to physical machine or virtual machine instance.
type VolumeSpec struct {
	*BaseModel
	ProjectId   string `json:"projectId,omitempty"`
	UserId      string `json:"userId,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	// Default unit of volume Size is GB.
	Size             int64  `json:"size,omitempty"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
	// +readOnly:true
	Status    string `json:"status,omitempty"`
	PoolId    string `json:"poolId,omitempty"`
	ProfileId string `json:"profileId,omitempty"`
	// Metadata should be kept until the scemantics between opensds volume
	// and backend storage resouce description are clear.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`
}

type VolumeAttachmentSpec struct {
	*BaseModel
	ProjectId  string `json:"projectId,omitempty"`
	UserId     string `json:"userId,omitempty"`
	VolumeId   string `json:"volumeId,omitempty"`
	Mountpoint string `json:"mountpoint,omitempty"`
	// +readOnly:true
	Status string `json:"status,omitempty"`
	// Metadata should be kept until the scemantics between opensds volume
	// attachment and backend attached storage resouce description are clear.
	// +optional
	Metadata       map[string]string `json:"metadata,omitempty"`
	HostInfo       `json:"hostInfo,omitempty"`
	ConnectionInfo `json:"connectionInfo,omitempty"`
}

// HostInfo is a structure for all properties of host when create a volume
// attachment.
type HostInfo struct {
	Platform  string `json:"platform,omitempty"`
	OsType    string `json:"osType,omitempty"`
	Ip        string `json:"ip,omitempty"`
	Host      string `json:"host,omitempty"`
	Initiator string `json:"initiator,omitempty"`
}

// ConnectionInfo is a structure for all properties of connection when
// create a volume attachment.
type ConnectionInfo struct {
	DriverVolumeType     string                 `json:"driverVolumeType,omitempty"`
	ConnectionData       map[string]interface{} `json:"data,omitempty"`
	AdditionalProperties map[string]interface{} `json:"additionalProperties,omitempty"`
}

func (con *ConnectionInfo) EncodeConnectionData() []byte {
	conBody, _ := json.Marshal(&con.ConnectionData)
	return conBody
}

type VolumeSnapshotSpec struct {
	*BaseModel
	ProjectId   string `json:"projectId,omitempty"`
	UserId      string `json:"userId,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	// Default unit of snapshot Size is GB.
	Size     int64  `json:"size,omitempty"`
	Status   string `json:"status,omitempty"`
	VolumeId string `json:"volumeId,omitempty"`
	// Metadata should be kept until the scemantics between opensds volume
	// snapshot and backend storage resouce snapshot description are clear.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`
}
