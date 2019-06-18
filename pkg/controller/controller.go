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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/drivers/utils/config"
	osdsCtx "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller/dr"
	"github.com/opensds/opensds/pkg/controller/fileshare"
	"github.com/opensds/opensds/pkg/controller/metrics"
	"github.com/opensds/opensds/pkg/controller/policy"
	"github.com/opensds/opensds/pkg/controller/selector"
	"github.com/opensds/opensds/pkg/controller/volume"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils"
	"google.golang.org/grpc"
)

const (
	CREATE_LIFECIRCLE_FLAG = iota + 1
	GET_LIFECIRCLE_FLAG
	LIST_LIFECIRCLE_FLAG
	DELETE_LIFECIRCLE_FLAG
	EXTEND_LIFECIRCLE_FLAG
)

func NewController(port string) *Controller {
	volCtrl := volume.NewController()
	fileShareCtrl := fileshare.NewController()
	metricsCtrl := metrics.NewController()
	return &Controller{
		selector:            selector.NewSelector(),
		volumeController:    volCtrl,
		fileshareController: fileShareCtrl,
		metricsController:   metricsCtrl,
		drController:        dr.NewController(volCtrl),
		Port:                port,
	}
}

type Controller struct {
	selector            selector.Selector
	volumeController    volume.Controller
	fileshareController fileshare.Controller
	metricsController   metrics.Controller
	drController        dr.Controller
	policyController    policy.Controller

	Port string
}

// Run method would start the listen mechanism of controller module.
func (c *Controller) Run() error {
	// New Grpc Server
	s := grpc.NewServer()
	// Register controller service.
	pb.RegisterControllerServer(s, c)
	pb.RegisterFileShareControllerServer(s, c)

	// Listen the controller server port.
	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalf("failed to listen: %+v", err)
		return err
	}

	log.Info("Controller server initialized! Start listening on port:", lis.Addr())

	// Start controller server watching loop.
	defer s.Stop()
	return s.Serve(lis)
}

