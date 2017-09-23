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
	log "github.com/golang/glog"
	"net"
	"strings"

	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/dock"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	api "github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"

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
func (ds *dockServer) CreateVolume(ctx context.Context, req *pb.DockRequest) (*pb.DockResponse, error) {
	log.Info("Dock server receive create volume request, vr =", req)

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

	// If volume uuid is null, generate it randomly.
	if vol.GetId() == "" {
		if ok := utils.NewSetter().SetUuid(vol); ok != nil {
			log.Error("When set volume uuid:", ok)
			return &pb.DockResponse{}, err
		}
	}

	// Set volume created time.
	if ok := utils.NewSetter().SetCreatedTimeStamp(vol); ok != nil {
		log.Error("When set volume created time:", ok)
		return &pb.DockResponse{}, err
	}

	vol.ProfileId = req.GetProfileId()
	vol.PoolId = req.GetPoolId()

	result, err := db.C.CreateVolume(vol)
	if err != nil {
		log.Error("When create volume in db module:", err)
		return &pb.DockResponse{}, err
	}

	volBody, _ := json.Marshal(result)
	return &pb.DockResponse{
		Status:  "Success",
		Message: string(volBody),
	}, nil
}

// GetVolume implements opensds.DockServer
func (ds *dockServer) GetVolume(ctx context.Context, req *pb.DockRequest) (*pb.DockResponse, error) {
	log.Info("Dock server receive get volume request, vr =", req)

	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Error("When parsing dock info:", err)
	}

	result, err := dock.NewDockHub(dck.GetDriverName()).GetVolume(req.GetVolumeId())
	if err != nil {
		log.Error("When get volume in dock module:", err)
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

	if err := db.C.DeleteVolume(req.GetVolumeId()); err != nil {
		log.Error("Error occured in dock module when delete volume in db:", err)
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

	// If volume attachment uuid is null, generate it randomly.
	if atc.GetId() == "" {
		if ok := utils.NewSetter().SetUuid(atc); ok != nil {
			log.Error("When set volume attachment uuid:", ok)
			return &pb.DockResponse{}, err
		}
	}

	// Set volume attachment created time.
	if ok := utils.NewSetter().SetCreatedTimeStamp(atc); ok != nil {
		log.Error("When set volume attachment created time:", ok)
		return &pb.DockResponse{}, err
	}

	result, err := db.C.CreateVolumeAttachment(req.GetVolumeId(), atc)
	if err != nil {
		log.Error("Error occured in dock module when create volume attachment in db:", err)
		return &pb.DockResponse{}, err
	}

	atcBody, _ := json.Marshal(result)
	return &pb.DockResponse{
		Status:  "Success",
		Message: string(atcBody),
	}, nil
}

// UpdateVolumeAttachment implements opensds.DockServer
func (ds *dockServer) UpdateVolumeAttachment(ctx context.Context, req *pb.DockRequest) (*pb.DockResponse, error) {
	log.Info("Dock server receive update volume attachment request, vr =", req)

	var dck, hostInfo = &api.DockSpec{}, &api.HostInfo{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Error("When parsing dock info:", err)
	}
	if err := json.Unmarshal([]byte(req.GetHostInfo()), hostInfo); err != nil {
		log.Error("Error occured in dock module when parsing host info:", err)
		return &pb.DockResponse{}, err
	}

	err := dock.NewDockHub(dck.GetDriverName()).UpdateVolumeAttachment(
		req.GetVolumeId(),
		hostInfo.Host,
		req.GetMountpoint())
	if err != nil {
		log.Error("Error occured in dock module when update volume attachment:", err)
		return &pb.DockResponse{}, err
	}

	result, err := db.C.UpdateVolumeAttachment(
		req.GetVolumeId(),
		req.GetAttachmentId(),
		req.GetMountpoint(),
		hostInfo)
	if err != nil {
		log.Error("Error occured in dock module when update volume attachment in db:", err)
		return &pb.DockResponse{}, err
	}

	atcBody, _ := json.Marshal(result)
	return &pb.DockResponse{
		Status:  "Success",
		Message: string(atcBody),
	}, nil
}

// DeleteVolumeAttachment implements opensds.DockServer
func (ds *dockServer) DeleteVolumeAttachment(ctx context.Context, req *pb.DockRequest) (*pb.DockResponse, error) {
	log.Info("Dock server receive delete volume attachment request, vr =", req)

	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Error("When parsing dock info:", err)
	}

	if err := dock.NewDockHub(dck.GetDriverName()).DeleteVolumeAttachment(req.GetVolumeId()); err != nil {
		log.Error("Error occured in dock module when delete volume attachment:", err)
		if strings.Contains(err.Error(), "The status of volume is not in-use") {
			if err = db.C.DeleteVolumeAttachment(req.GetVolumeId(), req.GetAttachmentId()); err != nil {
				log.Error("Error occured in dock module when delete volume attachment in db:", err)
				return &pb.DockResponse{}, err
			}
		}
		return &pb.DockResponse{}, err
	}

	if err := db.C.DeleteVolumeAttachment(req.GetVolumeId(), req.GetAttachmentId()); err != nil {
		log.Error("Error occured in dock module when delete volume attachment in db:", err)
		return &pb.DockResponse{}, err
	}

	return &pb.DockResponse{
		Status:  "Success",
		Message: "Delete volume attachment success",
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

	// If volume snapshot uuid is null, generate it randomly.
	if snp.GetId() == "" {
		if ok := utils.NewSetter().SetUuid(snp); ok != nil {
			log.Error("When set volume snapshot uuid:", ok)
			return &pb.DockResponse{}, err
		}
	}

	// Set volume snapshot created time.
	if ok := utils.NewSetter().SetCreatedTimeStamp(snp); ok != nil {
		log.Error("When set volume snapshot created time:", ok)
		return &pb.DockResponse{}, err
	}

	result, err := db.C.CreateVolumeSnapshot(snp)
	if err != nil {
		log.Error("Error occured in dock module when create volume snapshot in db:", err)
		return &pb.DockResponse{}, err
	}

	snpBody, _ := json.Marshal(result)
	return &pb.DockResponse{
		Status:  "Success",
		Message: string(snpBody),
	}, nil
}

// GetVolumeSnapshot implements opensds.DockServer
func (ds *dockServer) GetVolumeSnapshot(ctx context.Context, req *pb.DockRequest) (*pb.DockResponse, error) {
	log.Info("Dock server receive get volume snapshot request, vr =", req)

	var dck = &api.DockSpec{}
	if err := json.Unmarshal([]byte(req.GetDockInfo()), dck); err != nil {
		log.Error("When parsing dock info:", err)
	}

	result, err := dock.NewDockHub(dck.GetDriverName()).GetSnapshot(req.GetSnapshotId())
	if err != nil {
		log.Error("Error occured in dock module when get snapshot:", err)
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

	if err := db.C.DeleteVolumeSnapshot(req.GetSnapshotId()); err != nil {
		log.Error("Error occured in dock module when delete volume snapshot in db:", err)
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
