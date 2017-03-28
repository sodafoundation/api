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

	"github.com/opensds/opensds/pkg/controller/api"
	"github.com/opensds/opensds/pkg/controller/api/v1/shares"

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

func (this *ShareController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	resourceType := this.Ctx.Input.Param(":resource")
	reqBody, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		log.Println("Read share request body failed:", err)
		resBody, _ := json.Marshal("Read share request body failed!")
		this.Ctx.Output.Body(resBody)
	}

	shareRequest := &shares.ShareRequest{
		ResourceType: resourceType,
	}
	if err = json.Unmarshal(reqBody, shareRequest); err != nil {
		log.Println("Parse share request body failed:", err)
		resBody, _ := json.Marshal("Parse share request body failed!")
		this.Ctx.Output.Body(resBody)
	}

	result, err := shares.CreateShare(shareRequest)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Create share failed!")
		this.Ctx.Output.Body(resBody)
	} else {
		if reflect.DeepEqual(result, falseShareResponse) {
			log.Println("Create share failed!")
			resBody, _ := json.Marshal("Create share failed!")
			this.Ctx.Output.Body(resBody)
		} else {
			resBody, _ := json.Marshal(result)
			this.Ctx.Output.Body(resBody)
		}
	}
}

func (this *ShareController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	resourceType := this.Ctx.Input.Param(":resource")

	shareRequest := &shares.ShareRequest{
		ResourceType: resourceType,
		AllowDetails: false,
	}
	result, err := shares.ListShares(shareRequest)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("List shares failed!")
		this.Ctx.Output.Body(resBody)
	} else {
		if reflect.DeepEqual(result, falseAllSharesResponse) {
			log.Println("List shares failed!")
			resBody, _ := json.Marshal("List shares failed!")
			this.Ctx.Output.Body(resBody)
		} else {
			resBody, _ := json.Marshal(result)
			this.Ctx.Output.Body(resBody)
		}
	}
}

func (this *ShareController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	resBody, _ := json.Marshal("Not supported!")
	this.Ctx.Output.Body(resBody)
}

func (this *ShareController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	resourceType := this.Ctx.Input.Param(":resource")
	reqBody, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		log.Println("Read share request body failed:", err)
		resBody, _ := json.Marshal("Read share request body failed!")
		this.Ctx.Output.Body(resBody)
	}

	shareRequest := &shares.ShareRequest{
		ResourceType: resourceType,
	}
	if err = json.Unmarshal(reqBody, shareRequest); err != nil {
		log.Println("Parse share request body failed:", err)
		resBody, _ := json.Marshal("Parse share request body failed!")
		this.Ctx.Output.Body(resBody)
	}

	result := shares.DeleteShare(shareRequest)
	resBody, _ := json.Marshal(result)
	this.Ctx.Output.Body(resBody)
}

func GetShare(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	resourceType := ctx.Input.Param(":resource")
	volId := ctx.Input.Param(":id")

	shareRequest := &shares.ShareRequest{
		ResourceType: resourceType,
		Id:           volId,
	}
	result, err := shares.GetShare(shareRequest)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Get share failed!")
		ctx.Output.Body(resBody)
	} else {
		if reflect.DeepEqual(result, falseAllSharesResponse) {
			log.Println("Get share failed!")
			resBody, _ := json.Marshal("Get share failed!")
			ctx.Output.Body(resBody)
		} else {
			resBody, _ := json.Marshal(result)
			ctx.Output.Body(resBody)
		}
	}
}

func AttachShare(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	reqBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Read share request body failed:", err)
		resBody, _ := json.Marshal("Read share request body failed!")
		ctx.Output.Body(resBody)
	}

	shareRequest := &shares.ShareRequest{}
	if err = json.Unmarshal(reqBody, shareRequest); err != nil {
		log.Println("Parse share request body failed:", err)
		resBody, _ := json.Marshal("Parse share request body failed!")
		ctx.Output.Body(resBody)
	}

	result := shares.AttachShare(shareRequest)
	resBody, _ := json.Marshal(result)
	ctx.Output.Body(resBody)
}

func DetachShare(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	reqBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Read share request body failed:", err)
		resBody, _ := json.Marshal("Read share request body failed!")
		ctx.Output.Body(resBody)
	}

	shareRequest := &shares.ShareRequest{}
	if err = json.Unmarshal(reqBody, shareRequest); err != nil {
		log.Println("Parse share request body failed:", err)
		resBody, _ := json.Marshal("Parse share request body failed!")
		ctx.Output.Body(resBody)
	}

	result := shares.DetachShare(shareRequest)
	resBody, _ := json.Marshal(result)
	ctx.Output.Body(resBody)
}

func MountShare(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	reqBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Read share request body failed:", err)
		resBody, _ := json.Marshal("Read share request body failed!")
		ctx.Output.Body(resBody)
	}

	shareRequest := &shares.ShareRequest{}
	if err = json.Unmarshal(reqBody, shareRequest); err != nil {
		log.Println("Parse share request body failed:", err)
		resBody, _ := json.Marshal("Parse share request body failed!")
		ctx.Output.Body(resBody)
	}

	result := shares.MountShare(shareRequest)
	resBody, _ := json.Marshal(result)
	ctx.Output.Body(resBody)
}

func UnmountShare(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	reqBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Read share request body failed:", err)
		resBody, _ := json.Marshal("Read share request body failed!")
		ctx.Output.Body(resBody)
	}

	shareRequest := &shares.ShareRequest{}
	if err = json.Unmarshal(reqBody, shareRequest); err != nil {
		log.Println("Parse share request body failed:", err)
		resBody, _ := json.Marshal("Parse share request body failed!")
		ctx.Output.Body(resBody)
	}

	result := shares.UnmountShare(shareRequest)
	resBody, _ := json.Marshal(result)
	ctx.Output.Body(resBody)
}
