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

	pools "github.com/opensds/opensds/pkg/apiserver"

	"github.com/astaxie/beego"
)

type PoolController struct {
	beego.Controller
}

func (this *PoolController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *PoolController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	pr := &pools.PoolRequest{}
	result, err := pools.ListPools(pr)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("List storage pools failed!")
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(200)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *PoolController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *PoolController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}

type SpecifiedPoolController struct {
	beego.Controller
}

func (this *SpecifiedPoolController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *SpecifiedPoolController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	pr := &pools.PoolRequest{
		Id: this.Ctx.Input.Param(":id"),
	}
	result, err := pools.GetPool(pr)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Get storage pool failed!")
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(200)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *SpecifiedPoolController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *SpecifiedPoolController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}
