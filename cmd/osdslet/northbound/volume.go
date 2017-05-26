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
	volumes "github.com/opensds/opensds/pkg/apiserver"

	"github.com/astaxie/beego"
)

type VolumeController struct {
	beego.Controller
}

func (this *VolumeController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	reqBody, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		log.Println("Read volume request body failed:", err)
		resBody, _ := json.Marshal("Read volume request body failed!")
		this.Ctx.Output.Body(resBody)
	}

	volumeRequest := &volumes.VolumeRequest{}
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
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *VolumeController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	volumeRequest := &volumes.VolumeRequest{}
	result, err := volumes.ListVolumes(volumeRequest)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("List volumes failed!")
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *VolumeController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *VolumeController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}

type SpecifiedVolumeController struct {
	beego.Controller
}

func (this *SpecifiedVolumeController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *SpecifiedVolumeController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	volumeRequest := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			Id: this.Ctx.Input.Param(":id"),
		},
	}
	result, err := volumes.GetVolume(volumeRequest)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Get volume failed!")
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *SpecifiedVolumeController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *SpecifiedVolumeController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	volId := this.Ctx.Input.Param(":id")
	reqBody, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		log.Println("Read volume request body failed:", err)
		resBody, _ := json.Marshal("Read volume request body failed!")
		this.Ctx.Output.Body(resBody)
	}

	volumeRequest := &volumes.VolumeRequest{}
	if err = json.Unmarshal(reqBody, volumeRequest); err != nil {
		log.Println("Parse volume request body failed:", err)
		resBody, _ := json.Marshal("Parse volume request body failed!")
		this.Ctx.Output.Body(resBody)
	}
	volumeRequest.Schema.Id = volId

	result := volumes.DeleteVolume(volumeRequest)
	resBody, _ := json.Marshal(result)
	this.Ctx.Output.Body(resBody)
}

type VolumeAttachmentController struct {
	beego.Controller
}

func (this *VolumeAttachmentController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	reqBody, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		log.Println("Read volume request body failed:", err)
		resBody, _ := json.Marshal("Read volume request body failed!")
		this.Ctx.Output.Body(resBody)
	}

	volumeRequest := &volumes.VolumeRequest{}
	if err = json.Unmarshal(reqBody, volumeRequest); err != nil {
		log.Println("Parse volume request body failed:", err)
		resBody, _ := json.Marshal("Parse volume request body failed!")
		this.Ctx.Output.Body(resBody)
	}

	result, err := volumes.CreateVolumeAttachment(volumeRequest)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Create volume attachment failed!")
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *VolumeAttachmentController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	volId := this.GetString("volId")
	vr := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			Id: volId,
		},
	}
	result, err := volumes.ListVolumeAttachments(vr)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("List volume attachments failed!")
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *VolumeAttachmentController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *VolumeAttachmentController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}

type SpecifiedVolumeAttachmentController struct {
	beego.Controller
}

func (this *SpecifiedVolumeAttachmentController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *SpecifiedVolumeAttachmentController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	volId := this.GetString("volId")
	attachmentId := this.GetString("id")
	vr := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			Id:           volId,
			AttachmentId: attachmentId,
		},
	}
	result, err := volumes.GetVolumeAttachment(vr)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Get volume attachment failed!")
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *SpecifiedVolumeAttachmentController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	volId := this.GetString("volId")
	attachmentId := this.GetString("id")

	reqBody, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		log.Println("Read volume request body failed:", err)
		resBody, _ := json.Marshal("Read volume request body failed!")
		this.Ctx.Output.Body(resBody)
	}

	vr := &volumes.VolumeRequest{}
	if err = json.Unmarshal(reqBody, vr); err != nil {
		log.Println("Parse volume request body failed:", err)
		resBody, _ := json.Marshal("Parse volume request body failed!")
		this.Ctx.Output.Body(resBody)
	}
	vr.Schema.Id, vr.Schema.AttachmentId = volId, attachmentId

	result, err := volumes.GetVolumeAttachment(vr)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Update volume attachment failed!")
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *SpecifiedVolumeAttachmentController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	volId := this.GetString("volId")
	attachmentId := this.GetString("id")
	vr := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			Id:           volId,
			AttachmentId: attachmentId,
		},
	}

	result := volumes.DeleteVolumeAttachment(vr)
	resBody, _ := json.Marshal(result)
	this.Ctx.Output.Body(resBody)
}

type VolumeSnapshotController struct {
	beego.Controller
}

func (this *VolumeSnapshotController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	reqBody, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		log.Println("Read volume request body failed:", err)
		resBody, _ := json.Marshal("Read volume request body failed!")
		this.Ctx.Output.Body(resBody)
	}

	volumeRequest := &volumes.VolumeRequest{}
	if err = json.Unmarshal(reqBody, volumeRequest); err != nil {
		log.Println("Parse volume request body failed:", err)
		resBody, _ := json.Marshal("Parse volume request body failed!")
		this.Ctx.Output.Body(resBody)
	}

	result, err := volumes.CreateVolumeSnapshot(volumeRequest)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Create volume attachment failed!")
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *VolumeSnapshotController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	vr := &volumes.VolumeRequest{}
	result, err := volumes.ListVolumeSnapshots(vr)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("List volume snapshots failed!")
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *VolumeSnapshotController) Put() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *VolumeSnapshotController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}

type SpecifiedVolumeSnapshotController struct {
	beego.Controller
}

func (this *SpecifiedVolumeSnapshotController) Post() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Body([]byte("Not supported!"))
}

func (this *SpecifiedVolumeSnapshotController) Get() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	snapshotId := this.GetString("id")
	vr := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			SnapshotId: snapshotId,
		},
	}
	result, err := volumes.GetVolumeSnapshot(vr)
	if err != nil {
		log.Println(err)
		resBody, _ := json.Marshal("Get volume snapshot failed!")
		this.Ctx.Output.Body(resBody)
	} else {
		resBody, _ := json.Marshal(result)
		this.Ctx.Output.Body(resBody)
	}
}

func (this *SpecifiedVolumeSnapshotController) Delete() {
	this.Ctx.Output.Header("Content-Type", "application/json")
	this.Ctx.Output.ContentType("application/json")

	volId := this.GetString("volId")
	snapshotId := this.GetString("id")
	vr := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			Id:         volId,
			SnapshotId: snapshotId,
		},
	}

	result := volumes.DeleteVolumeSnapshot(vr)
	resBody, _ := json.Marshal(result)
	this.Ctx.Output.Body(resBody)
}
