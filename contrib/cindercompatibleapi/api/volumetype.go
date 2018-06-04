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

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/astaxie/beego"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/cindercompatibleapi/converter"
	"github.com/opensds/opensds/pkg/model"
)

// TypePortal ...
type TypePortal struct {
	beego.Controller
}

// DefaultTypeName ...
var DefaultTypeName = "default"

// UpdateType ...
func (portal *TypePortal) UpdateType() {
	id := portal.Ctx.Input.Param(":volumeTypeId")
	var cinderReq = converter.UpdateTypeReqSpec{}
	if err := json.NewDecoder(portal.Ctx.Request.Body).Decode(&cinderReq); err != nil {
		reason := fmt.Sprintf("Update a volume type, parse request body failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	profile, err := converter.UpdateTypeReq(&cinderReq)
	if err != nil {
		reason := fmt.Sprintf("Update a volume type failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	profile, err = client.UpdateProfile(id, profile)
	if err != nil {
		reason := fmt.Sprintf("Update a volume type failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.UpdateTypeResp(profile)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Update a volume type, marshal result failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// AddExtraProperty ...
func (portal *TypePortal) AddExtraProperty() {
	id := portal.Ctx.Input.Param(":volumeTypeId")
	var cinderReq = converter.AddExtraReqSpec{}
	if err := json.NewDecoder(portal.Ctx.Request.Body).Decode(&cinderReq); err != nil {
		reason := fmt.Sprintf("Create or update extra specs for volume type, parse request body failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	profileExtra := converter.AddExtraReq(&cinderReq)
	profileExtra, err := client.AddExtraProperty(id, profileExtra)
	if err != nil {
		reason := fmt.Sprintf("Create or update extra specs for volume type failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.AddExtraResp(profileExtra)
	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Create or update extra specs for volume type, marshal result failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// ListExtraProperties ...
func (portal *TypePortal) ListExtraProperties() {
	id := portal.Ctx.Input.Param(":volumeTypeId")
	profileExtra, err := client.ListExtraProperties(id)

	if err != nil {
		reason := fmt.Sprintf("Show all extra specifications for volume type failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.ShowAllExtraResp(profileExtra)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Show all extra specifications for volume type, marshal result failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// ShowExtraProperty ...
func (portal *TypePortal) ShowExtraProperty() {
	id := portal.Ctx.Input.Param(":volumeTypeId")
	profileExtra, err := client.ListExtraProperties(id)

	if err != nil {
		reason := fmt.Sprintf("Show extra specification for volume type failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	key := portal.Ctx.Input.Param(":key")
	result := converter.ShowExtraResp(key, profileExtra)
	if nil == (*result) {
		reason := "The key name of the extra spec for the volume type can not be found"
		portal.Ctx.Output.SetStatus(http.StatusNotFound)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	body, err := json.Marshal(result)

	if err != nil {
		reason := fmt.Sprintf("Show extra specification for volume type, marshal result failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// UpdateExtraProperty ...
func (portal *TypePortal) UpdateExtraProperty() {
	id := portal.Ctx.Input.Param(":volumeTypeId")
	key := portal.Ctx.Input.Param(":key")
	var cinderReq = converter.UpdateExtraReqSpec{}

	if err := json.NewDecoder(portal.Ctx.Request.Body).Decode(&cinderReq); err != nil {
		reason := fmt.Sprintf("Update extra specification for volume type, parse request body failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	profileExtra, err := converter.UpdateExtraReq(key, &cinderReq)
	if err != nil {
		reason := fmt.Sprintf("Update extra specification for volume type failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	profileExtra, err = client.AddExtraProperty(id, profileExtra)
	if err != nil {
		reason := fmt.Sprintf("Update extra specification for volume type failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.UpdateExtraResp(key, profileExtra)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Update extra specification for volume type, marshal result failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return

}

// DeleteExtraProperty ...
func (portal *TypePortal) DeleteExtraProperty() {
	id := portal.Ctx.Input.Param(":volumeTypeId")
	key := portal.Ctx.Input.Param(":key")
	err := client.RemoveExtraProperty(id, key)

	if err != nil {
		reason := fmt.Sprintf("Delete extra specification for volume type failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusAccepted)
	return
}

// GetType ...
func (portal *TypePortal) GetType() {
	id := portal.Ctx.Input.Param(":volumeTypeId")
	DefaultName := os.Getenv("DEFAULT_VOLUME_TYPE_NAME")
	if ("" != DefaultName) && (DefaultTypeName != DefaultName) {
		DefaultTypeName = DefaultName
		log.Info("DefaultTypeName = " + DefaultTypeName)
	}

	var profile *model.ProfileSpec

	if "default" != id {
		foundProfile, err := client.GetProfile(id)

		if err != nil {
			reason := fmt.Sprintf("Get profile failed: %v", err)
			portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
			portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
			log.Error(reason)
			return
		}

		profile = foundProfile
	} else {
		profiles, err := client.ListProfiles()
		if err != nil {
			reason := fmt.Sprintf("List profiles failed: %v", err)
			portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
			portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
			log.Error(reason)
			return
		}

		for _, v := range profiles {
			if DefaultTypeName == v.Name {
				profile = v
			}
		}

		if nil == profile {
			reason := "Default volume type can not be found"
			portal.Ctx.Output.SetStatus(http.StatusNotFound)
			portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
			log.Error(reason)
			return
		}
	}

	result := converter.ShowTypeResp(profile)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Show volume type detail, marshal result failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// DeleteType ...
func (portal *TypePortal) DeleteType() {
	id := portal.Ctx.Input.Param(":volumeTypeId")
	err := client.DeleteProfile(id)

	if err != nil {
		reason := fmt.Sprintf("Delete a volume type failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusAccepted)
	return
}

// ListTypes ...
func (portal *TypePortal) ListTypes() {
	profiles, err := client.ListProfiles()
	if err != nil {
		reason := fmt.Sprintf("List all volume types failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.ListTypesResp(profiles)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("List all volume types, marshal result failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// CreateType ...
func (portal *TypePortal) CreateType() {
	var cinderReq = converter.CreateTypeReqSpec{}
	if err := json.NewDecoder(portal.Ctx.Request.Body).Decode(&cinderReq); err != nil {
		reason := fmt.Sprintf("Create a volume type, parse request body failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	profile, err := converter.CreateTypeReq(&cinderReq)
	if err != nil {
		reason := fmt.Sprintf("Create a volume type failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	profile, err = client.CreateProfile(profile)
	if err != nil {
		reason := fmt.Sprintf("Create a volume type failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.CreateTypeResp(profile)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Create a volume type, marshal result failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return
}
