// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
This module implements a entry into the OpenSDS northbound service.

*/

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/cindercompatibleapi/converter"

	"github.com/opensds/opensds/pkg/model"
)

// VolumePortal ...
type VolumePortal struct {
	beego.Controller
}

// ListVolumesDetails ...
func (portal *VolumePortal) ListVolumesDetails() {
	volumes, err := client.ListVolumes()
	if err != nil {
		reason := fmt.Sprintf("List accessible volumes with details failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.ListVolumesDetailsResp(volumes)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("List accessible volumes with details, marshal result failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// CreateVolume ...
func (portal *VolumePortal) CreateVolume() {
	var cinderReq = converter.CreateVolumeReqSpec{}

	if err := json.NewDecoder(portal.Ctx.Request.Body).Decode(&cinderReq); err != nil {
		reason := fmt.Sprintf("Create a volume, parse request body failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	volume, err := converter.CreateVolumeReq(&cinderReq)
	if err != nil {
		reason := fmt.Sprintf("Create a volume failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	volume, err = client.CreateVolume(volume)
	if err != nil {
		reason := fmt.Sprintf("Create a volume failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.CreateVolumeResp(volume)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Create a volume, marshal result failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusAccepted)
	portal.Ctx.Output.Body(body)
	return
}

// ListVolumes ...
func (portal *VolumePortal) ListVolumes() {
	volumes, err := client.ListVolumes()
	if err != nil {
		reason := fmt.Sprintf("List accessible volumes failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.ListVolumesResp(volumes)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("List accessible volumes, marshal result failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// GetVolume ...
func (portal *VolumePortal) GetVolume() {
	id := portal.Ctx.Input.Param(":volumeId")
	volume, err := client.GetVolume(id)

	if err != nil {
		reason := fmt.Sprintf("Show a volume's details failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.ShowVolumeResp(volume)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Show a volume's details, marshal result failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// UpdateVolume ...
func (portal *VolumePortal) UpdateVolume() {
	id := portal.Ctx.Input.Param(":volumeId")
	var cinderReq = converter.UpdateVolumeReqSpec{}

	if err := json.NewDecoder(portal.Ctx.Request.Body).Decode(&cinderReq); err != nil {
		reason := fmt.Sprintf("Update a volume, parse request body failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	volume, err := converter.UpdateVolumeReq(&cinderReq)
	if err != nil {
		reason := fmt.Sprintf("Update a volume failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
		portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	volume, err = client.UpdateVolume(id, volume)

	if err != nil {
		reason := fmt.Sprintf("Update a volume failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	result := converter.UpdateVolumeResp(volume)
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Update a volume, marshal result failed: %s", err.Error())
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusOK)
	portal.Ctx.Output.Body(body)
	return
}

// DeleteVolume ...
func (portal *VolumePortal) DeleteVolume() {
	id := portal.Ctx.Input.Param(":volumeId")
	volume := model.VolumeSpec{}

	err := client.DeleteVolume(id, &volume)

	if err != nil {
		reason := fmt.Sprintf("Delete a volume failed: %v", err)
		portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	portal.Ctx.Output.SetStatus(http.StatusAccepted)
	return
}

// VolumeAction ...
func (portal *VolumePortal) VolumeAction() {
	id := portal.Ctx.Input.Param(":volumeId")
	byts, err := ioutil.ReadAll(portal.Ctx.Request.Body)
	if err != nil {
		reason := fmt.Sprintf("Volume actions failed: request body is incorrect")
		log.Error(reason)
		portal.Ctx.Output.SetStatus(http.StatusNotFound)
		portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		return
	}

	rawBodyText := string(byts)

	// No actual operation is currently done
	if `{"os-reserve": null}` == rawBodyText {
		portal.Ctx.Output.SetStatus(http.StatusAccepted)
		return
	}

	if strings.HasPrefix(rawBodyText, `{"os-initialize_connection"`) {
		var cinderReq = converter.InitializeConnectionReqSpec{}
		err = json.Unmarshal([]byte(rawBodyText), &cinderReq)

		if err != nil {
			reason := fmt.Sprintf("Initialize connection, parse request body failed: %s", err.Error())
			portal.Ctx.Output.SetStatus(model.ErrorBadRequest)
			portal.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
			log.Error(reason)
			return
		}

		attachment := converter.InitializeConnectionReq(&cinderReq, id)
		attachment, err := client.CreateVolumeAttachment(attachment)

		if err != nil {
			reason := fmt.Sprintf("Initialize connection failed: %s", err.Error())
			portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
			portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
			log.Error(reason)
			return
		}

		isAvailable := false
		sum := 0

		for {
			sum++
			time.Sleep(1e9)
			attachment, _ = client.GetVolumeAttachment(attachment.Id)
			if ("available" == attachment.Status) && ("" != attachment.ConnectionInfo.DriverVolumeType) &&
				//(nil != attachment.ConnectionInfo.ConnectionData["authPassword"]) &&
				(nil != attachment.ConnectionInfo.ConnectionData["targetDiscovered"]) &&
				//(nil != attachment.ConnectionInfo.ConnectionData["encrypted"]) &&
				(nil != attachment.ConnectionInfo.ConnectionData["targetIQN"]) &&
				(nil != attachment.ConnectionInfo.ConnectionData["targetPortal"]) &&
				//(nil != attachment.ConnectionInfo.ConnectionData["volumeId"]) &&
				(nil != attachment.ConnectionInfo.ConnectionData["targetLun"]) {
				//(nil != attachment.ConnectionInfo.ConnectionData["accessMode"]) &&
				//(nil != attachment.ConnectionInfo.ConnectionData["authUserName"]) &&
				//(nil != attachment.ConnectionInfo.ConnectionData["authMethod"]) {
				isAvailable = true
				break
			}

			//The maximum waiting time is 10 seconds
			if 10 == sum {
				break
			}
		}

		if false == isAvailable {
			reason := fmt.Sprintf("Initialize connection, attachment is not available or connectionInfo is incorrect")
			attachmentByts, _ := json.Marshal(attachment)
			fmt.Println(string(attachmentByts))
			portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
			portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
			log.Error(reason)
			return
		}

		result := converter.InitializeConnectionResp(attachment)
		body, err := json.Marshal(result)
		if err != nil {
			reason := fmt.Sprintf("Initialize connection, marshal result failed: %s", err.Error())
			portal.Ctx.Output.SetStatus(model.ErrorInternalServer)
			portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
			log.Error(reason)
			return
		}

		portal.Ctx.Output.SetStatus(http.StatusOK)
		portal.Ctx.Output.Body(body)
		return
	}

	if strings.HasPrefix(rawBodyText, `{"os-terminate_connection":`) {
		portal.Ctx.Output.SetStatus(http.StatusAccepted)
		return
	}

	if `{"os-unreserve": null}` == rawBodyText {
		portal.Ctx.Output.SetStatus(http.StatusAccepted)
		return
	}

	if strings.HasPrefix(rawBodyText, `{"os-attach":`) {
		portal.Ctx.Output.SetStatus(http.StatusAccepted)
		return
	}

	if strings.HasPrefix(rawBodyText, `{"os-detach":`) {
		portal.Ctx.Output.SetStatus(http.StatusAccepted)
		return
	}

	if strings.HasPrefix(rawBodyText, `{"os-begin_detaching":`) {
		portal.Ctx.Output.SetStatus(http.StatusAccepted)
		return
	}

	reason := fmt.Sprintf("Volume actions failed: the body of the request is wrong or not currently supported")
	log.Error("Volume actions failed: " + rawBodyText + " is incorrect")
	portal.Ctx.Output.SetStatus(http.StatusNotFound)
	portal.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
	return
}
