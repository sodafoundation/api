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

package policyengine

import (
	"encoding/json"
	"errors"
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
	"github.com/opensds/opensds/pkg/controller/metadata/profile"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
)

const (
	POLICY_TYPE_MAPPING_TABLE = `{
		"iops": "feature",
		"latency": "feature",
		"highAvailability": "feature",
		"intervalSnapshot": "operation",
		"intervalBackup": "operation",
		"deleteSnapshotPolicy": "operation"
	}`
	POLICY_LIFECIRCLE_TABLE = `{
		"iops": 1,
		"latency": 1,
		"highAvailability": 1,
		"intervalSnapshot": 1,
		"intervalBackup": 1,
		"deleteSnapshotPolicy": 4
	}`
)

var PolicyTypeMappingTable map[string]string
var PolicyLifecircleTable map[string]int

func Init() {
	json.Unmarshal([]byte(POLICY_TYPE_MAPPING_TABLE), &PolicyTypeMappingTable)
	json.Unmarshal([]byte(POLICY_LIFECIRCLE_TABLE), &PolicyLifecircleTable)
}

func IsProfileSupported(desiredProfile *api.StorageProfile) bool {
	profiles, err := profile.ListProfiles()
	if err != nil {
		log.Println("[Error] When list profiles:", err)
		return false
	}

	var isSupported bool
	for _, profile := range profiles.Profiles {
		// If profile name provided, find if the same one exists in profile tables
		if profile.Name == desiredProfile.Name && profile.Name != "" {
			desiredProfile.BackendDriver = profile.BackendDriver
			desiredProfile.StorageTags = profile.StorageTags
			return true
		}
		// If backend type is not same, move to the next profile
		if profile.BackendDriver != desiredProfile.BackendDriver {
			continue
		}
		// Find if the desired storage tags are contained in any profile
		isSupported = true
		for tag := range desiredProfile.StorageTags {
			if !Contained(tag, profile.StorageTags) {
				isSupported = false
			}
		}
	}
	return isSupported
}

func FindPolicyType(policy string) (string, error) {
	for key := range PolicyTypeMappingTable {
		if key == policy {
			return PolicyTypeMappingTable[key], nil
		}
	}

	return "", errors.New("The policy type of " + policy + " not supported")
}

type StorageTag struct {
	syncTag  map[string]string
	asyncTag map[string]string
}

func NewStorageTag(tags map[string]string, flag int) *StorageTag {
	var st = &StorageTag{
		syncTag:  map[string]string{},
		asyncTag: map[string]string{},
	}
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

func ExecuteSyncPolicy(st *StorageTag, vr *pb.VolumeRequest) error {
	swf, err := RegisterSynchronizedWorkflow(vr, st.syncTag)
	if err != nil {
		return err
	}

	if err = ExecuteSynchronizedWorkflow(swf); err != nil {
		return err
	}
	return nil
}

func ExecuteAsyncPolicy(vr *pb.VolumeRequest, st *StorageTag, in string, errChan chan error) {
	awf, err := RegisterAsynchronizedWorkflow(vr, st.asyncTag, in)
	if err != nil {
		errChan <- err
	}

	defer close(errChan)

	if err = ExecuteAsynchronizedWorkflow(awf); err != nil {
		errChan <- err
	}
	errChan <- nil
}
