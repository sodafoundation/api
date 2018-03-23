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

	"time"

	log "github.com/golang/glog"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller/policy"
	"github.com/opensds/opensds/pkg/controller/selector"
	"github.com/opensds/opensds/pkg/controller/volume"
	"github.com/opensds/opensds/pkg/db"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/constants"
	"github.com/satori/go.uuid"
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

func (c *Controller) AsynCreateVolume(ctx *c.Context, in *model.VolumeSpec) (*model.VolumeSpec, *model.StoragePoolSpec, *model.ProfileSpec, error) {
	var profile *model.ProfileSpec
	var err error

	if in.Id == "" {
		in.Id = uuid.NewV4().String()
	}
	if in.ProfileId == "" {
		log.Warning("Use default profile when user doesn't specify profile.")
		profile, err = db.C.GetDefaultProfile(ctx)
	} else {
		profile, err = db.C.GetProfile(ctx, in.ProfileId)
	}

	if err != nil {
		log.Error("Get profile failed: ", err)
		return nil, nil, nil, err
	}

	if in.Size <= 0 {
		errMsg := fmt.Sprintf("Invalid volume size: %d", in.Size)
		log.Error(errMsg)
		return nil, nil, nil, errors.New(errMsg)
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
	filterRequest["freeCapacity"] = ">= " + strconv.Itoa(int(in.Size))
	filterRequest["availabilityZone"] = in.AvailabilityZone

	polInfo, err := c.selector.SelectSupportedPool(filterRequest)

	in.PoolId = polInfo.Id
	if in.CreatedAt == "" {
		in.CreatedAt = time.Now().Format(constants.TimeFormat)
	}

	vol := &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id:        in.Id,
			CreatedAt: in.CreatedAt,
		},
		UserId:           in.UserId,
		TenantId:         in.TenantId,
		Name:             in.Name,
		Description:      in.Description,
		Size:             in.Size,
		AvailabilityZone: in.AvailabilityZone,
		Status:           model.VOLUME_CREATING,
		ProfileId:        profile.Id,
		PoolId:           polInfo.Id,
	}
	result, err := db.C.CreateVolume(ctx, vol)
	if err != nil {
		log.Error("When add volume to db:", err)
		return nil, nil, nil, err
	}

	return result, polInfo, profile, nil
}

func (c *Controller) CreateVolume(ctx *c.Context, args ...interface{}) {
	in := args[0].(*model.VolumeSpec)
	polInfo := args[1].(*model.StoragePoolSpec)
	profile := args[2].(*model.ProfileSpec)
	errchanVolume := args[3].(chan error)

	defer close(errchanVolume)

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
		DockId:           dockInfo.Id,
		DriverName:       dockInfo.DriverName,
		Context:          ctx.ToJson(),
	}

	result, err := c.volumeController.CreateVolume(opt)
	if err != nil {
		log.Error("When create volume:", err.Error())
		errchanVolume <- err
		return
	}

	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(profile)
	c.policyController.Setup(CREATE_LIFECIRCLE_FLAG)
	c.policyController.SetDock(dockInfo)

	var errChanPolicy = make(chan error, 1)
	volBody, _ := json.Marshal(result)
	go c.policyController.ExecuteAsyncPolicy(opt, string(volBody), errChanPolicy)

	if err := <-errChanPolicy; err != nil {
		log.Error("When execute async policy:", err)
		errchanVolume <- err
		return
	}

	errchanVolume <- nil
}

//Just modify the state of the volume to be deleted in the DB, the real deletion in another thread
func (c *Controller) AsynDeleteVolume(ctx *c.Context, in *model.VolumeSpec) error {
	if in.Status != model.VOLUME_AVAILABLE {
		errMsg := "Only the volume with the status available can be deleted"
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	in.Status = model.VOLUME_DELETING
	_, err := db.C.UpdateVolume(ctx, in)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) DeleteVolume(ctx *c.Context, in *model.VolumeSpec, errchanvol chan error) {
	defer close(errchanvol)
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
		DockId:     dockInfo.Id,
		DriverName: dockInfo.DriverName,
		Context:    ctx.ToJson(),
	}

	var errChan = make(chan error, 1)
	go c.policyController.ExecuteAsyncPolicy(opt, "", errChan)

	if err := <-errChan; err != nil {
		log.Error("When execute async policy:", err)
		errchanvol <- err
		return
	}

	err = c.volumeController.DeleteVolume(opt)
	if err != nil {
		errchanvol <- err
		return
	}
	errchanvol <- nil
}

// ExtendVolume ...
func (c *Controller) ExtendVolume(ctx *c.Context, in *model.VolumeSpec, errchanVolume chan error) {
	prf, err := db.C.GetProfile(ctx, in.ProfileId)
	if err != nil {
		log.Error("when search profile in db:", err.Error())
		errchanVolume <- err
		return
	}
	defer close(errchanVolume)
	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(prf)
	c.policyController.Setup(EXTEND_LIFECIRCLE_FLAG)

	dockInfo, err := db.C.GetDockByPoolId(ctx, in.PoolId)
	if err != nil {
		log.Error("When search dock in db by pool id: ", err.Error())
		errchanVolume <- err
		return

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
		errchanVolume <- err
		return
	}

	volBody, _ := json.Marshal(result)
	var errChan = make(chan error, 1)
	go c.policyController.ExecuteAsyncPolicy(opt, string(volBody), errChan)

	if err := <-errChan; err != nil {
		log.Error("When execute async policy:", err.Error())
		errchanVolume <- err
		return
	}

	errchanVolume <- err
}

