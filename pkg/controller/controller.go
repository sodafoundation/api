// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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
	c "github.com/opensds/opensds/pkg/context"
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
	EXTEND_LIFECIRCLE_FLAG
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

func (c *Controller) CreateVolume(ctx *c.Context, in *model.VolumeSpec) (*model.VolumeSpec, error) {
	var profile *model.ProfileSpec
	var err error

	if in.ProfileId == "" {
		log.Warning("Use default profile when user doesn't specify profile.")
		profile, err = db.C.GetDefaultProfile(ctx)
	} else {
		profile, err = db.C.GetProfile(ctx, in.ProfileId)
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
	dockInfo, err := db.C.GetDock(ctx, polInfo.DockId)
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
		Context:          ctx.ToJson(),
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

func (c *Controller) DeleteVolume(ctx *c.Context, in *model.VolumeSpec) error {
	prf, err := db.C.GetProfile(ctx, in.ProfileId)
	if err != nil {
		log.Error("when search profile in db:", err)
		return err
	}

	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(prf)
	c.policyController.Setup(DELETE_LIFECIRCLE_FLAG)

	dockInfo, err := db.C.GetDockByPoolId(ctx, in.PoolId)
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
		Context:    ctx.ToJson(),
	}

	var errChan = make(chan error, 1)
	go c.policyController.ExecuteAsyncPolicy(opt, "", errChan)

	if err := <-errChan; err != nil {
		log.Error("When execute async policy:", err)
		return err
	}

	return c.volumeController.DeleteVolume(opt)
}

// ExtendVolume ...
func (c *Controller) ExtendVolume(ctx *c.Context, in *model.VolumeSpec) (*model.VolumeSpec, error) {
	prf, err := db.C.GetProfile(ctx, in.ProfileId)
	if err != nil {
		log.Error("when search profile in db:", err)
		return nil, err
	}

	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(prf)
	c.policyController.Setup(EXTEND_LIFECIRCLE_FLAG)

	dockInfo, err := db.C.GetDockByPoolId(ctx, in.PoolId)
	if err != nil {
		log.Error("When search dock in db by pool id: ", err)
		return nil, err
	}
	c.policyController.SetDock(dockInfo)
	c.volumeController.SetDock(dockInfo)

	opt := &pb.ExtendVolumeOpts{
		Id:         in.Id,
		Size:       in.Size,
		Metadata:   in.Metadata,
		DockId:     dockInfo.Id,
		DriverName: dockInfo.DriverName,
		Context:    ctx.ToJson(),
	}

	result, err := c.volumeController.ExtendVolume(opt)
	if err != nil {
		return nil, err
	}

	volBody, _ := json.Marshal(result)
	var errChan = make(chan error, 1)
	go c.policyController.ExecuteAsyncPolicy(opt, string(volBody), errChan)

	if err := <-errChan; err != nil {
		log.Error("When execute async policy:", err)
		return nil, err
	}

	return result, nil
}

func (c *Controller) CreateVolumeAttachment(ctx *c.Context, in *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	vol, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in create volume attachment method: ", err)
		return nil, err
	}
	dockInfo, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
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
			Context:    ctx.ToJson(),
		},
	)
}

func (c *Controller) UpdateVolumeAttachment(in *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (c *Controller) DeleteVolumeAttachment(ctx *c.Context, in *model.VolumeAttachmentSpec) error {
	vol, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in delete volume attachment method: ", err)
		return err
	}
	dockInfo, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
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
			Context:    ctx.ToJson(),
		},
	)
}

func (c *Controller) CreateVolumeSnapshot(ctx *c.Context, in *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	vol, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in create volume snapshot method: ", err)
		return nil, err
	}

	dockInfo, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
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
			Context:     ctx.ToJson(),
		},
	)
}

func (c *Controller) DeleteVolumeSnapshot(ctx *c.Context, in *model.VolumeSnapshotSpec) error {
	vol, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in delete volume snapshot method: ", err)
		return err
	}
	dockInfo, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
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
			Context:    ctx.ToJson(),
		},
	)
}
