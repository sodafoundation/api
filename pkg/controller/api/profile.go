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
This module implements the operation of profile resources.

*/

package api

import (
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
	"github.com/opensds/opensds/pkg/controller/metadata/profile"
)

type ProfileDeliver interface {
	getProfile(name string) (*api.StorageProfile, error)

	listProfiles() (*api.StorageProfiles, error)
}

// ProfileRequest is a structure for all properties of
// a profile request
type ProfileRequest struct{}

func (pr ProfileRequest) getProfile(name string) (*api.StorageProfile, error) {
	return profile.GetProfile(name)
}

func (pr ProfileRequest) listProfiles() (*api.StorageProfiles, error) {
	return profile.ListProfiles()
}

func GetProfile(pd ProfileDeliver, name string) (api.StorageProfile, error) {
	pf, err := pd.getProfile(name)
	if err != nil {
		log.Printf("Get profile %s error: %v\n", name, err)
		return api.StorageProfile{}, err
	}

	return *pf, nil
}

func ListProfiles(pd ProfileDeliver) (api.StorageProfiles, error) {
	pfs, err := pd.listProfiles()
	if err != nil {
		log.Printf("List profiles error: %v\n", err)
		return api.StorageProfiles{}, err
	}

	return *pfs, nil
}
