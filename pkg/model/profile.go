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
	"sort"
	"strings"
)

// An OpenSDS profile is identified by a unique name and ID. With adding
// extra properties, each profile can contains a set of tags of storage
// capabilities which are desirable features for a class of applications.
type ProfileSpec struct {
	*BaseModel
	// The uuid of project
	// + readOnly
	TenantId string `json:"TenantId"`

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

var profileSortKey string

type ProfileSlice []*ProfileSpec

func (profile ProfileSlice) Len() int { return len(profile) }

func (profile ProfileSlice) Swap(i, j int) { profile[i], profile[j] = profile[j], profile[i] }

func (profile ProfileSlice) Less(i, j int) bool {
	switch profileSortKey {

	case "ID":
		return profile[i].Id < profile[j].Id
	case "NAME":
		return profile[i].Name < profile[j].Name
	case "DESCRIPTION":
		return profile[i].Description < profile[j].Description
	}
	return false
}

func (c *ProfileSpec) FindValue(k string, p *ProfileSpec) string {
	switch k {
	case "Id":
		return p.Id
	case "CreatedAt":
		return p.CreatedAt
	case "UpdatedAt":
		return p.UpdatedAt
	case "Name":
		return p.Name
	case "Description":
		return p.Description
	case "StorageType":
		return p.StorageType
	}
	return ""
}

func (c *ProfileSpec) SortList(profiles []*ProfileSpec, sortKey, sortDir string) []*ProfileSpec {

	profileSortKey = sortKey

	if strings.EqualFold(sortDir, "asc") {
		sort.Sort(ProfileSlice(profiles))
	} else {
		sort.Sort(sort.Reverse(ProfileSlice(profiles)))
	}
	return profiles
}
