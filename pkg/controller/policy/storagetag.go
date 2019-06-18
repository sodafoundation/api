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

package policy

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/utils"
)

const (
	POLICY_TYPE_MAPPING_TABLE = `{
		"diskType": "feature",
		"thinProvision": "feature",
		"highAvailability": "feature",
		"intervalSnapshot": "operation",
		"deleteSnapshotPolicy": "operation"
	}`
	POLICY_LIFECIRCLE_TABLE = `{
		"thinProvision": 1,
		"highAvailability": 1,
		"intervalSnapshot": 1,
		"deleteSnapshotPolicy": 4
	}`
)

// PolicyTypeMappingTable
var PolicyTypeMappingTable map[string]string

// PolicyLifecircleTable
var PolicyLifecircleTable map[string]int

func init() {
	json.Unmarshal([]byte(POLICY_TYPE_MAPPING_TABLE), &PolicyTypeMappingTable)
	json.Unmarshal([]byte(POLICY_LIFECIRCLE_TABLE), &PolicyLifecircleTable)
}

// IsStorageTagSupported
func IsStorageTagSupported(tags map[string]string) bool {
	for key := range tags {
		if PolicyTypeMappingTable[key] != "operation" {
			return false
		}
	}
	return true
}

// FindPolicyType
func FindPolicyType(policy string) (string, error) {
	if !utils.Contained(policy, PolicyTypeMappingTable) {
		return "", errors.New("The policy type of " + policy + " not supported")
	}

	return PolicyTypeMappingTable[policy], nil
}

// StorageTag
type StorageTag struct {
	syncTag  map[string]interface{}
	asyncTag map[string]string
}

// NewStorageTag
func NewStorageTag(tags map[string]interface{}, flag int) *StorageTag {
	var st = &StorageTag{
		syncTag:  make(map[string]interface{}),
		asyncTag: make(map[string]string),
	}

	// Screen storage tags through life circle flag
	for key := range tags {
		if flag != PolicyLifecircleTable[key] {
			delete(tags, key)
		}
	}
	// Divide all tags into sync and async part
	for key := range tags {
		pType, err := FindPolicyType(key)
		if err != nil {
			log.Error("When parse storage tag:", err)
		}
		switch pType {
		case "feature":
			st.syncTag[key] = tags[key]
		case "operation":
			st.asyncTag[key] = fmt.Sprint(tags[key])
		}
	}
	return st
}

// GetSyncTag
func (st *StorageTag) GetSyncTag() map[string]interface{} {
	return st.syncTag
}

// GetAsyncTag
func (st *StorageTag) GetAsyncTag() map[string]string {
	return st.asyncTag
}
