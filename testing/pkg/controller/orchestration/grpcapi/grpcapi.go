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
This module implements the api module client of messaging service.

*/

package grpcapi

import (
	"fmt"

	client "github.com/opensds/opensds/testing/pkg/grpc/dock/fake_client"
	pb "github.com/opensds/opensds/testing/pkg/grpc/fake_opensds"
)

func CreateVolume(vr *pb.VolumeRequest) *pb.Response {
	var falseResp *pb.Response

	resp, err := client.CreateVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func GetVolume(vr *pb.VolumeRequest) *pb.Response {
	var falseResp *pb.Response

	resp, err := client.GetVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func ListVolumes(vr *pb.VolumeRequest) *pb.Response {
	var falseResp *pb.Response

	resp, err := client.ListVolumes(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DeleteVolume(vr *pb.VolumeRequest) *pb.Response {
	var falseResp *pb.Response

	resp, err := client.DeleteVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func AttachVolume(vr *pb.VolumeRequest) *pb.Response {
	var falseResp *pb.Response

	resp, err := client.AttachVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DetachVolume(vr *pb.VolumeRequest) *pb.Response {
	var falseResp *pb.Response

	resp, err := client.DetachVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func MountVolume(vr *pb.VolumeRequest) *pb.Response {
	var falseResp *pb.Response

	resp, err := client.MountVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func UnmountVolume(vr *pb.VolumeRequest) *pb.Response {
	var falseResp *pb.Response

	resp, err := client.UnmountVolume(vr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func CreateShare(sr *pb.ShareRequest) *pb.Response {
	var falseResp *pb.Response

	resp, err := client.CreateShare(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func GetShare(sr *pb.ShareRequest) *pb.Response {
	var falseResp *pb.Response

	resp, err := client.GetShare(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func ListShares(sr *pb.ShareRequest) *pb.Response {
	var falseResp *pb.Response

	resp, err := client.ListShares(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}

func DeleteShare(sr *pb.ShareRequest) *pb.Response {
	var falseResp *pb.Response

	resp, err := client.DeleteShare(sr)
	if err != nil {
		falseResp.Status = "Failure"
		falseResp.Error = fmt.Sprintf("%v", err)
		return falseResp
	}

	return resp
}
