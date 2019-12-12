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

package model

import "encoding/json"

// VolumeAttachmentSpec is a description of volume attached resource.
type VolumeAttachmentSpec struct {
	*BaseModel

	// The uuid of the project that the volume belongs to.
	TenantId string `json:"tenantId,omitempty"`

	// The uuid of the user that the volume belongs to.
	// +optional
	UserId string `json:"userId,omitempty"`

	// The uuid of the host which the attachment belongs to.
	HostId string `json:"hostId,omitempty"`

	// The uuid of the volume which the attachment belongs to.
	VolumeId string `json:"volumeId,omitempty"`

	// The status of the attachment.
	// One of: "attaching", "attached", "error", etc.
	Status string `json:"status,omitempty"`

	// The locaility when the volume was attached to a host.
	Mountpoint string `json:"mountpoint,omitempty"`

	// read-only (‘ro’) or read-and-write (‘rw’), default is ‘rw’
	AttachMode string `json:"attachMode,omitempty"`

	// The protocol
	AccessProtocol string `json:"accessProtocol,omitempty"`

	// See details in `ConnectionInfo`
	ConnectionInfo `json:"connectionInfo,omitempty"`
}

// ConnectionInfo is a structure for all properties of connection when
// create a volume attachment.
type ConnectionInfo struct {
	DriverVolumeType     string                 `json:"driverVolumeType,omitempty"`
	ConnectionData       map[string]interface{} `json:"data,omitempty"`
	AdditionalProperties map[string]interface{} `json:"additionalProperties,omitempty"`
}

// EncodeConnectionData will marshal itself to byte
func (con *ConnectionInfo) EncodeConnectionData() []byte {
	conBody, _ := json.Marshal(&con.ConnectionData)
	return conBody
}
