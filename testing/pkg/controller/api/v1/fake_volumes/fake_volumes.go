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
	"errors"
	"log"

	"github.com/opensds/opensds/testing/pkg/controller/api"
	"github.com/opensds/opensds/testing/pkg/controller/api/grpcapi"
	pb "github.com/opensds/opensds/testing/pkg/grpc/fake_opensds"
)

type VolumeRequestDeliver interface {
	createVolume() *pb.Response

	getVolume() *pb.Response

	listVolumes() *pb.Response

	deleteVolume() *pb.Response

	attachVolume() *pb.Response

	detachVolume() *pb.Response

	mountVolume() *pb.Response

	unmountVolume() *pb.Response
}

// VolumeRequest is a structure for all properties of
// a volume request
type VolumeRequest struct {
	ResourceType string `json:"resourcetType,omitempty"`
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Size         int32  `json:"size"`
	AllowDetails bool   `json:"allowDetails"`

	ActionType string `json:"actionType,omitempty"`
	Host       string `json:"host,omitempty"`
	Device     string `json:"device,omitempty"`
	Attachment string `json:"attachment,omitempty"`
	MountDir   string `json:"mountDir,omitempty"`
	FsType     string `json:"fsType,omitempty"`
}

func (vr VolumeRequest) createVolume() *pb.Response {
	return grpcapi.CreateVolume(vr.ResourceType, vr.Name, vr.Size)
}

func (vr VolumeRequest) getVolume() *pb.Response {
	return grpcapi.GetVolume(vr.ResourceType, vr.Id)
}

func (vr VolumeRequest) listVolumes() *pb.Response {
	return grpcapi.ListVolumes(vr.ResourceType, vr.AllowDetails)
}

func (vr VolumeRequest) deleteVolume() *pb.Response {
	return grpcapi.DeleteVolume(vr.ResourceType, vr.Id)
}

func (vr VolumeRequest) attachVolume() *pb.Response {
	return grpcapi.AttachVolume(vr.ResourceType, vr.Id, vr.Host, vr.Device)
}

func (vr VolumeRequest) detachVolume() *pb.Response {
	return grpcapi.DetachVolume(vr.ResourceType, vr.Id, vr.Attachment)
}

func (vr VolumeRequest) mountVolume() *pb.Response {
	return grpcapi.MountVolume(vr.MountDir, vr.Device, vr.FsType)
}

func (vr VolumeRequest) unmountVolume() *pb.Response {
	return grpcapi.UnmountVolume(vr.MountDir)
}

func CreateVolume(vrd VolumeRequestDeliver) (api.VolumeResponse, error) {
	var nullResponse api.VolumeResponse

	result := vrd.createVolume()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("Create volume error:", err)
		return nullResponse, err
	}

	var volumeResponse api.VolumeResponse
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &volumeResponse); err != nil {
		return nullResponse, err
	}
	return volumeResponse, nil
}

func GetVolume(vrd VolumeRequestDeliver) (api.VolumeDetailResponse, error) {
	var nullResponse api.VolumeDetailResponse

	result := vrd.getVolume()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("Get volume error:", err)
		return nullResponse, err
	}

	var volumeDetailResponse api.VolumeDetailResponse
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &volumeDetailResponse); err != nil {
		return nullResponse, err
	}
	return volumeDetailResponse, nil
}

func ListVolumes(vrd VolumeRequestDeliver) ([]api.VolumeResponse, error) {
	var nullResponses []api.VolumeResponse

	result := vrd.listVolumes()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("List all volumes error:", err)
		return nullResponses, err
	}

	var volumesResponse []api.VolumeResponse
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &volumesResponse); err != nil {
		return nullResponses, err
	}
	return volumesResponse, nil
}

func DeleteVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	var defaultResponse api.DefaultResponse

	result := vrd.deleteVolume()
	if result.GetStatus() == "Failure" {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = result.GetError()
		log.Println("Delete volume error:", defaultResponse.Error)
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}

func AttachVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	var defaultResponse api.DefaultResponse

	result := vrd.attachVolume()
	if result.GetStatus() == "Failure" {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = result.GetError()
		log.Println("Attach volume error:", defaultResponse.Error)
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}

func DetachVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	var defaultResponse api.DefaultResponse

	result := vrd.detachVolume()
	if result.GetStatus() == "Failure" {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = result.GetError()
		log.Println("Detach volume error:", defaultResponse.Error)
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}

func MountVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	var defaultResponse api.DefaultResponse

	result := vrd.mountVolume()
	if result.GetStatus() == "Failure" {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = result.GetError()
		log.Println("Mount volume error:", defaultResponse.Error)
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}

func UnmountVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	var defaultResponse api.DefaultResponse

	result := vrd.unmountVolume()
	if result.GetStatus() == "Failure" {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = result.GetError()
		log.Println("Unmount volume error:", defaultResponse.Error)
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}
