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
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
)

// Selector is an interface that exposes some operation of different selectors.
type Selector interface {
	SelectSupportedPool(tags map[string]interface{}) (*model.StoragePoolSpec, error)
}

type selector struct{}

var filterChain Filter

func init() {
	diskTypeFilter := &DiskTypeFilter{}
	compressFilter := &CompressFilter{Next: diskTypeFilter}
	dedupeFilter := &DedupeFilter{Next: compressFilter}
	thinFilter := &ThinFilter{Next: dedupeFilter}
	capacityFilter := &CapacityFilter{Next: thinFilter}
	azFilter := &AZFilter{Next: capacityFilter}

	filterChain = azFilter
}

// NewSelector method creates a new selector structure and return its pointer.
func NewSelector() Selector {
	return &selector{}
}

// SelectSupportedPool
func (s *selector) SelectSupportedPool(tags map[string]interface{}) (*model.StoragePoolSpec, error) {
	pools, err := db.C.ListPools()
	if err != nil {
		log.Error("When list pools in resources SelectSupportedPool: ", err)
		return nil, err
	}

	supportedPools, err := filterChain.Handle(tags, pools)
	if err != nil {
		log.Error("Filter supported pools failed: ", err)
		return nil, err
	}

	// Now, we just return the firt supported pool which will be improved in
	// the future.
	return supportedPools[0], nil

}
