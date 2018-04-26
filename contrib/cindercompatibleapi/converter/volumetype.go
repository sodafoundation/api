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
This module implements a entry into the OpenSDS northbound service.
*/

package converter

import (
	"errors"

	"github.com/opensds/opensds/pkg/model"
)

// *******************Create a volume type*******************

// CreateTypeReqSpec ...
type CreateTypeReqSpec struct {
	VolumeType CreateReqVolumeType `json:"volume_type,omitempty"`
}

// CreateReqVolumeType ...
type CreateReqVolumeType struct {
	Name           string    `json:"name"`
	AccessIsPublic bool      `json:"os-volume-type-access:is_public,omitempty"`
	Description    string    `json:"description,omitempty"`
	Extras         ExtraSpec `json:"extra_specs,omitempty"`
}

// CreateTypeRespSpec ...
type CreateTypeRespSpec struct {
	VolumeType CreateRespVolumeType `json:"volume_type,omitempty"`
}

// CreateRespVolumeType ...
type CreateRespVolumeType struct {
	IsPublic       bool      `json:"is_public,omitempty"`
	Extras         ExtraSpec `json:"extra_specs,omitempty"`
	Description    string    `json:"description,omitempty"`
	Name           string    `json:"name,omitempty"`
	ID             string    `json:"id,omitempty"`
	AccessIsPublic bool      `json:"os-volume-type-access:is_public,omitempty"`
}

// CreateTypeReq ...
func CreateTypeReq(cinderReq *CreateTypeReqSpec) (*model.ProfileSpec, error) {
	profile := model.ProfileSpec{}

	profile.Name = cinderReq.VolumeType.Name
	if false == cinderReq.VolumeType.AccessIsPublic {
		return nil, errors.New("OpenSDS does not support os-volume-type-access:is_public = false")
	}
	profile.Description = cinderReq.VolumeType.Description
	profile.Extras = *(CinderExtraToOpenSDSExtra(&(cinderReq.VolumeType.Extras)))

	// The storageType can be block, file, object, default is block
	profile.StorageType = "block"

	return &profile, nil
}

// CinderExtraToOpenSDSExtra ...
func CinderExtraToOpenSDSExtra(typeExtra *ExtraSpec) *model.ExtraSpec {
	var profileExtras model.ExtraSpec
	profileExtras = make(map[string]interface{})
	for key, value := range *typeExtra {
		profileExtras[key] = value
	}

	return &profileExtras
}

// CreateTypeResp ...
func CreateTypeResp(profile *model.ProfileSpec) *CreateTypeRespSpec {
	resp := CreateTypeRespSpec{}
	resp.VolumeType.IsPublic = true
	resp.VolumeType.Extras = *(OpenSDSExtraToCinderExtra(&(profile.Extras)))
	resp.VolumeType.Description = profile.Description
	resp.VolumeType.Name = profile.Name
	resp.VolumeType.ID = profile.BaseModel.Id
	resp.VolumeType.AccessIsPublic = true

	return &resp
}

// *******************Update a volume type*******************

// UpdateTypeReqSpec ...
type UpdateTypeReqSpec struct {
	VolumeType UpdateReqVolumeType `json:"volume_type"`
}

// UpdateReqVolumeType ...
type UpdateReqVolumeType struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IsPublic    bool   `json:"is_public,omitempty"`
}

// UpdateTypeRespSpec ...
type UpdateTypeRespSpec struct {
	VolumeType UpdateRespVolumeType `json:"volume_type"`
}

