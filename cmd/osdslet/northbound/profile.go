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

package northbound

import (
	"encoding/json"
	"log"

	profiles "github.com/opensds/opensds/pkg/apiserver"

	"github.com/astaxie/beego"
)

type ProfileController struct {
	beego.Controller
}

func (this *ProfileController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *ProfileController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	pr := &profiles.ProfileRequest{}
	result, err := profiles.ListProfiles(pr)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("List storage profiles failed!")
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(200)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *ProfileController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *ProfileController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}

type SpecifiedProfileController struct {
	beego.Controller
}

func (this *SpecifiedProfileController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *SpecifiedProfileController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	pr := &profiles.ProfileRequest{
		Id: this.Ctx.Input.Param(":id"),
	}
	result, err := profiles.GetProfile(pr)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Get storage profile failed!")
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(200)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *SpecifiedProfileController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *SpecifiedProfileController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}
