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
	"github.com/opensds/opensds/pkg/controller/api/v1/volumes"

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

func (this *VolumeController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	resourceType := this.Ctx.Input.Param(":resource")
	reqBody, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		log.Println("Read volume request body failed:", err)
		resBody, _ := json.Marshal("Read volume request body failed!")
		this.Ctx.Output.Body(resBody)
	}

	volumeRequest := &volumes.VolumeRequest{
		ResourceType: resourceType,
	}
	if err = json.Unmarshal(reqBody, volumeRequest); err != nil {
		log.Println("Parse volume request body failed:", err)
		resBody, _ := json.Marshal("Parse volume request body failed!")
		this.Ctx.Output.Body(resBody)
	}

	result, err := volumes.CreateVolume(volumeRequest)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Create volume failed!")
		this.Ctx.Output.Body(resBody)
	} else {
		if reflect.DeepEqual(result, falseVolumeResponse) {
			log.Println("Create volume failed!")
			resBody, _ := json.Marshal("Create volume failed!")
			this.Ctx.Output.Body(resBody)
		} else {
			resBody, _ := json.Marshal(result)
			this.Ctx.Output.Body(resBody)
		}
	}
}

func (this *VolumeController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	resourceType := this.Ctx.Input.Param(":resource")

	volumeRequest := &volumes.VolumeRequest{
		ResourceType: resourceType,
		AllowDetails: false,
	}
	result, err := volumes.ListVolumes(volumeRequest)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("List volumes failed!")
		this.Ctx.Output.Body(resBody)
	} else {
		if reflect.DeepEqual(result, falseAllVolumesResponse) {
			log.Println("List volumes failed!")
			resBody, _ := json.Marshal("List volumes failed!")
			this.Ctx.Output.Body(resBody)
		} else {
			resBody, _ := json.Marshal(result)
			this.Ctx.Output.Body(resBody)
		}
	}
}

func (this *VolumeController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	resBody, _ := json.Marshal("Not supported!")
	this.Ctx.Output.Body(resBody)
}

func (this *VolumeController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	resourceType := this.Ctx.Input.Param(":resource")
	reqBody, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		log.Println("Read volume request body failed:", err)
		resBody, _ := json.Marshal("Read volume request body failed!")
		this.Ctx.Output.Body(resBody)
	}

	volumeRequest := &volumes.VolumeRequest{
		ResourceType: resourceType,
	}
	if err = json.Unmarshal(reqBody, volumeRequest); err != nil {
		log.Println("Parse volume request body failed:", err)
		resBody, _ := json.Marshal("Parse volume request body failed!")
		this.Ctx.Output.Body(resBody)
	}

	result := volumes.DeleteVolume(volumeRequest)
	resBody, _ := json.Marshal(result)
	this.Ctx.Output.Body(resBody)
}

func GetVolume(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	resourceType := ctx.Input.Param(":resource")
	volId := ctx.Input.Param(":id")

	volumeRequest := &volumes.VolumeRequest{
		ResourceType: resourceType,
		Id:           volId,
	}
	result, err := volumes.GetVolume(volumeRequest)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Get volume failed!")
		ctx.Output.Body(resBody)
	} else {
		if reflect.DeepEqual(result, falseAllVolumesResponse) {
			log.Println("Get volume failed!")
			resBody, _ := json.Marshal("Get volume failed!")
			ctx.Output.Body(resBody)
		} else {
			resBody, _ := json.Marshal(result)
			ctx.Output.Body(resBody)
		}
	}
}

func AttachVolume(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	reqBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Read volume request body failed:", err)
		resBody, _ := json.Marshal("Read volume request body failed!")
		ctx.Output.Body(resBody)
	}

	volumeRequest := &volumes.VolumeRequest{}
	if err = json.Unmarshal(reqBody, volumeRequest); err != nil {
		log.Println("Parse volume request body failed:", err)
		resBody, _ := json.Marshal("Parse volume request body failed!")
		ctx.Output.Body(resBody)
	}

	result := volumes.AttachVolume(volumeRequest)
	resBody, _ := json.Marshal(result)
	ctx.Output.Body(resBody)
}

func DetachVolume(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	reqBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Read volume request body failed:", err)
		resBody, _ := json.Marshal("Read volume request body failed!")
		ctx.Output.Body(resBody)
	}

	volumeRequest := &volumes.VolumeRequest{}
	if err = json.Unmarshal(reqBody, volumeRequest); err != nil {
		log.Println("Parse volume request body failed:", err)
		resBody, _ := json.Marshal("Parse volume request body failed!")
		ctx.Output.Body(resBody)
	}

	result := volumes.DetachVolume(volumeRequest)
	resBody, _ := json.Marshal(result)
	ctx.Output.Body(resBody)
}

func MountVolume(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	reqBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Read volume request body failed:", err)
		resBody, _ := json.Marshal("Read volume request body failed!")
		ctx.Output.Body(resBody)
	}

	volumeRequest := &volumes.VolumeRequest{}
	if err = json.Unmarshal(reqBody, volumeRequest); err != nil {
		log.Println("Parse volume request body failed:", err)
		resBody, _ := json.Marshal("Parse volume request body failed!")
		ctx.Output.Body(resBody)
	}

	result := volumes.MountVolume(volumeRequest)
	resBody, _ := json.Marshal(result)
	ctx.Output.Body(resBody)
}

func UnmountVolume(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	reqBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Read volume request body failed:", err)
		resBody, _ := json.Marshal("Read volume request body failed!")
		ctx.Output.Body(resBody)
	}

	volumeRequest := &volumes.VolumeRequest{}
	if err = json.Unmarshal(reqBody, volumeRequest); err != nil {
		log.Println("Parse volume request body failed:", err)
		resBody, _ := json.Marshal("Parse volume request body failed!")
		ctx.Output.Body(resBody)
	}

	result := volumes.UnmountVolume(volumeRequest)
	resBody, _ := json.Marshal(result)
	ctx.Output.Body(resBody)
}
