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
	"github.com/opensds/opensds/pkg/controller/selector"
	"github.com/opensds/opensds/pkg/controller/volume"
	"github.com/opensds/opensds/pkg/db"
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
	var c = &Controller{}

	// If volume input is not null, the controller will add policy orchestration
	// to manage volume resource.
	if vol != nil {
		prf, err := SearchProfile(vol.GetProfileId(), db.C)
		if err != nil {
			log.Error("when search profiles in db:", err)
			return nil, err
		}

		c.volume = vol
		c.profile = prf

		// Generate CreateVolumeOpts by parsing input VolumeSpec.
		c.createVolumeOpts = func(vol *model.VolumeSpec) *pb.CreateVolumeOpts {
			return &pb.CreateVolumeOpts{
				Id:               vol.GetId(),
				Name:             vol.GetName(),
				Description:      vol.GetDescription(),
				Size:             vol.GetSize(),
				AvailabilityZone: vol.GetAvailabilityZone(),
				ProfileId:        vol.GetProfileId(),
			}
		}(vol)
		// Generate DeleteVolumeOpts by parsing input VolumeSpec.
		c.deleteVolumeOpts = func(vol *model.VolumeSpec) *pb.DeleteVolumeOpts {
			return &pb.DeleteVolumeOpts{
				Id: vol.GetId(),
			}
		}(vol)

		// Initialize policy controller when profile is specified.
		c.policyController = policy.NewController(c.profile)
	}
	if atc != nil {
		// Generate CreateAttachment by parsing input VolumeAttachmentSpec.
		c.createAttachmentOpts = func(atc *model.VolumeAttachmentSpec) *pb.CreateAttachmentOpts {
			return &pb.CreateAttachmentOpts{
				Id:       atc.GetId(),
				VolumeId: atc.GetVolumeId(),
			}
		}(atc)
	}
	if snp != nil {
		c.volSnapshot = snp

		// Generate CreateVolumeSnapshotOpts by parsing input VolumeSnapshotSpec.
		c.createVolumeSnapshotOpts = func(snp *model.VolumeSnapshotSpec) *pb.CreateVolumeSnapshotOpts {
			return &pb.CreateVolumeSnapshotOpts{
				Id:          snp.GetId(),
				Name:        snp.GetName(),
				Description: snp.GetDescription(),
				Size:        snp.GetSize(),
				VolumeId:    snp.GetVolumeId(),
			}
		}(snp)
		// Generate DeleteVolumeSnapshotOpts by parsing input VolumeSnapshotSpec.
		c.deleteVolumeSnapshotOpts = func(snp *model.VolumeSnapshotSpec) *pb.DeleteVolumeSnapshotOpts {
			return &pb.DeleteVolumeSnapshotOpts{
				Id: snp.GetId(),
			}
		}(snp)
	}

	// Initialize selector controller.
	c.Selector = selector.NewSelector()
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
	selector.Selector

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

	polInfo, err := c.SelectSupportedPool(c.policyController.StorageTag().GetSyncTag())
	if err != nil {
		log.Error("When search supported pool resource:", err)
		return nil, err
	}

	dockInfo, err := c.SelectDock(polInfo)
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return nil, err
	}

	c.createVolumeOpts.PoolId = polInfo.GetId()
	c.createVolumeOpts.DockId = dockInfo.GetId()
	c.createVolumeOpts.DriverName = dockInfo.GetDriverName()

	c.policyController.SetDock(dockInfo)
	c.volumeController.SetDock(dockInfo)

	result, err := c.volumeController.CreateVolume()
	if err != nil {
		return nil, err
	}

	var errChan = make(chan error, 1)
	volBody, _ := json.Marshal(result)
	go c.policyController.ExecuteAsyncPolicy(c.createVolumeOpts, string(volBody), errChan)

	return result, nil
}

func (c *Controller) DeleteVolume() *model.Response {
	// Select the storage tag according to the lifecycle flag.
	c.policyController.Setup(DELETE_LIFECIRCLE_FLAG)

	dockInfo, err := c.SelectDock(c.deleteVolumeOpts.GetId())
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return &model.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}

	c.deleteVolumeOpts.DockId = dockInfo.GetId()
	c.deleteVolumeOpts.DriverName = dockInfo.GetDriverName()

	c.policyController.SetDock(dockInfo)
	c.volumeController.SetDock(dockInfo)

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
	dockInfo, err := c.SelectDock(c.createAttachmentOpts.GetVolumeId())
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return nil, err
	}

	c.createAttachmentOpts.DockId = dockInfo.GetId()
	c.createAttachmentOpts.DriverName = dockInfo.GetDriverName()

	c.volumeController.SetDock(dockInfo)

	return c.volumeController.CreateVolumeAttachment()
}

func (c *Controller) UpdateVolumeAttachment() (*model.VolumeAttachmentSpec, error) {
	dockInfo, err := c.SelectDock(c.createAttachmentOpts.GetVolumeId())
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return nil, err
	}

	c.createAttachmentOpts.DockId = dockInfo.GetId()
	c.createAttachmentOpts.DriverName = dockInfo.GetDriverName()

	c.volumeController.SetDock(dockInfo)

	return c.volumeController.UpdateVolumeAttachment()
}

func (c *Controller) DeleteVolumeAttachment() *model.Response {
	dockInfo, err := c.SelectDock(c.createAttachmentOpts.GetVolumeId())
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return &model.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}

	c.createAttachmentOpts.DockId = dockInfo.GetId()
	c.createAttachmentOpts.DriverName = dockInfo.GetDriverName()

	c.volumeController.SetDock(dockInfo)

	return c.volumeController.DeleteVolumeAttachment()
}

func (c *Controller) CreateVolumeSnapshot() (*model.VolumeSnapshotSpec, error) {
	dockInfo, err := c.SelectDock(c.createVolumeSnapshotOpts.GetVolumeId())
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return nil, err
	}

	c.createVolumeSnapshotOpts.DockId = dockInfo.GetId()
	c.createVolumeSnapshotOpts.DriverName = dockInfo.GetDriverName()

	c.volumeController.SetDock(dockInfo)

	return c.volumeController.CreateVolumeSnapshot()
}

func (c *Controller) DeleteVolumeSnapshot() *model.Response {
	dockInfo, err := c.SelectDock(c.volSnapshot.VolumeId)
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return &model.Response{
			Status: "Failure",
			Error:  fmt.Sprint(err),
		}
	}

	c.deleteVolumeSnapshotOpts.DockId = dockInfo.GetId()
	c.deleteVolumeSnapshotOpts.DriverName = dockInfo.GetDriverName()

	c.volumeController.SetDock(dockInfo)

	return c.volumeController.DeleteVolumeSnapshot()
}
