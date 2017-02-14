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
	"github.com/opensds/opensds/pkg/api/volumes"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var falseVolumeResponse api.VolumeResponse
var falseVolumeDetailResponse api.VolumeDetailResponse
var falseAllVolumesResponse []api.VolumeResponse
var falseAllVolumesDetailResponse api.VolumeDetailResponse

type VolumeController struct {
	beego.Controller
}

func (this *VolumeController) Get() {
	resourceType := this.Ctx.Input.Param(":resource")
	id := this.Ctx.Input.Param(":id")

	result, err := volumes.Show(resourceType, id)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString("Show volume failed!")
	} else {
		if reflect.DeepEqual(result, falseVolumeResponse) {
			log.Println("Show volume failed!")
			this.Ctx.WriteString("Show volume failed!")
		} else {
			rbody, _ := json.Marshal(result)
			this.Ctx.WriteString(string(rbody))
		}
	}
}

func (this *VolumeController) Put() {
	this.Ctx.WriteString("Not finished!")
}

func (this *VolumeController) Delete() {
	resourceType := this.Ctx.Input.Param(":resource")
	id := this.Ctx.Input.Param(":id")

	result, err := volumes.Delete(resourceType, id)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString("Delete volume failed!")
	} else {
		if result == "" {
			log.Println("Delete volume failed!")
			this.Ctx.WriteString("Delete volume failed!")
		} else {
			this.Ctx.WriteString(result)
		}
	}
}

func PostVolume(ctx *context.Context) {
	resourceType := ctx.Input.Param(":resource")
	rbody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Read volume request body failed:", err)
		ctx.WriteString("Read volume request body failed!")
	}

	var volumeRequest api.VolumeRequest
	volumeRequest.Resource_type = resourceType
	if err = json.Unmarshal(rbody, &volumeRequest); err != nil {
		log.Println("Parse volume request body failed:", err)
		ctx.WriteString("Parse volume request body failed!")
	}

	result, err := volumes.Create(volumeRequest.Resource_type,
		volumeRequest.Name,
		volumeRequest.Size)
	if err != nil {
		log.Println(err)
		ctx.WriteString("Create volume failed!")
	} else {
		if reflect.DeepEqual(result, falseVolumeResponse) {
			log.Println("Create volume failed!")
			ctx.WriteString("Create volume failed!")
		} else {
			rbody, _ := json.Marshal(result)
			ctx.WriteString(string(rbody))
		}
	}
}

func GetAllVolumes(ctx *context.Context) {
	resourceType := ctx.Input.Param(":resource")
	allowDetails := false

	result, err := volumes.List(resourceType, allowDetails)
	if err != nil {
		log.Println(err)
		ctx.WriteString("List volumes failed!")
	} else {
		if reflect.DeepEqual(result, falseAllVolumesResponse) {
			log.Println("List volumes failed!")
			ctx.WriteString("List volumes failed!")
		} else {
			rbody, _ := json.Marshal(result)
			ctx.WriteString(string(rbody))
		}
	}
}
