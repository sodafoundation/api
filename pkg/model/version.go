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

type VersionSpec struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}

func (ver *VersionSpec) GetName() string {
	return ver.Name
}

func (ver *VersionSpec) GetDescription() string {
	return ver.Description
}

func (ver *VersionSpec) GetStatus() string {
	return ver.Status
}

func (ver *VersionSpec) GetUpdatedTime() string {
	return ver.UpdatedAt
}
