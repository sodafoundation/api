/*
 * This source file has been modified by Huawei Technologies Co., Ltd.
 * Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
 *
 * Copyright 2015, Google Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *     * Neither the name of Google Inc. nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

/*
This module implements the entry into operations of storageDock module.

*/

package server

import (
	"encoding/json"
	"fmt"
	"net"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/dock"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// dockServer is used to implement pb.DockServer
type dockServer struct {
	Port string
}

// NewDockServer returns a dockServer instance.
func NewDockServer(port string) *dockServer {
	return &dockServer{
		Port: port,
	}
}

func (ds *dockServer) Run() error {
	// New Grpc Server
	s := grpc.NewServer()
	// Register dock service.
	pb.RegisterProvisionDockServer(s, ds)
	pb.RegisterAttachDockServer(s, ds)

	// Listen the dock server port.
	lis, err := net.Listen("tcp", ds.Port)
	if err != nil {
		log.Fatalf("failed to listen: %+v", err)
		return err
	}

	log.Info("Dock server initialized! Start listening on port:", lis.Addr())

	// Start dock server watching loop.
	defer s.Stop()
	return s.Serve(lis)
}

// CreateVolume implements pb.DockServer.CreateVolume
func (ds *dockServer) CreateVolume(ctx context.Context, opt *pb.CreateVolumeOpts) (*pb.GenericResponse, error) {
	var res pb.GenericResponse

	log.Info("Dock server receive create volume request, vr =", opt)

	vol, err := dock.Brain.CreateVolume(opt)
	if err != nil {
		log.Error("When create volume in dock module:", err)

		res.Reply = GenericResponseError("400", fmt.Sprint(err))
		return &res, err
	}

	res.Reply = GenericResponseResult(vol)
	return &res, nil
}

// DeleteVolume implements pb.DockServer.DeleteVolume
func (ds *dockServer) DeleteVolume(ctx context.Context, opt *pb.DeleteVolumeOpts) (*pb.GenericResponse, error) {
	var res pb.GenericResponse

	log.Info("Dock server receive delete volume request, vr =", opt)

	if err := dock.Brain.DeleteVolume(opt); err != nil {
		log.Error("Error occurred in dock module when delete volume:", err)

		res.Reply = GenericResponseError("400", fmt.Sprint(err))
		return &res, err
	}

	res.Reply = GenericResponseResult("")
	return &res, nil
}

// ExtendVolume implements pb.DockServer.ExtendVolume
func (ds *dockServer) ExtendVolume(ctx context.Context, opt *pb.ExtendVolumeOpts) (*pb.GenericResponse, error) {
	var res pb.GenericResponse

	log.Info("Dock server receive extend volume request, vr =", opt)

	vol, err := dock.Brain.ExtendVolume(opt)
	if err != nil {
		log.Error("When extend volume in dock module:", err)

		res.Reply = GenericResponseError("400", fmt.Sprint(err))
		return &res, err
	}

	res.Reply = GenericResponseResult(vol)
	return &res, nil
}

// CreateAttachment implements pb.DockServer.CreateAttachment
func (ds *dockServer) CreateAttachment(ctx context.Context, opt *pb.CreateAttachmentOpts) (*pb.GenericResponse, error) {
	var res pb.GenericResponse

	log.Info("Dock server receive create volume attachment request, vr =", opt)

	atc, err := dock.Brain.CreateVolumeAttachment(opt)
	if err != nil {
		log.Error("Error occurred in dock module when create volume attachment:", err)

		res.Reply = GenericResponseError("400", fmt.Sprint(err))
		return &res, err
	}

	res.Reply = GenericResponseResult(atc)
	return &res, nil
}

// DeleteAttachment implements pb.DockServer.DeleteAttachment
func (ds *dockServer) DeleteAttachment(ctx context.Context, opt *pb.DeleteAttachmentOpts) (*pb.GenericResponse, error) {
	var res pb.GenericResponse

	log.Info("Dock server receive delete volume attachment request, vr =", opt)

	if err := dock.Brain.DeleteVolumeAttachment(opt); err != nil {
		log.Error("Error occurred in dock module when delete volume attachment:", err)

		res.Reply = GenericResponseError("400", fmt.Sprint(err))
		return &res, err
	}

	res.Reply = GenericResponseResult("")
	return &res, nil
}

// CreateVolumeSnapshot implements pb.DockServer.CreateVolumeSnapshot
func (ds *dockServer) CreateVolumeSnapshot(ctx context.Context, opt *pb.CreateVolumeSnapshotOpts) (*pb.GenericResponse, error) {
	var res pb.GenericResponse

	log.Info("Dock server receive create volume snapshot request, vr =", opt)

	snp, err := dock.Brain.CreateSnapshot(opt)
	if err != nil {
		log.Error("Error occurred in dock module when create snapshot:", err)
		res.Reply = GenericResponseError("400", fmt.Sprint(err))
		return &res, err
	}

	res.Reply = GenericResponseResult(snp)
	return &res, nil
}

// DeleteVolumeSnapshot implements pb.DockServer.DeleteVolumeSnapshot
func (ds *dockServer) DeleteVolumeSnapshot(ctx context.Context, opt *pb.DeleteVolumeSnapshotOpts) (*pb.GenericResponse, error) {
	var res pb.GenericResponse

	log.Info("Dock server receive delete volume snapshot request, vr =", opt)

	if err := dock.Brain.DeleteSnapshot(opt); err != nil {
		log.Error("Error occurred in dock module when delete snapshot:", err)

		res.Reply = GenericResponseError("400", fmt.Sprint(err))
		return &res, err
	}

	res.Reply = GenericResponseResult("")
	return &res, nil
}

// AttachVolume implements pb.DockServer.AttachVolume
func (ds *dockServer) AttachVolume(ctx context.Context, opt *pb.AttachVolumeOpts) (*pb.GenericResponse, error) {
	var res pb.GenericResponse

	log.Info("Dock server receive attach volume request, vr =", opt)

	atc, err := dock.Brain.AttachVolume(opt)
	if err != nil {
		log.Error("Error occurred in dock module when attach volume:", err)

		res.Reply = GenericResponseError("400", fmt.Sprint(err))
		return &res, err
	}

	res.Reply = GenericResponseResult(atc)
	return &res, nil
}

// DetachVolume implements pb.DockServer.DetachVolume
func (ds *dockServer) DetachVolume(ctx context.Context, opt *pb.DetachVolumeOpts) (*pb.GenericResponse, error) {
	var res pb.GenericResponse

	log.Info("Dock server receive detach volume request, vr =", opt)

	if err := dock.Brain.DetachVolume(opt); err != nil {
		log.Error("Error occurred in dock module when detach volume:", err)

		res.Reply = GenericResponseError("400", fmt.Sprint(err))
		return &res, err
	}

	res.Reply = GenericResponseResult("")
	return &res, nil
}

func GenericResponseResult(message interface{}) *pb.GenericResponse_Result_ {
	var msg string
	switch message.(type) {
	case string:
		msg = message.(string)
	default:
		msgJSON, _ := json.Marshal(message)
		msg = string(msgJSON)
	}

	return &pb.GenericResponse_Result_{
		Result: &pb.GenericResponse_Result{
			Message: msg,
		},
	}
}

func GenericResponseError(code, description string) *pb.GenericResponse_Error_ {
	return &pb.GenericResponse_Error_{
		Error: &pb.GenericResponse_Error{
			Code:        code,
			Description: description,
		},
	}
}
