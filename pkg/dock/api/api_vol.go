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
This module implements the entry into operations of storageDock module.

*/

package api

import (
	"fmt"
	"log"

	"github.com/opensds/opensds/pkg/dock"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
)

func CreateVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	result, err := dock.CreateVolume(vr.GetResoureType(),
		vr.GetName(),
		vr.GetVolumeType(),
		vr.GetSize())
	if err != nil {
		log.Println("Error occured in adapter module when create volume:", err)
		resp := &pb.Response{
			Status: "Failure",
			Error:  fmt.Sprintf("%v", err),
		}
		return resp, nil
	}

	resp := &pb.Response{
		Status:  "Success",
		Message: result,
	}
	return resp, nil
}

func GetVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	result, err := dock.GetVolume(vr.GetResoureType(), vr.GetId())
	if err != nil {
		log.Println("Error occured in adapter module when get volume:", err)
		resp := &pb.Response{
			Status: "Failure",
			Error:  fmt.Sprintf("%v", err),
		}
		return resp, nil
	}

	resp := &pb.Response{
		Status:  "Success",
		Message: result,
	}
	return resp, nil
}

func ListVolumes(vr *pb.VolumeRequest) (*pb.Response, error) {
	result, err := dock.GetAllVolumes(vr.GetResoureType(),
		vr.GetAllowDetails())
	if err != nil {
		log.Println("Error occured in adapter module when list volumes:", err)
		resp := &pb.Response{
			Status: "Failure",
			Error:  fmt.Sprintf("%v", err),
		}
		return resp, nil
	}

	resp := &pb.Response{
		Status:  "Success",
		Message: result,
	}
	return resp, nil
}

func DeleteVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	result, err := dock.DeleteVolume(vr.GetResoureType(),
		vr.GetId())
	if err != nil {
		log.Println("Error occured in adapter module when delete volume:", err)
		resp := &pb.Response{
			Status: "Failure",
			Error:  fmt.Sprintf("%v", err),
		}
		return resp, nil
	}

	resp := &pb.Response{
		Status:  "Success",
		Message: result,
	}
	return resp, nil
}

func AttachVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	result, err := dock.AttachVolume(vr.GetResoureType(),
		vr.GetId(),
		vr.GetHost(),
		vr.GetDevice())
	if err != nil {
		log.Println("Error occured in adapter module when attach volume:", err)
		resp := &pb.Response{
			Status: "Failure",
			Error:  fmt.Sprintf("%v", err),
		}
		return resp, nil
	}

	resp := &pb.Response{
		Status:  "Success",
		Message: result,
	}
	return resp, nil
}

func DetachVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	result, err := dock.DetachVolume(vr.GetResoureType(),
		vr.GetId(),
		vr.GetAttachment())
	if err != nil {
		log.Println("Error occured in adapter module when detach volume:", err)
		resp := &pb.Response{
			Status: "Failure",
			Error:  fmt.Sprintf("%v", err),
		}
		return resp, nil
	}

	resp := &pb.Response{
		Status:  "Success",
		Message: result,
	}
	return resp, nil
}

func MountVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	result, err := dock.MountVolume(vr.GetMountDir(),
		vr.GetDevice(),
		vr.GetFsType())
	if err != nil {
		log.Println("Error occured in adapter module when mount volume:", err)
		resp := &pb.Response{
			Status: "Failure",
			Error:  fmt.Sprintf("%v", err),
		}
		return resp, nil
	}

	resp := &pb.Response{
		Status:  "Success",
		Message: result,
	}
	return resp, nil
}

func UnmountVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	result, err := dock.UnmountVolume(vr.GetMountDir())
	if err != nil {
		log.Println("Error occured in adapter module when unmount volume:", err)
		resp := &pb.Response{
			Status: "Failure",
			Error:  fmt.Sprintf("%v", err),
		}
		return resp, nil
	}

	resp := &pb.Response{
		Status:  "Success",
		Message: result,
	}
	return resp, nil
}
