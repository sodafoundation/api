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

package fake_server

import (
	"log"
	"net"

	dockApi "github.com/opensds/opensds/testing/pkg/dock/api"
	pb "github.com/opensds/opensds/testing/pkg/grpc/fake_opensds"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// dockServer is used to implement opensds.DockServer.
type dockServer struct {
	Port string `json:"port"`
}

// NewDockServer returns an dockServer instance.
func NewDockServer(port string) *dockServer {
	return &dockServer{
		Port: port,
	}
}

// CreateVolume implements opensds.DockServer
func (ds *dockServer) CreateVolume(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Dock server receive create volume request, vr =", in)
	return dockApi.CreateVolume(in)
}

// GetVolume implements opensds.DockServer
func (ds *dockServer) GetVolume(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Dock server receive get volume request, vr =", in)
	return dockApi.GetVolume(in)
}

// ListVolumes implements opensds.DockServer
func (ds *dockServer) ListVolumes(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Dock server receive list volumes request, vr =", in)
	return dockApi.ListVolumes(in)
}

// DeleteVolume implements opensds.DockServer
func (ds *dockServer) DeleteVolume(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Dock server receive delete volume request, vr =", in)
	return dockApi.DeleteVolume(in)
}

// AttachVolume implements opensds.DockServer
func (ds *dockServer) AttachVolume(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Dock server receive attach volume request, vr =", in)
	return dockApi.AttachVolume(in)
}

// DetachVolume implements opensds.DockServer
func (ds *dockServer) DetachVolume(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Dock server receive detach volume request, vr =", in)
	return dockApi.DetachVolume(in)
}

// MountVolume implements opensds.DockServer
func (ds *dockServer) MountVolume(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Dock server receive mount volume request, vr =", in)
	return dockApi.MountVolume(in)
}

// UnmountVolume implements opensds.DockServer
func (ds *dockServer) UnmountVolume(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Dock server receive unmount volume request, vr =", in)
	return dockApi.UnmountVolume(in)
}

// CreateShare implements opensds.DockServer
func (ds *dockServer) CreateShare(ctx context.Context, in *pb.ShareRequest) (*pb.Response, error) {
	log.Println("Dock server receive create share request, sr =", in)
	return dockApi.CreateShare(in)
}

// GetShare implements opensds.DockServer
func (ds *dockServer) GetShare(ctx context.Context, in *pb.ShareRequest) (*pb.Response, error) {
	log.Println("Dock server receive get share request, sr =", in)
	return dockApi.GetShare(in)
}

// ListShares implements opensds.DockServer
func (ds *dockServer) ListShares(ctx context.Context, in *pb.ShareRequest) (*pb.Response, error) {
	log.Println("Dock server receive list shares request, sr =", in)
	return dockApi.ListShares(in)
}

// DeleteShare implements opensds.DockServer
func (ds *dockServer) DeleteShare(ctx context.Context, in *pb.ShareRequest) (*pb.Response, error) {
	log.Println("Dock server receive delete share request, sr =", in)
	return dockApi.DeleteShare(in)
}

func (ds *dockServer) ListenAndServe() {
	lis, err := net.Listen("tcp", ds.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Dock server initialized! Start listening on port:", ds.Port)

	gs := grpc.NewServer()
	pb.RegisterDockServer(gs, ds)
	gs.Serve(lis)

	defer gs.Stop()
}
