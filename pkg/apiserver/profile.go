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

package api

import (
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
	db "github.com/opensds/opensds/pkg/db/api"
)

type ProfileRequestDeliver interface {
	getProfile() (*api.StorageProfile, error)

	listProfiles() (*[]api.StorageProfile, error)
}

// ProfileRequest is a structure for all properties of
// a storage profile request
type ProfileRequest struct {
	Id string `json:"id"`
}

func (pr *ProfileRequest) getProfile() (*api.StorageProfile, error) {
	return db.GetProfile(pr.Id)
}

func (pr *ProfileRequest) listProfiles() (*[]api.StorageProfile, error) {
	return db.ListProfiles()
}

func GetProfile(prd ProfileRequestDeliver) (api.StorageProfile, error) {
	var nullResponse = api.StorageProfile{}

	prf, err := prd.getProfile()
	if err != nil {
		log.Println("Get profile error:", err)
		return nullResponse, err
	}
	return *prf, nil
}

func ListProfiles(prd ProfileRequestDeliver) ([]api.StorageProfile, error) {
	var nullResponses = []api.StorageProfile{}

	prfs, err := prd.listProfiles()
	if err != nil {
		log.Println("List all profiles error:", err)
		return nullResponses, err
	}
	return *prfs, nil
}
