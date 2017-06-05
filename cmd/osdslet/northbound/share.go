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
	"io/ioutil"
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
	shares "github.com/opensds/opensds/pkg/apiserver"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type ShareController struct {
	beego.Controller
}

func (this *ShareController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	reqBody, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		log.Println("Read share request body failed:", err)
		resBody, _ := json.Marshal("Read share request body failed!")
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body(resBody)
	}

	shareRequest := &shares.ShareRequest{}
	if err = json.Unmarshal(reqBody, shareRequest); err != nil {
		log.Println("Parse share request body failed:", err)
		resBody, _ := json.Marshal("Parse share request body failed!")
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body(resBody)
	}

	result, err := shares.CreateShare(shareRequest)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Create share failed!")
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(201)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *ShareController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	resourceType := this.GetString("resource")

	shareRequest := &shares.ShareRequest{
		Profile: &api.StorageProfile{
			BackendDriver: resourceType,
		},
	}
	result, err := shares.ListShares(shareRequest)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("List shares failed!")
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(200)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *ShareController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *ShareController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}

type SpecifiedShareController struct {
	beego.Controller
}

func (this *SpecifiedShareController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *SpecifiedShareController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	shrId := this.Ctx.Input.Param(":id")
	resourceType := this.GetString("resource")

	shareRequest := &shares.ShareRequest{
		Schema: &api.ShareOperationSchema{
			Id: shrId,
		},
		Profile: &api.StorageProfile{
			BackendDriver: resourceType,
		},
	}

	result, err := shares.GetShare(shareRequest)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Get share failed!")
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.SetStatus(200)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *SpecifiedShareController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.SetStatus(501)
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *SpecifiedShareController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	shrId := this.Ctx.Input.Param(":id")
	reqBody, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		log.Println("Read share request body failed:", err)
		resBody, _ := json.Marshal("Read share request body failed!")
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body(resBody)
	}

	shareRequest := &shares.ShareRequest{}
	if err = json.Unmarshal(reqBody, shareRequest); err != nil {
		log.Println("Parse share request body failed:", err)
		resBody, _ := json.Marshal("Parse share request body failed!")
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body(resBody)
	}
	shareRequest.Schema.Id = shrId

	result := shares.DeleteShare(shareRequest)
	resBody, _ := json.Marshal(result)
	this.Ctx.Output.SetStatus(201)
	this.Ctx.Output.Body(resBody)
}

func AttachShare(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	reqBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Read share request body failed:", err)
		resBody, _ := json.Marshal("Read share request body failed!")
		ctx.Output.SetStatus(500)
		ctx.Output.Body(resBody)
	}

	shareRequest := &shares.ShareRequest{}
	if err = json.Unmarshal(reqBody, shareRequest); err != nil {
		log.Println("Parse share request body failed:", err)
		resBody, _ := json.Marshal("Parse share request body failed!")
		ctx.Output.SetStatus(500)
		ctx.Output.Body(resBody)
	}

	result := shares.AttachShare(shareRequest)
	resBody, _ := json.Marshal(result)
	ctx.Output.SetStatus(201)
	ctx.Output.Body(resBody)
}

func DetachShare(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	reqBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Read share request body failed:", err)
		resBody, _ := json.Marshal("Read share request body failed!")
		ctx.Output.SetStatus(500)
		ctx.Output.Body(resBody)
	}

	shareRequest := &shares.ShareRequest{}
	if err = json.Unmarshal(reqBody, shareRequest); err != nil {
		log.Println("Parse share request body failed:", err)
		resBody, _ := json.Marshal("Parse share request body failed!")
		ctx.Output.SetStatus(500)
		ctx.Output.Body(resBody)
	}

	result := shares.DetachShare(shareRequest)
	resBody, _ := json.Marshal(result)
	ctx.Output.SetStatus(201)
	ctx.Output.Body(resBody)
}
