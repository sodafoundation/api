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
	"strconv"

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

func (c *Controller) CreateVolume(ctx *c.Context, in *model.VolumeSpec, errchanVolume chan error) {
	var err error
	var profile *model.ProfileSpec

	if in.ProfileId == "" {
		log.Warning("Use default profile when user doesn't specify profile.")
		profile, err = db.C.GetDefaultProfile(ctx)
	} else {
		profile, err = db.C.GetProfile(ctx, in.ProfileId)
	}
	if err != nil {
		log.Error("Get profile failed: ", err)
		errchanVolume <- err
		return
	}

	var filterRequest map[string]interface{}
	if profile.Extras != nil {
		filterRequest = profile.Extras
	} else {
		filterRequest = make(map[string]interface{})
	}
	filterRequest["freeCapacity"] = ">= " + strconv.Itoa(int(in.Size))
	filterRequest["availabilityZone"] = in.AvailabilityZone

	polInfo, err := c.selector.SelectSupportedPool(filterRequest)
	if err != nil {
		errchanVolume <- err
		return
	}

	dockInfo, err := db.C.GetDock(ctx, polInfo.DockId)
	if err != nil {
		log.Error("When search supported dock resource:", err.Error())
		errchanVolume <- err
		return
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
		DriverName:       dockInfo.DriverName,
		Context:          ctx.ToJson(),
	}

	result, err := c.volumeController.CreateVolume(opt)
	if err != nil {
		//Change the status of the volume to error when the creation faild
		if errUpdate := c.UpdateStatus(ctx, in, model.VOLUME_ERROR); errUpdate != nil {
			errchanVolume <- errUpdate
			return
		}
		log.Error("When create volume:", err.Error())
		errchanVolume <- err
		return
	}

	result.PoolId, result.ProfileId = opt.GetPoolId(), opt.GetProfileId()

	// Update the volume data in database.
	if err = c.UpdateStatus(ctx, result, model.VOLUME_AVAILABLE); err != nil {
		errchanVolume <- err
		return
	}

	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(profile)
	c.policyController.Setup(CREATE_LIFECIRCLE_FLAG)
	c.policyController.SetDock(dockInfo)

	var errChanPolicy = make(chan error, 1)
	defer close(errChanPolicy)
	volBody, _ := json.Marshal(result)
	go c.policyController.ExecuteAsyncPolicy(opt, string(volBody), errChanPolicy)
	if err := <-errChanPolicy; err != nil {
		log.Error("When execute async policy:", err)
		errchanVolume <- err
		return
	}
	errchanVolume <- nil
}

func (c *Controller) DeleteVolume(ctx *c.Context, in *model.VolumeSpec, errchanvol chan error) {
	prf, err := db.C.GetProfile(ctx, in.ProfileId)
	if err != nil {
		log.Error("when search profile in db:", err)
		errchanvol <- err
		return
	}

	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(prf)
	c.policyController.Setup(DELETE_LIFECIRCLE_FLAG)

	dockInfo, err := db.C.GetDockByPoolId(ctx, in.PoolId)
	if err != nil {
		log.Error("When search dock in db by pool id: ", err)
		errchanvol <- err
		return
	}
	c.policyController.SetDock(dockInfo)
	c.volumeController.SetDock(dockInfo)

	opt := &pb.DeleteVolumeOpts{
		Id:         in.Id,
		Metadata:   in.Metadata,
		DriverName: dockInfo.DriverName,
		Context:    ctx.ToJson(),
	}

	var errChan = make(chan error, 1)
	defer close(errChan)
	go c.policyController.ExecuteAsyncPolicy(opt, "", errChan)

	if err := <-errChan; err != nil {
		log.Error("When execute async policy:", err)
		errchanvol <- err
		return
	}

	err = c.volumeController.DeleteVolume(opt)
	if err != nil {
		if errUpdate := c.UpdateStatus(ctx, in, model.VOLUEM_ERROR_DELETING); errUpdate != nil {
			errchanvol <- errUpdate
			return
		}
		errchanvol <- err
		return
	}
	if err = db.C.DeleteVolume(ctx, opt.GetId()); err != nil {
		log.Error("Error occurred in dock module when delete volume in db:", err.Error())
		errchanvol <- err
		return
	}
	errchanvol <- nil
}

