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

	"github.com/astaxie/beego"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/policy"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
)

type ProfilePortal struct {
	beego.Controller
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
		reason := fmt.Sprintf("Parse profile request body failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Call db api module to handle create profile request.
	result, err := db.C.CreateProfile(c.GetContext(this.Ctx), &profile)
	if err != nil {
		reason := fmt.Sprintf("Create profile failed: %s", err.Error())
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
	result, err := db.C.ListProfiles(c.GetContext(this.Ctx))
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
		reason := fmt.Sprintf("Get profile failed: %s", err.Error())
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

func (this *ProfilePortal) AddExtraProperty() {
	if !policy.Authorize(this.Ctx, "profile:add_extra_property") {
		return
	}
	var extra model.ExtraSpec
	id := this.Ctx.Input.Param(":profileId")

	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&extra); err != nil {
		reason := fmt.Sprintf("Parse extra request body failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result, err := db.C.AddExtraProperty(c.GetContext(this.Ctx), id, extra)
	if err != nil {
		reason := fmt.Sprintf("Create extra property failed: %s", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal extra property added result failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *ProfilePortal) ListExtraProperties() {
	if !policy.Authorize(this.Ctx, "profile:list_extra_properties") {
		return
	}
	id := this.Ctx.Input.Param(":profileId")

	result, err := db.C.ListExtraProperties(c.GetContext(this.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("List extra properties failed: %s", err)
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal extra properties listed result failed: %v", err)
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *ProfilePortal) RemoveExtraProperty() {
	if !policy.Authorize(this.Ctx, "profile:remove_extra_property") {
		return
	}
	id := this.Ctx.Input.Param(":profileId")
	extraKey := this.Ctx.Input.Param(":extraKey")

	if err := db.C.RemoveExtraProperty(c.GetContext(this.Ctx), id, extraKey); err != nil {
		reason := fmt.Sprintf("Remove extra property failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	return
}
