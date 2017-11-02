// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package volume

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/dock/client"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type fakeClient struct {
	TargetPlace string
}

func NewFakeClient(address string) client.Client {
	return &fakeClient{
		TargetPlace: address,
	}
}

func (fc *fakeClient) Update(dockInfo *model.DockSpec) error {
	return nil
}

func (fc *fakeClient) Close() {
	return
}

// Create a volume
func (fc *fakeClient) CreateVolume(ctx context.Context, in *pb.CreateVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	volBody, _ := json.Marshal(&sampleVolume)

	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: string(volBody),
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

// Create a volume attachment
func (fc *fakeClient) CreateAttachment(ctx context.Context, in *pb.CreateAttachmentOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	volBody, _ := json.Marshal(&sampleAttachment)

	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: string(volBody),
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
	volBody, _ := json.Marshal(&sampleSnapshot)

	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: string(volBody),
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
		Client: NewFakeClient(""),
		//Request: req,
	}
}

func TestCreateVolume(t *testing.T) {
	fc := NewFakeController( /*&pb.DockRequest{}*/ )
	var expected = &sampleVolume

	result, err := fc.CreateVolume()
	if err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteVolume(t *testing.T) {
	fc := NewFakeController( /*&pb.DockRequest{}*/ )
	var expected = &model.Response{Status: "Success"}

	result := fc.DeleteVolume()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestCreateVolumeAttachment(t *testing.T) {
	fc := NewFakeController( /*&pb.DockRequest{}*/ )
	var expected = &sampleAttachment

	result, err := fc.CreateVolumeAttachment()
	if err != nil {
		t.Errorf("Failed to create volume attachment, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestCreateVolumeSnapshot(t *testing.T) {
	fc := NewFakeController( /*&pb.DockRequest{}*/ )
	var expected = &sampleSnapshot

	result, err := fc.CreateVolumeSnapshot()
	if err != nil {
		t.Errorf("Failed to create volume snapshot, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteVolumeSnapshot(t *testing.T) {
	fc := NewFakeController( /*&pb.DockRequest{}*/ )
	var expected = &model.Response{Status: "Success"}

	result := fc.DeleteVolumeSnapshot()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

var (
	sampleVolume = model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id:        "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
			CreatedAt: "2017-08-02T09:17:05",
		},
		Name:        "fake-volume",
		Description: "fake volume for testing",
		Size:        1,
		PoolId:      "80287bf8-66de-11e7-b031-f3b0af1675ba",
	}

	sampleAttachment = model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		},
		Name:        "fake-volume-attachment",
		Description: "fake volume attachment for testing",
		VolumeId:    "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}

	sampleModifiedAttachment = model.VolumeAttachmentSpec{
		BaseModel: &model.BaseModel{
			Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		},
		Name:        "modified-fake-volume-attachment",
		Description: "modified fake volume attachment for testing",
		VolumeId:    "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}

	sampleSnapshot = model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		},
		Name:        "fake-volume-snapshot",
		Description: "fake volume snapshot for testing",
		VolumeId:    "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}
)
