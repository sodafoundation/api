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
	"fmt"

	"github.com/opensds/opensds/pkg/grpc/controller/orchestration/client"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
)

func CreateVolume(dockId, resourceType, name, volumeType string, size int32) *pb.Response {
	var falseResp *pb.Response

	vr := &pb.VolumeRequest{
		DockId:      dockId,
		ResoureType: resourceType,
		Name:        name,
		VolumeType:  volumeType,
		Size:        size,
	}

	resp, err := client.CreateVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func GetVolume(dockId, resourceType string, volID string) *pb.Response {
	var falseResp *pb.Response

	vr := &pb.VolumeRequest{
		DockId:      dockId,
		ResoureType: resourceType,
		Id:          volID,
	}

	resp, err := client.GetVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func ListVolumes(dockId, resourceType string, allowDetails bool) *pb.Response {
	var falseResp *pb.Response

	vr := &pb.VolumeRequest{
		DockId:       dockId,
		ResoureType:  resourceType,
		AllowDetails: allowDetails,
	}

	resp, err := client.ListVolumes(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DeleteVolume(dockId, resourceType string, volID string) *pb.Response {
	var falseResp *pb.Response

	vr := &pb.VolumeRequest{
		DockId:      dockId,
		ResoureType: resourceType,
		Id:          volID,
	}

	resp, err := client.DeleteVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func AttachVolume(dockId, resourceType, volID, host, device string) *pb.Response {
	var falseResp *pb.Response

	vr := &pb.VolumeRequest{
		DockId:      dockId,
		ResoureType: resourceType,
		Id:          volID,
		Host:        host,
		Device:      device,
	}

	resp, err := client.AttachVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DetachVolume(dockId, resourceType, volID, attachment string) *pb.Response {
	var falseResp *pb.Response

	vr := &pb.VolumeRequest{
		DockId:      dockId,
		ResoureType: resourceType,
		Id:          volID,
		Attachment:  attachment,
	}

	resp, err := client.DetachVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func MountVolume(dockId, resourceType, mountDir, device, fsType string) *pb.Response {
	var falseResp *pb.Response

	vr := &pb.VolumeRequest{
		DockId:      dockId,
		ResoureType: resourceType,
		MountDir:    mountDir,
		Device:      device,
		FsType:      fsType,
	}

	resp, err := client.MountVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func UnmountVolume(dockId, resourceType, mountDir string) *pb.Response {
	var falseResp *pb.Response

	vr := &pb.VolumeRequest{
		DockId:      dockId,
		ResoureType: resourceType,
		MountDir:    mountDir,
	}

	resp, err := client.UnmountVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func CreateShare(dockId, resourceType, name, shrType, shrProto string, size int32) *pb.Response {
	var falseResp *pb.Response

	sr := &pb.ShareRequest{
		DockId:      dockId,
		ResoureType: resourceType,
		Name:        name,
		ShareType:   shrType,
		ShareProto:  shrProto,
		Size:        size,
	}

	resp, err := client.CreateShare(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func GetShare(dockId, resourceType string, shrID string) *pb.Response {
	var falseResp *pb.Response

	sr := &pb.ShareRequest{
		DockId:      dockId,
		ResoureType: resourceType,
		Id:          shrID,
	}

	resp, err := client.GetShare(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func ListShares(dockId, resourceType string, allowDetails bool) *pb.Response {
	var falseResp *pb.Response

	sr := &pb.ShareRequest{
		DockId:       dockId,
		ResoureType:  resourceType,
		AllowDetails: allowDetails,
	}

	resp, err := client.ListShares(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DeleteShare(dockId, resourceType string, shrID string) *pb.Response {
	var falseResp *pb.Response

	sr := &pb.ShareRequest{
		DockId:      dockId,
		ResoureType: resourceType,
		Id:          shrID,
	}

	resp, err := client.DeleteShare(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}
