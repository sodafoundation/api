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

// DockSpec is initialized by specific driver configuration. Each backend
// can be regarded as a docking service between SDS controller and storage
// service.
type DockSpec struct {
	*BaseModel
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	// +readOnly:true
	Status      string `json:"status,omitempty"`
	StorageType string `json:"storageType,omitempty"`
	Endpoint    string `json:"endpoint,omitempty"`
	DriverName  string `json:"driverName,omitempty"`
}
