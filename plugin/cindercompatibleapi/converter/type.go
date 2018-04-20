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
	"github.com/opensds/opensds/plugin/CinderCompatibleAPI/cindermodel"
)

// *******************Create*******************

// CreateTypeReq ...
func CreateTypeReq(cinderReq *cindermodel.CreateTypeReqSpec) (*model.ProfileSpec, error) {
	profile := model.ProfileSpec{}

	profile.Name = cinderReq.VolumeType.Name
	if false == cinderReq.VolumeType.AccessIsPublic {
		return nil, errors.New("When creating a volume type, opensds does not support os-volume-type-access:is_public = false")
	}
	profile.Description = cinderReq.VolumeType.Description
	profile.Extras = *(CinderExtraToOpenSDSExtra(&(cinderReq.VolumeType.Extras)))

	// The storageType can be block, file, object, default is block
	profile.StorageType = "block"

	return &profile, nil
}

// CinderExtraToOpenSDSExtra ...
func CinderExtraToOpenSDSExtra(typeExtra *cindermodel.ExtraSpec) *model.ExtraSpec {
	var profileExtras model.ExtraSpec
	profileExtras = make(map[string]interface{})
	for key, value := range *typeExtra {
		profileExtras[key] = value
	}

	return &profileExtras
}

// CreateTypeResp ...
func CreateTypeResp(profile *model.ProfileSpec) *cindermodel.CreateTypeRespSpec {
	resp := cindermodel.CreateTypeRespSpec{}
	resp.VolumeType.IsPublic = true
	resp.VolumeType.Extras = *(OpenSDSExtraToCinderExtra(&(profile.Extras)))
	resp.VolumeType.Description = profile.Description
	resp.VolumeType.Name = profile.Name
	resp.VolumeType.ID = profile.BaseModel.Id
	resp.VolumeType.AccessIsPublic = true

	return &resp
}

// *******************Update*******************

// UpdateTypeReq ...
func UpdateTypeReq(cinderReq *cindermodel.UpdateTypeReqSpec) (*model.ProfileSpec, error) {
	profile := model.ProfileSpec{}
	profile.Name = cinderReq.VolumeType.Name
	profile.Description = cinderReq.VolumeType.Description
	if false == cinderReq.VolumeType.IsPublic {
		return nil, errors.New("When updating a volume type, opensds does not support is_public = false")
	}

	return &profile, nil
}

// UpdateTypeResp ...
func UpdateTypeResp(profile *model.ProfileSpec) *cindermodel.UpdateTypeRespSpec {
	resp := cindermodel.UpdateTypeRespSpec{}
	resp.VolumeType.IsPublic = true
	resp.VolumeType.Extras = *(OpenSDSExtraToCinderExtra(&(profile.Extras)))
	resp.VolumeType.Description = profile.Description
	resp.VolumeType.Name = profile.Name
	resp.VolumeType.ID = profile.BaseModel.Id

	return &resp
}

// *******************Create or update extra*******************

// AddExtraReq ...
func AddExtraReq(cinderReq *cindermodel.AddExtraReqSpec) *model.ExtraSpec {
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
func AddExtraResp(profileExtras *model.ExtraSpec) *cindermodel.AddExtraRespSpec {
	var resp cindermodel.AddExtraRespSpec

	if len(*profileExtras) >= 1 {
		resp.Extras = make(map[string]interface{})
		for key, value := range *profileExtras {
			resp.Extras[key] = value
		}
	}

	return &resp
}

// *******************Show all extra*******************

// ShowAllExtraResp ...
func ShowAllExtraResp(profileExtras *model.ExtraSpec) *cindermodel.ShowAllExtraRespSpec {
	var resp cindermodel.ShowAllExtraRespSpec

	if len(*profileExtras) >= 1 {
		resp.Extras = make(map[string]interface{})
		for key, value := range *profileExtras {
			resp.Extras[key] = value
		}
	}

	return &resp
}

// *******************Show extra*******************

//ShowExtraResp ...
func ShowExtraResp(reqkey string, profileExtras *model.ExtraSpec) *cindermodel.ShowExtraRespSpec {
	var resp cindermodel.ShowExtraRespSpec

	if (len(*profileExtras) >= 1) && (nil != (*profileExtras)[reqkey]) {
		resp = make(map[string]interface{})
		resp[reqkey] = (*profileExtras)[reqkey]
	}

	return &resp
}

// *******************Update extra*******************

// UpdateExtraReq ...
func UpdateExtraReq(reqkey string, cinderReq *cindermodel.UpdateExtraReqSpec) (*model.ExtraSpec, error) {
	var profileExtras model.ExtraSpec

	if (1 == len(*cinderReq)) && (nil != (*cinderReq)[reqkey]) {
		profileExtras = make(map[string]interface{})
		profileExtras[reqkey] = (*cinderReq)[reqkey]
	} else {
		return nil, errors.New("The bady of the request is wrong")
	}

	return &profileExtras, nil
}

// UpdateExtraResp ...
func UpdateExtraResp(reqkey string, profileExtras *model.ExtraSpec) *cindermodel.UpdateExtraRespSpec {
	var resp cindermodel.UpdateExtraRespSpec

	if (len(*profileExtras) >= 1) && (nil != (*profileExtras)[reqkey]) {
		resp = make(map[string]interface{})
		resp[reqkey] = (*profileExtras)[reqkey]
	}

	return &resp
}

// *******************Show Type*******************

// ShowTypeResp ...
func ShowTypeResp(profile *model.ProfileSpec) *cindermodel.ShowTypeRespSpec {
	resp := cindermodel.ShowTypeRespSpec{}
	resp.VolumeType.IsPublic = true
	resp.VolumeType.Extras = *(OpenSDSExtraToCinderExtra(&(profile.Extras)))
	resp.VolumeType.Description = profile.Description
	resp.VolumeType.Name = profile.Name
	resp.VolumeType.ID = profile.BaseModel.Id

	return &resp
}

// *******************List Type*******************

// ListTypeResp ...
func ListTypeResp(profiles []*model.ProfileSpec) *cindermodel.ListTypeRespSpec {
	var resp cindermodel.ListTypeRespSpec
	var volumeType cindermodel.VolumeTypeOfListType

	if 0 == len(profiles) {
		resp.VolumeTypes = make([]cindermodel.VolumeTypeOfListType, 0, 0)
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
func OpenSDSExtraToCinderExtra(profileExtras *model.ExtraSpec) *cindermodel.ExtraSpec {
	var typeExtra cindermodel.ExtraSpec

	if len(*profileExtras) >= 1 {
		typeExtra = make(map[string]interface{})
		for key, value := range *profileExtras {
			typeExtra[key] = value
		}
	}

	return &typeExtra
}
