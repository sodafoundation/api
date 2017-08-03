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

package volume

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/opensds/opensds/pkg/grpc/dock/client"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
	api "github.com/opensds/opensds/pkg/model"
)

type Controller interface {
	CreateVolume() (*api.VolumeSpec, error)

	DeleteVolume() *api.Response

	CreateVolumeAttachment() (*api.VolumeAttachmentSpec, error)

	UpdateVolumeAttachment() (*api.VolumeAttachmentSpec, error)

	DeleteVolumeAttachment() *api.Response

	CreateVolumeSnapshot() (*api.VolumeSnapshotSpec, error)

	DeleteVolumeSnapshot() *api.Response
}

func NewController(request *pb.DockRequest) Controller {
	return &controller{
		Request: request,
	}
}

type controller struct {
	Request *pb.DockRequest
}

func (c *controller) CreateVolume() (*api.VolumeSpec, error) {
	result, err := client.CreateVolume(c.Request)
	if err != nil {
		log.Println("[Error] Create volume failed in volume controller:", err)
		return &api.VolumeSpec{}, err
	}

	var vol = &api.VolumeSpec{}
	if err = json.Unmarshal([]byte(result.GetMessage()), vol); err != nil {
		log.Println("[Error] Create volume failed in volume controller:", err)
		return &api.VolumeSpec{}, err
	}
	return vol, nil
}

func (c *controller) DeleteVolume() *api.Response {
	result, err := client.DeleteVolume(c.Request)
	if err != nil {
		log.Println("[Error] Delete volume failed in volume controller:", err)
		return &api.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}
	return &api.Response{
		Status:  result.GetStatus(),
		Message: result.GetMessage(),
	}
}

func (c *controller) CreateVolumeAttachment() (*api.VolumeAttachmentSpec, error) {
	result, err := client.CreateVolumeAttachment(c.Request)
	if err != nil {
		log.Println("[Error] Create volume failed in volume controller:", err)
		return &api.VolumeAttachmentSpec{}, err
	}

	var atc = &api.VolumeAttachmentSpec{}
	if err = json.Unmarshal([]byte(result.GetMessage()), atc); err != nil {
		log.Println("[Error] Create volume failed in volume controller:", err)
		return &api.VolumeAttachmentSpec{}, err
	}
	return atc, nil
}

func (c *controller) UpdateVolumeAttachment() (*api.VolumeAttachmentSpec, error) {
	result, err := client.UpdateVolumeAttachment(c.Request)
	if err != nil {
		log.Println("[Error] Update volume attachment failed in volume controller:", err)
		return &api.VolumeAttachmentSpec{}, err
	}

	var atc = &api.VolumeAttachmentSpec{}
	if err = json.Unmarshal([]byte(result.GetMessage()), atc); err != nil {
		log.Println("[Error] Update volume attachment failed in volume controller:", err)
		return &api.VolumeAttachmentSpec{}, err
	}
	return atc, nil
}

func (c *controller) DeleteVolumeAttachment() *api.Response {
	result, err := client.DeleteVolumeAttachment(c.Request)
	if err != nil {
		log.Println("[Error] Delete volume attachment failed in volume controller:", err)
		return &api.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}
	return &api.Response{
		Status:  result.GetStatus(),
		Message: result.GetMessage(),
	}
}

func (c *controller) CreateVolumeSnapshot() (*api.VolumeSnapshotSpec, error) {
	result, err := client.CreateVolumeSnapshot(c.Request)
	if err != nil {
		log.Println("[Error] Create volume snapshot failed in volume controller:", err)
		return &api.VolumeSnapshotSpec{}, err
	}

	var snp = &api.VolumeSnapshotSpec{}
	if err = json.Unmarshal([]byte(result.GetMessage()), snp); err != nil {
		log.Println("[Error] Create volume snapshot failed in volume controller:", err)
		return &api.VolumeSnapshotSpec{}, err
	}
	return snp, nil
}

func (c *controller) DeleteVolumeSnapshot() *api.Response {
	result, err := client.DeleteVolumeSnapshot(c.Request)
	if err != nil {
		log.Println("[Error] Delete volume snapshot failed in volume controller:", err)
		return &api.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}
	return &api.Response{
		Status:  result.GetStatus(),
		Message: result.GetMessage(),
	}
}
