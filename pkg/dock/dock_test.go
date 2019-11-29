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

/*
This module implements the entry into operations of storageDock module.
*/

package dock

import (
	"context"
	"reflect"
	"testing"

	"github.com/opensds/opensds/contrib/drivers"
	"github.com/opensds/opensds/contrib/drivers/filesharedrivers"
	"github.com/opensds/opensds/pkg/dock/discovery"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	data "github.com/opensds/opensds/testutils/collection"
)

func NewFakeDockServer() *dockServer {
	return &dockServer{
		Port:       "50050",
		Discoverer: discovery.NewDockDiscoverer(model.DockTypeProvioner),
	}
}

func NewFakeAttachDockServer() *dockServer {
	return &dockServer{
		Port:       "50050",
		Discoverer: discovery.NewDockDiscoverer(model.DockTypeAttacher),
	}
}

func TestNewDockServer(t *testing.T) {
	type args struct {
		dockType string
		port     string
	}
	tests := []struct {
		name string
		args args
		want *dockServer
	}{
		{
			name: "Provisioner docktype test",
			args: args{model.DockTypeProvioner, "50050"},
			want: NewFakeDockServer(),
		},
		{
			name: "Attacher docktype test",
			args: args{model.DockTypeAttacher, "50050"},
			want: NewFakeAttachDockServer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDockServer(tt.args.dockType, tt.args.port)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDockServer() = %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_dockServer_CreateFileShareAcl(t *testing.T) {
	type fields struct {
		Port            string
		Discoverer      discovery.DockDiscoverer
		Driver          drivers.VolumeDriver
		MetricDriver    drivers.MetricDriver
		FileShareDriver filesharedrivers.FileShareDriver
	}
	type args struct {
		ctx context.Context
		opt *pb.CreateFileShareAclOpts
	}
	var req = &pb.CreateFileShareAclOpts{
		Id:               "d2975ebe-d82c-430f-b28e-f373746a71ca",
		Description:      "This is a sample Acl for testing",
		Type:             "ip",
		AccessTo:         "10.21.23.10",
		AccessCapability: []string{"Read", "Write"},
	}
	want1 := &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: data.ByteFileShareAcl,
			},
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.GenericResponse
		wantErr bool
	}{
		{name: "Create file share acl dock test", args: args{
			ctx: context.Background(),
			opt: req,
		}, want: want1, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &dockServer{
				Port:            tt.fields.Port,
				Discoverer:      tt.fields.Discoverer,
				Driver:          tt.fields.Driver,
				MetricDriver:    tt.fields.MetricDriver,
				FileShareDriver: tt.fields.FileShareDriver,
			}
			_, err := ds.CreateFileShareAcl(tt.args.ctx, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("dockServer.CreateFileShareAcl() error = %v", err)
			}
		})
	}
}

func Test_dockServer_DeleteFileShareAcl(t *testing.T) {
	type fields struct {
		Port            string
		Discoverer      discovery.DockDiscoverer
		Driver          drivers.VolumeDriver
		MetricDriver    drivers.MetricDriver
		FileShareDriver filesharedrivers.FileShareDriver
	}
	type args struct {
		ctx context.Context
		opt *pb.DeleteFileShareAclOpts
	}
	var req = &pb.DeleteFileShareAclOpts{
		Id:          "d2975ebe-d82c-430f-b28e-f373746a71ca",
		Description: "This is a sample Acl for testing",
	}
	want1 := &pb.GenericResponse{
		Reply: nil,
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.GenericResponse
		wantErr bool
	}{
		{name: "Delete file share acl dock test", args: args{
			ctx: context.Background(),
			opt: req,
		}, want: want1, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &dockServer{
				Port:            tt.fields.Port,
				Discoverer:      tt.fields.Discoverer,
				Driver:          tt.fields.Driver,
				MetricDriver:    tt.fields.MetricDriver,
				FileShareDriver: tt.fields.FileShareDriver,
			}
			_, err := ds.DeleteFileShareAcl(tt.args.ctx, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("dockServer.DeleteFileShareAcl() error = %v", err)
			}
		})
	}
}

