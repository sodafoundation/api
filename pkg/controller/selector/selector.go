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

package selector

import (
	"errors"

	log "github.com/golang/glog"

	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
)

type Selector interface {
	SelectSupportedPool(tags map[string]string) (*model.StoragePoolSpec, error)

	SelectDock(input interface{}) (*model.DockSpec, error)
}

type selector struct {
	storBox db.Client
}

func NewSelector() Selector {
	return &selector{
		storBox: db.C,
	}
}

func NewFakeSelector() Selector {
	return &selector{
		storBox: db.NewFakeDbClient(),
	}
}

func (s *selector) SelectSupportedPool(tags map[string]string) (*model.StoragePoolSpec, error) {
	pols, err := s.storBox.ListPools()
	if err != nil {
		log.Error("When list pool resources in db:", err)
		return nil, err
	}

	// Find if the desired storage tags are contained in any profile
	for _, pol := range pols {
		var isSupported = true
		for tag := range tags {
			if !utils.Contained(tag, pol.Parameters) {
				isSupported = false
				break
			}
			if pol.Parameters[tag] != "true" {
				isSupported = false
				break
			}
		}
		if isSupported {
			return pol, nil
		}
	}

	return nil, errors.New("No pool resource supported!")
}

func (s *selector) SelectDock(input interface{}) (*model.DockSpec, error) {
	dcks, err := s.storBox.ListDocks()
	if err != nil {
		log.Error("When list dock resources in db:", err)
		return nil, err
	}

	var pol *model.StoragePoolSpec

	switch input.(type) {
	case string:
		// If user specifies a volume id, then the selector will find the
		// storage pool by calling database.
		volID := input.(string)
		vol, err := s.storBox.GetVolume(volID)
		if err != nil {
			log.Errorf("When get volume %v in db: %v\n", input, err)
			return nil, err
		}

		pol, err = s.storBox.GetPool(vol.GetPoolId())
		if err != nil {
			log.Errorf("When get pool %s in db: %v\n", vol.GetPoolId(), err)
			return nil, err
		}
	case *model.StoragePoolSpec:
		pol = input.(*model.StoragePoolSpec)
	}

	for _, dck := range dcks {
		if dck.GetId() == pol.GetDockId() {
			return dck, nil
		}
	}
	return nil, errors.New("No dock resource supported!")
}
