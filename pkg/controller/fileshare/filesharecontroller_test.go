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

package fileshare

import (
	"context"
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/dock/client"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	. "github.com/opensds/opensds/testutils/collection"
	"google.golang.org/grpc"
)

type fakefileshareClient struct{}

func NewFakeClient() client.Client {
	return &fakefileshareClient{}
}

func (fc *fakefileshareClient) Connect(edp string) error {
	return nil
}

func (fc *fakefileshareClient) Close() {
	return
}

func (fc *fakefileshareClient) CreateFileShare(ctx context.Context, in *pb.CreateFileShareOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteFileShare,
			},
		},
	}, nil
}

func (fc *fakefileshareClient) CreateFileShareAcl(ctx context.Context, in *pb.CreateFileShareAclOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteFileShareAcl,
			},
		},
	}, nil
}

// DeleteFileShare provides a mock function with given fields: ctx, in, opts
func (fc *fakefileshareClient) DeleteFileShare(ctx context.Context, in *pb.DeleteFileShareOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteFileShare,
			},
		},
	}, nil
}

// DeleteFileShareAcl provides a mock function with given fields: ctx, in, opts
func (fc *fakefileshareClient) DeleteFileShareAcl(ctx context.Context, in *pb.DeleteFileShareAclOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteFileShareAcl,
			},
		},
	}, nil
}

func (fc *fakefileshareClient) CreateFileShareSnapshot(ctx context.Context, in *pb.CreateFileShareSnapshotOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteFileShareSnapshot,
			},
		},
	}, nil
}

func (fc *fakefileshareClient) DeleteFileShareSnapshot(ctx context.Context, in *pb.DeleteFileShareSnapshotOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteFileShareSnapshot,
			},
		},
	}, nil
}

func NewFakeController() Controller {
	return &controller{
		Client:   NewFakeClient(),
		DockInfo: &model.DockSpec{},
	}
}

func TestCreateFileShare(t *testing.T) {
	fc := NewFakeController()
	var expected = &SampleShares[0]

	result, err := fc.CreateFileShare(&pb.CreateFileShareOpts{})
	if err != nil {
		t.Errorf("failed to create fileshare, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v\n", expected, result)
	}
}

func TestDeleteFileShare(t *testing.T) {
	fc := NewFakeController()

	result := fc.DeleteFileShare(&pb.DeleteFileShareOpts{})
	if result != nil {
		t.Errorf("expected %v, got %v\n", nil, result)
	}
}

func TestCreateFileShareSnapshot(t *testing.T) {
	fc := NewFakeController()
	var expected = &SampleShareSnapshots[0]

	result, err := fc.CreateFileShareSnapshot(&pb.CreateFileShareSnapshotOpts{})
	if err != nil {
		t.Errorf("failed to create file share snapshot, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v\n", expected, result)
	}
}

func TestDeleteFileShareSnapshot(t *testing.T) {
	fc := NewFakeController()

	result := fc.DeleteFileShareSnapshot(&pb.DeleteFileShareSnapshotOpts{})
	if result != nil {
		t.Errorf("expected %v, got %v\n", nil, result)
	}
}

func TestCreateFileShareAcl(t *testing.T) {
	fc := NewFakeController()
	var expected = &SampleFileSharesAcl[0]

	result, err := fc.CreateFileShareAcl(&pb.CreateFileShareAclOpts{})
	if err != nil {
		t.Errorf("failed to create fileshare acl, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v\n", expected, result)
	}
}

func TestDeleteFileShareAcl(t *testing.T) {
	fc := NewFakeController()

	result := fc.DeleteFileShareAcl(&pb.DeleteFileShareAclOpts{})
	if result != nil {
		t.Errorf("expected %v, got %v\n", nil, result)
	}
}

func (fc *fakefileshareClient) CreateVolume(ctx context.Context, in *pb.CreateVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) DeleteVolume(ctx context.Context, in *pb.DeleteVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) ExtendVolume(ctx context.Context, in *pb.ExtendVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) CreateVolumeSnapshot(ctx context.Context, in *pb.CreateVolumeSnapshotOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) DeleteVolumeSnapshot(ctx context.Context, in *pb.DeleteVolumeSnapshotOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) CreateVolumeAttachment(ctx context.Context, in *pb.CreateVolumeAttachmentOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) DeleteVolumeAttachment(ctx context.Context, in *pb.DeleteVolumeAttachmentOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) CreateReplication(ctx context.Context, in *pb.CreateReplicationOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) DeleteReplication(ctx context.Context, in *pb.DeleteReplicationOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) EnableReplication(ctx context.Context, in *pb.EnableReplicationOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) DisableReplication(ctx context.Context, in *pb.DisableReplicationOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) FailoverReplication(ctx context.Context, in *pb.FailoverReplicationOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) CreateVolumeGroup(ctx context.Context, in *pb.CreateVolumeGroupOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) UpdateVolumeGroup(ctx context.Context, in *pb.UpdateVolumeGroupOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) DeleteVolumeGroup(ctx context.Context, in *pb.DeleteVolumeGroupOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) CollectMetrics(ctx context.Context, in *pb.CollectMetricsOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) AttachVolume(ctx context.Context, in *pb.AttachVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) DetachVolume(ctx context.Context, in *pb.DetachVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) GetMetrics(ctx context.Context, in *pb.GetMetricsOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakefileshareClient) GetUrls(ctx context.Context, in *pb.NoParams, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}
