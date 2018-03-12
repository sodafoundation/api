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

// A pool is discoveried and updated by a dock service. Each pool can be regarded
// as a physical storage pool or a virtual storage pool. It's a logical and
// atomic pool and can be abstracted from any storage platform.
type StoragePoolSpec struct {
	*BaseModel
	// The uuid of project
	// + readOnly
	ProjectId string `json:"projectId"`

	// The name of the pool.
	Name string `json:"name,omitempty"`

	// The description of the pool.
	// +optional
	Description string `json:"description,omitempty"`

	// The status of the pool.
	// One of: "available" or "unavailable".
	Status string `json:"status,omitempty"`

	// The uuid of the dock which the pool belongs to.
	DockId string `json:"dockId,omitempty"`

	// The locality that pool belongs to.
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// The total capacity of the pool.
	// Default unit of TotalCapacity is GB.
	TotalCapacity int64 `json:"totalCapacity,omitempty"`

	// The free capaicty of the pool.
	// Default unit of FreeCapacity is GB.
	FreeCapacity int64 `json:"freeCapacity,omitempty"`

	// The storage type of the dock.
	// One of: "block", "file" or "object".
	StorageType string `json:"storageType,omitempty"`

	// Map of keys and json object that represents the extra epecs
	// of the pool, such as supported capabilities.
	// +optional
	Extras ExtraSpec `json:"extras,omitempty"`
}
