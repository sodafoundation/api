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

package server

import (
	"log"
	"net"

	dockApi "github.com/opensds/opensds/pkg/dock/api"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
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
func (ds *dockServer) CreateVolume(ctx context.Context, in *pb.DockRequest) (*pb.DockResponse, error) {
	log.Println("Dock server receive create volume request, vr =", in)
	return dockApi.CreateVolume(in)
}

// GetVolume implements opensds.DockServer
func (ds *dockServer) GetVolume(ctx context.Context, in *pb.DockRequest) (*pb.DockResponse, error) {
	log.Println("Dock server receive get volume request, vr =", in)
	return dockApi.GetVolume(in)
}

// DeleteVolume implements opensds.DockServer
func (ds *dockServer) DeleteVolume(ctx context.Context, in *pb.DockRequest) (*pb.DockResponse, error) {
	log.Println("Dock server receive delete volume request, vr =", in)
	return dockApi.DeleteVolume(in)
}

// CreateVolumeAttachment implements opensds.DockServer
func (ds *dockServer) CreateVolumeAttachment(ctx context.Context, in *pb.DockRequest) (*pb.DockResponse, error) {
	log.Println("Dock server receive create volume attachment request, vr =", in)
	return dockApi.CreateVolumeAttachment(in)
}

// UpdateVolumeAttachment implements opensds.DockServer
func (ds *dockServer) UpdateVolumeAttachment(ctx context.Context, in *pb.DockRequest) (*pb.DockResponse, error) {
	log.Println("Dock server receive update volume attachment request, vr =", in)
	return dockApi.UpdateVolumeAttachment(in)
}

// DeleteVolumeAttachment implements opensds.DockServer
func (ds *dockServer) DeleteVolumeAttachment(ctx context.Context, in *pb.DockRequest) (*pb.DockResponse, error) {
	log.Println("Dock server receive delete volume attachment request, vr =", in)
	return dockApi.DeleteVolumeAttachment(in)
}

// CreateVolumeSnapshot implements opensds.DockServer
func (ds *dockServer) CreateVolumeSnapshot(ctx context.Context, in *pb.DockRequest) (*pb.DockResponse, error) {
	log.Println("Dock server receive create volume snapshot request, vr =", in)
	return dockApi.CreateVolumeSnapshot(in)
}

// GetVolumeSnapshot implements opensds.DockServer
func (ds *dockServer) GetVolumeSnapshot(ctx context.Context, in *pb.DockRequest) (*pb.DockResponse, error) {
	log.Println("Dock server receive get volume snapshot request, vr =", in)
	return dockApi.GetVolumeSnapshot(in)
}

// DeleteVolumeSnapshot implements opensds.DockServer
func (ds *dockServer) DeleteVolumeSnapshot(ctx context.Context, in *pb.DockRequest) (*pb.DockResponse, error) {
	log.Println("Dock server receive delete volume snapshot request, vr =", in)
	return dockApi.DeleteVolumeSnapshot(in)
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

		log.Println("Dock server initialized! Start listening on port:", ds.Port)

		// Start dock server watching loop.
		ds.Server.Serve(lis)

		defer ds.Server.Stop()
	default:
		log.Fatalln("Don't support this type!")
		return
	}
}
