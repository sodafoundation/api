// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

/*
This module implements a entry into the OpenSDS northbound service.

*/

package api

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/opa"
	"github.com/opensds/opensds/pkg/utils"
)

func init() {
	var input = []*model.ProfileSpec{}
	if err := opa.RegisterData(&input); err != nil {
		panic(err)
	}
}

type ProfilePortal struct {
	beego.Controller
}

func (this *ProfilePortal) CreateProfile() {
	var profile = model.ProfileSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&profile); err != nil {
		reason := fmt.Sprintf("Parse profile request body failed: %s", err.Error())
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	// If profile uuid and created time is null, generate it randomly.
	if err := utils.ValidateData(&profile, utils.S); err != nil {
		reason := fmt.Sprintf("Validate profile data failed: %s", err.Error())
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	// Call db api module to handle create profile request.
	if err := db.C.CreateProfile(&profile); err != nil {
		reason := fmt.Sprintf("Create profile failed: %s", err.Error())
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	if err := opa.PatchData(&profile, "add", "-"); err != nil {
		reason := fmt.Sprintf("Patch profile data failed: %s", err.Error())
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(&profile)
	if err != nil {
		reason := fmt.Sprintf("Marshal profile created result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)
	this.Ctx.Output.Body(body)
	return
}

func (this *ProfilePortal) ListProfiles() {
	result, err := db.C.ListProfiles()
	if err != nil {
		reason := fmt.Sprintf("List profiles failed: %s", err.Error())
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal profiles listed result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

type SpecifiedProfilePortal struct {
	beego.Controller
}

func (this *SpecifiedProfilePortal) GetProfile() {
	id := this.Ctx.Input.Param(":profileId")

	result, err := db.C.GetProfile(id)
	if err != nil {
		reason := fmt.Sprintf("Get profiles failed: %s", err.Error())
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal profile showed result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *SpecifiedProfilePortal) UpdateProfile() {
	var profile = model.ProfileSpec{
		BaseModel: &model.BaseModel{},
	}
	id := this.Ctx.Input.Param(":profileId")

	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(&profile); err != nil {
		reason := fmt.Sprintf("Parse profile request body failed: %s", err.Error())
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	result, err := db.C.UpdateProfile(id, &profile)
	if err != nil {
		reason := fmt.Sprintf("Update profiles failed: %s", err.Error())
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal profile updated result failed: %s", err.Error())
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)
	this.Ctx.Output.Body(body)
	return
}

func (this *SpecifiedProfilePortal) DeleteProfile() {
	id := this.Ctx.Input.Param(":profileId")

	if err := db.C.DeleteProfile(id); err != nil {
		reason := fmt.Sprintf("Delete profiles failed: %s", err.Error())
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusAccepted)
	this.Ctx.Output.Body([]byte("Delete profile success!"))
	return
}

type ProfileExtrasPortal struct {
	beego.Controller
}

func (this *ProfileExtrasPortal) AddExtraProperty() {
	var extra model.ExtraSpec
	id := this.Ctx.Input.Param(":profileId")

	if err := json.NewDecoder(this.Ctx.Request.Body).Decode(extra); err != nil {
		log.Error("Parse extra request body failed:", err)
		resBody, _ := json.Marshal("Parse extra request body failed!")
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(resBody)
		return
	}

	result, err := db.C.AddExtraProperty(id, extra)
	if err != nil {
		log.Error(err)
		resBody, _ := json.Marshal("Create extra property failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
		return
	}

	resBody, _ := json.Marshal(result)
	this.Ctx.Output.SetStatus(StatusAccepted)
	this.Ctx.Output.Body(resBody)
	return
}

func (this *ProfileExtrasPortal) ListExtraProperties() {
	id := this.Ctx.Input.Param(":profileId")

	result, err := db.C.ListExtraProperties(id)
	if err != nil {
		log.Error(err)
		resBody, _ := json.Marshal("List extra properties failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
		return
	}

	resBody, _ := json.Marshal(result)
	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(resBody)
	return
}

func (this *ProfileExtrasPortal) RemoveExtraProperty() {
	id := this.Ctx.Input.Param(":profileId")
	extraKey := this.Ctx.Input.Param(":extraKey")

	if err := db.C.RemoveExtraProperty(id, extraKey); err != nil {
		log.Error(err)
		resBody, _ := json.Marshal("Remove profile extra property failed: " + fmt.Sprint(err))
		this.Ctx.Output.SetStatus(StatusBadRequest)
		this.Ctx.Output.Body(resBody)
		return
	}

	resBody, _ := json.Marshal("Remove extra property success!")
	this.Ctx.Output.SetStatus(StatusAccepted)
	this.Ctx.Output.Body(resBody)
	return
}
