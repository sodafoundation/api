// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
This module implements the policy-based scheduling by parsing storage
profiles configured by admin.

*/

package selector

import (
	"errors"

	log "github.com/golang/glog"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
)

// Selector is an interface that exposes some operation of different selectors.
type Selector interface {
	SelectSupportedPool(map[string]interface{}) (*model.StoragePoolSpec, error)
	SelectSupportedPoolForVG(*model.VolumeGroupSpec) (*model.StoragePoolSpec, error)
}

type selector struct{}

// NewSelector method creates a new selector structure and return its pointer.
func NewSelector() Selector {
	return &selector{}
}

// SelectSupportedPool
func (s *selector) SelectSupportedPool(tags map[string]interface{}) (*model.StoragePoolSpec, error) {
	pools, err := db.C.ListPools(c.NewAdminContext())
	if err != nil {
		log.Error("When list pools in resources SelectSupportedPool: ", err)
		return nil, err
	}

	supportedPools, err := SelectSupportedPools(1, tags, pools)
	if err != nil {
		log.Error("Filter supported pools failed: ", err)
		return nil, err
	}

	// Now, we just return the first supported pool which will be improved in
	// the future.
	return supportedPools[0], nil

}

func (s *selector) SelectSupportedPoolForVG(in *model.VolumeGroupSpec) (*model.StoragePoolSpec, error) {
	var profiles []*model.ProfileSpec
	for _, profId := range in.Profiles {
		profile, err := db.C.GetProfile(c.NewAdminContext(), profId)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	pools, err := db.C.ListPools(c.NewAdminContext())
	if err != nil {
		return nil, err
	}

	var filterRequest map[string]interface{}

	for _, pool := range pools {

		var poolIsFound = true

		for _, profile := range profiles {

			if profile.Extras != nil {
				filterRequest = profile.Extras
			} else {
				filterRequest = make(map[string]interface{})
			}
			filterRequest["availabilityZone"] = in.AvailabilityZone

			isAvailable, err := IsAvailablePool(filterRequest, pool)
			if nil != err {
				return nil, err
			}
			if !isAvailable {
				poolIsFound = false
				break
			}
		}

		if poolIsFound {
			return pool, nil
		}
	}

	return nil, errors.New("No valid pool found for group.")
}

// SelectSupportedPools ...
func SelectSupportedPools(maxNum int, filterReq map[string]interface{}, pools []*model.StoragePoolSpec) ([]*model.StoragePoolSpec, error) {
	supportedPools := []*model.StoragePoolSpec{}

	for _, pool := range pools {
		if len(supportedPools) < maxNum {
			isAvailable, err := IsAvailablePool(filterReq, pool)
			if nil != err {
				return nil, err
			}

			if isAvailable {
				supportedPools = append(supportedPools, pool)
			}
		} else {
			break
		}
	}

	if len(supportedPools) > 0 {
		return supportedPools, nil
	}

	return nil, errors.New("no available pool to meet user's requirement")
}
