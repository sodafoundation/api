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

	"github.com/opensds/opensds/pkg/controller/orchestration/scheduler"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
)

func CreateShare(sr *pb.ShareRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(sr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	ss := &scheduler.ShareScheduler{
		DesiredProfile: profile,
	}
	return ss.CreateShare(sr)
}

func GetShare(sr *pb.ShareRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(sr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	ss := &scheduler.ShareScheduler{
		DesiredProfile: profile,
	}
	return ss.GetShare(sr)
}

func ListShares(sr *pb.ShareRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(sr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	ss := &scheduler.ShareScheduler{
		DesiredProfile: profile,
	}
	return ss.ListShares(sr)
}

func DeleteShare(sr *pb.ShareRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(sr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	ss := &scheduler.ShareScheduler{
		DesiredProfile: profile,
	}
	return ss.DeleteShare(sr)
}

func AttachShare(sr *pb.ShareRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(sr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	ss := &scheduler.ShareScheduler{
		DesiredProfile: profile,
	}
	return ss.AttachShare(sr)
}

func DetachShare(sr *pb.ShareRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(sr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	ss := &scheduler.ShareScheduler{
		DesiredProfile: profile,
	}
	return ss.DetachShare(sr)
}

func MountShare(sr *pb.ShareRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(sr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	ss := &scheduler.ShareScheduler{
		DesiredProfile: profile,
	}
	return ss.MountShare(sr)
}

func UnmountShare(sr *pb.ShareRequest) (*pb.Response, error) {
	if err := json.Unmarshal([]byte(sr.GetStorageProfile()), profile); err != nil {
		return &pb.Response{}, err
	}
	ss := &scheduler.ShareScheduler{
		DesiredProfile: profile,
	}
	return ss.UnmountShare(sr)
}