// CreateVolume implements pb.ControllerServer.CreateVolume
func (c *Controller) CreateVolume(contx context.Context, opt *pb.CreateVolumeOpts) (*pb.GenericResponse, error) {
	var err error
	var snap *model.VolumeSnapshotSpec
	var snapVol *model.VolumeSpec

	log.Info("Controller server receive create volume request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	prf := model.NewProfileFromJson(opt.Profile)

	if opt.SnapshotId != "" {
		snap, err = db.C.GetVolumeSnapshot(ctx, opt.SnapshotId)
		if err != nil {
			db.UpdateVolumeStatus(ctx, db.C, opt.Id, model.VolumeError)
			log.Error("get snapshot failed in create volume method: ", err)
			return pb.GenericResponseError(err), err
		}
		snapVol, err = db.C.GetVolume(ctx, snap.VolumeId)
		if err != nil {
			db.UpdateVolumeStatus(ctx, db.C, opt.Id, model.VolumeError)
			log.Error("get volume failed in create volume method: ", err)
			return pb.GenericResponseError(err), err
		}
		opt.SnapshotSize = snapVol.Size
		opt.PoolId = snapVol.PoolId
		opt.Metadata = utils.MergeStringMaps(opt.Metadata, snap.Metadata)
	}

	// This vol structure is currently fetched from database, but eventually
	// it will be removed after SelectSupportedPoolForVolume method in selector
	// is updated.
	vol, err := db.C.GetVolume(ctx, opt.Id)
	if err != nil {
		db.UpdateVolumeStatus(ctx, db.C, opt.Id, model.VolumeError)
		return pb.GenericResponseError(err), err
	}

	log.V(5).Infof("controller create volume:  get volume from db %+v", vol)

	polInfo, err := c.selector.SelectSupportedPoolForVolume(vol)
	if err != nil {
		db.UpdateVolumeStatus(ctx, db.C, opt.Id, model.VolumeError)
		return pb.GenericResponseError(err), err
	}

	// The default value of multi-attach is false, if it becomes true, then update into db
	log.V(5).Infof("update volume %+v", vol)

	if vol.MultiAttach {
		db.C.UpdateVolume(ctx, vol)
	}

	// whether specify a pool or not, opt's poolid and pool name should be
	// assigned by polInfo
	opt.PoolId = polInfo.Id
	opt.PoolName = polInfo.Name

	log.V(5).Infof("select pool %v and poolinfo : %v  for volume %+v", opt.PoolId, opt.PoolName, vol)

	dockInfo, err := db.C.GetDock(ctx, polInfo.DockId)
	if err != nil {
		db.UpdateVolumeStatus(ctx, db.C, opt.Id, model.VolumeError)
		log.Error("when search supported dock resource:", err.Error())
		return pb.GenericResponseError(err), err
	}
	c.volumeController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName

	log.V(5).Infof("selected driver name for create volume %+v", opt.DriverName)

	result, err := c.volumeController.CreateVolume(opt)
	if err != nil {
		// Change the status of the volume to error when the creation faild
		defer db.UpdateVolumeStatus(ctx, db.C, opt.Id, model.VolumeError)
		log.Error("when create volume:", err.Error())
		return pb.GenericResponseError(err), err
	}
	result.PoolId, result.ProfileId = opt.GetPoolId(), opt.GetProfileId()

	// Update the volume data in database.
	db.C.UpdateStatus(ctx, result, model.VolumeAvailable)

	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(prf)
	c.policyController.Setup(CREATE_LIFECIRCLE_FLAG)
	c.policyController.SetDock(dockInfo)

	var errChanPolicy = make(chan error, 1)
	defer close(errChanPolicy)
	volBody, _ := json.Marshal(result)
	go c.policyController.ExecuteAsyncPolicy(opt, string(volBody), errChanPolicy)
	if err := <-errChanPolicy; err != nil {
		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(result), nil
}

// DeleteVolume implements pb.ControllerServer.DeleteVolume
func (c *Controller) DeleteVolume(contx context.Context, opt *pb.DeleteVolumeOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive delete volume request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	prf := model.NewProfileFromJson(opt.Profile)
	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(prf)
	c.policyController.Setup(DELETE_LIFECIRCLE_FLAG)

	dockInfo, err := db.C.GetDockByPoolId(ctx, opt.PoolId)
	if err != nil {
		log.Error("when search dock in db by pool id: ", err)
		db.UpdateVolumeStatus(ctx, db.C, opt.Id, model.VolumeErrorDeleting)
		return pb.GenericResponseError(err), err
	}
	c.policyController.SetDock(dockInfo)
	c.volumeController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName

	var errChan = make(chan error, 1)
	defer close(errChan)
	go c.policyController.ExecuteAsyncPolicy(opt, "", errChan)

	if err := <-errChan; err != nil {
		log.Error("when execute async policy: ", err)
		db.UpdateVolumeStatus(ctx, db.C, opt.Id, model.VolumeErrorDeleting)
		return pb.GenericResponseError(err), err
	}

	if err = c.volumeController.DeleteVolume(opt); err != nil {
		db.UpdateVolumeStatus(ctx, db.C, opt.Id, model.VolumeErrorDeleting)
		return pb.GenericResponseError(err), err
	}

	if err = db.C.DeleteVolume(ctx, opt.GetId()); err != nil {
		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(nil), nil
}

// ExtendVolume implements pb.ControllerServer.ExtendVolume
func (c *Controller) ExtendVolume(contx context.Context, opt *pb.ExtendVolumeOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive extend volume request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	vol, err := db.C.GetVolume(ctx, opt.Id)
	if err != nil {
		log.Error("get volume failed in extend volume method: ", err.Error())
		return pb.GenericResponseError(err), err
	}

	// roll back size and status
	var rollBack = false
	defer func() {
		if rollBack {
			db.UpdateVolumeStatus(ctx, db.C, opt.Id, model.VolumeAvailable)
		}
	}()

	pool, err := db.C.GetPool(ctx, vol.PoolId)
	if nil != err {
		log.Error("get pool failed in extend volume method: ", err.Error())
		rollBack = true
		return pb.GenericResponseError(err), err
	}

	var newSize = opt.GetSize()
	if pool.FreeCapacity <= (newSize - vol.Size) {
		reason := fmt.Sprintf("pool free capacity(%d) < new size(%d) - old size(%d)",
			pool.FreeCapacity, newSize, vol.Size)
		rollBack = true
		return pb.GenericResponseError(reason), errors.New(reason)
	}
	opt.PoolId = pool.Id
	opt.PoolName = pool.Name
	prf := model.NewProfileFromJson(opt.Profile)

	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(prf)
	c.policyController.Setup(EXTEND_LIFECIRCLE_FLAG)

	dockInfo, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
	if err != nil {
		log.Error("when search dock in db by pool id: ", err.Error())
		rollBack = true
		return pb.GenericResponseError(err), err

	}
	c.policyController.SetDock(dockInfo)
	c.volumeController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName

	result, err := c.volumeController.ExtendVolume(opt)
	if err != nil {
		log.Error("extend volume failed: ", err.Error())
		rollBack = true
		return pb.GenericResponseError(err), err
	}

	// Update the volume data in database.
	result.Size = newSize
	result.PoolId, result.ProfileId = opt.GetPoolId(), opt.GetProfileId()
	db.C.UpdateStatus(ctx, result, model.VolumeAvailable)

	volBody, _ := json.Marshal(result)
	var errChan = make(chan error, 1)
	defer close(errChan)
	go c.policyController.ExecuteAsyncPolicy(opt, string(volBody), errChan)

	if err := <-errChan; err != nil {
		log.Error("when execute async policy:", err.Error())
		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(result), nil
}

// CreateVolumeAttachment implements pb.ControllerServer.CreateVolumeAttachment
func (c *Controller) CreateVolumeAttachment(contx context.Context, opt *pb.CreateVolumeAttachmentOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive create volume attachment request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	vol, err := db.C.GetVolume(ctx, opt.VolumeId)
	if err != nil {
		msg := fmt.Sprintf("get volume failed in create volume attachment method: %v", err)
		log.Error(msg)
		return pb.GenericResponseError(msg), err
	}

	opt.Metadata = utils.MergeStringMaps(opt.Metadata, vol.Metadata)

	pol, err := db.C.GetPool(ctx, vol.PoolId)
	if err != nil {
		msg := fmt.Sprintf("get pool failed in create volume attachment method: %v", err)
		log.Error(msg)
		return pb.GenericResponseError(msg), err
	}

	var protocol = pol.Extras.IOConnectivity.AccessProtocol
	if protocol == "" {
		// Default protocol is iscsi
		protocol = config.ISCSIProtocol
	}

	opt.AccessProtocol = protocol

	dockInfo, err := db.C.GetDock(ctx, pol.DockId)
	if err != nil {
		msg := fmt.Sprintf("when search supported dock resource: %v", err)
		log.Error(msg)
		return pb.GenericResponseError(msg), err
	}
	c.volumeController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName

	result, err := c.volumeController.CreateVolumeAttachment(opt)
	if err != nil {
		db.UpdateVolumeAttachmentStatus(ctx, db.C, opt.Id, model.VolumeAttachError)
		db.UpdateVolumeStatus(ctx, db.C, vol.Id, model.VolumeAvailable)
		msg := fmt.Sprintf("create volume attachment failed: %v", err)
		log.Error(msg)
		return pb.GenericResponseError(msg), err
	}

	result.AccessProtocol = protocol
	if vol.Status == model.VolumeAttaching {
		db.UpdateVolumeStatus(ctx, db.C, vol.Id, model.VolumeInUse)
	} else {
		msg := fmt.Sprintf("wrong volume status when volume attachment creation completed")
		log.Error(msg)
		return pb.GenericResponseError(msg), err
	}

	result.Status = model.VolumeAttachAvailable

	log.V(8).Infof("Create volume attachment successfully, the info is %v", result)
	// Save changes to db.
	db.C.UpdateVolumeAttachment(ctx, opt.Id, result)

	return pb.GenericResponseResult(result), nil
}

// DeleteVolumeAttachment implements pb.ControllerServer.DeleteVolumeAttachment
func (c *Controller) DeleteVolumeAttachment(contx context.Context, opt *pb.DeleteVolumeAttachmentOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive delete volume attachment request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	vol, err := db.C.GetVolume(ctx, opt.VolumeId)
	if err != nil {
		msg := fmt.Sprintf("get volume failed in delete volume attachment method: %v", err)
		log.Error(msg)
		return pb.GenericResponseError(msg), err
	}
	opt.Metadata = utils.MergeStringMaps(opt.Metadata, vol.Metadata)

	dockInfo, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
	if err != nil {
		msg := fmt.Sprintf("when search supported dock resource: %v", err)
		log.Error(msg)
		return pb.GenericResponseError(msg), err
	}

	c.volumeController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName

	if err = c.volumeController.DeleteVolumeAttachment(opt); err != nil {
		msg := fmt.Sprintf("delete volume attachment failed: %v", err)
		log.Error(msg)
		db.C.DeleteVolumeAttachment(ctx, opt.Id)
		return pb.GenericResponseError(msg), err
	}

	if err = db.C.DeleteVolumeAttachment(ctx, opt.Id); err != nil {
		msg := fmt.Sprintf("error occurred in dock module when delete volume attachment in db: %v", err)
		log.Error(msg)
		return pb.GenericResponseError(msg), err
	}

	db.UpdateVolumeStatus(ctx, db.C, vol.Id, model.VolumeAvailable)

	return pb.GenericResponseResult(nil), nil
}

// CreateVolumeSnapshot implements pb.ControllerServer.CreateVolumeSnapshot
func (c *Controller) CreateVolumeSnapshot(contx context.Context, opt *pb.CreateVolumeSnapshotOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive create volume snapshot request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	if opt.Metadata == nil {
		opt.Metadata = map[string]string{}
	}

	profile := model.NewProfileFromJson(opt.Profile)
	if profile.SnapshotProperties.Topology.Bucket != "" {
		opt.Metadata["bucket"] = profile.SnapshotProperties.Topology.Bucket
	}

	vol, err := db.C.GetVolume(ctx, opt.VolumeId)
	if err != nil {
		log.Error("get volume failed in create volume snapshot method: ", err)
		db.UpdateVolumeSnapshotStatus(ctx, db.C, opt.Id, model.VolumeSnapError)
		return pb.GenericResponseError(err), err
	}
	opt.Size = vol.Size
	opt.Metadata = utils.MergeStringMaps(opt.Metadata, vol.Metadata)

	dockInfo, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
	if err != nil {
		log.Error("when search supported dock resource: ", err)
		db.UpdateVolumeSnapshotStatus(ctx, db.C, opt.Id, model.VolumeSnapError)
		return pb.GenericResponseError(err), err
	}
	c.volumeController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName

	result, err := c.volumeController.CreateVolumeSnapshot(opt)
	if err != nil {
		db.UpdateVolumeSnapshotStatus(ctx, db.C, opt.Id, model.VolumeSnapError)
		return pb.GenericResponseError(err), err
	}

	db.C.UpdateStatus(ctx, result, model.VolumeSnapAvailable)
	return pb.GenericResponseResult(result), nil
}

// DeleteVolumeSnapshot implements pb.ControllerServer.DeleteVolumeSnapshot
func (c *Controller) DeleteVolumeSnapshot(contx context.Context, opt *pb.DeleteVolumeSnapshotOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive delete volume snapshot request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	vol, err := db.C.GetVolume(ctx, opt.VolumeId)
	if err != nil {
		log.Error("get volume failed in delete volume snapshot method: ", err)
		db.UpdateVolumeSnapshotStatus(ctx, db.C, opt.Id, model.VolumeSnapErrorDeleting)
		return pb.GenericResponseError(err), err
	}

	opt.Metadata = utils.MergeStringMaps(opt.Metadata, vol.Metadata)
	prf := model.NewProfileFromJson(opt.Profile)
	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(prf)
	c.policyController.Setup(DELETE_LIFECIRCLE_FLAG)

	dockInfo, err := db.C.GetDockByPoolId(ctx, vol.PoolId)
	if err != nil {
		log.Error("when search supported dock resource: ", err)
		db.UpdateVolumeSnapshotStatus(ctx, db.C, opt.Id, model.VolumeSnapErrorDeleting)
		return pb.GenericResponseError(err), err
	}
	c.volumeController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName

	if err = c.volumeController.DeleteVolumeSnapshot(opt); err != nil {
		log.Error("error occurred in controller module when delete volume snapshot: ", err)
		db.UpdateVolumeSnapshotStatus(ctx, db.C, opt.Id, model.VolumeSnapErrorDeleting)
		return pb.GenericResponseError(err), err
	}
	if err = db.C.DeleteVolumeSnapshot(ctx, opt.Id); err != nil {
		log.Error("error occurred in controller module when delete volume snapshot in db: ", err)
		db.UpdateVolumeSnapshotStatus(ctx, db.C, opt.Id, model.VolumeSnapErrorDeleting)
		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(nil), nil
}

// CreateReplication implements pb.ControllerServer.CreateReplication
func (c *Controller) CreateReplication(contx context.Context, opt *pb.CreateReplicationOpts) (*pb.GenericResponse, error) {
	// TODO: Get profile and do some policy action.

	log.Info("Controller server receive create volume replication request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	pvol, err := db.C.GetVolume(ctx, opt.PrimaryVolumeId)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationError)
		return pb.GenericResponseError(err), err
	}
	// TODO: If user does not provide the secondary volume. Do the following steps:
	// 1. Get profile from db.
	// 2. Use selector to choose backend.
	// 3. Create volume.
	// TODO: The secondary volume may be across region.
	svol, err := db.C.GetVolume(ctx, opt.SecondaryVolumeId)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationError)
		return pb.GenericResponseError(err), err
	}

	// This replica structure is currently fetched from database, but eventually
	// it will be removed after CreateReplication method in drController is
	// updated.
	replica, err := db.C.GetReplication(ctx, opt.Id)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationError)
		return pb.GenericResponseError(err), err
	}
	result, err := c.drController.CreateReplication(ctx, replica, pvol, svol)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationError)
		return pb.GenericResponseError(err), err
	}

	// update status ,driver data, metadata
	db.C.UpdateStatus(ctx, result, model.ReplicationAvailable)
	return pb.GenericResponseResult(result), nil
}

