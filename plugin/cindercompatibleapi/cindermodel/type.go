// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package cindermodel

// *******************Create*******************

// CreateTypeReqSpec ...
type CreateTypeReqSpec struct {
	VolumeType VolumeTypeOfCreateTypeReq `json:"volume_type,omitempty"`
}

// VolumeTypeOfCreateTypeReq ...
type VolumeTypeOfCreateTypeReq struct {
	Name           string    `json:"name"`
	AccessIsPublic bool      `json:"os-volume-type-access:is_public,omitempty"`
	Description    string    `json:"description,omitempty"`
	Extras         ExtraSpec `json:"extra_specs,omitempty"`
}

// CreateTypeRespSpec ...
type CreateTypeRespSpec struct {
	VolumeType VolumeTypeOfCreateTypeResp `json:"volume_type,omitempty"`
}

// VolumeTypeOfCreateTypeResp ...
type VolumeTypeOfCreateTypeResp struct {
	IsPublic       bool      `json:"is_public,omitempty"`
	Extras         ExtraSpec `json:"extra_specs,omitempty"`
	Description    string    `json:"description,omitempty"`
	Name           string    `json:"name,omitempty"`
	ID             string    `json:"id,omitempty"`
	AccessIsPublic bool      `json:"os-volume-type-access:is_public,omitempty"`
}

// *******************Update*******************

// UpdateTypeReqSpec ...
type UpdateTypeReqSpec struct {
	VolumeType VolumeTypeOfUpdateTypeReq `json:"volume_type"`
}

// VolumeTypeOfUpdateTypeReq ...
type VolumeTypeOfUpdateTypeReq struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IsPublic    bool   `json:"is_public,omitempty"`
}

// UpdateTypeRespSpec ...
type UpdateTypeRespSpec struct {
	VolumeType VolumeTypeOfUpdateTypeResp `json:"volume_type"`
}

// VolumeTypeOfUpdateTypeResp ...
type VolumeTypeOfUpdateTypeResp struct {
	IsPublic    bool      `json:"is_public"`
	Extras      ExtraSpec `json:"extra_specs"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	ID          string    `json:"id"`
}

// *******************Create or update extra*******************

// AddExtraReqSpec ...
type AddExtraReqSpec struct {
	Extras ExtraSpec `json:"extra_specs"`
}

// AddExtraRespSpec ...
type AddExtraRespSpec struct {
	Extras ExtraSpec `json:"extra_specs"`
}

// *******************Show all extra*******************

// ShowAllExtraRespSpec ...
type ShowAllExtraRespSpec struct {
	Extras ExtraSpec `json:"extra_specs"`
}

// *******************Show all extra*******************

// ShowExtraRespSpec ...
type ShowExtraRespSpec map[string]interface{}

// *******************Update extra*******************

// UpdateExtraReqSpec ...
type UpdateExtraReqSpec map[string]interface{}

// UpdateExtraRespSpec ...
type UpdateExtraRespSpec map[string]interface{}

// *******************Show volume type detail*******************

// ShowTypeRespSpec ...
type ShowTypeRespSpec struct {
	VolumeType VolumeTypeOfShowTypeResp `json:"volume_type"`
}

// VolumeTypeOfShowTypeResp ...
type VolumeTypeOfShowTypeResp struct {
	IsPublic    bool      `json:"is_public"`
	Extras      ExtraSpec `json:"extra_specs"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	ID          string    `json:"id"`
}

// *******************List all volume types*******************

// ListTypeRespSpec ...
type ListTypeRespSpec struct {
	VolumeTypes []VolumeTypeOfListType `json:"volume_types"`
}

// VolumeTypeOfListType ...
type VolumeTypeOfListType struct {
	Extras         ExtraSpec `json:"extra_specs"`
	Name           string    `json:"name"`
	AccessIsPublic bool      `json:"os-volume-type-access:is_public"`
	IsPublic       bool      `json:"is_public"`
	ID             string    `json:"id"`
	Description    string    `json:"description"`
}