// UpdateRespVolumeType ...
type UpdateRespVolumeType struct {
	IsPublic    bool      `json:"is_public"`
	Extras      ExtraSpec `json:"extra_specs"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	ID          string    `json:"id"`
}

// UpdateTypeReq ...
func UpdateTypeReq(cinderReq *UpdateTypeReqSpec) (*model.ProfileSpec, error) {
	profile := model.ProfileSpec{}
	profile.Name = cinderReq.VolumeType.Name
	profile.Description = cinderReq.VolumeType.Description
	if false == cinderReq.VolumeType.IsPublic {
		return nil, errors.New("OpenSDS does not support is_public = false")
	}

	return &profile, nil
}

// UpdateTypeResp ...
func UpdateTypeResp(profile *model.ProfileSpec) *UpdateTypeRespSpec {
	resp := UpdateTypeRespSpec{}
	resp.VolumeType.IsPublic = true
	resp.VolumeType.Extras = *(OpenSDSExtraToCinderExtra(&(profile.Extras)))
	resp.VolumeType.Description = profile.Description
	resp.VolumeType.Name = profile.Name
	resp.VolumeType.ID = profile.BaseModel.Id

	return &resp
}

// *******************Create or update extra specs for volume type*******************

// AddExtraReqSpec ...
type AddExtraReqSpec struct {
	Extras ExtraSpec `json:"extra_specs"`
}

// AddExtraRespSpec ...
type AddExtraRespSpec struct {
	Extras ExtraSpec `json:"extra_specs"`
}

// AddExtraReq ...
func AddExtraReq(cinderReq *AddExtraReqSpec) *model.ExtraSpec {
	var profileExtras model.ExtraSpec

	if len(cinderReq.Extras) >= 1 {
		profileExtras = make(map[string]interface{})
		for key, value := range cinderReq.Extras {
			profileExtras[key] = value
		}
	}

	return &profileExtras
}

// AddExtraResp ...
func AddExtraResp(profileExtras *model.ExtraSpec) *AddExtraRespSpec {
	var resp AddExtraRespSpec

	if len(*profileExtras) >= 1 {
		resp.Extras = make(map[string]interface{})
		for key, value := range *profileExtras {
			resp.Extras[key] = value
		}
	}

	return &resp
}

// *******************Show all extra specifications for volume type*******************

// ShowAllExtraRespSpec ...
type ShowAllExtraRespSpec struct {
	Extras ExtraSpec `json:"extra_specs"`
}

// ShowAllExtraResp ...
func ShowAllExtraResp(profileExtras *model.ExtraSpec) *ShowAllExtraRespSpec {
	var resp ShowAllExtraRespSpec

	if len(*profileExtras) >= 1 {
		resp.Extras = make(map[string]interface{})
		for key, value := range *profileExtras {
			resp.Extras[key] = value
		}
	}

	return &resp
}

// *******************Show extra specification for volume type*******************

// ShowExtraRespSpec ...
type ShowExtraRespSpec map[string]interface{}

//ShowExtraResp ...
func ShowExtraResp(reqkey string, profileExtras *model.ExtraSpec) *ShowExtraRespSpec {
	var resp ShowExtraRespSpec

	if (len(*profileExtras) >= 1) && (nil != (*profileExtras)[reqkey]) {
		resp = make(map[string]interface{})
		resp[reqkey] = (*profileExtras)[reqkey]
	}

	return &resp
}

// *******************Update extra specification for volume type*******************

// UpdateExtraReqSpec ...
type UpdateExtraReqSpec map[string]interface{}

// UpdateExtraRespSpec ...
type UpdateExtraRespSpec map[string]interface{}

// UpdateExtraReq ...
func UpdateExtraReq(reqkey string, cinderReq *UpdateExtraReqSpec) (*model.ExtraSpec, error) {
	var profileExtras model.ExtraSpec

	if (1 == len(*cinderReq)) && (nil != (*cinderReq)[reqkey]) {
		profileExtras = make(map[string]interface{})
		profileExtras[reqkey] = (*cinderReq)[reqkey]
	} else {
		return nil, errors.New("The body of the request is wrong")
	}

	return &profileExtras, nil
}

// UpdateExtraResp ...
func UpdateExtraResp(reqkey string, profileExtras *model.ExtraSpec) *UpdateExtraRespSpec {
	var resp UpdateExtraRespSpec

	if (len(*profileExtras) >= 1) && (nil != (*profileExtras)[reqkey]) {
		resp = make(map[string]interface{})
		resp[reqkey] = (*profileExtras)[reqkey]
	}

	return &resp
}

// *******************Show volume type detail*******************

// ShowTypeRespSpec ...
type ShowTypeRespSpec struct {
	VolumeType ShowRespVolumeType `json:"volume_type"`
}

// ShowRespVolumeType ...
type ShowRespVolumeType struct {
	IsPublic    bool      `json:"is_public"`
	Extras      ExtraSpec `json:"extra_specs"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	ID          string    `json:"id"`
}

// ShowTypeResp ...
func ShowTypeResp(profile *model.ProfileSpec) *ShowTypeRespSpec {
	resp := ShowTypeRespSpec{}
	resp.VolumeType.IsPublic = true
	resp.VolumeType.Extras = *(OpenSDSExtraToCinderExtra(&(profile.Extras)))
	resp.VolumeType.Description = profile.Description
	resp.VolumeType.Name = profile.Name
	resp.VolumeType.ID = profile.BaseModel.Id

	return &resp
}

// *******************List all volume types*******************

// ListTypesRespSpec ...
type ListTypesRespSpec struct {
	VolumeTypes []ListRespVolumeType `json:"volume_types"`
}

// ListRespVolumeType ...
type ListRespVolumeType struct {
	Extras         ExtraSpec `json:"extra_specs"`
	Name           string    `json:"name"`
	AccessIsPublic bool      `json:"os-volume-type-access:is_public"`
	IsPublic       bool      `json:"is_public"`
	ID             string    `json:"id"`
	Description    string    `json:"description"`
}

// ListTypesResp ...
func ListTypesResp(profiles []*model.ProfileSpec) *ListTypesRespSpec {
	var resp ListTypesRespSpec
	var volumeType ListRespVolumeType

	if 0 == len(profiles) {
		resp.VolumeTypes = make([]ListRespVolumeType, 0, 0)
	} else {
		for _, profile := range profiles {
			volumeType.Extras = *(OpenSDSExtraToCinderExtra(&(profile.Extras)))
			volumeType.Name = profile.Name
			volumeType.AccessIsPublic = true
			volumeType.IsPublic = true
			volumeType.ID = profile.BaseModel.Id
			volumeType.Description = profile.Description

			resp.VolumeTypes = append(resp.VolumeTypes, volumeType)
		}
	}

	return &resp
}

// OpenSDSExtraToCinderExtra ...
func OpenSDSExtraToCinderExtra(profileExtras *model.ExtraSpec) *ExtraSpec {
	var typeExtra ExtraSpec

	if len(*profileExtras) >= 1 {
		typeExtra = make(map[string]interface{})
		for key, value := range *profileExtras {
			typeExtra[key] = value
		}
	}

	return &typeExtra
}