// DeleteReplication implements pb.ControllerServer.DeleteReplication
func (c *Controller) DeleteReplication(contx context.Context, opt *pb.DeleteReplicationOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive delete volume replication request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	pvol, err := db.C.GetVolume(ctx, opt.PrimaryVolumeId)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorDeleting)
		return pb.GenericResponseError(err), err
	}
	svol, err := db.C.GetVolume(ctx, opt.SecondaryVolumeId)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorDeleting)
		return pb.GenericResponseError(err), err
	}

	// This replica structure is currently fetched from database, but eventually
	// it will be removed after DeleteReplication method in drController is
	// updated.
	replica, err := db.C.GetReplication(ctx, opt.Id)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorDeleting)
		return pb.GenericResponseError(err), err
	}
	if err = c.drController.DeleteReplication(ctx, replica, pvol, svol); err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorDeleting)
		return pb.GenericResponseError(err), err
	}

	if err = db.C.DeleteReplication(ctx, opt.Id); err != nil {
		log.Error("error occurred in controller module when delete replication in db: ", err)
		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(nil), nil
}

// EnableReplication implements pb.ControllerServer.EnableReplication
func (c *Controller) EnableReplication(contx context.Context, opt *pb.EnableReplicationOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive enable volume replication request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	pvol, err := db.C.GetVolume(ctx, opt.PrimaryVolumeId)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorEnabling)
		return pb.GenericResponseError(err), err
	}
	svol, err := db.C.GetVolume(ctx, opt.SecondaryVolumeId)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorEnabling)
		return pb.GenericResponseError(err), err
	}

	// This replica structure is currently fetched from database, but eventually
	// it will be removed after EnableReplication method in drController is
	// updated.
	replica, err := db.C.GetReplication(ctx, opt.Id)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorEnabling)
		return pb.GenericResponseError(err), err
	}
	if err = c.drController.EnableReplication(ctx, replica, pvol, svol); err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorEnabling)
		return pb.GenericResponseError(err), err
	}

	db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationEnabled)
	return pb.GenericResponseResult(nil), nil
}

