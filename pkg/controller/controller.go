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

package controller

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/opensds/opensds/pkg/controller/policy"
	"github.com/opensds/opensds/pkg/controller/volume"
	"github.com/opensds/opensds/pkg/db"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
	api "github.com/opensds/opensds/pkg/model"
)

const (
	CREATE_LIFECIRCLE_FLAG = 1
	GET_LIFECIRCLE_FLAG    = 2
	LIST_LIFECIRCLE_FLAG   = 3
	DELETE_LIFECIRCLE_FLAG = 4
)

func NewControllerWithVolumeConfig(
	vol *api.VolumeSpec,
	atc *api.VolumeAttachmentSpec,
	snp *api.VolumeSnapshotSpec,
) (*Controller, error) {
	c := &Controller{
		request: &pb.DockRequest{},
	}

	// If volume input is not null, the controller will add policy orchestration
	// to manage volume resource.
	if vol != nil {
		prf, err := SearchProfile(db.C, vol.GetProfileId())
		if err != nil {
			log.Println("[Error] when search profiles in db:", err)
			return nil, err
		}

		c.profile = prf
		c.request.VolumeId = vol.GetId()
		c.request.VolumeName = vol.GetName()
		c.request.VolumeDescription = vol.GetDescription()
		c.request.VolumeSize = vol.GetSize()
		c.request.ProfileId = prf.GetId()
	}
	if atc != nil {
		c.request.AttachmentId = atc.GetId()
		c.request.AttachmentName = atc.GetName()
		c.request.AttachmentDescription = atc.GetDescription()
		c.request.VolumeId = atc.GetVolumeId()
	}
	if snp != nil {
		c.request.SnapshotId = snp.GetId()
		c.request.SnapshotName = snp.GetName()
		c.request.SnapshotDescription = snp.GetDescription()
		c.request.VolumeId = snp.GetVolumeId()
	}

	return c, nil
}

type Controller struct {
	volumeController volume.Controller
	policyController policy.Controller
	profile          *api.ProfileSpec
	request          *pb.DockRequest
}

func (c *Controller) CreateVolume() (*api.VolumeSpec, error) {
	// Initialize volume and policy controller.
	c.policyController = policy.NewController(c.profile)
	c.volumeController = volume.NewController(c.request)
	c.policyController.Setup(CREATE_LIFECIRCLE_FLAG)

	polInfo, err := SearchSupportedPool(db.C, c.policyController.StorageTag().GetSyncTag())
	if err != nil {
		log.Println("[Error] When search supported pool resource:", err)
		return &api.VolumeSpec{}, err
	}
	c.request.PoolId = polInfo.GetId()

	dckInfo, err := SearchDockByPool(db.C, polInfo)
	if err != nil {
		log.Println("[Error] When search supported dock resource:", err)
		return &api.VolumeSpec{}, err
	}
	dckBody, _ := json.Marshal(dckInfo)
	c.request.DockInfo = string(dckBody)

	result, err := c.volumeController.CreateVolume()
	if err != nil {
		return &api.VolumeSpec{}, err
	}

	var errChan = make(chan error, 1)
	volBody, _ := json.Marshal(result)
	go c.policyController.ExecuteAsyncPolicy(c.request, string(volBody), errChan)

	return result, nil
}

func (c *Controller) DeleteVolume() *api.Response {
	c.policyController = policy.NewController(c.profile)
	c.volumeController = volume.NewController(c.request)
	c.policyController.Setup(DELETE_LIFECIRCLE_FLAG)

	dckInfo, err := SearchDockByVolume(db.C, c.request.GetVolumeId())
	if err != nil {
		log.Println("[Error] When search supported dock resource:", err)
		return &api.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}
	dckBody, _ := json.Marshal(dckInfo)
	c.request.DockInfo = string(dckBody)

	var errChan = make(chan error, 1)
	go c.policyController.ExecuteAsyncPolicy(c.request, "", errChan)

	if err := <-errChan; err != nil {
		log.Println("[Error] When execute async policy:", err)
		return &api.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}

	return c.volumeController.DeleteVolume()
}

func (c *Controller) CreateVolumeAttachment() (*api.VolumeAttachmentSpec, error) {
	c.volumeController = volume.NewController(c.request)

	dckInfo, err := SearchDockByVolume(db.C, c.request.GetVolumeId())
	if err != nil {
		log.Println("[Error] When search supported dock resource:", err)
		return &api.VolumeAttachmentSpec{}, err
	}
	dckBody, _ := json.Marshal(dckInfo)
	c.request.DockInfo = string(dckBody)

	return c.volumeController.CreateVolumeAttachment()
}

func (c *Controller) UpdateVolumeAttachment() (*api.VolumeAttachmentSpec, error) {
	c.volumeController = volume.NewController(c.request)

	dckInfo, err := SearchDockByVolume(db.C, c.request.GetVolumeId())
	if err != nil {
		log.Println("[Error] When search supported dock resource:", err)
		return &api.VolumeAttachmentSpec{}, err
	}
	dckBody, _ := json.Marshal(dckInfo)
	c.request.DockInfo = string(dckBody)

	return c.volumeController.UpdateVolumeAttachment()
}

func (c *Controller) DeleteVolumeAttachment() *api.Response {
	c.volumeController = volume.NewController(c.request)

	dckInfo, err := SearchDockByVolume(db.C, c.request.GetVolumeId())
	if err != nil {
		log.Println("[Error] When search supported dock resource:", err)
		return &api.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}
	dckBody, _ := json.Marshal(dckInfo)
	c.request.DockInfo = string(dckBody)

	return c.volumeController.DeleteVolumeAttachment()
}

func (c *Controller) CreateVolumeSnapshot() (*api.VolumeSnapshotSpec, error) {
	c.volumeController = volume.NewController(c.request)

	dckInfo, err := SearchDockByVolume(db.C, c.request.GetVolumeId())
	if err != nil {
		log.Println("[Error] When search supported dock resource:", err)
		return &api.VolumeSnapshotSpec{}, err
	}
	dckBody, _ := json.Marshal(dckInfo)
	c.request.DockInfo = string(dckBody)

	return c.volumeController.CreateVolumeSnapshot()
}

func (c *Controller) DeleteVolumeSnapshot() *api.Response {
	c.volumeController = volume.NewController(c.request)

	dckInfo, err := SearchDockByVolume(db.C, c.request.GetVolumeId())
	if err != nil {
		log.Println("[Error] When search supported dock resource:", err)
		return &api.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}
	dckBody, _ := json.Marshal(dckInfo)
	c.request.DockInfo = string(dckBody)

	return c.volumeController.DeleteVolumeSnapshot()
}