// ExtendVolume ...
func (c *Controller) ExtendVolume(ctx *c.Context, volID string, newSize int64, errchanVolume chan error) {
	vol, err := db.C.GetVolume(ctx, volID)
	var volumeSize = vol.Size
	if err != nil {
		log.Error("Get volume failed in extend volume method: ", err.Error())
		errchanVolume <- err
		return
	}

	if newSize > vol.Size {
		pool, err := db.C.GetPool(ctx, vol.PoolId)
		if nil != err {
			log.Error("Get pool failed in extend volume method: ", err.Error())
			errchanVolume <- err
			return
		}

		if pool.FreeCapacity >= (newSize - vol.Size) {
			vol.Size = newSize
		} else {
			reason := fmt.Sprintf("pool free capacity(%d) < new size(%d) - old size(%d)",
				pool.FreeCapacity, newSize, vol.Size)
			errchanVolume <- errors.New(reason)
			return
		}
	} else {
		reason := fmt.Sprintf("new size(%d) <= old size(%d)", newSize, vol.Size)
		errchanVolume <- errors.New(reason)
		log.Error(reason)
		return
	}

	prf, err := db.C.GetProfile(ctx, vol.ProfileId)
	if err != nil {
		log.Error("when search profile in db:", err)
		errchanVolume <- err
		return
	}

	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(prf)
	c.policyController.Setup(EXTEND_LIFECIRCLE_FLAG)

	dockInfo, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
	if err != nil {
		log.Error("When search dock in db by pool id: ", err.Error())
		errchanVolume <- err
		return

	}
	c.policyController.SetDock(dockInfo)
	c.volumeController.SetDock(dockInfo)

	opt := &pb.ExtendVolumeOpts{
		Id:         vol.Id,
		Size:       vol.Size,
		Metadata:   vol.Metadata,
		DriverName: dockInfo.DriverName,
		Context:    ctx.ToJson(),
	}

	result, err := c.volumeController.ExtendVolume(opt)
	if err != nil {
		vol.Size = volumeSize
		if errUpdate := c.UpdateStatus(ctx, vol, model.VOLUME_ERROR); errUpdate != nil {
			errchanVolume <- errUpdate
			return
		}
		errchanVolume <- err
		return
	}
	result.PoolId, result.ProfileId = opt.GetPoolId(), opt.GetProfileId()

	// Update the volume data in database.
	if errUpdate := c.UpdateStatus(ctx, result, model.VOLUME_AVAILABLE); errUpdate != nil {
		errchanVolume <- errUpdate
		return
	}

	volBody, _ := json.Marshal(result)
	var errChan = make(chan error, 1)
	defer close(errChan)
	go c.policyController.ExecuteAsyncPolicy(opt, string(volBody), errChan)

	if err := <-errChan; err != nil {
		log.Error("When execute async policy:", err.Error())
		errchanVolume <- err
		return
	}

	errchanVolume <- nil
}

func (c *Controller) CreateVolumeAttachment(ctx *c.Context, in *model.VolumeAttachmentSpec, errchanVolAtm chan error) {
	vol, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in create volume attachment method: ", err)
		errchanVolAtm <- err
		return
	}
	dockInfo, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
	if err != nil {
		log.Error("When search supported dock resource:", err)
		errchanVolAtm <- err
		return
	}
	c.volumeController.SetDock(dockInfo)
	var atm = &pb.CreateAttachmentOpts{
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
		DriverName: dockInfo.DriverName,
		Context:    ctx.ToJson(),
	}
	result, err := c.volumeController.CreateVolumeAttachment(atm)
	if err != nil {
		if errUpdate := c.UpdateStatus(ctx, in, model.VOLUMEATM_ERROR); errUpdate != nil {
			errchanVolAtm <- errUpdate
			return
		}
		errchanVolAtm <- err
		return
	}
	if err = c.UpdateStatus(ctx, result, model.VOLUMEATM_AVAILABLE); err != nil {
		errchanVolAtm <- err
		return
	}
	errchanVolAtm <- nil
}

