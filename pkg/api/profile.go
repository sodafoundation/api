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

func (p *ProfilePortal) CreateProfile() {
	if !policy.Authorize(p.Ctx, "profile:create") {
		return
	}

	var profile = model.ProfileSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(p.Ctx.Request.Body).Decode(&profile); err != nil {
		reason := fmt.Sprintf("Parse profile request body failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Call db api module to handle create profile request.
	result, err := db.C.CreateProfile(c.GetContext(p.Ctx), &profile)
	if err != nil {
		reason := fmt.Sprintf("Create profile failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal profile created result failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorInternalServer)
		p.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	p.Ctx.Output.SetStatus(StatusOK)
	p.Ctx.Output.Body(body)
	return
}

func (p *ProfilePortal) ListProfiles() {
	if !policy.Authorize(p.Ctx, "profile:list") {
		return
	}

	m, err := p.GetParameters()
	if err != nil {
		reason := fmt.Sprintf("List profiles failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result, err := db.C.ListProfilesWithFilter(c.GetContext(p.Ctx), m)
	if err != nil {
		reason := fmt.Sprintf("List profiles failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal profiles listed result failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorInternalServer)
		p.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	p.Ctx.Output.SetStatus(StatusOK)
	p.Ctx.Output.Body(body)
	return
}

func (p *ProfilePortal) GetProfile() {
	if !policy.Authorize(p.Ctx, "profile:get") {
		return
	}
	id := p.Ctx.Input.Param(":profileId")

	result, err := db.C.GetProfile(c.GetContext(p.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("Get profile failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal profile got result failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorInternalServer)
		p.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	p.Ctx.Output.SetStatus(StatusOK)
	p.Ctx.Output.Body(body)
	return
}

func (p *ProfilePortal) UpdateProfile() {
	if !policy.Authorize(p.Ctx, "profile:update") {
		return
	}
	var profile = model.ProfileSpec{
		BaseModel: &model.BaseModel{},
	}
	id := p.Ctx.Input.Param(":profileId")

	if err := json.NewDecoder(p.Ctx.Request.Body).Decode(&profile); err != nil {
		reason := fmt.Sprintf("Parse profile request body failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result, err := db.C.UpdateProfile(c.GetContext(p.Ctx), id, &profile)
	if err != nil {
		reason := fmt.Sprintf("Update profiles failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal profile updated result failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorInternalServer)
		p.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	p.Ctx.Output.SetStatus(StatusOK)
	p.Ctx.Output.Body(body)
	return
}

func (p *ProfilePortal) DeleteProfile() {
	if !policy.Authorize(p.Ctx, "profile:delete") {
		return
	}
	id := p.Ctx.Input.Param(":profileId")
	ctx := c.GetContext(p.Ctx)
	profile, err := db.C.GetProfile(ctx, id)
	if err != nil {
		reason := fmt.Sprintf("Get profile failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	if err := db.C.DeleteProfile(ctx, profile.Id); err != nil {
		reason := fmt.Sprintf("Delete profiles failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	p.Ctx.Output.SetStatus(StatusOK)
	return
}

func (p *ProfilePortal) AddCustomProperty() {
	if !policy.Authorize(p.Ctx, "profile:add_custom_property") {
		return
	}
	var custom model.CustomPropertiesSpec
	id := p.Ctx.Input.Param(":profileId")

	if err := json.NewDecoder(p.Ctx.Request.Body).Decode(&custom); err != nil {
		reason := fmt.Sprintf("Parse custom properties request body failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result, err := db.C.AddCustomProperty(c.GetContext(p.Ctx), id, custom)
	if err != nil {
		reason := fmt.Sprintf("Add custom property failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal custom property added result failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorInternalServer)
		p.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	p.Ctx.Output.SetStatus(StatusOK)
	p.Ctx.Output.Body(body)
	return
}

func (p *ProfilePortal) ListCustomProperties() {
	if !policy.Authorize(p.Ctx, "profile:list_custom_properties") {
		return
	}
	id := p.Ctx.Input.Param(":profileId")

	result, err := db.C.ListCustomProperties(c.GetContext(p.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("List custom properties failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal custom properties listed result failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorInternalServer)
		p.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	p.Ctx.Output.SetStatus(StatusOK)
	p.Ctx.Output.Body(body)
	return
}

func (p *ProfilePortal) RemoveCustomProperty() {
	if !policy.Authorize(p.Ctx, "profile:remove_custom_property") {
		return
	}
	id := p.Ctx.Input.Param(":profileId")
	customKey := p.Ctx.Input.Param(":customKey")

	if err := db.C.RemoveCustomProperty(c.GetContext(p.Ctx), id, customKey); err != nil {
		reason := fmt.Sprintf("Remove custom property failed: %v", err)
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	p.Ctx.Output.SetStatus(StatusOK)
	return
}
