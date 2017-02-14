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
	"github.com/opensds/opensds/pkg/api/shares"

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
	resourceType := this.Ctx.Input.Param(":resource")
	id := this.Ctx.Input.Param(":id")

	result, err := shares.Show(resourceType, id)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString("Show share failed!")
	} else {
		if reflect.DeepEqual(result, falseShareResponse) {
			log.Println("Show share failed!")
			this.Ctx.WriteString("Show share failed!")
		} else {
			rbody, _ := json.Marshal(result)
			this.Ctx.WriteString(string(rbody))
		}
	}
}

func (this *ShareController) Put() {
	this.Ctx.WriteString("Not finished!")
}

func (this *ShareController) Delete() {
	resourceType := this.Ctx.Input.Param(":resource")
	id := this.Ctx.Input.Param(":id")

	result, err := shares.Delete(resourceType, id)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString("Delete share failed!")
	} else {
		if result == "" {
			log.Println("Delete share failed!")
			this.Ctx.WriteString("Delete share failed!")
		} else {
			this.Ctx.WriteString(result)
		}
	}
}

func PostShare(ctx *context.Context) {
	resourceType := ctx.Input.Param(":resource")
	rbody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Read share request body failed:", err)
		ctx.WriteString("Read share request body failed!")
	}

	var shareRequest api.ShareRequest
	shareRequest.Resource_type = resourceType
	if err = json.Unmarshal(rbody, &shareRequest); err != nil {
		log.Println("Parse volume request body failed:", err)
		ctx.WriteString("Parse volume request body failed!")
	}

	result, err := shares.Create(shareRequest.Resource_type,
		shareRequest.Name,
		shareRequest.Share_type,
		shareRequest.Share_proto,
		shareRequest.Size)
	if err != nil {
		log.Println(err)
		ctx.WriteString("Create share failed!")
	} else {
		if reflect.DeepEqual(result, falseShareResponse) {
			log.Println("Create share failed!")
			ctx.WriteString("Create share failed!")
		} else {
			rbody, _ := json.Marshal(result)
			ctx.WriteString(string(rbody))
		}
	}
}

func GetAllShares(ctx *context.Context) {
	resourceType := ctx.Input.Param(":resource")
	allowDetails := false

	result, err := shares.List(resourceType, allowDetails)
	if err != nil {
		log.Println(err)
		ctx.WriteString("List shares failed!")
	} else {
		if reflect.DeepEqual(result, falseAllSharesResponse) {
			log.Println("List shares failed!")
			ctx.WriteString("List shares failed!")
		} else {
			rbody, _ := json.Marshal(result)
			ctx.WriteString(string(rbody))
		}
	}
}
