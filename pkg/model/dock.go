// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

/*
This module implements the common data structure.

*/

package model

import (
	"encoding/json"
)

type DockSpec struct {
	*BaseModel
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Status      string                 `json:"status,omitempty"`
	StorageType string                 `json:"storageType,omitempty"`
	Endpoint    string                 `json:"endpoint,omitempty"`
	DriverName  string                 `json:"driverName,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

func (dck *DockSpec) GetName() string {
	return dck.Name
}

func (dck *DockSpec) GetDescription() string {
	return dck.Description
}

func (dck *DockSpec) GetStatus() string {
	return dck.Status
}

func (dck *DockSpec) GetStorageType() string {
	return dck.StorageType
}

func (dck *DockSpec) GetEndpoint() string {
	return dck.Endpoint
}

func (dck *DockSpec) GetDriverName() string {
	return dck.DriverName
}

func (dck *DockSpec) GetParameters() map[string]interface{} {
	return dck.Parameters
}

func (dck *DockSpec) EncodeParameters() []byte {
	parmBody, _ := json.Marshal(&dck.Parameters)
	return parmBody
}
