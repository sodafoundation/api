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
This module implements the api module client of grpc service.

*/

package grpcapi

import (
	"encoding/json"
	"fmt"

	api "github.com/opensds/opensds/pkg/api/v1"
	"github.com/opensds/opensds/pkg/grpc/controller/orchestration/client"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
)

var falseResp = &pb.Response{
	Status: "Failure",
}

func CreateVolume(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		VolumeName:     schema.Name,
		Size:           schema.Size,
	}

	resp, err := client.CreateVolume(vr)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DeleteVolume(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		VolumeId:       schema.Id,
	}

	resp, err := client.DeleteVolume(vr)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func CreateVolumeAttachment(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	pbody, err := json.Marshal(profile)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}
	hbody, err := json.Marshal(schema.HostInfo)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(pbody),
		HostInfo:       string(hbody),
		VolumeId:       schema.Id,
		DoLocalAttach:  schema.DoLocalAttach,
		MultiPath:      schema.MultiPath,
	}

	resp, err := client.CreateVolumeAttachment(vr)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func UpdateVolumeAttachment(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		VolumeId:       schema.Id,
		AttachmentId:   schema.AttachmentId,
		Mountpoint:     schema.Mountpoint,
	}

	resp, err := client.UpdateVolumeAttachment(vr)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DeleteVolumeAttachment(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		VolumeId:       schema.Id,
		AttachmentId:   schema.AttachmentId,
	}

	resp, err := client.DeleteVolumeAttachment(vr)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func CreateVolumeSnapshot(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile:      string(body),
		SnapshotName:        schema.SnapshotName,
		VolumeId:            schema.Id,
		SnapshotDescription: schema.Description,
	}

	resp, err := client.CreateVolumeSnapshot(vr)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DeleteVolumeSnapshot(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		VolumeId:       schema.Id,
		SnapshotId:     schema.SnapshotId,
	}

	resp, err := client.DeleteVolumeSnapshot(vr)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func CreateShare(schema *api.ShareOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	sr := &pb.ShareRequest{
		StorageProfile: string(body),
		Name:           schema.Name,
		ShareProto:     schema.ShareProto,
		Size:           schema.Size,
	}

	resp, err := client.CreateShare(sr)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func GetShare(schema *api.ShareOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	sr := &pb.ShareRequest{
		StorageProfile: string(body),
		Id:             schema.Id,
	}

	resp, err := client.GetShare(sr)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func ListShares(schema *api.ShareOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	sr := &pb.ShareRequest{
		StorageProfile: string(body),
	}

	resp, err := client.ListShares(sr)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DeleteShare(schema *api.ShareOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	sr := &pb.ShareRequest{
		StorageProfile: string(body),
		Id:             schema.Id,
	}

	resp, err := client.DeleteShare(sr)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func AttachShare(schema *api.ShareOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	sr := &pb.ShareRequest{
		StorageProfile: string(body),
		Id:             schema.Id,
	}

	resp, err := client.AttachShare(sr)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DetachShare(schema *api.ShareOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	sr := &pb.ShareRequest{
		StorageProfile: string(body),
		Device:         schema.Device,
	}

	resp, err := client.DetachShare(sr)
	if err != nil {
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}
