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

	"github.com/astaxie/beego"
	docks "github.com/opensds/opensds/pkg/apiserver"
)

type DockController struct {
	beego.Controller
}

func (this *DockController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *DockController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	dockRequest := &docks.DockRequest{}
	result, err := docks.ListDocks(dockRequest)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("List docks failed!")
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *DockController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *DockController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}

type SpecifiedDockController struct {
	beego.Controller
}

func (this *SpecifiedDockController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *SpecifiedDockController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	dockRequest := &docks.DockRequest{
		Id: this.Ctx.Input.Param(":id"),
	}
	result, err := docks.GetDock(dockRequest)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Get dock failed!")
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *SpecifiedDockController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *SpecifiedDockController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}
