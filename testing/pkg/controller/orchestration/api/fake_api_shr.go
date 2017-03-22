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

func CreateShare(sr *pb.ShareRequest) (*pb.Response, error) {
	return grpcapi.CreateShare(sr), nil
}

func GetShare(sr *pb.ShareRequest) (*pb.Response, error) {
	return grpcapi.GetShare(sr), nil
}

func ListShares(sr *pb.ShareRequest) (*pb.Response, error) {
	return grpcapi.ListShares(sr), nil
}

func DeleteShare(sr *pb.ShareRequest) (*pb.Response, error) {
	return grpcapi.DeleteShare(sr), nil
}
