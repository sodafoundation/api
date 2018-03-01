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

package selector

import (
	"errors"

	"github.com/opensds/opensds/pkg/model"
)

type Filter interface {
	Handle(request map[string]interface{}, pools []*model.StoragePoolSpec) ([]*model.StoragePoolSpec, error)
}

type AZFilter struct {
	Next Filter
}

func (filter *AZFilter) Handle(request map[string]interface{}, pools []*model.StoragePoolSpec) ([]*model.StoragePoolSpec, error) {
	availabilityZone, ok := request["availabilityZone"]
	if !ok {
		return nil, errors.New("The request doesn't contain availability zone property when filter pools in AZFilter.")
	}

	availablePools := []*model.StoragePoolSpec{}
	for _, pool := range pools {
		if availabilityZone == pool.AvailabilityZone {
			availablePools = append(availablePools, pool)
		}
	}

	if len(availablePools) == 0 {
		return nil, errors.New("No available pools to meet user's availability zone requirement.")
	} else if filter.Next != nil {
		return filter.Next.Handle(request, availablePools)
	}
	return availablePools, nil
}

type CapacityFilter struct {
	Next Filter
}

func (filter *CapacityFilter) Handle(request map[string]interface{}, pools []*model.StoragePoolSpec) ([]*model.StoragePoolSpec, error) {

	size, ok := request["size"]
	if !ok {
		return nil, errors.New("The request doesn't contain size property when filter pools in CapacityFilter.")
	}

	availablePools := []*model.StoragePoolSpec{}
	for _, pool := range pools {
		if val, ok := size.(int64); ok {
			if val <= pool.FreeCapacity {
				availablePools = append(availablePools, pool)
			}
		} else {
			return nil, errors.New("Invalid parameter when filter pools in CapacityFilter.")
		}

	}

	if len(availablePools) == 0 {
		return nil, errors.New("No available pools to meet user's capacity requirement.")
	} else if filter.Next != nil {
		return filter.Next.Handle(request, availablePools)
	}
	return availablePools, nil
}

type ThinFilter struct {
	Next Filter
}

func (filter *ThinFilter) Handle(request map[string]interface{}, pools []*model.StoragePoolSpec) ([]*model.StoragePoolSpec, error) {

	availablePools := []*model.StoragePoolSpec{}
	thinRequired, ok := request["thin"]
	if !ok {
		availablePools = pools
	} else {
		for _, pool := range pools {
			thinSupport, ok := pool.Extras["thin"]
			if !ok {
				thinSupport = false
			}

			if thinRequired == thinSupport {
				availablePools = append(availablePools, pool)
			}
		}
	}

	if len(availablePools) == 0 {
		return nil, errors.New("No available pools to meet user's thin requirement.")
	} else if filter.Next != nil {
		return filter.Next.Handle(request, availablePools)
	}
	return availablePools, nil
}

type CompressFilter struct {
	Next Filter
}

func (filter *CompressFilter) Handle(request map[string]interface{}, pools []*model.StoragePoolSpec) ([]*model.StoragePoolSpec, error) {

	availablePools := []*model.StoragePoolSpec{}
	compressRequired, ok := request["compress"]
	if !ok {
		availablePools = pools
	} else {
		for _, pool := range pools {
			compressSupport, ok := pool.Extras["compress"]
			if !ok {
				compressSupport = false
			}

			if compressRequired == compressSupport {
				availablePools = append(availablePools, pool)
			}
		}
	}

	if len(availablePools) == 0 {
		return nil, errors.New("No available pools to meet user's compress requirement.")
	} else if filter.Next != nil {
		return filter.Next.Handle(request, availablePools)
	}
	return availablePools, nil
}

type DedupeFilter struct {
	Next Filter
}

func (filter *DedupeFilter) Handle(request map[string]interface{}, pools []*model.StoragePoolSpec) ([]*model.StoragePoolSpec, error) {

	availablePools := []*model.StoragePoolSpec{}
	dedupeRequired, ok := request["dedupe"]
	if !ok {
		availablePools = pools
	} else {
		for _, pool := range pools {
			dedupeSupport, ok := pool.Extras["dedupe"]
			if !ok {
				dedupeSupport = false
			}

			if dedupeRequired == dedupeSupport {
				availablePools = append(availablePools, pool)
			}
		}
	}

	if len(availablePools) == 0 {
		return nil, errors.New("No available pools to meet user's dedupe requirement.")
	} else if filter.Next != nil {
		return filter.Next.Handle(request, availablePools)
	}
	return availablePools, nil
}

type DiskTypeFilter struct {
	Next Filter
}

func (filter *DiskTypeFilter) Handle(request map[string]interface{}, pools []*model.StoragePoolSpec) ([]*model.StoragePoolSpec, error) {

	availablePools := []*model.StoragePoolSpec{}
	diskTypeRequired, ok := request["diskType"]
	if !ok {
		availablePools = pools
	} else {
		for _, pool := range pools {
			diskTypeSupport, ok := pool.Extras["diskType"]
			if !ok {
				diskTypeSupport = false
			}

			if diskTypeRequired == diskTypeSupport {
				availablePools = append(availablePools, pool)
			}
		}
	}

	if len(availablePools) == 0 {
		return nil, errors.New("No available pools to meet user's diskType requirement.")
	} else if filter.Next != nil {
		return filter.Next.Handle(request, availablePools)
	}
	return availablePools, nil
}
