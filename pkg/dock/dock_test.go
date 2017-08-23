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

package dock

import (
	"reflect"
	"testing"

	api "github.com/opensds/opensds/pkg/model"
)

var (
	fd    = &DockHub{ResourceType: "default"}
	volID = "038754e8-87d6-11e7-a459-6f41395c65a5"
	snpID = "a463ce6e-87d6-11e7-8b27-737694075ac9"
)

func TestCreateVolume(t *testing.T) {
	var name = "ups-volume"
	var size int64 = 1
	var expectedResult = &api.VolumeSpec{BaseModel: &api.BaseModel{}}

	result, err := fd.CreateVolume(name, size)
	if err != nil {
		t.Errorf("Create volume failed, got err %v\n", err)
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected %v, got %v\n", expectedResult, result)
	}
}

func TestGetVolume(t *testing.T) {
	var expectedResult = &api.VolumeSpec{BaseModel: &api.BaseModel{}}

	result, err := fd.GetVolume(volID)
	if err != nil {
		t.Errorf("Get volume failed, got err %v\n", err)
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected %v, got %v\n", expectedResult, result)
	}
}

func TestDeleteVolume(t *testing.T) {

	if err := fd.DeleteVolume(volID); err != nil {
		t.Errorf("Delete volume failed, got err %v\n", err)
	}
}

func TestCreateVolumeAttachment(t *testing.T) {
	var doLocalAttach, multiPath bool
	var expectedResult = &api.VolumeAttachmentSpec{ConnectionInfo: &api.ConnectionInfo{}}

	result, err := fd.CreateVolumeAttachment(volID, doLocalAttach, multiPath, nil)
	if err != nil {
		t.Errorf("Create volume attachment failed, got err %v\n", err)
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected %v, got %v\n", expectedResult, result)
	}
}

func TestUpdateVolumeAttachment(t *testing.T) {
	var host, mountpoint = "localhost", "/mnt"

	if err := fd.UpdateVolumeAttachment(volID, host, mountpoint); err != nil {
		t.Errorf("Update volume attachment failed, got err %v\n", err)
	}
}

func TestDeleteVolumeAttachment(t *testing.T) {

	if err := fd.DeleteVolumeAttachment(volID); err != nil {
		t.Errorf("Delete volume attachment failed, got err %v\n", err)
	}
}

func TestCreateSnapshot(t *testing.T) {
	var name, description = "ups-volume", "fake volume for testing"
	var expectedResult = &api.VolumeSnapshotSpec{
		BaseModel: &api.BaseModel{},
		VolumeId:  volID,
	}

	result, err := fd.CreateSnapshot(name, volID, description)
	if err != nil {
		t.Errorf("Create volume snapshot failed, got err %v\n", err)
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected %v, got %v\n", expectedResult, result)
	}
}

func TestGetSnapshot(t *testing.T) {
	var expectedResult = &api.VolumeSnapshotSpec{BaseModel: &api.BaseModel{}}

	result, err := fd.GetSnapshot(snpID)
	if err != nil {
		t.Errorf("Get volume snapshot failed, got err %v\n", err)
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected %v, got %v\n", expectedResult, result)
	}
}

func TestDeleteSnapshot(t *testing.T) {

	if err := fd.DeleteSnapshot(snpID); err != nil {
		t.Errorf("Delete volume snapshot failed, got err %v\n", err)
	}
}
