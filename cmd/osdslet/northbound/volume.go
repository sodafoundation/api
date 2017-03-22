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
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	"github.com/opensds/opensds/pkg/api"
	"github.com/opensds/opensds/pkg/api/v1/volumes"

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
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	resourceType := this.Ctx.Input.Param(":resource")
	id := this.Ctx.Input.Param(":id")

	volumeRequest := volumes.VolumeRequest{
		ResourceType: resourceType,
		Id:           id,
	}
	result, err := volumes.GetVolume(volumeRequest)
	if err != nil {
		log.Println(err)
		rbody, _ := json.Marshal("Show volume failed!")
		this.Ctx.Output.Body(rbody)
	} else {
		if reflect.DeepEqual(result, falseVolumeResponse) {
			log.Println("Show volume failed!")
			rbody, _ := json.Marshal("Show volume failed!")
			this.Ctx.Output.Body(rbody)
		} else {
			rbody, _ := json.Marshal(result)
			this.Ctx.Output.Body(rbody)
		}
	}
}

func (this *VolumeController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	rbody, _ := json.Marshal("Not supported!")
	this.Ctx.Output.Body(rbody)
}

func (this *VolumeController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	resourceType := this.Ctx.Input.Param(":resource")
	id := this.Ctx.Input.Param(":id")

	volumeRequest := volumes.VolumeRequest{
		ResourceType: resourceType,
		Id:           id,
	}
	result := volumes.DeleteVolume(volumeRequest)
	rbody, _ := json.Marshal(result)
	this.Ctx.Output.Body(rbody)
}

func PostVolume(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	resourceType := ctx.Input.Param(":resource")
	rbody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Read volume request body failed:", err)
		rbody, _ := json.Marshal("Read volume request body failed!")
		ctx.Output.Body(rbody)
	}

	var volumeRequest volumes.VolumeRequest
	volumeRequest.ResourceType = resourceType
	if err = json.Unmarshal(rbody, &volumeRequest); err != nil {
		log.Println("Parse volume request body failed:", err)
		rbody, _ := json.Marshal("Parse volume request body failed!")
		ctx.Output.Body(rbody)
	}

	result, err := volumes.CreateVolume(volumeRequest)
	if err != nil {
		log.Println(err)
		rbody, _ := json.Marshal("Create volume failed!")
		ctx.Output.Body(rbody)
	} else {
		if reflect.DeepEqual(result, falseVolumeResponse) {
			log.Println("Create volume failed!")
			rbody, _ := json.Marshal("Create volume failed!")
			ctx.Output.Body(rbody)
		} else {
			rbody, _ := json.Marshal(result)
			ctx.Output.Body(rbody)
		}
	}
}

func GetAllVolumes(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	resourceType := ctx.Input.Param(":resource")

	volumeRequest := volumes.VolumeRequest{
		ResourceType: resourceType,
		AllowDetails: false,
	}
	result, err := volumes.ListVolumes(volumeRequest)
	if err != nil {
		log.Println(err)
		rbody, _ := json.Marshal("List volumes failed!")
		ctx.Output.Body(rbody)
	} else {
		if reflect.DeepEqual(result, falseAllVolumesResponse) {
			log.Println("List volumes failed!")
			rbody, _ := json.Marshal("List volumes failed!")
			ctx.Output.Body(rbody)
		} else {
			rbody, _ := json.Marshal(result)
			ctx.Output.Body(rbody)
		}
	}
}

func PostVolumeAction(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	resourceType := ctx.Input.Param(":resource")
	id := ctx.Input.Param(":id")
	rbody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Read volume request body failed:", err)
		rbody, _ := json.Marshal("Read volume request body failed!")
		ctx.Output.Body(rbody)
	}

	var volumeRequest volumes.VolumeRequest
	volumeRequest.ResourceType = resourceType
	volumeRequest.Id = id
	if err = json.Unmarshal(rbody, &volumeRequest); err != nil {
		log.Println("Parse volume request body failed:", err)
		rbody, _ := json.Marshal("Parse volume request body failed!")
		ctx.Output.Body(rbody)
	}

	switch volumeRequest.ActionType {
	case "attach":
		if volumeRequest.Host == "" {
			volumeRequest.Host, _ = os.Hostname()
		}
		if volumeRequest.Device == "" {
			volumeRequest.Device = "/mnt"
		}

		result := volumes.AttachVolume(volumeRequest)
		rbody, _ := json.Marshal(result)
		ctx.Output.Body(rbody)
	case "detach":
		result := volumes.DetachVolume(volumeRequest)
		rbody, _ := json.Marshal(result)
		ctx.Output.Body(rbody)
	case "mount":
		result := volumes.MountVolume(volumeRequest)
		rbody, _ := json.Marshal(result)
		ctx.Output.Body(rbody)
	case "unmount":
		result := volumes.UnmountVolume(volumeRequest)
		rbody, _ := json.Marshal(result)
		ctx.Output.Body(rbody)
	default:
		err := errors.New("The type of volume action is not correct!")
		log.Println(err)
		rbody, _ := json.Marshal("The type of volume action is not correct!")
		ctx.Output.Body(rbody)
	}
}
