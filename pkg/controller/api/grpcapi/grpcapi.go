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

var falseResp = &pb.Response{}

func CreateVolume(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		Name:           schema.Name,
		VolumeType:     schema.VolumeType,
		Size:           schema.Size,
	}

	resp, err := client.CreateVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func GetVolume(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		Id:             schema.Id,
	}

	resp, err := client.GetVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func ListVolumes(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		AllowDetails:   schema.AllowDetails,
	}

	resp, err := client.ListVolumes(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DeleteVolume(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		Id:             schema.Id,
	}

	resp, err := client.DeleteVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func AttachVolume(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		Id:             schema.Id,
	}

	resp, err := client.AttachVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DetachVolume(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		Device:         schema.Device,
	}

	resp, err := client.DetachVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func MountVolume(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		MountDir:       schema.MountDir,
		Device:         schema.Device,
		FsType:         schema.FsType,
	}

	resp, err := client.MountVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func UnmountVolume(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		MountDir:       schema.MountDir,
	}

	resp, err := client.UnmountVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func CreateVolumeSnapshot(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile:  string(body),
		DockId:          schema.DockId,
		SnapshotName:    schema.SnapshotName,
		Id:              schema.Id,
		Description:     schema.Description,
		ForceSnapshoted: schema.ForceSnapshoted,
	}

	resp, err := client.CreateVolumeSnapshot(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func GetVolumeSnapshot(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		SnapshotId:     schema.SnapshotId,
	}

	resp, err := client.GetVolumeSnapshot(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func ListVolumeSnapshots(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
	}

	resp, err := client.ListVolumeSnapshots(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DeleteVolumeSnapshot(schema *api.VolumeOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	vr := &pb.VolumeRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		SnapshotId:     schema.SnapshotId,
	}

	resp, err := client.DeleteVolumeSnapshot(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func CreateShare(schema *api.ShareOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	sr := &pb.ShareRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		Name:           schema.Name,
		ShareType:      schema.ShareType,
		ShareProto:     schema.ShareProto,
		Size:           schema.Size,
	}

	resp, err := client.CreateShare(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func GetShare(schema *api.ShareOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	sr := &pb.ShareRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		Id:             schema.Id,
	}

	resp, err := client.GetShare(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func ListShares(schema *api.ShareOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	sr := &pb.ShareRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		AllowDetails:   schema.AllowDetails,
	}

	resp, err := client.ListShares(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DeleteShare(schema *api.ShareOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	sr := &pb.ShareRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		Id:             schema.Id,
	}

	resp, err := client.DeleteShare(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func AttachShare(schema *api.ShareOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	sr := &pb.ShareRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		Id:             schema.Id,
	}

	resp, err := client.AttachShare(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DetachShare(schema *api.ShareOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	sr := &pb.ShareRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		Device:         schema.Device,
	}

	resp, err := client.DetachShare(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func MountShare(schema *api.ShareOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	sr := &pb.ShareRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		MountDir:       schema.MountDir,
		Device:         schema.Device,
		FsType:         schema.FsType,
	}

	resp, err := client.MountShare(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func UnmountShare(schema *api.ShareOperationSchema, profile *api.StorageProfile) *pb.Response {
	body, err := json.Marshal(profile)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	sr := &pb.ShareRequest{
		StorageProfile: string(body),
		DockId:         schema.DockId,
		MountDir:       schema.MountDir,
	}

	resp, err := client.UnmountShare(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}
