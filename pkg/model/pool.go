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

type StoragePoolSpec struct {
	*BaseModel
	Name             string                 `json:"name,omitempty"`
	Description      string                 `json:"description,omitempty"`
	Status           string                 `json:"status,omitempty"`
	DockId           string                 `json:"dockId,omitempty"`
	AvailabilityZone string                 `json:"availabilityZone,omitempty"`
	TotalCapacity    int64                  `json:"totalCapacity,omitempty"`
	FreeCapacity     int64                  `json:"freeCapacity,omitempty"`
	StorageType      string                 `json:"-"`
	Parameters       map[string]interface{} `json:"extras,omitempty"`
}

func (pol *StoragePoolSpec) GetName() string {
	return pol.Name
}

func (pol *StoragePoolSpec) GetDescription() string {
	return pol.Description
}

func (pol *StoragePoolSpec) GetStatus() string {
	return pol.Status
}

func (pol *StoragePoolSpec) GetDockId() string {
	return pol.DockId
}

func (pol *StoragePoolSpec) GetAvailability() string {
	return pol.AvailabilityZone
}

func (pol *StoragePoolSpec) GetTotalCapacity() int64 {
	return pol.TotalCapacity
}

func (pol *StoragePoolSpec) GetFreeCapacity() int64 {
	return pol.FreeCapacity
}

func (pol *StoragePoolSpec) GetStorageType() string {
	return pol.StorageType
}

func (pol *StoragePoolSpec) GetParameters() map[string]interface{} {
	return pol.Parameters
}

func (pol *StoragePoolSpec) EncodeParameters() []byte {
	parmBody, _ := json.Marshal(&pol.Parameters)
	return parmBody
}
