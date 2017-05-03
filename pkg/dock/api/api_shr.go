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
	"log"

	dock "github.com/opensds/opensds/pkg/dock/share"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
)

func CreateShare(sr *pb.ShareRequest) (*pb.Response, error) {
	result, err := dock.CreateShare(sr.GetResourceType(),
		sr.GetName(),
		sr.GetShareType(),
		sr.GetShareProto(),
		sr.GetSize())
	if err != nil {
		log.Println("Error occured in dock module when create share:", err)
		return &pb.Response{}, err
	}

	return &pb.Response{
		Status:  "Success",
		Message: result,
	}, nil
}

func GetShare(sr *pb.ShareRequest) (*pb.Response, error) {
	result, err := dock.GetShare(sr.GetResourceType(), sr.GetId())
	if err != nil {
		log.Println("Error occured in dock module when get share:", err)
		return &pb.Response{}, err
	}

	return &pb.Response{
		Status:  "Success",
		Message: result,
	}, nil
}

func ListShares(sr *pb.ShareRequest) (*pb.Response, error) {
	result, err := dock.GetAllShares(sr.GetResourceType(),
		sr.GetAllowDetails())
	if err != nil {
		log.Println("Error occured in dock module when list shares:", err)
		return &pb.Response{}, err
	}

	return &pb.Response{
		Status:  "Success",
		Message: result,
	}, nil
}

func DeleteShare(sr *pb.ShareRequest) (*pb.Response, error) {
	result, err := dock.DeleteShare(sr.GetResourceType(),
		sr.GetId())
	if err != nil {
		log.Println("Error occured in dock module when delete share:", err)
		return &pb.Response{}, err
	}

	return &pb.Response{
		Status:  "Success",
		Message: result,
	}, nil
}

func AttachShare(sr *pb.ShareRequest) (*pb.Response, error) {
	result, err := dock.AttachShare(sr.GetResourceType(),
		sr.GetId())
	if err != nil {
		log.Println("Error occured in dock module when attach share:", err)
		return &pb.Response{}, err
	}

	return &pb.Response{
		Status:  "Success",
		Message: result,
	}, nil
}

func DetachShare(sr *pb.ShareRequest) (*pb.Response, error) {
	result, err := dock.DetachShare(sr.GetResourceType(),
		sr.GetDevice())
	if err != nil {
		log.Println("Error occured in dock module when detach share:", err)
		return &pb.Response{}, err
	}

	return &pb.Response{
		Status:  "Success",
		Message: result,
	}, nil
}

func MountShare(sr *pb.ShareRequest) (*pb.Response, error) {
	result, err := dock.MountShare(sr.GetMountDir(),
		sr.GetDevice(),
		sr.GetFsType())
	if err != nil {
		log.Println("Error occured in dock module when mount share:", err)
		return &pb.Response{}, err
	}

	return &pb.Response{
		Status:  "Success",
		Message: result,
	}, nil
}

func UnmountShare(sr *pb.ShareRequest) (*pb.Response, error) {
	result, err := dock.UnmountShare(sr.GetMountDir())
	if err != nil {
		log.Println("Error occured in dock module when unmount share:", err)
		return &pb.Response{}, err
	}

	return &pb.Response{
		Status:  "Success",
		Message: result,
	}, nil
}
