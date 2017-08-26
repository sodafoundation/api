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
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/grpc/dock/client"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
	api "github.com/opensds/opensds/pkg/model"
)

func NewFakeController(req *pb.DockRequest) Controller {
	return &controller{
		Client:  client.NewFakeClient(""),
		Request: req,
	}
}

func TestCreateVolume(t *testing.T) {
	fc := NewFakeController(&pb.DockRequest{})
	var expected = &client.SampleVolume

	result, err := fc.CreateVolume()
	if err != nil {
		t.Errorf("Failed to create volume, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteVolume(t *testing.T) {
	fc := NewFakeController(&pb.DockRequest{})
	var expected = &api.Response{Status: "Success"}

	result := fc.DeleteVolume()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestCreateVolumeAttachment(t *testing.T) {
	fc := NewFakeController(&pb.DockRequest{})
	var expected = &client.SampleAttachment

	result, err := fc.CreateVolumeAttachment()
	if err != nil {
		t.Errorf("Failed to create volume attachment, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestUpdateVolumeAttachment(t *testing.T) {
	fc := NewFakeController(&pb.DockRequest{})
	var expected = &client.SampleModifiedAttachment

	result, err := fc.UpdateVolumeAttachment()
	if err != nil {
		t.Errorf("Failed to update volume attachment, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteVolumeAttachment(t *testing.T) {
	fc := NewFakeController(&pb.DockRequest{})
	var expected = &api.Response{Status: "Success"}

	result := fc.DeleteVolumeAttachment()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestCreateVolumeSnapshot(t *testing.T) {
	fc := NewFakeController(&pb.DockRequest{})
	var expected = &client.SampleSnapshot

	result, err := fc.CreateVolumeSnapshot()
	if err != nil {
		t.Errorf("Failed to create volume snapshot, err is %v\n", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}

func TestDeleteVolumeSnapshot(t *testing.T) {
	fc := NewFakeController(&pb.DockRequest{})
	var expected = &api.Response{Status: "Success"}

	result := fc.DeleteVolumeSnapshot()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v\n", expected, result)
	}
}
