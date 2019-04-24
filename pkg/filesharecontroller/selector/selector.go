// Copyright (c) 2019 OpenSDS Authors. All Rights Reserved.
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
// ## Refactor work will need to be done to consolidate the selector for block and file share
/*
This module implements the policy-based scheduling by parsing storage
profiles configured by admin.

*/

package selector

import (
	"errors"
	"strconv"

	log "github.com/golang/glog"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
)

// Selector is an interface that exposes some operation of different selectors.
type Selector interface {
	SelectSupportedPoolForFileShare(*model.FileShareSpec) (*model.StoragePoolSpec, error)
}

type selector struct{}

// NewSelector method creates a new selector structure and return its pointer.
func NewSelector() Selector {
	return &selector{}
}

// SelectSupportedPoolForFileShare
func (s *selector) SelectSupportedPoolForFileShare(in *model.FileShareSpec) (*model.StoragePoolSpec, error) {
	var prf *model.ProfileSpec
	var err error

	if in.ProfileId == "" {
		log.Warning("Use default profile when user doesn't specify profile.")
		prf, err = db.C.GetDefaultProfile(c.NewAdminContext())
	} else {
		prf, err = db.C.GetProfile(c.NewAdminContext(), in.ProfileId)
	}
	if err != nil {
		log.Error("Get profile failed: ", err)
		return nil, err
	}
	pools, err := db.C.ListPools(c.NewAdminContext())
	if err != nil {
		log.Error("When list pools in resources SelectSupportedPool: ", err)
		return nil, err
	}

	// Generate filter request according to the rules defined in profile.
	fltRequest := func(prf *model.ProfileSpec, in *model.FileShareSpec) map[string]interface{} {
		var filterRequest map[string]interface{}
		if !prf.CustomProperties.IsEmpty() {
			filterRequest = prf.CustomProperties
		} else {
			filterRequest = make(map[string]interface{})
		}
		// Insert some basic rules.
		filterRequest["freeCapacity"] = ">= " + strconv.Itoa(int(in.Size))
		if in.AvailabilityZone != "" {
			filterRequest["availabilityZone"] = in.AvailabilityZone
		} else {
			filterRequest["availabilityZone"] = "default"
		}
		if in.PoolId != "" {
			filterRequest["id"] = in.PoolId
		}
		// Insert some rules of provisioning properties.
		if pp := prf.ProvisioningProperties; !pp.IsEmpty() {
			if ds := pp.DataStorage; !ds.IsEmpty() {
				filterRequest["extras.dataStorage.isSpaceEfficient"] =
					"<is> " + strconv.FormatBool(ds.IsSpaceEfficient)
				if ds.ProvisioningPolicy != "" {
					filterRequest["extras.dataStorage.provisioningPolicy"] =
						ds.ProvisioningPolicy
				}
				if ds.RecoveryTimeObjective != 0 {
					filterRequest["extras.dataStorage.recoveryTimeObjective"] =
						"<= " + strconv.Itoa(int(ds.RecoveryTimeObjective))
				}
			}
			if ic := pp.IOConnectivity; !ic.IsEmpty() {
				if ic.AccessProtocol != "" {
					filterRequest["extras.ioConnectivity.accessProtocol"] =
						ic.AccessProtocol
				}
				if ic.MaxIOPS != 0 {
					filterRequest["extras.ioConnectivity.maxIOPS"] =
						">= " + strconv.Itoa(int(ic.MaxIOPS))
				}
				if ic.MaxBWS != 0 {
					filterRequest["extras.ioConnectivity.maxBWS"] =
						">= " + strconv.Itoa(int(ic.MaxBWS))
				}
			}
		}

		return filterRequest
	}(prf, in)

	supportedPools, err := SelectSupportedPools(1, fltRequest, pools)
	if err != nil {
		log.Error("Filter supported pools failed: ", err)
		return nil, err
	}
	// Now, we just return the first supported pool which will be improved in
	// the future.
	return supportedPools[0], nil
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
