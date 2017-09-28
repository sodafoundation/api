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

package controller

import (
	"errors"

	log "github.com/golang/glog"

	"github.com/opensds/opensds/pkg/db"
	api "github.com/opensds/opensds/pkg/model"
)

func GetProfile(prfId string, storBox db.Client) (*api.ProfileSpec, error) {
	prfs, err := storBox.ListProfiles()
	if err != nil {
		log.Error("When list profiles:", err)
		return nil, err
	}

	// If a user doesn't specify profile id, then a default profile will be
	// automatically assigned.
	if prfId == "" {
		for _, prf := range *prfs {
			if prf.GetName() == "default" {
				return &prf, nil
			}
		}
	} else {
		for _, prf := range *prfs {
			if prf.GetId() == prfId {
				return &prf, nil
			}
		}
	}

	return nil, errors.New("Can not find default profile in db!")
}