func (c *Controller) UpdateVolumeAttachment(in *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (c *Controller) DeleteVolumeAttachment(ctx *c.Context, in *model.VolumeAttachmentSpec, errchan chan error) {
	vol, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in delete volume attachment method: ", err)
		errchan <- err
		return
	}
	dockInfo, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
	if err != nil {
		log.Error("When search supported dock resource:", err)
		errchan <- err
		return
	}
	c.volumeController.SetDock(dockInfo)

	err = c.volumeController.DeleteVolumeAttachment(
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
			DriverName: dockInfo.DriverName,
			Context:    ctx.ToJson(),
		},
	)

	if err != nil {
		if errUpdate := c.UpdateStatus(ctx, in, model.VOLUMEATM_ERROR_DELETING); errUpdate != nil {
			errchan <- errUpdate
			return
		}
		errchan <- err
		return
	}
	if err := db.C.DeleteVolumeAttachment(ctx, in.Id); err != nil {
		log.Error("Error occurred in dock module when delete volume attachment in db:", err)
		errchan <- err
		return
	}

	errchan <- nil
}

func (c *Controller) CreateVolumeSnapshot(ctx *c.Context, in *model.VolumeSnapshotSpec, errchan chan error) {
	vol, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in create volume snapshot method: ", err)
		errchan <- err
		return
	}

	dockInfo, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
	if err != nil {
		log.Error("When search supported dock resource:", err)
		errchan <- err
		return
	}
	c.volumeController.SetDock(dockInfo)

	snp, err := c.volumeController.CreateVolumeSnapshot(
		&pb.CreateVolumeSnapshotOpts{
			Id:          in.Id,
			Name:        in.Name,
			Description: in.Description,
			VolumeId:    in.VolumeId,
			Size:        vol.Size,
			Metadata:    utils.MergeStringMaps(in.Metadata, vol.Metadata),
			DriverName:  dockInfo.DriverName,
			Context:     ctx.ToJson(),
		},
	)
	if err != nil {
		if errUpdate := c.UpdateStatus(ctx, in, model.VOLUMESNAP_ERROR); errUpdate != nil {
			errchan <- errUpdate
			return
		}
		errchan <- err
		return
	}
	if errUpdate := c.UpdateStatus(ctx, snp, model.VOLUMESNAP_AVAILABLE); errUpdate != nil {
		errchan <- errUpdate
		return
	}
	errchan <- nil
}

func (c *Controller) DeleteVolumeSnapshot(ctx *c.Context, in *model.VolumeSnapshotSpec, errchan chan error) {
	vol, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in delete volume snapshot method: ", err)
		errchan <- err
		return
	}
	dockInfo, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
	if err != nil {
		log.Error("When search supported dock resource:", err)
		errchan <- err
		return
	}
	c.volumeController.SetDock(dockInfo)

	err = c.volumeController.DeleteVolumeSnapshot(
		&pb.DeleteVolumeSnapshotOpts{
			Id:         in.Id,
			VolumeId:   in.VolumeId,
			Metadata:   utils.MergeStringMaps(in.Metadata, vol.Metadata),
			DriverName: dockInfo.DriverName,
			Context:    ctx.ToJson(),
		},
	)
	if err != nil {
		if errUpdate := c.UpdateStatus(ctx, in, model.VOLUMESNAP_ERROR_DELETING); errUpdate != nil {
			errchan <- errUpdate
			return
		}
		log.Error("Error occurred in dock module when delete volume snapshot in driver:", err)
		errchan <- err
		return
	}
	if err = db.C.DeleteVolumeSnapshot(ctx, in.Id); err != nil {
		log.Error("Error occurred in dock module when delete volume snapshot in db:", err)
		errchan <- err
		return
	}
	errchan <- nil
}

func (c *Controller) UpdateStatus(ctx *c.Context, in interface{}, status string) error {
	switch in.(type) {

	case *model.VolumeSnapshotSpec:
		snap := in.(*model.VolumeSnapshotSpec)
		snap.Status = status
		if _, errUpdate := db.C.UpdateVolumeSnapshot(ctx, snap.Id, snap); errUpdate != nil {
			log.Error("Error occurs when update volume snapshot status in db:", errUpdate.Error())
			return errUpdate
		}

	case *model.VolumeAttachmentSpec:
		attm := in.(*model.VolumeAttachmentSpec)
		attm.Status = status
		if _, errUpdate := db.C.UpdateVolumeAttachment(ctx, attm.Id, attm); errUpdate != nil {
			log.Error("Error occurred in dock module when update volume attachment status in db:", errUpdate)
			return errUpdate
		}

	case *model.VolumeSpec:
		vol := in.(*model.VolumeSpec)
		vol.Status = status
		if _, errUpdate := db.C.UpdateVolume(ctx, vol); errUpdate != nil {
			log.Error("When update volume status in db:", errUpdate.Error())
			return errUpdate
		}
	}
	return nil
}
