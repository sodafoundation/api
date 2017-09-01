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

package policy

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/opensds/opensds/pkg/utils"
)

const (
	POLICY_TYPE_MAPPING_TABLE = `{
		"iops": "feature",
		"thinProvision": "feature",
		"highAvailability": "feature",
		"intervalSnapshot": "operation",
		"deleteSnapshotPolicy": "operation"
	}`
	POLICY_LIFECIRCLE_TABLE = `{
		"iops": 1,
		"thinProvision": 1,
		"highAvailability": 1,
		"intervalSnapshot": 1,
		"deleteSnapshotPolicy": 4
	}`
)

var PolicyTypeMappingTable map[string]string
var PolicyLifecircleTable map[string]int

func init() {
	json.Unmarshal([]byte(POLICY_TYPE_MAPPING_TABLE), &PolicyTypeMappingTable)
	json.Unmarshal([]byte(POLICY_LIFECIRCLE_TABLE), &PolicyLifecircleTable)
}

func IsStorageTagSupported(tags map[string]string) bool {
	for key := range tags {
		if PolicyTypeMappingTable[key] != "operation" {
			return false
		}
	}
	return true
}

func FindPolicyType(policy string) (string, error) {
	if !utils.Contained(policy, PolicyTypeMappingTable) {
		return "", errors.New("The policy type of " + policy + " not supported")
	}

	return PolicyTypeMappingTable[policy], nil
}

type StorageTag struct {
	syncTag  map[string]string
	asyncTag map[string]string
}

func NewStorageTag(in map[string]interface{}, flag int) *StorageTag {
	var st = &StorageTag{
		syncTag:  map[string]string{},
		asyncTag: map[string]string{},
	}

	tags := MapValuetoString(in)
	// Screen storage tags through life circle flag
	for key := range tags {
		if flag != PolicyLifecircleTable[key] {
			delete(tags, key)
		}
	}
	// Devide all tags into sync and async part
	for key := range tags {
		pType, err := FindPolicyType(key)
		if err != nil {
			log.Println("[Error] When parse storage tag:", err)
		}
		switch pType {
		case "feature":
			st.syncTag[key] = tags[key]
		case "operation":
			st.asyncTag[key] = tags[key]
		}
	}
	return st
}

func (st *StorageTag) GetSyncTag() map[string]string {
	return st.syncTag
}

func (st *StorageTag) GetAsyncTag() map[string]string {
	return st.asyncTag
}

func MapValuetoString(in map[string]interface{}) map[string]string {
	var out = map[string]string{}

	for k, v := range in {
		switch v.(type) {
		case int:
			out[k] = strconv.Itoa(v.(int))
		case bool:
			out[k] = strconv.FormatBool(v.(bool))
		}
	}
	return out
}
