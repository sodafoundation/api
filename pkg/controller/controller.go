// Copyright 2017 The OpenSDS Authors.
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

package controller

import (
	"encoding/json"
	"errors"

	log "github.com/golang/glog"

	"github.com/opensds/opensds/pkg/controller/policy"
	"github.com/opensds/opensds/pkg/controller/selector"
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

var Brain *Controller

func NewController() *Controller {
	return &Controller{
		Selector:         selector.NewSelector(),
		volumeController: volume.NewController(),
	}
}

type Controller struct {
	selector.Selector

	volumeController volume.Controller
	policyController policy.Controller
}

func (c *Controller) CreateVolume(in *model.VolumeSpec) (*model.VolumeSpec, error) {
	var prfID = in.GetProfileId()

	prf, err := c.SelectProfile(prfID)
	if err != nil {
		log.Error("when search profiles in db:", err)
		return nil, err
	}

	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(prf)
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
	c.policyController.SetDock(dockInfo)
	c.volumeController.SetDock(dockInfo)

	opt := &pb.CreateVolumeOpts{
		Id:               in.GetId(),
		Name:             in.GetName(),
		Description:      in.GetDescription(),
		Size:             in.GetSize(),
		AvailabilityZone: in.GetAvailabilityZone(),
		ProfileId:        prfID,
		PoolId:           polInfo.GetId(),
		DockId:           dockInfo.GetId(),
		DriverName:       dockInfo.GetDriverName(),
	}
	result, err := c.volumeController.CreateVolume(opt)
	if err != nil {
		return nil, err
	}

	var errChan = make(chan error, 1)
	volBody, _ := json.Marshal(result)
	go c.policyController.ExecuteAsyncPolicy(opt, string(volBody), errChan)

	return result, nil
}

func (c *Controller) DeleteVolume(in *model.VolumeSpec) error {
	prf, err := c.SelectProfile(in.GetProfileId())
	if err != nil {
		log.Error("when search profiles in db:", err)
		return err
	}

	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(prf)
	c.policyController.Setup(DELETE_LIFECIRCLE_FLAG)

	dockInfo, err := c.SelectDock(in.GetId())
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return err
	}
	c.policyController.SetDock(dockInfo)
	c.volumeController.SetDock(dockInfo)

	opt := &pb.DeleteVolumeOpts{
		Id:         in.GetId(),
		Metadata:   in.GetMetadata(),
		DockId:     dockInfo.GetId(),
		DriverName: dockInfo.GetDriverName(),
	}

	var errChan = make(chan error, 1)
	go c.policyController.ExecuteAsyncPolicy(opt, "", errChan)

	if err := <-errChan; err != nil {
		log.Error("When execute async policy:", err)
		return err
	}

	return c.volumeController.DeleteVolume(opt)
}

func (c *Controller) CreateVolumeAttachment(in *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	dockInfo, err := c.SelectDock(in.GetVolumeId())
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return nil, err
	}
	c.volumeController.SetDock(dockInfo)

	return c.volumeController.CreateVolumeAttachment(
		&pb.CreateAttachmentOpts{
			Id:       in.GetId(),
			VolumeId: in.GetVolumeId(),
			HostInfo: &pb.HostInfo{
				Platform:  in.GetPlatform(),
				OsType:    in.GetOsType(),
				Ip:        in.GetIp(),
				Host:      in.GetHost(),
				Initiator: in.GetInitiator(),
			},
			Metadata:   in.GetMetadata(),
			DockId:     dockInfo.GetId(),
			DriverName: dockInfo.GetDriverName(),
		},
	)
}

func (c *Controller) UpdateVolumeAttachment(in *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (c *Controller) DeleteVolumeAttachment(in *model.VolumeAttachmentSpec) error {
	dockInfo, err := c.SelectDock(in.GetVolumeId())
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return err
	}
	c.volumeController.SetDock(dockInfo)

	return c.volumeController.DeleteVolumeAttachment(
		&pb.DeleteAttachmentOpts{
			Id:       in.GetId(),
			VolumeId: in.GetVolumeId(),
			HostInfo: &pb.HostInfo{
				Platform:  in.GetPlatform(),
				OsType:    in.GetOsType(),
				Ip:        in.GetIp(),
				Host:      in.GetHost(),
				Initiator: in.GetInitiator(),
			},
			Metadata:   in.GetMetadata(),
			DockId:     dockInfo.GetId(),
			DriverName: dockInfo.GetDriverName(),
		},
	)
}

func (c *Controller) CreateVolumeSnapshot(in *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	dockInfo, err := c.SelectDock(in.GetVolumeId())
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return nil, err
	}
	c.volumeController.SetDock(dockInfo)

	return c.volumeController.CreateVolumeSnapshot(
		&pb.CreateVolumeSnapshotOpts{
			Id:          in.GetId(),
			Name:        in.GetName(),
			Description: in.GetDescription(),
			Size:        in.GetSize(),
			VolumeId:    in.GetVolumeId(),
			Metadata:    in.GetMetadata(),
		},
	)
}

func (c *Controller) DeleteVolumeSnapshot(in *model.VolumeSnapshotSpec) error {
	dockInfo, err := c.SelectDock(in.GetVolumeId())
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return err
	}
	c.volumeController.SetDock(dockInfo)

	return c.volumeController.DeleteVolumeSnapshot(
		&pb.DeleteVolumeSnapshotOpts{
			Id:       in.GetId(),
			VolumeId: in.GetVolumeId(),
			Metadata: in.GetMetadata(),
		},
	)
}
