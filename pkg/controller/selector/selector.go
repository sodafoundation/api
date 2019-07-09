// Copyright 2017 The OpenSDS Authors.
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
	"strconv"

	log "github.com/golang/glog"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
)

// Selector is an interface that exposes some operation of different selectors.
type Selector interface {
	SelectSupportedPoolForVolume(*model.VolumeSpec) (*model.StoragePoolSpec, error)
	SelectSupportedPoolForVG(*model.VolumeGroupSpec) (*model.StoragePoolSpec, error)
	SelectSupportedPoolForFileShare(*model.FileShareSpec) (*model.StoragePoolSpec, error)
}

type selector struct{}

// NewSelector method creates a new selector structure and return its pointer.
func NewSelector() Selector {
	return &selector{}
}

// SelectSupportedPoolForVolume
func (s *selector) SelectSupportedPoolForVolume(vol *model.VolumeSpec) (*model.StoragePoolSpec, error) {
	var prf *model.ProfileSpec
	var err error
	if vol.ProfileId == "" {
		log.Warning("Use default profile when user doesn't specify profile.")
		prf, err = db.C.GetDefaultProfile(c.NewAdminContext())
	} else {
		prf, err = db.C.GetProfile(c.NewAdminContext(), vol.ProfileId)
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
	fltRequest := func(prf *model.ProfileSpec, vol *model.VolumeSpec) map[string]interface{} {
		var filterRequest map[string]interface{}
		filterRequest = prf.CustomProperties.GetCapabilitiesProperties()

		if v, ok := filterRequest["multiAttach"]; ok && v == "<is> true" {
			log.Info("change volume multiAttach flag to be true.")
			vol.MultiAttach = true
		}

		// Insert some basic rules.
		filterRequest["freeCapacity"] = ">= " + strconv.Itoa(int(vol.Size))
		if vol.AvailabilityZone != "" {
			filterRequest["availabilityZone"] = vol.AvailabilityZone
		} else {
			filterRequest["availabilityZone"] = "default"
		}
		if vol.PoolId != "" {
			filterRequest["id"] = vol.PoolId
		}
		filterRequest["storageType"] = prf.StorageType
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
		// Insert some rules of replication properties.
		if rp := prf.ReplicationProperties; !rp.IsEmpty() {
			if dp := rp.DataProtection; !dp.IsEmpty() {
				filterRequest["extras.dataProtection.isIsolated"] =
					"<is> " + strconv.FormatBool(dp.IsIsolated)
				if dp.RecoveryGeographicObject != "" {
					filterRequest["extras.dataProtection.recoveryGeographicObject"] =
						dp.RecoveryGeographicObject
				}
				if dp.RecoveryTimeObjective != "" {
					filterRequest["extras.dataProtection.recoveryTimeObjective"] =
						dp.RecoveryTimeObjective
				}
				if dp.ReplicaType != "" {
					filterRequest["extras.dataProtection.replicaType"] =
						dp.ReplicaType
				}
			}
		}
		return filterRequest
	}(prf, vol)

	log.Infof("The filter request for pool is %v", fltRequest)
	supportedPools, err := SelectSupportedPools(1, fltRequest, pools)
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
			filterRequest = profile.CustomProperties.GetCapabilitiesProperties()
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

// SelectSupportedPoolForFileShare
func (s *selector) SelectSupportedPoolForFileShare(in *model.FileShareSpec) (*model.StoragePoolSpec, error) {
	var prf *model.ProfileSpec
	var err error

	if in.ProfileId == "" {
		log.Warning("use default profile when user doesn't specify profile.")
		prf, err = db.C.GetDefaultProfileFileShare(c.NewAdminContext())
		if err != nil {
			log.Error("get profile failed: ", err)
			return nil, err
		}
	} else {
		prf, err = db.C.GetProfile(c.NewAdminContext(), in.ProfileId)
		if err != nil {
			log.Error("get profile failed: ", err)
			return nil, err
		}
	}

	pools, err := db.C.ListPools(c.NewAdminContext())
	if err != nil {
		log.Error("when list pools in resources SelectSupportedPool: ", err)
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

		filterRequest["storageType"] = prf.StorageType

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
				if !ds.IsEmptyStorageAccessCapability() {
					filterRequest["extras.dataStorage.storageAccessCapability"] =
						ds.StorageAccessCapability
				}
				if ds.MaxFileNameLengthBytes != 0 {
					filterRequest["extras.dataStorage.maxFileNameLengthBytes"] =
						"<= " + strconv.Itoa(int(ds.MaxFileNameLengthBytes))
				}
				if ds.CharacterCodeSet != "" {
					filterRequest["extras.dataStorage.characterCodeSet"] =
						ds.CharacterCodeSet
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
		log.Error("filter supported pools failed: ", err)
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