// DisableReplication implements pb.ControllerServer.DisableReplication
func (c *Controller) DisableReplication(contx context.Context, opt *pb.DisableReplicationOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive disable volume replication request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	pvol, err := db.C.GetVolume(ctx, opt.PrimaryVolumeId)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorDisabling)
		return pb.GenericResponseError(err), err
	}
	svol, err := db.C.GetVolume(ctx, opt.SecondaryVolumeId)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorDisabling)
		return pb.GenericResponseError(err), err
	}

	// This replica structure is currently fetched from database, but eventually
	// it will be removed after DisableReplication method in drController is
	// updated.
	replica, err := db.C.GetReplication(ctx, opt.Id)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorDisabling)
		return pb.GenericResponseError(err), err
	}
	if err = c.drController.DisableReplication(ctx, replica, pvol, svol); err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorDisabling)
		return pb.GenericResponseError(err), err
	}

	db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationDisabled)
	return pb.GenericResponseResult(nil), nil
}

// FailoverReplication implements pb.ControllerServer.FailoverReplication
func (c *Controller) FailoverReplication(contx context.Context, opt *pb.FailoverReplicationOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive failover volume replication request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	pvol, err := db.C.GetVolume(ctx, opt.PrimaryVolumeId)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorFailover)
		return pb.GenericResponseError(err), err
	}
	svol, err := db.C.GetVolume(ctx, opt.SecondaryVolumeId)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorFailover)
		return pb.GenericResponseError(err), err
	}

	var replicaStatus string
	var failover = &model.FailoverReplicationSpec{
		AllowAttachedVolume: opt.AllowAttachedVolume,
		SecondaryBackendId:  opt.SecondaryBackendId,
	}
	// This replica structure is currently fetched from database, but eventually
	// it will be removed after FailoverReplication method in drController is
	// updated.
	replica, err := db.C.GetReplication(ctx, opt.Id)
	if err != nil {
		db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorDisabling)
		return pb.GenericResponseError(err), err
	}
	err = c.drController.FailoverReplication(ctx, replica, failover, pvol, svol)
	if failover.SecondaryBackendId == model.ReplicationDefaultBackendId {
		if err != nil {
			db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorFailover)
			return pb.GenericResponseError(err), err
		}
		replicaStatus = model.ReplicationFailover
	} else {
		if err != nil {
			db.UpdateReplicationStatus(ctx, db.C, opt.Id, model.ReplicationErrorFailback)
			return pb.GenericResponseError(err), err
		}
		replicaStatus = model.ReplicationEnabled
	}

	db.UpdateReplicationStatus(ctx, db.C, opt.Id, replicaStatus)
	return pb.GenericResponseResult(nil), nil
}

