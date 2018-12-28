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
	"fmt"
	"net"

	log "github.com/golang/glog"
	c "github.com/opensds/opensds/pkg/context"
	pb "github.com/opensds/opensds/pkg/controller/proto"
	"github.com/opensds/opensds/pkg/model"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func NewCtlServer(port string) *ctlServer {
	return &ctlServer{
		Port: port,
	}
}

type ctlServer struct {
	Port string
}

// Run method would start the listen mechanism of controller module.
func (cs *ctlServer) Run() error {
	// New Grpc Server
	s := grpc.NewServer()
	// Register controller service.
	pb.RegisterControllerServer(s, cs)

	// Listen the controller server port.
	lis, err := net.Listen("tcp", cs.Port)
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
func (cs *ctlServer) CreateVolume(ctx context.Context, opt *pb.CreateVolumeOpts) (*pb.GenericResponse, error) {
	var req model.VolumeSpec

	log.Info("Controller server receive create volume request, vr =", opt)

	if err := json.Unmarshal([]byte(opt.Message), &req); err != nil {
		reason := fmt.Sprintf("Decode volume request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	// NOTE:The real volume creation process.
	// CreateVolume request is sent to the Dock. Dock will update volume status to "available"
	// after volume creation is completed.
	var errchan = make(chan error, 1)
	defer close(errchan)
	go Brain.CreateVolume(c.NewAdminContext(), &req, errchan)
	if err := <-errchan; err != nil {
		reason := fmt.Sprintf("Marshal volume created result failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}

	return nil, nil
}

// DeleteVolume implements pb.ControllerServer.DeleteVolume
func (cs *ctlServer) DeleteVolume(ctx context.Context, opt *pb.DeleteVolumeOpts) (*pb.GenericResponse, error) {
	var req model.VolumeSpec

	log.Info("Controller server receive delete volume request, vr =", opt)

	if err := json.Unmarshal([]byte(opt.Message), &req); err != nil {
		reason := fmt.Sprintf("Decode volume request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	// NOTE:The real volume deletion process.
	// DeleteVolume request is sent to the Dock. Dock will remove volume record
	// after volume deletion is completed.
	var errchan = make(chan error, 1)
	defer close(errchan)
	go Brain.DeleteVolume(c.NewAdminContext(), &req, errchan)
	if err := <-errchan; err != nil {
		reason := fmt.Sprintf("Delete volume failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}

	return nil, nil
}

// ExtendVolume implements pb.ControllerServer.ExtendVolume
func (cs *ctlServer) ExtendVolume(ctx context.Context, opt *pb.ExtendVolumeOpts) (*pb.GenericResponse, error) {
	var req model.ExtendVolumeSpec

	log.Info("Controller server receive extend volume request, vr =", opt)

	if err := json.Unmarshal([]byte(opt.Message), &req); err != nil {
		reason := fmt.Sprintf("Decode volume request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	// NOTE:The real volume extention process.
	// ExtendVolume request is sent to the Dock. Dock will update volume status to "available"
	// after volume extention is completed.
	var errchan = make(chan error, 1)
	defer close(errchan)
	go Brain.ExtendVolume(c.NewAdminContext(), opt.Id, req.NewSize, errchan)
	if err := <-errchan; err != nil {
		reason := fmt.Sprintf("Extend volume failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	return nil, nil
}

// CreateVolumeAttachment implements pb.ControllerServer.CreateVolumeAttachment
func (cs *ctlServer) CreateVolumeAttachment(ctx context.Context, opt *pb.CreateVolumeAttachmentOpts) (*pb.GenericResponse, error) {
	var req model.VolumeAttachmentSpec

	log.Info("Controller server receive create volume attachment request, vr =", opt)

	if err := json.Unmarshal([]byte(opt.Message), &req); err != nil {
		reason := fmt.Sprintf("Decode volume attachment request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	// NOTE:The real volume attachment creation process.
	// Volume attachment creation request is sent to the Dock. Dock will update volume attachment status to "available"
	// after volume attachment creation is completed.
	errchan := make(chan error, 1)
	defer close(errchan)
	go Brain.CreateVolumeAttachment(c.NewAdminContext(), &req, errchan)
	if err := <-errchan; err != nil {
		reason := fmt.Sprintf("Create volume attachment failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}

	return nil, nil
}

// DeleteVolumeAttachment implements pb.ControllerServer.DeleteVolumeAttachment
func (cs *ctlServer) DeleteVolumeAttachment(ctx context.Context, opt *pb.DeleteVolumeAttachmentOpts) (*pb.GenericResponse, error) {
	var req model.VolumeAttachmentSpec

	log.Info("Controller server receive delete volume attachment request, vr =", opt)

	if err := json.Unmarshal([]byte(opt.Message), &req); err != nil {
		reason := fmt.Sprintf("Decode volume attachment request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	// NOTE:The real volume attachment deletion process.
	// Volume attachment deletion request is sent to the Dock. Dock will delete volume attachment from database
	// or update its status to "errorDeleting" if volume connection termination failed.
	var errchan = make(chan error, 1)
	go Brain.DeleteVolumeAttachment(c.NewAdminContext(), &req, errchan)
	defer close(errchan)
	if err := <-errchan; err != nil {
		reason := fmt.Sprintf("Delete volume attachment failed: %v", err.Error())
		log.Error(reason)
		return nil, err
	}

	return nil, nil
}

// CreateVolumeSnapshot implements pb.ControllerServer.CreateVolumeSnapshot
func (cs *ctlServer) CreateVolumeSnapshot(ctx context.Context, opt *pb.CreateVolumeSnapshotOpts) (*pb.GenericResponse, error) {
	var req model.VolumeSnapshotSpec

	log.Info("Controller server receive create volume snapshot request, vr =", opt)

	if err := json.Unmarshal([]byte(opt.Message), &req); err != nil {
		reason := fmt.Sprintf("Decode volume snapshot request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	// NOTE:The real volume snapshot creation process.
	// Volume snapshot creation request is sent to the Dock. Dock will update volume snapshot status to "available"
	// after volume snapshot creation is completed.
	errchan := make(chan error, 1)
	defer close(errchan)
	go Brain.CreateVolumeSnapshot(c.NewAdminContext(), &req, errchan)
	if err := <-errchan; err != nil {
		reason := fmt.Sprintf("Create volume snapshot failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}

	return nil, nil
}

// DeleteVolumeSnapshot implements pb.ControllerServer.DeleteVolumeSnapshot
func (cs *ctlServer) DeleteVolumeSnapshot(ctx context.Context, opt *pb.DeleteVolumeSnapshotOpts) (*pb.GenericResponse, error) {
	var req model.VolumeSnapshotSpec

	log.Info("Controller server receive delete volume snapshot request, vr =", opt)

	if err := json.Unmarshal([]byte(opt.Message), &req); err != nil {
		reason := fmt.Sprintf("Decode volume snapshot request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	// NOTE:The real volume snapshot deletion process.
	// Volume snapshot deletion request is sent to the Dock. Dock will delete volume snapshot from database
	// or update its status to "errorDeleting" if volume connection termination failed.
	var errchan = make(chan error, 1)
	go Brain.DeleteVolumeSnapshot(c.NewAdminContext(), &req, errchan)
	defer close(errchan)
	if err := <-errchan; err != nil {
		reason := fmt.Sprintf("Delete volume snapshot failed: %v", err.Error())
		log.Error(reason)
		return nil, err
	}

	return nil, nil
}

// CreateReplication implements pb.ControllerServer.CreateReplication
func (cs *ctlServer) CreateReplication(ctx context.Context, opt *pb.CreateReplicationOpts) (*pb.GenericResponse, error) {
	var req model.ReplicationSpec

	log.Info("Controller server receive create volume replication request, vr =", opt)

	if err := json.Unmarshal([]byte(opt.Message), &req); err != nil {
		reason := fmt.Sprintf("Decode volume replication request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	// Call global controller variable to handle create replication request.
	if _, err := Brain.CreateReplication(c.NewAdminContext(), &req); err != nil {
		reason := fmt.Sprintf("Create replication failed: %v", err.Error())
		log.Error(reason)
		return nil, err
	}

	return nil, nil
}

// DeleteReplication implements pb.ControllerServer.DeleteReplication
func (cs *ctlServer) DeleteReplication(ctx context.Context, opt *pb.DeleteReplicationOpts) (*pb.GenericResponse, error) {
	var req model.ReplicationSpec

	log.Info("Controller server receive delete volume replication request, vr =", opt)

	if err := json.Unmarshal([]byte(opt.Message), &req); err != nil {
		reason := fmt.Sprintf("Decode volume replication request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	// Call global controller variable to handle delete replication request.
	if _, err := Brain.CreateReplication(c.NewAdminContext(), &req); err != nil {
		reason := fmt.Sprintf("Delete replication failed: %v", err.Error())
		log.Error(reason)
		return nil, err
	}

	return nil, nil
}

// Enable a replication
func (cs *ctlServer) EnableReplication(ctx context.Context, opt *pb.EnableReplicationOpts) (*pb.GenericResponse, error) {
	var req model.ReplicationSpec

	log.Info("Controller server receive enable volume replication request, vr =", opt)

	if err := json.Unmarshal([]byte(opt.Message), &req); err != nil {
		reason := fmt.Sprintf("Decode volume replication request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	// Call global controller variable to handle enable replication request.
	if err := Brain.EnableReplication(c.NewAdminContext(), &req); err != nil {
		reason := fmt.Sprintf("Enable replication failed: %v", err.Error())
		log.Error(reason)
		return nil, err
	}

	return nil, nil
}

// Disable a replication
func (cs *ctlServer) DisableReplication(ctx context.Context, opt *pb.DisableReplicationOpts) (*pb.GenericResponse, error) {
	var req model.ReplicationSpec

	log.Info("Controller server receive disable volume replication request, vr =", opt)

	if err := json.Unmarshal([]byte(opt.Message), &req); err != nil {
		reason := fmt.Sprintf("Decode volume replication request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	// Call global controller variable to handle disable replication request.
	if err := Brain.DisableReplication(c.NewAdminContext(), &req); err != nil {
		reason := fmt.Sprintf("Disable replication failed: %v", err.Error())
		log.Error(reason)
		return nil, err
	}

	return nil, nil
}

// Failover a replication
func (cs *ctlServer) FailoverReplication(ctx context.Context, opt *pb.FailoverReplicationOpts) (*pb.GenericResponse, error) {
	var req model.ReplicationSpec
	var foReq model.FailoverReplicationSpec

	log.Info("Controller server receive failover volume replication request, vr =", opt)

	if err := json.Unmarshal([]byte(opt.Message), &req); err != nil {
		reason := fmt.Sprintf("Decode volume replication request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	if err := json.Unmarshal([]byte(opt.FailoverMessage), &foReq); err != nil {
		reason := fmt.Sprintf("Decode volume replication request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	// Call global controller variable to handle failover replication request.
	if err := Brain.FailoverReplication(c.NewAdminContext(), &req, &foReq); err != nil {
		reason := fmt.Sprintf("Failover replication failed: %v", err.Error())
		log.Error(reason)
		return nil, err
	}

	return nil, nil
}

// Create a volume group
func (cs *ctlServer) CreateVolumeGroup(ctx context.Context, opt *pb.CreateVolumeGroupOpts) (*pb.GenericResponse, error) {
	var req model.VolumeGroupSpec

	log.Info("Controller server receive create volume group request, vr =", opt)

	if err := json.Unmarshal([]byte(opt.Message), &req); err != nil {
		reason := fmt.Sprintf("Decode volume group request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	// Call global controller variable to handle create volume group request.
	if err := Brain.CreateVolumeGroup(c.NewAdminContext(), &req); err != nil {
		reason := fmt.Sprintf("Create volume group failed: %v", err.Error())
		log.Error(reason)
		return nil, err
	}

	return nil, nil
}

// Update volume group
func (cs *ctlServer) UpdateVolumeGroup(ctx context.Context, opt *pb.UpdateVolumeGroupOpts) (*pb.GenericResponse, error) {
	var req model.VolumeGroupSpec

	log.Info("Controller server receive update volume group request, vr =", opt)

	if err := json.Unmarshal([]byte(opt.Message), &req); err != nil {
		reason := fmt.Sprintf("Decode volume group request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	// Call global controller variable to handle update volume group request.
	if err := Brain.UpdateVolumeGroup(c.NewAdminContext(), &req, opt.AddVolMessage, opt.RmVolMessage); err != nil {
		reason := fmt.Sprintf("Update volume group failed: %v", err.Error())
		log.Error(reason)
		return nil, err
	}

	return nil, nil
}

// Delete volume group
func (cs *ctlServer) DeleteVolumeGroup(ctx context.Context, opt *pb.DeleteVolumeGroupOpts) (*pb.GenericResponse, error) {
	var req model.VolumeGroupSpec

	log.Info("Controller server receive delete volume group request, vr =", opt)

	if err := json.Unmarshal([]byte(opt.Message), &req); err != nil {
		reason := fmt.Sprintf("Decode volume group request failed: %s", err.Error())
		log.Error(reason)
		return nil, err
	}
	// Call global controller variable to handle delete volume group request.
	if err := Brain.DeleteVolumeGroup(c.NewAdminContext(), &req); err != nil {
		reason := fmt.Sprintf("Delete volume group failed: %v", err.Error())
		log.Error(reason)
		return nil, err
	}

	return nil, nil
}
