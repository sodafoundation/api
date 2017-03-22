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

	orchApi "github.com/opensds/opensds/pkg/controller/orchestration/api"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// orchServer is used to implement opensds.OrchesrationServer.
type orchServer struct {
	Port string `json:"port"`
}

// NewOrchServer returns an orchServer instance.
func NewOrchServer(port string) *orchServer {
	return &orchServer{
		Port: port,
	}
}

// CreateVolume implements opensds.OrchestrationServer
func (os *orchServer) CreateVolume(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Orchestration server receive create volume request, vr =", in)
	return orchApi.CreateVolume(in)
}

// GetVolume implements opensds.OrchestrationServer
func (os *orchServer) GetVolume(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Orchestration server receive get volume request, vr =", in)
	return orchApi.GetVolume(in)
}

// ListVolumes implements opensds.OrchestrationServer
func (os *orchServer) ListVolumes(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Orchestration server receive list volumes request, vr =", in)
	return orchApi.ListVolumes(in)
}

// DeleteVolume implements opensds.OrchestrationServer
func (os *orchServer) DeleteVolume(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Orchestration server receive delete volume request, vr =", in)
	return orchApi.DeleteVolume(in)
}

// AttachVolume implements opensds.OrchestrationServer
func (os *orchServer) AttachVolume(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Orchestration server receive attach volume request, vr =", in)
	return orchApi.AttachVolume(in)
}

// DetachVolume implements opensds.OrchestrationServer
func (os *orchServer) DetachVolume(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Orchestration server receive detach volume request, vr =", in)
	return orchApi.DetachVolume(in)
}

// MountVolume implements opensds.OrchestrationServer
func (os *orchServer) MountVolume(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Orchestration server receive mount volume request, vr =", in)
	return orchApi.MountVolume(in)
}

// UnmountVolume implements opensds.OrchestrationServer
func (os *orchServer) UnmountVolume(ctx context.Context, in *pb.VolumeRequest) (*pb.Response, error) {
	log.Println("Orchestration server receive unmount volume request, vr =", in)
	return orchApi.UnmountVolume(in)
}

// CreateShare implements opensds.OrchestrationServer
func (os *orchServer) CreateShare(ctx context.Context, in *pb.ShareRequest) (*pb.Response, error) {
	log.Println("Orchestration server receive create share request, sr =", in)
	return orchApi.CreateShare(in)
}

// GetShare implements opensds.OrchestrationServer
func (os *orchServer) GetShare(ctx context.Context, in *pb.ShareRequest) (*pb.Response, error) {
	log.Println("Orchestration server receive get share request, sr =", in)
	return orchApi.GetShare(in)
}

// ListShares implements opensds.OrchestrationServer
func (os *orchServer) ListShares(ctx context.Context, in *pb.ShareRequest) (*pb.Response, error) {
	log.Println("Orchestration server receive list shares request, sr =", in)
	return orchApi.ListShares(in)
}

// DeleteShare implements opensds.OrchestrationServer
func (os *orchServer) DeleteShare(ctx context.Context, in *pb.ShareRequest) (*pb.Response, error) {
	log.Println("Orchestration server receive delete share request, sr =", in)
	return orchApi.DeleteShare(in)
}

func (os *orchServer) ListenAndServe() {
	lis, err := net.Listen("tcp", os.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Orchestration server initialized! Start listening on port:", os.Port)

	gs := grpc.NewServer()
	pb.RegisterOrchestrationServer(gs, os)
	gs.Serve(lis)

	defer gs.Stop()
}
