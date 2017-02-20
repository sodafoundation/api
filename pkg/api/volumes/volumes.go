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
This module implements the entry into CRUD operation of volumes.

*/

package volumes

import (
	"encoding/json"
	"log"

	"github.com/opensds/opensds/pkg/api"
	"github.com/opensds/opensds/pkg/api/rpcapi"
)

type VolumeRequestDeliver interface {
	createVolume() (string, error)

	getVolume() (string, error)

	getAllVolumes() (string, error)

	updateVolume() (string, error)

	deleteVolume() (string, error)

	attachVolume() (string, error)

	detachVolume() (string, error)

	mountVolume() (string, error)

	unmountVolume() (string, error)
}

// VolumeRequest is a structure for all properties of
// a volume request
type VolumeRequest struct {
	ResourceType string `json:"resourcetType,omitempty"`
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Size         int    `json:"size"`
	AllowDetails bool   `json:"allowDetails"`

	ActionType string `json:"actionType"`
	Host       string `json:"host,omitempty"`
	Device     string `json:"device"`
	Attachment string `json:"attachment,omitempty"`
	MountDir   string `json:"mountDir"`
	FsType     string `json:"fsType"`
}

func (vr VolumeRequest) createVolume() (string, error) {
	return rpcapi.CreateVolume(vr.ResourceType, vr.Name, vr.Size)
}

func (vr VolumeRequest) getVolume() (string, error) {
	return rpcapi.GetVolume(vr.ResourceType, vr.Id)
}

func (vr VolumeRequest) getAllVolumes() (string, error) {
	return rpcapi.GetAllVolumes(vr.ResourceType, vr.AllowDetails)
}

func (vr VolumeRequest) updateVolume() (string, error) {
	return rpcapi.UpdateVolume(vr.ResourceType, vr.Id, vr.Name)
}

func (vr VolumeRequest) deleteVolume() (string, error) {
	return rpcapi.DeleteVolume(vr.ResourceType, vr.Id)
}

func (vr VolumeRequest) attachVolume() (string, error) {
	return rpcapi.AttachVolume(vr.ResourceType, vr.Id, vr.Host, vr.Device)
}

func (vr VolumeRequest) detachVolume() (string, error) {
	return rpcapi.DetachVolume(vr.ResourceType, vr.Id, vr.Attachment)
}

func (vr VolumeRequest) mountVolume() (string, error) {
	return rpcapi.MountVolume(vr.MountDir, vr.Device, vr.Id, vr.FsType)
}

func (vr VolumeRequest) unmountVolume() (string, error) {
	return rpcapi.UnmountVolume(vr.MountDir)
}

func Create(vrd VolumeRequestDeliver) (api.VolumeResponse, error) {
	var nullResponse api.VolumeResponse

	result, err := vrd.createVolume()
	if err != nil {
		log.Println("Create volume error: ", err)
		return nullResponse, err
	}

	var volumeResponse api.VolumeResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &volumeResponse); err != nil {
		return nullResponse, err
	}
	return volumeResponse, nil
}

func Show(vrd VolumeRequestDeliver) (api.VolumeDetailResponse, error) {
	var nullResponse api.VolumeDetailResponse

	result, err := vrd.getVolume()
	if err != nil {
		log.Println("Show volume error: ", err)
		return nullResponse, err
	}

	var volumeDetailResponse api.VolumeDetailResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &volumeDetailResponse); err != nil {
		return nullResponse, err
	}
	return volumeDetailResponse, nil
}

func List(vrd VolumeRequestDeliver) ([]api.VolumeResponse, error) {
	var nullResponses []api.VolumeResponse

	result, err := vrd.getAllVolumes()
	if err != nil {
		log.Println("List volumes error: ", err)
		return nullResponses, err
	}

	var volumesResponse []api.VolumeResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &volumesResponse); err != nil {
		return nullResponses, err
	}
	return volumesResponse, nil
}

func Update(vrd VolumeRequestDeliver) (api.VolumeResponse, error) {
	var nullResponse api.VolumeResponse

	result, err := vrd.updateVolume()
	if err != nil {
		log.Println("Update volume error: ", err)
		return nullResponse, err
	}

	var volumeResponse api.VolumeResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &volumeResponse); err != nil {
		return nullResponse, err
	}
	return volumeResponse, nil
}

func Delete(vrd VolumeRequestDeliver) (string, error) {
	result, err := vrd.deleteVolume()
	if err != nil {
		log.Println("Delete volume error: ", err)
		return "", err
	}
	return result, nil
}

func Attach(vrd VolumeRequestDeliver) (string, error) {
	result, err := vrd.attachVolume()
	if err != nil {
		log.Println("Attach volume error: ", err)
		return "", err
	}
	return result, nil
}

func Detach(vrd VolumeRequestDeliver) (string, error) {
	result, err := vrd.detachVolume()
	if err != nil {
		log.Println("Detach volume error: ", err)
		return "", err
	}
	return result, nil
}

func Mount(vrd VolumeRequestDeliver) (string, error) {
	result, err := vrd.mountVolume()
	if err != nil {
		log.Println("Mount volume error: ", err)
		return "", err
	}
	return result, nil
}

func Unmount(vrd VolumeRequestDeliver) (string, error) {
	result, err := vrd.unmountVolume()
	if err != nil {
		log.Println("Unmount volume error: ", err)
		return "", err
	}
	return result, nil
}
