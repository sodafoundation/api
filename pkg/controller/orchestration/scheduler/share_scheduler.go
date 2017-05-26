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
This module implements the policy-based scheduling by parsing storage
profiles configured by admin.

*/

package scheduler

import (
	api "github.com/opensds/opensds/pkg/api/v1"
	"github.com/opensds/opensds/pkg/grpc/dock/client"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
)

type ShareScheduler struct {
	DesiredProfile *api.StorageProfile
}

func (ss *ShareScheduler) CreateShare(sr *pb.ShareRequest) (*pb.Response, error) {
	sr.ResourceType = ss.DesiredProfile.BackendDriver
	return client.CreateShare(sr)
}

func (ss *ShareScheduler) GetShare(sr *pb.ShareRequest) (*pb.Response, error) {
	sr.ResourceType = ss.DesiredProfile.BackendDriver
	return client.GetShare(sr)
}

func (ss *ShareScheduler) ListShares(sr *pb.ShareRequest) (*pb.Response, error) {
	sr.ResourceType = ss.DesiredProfile.BackendDriver
	return client.ListShares(sr)
}

func (ss *ShareScheduler) DeleteShare(sr *pb.ShareRequest) (*pb.Response, error) {
	sr.ResourceType = ss.DesiredProfile.BackendDriver
	return client.DeleteShare(sr)
}

func (ss *ShareScheduler) AttachShare(sr *pb.ShareRequest) (*pb.Response, error) {
	sr.ResourceType = ss.DesiredProfile.BackendDriver
	return client.AttachShare(sr)
}

func (ss *ShareScheduler) DetachShare(sr *pb.ShareRequest) (*pb.Response, error) {
	sr.ResourceType = ss.DesiredProfile.BackendDriver
	return client.DetachShare(sr)
}
