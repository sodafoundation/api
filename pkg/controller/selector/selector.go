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
	"sync"

	log "github.com/golang/glog"

	"github.com/opensds/opensds/pkg/db"
	api "github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/opa"
)

type InputData interface{}

type Selector interface {
	SelectSupportedPool(prf *api.ProfileSpec) (*api.StoragePoolSpec, error)

	SelectDock(input interface{}) (*api.DockSpec, error)
}

type selector struct {
	storBox db.Client
	lock    sync.Mutex
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

func (s *selector) SelectSupportedPool(prf *api.ProfileSpec) (*api.StoragePoolSpec, error) {
	prfs, err := s.storBox.ListProfiles()
	if err != nil {
		log.Error("When list profiles resources in db:", err)
		return nil, err
	}
	pols, err := s.storBox.ListPools()
	if err != nil {
		log.Error("When list pool resources in db:", err)
		return nil, err
	}

	if err = s.injectData(&prfs, &pols); err != nil {
		return nil, err
	}

	return opa.GetPoolData(prf.GetName())
}

func (s *selector) SelectDock(input interface{}) (*api.DockSpec, error) {
	dcks, err := s.storBox.ListDocks()
	if err != nil {
		log.Error("When list dock resources in db:", err)
		return nil, err
	}
	pols, err := s.storBox.ListPools()
	if err != nil {
		log.Error("When list pool resources in db:", err)
		return nil, err
	}

	if err = s.injectData(&dcks, &pols); err != nil {
		return nil, err
	}

	var polId string
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

		polId = vol.GetPoolId()
	case *api.StoragePoolSpec:
		polId = input.(*api.StoragePoolSpec).GetId()
	}

	return opa.GetDockData(polId)
}

func (s *selector) injectData(input ...InputData) error {
	var err error

	log.Info("Waiting for injecting data...")
	s.lock.Lock()
	defer s.lock.Unlock()
	log.Info("Start injecting data...")

	for _, v := range input {
		if err = opa.RegisterData(v); err != nil {
			log.Error("Fail to register data in opa:", err)
			return err
		}
	}

	return nil
}
