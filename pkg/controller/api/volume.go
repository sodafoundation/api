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
This module implements the entry into CRUD operation of volumes.

*/

package api

import (
	"encoding/json"
	"errors"
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
	"github.com/opensds/opensds/pkg/controller/api/grpcapi"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
)

var defaultResponse api.DefaultResponse
var nullSnapshotResponse api.VolumeSnapshotResponse
var nullSnapshotResponses []api.VolumeSnapshotResponse

type VolumeRequestDeliver interface {
	createVolume() *pb.Response

	getVolume() *pb.Response

	listVolumes() *pb.Response

	deleteVolume() *pb.Response

	attachVolume() *pb.Response

	detachVolume() *pb.Response

	mountVolume() *pb.Response

	unmountVolume() *pb.Response

	createVolumeSnapshot() *pb.Response

	getVolumeSnapshot() *pb.Response

	listVolumeSnapshots() *pb.Response

	deleteVolumeSnapshot() *pb.Response
}

// VolumeRequest is a structure for all properties of
// a volume request
type VolumeRequest struct {
	Schema  *api.VolumeOperationSchema `json:"schema"`
	Profile *api.StorageProfile        `json:"profile"`
}

func (vr *VolumeRequest) createVolume() *pb.Response {
	return grpcapi.CreateVolume(vr.Schema, vr.Profile)
}

func (vr *VolumeRequest) getVolume() *pb.Response {
	return grpcapi.GetVolume(vr.Schema, vr.Profile)
}

func (vr *VolumeRequest) listVolumes() *pb.Response {
	return grpcapi.ListVolumes(vr.Schema, vr.Profile)
}

func (vr *VolumeRequest) deleteVolume() *pb.Response {
	return grpcapi.DeleteVolume(vr.Schema, vr.Profile)
}

func (vr *VolumeRequest) attachVolume() *pb.Response {
	return grpcapi.AttachVolume(vr.Schema, vr.Profile)
}

func (vr *VolumeRequest) detachVolume() *pb.Response {
	return grpcapi.DetachVolume(vr.Schema, vr.Profile)
}

func (vr *VolumeRequest) mountVolume() *pb.Response {
	return grpcapi.MountVolume(vr.Schema, vr.Profile)
}

func (vr *VolumeRequest) unmountVolume() *pb.Response {
	return grpcapi.UnmountVolume(vr.Schema, vr.Profile)
}

func (vr VolumeRequest) createVolumeSnapshot() *pb.Response {
	return grpcapi.CreateVolumeSnapshot(vr.Schema, vr.Profile)
}

func (vr VolumeRequest) getVolumeSnapshot() *pb.Response {
	return grpcapi.GetVolumeSnapshot(vr.Schema, vr.Profile)
}

func (vr VolumeRequest) listVolumeSnapshots() *pb.Response {
	return grpcapi.ListVolumeSnapshots(vr.Schema, vr.Profile)
}

func (vr VolumeRequest) deleteVolumeSnapshot() *pb.Response {
	return grpcapi.DeleteVolumeSnapshot(vr.Schema, vr.Profile)
}

func CreateVolume(vrd VolumeRequestDeliver) (api.VolumeResponse, error) {
	var nullResponse api.VolumeResponse

	result := vrd.createVolume()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("Create volume error:", err)
		return nullResponse, err
	}

	var volumeResponse api.VolumeResponse
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &volumeResponse); err != nil {
		return nullResponse, err
	}
	return volumeResponse, nil
}

func GetVolume(vrd VolumeRequestDeliver) (api.VolumeDetailResponse, error) {
	var nullResponse api.VolumeDetailResponse

	result := vrd.getVolume()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("Get volume error:", err)
		return nullResponse, err
	}

	var volumeDetailResponse api.VolumeDetailResponse
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &volumeDetailResponse); err != nil {
		return nullResponse, err
	}
	return volumeDetailResponse, nil
}

func ListVolumes(vrd VolumeRequestDeliver) ([]api.VolumeResponse, error) {
	var nullResponses []api.VolumeResponse

	result := vrd.listVolumes()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("List all volumes error:", err)
		return nullResponses, err
	}

	var volumesResponse []api.VolumeResponse
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &volumesResponse); err != nil {
		return nullResponses, err
	}
	return volumesResponse, nil
}

func DeleteVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	result := vrd.deleteVolume()
	if result.GetStatus() == "Failure" {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = result.GetError()
		log.Println("Delete volume error:", defaultResponse.Error)
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}

func AttachVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	result := vrd.attachVolume()
	if result.GetStatus() == "Failure" {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = result.GetError()
		log.Println("Attach volume error:", defaultResponse.Error)
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	defaultResponse.Message = result.GetMessage()
	return defaultResponse
}

func DetachVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	result := vrd.detachVolume()
	if result.GetStatus() == "Failure" {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = result.GetError()
		log.Println("Detach volume error:", defaultResponse.Error)
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}

func MountVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	result := vrd.mountVolume()
	if result.GetStatus() == "Failure" {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = result.GetError()
		log.Println("Mount volume error:", defaultResponse.Error)
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}

func UnmountVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	result := vrd.unmountVolume()
	if result.GetStatus() == "Failure" {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = result.GetError()
		log.Println("Unmount volume error:", defaultResponse.Error)
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}

func CreateVolumeSnapshot(vrd VolumeRequestDeliver) (api.VolumeSnapshotResponse, error) {
	result := vrd.createVolumeSnapshot()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("Create volume error:", err)
		return nullSnapshotResponse, err
	}

	var snapshotResponse api.VolumeSnapshotResponse
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &snapshotResponse); err != nil {
		return nullSnapshotResponse, err
	}
	return snapshotResponse, nil
}

func GetVolumeSnapshot(vrd VolumeRequestDeliver) (api.VolumeSnapshotResponse, error) {
	result := vrd.getVolumeSnapshot()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("Get volume error:", err)
		return nullSnapshotResponse, err
	}

	var snapshotResponse api.VolumeSnapshotResponse
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &snapshotResponse); err != nil {
		return nullSnapshotResponse, err
	}
	return snapshotResponse, nil
}

func ListVolumeSnapshots(vrd VolumeRequestDeliver) ([]api.VolumeSnapshotResponse, error) {
	result := vrd.listVolumeSnapshots()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("List all volumes error:", err)
		return nullSnapshotResponses, err
	}

	var snapshotResponses []api.VolumeSnapshotResponse
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &snapshotResponses); err != nil {
		return nullSnapshotResponses, err
	}
	return snapshotResponses, nil
}

func DeleteVolumeSnapshot(vrd VolumeRequestDeliver) api.DefaultResponse {
	result := vrd.deleteVolumeSnapshot()
	if result.GetStatus() == "Failure" {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = result.GetError()
		log.Println("Delete volume error:", defaultResponse.Error)
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}
