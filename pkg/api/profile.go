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
This module implements a entry into the OpenSDS northbound service.

*/

package api

import (
	"encoding/json"
	"fmt"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/policy"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
)

type ProfilePortal struct {
	BasePortal
}

func (this *ProfilePortal) CreateProfile() {
	if !policy.Authorize(this.Ctx, "profile:create") {
		return
	}

	var profile = model.ProfileSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&profile); err != nil {
		reason := fmt.Sprintf("Parse profile request body failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Call db api module to handle create profile request.
	result, err := db.C.CreateProfile(c.GetContext(this.Ctx), &profile)
	if err != nil {
		reason := fmt.Sprintf("Create profile failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal profile created result failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *ProfilePortal) ListProfiles() {
	if !policy.Authorize(this.Ctx, "profile:list") {
		return
	}

	m, err := this.GetParameters()
	if err != nil {
		reason := fmt.Sprintf("List profiles failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result, err := db.C.ListProfilesWithFilter(c.GetContext(this.Ctx), m)
	if err != nil {
		reason := fmt.Sprintf("List profiles failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal profiles listed result failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *ProfilePortal) GetProfile() {
	if !policy.Authorize(this.Ctx, "profile:get") {
		return
	}
	id := this.Ctx.Input.Param(":profileId")

	result, err := db.C.GetProfile(c.GetContext(this.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("Get profile failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal profile got result failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *ProfilePortal) UpdateProfile() {
	if !policy.Authorize(this.Ctx, "profile:update") {
		return
	}
	var profile = model.ProfileSpec{
		BaseModel: &model.BaseModel{},
	}
	id := this.Ctx.Input.Param(":profileId")

	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&profile); err != nil {
		reason := fmt.Sprintf("Parse profile request body failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result, err := db.C.UpdateProfile(c.GetContext(this.Ctx), id, &profile)
	if err != nil {
		reason := fmt.Sprintf("Update profiles failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal profile updated result failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *ProfilePortal) DeleteProfile() {
	if !policy.Authorize(this.Ctx, "profile:delete") {
		return
	}
	id := this.Ctx.Input.Param(":profileId")
	ctx := c.GetContext(this.Ctx)
	profile, err := db.C.GetProfile(ctx, id)
	if err != nil {
		reason := fmt.Sprintf("Get profile failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	if err := db.C.DeleteProfile(ctx, profile.Id); err != nil {
		reason := fmt.Sprintf("Delete profiles failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	return
}

func (this *ProfilePortal) AddCustomProperty() {
	if !policy.Authorize(this.Ctx, "profile:add_custom_property") {
		return
	}
	var custom model.CustomPropertiesSpec
	id := this.Ctx.Input.Param(":profileId")

	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&custom); err != nil {
		reason := fmt.Sprintf("Parse custom properties request body failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result, err := db.C.AddCustomProperty(c.GetContext(this.Ctx), id, custom)
	if err != nil {
		reason := fmt.Sprintf("Add custom property failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal custom property added result failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *ProfilePortal) ListCustomProperties() {
	if !policy.Authorize(this.Ctx, "profile:list_custom_properties") {
		return
	}
	id := this.Ctx.Input.Param(":profileId")

	result, err := db.C.ListCustomProperties(c.GetContext(this.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("List custom properties failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal custom properties listed result failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *ProfilePortal) RemoveCustomProperty() {
	if !policy.Authorize(this.Ctx, "profile:remove_custom_property") {
		return
	}
	id := this.Ctx.Input.Param(":profileId")
	customKey := this.Ctx.Input.Param(":customKey")

	if err := db.C.RemoveCustomProperty(c.GetContext(this.Ctx), id, customKey); err != nil {
		reason := fmt.Sprintf("Remove custom property failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	return
}
