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
	"reflect"

	"github.com/opensds/opensds/pkg/api"
	"github.com/opensds/opensds/pkg/api/v1/shares"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var falseShareResponse api.ShareResponse
var falseShareDetailResponse api.ShareDetailResponse
var falseAllSharesResponse []api.ShareResponse
var falseAllSharesDetailResponse []api.ShareDetailResponse

type ShareController struct {
	beego.Controller
}

func (this *ShareController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	resourceType := this.Ctx.Input.Param(":resource")
	id := this.Ctx.Input.Param(":id")

	shareRequest := shares.ShareRequest{
		ResourceType: resourceType,
		Id:           id,
	}
	result, err := shares.GetShare(shareRequest)
	if err != nil {
		log.Println(err)
		rbody, _ := json.Marshal("Show share failed!")
		this.Ctx.Output.Body(rbody)
	} else {
		if reflect.DeepEqual(result, falseShareResponse) {
			log.Println("Show share failed!")
			rbody, _ := json.Marshal("Show share failed!")
			this.Ctx.Output.Body(rbody)
		} else {
			rbody, _ := json.Marshal(result)
			this.Ctx.Output.Body(rbody)
		}
	}
}

func (this *ShareController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	rbody, _ := json.Marshal("Not supported!")
	this.Ctx.Output.Body(rbody)
}

func (this *ShareController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	resourceType := this.Ctx.Input.Param(":resource")
	id := this.Ctx.Input.Param(":id")

	shareRequest := shares.ShareRequest{
		ResourceType: resourceType,
		Id:           id,
	}
	result := shares.DeleteShare(shareRequest)
	rbody, _ := json.Marshal(result)
	this.Ctx.Output.Body(rbody)
}

func PostShare(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	resourceType := ctx.Input.Param(":resource")
	rbody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Read share request body failed:", err)
		rbody, _ := json.Marshal("Read share request body failed!")
		ctx.Output.Body(rbody)
	}

	var shareRequest shares.ShareRequest
	shareRequest.ResourceType = resourceType
	if err = json.Unmarshal(rbody, &shareRequest); err != nil {
		log.Println("Parse volume request body failed:", err)
		rbody, _ := json.Marshal("Parse share request body failed!")
		ctx.Output.Body(rbody)
	}

	result, err := shares.CreateShare(shareRequest)
	if err != nil {
		log.Println(err)
		rbody, _ := json.Marshal("Create share failed!")
		ctx.Output.Body(rbody)
	} else {
		if reflect.DeepEqual(result, falseShareResponse) {
			log.Println("Create share failed!")
			rbody, _ := json.Marshal("Create share failed!")
			ctx.Output.Body(rbody)
		} else {
			rbody, _ := json.Marshal(result)
			ctx.Output.Body(rbody)
		}
	}
}

func GetAllShares(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	resourceType := ctx.Input.Param(":resource")

	shareRequest := shares.ShareRequest{
		ResourceType: resourceType,
		AllowDetails: false,
	}
	result, err := shares.ListShares(shareRequest)
	if err != nil {
		log.Println(err)
		rbody, _ := json.Marshal("List shares failed!")
		ctx.Output.Body(rbody)
	} else {
		if reflect.DeepEqual(result, falseAllSharesResponse) {
			log.Println("List shares failed!")
			rbody, _ := json.Marshal("List shares failed!")
			ctx.Output.Body(rbody)
		} else {
			rbody, _ := json.Marshal(result)
			ctx.Output.Body(rbody)
		}
	}
}