// CreateVolumeGroup implements pb.ControllerServer.CreateVolumeGroup
func (c *Controller) CreateVolumeGroup(contx context.Context, opt *pb.CreateVolumeGroupOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive create volume group request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	// This vg structure is currently fetched from database, but eventually
	// it will be removed after SelectSupportedPoolForVG method in selector
	// is updated.
	vg, err := db.C.GetVolumeGroup(ctx, opt.Id)
	if err != nil {
		db.UpdateVolumeGroupStatus(ctx, db.C, opt.Id, model.VolumeGroupError)
		return pb.GenericResponseError(err), err
	}
	polInfo, err := c.selector.SelectSupportedPoolForVG(vg)
	if err != nil {
		log.Error("no valid pool find for group: ", err)
		db.UpdateVolumeGroupStatus(ctx, db.C, opt.Id, model.VolumeGroupError)
		return pb.GenericResponseError(err), err
	}
	opt.PoolId = polInfo.Id

	dockInfo, err := db.C.GetDock(ctx, polInfo.DockId)
	if err != nil {
		log.Error("no valid dock find for group: ", err)
		db.UpdateVolumeGroupStatus(ctx, db.C, opt.Id, model.VolumeGroupError)
		return pb.GenericResponseError(err), err
	}
	c.volumeController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName

	result, err := c.volumeController.CreateVolumeGroup(opt)
	if err != nil {
		db.UpdateVolumeGroupStatus(ctx, db.C, opt.Id, model.VolumeGroupError)
		return pb.GenericResponseError(err), err
	}
	result.PoolId = polInfo.Id

	// Update group id in the volumes
	for _, addVolId := range opt.AddVolumes {
		if _, err = db.C.UpdateVolume(ctx, &model.VolumeSpec{
			BaseModel: &model.BaseModel{Id: addVolId},
			GroupId:   opt.GetId(),
		}); err != nil {
			return pb.GenericResponseError(err), err
		}
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
	db.C.UpdateStatus(ctx, result, model.VolumeGroupAvailable)
	return pb.GenericResponseResult(result), nil
}

// UpdateVolumeGroup implements pb.ControllerServer.UpdateVolumeGroup
func (c *Controller) UpdateVolumeGroup(contx context.Context, opt *pb.UpdateVolumeGroupOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive update volume group request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	dock, err := db.C.GetDockByPoolId(ctx, opt.PoolId)
	if err != nil {
		db.UpdateVolumeGroupStatus(ctx, db.C, opt.Id, model.VolumeGroupError)
		return pb.GenericResponseError(err), err
	}
	c.volumeController.SetDock(dock)
	opt.DriverName = dock.DriverName

	vg, err := c.volumeController.UpdateVolumeGroup(opt)
	if err != nil {
		log.Error("when create volume group: ", err)
		db.UpdateVolumeGroupStatus(ctx, db.C, opt.Id, model.VolumeGroupError)

		//for _, addVol := range opt.AddVolumes {
		//	db.UpdateVolumeStatus(ctx, db.C, addVol, model.VolumeError)
		//}
		//for _, rmVol := range opt.RemoveVolumes {
		//	db.UpdateVolumeStatus(ctx, db.C, rmVol, model.VolumeError)
		//}

		return pb.GenericResponseError(err), err
	}

	// Update group id in the volumes
	for _, addVolId := range opt.AddVolumes {
		if _, err = db.C.UpdateVolume(ctx, &model.VolumeSpec{
			BaseModel: &model.BaseModel{Id: addVolId},
			GroupId:   opt.GetId(),
		}); err != nil {
			return pb.GenericResponseError(err), err
		}
	}

	for _, rmVolId := range opt.RemoveVolumes {
		if _, err = db.C.UpdateVolume(ctx, &model.VolumeSpec{
			BaseModel: &model.BaseModel{Id: rmVolId},
			GroupId:   "",
		}); err != nil {
			return pb.GenericResponseError(err), err
		}
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
	db.C.UpdateStatus(ctx, vg, model.VolumeGroupAvailable)
	return pb.GenericResponseResult(vg), nil
}

// DeleteVolumeGroup implements pb.ControllerServer.DeleteVolumeGroup
func (c *Controller) DeleteVolumeGroup(contx context.Context, opt *pb.DeleteVolumeGroupOpts) (*pb.GenericResponse, error) {
	ctx := osdsCtx.NewContextFromJson(opt.GetContext())

	log.Info("Controller server receive delete volume group request, vr =", opt)

	dock, err := db.C.GetDockByPoolId(ctx, opt.PoolId)
	if err != nil {
		db.UpdateVolumeGroupStatus(ctx, db.C, opt.Id, model.VolumeGroupErrorDeleting)
		return pb.GenericResponseError(err), err
	}
	c.volumeController.SetDock(dock)
	opt.DriverName = dock.DriverName

	if err = c.volumeController.DeleteVolumeGroup(opt); err != nil {
		log.Error("when delete volume group: ", err)
		db.UpdateVolumeGroupStatus(ctx, db.C, opt.Id, model.VolumeGroupErrorDeleting)
		return pb.GenericResponseError(err), err

	}

	if err = db.C.DeleteVolumeGroup(ctx, opt.Id); err != nil {
		log.Error("error occurred in controller module when delete volume group in db: ", err)
		db.UpdateVolumeGroupStatus(ctx, db.C, opt.Id, model.VolumeGroupErrorDeleting)
		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(nil), nil
}

// CreateFileShare implements pb.ControllerServer.CreateFileShare
func (c *Controller) CreateFileShare(contx context.Context, opt *pb.CreateFileShareOpts) (*pb.GenericResponse, error) {
	var err error

	log.Info("Controller server receive create file share request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	prf := model.NewProfileFromJson(opt.Profile)

	// This file share structure is currently fetched from database, but eventually
	// it will be removed after SelectSupportedPoolForFileShare method in selector
	// is updated.
	fileshare, err := db.C.GetFileShare(ctx, opt.Id)
	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareError)
		return pb.GenericResponseError(err), err
	}

	log.V(5).Infof("controller create fileshare:  get fileshare from db %+v", fileshare)

	polInfo, err := c.selector.SelectSupportedPoolForFileShare(fileshare)

	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareError)
		return pb.GenericResponseError(err), err
	}

	log.V(5).Infof("controller create fileshare:  selected poolInfo %+v", polInfo)
	// whether specify a pool or not, opt's poolid and pool name should be
	// assigned by polInfo
	opt.PoolId = polInfo.Id
	opt.PoolName = polInfo.Name

	dockInfo, err := db.C.GetDock(ctx, polInfo.DockId)
	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareError)
		log.Error("when search supported dock resource:", err.Error())
		return pb.GenericResponseError(err), err
	}
	c.fileshareController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName

	log.V(5).Infof("controller create fleshare:  selected Driver name %+v", opt.DriverName)

	result, err := c.fileshareController.CreateFileShare((*pb.CreateFileShareOpts)(opt))
	if err != nil {
		// Change the status of the file share to error when the creation faild
		defer db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareError)
		log.Error("when create file share:", err.Error())
		return pb.GenericResponseError(err), err
	}
	result.PoolId, result.ProfileId = opt.GetPoolId(), prf.Id

	// Update the file share data in database.
	log.V(5).Infof("file share creation result is %v", result)
	if err := db.C.UpdateStatus(ctx, result, model.FileShareAvailable); err != nil {
		return nil, err
	}

	return pb.GenericResponseResult(result), nil
}

