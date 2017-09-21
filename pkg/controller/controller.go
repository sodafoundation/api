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
	pb "github.com/opensds/opensds/pkg/dock/proto"
	api "github.com/opensds/opensds/pkg/model"
)

const (
	CREATE_LIFECIRCLE_FLAG = iota + 1
	GET_LIFECIRCLE_FLAG
	LIST_LIFECIRCLE_FLAG
	DELETE_LIFECIRCLE_FLAG
)

func NewControllerWithVolumeConfig(
	vol *api.VolumeSpec,
	atc *api.VolumeAttachmentSpec,
	snp *api.VolumeSnapshotSpec,
) (*Controller, error) {
	c := &Controller{
		request: &pb.DockRequest{},
	}

	c.searcher = NewDbSearcher()

	// If volume input is not null, the controller will add policy orchestration
	// to manage volume resource.
	if vol != nil {
		prf, err := c.searcher.SearchProfile(vol.GetProfileId())
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

		// Initialize policy controller when profile is specified.
		c.policyController = policy.NewController(c.profile)
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

	// Initialize volume controller.
	c.volumeController = volume.NewController(c.request)

	return c, nil
}

type Controller struct {
	searcher         Searcher
	volumeController volume.Controller
	policyController policy.Controller

	profile *api.ProfileSpec
	request *pb.DockRequest
}

func (c *Controller) CreateVolume() (*api.VolumeSpec, error) {
	// Select the storage tag according to the lifecycle flag.
	c.policyController.Setup(CREATE_LIFECIRCLE_FLAG)

	polInfo, err := c.searcher.SearchSupportedPool(c.policyController.StorageTag().GetSyncTag())
	if err != nil {
		log.Println("[Error] When search supported pool resource:", err)
		return &api.VolumeSpec{}, err
	}
	c.request.PoolId = polInfo.GetId()

	dckInfo, err := c.searcher.SearchDockByPool(polInfo)
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
	// Select the storage tag according to the lifecycle flag.
	c.policyController.Setup(DELETE_LIFECIRCLE_FLAG)

	dckInfo, err := c.searcher.SearchDockByVolume(c.request.GetVolumeId())
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

	dckInfo, err := c.searcher.SearchDockByVolume(c.request.GetVolumeId())
	if err != nil {
		log.Println("[Error] When search supported dock resource:", err)
		return &api.VolumeAttachmentSpec{}, err
	}
	dckBody, _ := json.Marshal(dckInfo)
	c.request.DockInfo = string(dckBody)

	return c.volumeController.CreateVolumeAttachment()
}

func (c *Controller) UpdateVolumeAttachment() (*api.VolumeAttachmentSpec, error) {

	dckInfo, err := c.searcher.SearchDockByVolume(c.request.GetVolumeId())
	if err != nil {
		log.Println("[Error] When search supported dock resource:", err)
		return &api.VolumeAttachmentSpec{}, err
	}
	dckBody, _ := json.Marshal(dckInfo)
	c.request.DockInfo = string(dckBody)

	return c.volumeController.UpdateVolumeAttachment()
}

func (c *Controller) DeleteVolumeAttachment() *api.Response {

	dckInfo, err := c.searcher.SearchDockByVolume(c.request.GetVolumeId())
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

	dckInfo, err := c.searcher.SearchDockByVolume(c.request.GetVolumeId())
	if err != nil {
		log.Println("[Error] When search supported dock resource:", err)
		return &api.VolumeSnapshotSpec{}, err
	}
	dckBody, _ := json.Marshal(dckInfo)
	c.request.DockInfo = string(dckBody)

	return c.volumeController.CreateVolumeSnapshot()
}

func (c *Controller) DeleteVolumeSnapshot() *api.Response {

	dckInfo, err := c.searcher.SearchDockByVolume(c.request.GetVolumeId())
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
