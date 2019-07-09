// Copyright 2019 The OpenSDS Authors.
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

package controller

import (
	"context"
	"net"
	"reflect"
	"runtime"

	log "github.com/golang/glog"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"google.golang.org/grpc"
)

func NewGrpcServer(port string) *GrpcServer {
	ctrl := NewController()
	return &GrpcServer{
		Controller: ctrl,
		Port:       port,
	}
}

type GrpcServer struct {
	*Controller
	Port string
}

// Run method would start the listen mechanism of controller module.
func (g *GrpcServer) Run() error {
	// New Grpc Server
	s := grpc.NewServer()
	// Register controller service.
	pb.RegisterControllerServer(s, g)
	pb.RegisterFileShareControllerServer(s, g)

	// Listen the controller server port.
	lis, err := net.Listen("tcp", g.Port)
	if err != nil {
		log.Fatalf("failed to listen: %+v", err)
		return err
	}

	log.Info("osdslet server initialized! Start listening on port:", lis.Addr())

	// Start controller server watching loop.
	defer s.Stop()
	return s.Serve(lis)
}

// AsyncDecorator is used to provide common method to print info to logs when got error.
// It will check the wrapped function parameter type and number, also the return value number.
func AsyncDecorator(fn interface{}, args ...interface{}) {
	f := reflect.ValueOf(fn)
	if f.Type().NumIn() != len(args) {
		log.Errorf("incorrect number of parameter(s) for function %v!\n",
			runtime.FuncForPC(f.Pointer()).Name())
		runtime.Goexit()
	}
	for i := 0; i < f.Type().NumIn(); i++ {
		if f.Type().In(i) != reflect.TypeOf(args[i]) && !reflect.TypeOf(args[i]).ConvertibleTo(f.Type().In(i)) {
			log.Errorf("parameter(s) for function %v is wrong type (should be %v)\n",
				runtime.FuncForPC(f.Pointer()).Name(), f.Type().In(i))
			runtime.Goexit()
		}
	}
	inputs := make([]reflect.Value, len(args))
	for k, in := range args {
		inputs[k] = reflect.ValueOf(in)
	}

	out := f.Call(inputs)

	// Wrapped function return value number must equal to 2
	if len(out) != 2 {
		log.Errorf("incorrect number of return value(s)\n")
		runtime.Goexit()
	}
	if !out[1].IsNil() {
		log.Errorf("call function '%v' failed: %v", runtime.FuncForPC(f.Pointer()).Name(), out[1].Interface())
	}
}

func (g *GrpcServer) CreateVolume(contx context.Context, opt *pb.CreateVolumeOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.CreateVolume, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// DeleteVolume implements pb.ControllerServer.DeleteVolume
func (g *GrpcServer) DeleteVolume(contx context.Context, opt *pb.DeleteVolumeOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.DeleteVolume, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// ExtendVolume implements pb.ControllerServer.ExtendVolume
func (g *GrpcServer) ExtendVolume(contx context.Context, opt *pb.ExtendVolumeOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.ExtendVolume, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// CreateVolumeAttachment implements pb.ControllerServer.CreateVolumeAttachment
func (g *GrpcServer) CreateVolumeAttachment(contx context.Context, opt *pb.CreateVolumeAttachmentOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.CreateVolumeAttachment, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// DeleteVolumeAttachment implements pb.ControllerServer.DeleteVolumeAttachment
func (g *GrpcServer) DeleteVolumeAttachment(contx context.Context, opt *pb.DeleteVolumeAttachmentOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.DeleteVolumeAttachment, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// CreateVolumeSnapshot implements pb.ControllerServer.CreateVolumeSnapshot
func (g *GrpcServer) CreateVolumeSnapshot(contx context.Context, opt *pb.CreateVolumeSnapshotOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.CreateVolumeSnapshot, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// DeleteVolumeSnapshot implements pb.ControllerServer.DeleteVolumeSnapshot
func (g *GrpcServer) DeleteVolumeSnapshot(contx context.Context, opt *pb.DeleteVolumeSnapshotOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.DeleteVolumeSnapshot, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// CreateReplication implements pb.ControllerServer.CreateReplication
func (g *GrpcServer) CreateReplication(contx context.Context, opt *pb.CreateReplicationOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.CreateReplication, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// DeleteReplication implements pb.ControllerServer.DeleteReplication
func (g *GrpcServer) DeleteReplication(contx context.Context, opt *pb.DeleteReplicationOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.DeleteReplication, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// EnableReplication implements pb.ControllerServer.EnableReplication
func (g *GrpcServer) EnableReplication(contx context.Context, opt *pb.EnableReplicationOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.EnableReplication, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// DisableReplication implements pb.ControllerServer.DisableReplication
func (g *GrpcServer) DisableReplication(contx context.Context, opt *pb.DisableReplicationOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.DisableReplication, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// FailoverReplication implements pb.ControllerServer.FailoverReplication
func (g *GrpcServer) FailoverReplication(contx context.Context, opt *pb.FailoverReplicationOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.FailoverReplication, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// CreateVolumeGroup implements pb.ControllerServer.CreateVolumeGroup
func (g *GrpcServer) CreateVolumeGroup(contx context.Context, opt *pb.CreateVolumeGroupOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.CreateVolumeGroup, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// UpdateVolumeGroup implements pb.ControllerServer.UpdateVolumeGroup
func (g *GrpcServer) UpdateVolumeGroup(contx context.Context, opt *pb.UpdateVolumeGroupOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.UpdateVolumeGroup, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// DeleteVolumeGroup implements pb.ControllerServer.DeleteVolumeGroup
func (g *GrpcServer) DeleteVolumeGroup(contx context.Context, opt *pb.DeleteVolumeGroupOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.DeleteVolumeGroup, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// CreateFileShare implements pb.ControllerServer.CreateFileShare
func (g *GrpcServer) CreateFileShare(contx context.Context, opt *pb.CreateFileShareOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.CreateFileShare, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// DeleteFileShare implements pb.ControllerServer.DeleteFileShare
func (g *GrpcServer) DeleteFileShare(contx context.Context, opt *pb.DeleteFileShareOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.DeleteFileShare, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// CreateFileShare implements pb.ControllerServer.CreateFileShareAcl
func (g *GrpcServer) CreateFileShareAcl(contx context.Context, opt *pb.CreateFileShareAclOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.CreateFileShareAcl, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// DeleteFileShareAcl implements pb.ControllerServer.DeleteFileShareAcl
func (g *GrpcServer) DeleteFileShareAcl(contx context.Context, opt *pb.DeleteFileShareAclOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.DeleteFileShareAcl, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// CreateFileShareSnapshot implements pb.ControllerServer.CreateFileShareSnapshot
func (g *GrpcServer) CreateFileShareSnapshot(contx context.Context, opt *pb.CreateFileShareSnapshotOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.CreateFileShareSnapshot, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}

// DeleteFileshareSnapshot implements pb.ControllerServer.DeleteFileshareSnapshot
func (g *GrpcServer) DeleteFileShareSnapshot(contx context.Context, opt *pb.DeleteFileShareSnapshotOpts) (*pb.GenericResponse, error) {
	go AsyncDecorator(g.Controller.DeleteFileShareSnapshot, contx, opt)
	return pb.GenericResponseResult("grpc cast success"), nil
}
