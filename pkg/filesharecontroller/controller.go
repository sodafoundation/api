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
	"github.com/opensds/opensds/pkg/filesharecontroller/fileshare"
	"net"
	log "github.com/golang/glog"
	osdsCtx "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/filesharecontroller/selector"
	"github.com/opensds/opensds/pkg/filesharecontroller/policy"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/fileshareproto"
	"golang.org/x/net/context"
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
	fileshareCtrl := fileshare.NewController()
	return &Controller{
		selector:         selector.NewSelector(),
		fileshareController: fileshareCtrl,
		Port:             port,
	}
}

type Controller struct {
	selector         selector.Selector
	fileshareController fileshare.Controller
	policyController policy.Controller
	Port string
}


// Run method would start the listen mechanism of controller module.
func (c *Controller) Run() error {
	// New Grpc Server
	s := grpc.NewServer()
	// Register controller service.
	pb.RegisterControllerServer(s, c)

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

// CreateFileShare implements pb.ControllerServer.CreateFileShare
func (c *Controller) CreateFileShare(contx context.Context, opt *pb.CreateFileShareOpts) (*pb.GenericResponse, error) {
	var err error
	var prf *model.ProfileSpec


	log.Info("Controller server receive create fileshare request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	if opt.ProfileId == "" {
		log.Warning("Use default profile when user doesn't specify profile.")
		prf, err = db.C.GetDefaultProfile(ctx)
		opt.ProfileId = prf.Id
	} else {
		prf, err = db.C.GetProfile(ctx, opt.ProfileId)
	}
	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareError)
		log.Error("get profile failed: ", err)
		return pb.GenericResponseError(err), err
	}

	// This fileshare structure is currently fetched from database, but eventually
	// it will be removed after SelectSupportedPoolForVolume method in selector
	// is updated.
	vol, err := db.C.GetFileShare(ctx, opt.Id)
	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareError)
		return pb.GenericResponseError(err), err
	}
	polInfo, err := c.selector.SelectSupportedPoolForFileShare(vol)
	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareError)
		return pb.GenericResponseError(err), err
	}
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

	result, err := c.fileshareController.CreateFileShare(opt)
	if err != nil {
		// Change the status of the file share to error when the creation faild
		defer db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareError)
		log.Error("when create file share:", err.Error())
		return pb.GenericResponseError(err), err
	}
	result.PoolId, result.ProfileId = opt.GetPoolId(), opt.GetProfileId()

	// Update the fileshare data in database.
	db.C.UpdateStatus(ctx, result, model.FileShareAvailable)

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

// DeleteFileShare implements pb.ControllerServer.DeleteFileShare
func (c *Controller) DeleteFileShare(contx context.Context, opt *pb.DeleteFileShareOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive delete file share request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	prf, err := db.C.GetProfile(ctx, opt.ProfileId)
	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareErrorDeleting)
		log.Error("when search profile in db:", err)
		return pb.GenericResponseError(err), err
	}

	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(prf)
	c.policyController.Setup(DELETE_LIFECIRCLE_FLAG)

	dockInfo, err := db.C.GetDockByPoolId(ctx, opt.PoolId)
	if err != nil {
		log.Error("when search dock in db by pool id: ", err)
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareErrorDeleting)
		return pb.GenericResponseError(err), err
	}
	c.policyController.SetDock(dockInfo)
	c.fileshareController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName

	var errChan = make(chan error, 1)
	defer close(errChan)
	go c.policyController.ExecuteAsyncPolicy(opt, "", errChan)

	if err := <-errChan; err != nil {
		log.Error("when execute async policy: ", err)
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareErrorDeleting)
		return pb.GenericResponseError(err), err
	}

	if err = c.fileshareController.DeleteFileShare(opt); err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareErrorDeleting)
		return pb.GenericResponseError(err), err
	}
	if err = db.C.DeleteFileShare(ctx, opt.GetId()); err != nil {
		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(nil), nil
}

