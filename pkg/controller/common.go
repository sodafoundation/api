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
	"log"
	"reflect"

	"github.com/opensds/opensds/pkg/db"
	api "github.com/opensds/opensds/pkg/model"
)

func SearchProfile(c db.Client, prfId string) (*api.ProfileSpec, error) {
	prfs, err := c.ListProfiles()
	if err != nil {
		log.Println("[Error] When list profiles:", err)
		return &api.ProfileSpec{}, err
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

	return &api.ProfileSpec{}, errors.New("Can not find default profile in db!")
}

func SearchSupportedPool(c db.Client, tags map[string]string) (*api.StoragePoolSpec, error) {
	pols, err := c.ListPools()
	if err != nil {
		log.Println("[Error] When list pool resources in db:", err)
		return &api.StoragePoolSpec{}, err
	}

	// Find if the desired storage tags are contained in any profile
	for _, pol := range *pols {
		var isSupported = true
		for tag := range tags {
			if !Contained(tag, pol.Parameters) {
				isSupported = false
				break
			}
			if pol.Parameters[tag] != "true" {
				isSupported = false
				break
			}
		}
		if isSupported {
			return &pol, nil
		}
	}

	return &api.StoragePoolSpec{}, errors.New("No pool resource supported!")
}

func SearchDockByPool(c db.Client, pol *api.StoragePoolSpec) (*api.DockSpec, error) {
	dcks, err := c.ListDocks()
	if err != nil {
		log.Println("[Error] When list dock resources in db:", err)
		return &api.DockSpec{}, err
	}

	for _, dck := range *dcks {
		if dck.GetId() == pol.GetDockId() {
			return &dck, nil
		}
	}
	return &api.DockSpec{}, errors.New("No dock resource supported!")
}

func SearchDockByVolume(c db.Client, volId string) (*api.DockSpec, error) {
	vol, err := c.GetVolume(volId)
	if err != nil {
		log.Printf("[Error] When get volume %s in db: %v\n", volId, err)
	}
	pols, err := c.ListPools()
	if err != nil {
		log.Println("[Error] When list pool resources in db:", err)
		return &api.DockSpec{}, err
	}
	dcks, err := c.ListDocks()
	if err != nil {
		log.Println("[Error] When list dock resources in db:", err)
		return &api.DockSpec{}, err
	}

	for _, pol := range *pols {
		if pol.GetId() == vol.GetPoolId() {
			for _, dck := range *dcks {
				if dck.GetId() == pol.GetDockId() {
					return &dck, nil
				}
			}
			return &api.DockSpec{}, errors.New("No dock resource supported!")
		}
	}

	return &api.DockSpec{}, errors.New("No pool resource supported!")
}

func Contained(obj, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	default:
		return false
	}
	return false
}
