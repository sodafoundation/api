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
This module implements the host data structure.

*/

package model

import "encoding/json"

// HostSpec is a comupter system which can be discoveried manually in agentless
// mode or automatically in Agent mode.
// It's a consumer of volume or file share from storage.
type HostSpec struct {
	*BaseModel

	// The uuid of the project that the host belongs to.
	TenantId string `json:"tenantId,omitempty"`

	// The uuid of the user that the host belongs to.
	// +optional
	UserId string `json:"userId,omitempty"`

	// The name of the host.
	// Only numbers, letters, '-', '_', '.' in ASCII characters are allowed.
	HostName string `json:"hostName,omitempty"`

	// The OS type of the host.
	OsType string `json:"osType,omitempty"`

	// The way to access host, system will access host to get more information
	// and install agent if accessMode is 'Agent'.
	// 'port', 'username'and 'password' are requried in 'Agent' mode.
	// Enum: [agent agentless]
	AccessMode string `json:"accessMode,omitempty"`

	// The locality that pool belongs to.
	IP string `json:"ip,omitempty"`

	// The accessible port for user connecting
	// +optional
	Port int64 `json:"port,omitempty"`

	// username
	// +optional
	Username string `json:"username,omitempty"`

	// password
	// +optional
	Password string `json:"password,omitempty"`

	// availability zones
	// +optional
	AvailabilityZones []string `json:"availabilityZones"`

	// initiators
	// +optional
	Initiators []*Initiator `json:"initiators"`
}

// Initiator can include any port which is used to connect storage
type Initiator struct {

	// port name
	PortName string `json:"portName,omitempty"`

	// protocol
	// Enum: [iSCSI FC]
	Protocol string `json:"protocol,omitempty"`
}

// MarshalJSON to remove sensitive data
func (m HostSpec) MarshalJSON() ([]byte, error) {
	type hostResp HostSpec
	resp := hostResp(m)
	resp.Password = ""
	return json.Marshal(resp)
}