// AsynExtendVolume ...
func (c *Controller) AsynExtendVolume(ctx *c.Context, in *model.VolumeSpec) (*model.VolumeSpec, error) {
	if in.Status != model.VOLUME_AVAILABLE {
		errMsg := "The status of the volume to be extended must be available"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	in.Status = model.VOLUME_EXTENDING
	// Store the volume data into database.
	result, err := db.C.ExtendVolume(ctx, in)
	if err != nil {
		log.Error("When extend volume in db module:", err)
		return nil, err
	}
	return result, nil
}

func (c *Controller) AsynCreateVolumeAttachment(ctx *c.Context, in *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	vol, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in create volume attachment method: ", err)
		return nil, err
	}
	if vol.Status != model.VOLUME_AVAILABLE {
		errMsg := "Only the status of volume is available, attachment can be created"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if in.Id == "" {
		in.Id = uuid.NewV4().String()
	}
	if in.CreatedAt == "" {
		in.CreatedAt = time.Now().Format(constants.TimeFormat)
	}
	if len(in.AdditionalProperties) == 0 {
		in.AdditionalProperties = map[string]interface{}{"attachment": "attachment"}
	}
	if len(in.ConnectionData) == 0 {
		in.ConnectionData = map[string]interface{}{"attachment": "attachment"}
	}
	if in.Platform == "" {
		in.Platform = ""
	}
	if in.OsType == "" {
		in.OsType = ""
	}
	if in.Ip == "" {
		in.Ip = ""
	}
	if in.Host == "" {
		in.Host = ""
	}
	if in.Initiator == "" {
		in.Initiator = ""
	}

	var atc = &model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id:        in.Id,
			CreatedAt: in.CreatedAt,
		},
		VolumeId: in.VolumeId,
		HostInfo: model.HostInfo{
			Platform:  in.Platform,
			OsType:    in.OsType,
			Ip:        in.Ip,
			Host:      in.Host,
			Initiator: in.Initiator,
		},
		Status:         model.VOLUMEATM_CREATING,
		Metadata:       utils.MergeStringMaps(in.Metadata, vol.Metadata),
		ConnectionInfo: in.ConnectionInfo,
	}

	result, err := db.C.CreateVolumeAttachment(ctx, atc)
	if err != nil {
		log.Error("Error occurred in dock module when create volume attachment in db:", err)
		return nil, err
	}
	return result, nil
}

func (c *Controller) CreateVolumeAttachment(ctx *c.Context, in *model.VolumeAttachmentSpec, errchanVolAtm chan error) {
	defer close(errchanVolAtm)
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
		DockId:     dockInfo.Id,
		DriverName: dockInfo.DriverName,
		Context:    ctx.ToJson(),
	}
	_, err = c.volumeController.CreateVolumeAttachment(atm)

	if err != nil {
		errchanVolAtm <- err
		return
	}
	errchanVolAtm <- nil
}

func (c *Controller) UpdateVolumeAttachment(in *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (c *Controller) DeleteVolumeAttachment(ctx *c.Context, in *model.VolumeAttachmentSpec, errchan chan error) {
	defer close(errchan)
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
			DockId:     dockInfo.Id,
			DriverName: dockInfo.DriverName,
			Context:    ctx.ToJson(),
		},
	)

	if err != nil {
		errchan <- err
		return
	}
	errchan <- nil
}

func (c *Controller) AsynCreateVolumeSnapshot(ctx *c.Context, in *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	vol, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		log.Error("Get volume failed in create volume snapshot method: ", err)
		return nil, err
	}
	if vol.Status != model.VOLUME_AVAILABLE && vol.Status != model.VOLUME_IN_USE {
		var errMsg = "Only the status of volume is available or in-use, the snapshot can be created"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	if in.Id == "" {
		in.Id = uuid.NewV4().String()
	}

	if in.CreatedAt == "" {
		in.CreatedAt = time.Now().Format(constants.TimeFormat)
	}

	var snap = &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id:        in.Id,
			CreatedAt: in.CreatedAt,
		},
		Name:        in.Name,
		Description: in.Description,
		VolumeId:    in.VolumeId,
		Size:        vol.Size,
		Metadata:    utils.MergeStringMaps(in.Metadata, vol.Metadata),
		Status:      model.VOLUMESNAP_CREATING,
	}

	result, err := db.C.CreateVolumeSnapshot(ctx, snap)
	if err != nil {
		log.Error("Error occurred in dock module when create volume snapshot in db:", err)
		return nil, err
	}
	return result, nil
}

func (c *Controller) CreateVolumeSnapshot(ctx *c.Context, in *model.VolumeSnapshotSpec, errchan chan error) {
	defer close(errchan)
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

	_, err = c.volumeController.CreateVolumeSnapshot(
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
	if err != nil {
		errchan <- err
		return
	}
	errchan <- nil
}

func (c *Controller) AsynDeleteVolumeSnapshot(ctx *c.Context, in *model.VolumeSnapshotSpec) error {
	if in.Status != model.VOLUMESNAP_AVAILABLE {
		errMsg := "Only the volume snapshot with the status available can be deleted"
		log.Error(errMsg)
		return errors.New(errMsg)
	}
	in.Status = model.VOLUMESNAP_DELETING
	_, err := db.C.UpdateVolumeSnapshot(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) DeleteVolumeSnapshot(ctx *c.Context, in *model.VolumeSnapshotSpec, errchan chan error) {
	defer close(errchan)
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
			DockId:     dockInfo.Id,
			DriverName: dockInfo.DriverName,
			Context:    ctx.ToJson(),
		},
	)
	if err != nil {
		errchan <- err
		return
	}
	errchan <- nil
}
