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
	"github.com/opensds/opensds/pkg/apiserver/grpcapi"
	db "github.com/opensds/opensds/pkg/db/api"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
)

type VolumeRequestDeliver interface {
	createVolume() *pb.Response

	getVolume() (*api.VolumeResponse, error)

	listVolumes() (*[]api.VolumeResponse, error)

	deleteVolume() *pb.Response

	createVolumeAttachment() *pb.Response

	getVolumeAttachment() (*api.VolumeAttachment, error)

	listVolumeAttachments() (*[]api.VolumeAttachment, error)

	updateVolumeAttachment() *pb.Response

	deleteVolumeAttachment() *pb.Response

	createVolumeSnapshot() *pb.Response

	getVolumeSnapshot() (*api.VolumeSnapshot, error)

	listVolumeSnapshots() (*[]api.VolumeSnapshot, error)

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

func (vr *VolumeRequest) getVolume() (*api.VolumeResponse, error) {
	return db.GetVolume(vr.Schema.Id)
}

func (vr *VolumeRequest) listVolumes() (*[]api.VolumeResponse, error) {
	return db.ListVolumes()
}

func (vr *VolumeRequest) deleteVolume() *pb.Response {
	return grpcapi.DeleteVolume(vr.Schema, vr.Profile)
}

func (vr *VolumeRequest) createVolumeAttachment() *pb.Response {
	return grpcapi.CreateVolumeAttachment(vr.Schema, vr.Profile)
}

func (vr *VolumeRequest) getVolumeAttachment() (*api.VolumeAttachment, error) {
	return db.GetVolumeAttachment(vr.Schema.Id, vr.Schema.AttachmentId)
}

func (vr *VolumeRequest) listVolumeAttachments() (*[]api.VolumeAttachment, error) {
	return db.ListVolumeAttachments(vr.Schema.Id)
}

func (vr *VolumeRequest) updateVolumeAttachment() *pb.Response {
	return grpcapi.UpdateVolumeAttachment(vr.Schema, vr.Profile)
}

func (vr *VolumeRequest) deleteVolumeAttachment() *pb.Response {
	return grpcapi.DeleteVolumeAttachment(vr.Schema, vr.Profile)
}

func (vr VolumeRequest) createVolumeSnapshot() *pb.Response {
	return grpcapi.CreateVolumeSnapshot(vr.Schema, vr.Profile)
}

func (vr VolumeRequest) getVolumeSnapshot() (*api.VolumeSnapshot, error) {
	return db.GetVolumeSnapshot(vr.Schema.SnapshotId)
}

func (vr VolumeRequest) listVolumeSnapshots() (*[]api.VolumeSnapshot, error) {
	return db.ListVolumeSnapshots()
}

func (vr VolumeRequest) deleteVolumeSnapshot() *pb.Response {
	return grpcapi.DeleteVolumeSnapshot(vr.Schema, vr.Profile)
}

func CreateVolume(vrd VolumeRequestDeliver) (api.VolumeResponse, error) {
	var nullResponse = api.VolumeResponse{}

	result := vrd.createVolume()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("Create volume error:", err)
		return nullResponse, err
	}

	var volumeResponse = api.VolumeResponse{}
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &volumeResponse); err != nil {
		return nullResponse, err
	}
	return volumeResponse, nil
}

func GetVolume(vrd VolumeRequestDeliver) (api.VolumeResponse, error) {
	var nullResponse = api.VolumeResponse{}

	vol, err := vrd.getVolume()
	if err != nil {
		log.Println("Get volume error:", err)
		return nullResponse, err
	}

	return *vol, nil
}

func ListVolumes(vrd VolumeRequestDeliver) ([]api.VolumeResponse, error) {
	var nullResponses = []api.VolumeResponse{}

	vols, err := vrd.listVolumes()
	if err != nil {
		log.Println("List all volumes error:", err)
		return nullResponses, err
	}

	return *vols, nil
}

func DeleteVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	var defaultResponse = api.DefaultResponse{}

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

func CreateVolumeAttachment(vrd VolumeRequestDeliver) (api.VolumeAttachment, error) {
	var nullVolumeAttachment = api.VolumeAttachment{}

	result := vrd.createVolumeAttachment()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("Create volume error:", err)
		return nullVolumeAttachment, err
	}

	var volumeAttachment = api.VolumeAttachment{}
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &volumeAttachment); err != nil {
		return nullVolumeAttachment, err
	}
	return volumeAttachment, nil
}

func GetVolumeAttachment(vrd VolumeRequestDeliver) (api.VolumeAttachment, error) {
	var nullVolumeAttachment = api.VolumeAttachment{}

	atc, err := vrd.getVolumeAttachment()
	if err != nil {
		log.Println("Get volume error:", err)
		return nullVolumeAttachment, err
	}

	return *atc, nil
}

func ListVolumeAttachments(vrd VolumeRequestDeliver) ([]api.VolumeAttachment, error) {
	var nullVolumeAttachments = []api.VolumeAttachment{}

	atcs, err := vrd.listVolumeAttachments()
	if err != nil {
		log.Println("List all volumes error:", err)
		return nullVolumeAttachments, err
	}

	return *atcs, nil
}

func UpdateVolumeAttachment(vrd VolumeRequestDeliver) (api.VolumeAttachment, error) {
	var nullVolumeAttachment = api.VolumeAttachment{}

	result := vrd.updateVolumeAttachment()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("Get volume error:", err)
		return nullVolumeAttachment, err
	}

	var volumeAttachment = api.VolumeAttachment{}
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &volumeAttachment); err != nil {
		return nullVolumeAttachment, err
	}
	return volumeAttachment, nil
}

func DeleteVolumeAttachment(vrd VolumeRequestDeliver) api.DefaultResponse {
	var defaultResponse = api.DefaultResponse{}

	result := vrd.deleteVolumeAttachment()
	if result.GetStatus() == "Failure" {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = result.GetError()
		log.Println("Delete volume error:", defaultResponse.Error)
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}

func CreateVolumeSnapshot(vrd VolumeRequestDeliver) (api.VolumeSnapshot, error) {
	var nullSnapshot = api.VolumeSnapshot{}

	result := vrd.createVolumeSnapshot()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("Create volume error:", err)
		return nullSnapshot, err
	}

	var snapshot = api.VolumeSnapshot{}
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &snapshot); err != nil {
		return nullSnapshot, err
	}
	return snapshot, nil
}

func GetVolumeSnapshot(vrd VolumeRequestDeliver) (api.VolumeSnapshot, error) {
	var nullSnapshot = api.VolumeSnapshot{}

	snp, err := vrd.getVolumeSnapshot()
	if err != nil {
		log.Println("Get volume error:", err)
		return nullSnapshot, err
	}

	return *snp, nil
}

func ListVolumeSnapshots(vrd VolumeRequestDeliver) ([]api.VolumeSnapshot, error) {
	var nullSnapshots = []api.VolumeSnapshot{}

	snps, err := vrd.listVolumeSnapshots()
	if err != nil {
		log.Println("List all volumes error:", err)
		return nullSnapshots, err
	}

	return *snps, nil
}

func DeleteVolumeSnapshot(vrd VolumeRequestDeliver) api.DefaultResponse {
	var defaultResponse = api.DefaultResponse{}

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
