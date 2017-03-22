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

func CreateShare(sr *pb.ShareRequest) (*pb.Response, error) {
	result, err := dock.CreateShare(sr.GetResoureType(),
		sr.GetName(),
		sr.GetShareType(),
		sr.GetShareProto(),
		sr.GetSize())
	if err != nil {
		log.Println("Error occured in adapter module when create share:", err)
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

func GetShare(sr *pb.ShareRequest) (*pb.Response, error) {
	result, err := dock.GetShare(sr.GetResoureType(), sr.GetId())
	if err != nil {
		log.Println("Error occured in adapter module when get share:", err)
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

func ListShares(sr *pb.ShareRequest) (*pb.Response, error) {
	result, err := dock.GetAllShares(sr.GetResoureType(),
		sr.GetAllowDetails())
	if err != nil {
		log.Println("Error occured in adapter module when list shares:", err)
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

func DeleteShare(sr *pb.ShareRequest) (*pb.Response, error) {
	result, err := dock.DeleteShare(sr.GetResoureType(),
		sr.GetId())
	if err != nil {
		log.Println("Error occured in adapter module when delete share:", err)
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
