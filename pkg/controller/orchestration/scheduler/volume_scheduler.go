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
	"encoding/json"
	"errors"
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
	db "github.com/opensds/opensds/pkg/db/api"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
)

func SearchProfile(profileName string) (*api.StorageProfile, error) {
	if profileName == "" {
		profileName = "default"
	}
	profiles, err := db.ListProfiles()
	if err != nil {
		log.Println("[Error] When list profiles:", err)
		return &api.StorageProfile{}, err
	}

	for _, profile := range *profiles {
		if profile.Name == profileName {
			return &profile, nil
		}
	}
	return &api.StorageProfile{}, errors.New("Can not find default profile in db!")
}

func SearchPoolByVolume(volId string) (*api.StoragePool, error) {
	vol, err := db.GetVolume(volId)
	if err != nil {
		log.Printf("[Error] When get volume %s in db: %v\n", volId, err)
	}
	pols, err := db.ListPools()
	if err != nil {
		log.Println("[Error] When list pool resources in db:", err)
		return &api.StoragePool{}, err
	}

	for _, pol := range *pols {
		if pol.Name == vol.PoolName {
			return &pol, nil
		}
	}
	return &api.StoragePool{}, errors.New("No pool resource supported!")
}

func SearchSupportedPool(tags map[string]string, usedCapacity int32) (*api.StoragePool, error) {
	pols, err := db.ListPools()
	if err != nil {
		log.Println("[Error] When list pool resources in db:", err)
		return &api.StoragePool{}, err
	}

	for _, pol := range *pols {
		// Find if the desired storage tags are contained in any profile
		var isSupported = true
		if pol.FreeCapacity < int64(usedCapacity) {
			continue
		}
		for tag := range tags {
			if !Contained(tag, pol.StorageTag) {
				isSupported = false
				break
			}
			if pol.StorageTag[tag] != "true" {
				isSupported = false
				break
			}
		}
		if isSupported {
			return &pol, nil
		}
	}

	return &api.StoragePool{}, errors.New("No pool resource supported!")
}

func UpdateDockInfo(vr *pb.VolumeRequest, polInfo *api.StoragePool) error {
	dcks, err := db.ListDocks()
	if err != nil {
		log.Printf("[Error] When list dock resources in db:", err)
		return err
	}

	for _, dck := range *dcks {
		if dck.Name == polInfo.DockName {
			dckBody, _ := json.Marshal(dck)
			vr.DockInfo, vr.PoolName = string(dckBody), polInfo.Name
			return nil
		}
	}
	return errors.New("No dock resource supported!")
}
