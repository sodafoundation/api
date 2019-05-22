// Copyright 2019 The OpenSDS Authors.
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

package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opensds/opensds/pkg/api/policy"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/constants"
)

type ProfilePortal struct {
	BasePortal
}

var StorageAccessCapabilities = []string{constants.Read, constants.Write, constants.Execute}

func (p *ProfilePortal) CreateProfile() {
	if !policy.Authorize(p.Ctx, "profile:create") {
		return
	}

	var profile = model.ProfileSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(p.Ctx.Request.Body).Decode(&profile); err != nil {
		errMsg := fmt.Sprintf("parse profile request body failed: %v", err)
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// Validate StorageType is block or file
	stype := strings.ToLower(profile.StorageType)
	switch stype {
	case constants.Block:
		break
	case constants.File:
		if pp := profile.ProvisioningProperties; pp.IsEmpty() {
			errMsg := fmt.Sprintf("parse profile request body failed : Please provide ProvisionProperties")
			p.ErrorHandle(model.ErrorBadRequest, errMsg)
			return
		} else {
			if ds := pp.DataStorage; !ds.IsEmpty() {
				if len(ds.StorageAccessCapability) == 0 {
					errMsg := fmt.Sprintf("parse profile request body failed %v is invalid storageaccesscapability",
						ds.StorageAccessCapability)
					p.ErrorHandle(model.ErrorBadRequest, errMsg)
					return
				} else {
					for _, x := range ds.StorageAccessCapability {
						if !utils.Contains(StorageAccessCapabilities, x) {
							errMsg := fmt.Sprintf("parse profile request body failed %v is invalid storageaccesscapability",
								ds.StorageAccessCapability)
							p.ErrorHandle(model.ErrorBadRequest, errMsg)
							return
						}
					}
				}

			}
		}
		break
	default:
		errMsg := fmt.Sprintf("parse profile request body failed : %v is invalid storagetype", stype)
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// Call db api module to handle create profile request.
	result, err := db.C.CreateProfile(c.GetContext(p.Ctx), &profile)
	if err != nil {
		errMsg := fmt.Sprintf("create profile failed: %v", err)
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal profile created result failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, body)
	return
}

func (p *ProfilePortal) ListProfiles() {
	if !policy.Authorize(p.Ctx, "profile:list") {
		return
	}

	m, err := p.GetParameters()
	if err != nil {
		errMsg := fmt.Sprintf("list profiles failed: %v", err)
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	result, err := db.C.ListProfilesWithFilter(c.GetContext(p.Ctx), m)
	if err != nil {
		errMsg := fmt.Sprintf("list profiles failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal profiles listed result failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, body)
	return
}

func (p *ProfilePortal) GetProfile() {
	if !policy.Authorize(p.Ctx, "profile:get") {
		return
	}
	id := p.Ctx.Input.Param(":profileId")

	result, err := db.C.GetProfile(c.GetContext(p.Ctx), id)
	if err != nil {
		errMsg := fmt.Sprintf("profile %s not found: %v", id, err)
		p.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal profile got result failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, body)
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
		errMsg := fmt.Sprintf("parse profile request body failed: %v", err)
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	result, err := db.C.UpdateProfile(c.GetContext(p.Ctx), id, &profile)
	if err != nil {
		errMsg := fmt.Sprintf("update profiles failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal profile updated result failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, body)
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
		errMsg := fmt.Sprintf("profile %s not found: %v", id, err)
		p.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	if err = db.C.DeleteProfile(ctx, profile.Id); err != nil {
		errMsg := fmt.Sprintf("delete profiles failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, nil)
	return
}

func (p *ProfilePortal) AddCustomProperty() {

	if !policy.Authorize(p.Ctx, "profile:add_custom_property") {
		return
	}
	var custom model.CustomPropertiesSpec
	id := p.Ctx.Input.Param(":profileId")

	if err := json.NewDecoder(p.Ctx.Request.Body).Decode(&custom); err != nil {
		errMsg := fmt.Sprintf("parse custom properties request body failed: %v", err)
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	result, err := db.C.AddCustomProperty(c.GetContext(p.Ctx), id, custom)
	if err != nil {
		errMsg := fmt.Sprintf("add custom property failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal custom property added result failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, body)
	return
}

func (p *ProfilePortal) ListCustomProperties() {

	if !policy.Authorize(p.Ctx, "profile:list_custom_properties") {
		return
	}
	id := p.Ctx.Input.Param(":profileId")

	result, err := db.C.ListCustomProperties(c.GetContext(p.Ctx), id)
	if err != nil {
		errMsg := fmt.Sprintf("list custom properties failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal custom properties listed result failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, body)
	return
}

func (p *ProfilePortal) RemoveCustomProperty() {

	if !policy.Authorize(p.Ctx, "profile:remove_custom_property") {
		return
	}
	id := p.Ctx.Input.Param(":profileId")
	customKey := p.Ctx.Input.Param(":customKey")

	if err := db.C.RemoveCustomProperty(c.GetContext(p.Ctx), id, customKey); err != nil {
		errMsg := fmt.Sprintf("remove custom property failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, nil)
	return
}
