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
	_ "log"

	"github.com/opensds/opensds/testing/pkg/controller/orchestration/grpcapi"
	pb "github.com/opensds/opensds/testing/pkg/grpc/fake_opensds"
)

func CreateVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	return grpcapi.CreateVolume(vr), nil
}

func GetVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	return grpcapi.GetVolume(vr), nil
}

func ListVolumes(vr *pb.VolumeRequest) (*pb.Response, error) {
	return grpcapi.ListVolumes(vr), nil
}

func DeleteVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	return grpcapi.DeleteVolume(vr), nil
}

func AttachVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	return grpcapi.AttachVolume(vr), nil
}

func DetachVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	return grpcapi.DetachVolume(vr), nil
}

func MountVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	return grpcapi.MountVolume(vr), nil
}

func UnmountVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	return grpcapi.UnmountVolume(vr), nil
}
