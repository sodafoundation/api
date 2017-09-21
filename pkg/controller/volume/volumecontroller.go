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
This module implements a entry into the OpenSDS volume controller service.

*/

package volume

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/opensds/opensds/pkg/dock/client"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	api "github.com/opensds/opensds/pkg/model"
	"golang.org/x/net/context"
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

func NewController(req *pb.DockRequest) Controller {
	return &controller{
		Client:  client.NewClient(),
		Request: req,
	}
}

type controller struct {
	client.Client

	Request *pb.DockRequest
}

func (c *controller) CreateVolume() (*api.VolumeSpec, error) {
	if err := c.Client.Update(c.Request.GetDockInfo()); err != nil {
		log.Println("[Error] When parsing dock info:", err)
		return nil, err
	}

	result, err := c.Client.CreateVolume(context.Background(), c.Request)
	if err != nil {
		log.Println("[Error] Create volume failed in volume controller:", err)
		return &api.VolumeSpec{}, err
	}
	defer c.Client.Close()

	var vol = &api.VolumeSpec{}
	if err = json.Unmarshal([]byte(result.GetMessage()), vol); err != nil {
		log.Println("[Error] Create volume failed in volume controller:", err)
		return &api.VolumeSpec{}, err
	}
	return vol, nil
}

func (c *controller) DeleteVolume() *api.Response {
	if err := c.Client.Update(c.Request.GetDockInfo()); err != nil {
		log.Println("[Error] When parsing dock info:", err)
		return nil
	}

	result, err := c.Client.DeleteVolume(context.Background(), c.Request)
	if err != nil {
		log.Println("[Error] Delete volume failed in volume controller:", err)
		return &api.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}
	defer c.Client.Close()

	return &api.Response{
		Status:  result.GetStatus(),
		Message: result.GetMessage(),
	}
}

func (c *controller) CreateVolumeAttachment() (*api.VolumeAttachmentSpec, error) {
	if err := c.Client.Update(c.Request.GetDockInfo()); err != nil {
		log.Println("[Error] When parsing dock info:", err)
		return nil, err
	}

	result, err := c.Client.CreateVolumeAttachment(context.Background(), c.Request)
	if err != nil {
		log.Println("[Error] Create volume failed in volume controller:", err)
		return &api.VolumeAttachmentSpec{}, err
	}
	defer c.Client.Close()

	var atc = &api.VolumeAttachmentSpec{}
	if err = json.Unmarshal([]byte(result.GetMessage()), atc); err != nil {
		log.Println("[Error] Create volume failed in volume controller:", err)
		return &api.VolumeAttachmentSpec{}, err
	}
	return atc, nil
}

func (c *controller) UpdateVolumeAttachment() (*api.VolumeAttachmentSpec, error) {
	if err := c.Client.Update(c.Request.GetDockInfo()); err != nil {
		log.Println("[Error] When parsing dock info:", err)
		return nil, err
	}

	result, err := c.Client.UpdateVolumeAttachment(context.Background(), c.Request)
	if err != nil {
		log.Println("[Error] Update volume attachment failed in volume controller:", err)
		return &api.VolumeAttachmentSpec{}, err
	}
	defer c.Client.Close()

	var atc = &api.VolumeAttachmentSpec{}
	if err = json.Unmarshal([]byte(result.GetMessage()), atc); err != nil {
		log.Println("[Error] Update volume attachment failed in volume controller:", err)
		return &api.VolumeAttachmentSpec{}, err
	}
	return atc, nil
}

func (c *controller) DeleteVolumeAttachment() *api.Response {
	if err := c.Client.Update(c.Request.GetDockInfo()); err != nil {
		log.Println("[Error] When parsing dock info:", err)
		return nil
	}

	result, err := c.Client.DeleteVolumeAttachment(context.Background(), c.Request)
	if err != nil {
		log.Println("[Error] Delete volume attachment failed in volume controller:", err)
		return &api.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}
	defer c.Client.Close()

	return &api.Response{
		Status:  result.GetStatus(),
		Message: result.GetMessage(),
	}
}

func (c *controller) CreateVolumeSnapshot() (*api.VolumeSnapshotSpec, error) {
	if err := c.Client.Update(c.Request.GetDockInfo()); err != nil {
		log.Println("[Error] When parsing dock info:", err)
		return nil, err
	}

	result, err := c.Client.CreateVolumeSnapshot(context.Background(), c.Request)
	if err != nil {
		log.Println("[Error] Create volume snapshot failed in volume controller:", err)
		return &api.VolumeSnapshotSpec{}, err
	}
	defer c.Client.Close()

	var snp = &api.VolumeSnapshotSpec{}
	if err = json.Unmarshal([]byte(result.GetMessage()), snp); err != nil {
		log.Println("[Error] Create volume snapshot failed in volume controller:", err)
		return &api.VolumeSnapshotSpec{}, err
	}
	return snp, nil
}

func (c *controller) DeleteVolumeSnapshot() *api.Response {
	if err := c.Client.Update(c.Request.GetDockInfo()); err != nil {
		log.Println("[Error] When parsing dock info:", err)
		return nil
	}

	result, err := c.Client.DeleteVolumeSnapshot(context.Background(), c.Request)
	if err != nil {
		log.Println("[Error] Delete volume snapshot failed in volume controller:", err)
		return &api.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}
	defer c.Client.Close()

	return &api.Response{
		Status:  result.GetStatus(),
		Message: result.GetMessage(),
	}
}
