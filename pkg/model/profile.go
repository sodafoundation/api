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

import (
	"encoding/json"
)

// An OpenSDS profile is identified by a unique name and ID. With adding
// extra properties, each profile can contains a set of tags of storage
// capabilities which are desirable features for a class of applications.
type ProfileSpec struct {
	*BaseModel

	// The name of the profile.
	Name string `json:"name,omitempty"`

	// The description of the profile.
	// +optional
	Description string `json:"description,omitempty"`

	// The storage type of the profile.
	// One of: "block", "file" or "object".
	StorageType string `json:"storageType,omitempty"`

	// Map of keys and json object that represents the extra specs
	// of the profile, such as requested capabilities.
	// +optional
	Extras ExtraSpec `json:"extras,omitempty"`
}

// ExtraSpec is a dictionary object that contains unique keys and json
// objects.
type ExtraSpec map[string]interface{}

func (ext ExtraSpec) Encode() []byte {
	parmBody, _ := json.Marshal(&ext)
	return parmBody
}