// DeleteFileShare implements pb.ControllerServer.DeleteFileShare
func (c *Controller) DeleteFileShare(contx context.Context, opt *pb.DeleteFileShareOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive delete file share request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())

	dockInfo, err := db.C.GetDockByPoolId(ctx, opt.PoolId)
	if err != nil {
		log.Error("when search dock in db by pool id: ", err)
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareErrorDeleting)
		return pb.GenericResponseError(err), err
	}

	c.fileshareController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName

	if err = c.fileshareController.DeleteFileShare(opt); err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareErrorDeleting)
		return pb.GenericResponseError(err), err
	}

	if err = db.C.DeleteFileShare(ctx, opt.GetId()); err != nil {
		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(nil), err
}

// CreateFileShare implements pb.ControllerServer.CreateFileShareAcl
func (c *Controller) CreateFileShareAcl(contx context.Context, opt *pb.CreateFileShareAclOpts) (*pb.GenericResponse, error) {

	log.Info("controller server receive create file share acl request, vr =", opt)
	ctx := osdsCtx.NewContextFromJson(opt.GetContext())

	fileshare, err := db.C.GetFileShare(ctx, opt.FileshareId)
	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareError)
		return pb.GenericResponseError(err), err
	}

	dockInfo, err := db.C.GetDockByPoolId(ctx, fileshare.PoolId)
	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareError)
		log.Error("when search supported dock resource:", err.Error())
		return pb.GenericResponseError(err), err
	}
	c.fileshareController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName
	opt.Name = fileshare.Name

	result, err := c.fileshareController.CreateFileShareAcl((*pb.CreateFileShareAclOpts)(opt))
	if err != nil {
		// Change the status of the file share acl to error when the creation faild
		defer db.UpdateFileShareStatus(ctx, db.C, fileshare.Id, model.FileShareError)
		log.Error("when create file share acl:", err.Error())
		return pb.GenericResponseError(err), err
	}

	db.C.UpdateFileShareAcl(ctx, result)
	return pb.GenericResponseResult(result), nil
}