func Test_dockServer_CreateFileShare(t *testing.T) {
	type fields struct {
		Port            string
		Discoverer      discovery.DockDiscoverer
		Driver          drivers.VolumeDriver
		MetricDriver    drivers.MetricDriver
		FileShareDriver filesharedrivers.FileShareDriver
	}
	type args struct {
		ctx context.Context
		opt *pb.CreateFileShareOpts
	}
	prf := &data.SampleFileShareProfiles[0]
	var req = &pb.CreateFileShareOpts{
		Id:          "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:        "sample-fileshare",
		Description: "This is a sample fileshare for testing",
		Size:        1,
		PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
		Profile:     prf.ToJson(),
	}
	want1 := &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: data.ByteFileShare,
			},
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.GenericResponse
		wantErr bool
	}{
		{name: "Create file share dock test", args: args{
			ctx: context.Background(),
			opt: req,
		}, want: want1, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &dockServer{
				Port:            tt.fields.Port,
				Discoverer:      tt.fields.Discoverer,
				Driver:          tt.fields.Driver,
				MetricDriver:    tt.fields.MetricDriver,
				FileShareDriver: tt.fields.FileShareDriver,
			}
			_, err := ds.CreateFileShare(tt.args.ctx, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("dockServer.CreateFileShare() failed error = %v", err)
			}
		})
	}
}

func Test_dockServer_DeleteFileShare(t *testing.T) {
	type fields struct {
		Port            string
		Discoverer      discovery.DockDiscoverer
		Driver          drivers.VolumeDriver
		MetricDriver    drivers.MetricDriver
		FileShareDriver filesharedrivers.FileShareDriver
	}
	type args struct {
		ctx context.Context
		opt *pb.DeleteFileShareOpts
	}
	var req = &pb.DeleteFileShareOpts{
		Id:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name: "sample-fileshare",
	}
	want1 := &pb.GenericResponse{
		Reply: nil,
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.GenericResponse
		wantErr bool
	}{
		{name: "Delete file share dock test", args: args{
			ctx: context.Background(),
			opt: req,
		}, want: want1, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &dockServer{
				Port:            tt.fields.Port,
				Discoverer:      tt.fields.Discoverer,
				Driver:          tt.fields.Driver,
				MetricDriver:    tt.fields.MetricDriver,
				FileShareDriver: tt.fields.FileShareDriver,
			}
			_, err := ds.DeleteFileShare(tt.args.ctx, tt.args.opt)
			if err != nil {
				t.Errorf("dockServer.DeleteFileShare() error = %v", err)
			}
		})
	}
}

func Test_dockServer_CreateFileShareSnapshot(t *testing.T) {
	type fields struct {
		Port            string
		Discoverer      discovery.DockDiscoverer
		Driver          drivers.VolumeDriver
		MetricDriver    drivers.MetricDriver
		FileShareDriver filesharedrivers.FileShareDriver
	}
	type args struct {
		ctx context.Context
		opt *pb.CreateFileShareSnapshotOpts
	}
	var req = &pb.CreateFileShareSnapshotOpts{
		Id:          "3769855c-a102-11e7-b772-17b880d2f537",
		FileshareId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Name:        "sample-snapshot-01",
		Description: "This is the first sample snapshot for testing",
		Size:        int64(1),
	}
	want1 := &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: data.ByteFileShareSnapshot,
			},
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.GenericResponse
		wantErr bool
	}{
		{name: "Create file share snapshot dock test", args: args{
			ctx: context.Background(),
			opt: req,
		}, want: want1, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &dockServer{
				Port:            tt.fields.Port,
				Discoverer:      tt.fields.Discoverer,
				Driver:          tt.fields.Driver,
				MetricDriver:    tt.fields.MetricDriver,
				FileShareDriver: tt.fields.FileShareDriver,
			}
			_, err := ds.CreateFileShareSnapshot(tt.args.ctx, tt.args.opt)
			if err != nil {
				t.Errorf("dockServer.CreateFileShareSnapshot() failed error = %v", err)
			}

		})
	}
}

func Test_dockServer_DeleteFileShareSnapshot(t *testing.T) {
	type fields struct {
		Port            string
		Discoverer      discovery.DockDiscoverer
		Driver          drivers.VolumeDriver
		MetricDriver    drivers.MetricDriver
		FileShareDriver filesharedrivers.FileShareDriver
	}
	type args struct {
		ctx context.Context
		opt *pb.DeleteFileShareSnapshotOpts
	}
	req := &pb.DeleteFileShareSnapshotOpts{
		Id:          "3769855c-a102-11e7-b772-17b880d2f537",
		FileshareId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
	}
	want1 := &pb.GenericResponse{
		Reply: nil,
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.GenericResponse
		wantErr bool
	}{
		{name: "Delete file share snapshot dock test", args: args{
			ctx: context.Background(),
			opt: req,
		}, want: want1, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &dockServer{
				Port:            tt.fields.Port,
				Discoverer:      tt.fields.Discoverer,
				Driver:          tt.fields.Driver,
				MetricDriver:    tt.fields.MetricDriver,
				FileShareDriver: tt.fields.FileShareDriver,
			}
			_, err := ds.DeleteFileShareSnapshot(tt.args.ctx, tt.args.opt)
			if err != nil {
				t.Errorf("dockServer.DeleteFileShareSnapshot() error = %v", err)
			}
		})
	}
}
