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

/*
This file is only used for testing.

*/

package client

import (
	"encoding/json"

	pb "github.com/opensds/opensds/pkg/grpc/opensds"
	api "github.com/opensds/opensds/pkg/model"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type fakeClient struct {
	TargetPlace string
}

func NewFakeClient(address string) Client {
	return &fakeClient{
		TargetPlace: address,
	}
}

func (fc *fakeClient) Update(dockInfo string) error {
	return nil
}

func (fc *fakeClient) Close() {
	return
}

// Create a volume
func (fc *fakeClient) CreateVolume(ctx context.Context, in *pb.DockRequest, opts ...grpc.CallOption) (*pb.DockResponse, error) {
	volBody, _ := json.Marshal(&SampleVolume)

	return &pb.DockResponse{
		Status:  "Success",
		Message: string(volBody),
	}, nil
}

// Get a volume
func (fc *fakeClient) GetVolume(ctx context.Context, in *pb.DockRequest, opts ...grpc.CallOption) (*pb.DockResponse, error) {
	volBody, _ := json.Marshal(&SampleVolume)

	return &pb.DockResponse{
		Status:  "Success",
		Message: string(volBody),
	}, nil
}

// Delete a volume
func (fc *fakeClient) DeleteVolume(ctx context.Context, in *pb.DockRequest, opts ...grpc.CallOption) (*pb.DockResponse, error) {
	return &pb.DockResponse{
		Status: "Success",
	}, nil
}

// Create a volume attachment
func (fc *fakeClient) CreateVolumeAttachment(ctx context.Context, in *pb.DockRequest, opts ...grpc.CallOption) (*pb.DockResponse, error) {
	volBody, _ := json.Marshal(&SampleAttachment)

	return &pb.DockResponse{
		Status:  "Success",
		Message: string(volBody),
	}, nil
}

// Update a volume attachment
func (fc *fakeClient) UpdateVolumeAttachment(ctx context.Context, in *pb.DockRequest, opts ...grpc.CallOption) (*pb.DockResponse, error) {
	volBody, _ := json.Marshal(&SampleModifiedAttachment)

	return &pb.DockResponse{
		Status:  "Success",
		Message: string(volBody),
	}, nil
}

// Delete a volume attachment
func (fc *fakeClient) DeleteVolumeAttachment(ctx context.Context, in *pb.DockRequest, opts ...grpc.CallOption) (*pb.DockResponse, error) {
	return &pb.DockResponse{
		Status: "Success",
	}, nil
}

// Create a volume snapshot
func (fc *fakeClient) CreateVolumeSnapshot(ctx context.Context, in *pb.DockRequest, opts ...grpc.CallOption) (*pb.DockResponse, error) {
	volBody, _ := json.Marshal(&SampleSnapshot)

	return &pb.DockResponse{
		Status:  "Success",
		Message: string(volBody),
	}, nil
}

// Get a volume snapshot
func (fc *fakeClient) GetVolumeSnapshot(ctx context.Context, in *pb.DockRequest, opts ...grpc.CallOption) (*pb.DockResponse, error) {
	volBody, _ := json.Marshal(&SampleSnapshot)

	return &pb.DockResponse{
		Status:  "Success",
		Message: string(volBody),
	}, nil
}

// Delete a volume snapshot
func (fc *fakeClient) DeleteVolumeSnapshot(ctx context.Context, in *pb.DockRequest, opts ...grpc.CallOption) (*pb.DockResponse, error) {
	return &pb.DockResponse{
		Status: "Success",
	}, nil
}

var (
	SampleVolume = api.VolumeSpec{
		BaseModel: &api.BaseModel{
			Id:        "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
			CreatedAt: "2017-08-02T09:17:05",
		},
		Name:        "fake-volume",
		Description: "fake volume for testing",
		Size:        1,
		PoolId:      "80287bf8-66de-11e7-b031-f3b0af1675ba",
	}

	SampleAttachment = api.VolumeAttachmentSpec{
		BaseModel: &api.BaseModel{
			Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		},
		Name:        "fake-volume-attachment",
		Description: "fake volume attachment for testing",
		VolumeId:    "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}

	SampleModifiedAttachment = api.VolumeAttachmentSpec{
		BaseModel: &api.BaseModel{
			Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		},
		Name:        "modified-fake-volume-attachment",
		Description: "modified fake volume attachment for testing",
		VolumeId:    "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}

	SampleSnapshot = api.VolumeSnapshotSpec{
		BaseModel: &api.BaseModel{
			Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		},
		Name:        "fake-volume-snapshot",
		Description: "fake volume snapshot for testing",
		VolumeId:    "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
	}
)
