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

	log "github.com/golang/glog"

	"github.com/opensds/opensds/pkg/controller/policy"
	"github.com/opensds/opensds/pkg/controller/volume"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
)

const (
	CREATE_LIFECIRCLE_FLAG = iota + 1
	GET_LIFECIRCLE_FLAG
	LIST_LIFECIRCLE_FLAG
	DELETE_LIFECIRCLE_FLAG
)

func NewControllerWithVolumeConfig(
	vol *model.VolumeSpec,
	atc *model.VolumeAttachmentSpec,
	snp *model.VolumeSnapshotSpec,
) (*Controller, error) {
	c := &Controller{
		createVolumeOpts:         &pb.CreateVolumeOpts{},
		deleteVolumeOpts:         &pb.DeleteVolumeOpts{},
		createVolumeSnapshotOpts: &pb.CreateVolumeSnapshotOpts{},
		deleteVolumeSnapshotOpts: &pb.DeleteVolumeSnapshotOpts{},
		createAttachmentOpts:     &pb.CreateAttachmentOpts{},
	}

	c.searcher = NewDbSearcher()

	// If volume input is not null, the controller will add policy orchestration
	// to manage volume resource.
	if vol != nil {
		c.volume = vol
		prf, err := c.searcher.SearchProfile(vol.GetProfileId())
		if err != nil {
			log.Error("when search profiles in db:", err)
			return nil, err
		}

		c.profile = prf
		c.createVolumeOpts.Id = vol.GetId()
		c.createVolumeOpts.Name = vol.GetName()
		c.createVolumeOpts.Description = vol.GetDescription()
		c.createVolumeOpts.Size = vol.GetSize()
		c.createVolumeOpts.ProfileId = prf.GetId()

		c.deleteVolumeOpts.Id = vol.GetId()

		// Initialize policy controller when profile is specified.
		c.policyController = policy.NewController(c.profile)
	}
	if atc != nil {

		c.createAttachmentOpts.VolumeId = atc.GetVolumeId()
		c.createAttachmentOpts.Id = atc.GetId()
	}
	if snp != nil {
		c.volSnapshot = snp
		c.createVolumeSnapshotOpts.Id = snp.GetId()
		c.createVolumeSnapshotOpts.Name = snp.GetId()
		c.createVolumeSnapshotOpts.Description = snp.GetDescription()
		c.createVolumeSnapshotOpts.VolumeId = snp.GetVolumeId()
		c.createVolumeSnapshotOpts.Size = snp.GetSize()

		c.deleteVolumeSnapshotOpts.Id = snp.GetId()
	}

	// Initialize volume controller.
	c.volumeController = volume.NewController(
		c.createVolumeOpts,
		c.deleteVolumeOpts,
		c.createVolumeSnapshotOpts,
		c.deleteVolumeSnapshotOpts,
		c.createAttachmentOpts)

	return c, nil
}

type Controller struct {
	searcher                 Searcher
	volumeController         volume.Controller
	policyController         policy.Controller
	profile                  *model.ProfileSpec
	volume                   *model.VolumeSpec
	volSnapshot              *model.VolumeSnapshotSpec
	createVolumeOpts         *pb.CreateVolumeOpts
	deleteVolumeOpts         *pb.DeleteVolumeOpts
	createVolumeSnapshotOpts *pb.CreateVolumeSnapshotOpts
	deleteVolumeSnapshotOpts *pb.DeleteVolumeSnapshotOpts
	createAttachmentOpts     *pb.CreateAttachmentOpts
}

func (c *Controller) CreateVolume() (*model.VolumeSpec, error) {
	// Select the storage tag according to the lifecycle flag.
	c.policyController.Setup(CREATE_LIFECIRCLE_FLAG)

	polInfo, err := c.searcher.SearchSupportedPool(c.policyController.StorageTag().GetSyncTag())
	if err != nil {
		log.Error("When search supported pool resource:", err)
		return &model.VolumeSpec{}, err
	}
	c.createVolumeOpts.PoolId = polInfo.GetId()

	dockInfo, err := c.searcher.SearchDockByPool(polInfo)
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return &model.VolumeSpec{}, err
	}
	c.createVolumeOpts.DockId = dockInfo.Id
	c.policyController.SetDock(dockInfo)
	c.volumeController.SetDock(dockInfo)

	result, err := c.volumeController.CreateVolume()
	if err != nil {
		return &model.VolumeSpec{}, err
	}

	var errChan = make(chan error, 1)
	volBody, _ := json.Marshal(result)
	go c.policyController.ExecuteAsyncPolicy(c.createVolumeOpts, string(volBody), errChan)

	return result, nil
}

func (c *Controller) DeleteVolume() *model.Response {
	// Select the storage tag according to the lifecycle flag.
	c.policyController.Setup(DELETE_LIFECIRCLE_FLAG)

	dockInfo, err := c.searcher.SearchDockByVolume(c.deleteVolumeOpts.GetId())
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return &model.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}
	c.policyController.SetDock(dockInfo)
	c.volumeController.SetDock(dockInfo)
	c.deleteVolumeOpts.DockId = dockInfo.Id

	var errChan = make(chan error, 1)
	go c.policyController.ExecuteAsyncPolicy(c.deleteVolumeOpts, "", errChan)

	if err := <-errChan; err != nil {
		log.Error("When execute async policy:", err)
		return &model.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}

	return c.volumeController.DeleteVolume()
}

func (c *Controller) CreateVolumeAttachment() (*model.VolumeAttachmentSpec, error) {

	dockInfo, err := c.searcher.SearchDockByVolume(c.createAttachmentOpts.GetVolumeId())
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return &model.VolumeAttachmentSpec{}, err
	}

	c.createAttachmentOpts.DockId = dockInfo.Id
	c.volumeController.SetDock(dockInfo)

	return c.volumeController.CreateVolumeAttachment()
}

func (c *Controller) UpdateVolumeAttachment() (*model.VolumeAttachmentSpec, error) {

	dockInfo, err := c.searcher.SearchDockByVolume(c.createAttachmentOpts.GetVolumeId())
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return &model.VolumeAttachmentSpec{}, err
	}

	c.createAttachmentOpts.DockId = dockInfo.Id
	c.volumeController.SetDock(dockInfo)

	return c.volumeController.UpdateVolumeAttachment()
}

func (c *Controller) DeleteVolumeAttachment() *model.Response {

	dockInfo, err := c.searcher.SearchDockByVolume(c.createAttachmentOpts.GetVolumeId())
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return &model.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}

	c.createAttachmentOpts.DockId = dockInfo.Id
	c.volumeController.SetDock(dockInfo)

	return c.volumeController.DeleteVolumeAttachment()
}

func (c *Controller) CreateVolumeSnapshot() (*model.VolumeSnapshotSpec, error) {

	dockInfo, err := c.searcher.SearchDockByVolume(c.createVolumeSnapshotOpts.GetVolumeId())
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return &model.VolumeSnapshotSpec{}, err
	}

	c.createVolumeSnapshotOpts.DockId = dockInfo.Id
	c.volumeController.SetDock(dockInfo)

	return c.volumeController.CreateVolumeSnapshot()
}

func (c *Controller) DeleteVolumeSnapshot() *model.Response {

	dockInfo, err := c.searcher.SearchDockByVolume(c.volSnapshot.VolumeId)
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return &model.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}

	c.deleteVolumeSnapshotOpts.DockId = dockInfo.Id
	c.volumeController.SetDock(dockInfo)

	return c.volumeController.DeleteVolumeSnapshot()
}
