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
	"errors"
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
	"github.com/opensds/opensds/pkg/controller/orchestration/policyengine"
	"github.com/opensds/opensds/pkg/grpc/dock/client"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
)

const (
	CREATE_LIFECIRCLE_FLAG = 1
	GET_LIFECIRCLE_FLAG    = 2
	LIST_LIFECIRCLE_FLAG   = 3
	DELETE_LIFECIRCLE_FLAG = 4
)

func init() {
	policyengine.Init()
}

type VolumeScheduler struct {
	DesiredProfile *api.StorageProfile
}

func (vs *VolumeScheduler) CreateVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	if vs.DesiredProfile.Name != "" || len(vs.DesiredProfile.StorageTags) != 0 {
		if !policyengine.IsProfileSupported(vs.DesiredProfile) {
			log.Printf("[Error] Profile %+v not supported\n", vs.DesiredProfile)
			return &pb.Response{}, errors.New("Profile not supported")
		}
	}

	st := policyengine.NewStorageTag(vs.DesiredProfile.StorageTags, CREATE_LIFECIRCLE_FLAG)
	vr.ResourceType = vs.DesiredProfile.BackendDriver
	vr.DockId = getDockId(vs.DesiredProfile.BackendDriver)

	result, err := client.CreateVolume(vr)
	if err != nil {
		return &pb.Response{}, err
	}

	var errChan = make(chan error, 1)
	go policyengine.ExecuteAsyncPolicy(vr, st, result.GetMessage(), errChan)

	return result, nil
}

func (vs *VolumeScheduler) GetVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	vr.ResourceType = vs.DesiredProfile.BackendDriver
	vr.DockId = getDockId(vs.DesiredProfile.BackendDriver)
	return client.GetVolume(vr)
}

func (vs *VolumeScheduler) ListVolumes(vr *pb.VolumeRequest) (*pb.Response, error) {
	vr.ResourceType = vs.DesiredProfile.BackendDriver
	vr.DockId = getDockId(vs.DesiredProfile.BackendDriver)
	return client.ListVolumes(vr)
}

func (vs *VolumeScheduler) DeleteVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	if vs.DesiredProfile.Name != "" || len(vs.DesiredProfile.StorageTags) != 0 {
		if !policyengine.IsProfileSupported(vs.DesiredProfile) {
			log.Printf("[Error] Profile %+v not supported\n", vs.DesiredProfile)
			return &pb.Response{}, errors.New("Profile not supported")
		}

		st := policyengine.NewStorageTag(vs.DesiredProfile.StorageTags, DELETE_LIFECIRCLE_FLAG)
		vr.ResourceType = vs.DesiredProfile.BackendDriver
		vr.DockId = getDockId(vs.DesiredProfile.BackendDriver)

		var errChan = make(chan error, 1)
		go policyengine.ExecuteAsyncPolicy(vr, st, "", errChan)

		if err := <-errChan; err != nil {
			log.Println("[Error] When execute async policy:", err)
			return &pb.Response{}, err
		}
	}

	vr.ResourceType = vs.DesiredProfile.BackendDriver
	vr.DockId = getDockId(vs.DesiredProfile.BackendDriver)
	return client.DeleteVolume(vr)
}

func (vs *VolumeScheduler) AttachVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	vr.ResourceType = vs.DesiredProfile.BackendDriver
	vr.DockId = getDockId(vs.DesiredProfile.BackendDriver)
	return client.AttachVolume(vr)
}

func (vs *VolumeScheduler) DetachVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	vr.ResourceType = vs.DesiredProfile.BackendDriver
	vr.DockId = getDockId(vs.DesiredProfile.BackendDriver)
	return client.DetachVolume(vr)
}

func (vs *VolumeScheduler) MountVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	vr.ResourceType = vs.DesiredProfile.BackendDriver
	vr.DockId = getDockId(vs.DesiredProfile.BackendDriver)
	return client.MountVolume(vr)
}

func (vs *VolumeScheduler) UnmountVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	vr.ResourceType = vs.DesiredProfile.BackendDriver
	vr.DockId = getDockId(vs.DesiredProfile.BackendDriver)
	return client.UnmountVolume(vr)
}

func (vs *VolumeScheduler) CreateVolumeSnapshot(vr *pb.VolumeRequest) (*pb.Response, error) {
	vr.ResourceType = vs.DesiredProfile.BackendDriver
	vr.DockId = getDockId(vs.DesiredProfile.BackendDriver)
	return client.CreateVolumeSnapshot(vr)
}

func (vs *VolumeScheduler) GetVolumeSnapshot(vr *pb.VolumeRequest) (*pb.Response, error) {
	vr.ResourceType = vs.DesiredProfile.BackendDriver
	vr.DockId = getDockId(vs.DesiredProfile.BackendDriver)
	return client.GetVolumeSnapshot(vr)
}

func (vs *VolumeScheduler) ListVolumeSnapshots(vr *pb.VolumeRequest) (*pb.Response, error) {
	vr.ResourceType = vs.DesiredProfile.BackendDriver
	vr.DockId = getDockId(vs.DesiredProfile.BackendDriver)
	return client.ListVolumeSnapshots(vr)
}

func (vs *VolumeScheduler) DeleteVolumeSnapshot(vr *pb.VolumeRequest) (*pb.Response, error) {
	vr.ResourceType = vs.DesiredProfile.BackendDriver
	vr.DockId = getDockId(vs.DesiredProfile.BackendDriver)
	return client.DeleteVolumeSnapshot(vr)
}
