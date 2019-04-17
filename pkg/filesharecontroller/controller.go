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
	_ "encoding/json"
	"github.com/opensds/opensds/pkg/filesharecontroller/fileshare"
	"fmt"

	"net"

	log "github.com/golang/glog"
	osdsCtx "github.com/opensds/opensds/pkg/context"

	"github.com/opensds/opensds/pkg/filesharecontroller/selector"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/fileshareproto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	fmt.Sprintf("router calls respective controller function here it is CreateFileShare")
	var err error
	var prf *model.ProfileSpec
	//var snap *model.FileShareSnapshotSpec
	//var snapVol *model.VolumeSpec

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
		db.UpdateFileShareStatus(ctx, db.C, model.FileShareError)
		log.Error("get profile failed: ", err)
		return pb.GenericResponseError(err), err
	}

	// This file share structure is currently fetched from database, but eventually
	// it will be removed after SelectSupportedPoolForFileShare method in selector
	// is updated.
	//db.C.GetVolume(ctx, opt.Id)
	//fileshare, err := db.C.ListFileShare(ctx, opt.Id)
	fileshare, err := db.C.GetFileShare(ctx, opt.Id)
	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, model.FileShareError)
		return pb.GenericResponseError(err), err
	}
	polInfo, err := c.selector.SelectSupportedPoolForFileShare(fileshare)
	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, model.FileShareError)
		return pb.GenericResponseError(err), err
	}
	// whether specify a pool or not, opt's poolid and pool name should be
	// assigned by polInfo
	opt.PoolId = polInfo.Id
	opt.PoolName = polInfo.Name

	dockInfo, err := db.C.GetDock(ctx, polInfo.DockId)
	if err != nil {
		db.UpdateFileShareStatus(ctx, db.C, model.FileShareError)
		log.Error("when search supported dock resource:", err.Error())
		return pb.GenericResponseError(err), err
	}
	c.fileshareController.SetDock(dockInfo)
	opt.DriverName = dockInfo.DriverName

	result, err := c.fileshareController.CreateFileShare((*pb.CreateFileShareOpts)(opt))
	if err != nil {
		// Change the status of the file share to error when the creation faild
		defer db.UpdateFileShareStatus(ctx, db.C, model.FileShareError)
		log.Error("when create file share:", err.Error())
		return pb.GenericResponseError(err), err
	}
	result.PoolId, result.ProfileId = opt.GetPoolId(), opt.GetProfileId()

	// Update the file share data in database.
	db.C.UpdateStatus(ctx, result, model.FileShareAvailable)

	
	return pb.GenericResponseResult(result), nil
}