// DeleteFileShareAcl implements pb.ControllerServer.DeleteFileShareAcl
func (c *Controller) DeleteFileShareAcl(contx context.Context, opt *pb.DeleteFileShareAclOpts) (*pb.GenericResponse, error) {

	log.Info("controller server receive delete file share acl request, vr =", opt)
	ctx := osdsCtx.NewContextFromJson(opt.GetContext())

	fileshare, err := db.C.GetFileShare(ctx, opt.FileshareId)
	if err != nil {
		return pb.GenericResponseError(err), err
	}
	dockInfo, err := db.C.GetDockByPoolId(ctx, fileshare.PoolId)
	if err != nil {
		log.Error("when search supported dock resource:", err.Error())
		return pb.GenericResponseError(err), err
	}
	c.fileshareController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName
	opt.Name = fileshare.Name

	if err = c.fileshareController.DeleteFileShareAcl((*pb.DeleteFileShareAclOpts)(opt)); err != nil {
		// Change the status of the file share acl to error when the creation faild
		log.Error("when delete file share acl:", err.Error())
		return pb.GenericResponseError(err), err
	}
	if err = db.C.DeleteFileShareAcl(ctx, opt.Id); err != nil {
		log.Error("error occurred in controller module when delete file share acl in db: ", err)
		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(nil), nil
}

// CreateFileShareSnapshot implements pb.ControllerServer.CreateFileShareSnapshot
func (c *Controller) CreateFileShareSnapshot(contx context.Context, opt *pb.CreateFileShareSnapshotOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive create file share snapshot request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	if opt.Metadata == nil {
		opt.Metadata = map[string]string{}
	}

	profile := model.NewProfileFromJson(opt.Profile)
	if profile.SnapshotProperties.Topology.Bucket != "" {
		opt.Metadata["bucket"] = profile.SnapshotProperties.Topology.Bucket
	}

	fileshare, err := db.C.GetFileShare(ctx, opt.FileshareId)
	if err != nil {
		log.Error("get file share failed in create file share snapshot method: ", err)
		db.UpdateFileShareSnapshotStatus(ctx, db.C, opt.Id, model.FileShareSnapError)
		return pb.GenericResponseError(err), err
	}
	opt.Size = fileshare.Size
	dockInfo, err := db.C.GetDockByPoolId(ctx, fileshare.PoolId)
	if err != nil {
		log.Error("when search supported dock resource: ", err)
		db.UpdateFileShareSnapshotStatus(ctx, db.C, opt.Id, model.FileShareSnapError)
		return pb.GenericResponseError(err), err
	}
	c.fileshareController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName

	result, err := c.fileshareController.CreateFileShareSnapshot(opt)
	if err != nil {
		db.UpdateFileShareSnapshotStatus(ctx, db.C, opt.Id, model.FileShareSnapError)
		return pb.GenericResponseError(err), err
	}

	db.C.UpdateStatus(ctx, result, model.FileShareSnapAvailable)
	return pb.GenericResponseResult(result), nil
}

// DeleteFileshareSnapshot implements pb.ControllerServer.DeleteFileshareSnapshot
func (c *Controller) DeleteFileShareSnapshot(contx context.Context, opt *pb.DeleteFileShareSnapshotOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive delete file share snapshot request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	fileshare, err := db.C.GetFileShare(ctx, opt.FileshareId)
	if err != nil {
		log.Error("get file share failed in delete file share snapshot method: ", err)
		db.UpdateFileShareSnapshotStatus(ctx, db.C, opt.Id, model.FileShareSnapErrorDeleting)
		return pb.GenericResponseError(err), err
	}

	dockInfo, err := db.C.GetDockByPoolId(ctx, fileshare.PoolId)
	if err != nil {
		log.Error("when search supported dock resource: ", err)
		db.UpdateFileShareSnapshotStatus(ctx, db.C, opt.Id, model.FileShareSnapErrorDeleting)
		return pb.GenericResponseError(err), err
	}
	c.fileshareController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName

	if err = c.fileshareController.DeleteFileShareSnapshot(opt); err != nil {
		log.Error("error occurred in controller module when delete file share snapshot: ", err)
		db.UpdateFileShareSnapshotStatus(ctx, db.C, opt.Id, model.FileShareSnapErrorDeleting)
		return pb.GenericResponseError(err), err
	}
	if err = db.C.DeleteFileShareSnapshot(ctx, opt.Id); err != nil {
		log.Error("error occurred in controller module when delete file share snapshot in db: ", err)
		db.UpdateFileShareSnapshotStatus(ctx, db.C, opt.Id, model.FileShareSnapErrorDeleting)
		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(nil), nil
}

func (c *Controller) GetMetrics(context context.Context, opt *pb.GetMetricsOpts) (*pb.GenericResponse, error) {
	log.Info("in controller get metrics methods")

	var result []*model.MetricSpec
	var err error

	if opt.StartTime == "" && opt.EndTime == "" {
		// no start and end time specified, get the latest value of this metric
		result, err = c.metricsController.GetLatestMetrics(opt)
	} else if opt.StartTime == opt.EndTime {
		// same start and end time specified, get the value of this metric at that timestamp
		result, err = c.metricsController.GetInstantMetrics(opt)
	} else {
		// range of start and end time is specified
		result, err = c.metricsController.GetRangeMetrics(opt)
	}

	if err != nil {
		log.Errorf("get metrics failed: %s\n", err.Error())

		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(result), err
}

func (c *Controller) CollectMetrics(context context.Context, opt *pb.CollectMetricsOpts) (*pb.GenericResponse, error) {
	log.V(5).Info("in controller collect metrics methods")

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	dockSpec, err := db.C.ListDocks(ctx)
	if err != nil {
		log.Errorf("list dock  failed in CollectMetrics method: %s", err.Error())
		return pb.GenericResponseError(err), err
	}
	for i, d := range dockSpec {
		if d.DriverName == opt.DriverName {
			log.Infof("driver found driver: %s", d.DriverName)
			dockInfo, err := db.C.GetDock(ctx, dockSpec[i].BaseModel.Id)
			if err != nil {
				log.Errorf("error %s when search dock in db by dock id: %s", err.Error(), dockSpec[i].BaseModel.Id)
				return pb.GenericResponseError(err), err

			}
			c.metricsController.SetDock(dockInfo)
			result, err := c.metricsController.CollectMetrics(opt)
			if err != nil {
				log.Errorf("collectMetrics failed: %s", err.Error())

				return pb.GenericResponseError(err), err
			}

			return pb.GenericResponseResult(result), nil
		}
	}
	return nil, nil
}

func (c *Controller) GetUrls(context.Context, *pb.NoParams) (*pb.GenericResponse, error) {
	log.V(5).Info("in controller get urls method")

	var result *map[string]model.UrlDesc
	var err error

	result, err = c.metricsController.GetUrls()

	// make return array
	arrUrls := make([]model.UrlSpec, 0)

	for k, v := range *result {
		// make each url spec
		urlSpec := model.UrlSpec{}
		urlSpec.Name = k
		urlSpec.Url = v.Url
		urlSpec.Desc = v.Desc
		// add to the array
		arrUrls = append(arrUrls, urlSpec)
	}

	if err != nil {
		log.Errorf("get urls failed: %s\n", err.Error())
		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(arrUrls), err
}
