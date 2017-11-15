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

type ProfileSpec struct {
	*BaseModel
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	StorageType string    `json:"storageType,omitempty"`
	Extra       ExtraSpec `json:"extras,omitempty"`
}

func (prf *ProfileSpec) GetName() string {
	return prf.Name
}

func (prf *ProfileSpec) GetDescription() string {
	return prf.Description
}

func (prf *ProfileSpec) GetStorageType() string {
	return prf.StorageType
}

type ExtraSpec map[string]interface{}
