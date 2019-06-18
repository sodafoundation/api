// Copyright 2017 The OpenSDS Authors.
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

package volume

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

type fakeClient struct{}

func NewFakeClient() client.Client {
	return &fakeClient{}
}

func (fc *fakeClient) Connect(edp string) error {
	return nil
}

func (fc *fakeClient) Close() {
	return
}

// Create a volume
func (fc *fakeClient) CreateVolume(ctx context.Context, in *pb.CreateVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteVolume,
			},
		},
	}, nil
}

// Delete a volume
func (fc *fakeClient) DeleteVolume(ctx context.Context, in *pb.DeleteVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// Extend a volume
func (fc *fakeClient) ExtendVolume(ctx context.Context, in *pb.ExtendVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteVolume,
			},
		},
	}, nil
}

// Create a volume attachment
func (fc *fakeClient) CreateVolumeAttachment(ctx context.Context, in *pb.CreateVolumeAttachmentOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteAttachment,
			},
		},
	}, nil
}

func (fc *fakeClient) DeleteVolumeAttachment(ctx context.Context, in *pb.DeleteVolumeAttachmentOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// Create a volume snapshot
func (fc *fakeClient) CreateVolumeSnapshot(ctx context.Context, in *pb.CreateVolumeSnapshotOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteSnapshot,
			},
		},
	}, nil
}

// Delete a volume snapshot
func (fc *fakeClient) DeleteVolumeSnapshot(ctx context.Context, in *pb.DeleteVolumeSnapshotOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// Create a volume snapshot
func (fc *fakeClient) CreateVolumeGroup(ctx context.Context, in *pb.CreateVolumeGroupOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteVolumeGroup,
			},
		},
	}, nil
}

// Create a volume snapshot
func (fc *fakeClient) UpdateVolumeGroup(ctx context.Context, in *pb.UpdateVolumeGroupOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteVolumeGroup,
			},
		},
	}, nil
}

// Delete a volume snapshot
func (fc *fakeClient) DeleteVolumeGroup(ctx context.Context, in *pb.DeleteVolumeGroupOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// Attach a volume
func (fc *fakeClient) AttachVolume(ctx context.Context, in *pb.AttachVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: "",
			},
		},
	}, nil
}

// Detach a volume
func (fc *fakeClient) DetachVolume(ctx context.Context, in *pb.DetachVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// Create a volume attachment
func (fc *fakeClient) CreateReplication(ctx context.Context, in *pb.CreateReplicationOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteReplication,
			},
		},
	}, nil
}

// Delete a replication
func (fc *fakeClient) DeleteReplication(ctx context.Context, in *pb.DeleteReplicationOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// Enable a replication
func (fc *fakeClient) EnableReplication(ctx context.Context, in *pb.EnableReplicationOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// Disable a replication
func (fc *fakeClient) DisableReplication(ctx context.Context, in *pb.DisableReplicationOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// Failover a replication
func (fc *fakeClient) FailoverReplication(ctx context.Context, in *pb.FailoverReplicationOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

func (fc *fakeClient) CollectMetrics(ctx context.Context, in *pb.CollectMetricsOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

func NewFakeController() Controller {
	return &controller{
		Client:   NewFakeClient(),
		DockInfo: &model.DockSpec{},
	}
}

func TestCreateVolume(t *testing.T) {
	fc := NewFakeController()
	var expected = &SampleVolumes[0]

	result, err := fc.CreateVolume(&pb.CreateVolumeOpts{})
	if err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteVolume(t *testing.T) {
	fc := NewFakeController()

	result := fc.DeleteVolume(&pb.DeleteVolumeOpts{})
	if result != nil {
		t.Errorf("Expected %v, got %v\n", nil, result)
	}
}

func TestExtendVolume(t *testing.T) {
	fc := NewFakeController()
	var expected = &SampleVolumes[0]

	result, err := fc.ExtendVolume(&pb.ExtendVolumeOpts{})
	if err != nil {
		t.Errorf("Failed to extend volume, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestCreateVolumeAttachment(t *testing.T) {
	fc := NewFakeController()
	var expected = &SampleAttachments[0]

	result, err := fc.CreateVolumeAttachment(&pb.CreateVolumeAttachmentOpts{})
	if err != nil {
		t.Errorf("Failed to create volume attachment, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteVolumeAttachment(t *testing.T) {
	fc := NewFakeController()

	result := fc.DeleteVolumeAttachment(&pb.DeleteVolumeAttachmentOpts{})
	if result != nil {
		t.Errorf("Expected %v, got %v\n", nil, result)
	}
}

func TestCreateVolumeSnapshot(t *testing.T) {
	fc := NewFakeController()
	var expected = &SampleSnapshots[0]

	result, err := fc.CreateVolumeSnapshot(&pb.CreateVolumeSnapshotOpts{})
	if err != nil {
		t.Errorf("Failed to create volume snapshot, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteVolumeSnapshot(t *testing.T) {
	fc := NewFakeController()

	result := fc.DeleteVolumeSnapshot(&pb.DeleteVolumeSnapshotOpts{})
	if result != nil {
		t.Errorf("Expected %v, got %v\n", nil, result)
	}
}

func TestCreateReplication(t *testing.T) {
	fc := NewFakeController()
	var expected = &SampleReplications[0]

	result, err := fc.CreateReplication(&pb.CreateReplicationOpts{})
	if err != nil {
		t.Errorf("Failed to create replication, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteReplication(t *testing.T) {
	fc := NewFakeController()

	result := fc.DeleteReplication(&pb.DeleteReplicationOpts{})
	if result != nil {
		t.Errorf("Expected %v, got %v\n", nil, result)
	}
}

func TestEnableReplication(t *testing.T) {
	fc := NewFakeController()

	result := fc.EnableReplication(&pb.EnableReplicationOpts{})
	if result != nil {
		t.Errorf("Expected %v, got %v\n", nil, result)
	}
}

func TestDisableReplication(t *testing.T) {
	fc := NewFakeController()

	result := fc.DisableReplication(&pb.DisableReplicationOpts{})
	if result != nil {
		t.Errorf("Expected %v, got %v\n", nil, result)
	}
}

func (fc *fakeClient) CreateFileShare(ctx context.Context, in *pb.CreateFileShareOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakeClient) CreateFileShareAcl(ctx context.Context, in *pb.CreateFileShareAclOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakeClient) DeleteFileShareAcl(ctx context.Context, in *pb.DeleteFileShareAclOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

// DeleteFileShare provides a mock function with given fields: ctx, in, opts
func (fc *fakeClient) DeleteFileShare(ctx context.Context, in *pb.DeleteFileShareOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakeClient) CreateFileShareSnapshot(ctx context.Context, in *pb.CreateFileShareSnapshotOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakeClient) DeleteFileShareSnapshot(ctx context.Context, in *pb.DeleteFileShareSnapshotOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}
