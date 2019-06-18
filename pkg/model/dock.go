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

const (
	DockTypeProvioner = "provisioner"
	DockTypeAttacher  = "attacher"
)

// DockSpec is initialized by specific driver configuration. Each backend
// can be regarded as a docking service between SDS controller and storage
// service.
type DockSpec struct {
	*BaseModel

	// The type of the dock(provisioner/attacher)
	Type string `json:"type,omitempty"`

	// The name of the dock.
	Name string `json:"name,omitempty"`

	// The description of the dock.
	// +optional
	Description string `json:"description,omitempty"`

	// The status of the dock.
	// One of: "available" or "unavailable".
	Status string `json:"status,omitempty"`

	// The storage type of the dock.
	// One of: "block", "file" or "object".
	StorageType string `json:"storageType,omitempty"`

	// Endpoint represents the dock server's access address.
	Endpoint string `json:"endpoint,omitempty"`

	// NodeId represents the identification of the host, it can be considered
	// as instance id or hostname.
	NodeId string `json:"nodeId,omitempty"`

	// DriverName represents the dock provider.
	// Currently One of: "cinder", "ceph", "lvm", "default".
	DriverName string `json:"driverName,omitempty"`

	// Metadata should be kept until the scemantics between opensds volume
	// attachment and backend attached storage resouce description are clear.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`
}
