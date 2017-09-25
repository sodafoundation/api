/*
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
	"net"

	log "github.com/golang/glog"

	"github.com/opensds/opensds/pkg/dock"
	pb "github.com/opensds/opensds/pkg/dock/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// dockServer is used to implement opensds.DockServer.
type dockServer struct {
	Server *grpc.Server
	Port   string
}

// NewDockServer returns an dockServer instance.
func NewDockServer(port string) pb.DockServer {
	// Construct dock server.
	gs := grpc.NewServer()
	ds := &dockServer{
		Server: gs,
		Port:   port,
	}

	// Register dock server.
	pb.RegisterDockServer(gs, ds)

	return ds
}

// CreateVolume implements opensds.DockServer
func (ds *dockServer) CreateVolume(ctx context.Context, opt *pb.CreateVolumeOpts) (*pb.GenericResponse, error) {
	log.Info("Dock server receive create volume request, vr =", opt)

	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Error("When parsing dock info:", err)
		return &pb.DockResponse{}, err
	}

	vol, err := dock.NewDockHub(dck.GetDriverName()).CreateVolume(
		req.GetVolumeName(),
		req.GetVolumeSize())
	if err != nil {
		log.Error("When create volume in dock module:", err)
		return &pb.DockResponse{}, err
	}

	volBody, _ := json.Marshal(result)
	return &pb.DockResponse{
		Status:  "Success",
		Message: string(volBody),
	}, nil
}

// DeleteVolume implements opensds.DockServer
func (ds *dockServer) DeleteVolume(ctx context.Context, req *pb.DockRequest) (*pb.DockResponse, error) {
	log.Info("Dock server receive delete volume request, vr =", req)

	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Error("When parsing dock info:", err)
	}

	if err := dock.NewDockHub(dck.GetDriverName()).DeleteVolume(req.GetVolumeId()); err != nil {
		log.Error("Error occured in dock module when delete volume:", err)
		return &pb.DockResponse{}, err
	}

	return &pb.DockResponse{
		Status:  "Success",
		Message: "Delete volume success",
	}, nil
}

// CreateVolumeAttachment implements opensds.DockServer
func (ds *dockServer) CreateVolumeAttachment(ctx context.Context, req *pb.DockRequest) (*pb.DockResponse, error) {
	log.Info("Dock server receive create volume attachment request, vr =", req)

	var dck, hostInfo = &api.DockSpec{}, &api.HostInfo{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Error("When parsing dock info:", err)
	}
	if err := json.Unmarshal([]byte(req.GetHostInfo()), hostInfo); err != nil {
		log.Error("Error occured in dock module when parsing host info:", err)
		return &pb.DockResponse{}, err
	}

	atc, err := dock.NewDockHub(dck.GetDriverName()).CreateVolumeAttachment(
		req.GetVolumeId(),
		req.GetDoLocalAttach(),
		req.GetMultiPath(),
		hostInfo)
	if err != nil {
		log.Error("Error occured in dock module when create volume attachment:", err)
		return &pb.DockResponse{}, err
	}

	atcBody, _ := json.Marshal(result)
	return &pb.DockResponse{
		Status:  "Success",
		Message: string(atcBody),
	}, nil
}

// CreateVolumeSnapshot implements opensds.DockServer
func (ds *dockServer) CreateVolumeSnapshot(ctx context.Context, req *pb.DockRequest) (*pb.DockResponse, error) {
	log.Info("Dock server receive create volume snapshot request, vr =", req)

	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Error("When parsing dock info:", err)
	}

	snp, err := dock.NewDockHub(dck.GetDriverName()).CreateSnapshot(
		req.GetSnapshotName(),
		req.GetVolumeId(),
		req.GetSnapshotDescription())
	if err != nil {
		log.Error("Error occured in dock module when create snapshot:", err)
		return &pb.DockResponse{}, err
	}

	snpBody, _ := json.Marshal(result)
	return &pb.DockResponse{
		Status:  "Success",
		Message: string(snpBody),
	}, nil
}

// DeleteVolumeSnapshot implements opensds.DockServer
func (ds *dockServer) DeleteVolumeSnapshot(ctx context.Context, req *pb.DockRequest) (*pb.DockResponse, error) {
	log.Info("Dock server receive delete volume snapshot request, vr =", req)

	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Error("When parsing dock info:", err)
	}

	if err := dock.NewDockHub(dck.GetDriverName()).DeleteSnapshot(req.GetSnapshotId()); err != nil {
		log.Error("Error occured in dock module when delete snapshot:", err)
		return &pb.DockResponse{}, err
	}

	return &pb.DockResponse{
		Status:  "Success",
		Message: "Delete snapshot success",
	}, nil
}

func ListenAndServe(srv pb.DockServer) {
	// Find whether the type of input is supported.
	switch srv.(type) {
	case *dockServer:
		ds := srv.(*dockServer)

		// Listen the dock server port.
		lis, err := net.Listen("tcp", ds.Port)
		if err != nil {
			log.Fatalf("failed to listen: %+v", err)
			return
		}

		log.Info("Dock server initialized! Start listening on port:", ds.Port)

		// Start dock server watching loop.
		ds.Server.Serve(lis)

		defer ds.Server.Stop()
	default:
		log.Fatalln("Don't support this type!")
		return
	}
}