/*
// CreateFileShare implements pb.ControllerServer.CreateFileShare
func (c *Controller) CreateFileShare(contx context.Context, opt *pb.CreateFileShareOpts) (*pb.GenericResponse, error) {
	fmt.Sprintf("router calls respective controller function here it is CreateFileShare")
	var err error
	var prf *model.ProfileSpec

	log.Info("Controller server receive create file share request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	if opt.ProfileId == "" {
		log.Warning("Use default profile when user doesn't specify profile.")
		prf, err = db.C.GetDefaultProfile(ctx)
		opt.ProfileId = prf.Id
	} else {
		prf, err = db.C.GetProfile(ctx, opt.ProfileId)
	}
	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareError)
		log.Error("get profile failed: ", err)
		return pb.GenericResponseError(err), err
	}

	// This file share structure is currently fetched from database, but eventually
	// it will be removed after SelectSupportedPoolForFileShare method in selector
	// is updated.
	fileshare, err := db.C.GetFileShare(ctx, opt.Id)
	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareError)
		return pb.GenericResponseError(err), err
	}
	polInfo, err := c.selector.SelectSupportedPoolForFileShare(fileshare)
	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareError)
		return pb.GenericResponseError(err), err
	}
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

	result, err := c.fileshareController.CreateFileShare((*pb.CreateFileShareOpts)(opt))
	if err != nil {
		// Change the status of the file share to error when the creation faild
		defer db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareError)
		log.Error("when create file share:", err.Error())
		return pb.GenericResponseError(err), err
	}
	result.PoolId, result.ProfileId = opt.GetPoolId(), opt.GetProfileId()

	// Update the file share data in database.
	db.C.UpdateStatus(ctx, result, model.FileShareAvailable)

	
	return pb.GenericResponseResult(result), nil
}

// DeleteFileShare implements pb.ControllerServer.DeleteFileShare
func (c *Controller) DeleteFileShare(contx context.Context, opt *pb.DeleteFileShareOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive delete file share request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	prf, err := db.C.GetProfile(ctx, opt.ProfileId)
	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareErrorDeleting)
		log.Error("when search profile in db:", err)
		return pb.GenericResponseError(err), err
	}


	// Select the storage tag according to the lifecycle flag.
	c.policyController = policy.NewController(prf)
	c.policyController.Setup(DELETE_LIFECIRCLE_FLAG)

	dockInfo, err := db.C.GetDockByPoolId(ctx, opt.PoolId)
	if err != nil {
		log.Error("when search dock in db by pool id: ", err)
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareErrorDeleting)
		return pb.GenericResponseError(err), err
	}

	c.fileshareController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName


	var errChan = make(chan error, 1)
	defer close(errChan)
	go c.policyController.ExecuteAsyncPolicy(opt, "", errChan)

	if err := <-errChan; err != nil {
		log.Error("when execute async policy: ", err)
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareErrorDeleting)
		return pb.GenericResponseError(err), err
	}

	if err = c.fileshareController.DeleteFileShare(opt); err != nil {
		db.UpdateFileShareStatus(ctx, db.C, opt.Id, model.FileShareErrorDeleting)
		return pb.GenericResponseError(err), err
	}
	if err = db.C.DeleteFileShare(ctx, opt.GetId()); err != nil {
		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(nil), nil
}
*/

// UpdateVolumeGroup implements pb.ControllerServer.UpdateVolumeGroup
func (c *Controller) UpdateFileShare(contx context.Context, opt *pb.UpdateFileShareOpts) (*pb.GenericResponse, error) {

	log.Info("Controller server receive update volume group request, vr =", opt)

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	dock, err := db.C.GetDockByPoolId(ctx, opt.PoolId)
	if err != nil {
		db.UpdateVolumeGroupStatus(ctx, db.C, opt.Id, model.FileShareError)
		return pb.GenericResponseError(err), err
	}
	c.fileshareController.SetDock(dock)
	opt.DriverName = dock.DriverName

	vg, err := c.fileshareController.UpdateFileShare(opt)
	if err != nil {
		log.Error("when create volume group: ", err)
		db.UpdateVolumeGroupStatus(ctx, db.C, opt.Id, model.FileShareError)

		//for _, addVol := range opt.AddVolumes {
		//	db.UpdateVolumeStatus(ctx, db.C, addVol, model.VolumeError)
		//}
		//for _, rmVol := range opt.RemoveVolumes {
		//	db.UpdateVolumeStatus(ctx, db.C, rmVol, model.VolumeError)
		//}

		return pb.GenericResponseError(err), err
	}

	/*
	// Update group id in the volumes
	for _, addVolId := range opt.AddFileShares {
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
	*/
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
	db.C.UpdateStatus(ctx, vg, model.FileShareAvailable)
	return pb.GenericResponseResult(vg), nil
}

