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
	"fmt"

	log "github.com/golang/glog"

	"github.com/opensds/opensds/pkg/controller/policy"
	"github.com/opensds/opensds/pkg/controller/selector"
	"github.com/opensds/opensds/pkg/controller/volume"
	"github.com/opensds/opensds/pkg/db"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
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
		selector:         selector.NewSelector(),
		volumeController: volume.NewController(),
	}
}

type Controller struct {
	selector         selector.Selector
	volumeController volume.Controller
	policyController policy.Controller
}

func (c *Controller) CreateVolume(in *model.VolumeSpec) (*model.VolumeSpec, error) {
	var profile *model.ProfileSpec
	var err error

	if in.ProfileId == "" {
		log.Warning("Use default profile when user doesn't specify profile.")
		profile, err = db.C.GetDefaultProfile()
	} else {
		profile, err = db.C.GetProfile(in.ProfileId)
	}
	if err != nil {
		log.Error("Get profile failed: ", err)
		return nil, err
	}

	if in.Size <= 0 {
		errMsg := fmt.Sprintf("Invalid volume size: %d", in.Size)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	if in.AvailabilityZone == "" {
		log.Warning("Use default availability zone when user doesn't specify availabilityZone.")
		in.AvailabilityZone = "default"
	}

	var filterRequest map[string]interface{}
	if profile.Extras != nil {
		filterRequest = profile.Extras
	} else {
		filterRequest = make(map[string]interface{})
	}
	filterRequest["size"] = in.Size
	filterRequest["availabilityZone"] = in.AvailabilityZone

	polInfo, err := c.selector.SelectSupportedPool(filterRequest)
	if err != nil {
		log.Error("When search supported pool resource:", err)
		return nil, err
	}
	dockInfo, err := db.C.GetDock(polInfo.DockId)
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return nil, err
	}

	c.volumeController.SetDock(dockInfo)
	opt := &pb.CreateVolumeOpts{
		Id:               in.Id,
		Name:             in.Name,
		Description:      in.Description,
		Size:             in.Size,
		AvailabilityZone: in.AvailabilityZone,
		ProfileId:        profile.Id,
		PoolId:           polInfo.Id,
		PoolName:         polInfo.Name,
		DockId:           dockInfo.Id,
		DriverName:       dockInfo.DriverName,
	}
	result, err := c.volumeController.CreateVolume(opt)
	if err != nil {
		return nil, err
	}

	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(profile)
	c.policyController.Setup(CREATE_LIFECIRCLE_FLAG)
	c.policyController.SetDock(dockInfo)

	var errChan = make(chan error, 1)
	volBody, _ := json.Marshal(result)
	go c.policyController.ExecuteAsyncPolicy(opt, string(volBody), errChan)

	return result, nil
}

func (c *Controller) DeleteVolume(in *model.VolumeSpec) error {
	prf, err := db.C.GetProfile(in.ProfileId)
	if err != nil {
		log.Error("when search profile in db:", err)
		return err
	}

	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(prf)
	c.policyController.Setup(DELETE_LIFECIRCLE_FLAG)

	dockInfo, err := db.C.GetDockByPoolId(in.PoolId)
	if err != nil {
		log.Error("When search dock in db by pool id: ", err)
		return err
	}
	c.policyController.SetDock(dockInfo)
	c.volumeController.SetDock(dockInfo)

	opt := &pb.DeleteVolumeOpts{
		Id:         in.Id,
		Metadata:   in.Metadata,
		DockId:     dockInfo.Id,
		DriverName: dockInfo.DriverName,
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
	vol, err := db.C.GetVolume(in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in create volume attachment method: ", err)
		return nil, err
	}
	dockInfo, err := db.C.GetDockByPoolId(vol.PoolId)
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return nil, err
	}
	c.volumeController.SetDock(dockInfo)

	return c.volumeController.CreateVolumeAttachment(
		&pb.CreateAttachmentOpts{
			Id:       in.Id,
			VolumeId: in.VolumeId,
			HostInfo: &pb.HostInfo{
				Platform:  in.Platform,
				OsType:    in.OsType,
				Ip:        in.Ip,
				Host:      in.Host,
				Initiator: in.Initiator,
			},
			Metadata:   utils.MergeStringMaps(in.Metadata, vol.Metadata),
			DockId:     dockInfo.Id,
			DriverName: dockInfo.DriverName,
		},
	)
}

func (c *Controller) UpdateVolumeAttachment(in *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (c *Controller) DeleteVolumeAttachment(in *model.VolumeAttachmentSpec) error {
	vol, err := db.C.GetVolume(in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in delete volume attachment method: ", err)
		return err
	}
	dockInfo, err := db.C.GetDockByPoolId(vol.PoolId)
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return err
	}
	c.volumeController.SetDock(dockInfo)

	return c.volumeController.DeleteVolumeAttachment(
		&pb.DeleteAttachmentOpts{
			Id:       in.Id,
			VolumeId: in.VolumeId,
			HostInfo: &pb.HostInfo{
				Platform:  in.Platform,
				OsType:    in.OsType,
				Ip:        in.Ip,
				Host:      in.Host,
				Initiator: in.Initiator,
			},
			Metadata:   utils.MergeStringMaps(in.Metadata, vol.Metadata),
			DockId:     dockInfo.Id,
			DriverName: dockInfo.DriverName,
		},
	)
}

func (c *Controller) CreateVolumeSnapshot(in *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	vol, err := db.C.GetVolume(in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in create volume snapshot method: ", err)
		return nil, err
	}

	dockInfo, err := db.C.GetDockByPoolId(vol.PoolId)
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return nil, err
	}
	c.volumeController.SetDock(dockInfo)

	return c.volumeController.CreateVolumeSnapshot(
		&pb.CreateVolumeSnapshotOpts{
			Id:          in.Id,
			Name:        in.Name,
			Description: in.Description,
			VolumeId:    in.VolumeId,
			Size:        vol.Size,
			Metadata:    utils.MergeStringMaps(in.Metadata, vol.Metadata),
			DockId:      dockInfo.Id,
			DriverName:  dockInfo.DriverName,
		},
	)
}

func (c *Controller) DeleteVolumeSnapshot(in *model.VolumeSnapshotSpec) error {
	vol, err := db.C.GetVolume(in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in delete volume snapshot method: ", err)
		return err
	}
	dockInfo, err := db.C.GetDockByPoolId(vol.PoolId)
	if err != nil {
		log.Error("When search supported dock resource:", err)
		return err
	}
	c.volumeController.SetDock(dockInfo)

	return c.volumeController.DeleteVolumeSnapshot(
		&pb.DeleteVolumeSnapshotOpts{
			Id:         in.Id,
			VolumeId:   in.VolumeId,
			Metadata:   utils.MergeStringMaps(in.Metadata, vol.Metadata),
			DockId:     dockInfo.Id,
			DriverName: dockInfo.DriverName,
		},
	)
}
