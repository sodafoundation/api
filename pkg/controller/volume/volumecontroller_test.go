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

package volume

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/dock/client"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/testutils/collection"
	"golang.org/x/net/context"
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
func (fc *fakeClient) CreateAttachment(ctx context.Context, in *pb.CreateAttachmentOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteAttachment,
			},
		},
	}, nil
}

func (fc *fakeClient) DeleteAttachment(ctx context.Context, in *pb.DeleteAttachmentOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
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

	result, err := fc.CreateVolumeAttachment(&pb.CreateAttachmentOpts{})
	if err != nil {
		t.Errorf("Failed to create volume attachment, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteVolumeAttachment(t *testing.T) {
	fc := NewFakeController()

	result := fc.DeleteVolumeAttachment(&pb.DeleteAttachmentOpts{})
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
