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
	"github.com/opensds/opensds/pkg/controller/dr"
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
	volCtrl := volume.NewController()
	return &Controller{
		selector:         selector.NewSelector(),
		volumeController: volCtrl,
		drController:     dr.NewController(volCtrl),
	}
}

type Controller struct {
	selector         selector.Selector
	volumeController volume.Controller
	drController     dr.Controller
	policyController policy.Controller
}

func (c *Controller) CreateVolume(ctx *c.Context, in *model.VolumeSpec, errchanVolume chan error) {
	var err error
	var profile *model.ProfileSpec
	var snap *model.VolumeSnapshotSpec
	var snapVol *model.VolumeSpec
	var snapSize int64

	if in.SnapshotId != "" {
		snap, err = db.C.GetVolumeSnapshot(ctx, in.SnapshotId)
		if err != nil {
			log.Error("Get snapshot failed in create volume method: ", err)
			if errUpdate := db.C.UpdateStatus(ctx, in, model.VolumeError); errUpdate != nil {
				errchanVolume <- errUpdate
				return
			}
			errchanVolume <- err
			return
		}
		snapVol, err = db.C.GetVolume(ctx, snap.VolumeId)
		if err != nil {
			log.Error("Get volume failed in create volume method: ", err)
			if errUpdate := db.C.UpdateStatus(ctx, in, model.VolumeError); errUpdate != nil {
				errchanVolume <- errUpdate
				return
			}
			errchanVolume <- err
			return
		}
		snapSize = snapVol.Size
	}

	if in.ProfileId == "" {
		log.Warning("Use default profile when user doesn't specify profile.")
		profile, err = db.C.GetDefaultProfile(ctx)
	} else {
		profile, err = db.C.GetProfile(ctx, in.ProfileId)
	}
	if err != nil {
		log.Error("Get profile failed: ", err)
		if errUpdate := db.C.UpdateStatus(ctx, in, model.VolumeError); errUpdate != nil {
			errchanVolume <- errUpdate
			return
		}
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
	if snapVol != nil {
		filterRequest["id"] = snapVol.PoolId
	}

	polInfo, err := c.selector.SelectSupportedPool(filterRequest)
	if err != nil {
		if errUpdate := db.C.UpdateStatus(ctx, in, model.VolumeError); errUpdate != nil {
			errchanVolume <- errUpdate
			return
		}
		errchanVolume <- err
		return
	}

	dockInfo, err := db.C.GetDock(ctx, polInfo.DockId)
	if err != nil {
		log.Error("When search supported dock resource:", err.Error())
		if errUpdate := db.C.UpdateStatus(ctx, in, model.VolumeError); errUpdate != nil {
			errchanVolume <- errUpdate
			return
		}
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
		SnapshotId:       in.SnapshotId,
		SnapshotSize:     snapSize,
		PoolName:         polInfo.Name,
		DriverName:       dockInfo.DriverName,
		Context:          ctx.ToJson(),
	}

	result, err := c.volumeController.CreateVolume(opt)
	if err != nil {
		//Change the status of the volume to error when the creation faild
		if errUpdate := db.C.UpdateStatus(ctx, in, model.VolumeError); errUpdate != nil {
			errchanVolume <- errUpdate
			return
		}
		log.Error("When create volume:", err.Error())
		errchanVolume <- err
		return
	}

	result.PoolId, result.ProfileId = opt.GetPoolId(), opt.GetProfileId()

	// Update the volume data in database.
	if err = db.C.UpdateStatus(ctx, result, model.VolumeAvailable); err != nil {
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
		if errUpdate := db.C.UpdateStatus(ctx, in, model.VolumeErrorDeleting); errUpdate != nil {
			errchanvol <- errUpdate
			return
		}

		errchanvol <- err
		return
	}

	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(prf)
	c.policyController.Setup(DELETE_LIFECIRCLE_FLAG)

	dockInfo, err := db.C.GetDockByPoolId(ctx, in.PoolId)
	if err != nil {
		log.Error("When search dock in db by pool id: ", err)
		if errUpdate := db.C.UpdateStatus(ctx, in, model.VolumeErrorDeleting); errUpdate != nil {
			errchanvol <- errUpdate
			return
		}

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
		if errUpdate := db.C.UpdateStatus(ctx, in, model.VolumeErrorDeleting); errUpdate != nil {
			errchanvol <- errUpdate
			return
		}

		errchanvol <- err
		return
	}

	err = c.volumeController.DeleteVolume(opt)
	if err != nil {
		if errUpdate := db.C.UpdateStatus(ctx, in, model.VolumeErrorDeleting); errUpdate != nil {
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
		if errUpdate := db.C.UpdateStatus(ctx, vol, model.VolumeError); errUpdate != nil {
			errchanVolume <- errUpdate
			return
		}
		errchanVolume <- err
		return
	}
	result.PoolId, result.ProfileId = opt.GetPoolId(), opt.GetProfileId()

	// Update the volume data in database.
	if errUpdate := db.C.UpdateStatus(ctx, result, model.VolumeAvailable); errUpdate != nil {
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

	pol, err := db.C.GetPool(ctx, vol.PoolId)
	if err != nil {
		log.Error("Get pool failed in create volume attachment method: ", err)
		errchanVolAtm <- err
		return
	}

	var protocol = pol.Extras.IOConnectivity.AccessProtocol
	if protocol == "" {
		// Default protocol is iscsi
		protocol = "iscsi"
	}

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
		AccessProtocol: protocol,
		Metadata:       utils.MergeStringMaps(in.Metadata, vol.Metadata),
		DriverName:     dockInfo.DriverName,
		Context:        ctx.ToJson(),
	}
	result, err := c.volumeController.CreateVolumeAttachment(atm)
	if err != nil {
		if errUpdate := db.C.UpdateStatus(ctx, in, model.VolumeAttachError); errUpdate != nil {
			errchanVolAtm <- errUpdate
			return
		}
		errchanVolAtm <- err
		return
	}
	if err = db.C.UpdateStatus(ctx, result, model.VolumeAttachAvailable); err != nil {
		errchanVolAtm <- err
		return
	}
	result.Status = model.VolumeAttachAvailable
	result.AccessProtocol = protocol
	if _, err = db.C.UpdateVolumeAttachment(ctx, result.Id, result); err != nil {
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
			AccessProtocol: in.AccessProtocol,
			Metadata:       utils.MergeStringMaps(in.Metadata, vol.Metadata),
			DriverName:     dockInfo.DriverName,
			Context:        ctx.ToJson(),
		},
	)

	if err != nil {
		if errUpdate := db.C.UpdateStatus(ctx, in, model.VolumeAttachErrorDeleting); errUpdate != nil {
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
		if errUpdate := db.C.UpdateStatus(ctx, in, model.VolumeSnapError); errUpdate != nil {
			errchan <- errUpdate
			return
		}
		errchan <- err
		return
	}
	if errUpdate := db.C.UpdateStatus(ctx, snp, model.VolumeSnapAvailable); errUpdate != nil {
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
		if errUpdate := db.C.UpdateStatus(ctx, in, model.VolumeSnapErrorDeleting); errUpdate != nil {
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

func (c *Controller) CreateVolumeGroup(ctx *c.Context, in *model.VolumeGroupSpec) error {
	polInfo, err := c.selector.SelectSupportedPoolForVG(in)
	if err != nil {
		msg := "No valid pool find for group"
		log.Error(msg)
		if err = db.C.UpdateStatus(ctx, in, model.VolumeGroupError); err != nil {
			return err
		}
		return errors.New(msg)
	}
	dockInfo, err := db.C.GetDock(ctx, polInfo.DockId)
	if err != nil {
		msg := "No valid dock find for group"
		log.Error(msg)
		if err = db.C.UpdateStatus(ctx, in, model.VolumeGroupError); err != nil {
			return err
		}
		return errors.New(msg)
	}

	c.volumeController.SetDock(dockInfo)

	opt := &pb.CreateVolumeGroupOpts{
		Id:               in.Id,
		Name:             in.Name,
		Description:      in.Description,
		AvailabilityZone: in.AvailabilityZone,
		DriverName:       dockInfo.DriverName,
		PoolId:           polInfo.Id,
		Context:          ctx.ToJson(),
	}

	_, err = c.volumeController.CreateVolumeGroup(opt)
	if err != nil {
		return err
	}

	// TODO Policy controller for the vg need to be modified.
	//	// Select the storage tag according to the lifecycle flag.
	//	c.policyController = policy.NewController(profile)
	//	c.policyController.Setup(CREATE_LIFECIRCLE_FLAG)
	//	c.policyController.SetDock(dockInfo)

	//	var errChanPolicy = make(chan error, 1)
	//	defer close(errChanPolicy)
	//	volBody, _ := json.Marshal(result)
	//	go c.policyController.ExecuteAsyncPolicy(opt, string(volBody), errChanPolicy)
	//	if err := <-errChanPolicy; err != nil {
	//		log.Error("When execute async policy:", err)
	//		errchanVolume <- err
	//		return
	//	}
	return nil
}

func (c *Controller) CreateReplication(ctx *c.Context, in *model.ReplicationSpec) (*model.ReplicationSpec, error) {
	// TODO: Get profile and do some policy action.

	pvol, err := db.C.GetVolume(ctx, in.PrimaryVolumeId)
	if err != nil {
		return nil, err
	}
	// TODO: If user does not provide the secondary volume. Do the following steps:
	// 1. Get profile from db.
	// 2. Use selector to choose backend.
	// 3. Create volume.
	// TODO: The secondary volume may be across region.
	svol, err := db.C.GetVolume(ctx, in.SecondaryVolumeId)
	if err != nil {
		return nil, err
	}

	result, err := c.drController.CreateReplication(ctx, in, pvol, svol)
	result.ReplicationStatus = model.ReplicationEnabled
	if err != nil {
		result.ReplicationStatus = model.ReplicationError
	}

	// update status ,driver data, metadata
	db.C.UpdateReplication(ctx, result.Id, result)
	return result, err
}

func (c *Controller) DeleteReplication(ctx *c.Context, in *model.ReplicationSpec) error {

	pvol, err := db.C.GetVolume(ctx, in.PrimaryVolumeId)
	if err != nil {
		return err
	}
	svol, err := db.C.GetVolume(ctx, in.SecondaryVolumeId)
	if err != nil {
		return err
	}

	err = c.drController.DeleteReplication(ctx, in, pvol, svol)
	if err != nil {
		db.C.UpdateStatus(ctx, in, model.ReplicationErrorDeleting)
	}
	return err
}

func (c *Controller) EnableReplication(ctx *c.Context, in *model.ReplicationSpec) error {
	pvol, err := db.C.GetVolume(ctx, in.PrimaryVolumeId)
	if err != nil {
		return err
	}
	svol, err := db.C.GetVolume(ctx, in.SecondaryVolumeId)
	if err != nil {
		return err
	}

	err = c.drController.EnableReplication(ctx, in, pvol, svol)
	in.ReplicationStatus = model.ReplicationEnabled
	if err != nil {
		in.ReplicationStatus = model.ReplicationErrorEnabling
	}
	if _, err := db.C.UpdateReplication(ctx, in.Id, in); err != nil {
		log.Error("update replication in db error, ", err)
	}
	return err
}

func (c *Controller) DisableReplication(ctx *c.Context, in *model.ReplicationSpec) error {
	pvol, err := db.C.GetVolume(ctx, in.PrimaryVolumeId)
	if err != nil {
		return err
	}
	svol, err := db.C.GetVolume(ctx, in.SecondaryVolumeId)
	if err != nil {
		return err
	}

	err = c.drController.DisableReplication(ctx, in, pvol, svol)
	in.ReplicationStatus = model.ReplicationDisabled
	if err != nil {
		in.ReplicationStatus = model.ReplicationErrorDisabling
	}
	if _, err := db.C.UpdateReplication(ctx, in.Id, in); err != nil {
		log.Error("update replication in db error, ", err)
	}

	return err
}

func (c *Controller) FailoverReplication(ctx *c.Context, replication *model.ReplicationSpec, failover *model.FailoverReplicationSpec) error {
	pvol, err := db.C.GetVolume(ctx, replication.PrimaryVolumeId)
	if err != nil {
		return err
	}
	svol, err := db.C.GetVolume(ctx, replication.SecondaryVolumeId)
	if err != nil {
		return err
	}

	err = c.drController.FailoverReplication(ctx, replication, failover, pvol, svol)
	if failover.SecondaryBackendId == model.ReplicationDefaultBackendId {
		if err != nil {
			replication.ReplicationStatus = model.ReplicationErrorFailover
		} else {
			replication.ReplicationStatus = model.ReplicationFailover
		}
	} else {
		if err != nil {
			replication.ReplicationStatus = model.ReplicationErrorFailback
		} else {
			replication.ReplicationStatus = model.ReplicationEnabled
		}
	}

	if _, err := db.C.UpdateReplication(ctx, replication.Id, replication); err != nil {
		log.Error("update replication in db error, ", err)
	}
	return err
}

func (c *Controller) UpdateVolumeGroup(ctx *c.Context, vg *model.VolumeGroupSpec, addVolumes []string, removeVolumes []string) error {
	dock, err := db.C.GetDockByPoolId(ctx, vg.PoolId)
	if err != nil {
		return err
	}

	c.volumeController.SetDock(dock)

	opt := &pb.UpdateVolumeGroupOpts{
		Id:            vg.Id,
		DriverName:    dock.DriverName,
		AddVolumes:    addVolumes,
		RemoveVolumes: removeVolumes,
		Context:       ctx.ToJson(),
	}

	err = c.volumeController.UpdateVolumeGroup(opt)
	if err != nil {
		log.Error("When create volume group:", err)
		return err
	}

	// TODO Policy controller for the vg need to be modified.
	//	// Select the storage tag according to the lifecycle flag.
	//	c.policyController = policy.NewController(profile)
	//	c.policyController.Setup(CREATE_LIFECIRCLE_FLAG)
	//	c.policyController.SetDock(dockInfo)

	//	var errChanPolicy = make(chan error, 1)
	//	defer close(errChanPolicy)
	//	volBody, _ := json.Marshal(result)
	//	go c.policyController.ExecuteAsyncPolicy(opt, string(volBody), errChanPolicy)
	//	if err := <-errChanPolicy; err != nil {
	//		log.Error("When execute async policy:", err)
	//		errchanVolume <- err
	//		return
	//	}
	return nil
}

func (c *Controller) DeleteVolumeGroup(ctx *c.Context, vg *model.VolumeGroupSpec) error {
	dock, err := db.C.GetDockByPoolId(ctx, vg.PoolId)
	if err != nil {
		return err
	}

	c.volumeController.SetDock(dock)

	opt := &pb.DeleteVolumeGroupOpts{
		Id:         vg.Id,
		DriverName: dock.DriverName,
		Context:    ctx.ToJson(),
	}

	err = c.volumeController.DeleteVolumeGroup(opt)
	if err != nil {
		log.Error("When delete volume group:", err)
		return err
	}

	return nil
}
