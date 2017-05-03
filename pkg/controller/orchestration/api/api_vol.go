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
This module implements the enry into the operations of orchestration module.

Request about volume operation will be passed to the grpc client and requests
about other resources (database, fileSystem, etc) will be passed to metaData
service module.

*/

package api

import (
	"encoding/json"

	api "github.com/opensds/opensds/pkg/api/v1"
	"github.com/opensds/opensds/pkg/controller/orchestration/scheduler"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
)

var profile = &api.StorageProfile{}

func CreateVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(vr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	vs := &scheduler.VolumeScheduler{
		DesiredProfile: profile,
	}
	return vs.CreateVolume(vr)
}

func GetVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(vr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	vs := &scheduler.VolumeScheduler{
		DesiredProfile: profile,
	}
	return vs.GetVolume(vr)
}

func ListVolumes(vr *pb.VolumeRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(vr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	vs := &scheduler.VolumeScheduler{
		DesiredProfile: profile,
	}
	return vs.ListVolumes(vr)
}

func DeleteVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(vr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	vs := &scheduler.VolumeScheduler{
		DesiredProfile: profile,
	}
	return vs.DeleteVolume(vr)
}

func AttachVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(vr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	vs := &scheduler.VolumeScheduler{
		DesiredProfile: profile,
	}
	return vs.AttachVolume(vr)
}

func DetachVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(vr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	vs := &scheduler.VolumeScheduler{
		DesiredProfile: profile,
	}
	return vs.DetachVolume(vr)
}

func MountVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(vr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	vs := &scheduler.VolumeScheduler{
		DesiredProfile: profile,
	}
	return vs.MountVolume(vr)
}

func UnmountVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(vr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	vs := &scheduler.VolumeScheduler{
		DesiredProfile: profile,
	}
	return vs.UnmountVolume(vr)
}

func CreateVolumeSnapshot(vr *pb.VolumeRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(vr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	vs := &scheduler.VolumeScheduler{
		DesiredProfile: profile,
	}
	return vs.CreateVolumeSnapshot(vr)
}

func GetVolumeSnapshot(vr *pb.VolumeRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(vr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	vs := &scheduler.VolumeScheduler{
		DesiredProfile: profile,
	}
	return vs.GetVolumeSnapshot(vr)
}

func ListVolumeSnapshots(vr *pb.VolumeRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(vr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	vs := &scheduler.VolumeScheduler{
		DesiredProfile: profile,
	}
	return vs.ListVolumeSnapshots(vr)
}

func DeleteVolumeSnapshot(vr *pb.VolumeRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(vr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	vs := &scheduler.VolumeScheduler{
		DesiredProfile: profile,
	}
	return vs.DeleteVolumeSnapshot(vr)
}
