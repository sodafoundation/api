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
This module implements the policy-based orchestration by parsing storage
profiles configured by admin.

*/

package orchestration

import (
	_ "errors"
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
	"github.com/opensds/opensds/pkg/controller/orchestration/policyengine"
	"github.com/opensds/opensds/pkg/controller/orchestration/scheduler"
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

type VolumeOrchestrator struct {
	DesiredProfile *api.StorageProfile
}

func (vo *VolumeOrchestrator) CreateVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	var err error
	vo.DesiredProfile, err = scheduler.SearchProfile(vo.DesiredProfile.Name)
	if err != nil {
		return &pb.Response{}, err
	}
	/*
		if !policyengine.IsStorageTagSupported(vo.DesiredProfile.StorageTags) {
			return &pb.Response{}, errors.New("Storage tags not supported!")
		}
	*/

	st := policyengine.NewStorageTag(vo.DesiredProfile.StorageTags, CREATE_LIFECIRCLE_FLAG)
	polInfo, err := scheduler.SearchSupportedPool(st.GetSyncTag(), vr.GetSize())
	if err != nil {
		log.Printf("[Error] When search supported pool resource:", err)
		return &pb.Response{}, err
	}
	if err = scheduler.UpdateDockInfo(vr, polInfo); err != nil {
		log.Printf("[Error] When update dock in volume request:", err)
		return &pb.Response{}, err
	}

	result, err := client.CreateVolume(vr)
	if err != nil {
		return &pb.Response{}, err
	}

	var errChan = make(chan error, 1)
	go policyengine.ExecuteAsyncPolicy(vr, st, result.GetMessage(), errChan)

	return result, nil
}

func (vo *VolumeOrchestrator) DeleteVolume(vr *pb.VolumeRequest) (*pb.Response, error) {
	var err error
	vo.DesiredProfile, err = scheduler.SearchProfile(vo.DesiredProfile.Name)
	if err != nil {
		return &pb.Response{}, err
	}
	/*
		if !policyengine.IsStorageTagSupported(vo.DesiredProfile.StorageTags) {
			return &pb.Response{}, errors.New("Storage tags not supported!")
		}
	*/

	st := policyengine.NewStorageTag(vo.DesiredProfile.StorageTags, DELETE_LIFECIRCLE_FLAG)
	polInfo, err := scheduler.SearchPoolByVolume(vr.GetVolumeId())
	if err != nil {
		log.Printf("[Error] When search supported pool resource:", err)
		return &pb.Response{}, err
	}
	if err = scheduler.UpdateDockInfo(vr, polInfo); err != nil {
		log.Printf("[Error] When update dock in volume request:", err)
		return &pb.Response{}, err
	}

	var errChan = make(chan error, 1)
	go policyengine.ExecuteAsyncPolicy(vr, st, "", errChan)

	if err := <-errChan; err != nil {
		log.Println("[Error] When execute async policy:", err)
		return &pb.Response{}, err
	}

	return client.DeleteVolume(vr)
}

func (vo *VolumeOrchestrator) CreateVolumeAttachment(vr *pb.VolumeRequest) (*pb.Response, error) {
	polInfo, err := scheduler.SearchPoolByVolume(vr.GetVolumeId())
	if err != nil {
		log.Printf("[Error] When search supported pool resource:", err)
		return &pb.Response{}, err
	}
	if err = scheduler.UpdateDockInfo(vr, polInfo); err != nil {
		log.Printf("[Error] When update dock in volume request:", err)
		return &pb.Response{}, err
	}
	return client.CreateVolumeAttachment(vr)
}

func (vo *VolumeOrchestrator) UpdateVolumeAttachment(vr *pb.VolumeRequest) (*pb.Response, error) {
	polInfo, err := scheduler.SearchPoolByVolume(vr.GetVolumeId())
	if err != nil {
		log.Printf("[Error] When search supported pool resource:", err)
		return &pb.Response{}, err
	}
	if err = scheduler.UpdateDockInfo(vr, polInfo); err != nil {
		log.Printf("[Error] When update dock in volume request:", err)
		return &pb.Response{}, err
	}
	return client.UpdateVolumeAttachment(vr)
}

func (vo *VolumeOrchestrator) DeleteVolumeAttachment(vr *pb.VolumeRequest) (*pb.Response, error) {
	polInfo, err := scheduler.SearchPoolByVolume(vr.GetVolumeId())
	if err != nil {
		log.Printf("[Error] When search supported pool resource:", err)
		return &pb.Response{}, err
	}
	if err = scheduler.UpdateDockInfo(vr, polInfo); err != nil {
		log.Printf("[Error] When update dock in volume request:", err)
		return &pb.Response{}, err
	}
	return client.DeleteVolumeAttachment(vr)
}

func (vo *VolumeOrchestrator) CreateVolumeSnapshot(vr *pb.VolumeRequest) (*pb.Response, error) {
	polInfo, err := scheduler.SearchPoolByVolume(vr.GetVolumeId())
	if err != nil {
		log.Printf("[Error] When search supported pool resource:", err)
		return &pb.Response{}, err
	}
	if err = scheduler.UpdateDockInfo(vr, polInfo); err != nil {
		log.Printf("[Error] When update dock in volume request:", err)
		return &pb.Response{}, err
	}
	return client.CreateVolumeSnapshot(vr)
}

func (vo *VolumeOrchestrator) DeleteVolumeSnapshot(vr *pb.VolumeRequest) (*pb.Response, error) {
	polInfo, err := scheduler.SearchPoolByVolume(vr.GetVolumeId())
	if err != nil {
		log.Printf("[Error] When search supported pool resource:", err)
		return &pb.Response{}, err
	}
	if err = scheduler.UpdateDockInfo(vr, polInfo); err != nil {
		log.Printf("[Error] When update dock in volume request:", err)
		return &pb.Response{}, err
	}
	return client.DeleteVolumeSnapshot(vr)
}
