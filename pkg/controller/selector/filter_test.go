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
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
)

func TestGetPoolCapabilityMap(t *testing.T) {
	Pool := model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4486139-78d5-462d-a7b9-fdaf6c797e11",
			CreatedAt: "2017-10-24T15:04:05",
		},
		FreeCapacity:     int64(50),
		AvailabilityZone: "az1",
		Extras: model.StoragePoolExtraSpec{
			DataStorage: model.DataStorageLoS{
				RecoveryTimeObjective: 3,
				ProvisioningPolicy:    "Thin",
				IsSpaceEfficient:      true,
			},
			IOConnectivity: model.IOConnectivityLoS{
				AccessProtocol: "rbd",
				MaxIOPS:        1,
				MaxBWS:         1000,
			},
			DataProtection: model.DataProtectionLos{},
			Advanced: model.ExtraSpec{
				"thin":     true,
				"dedupe":   true,
				"compress": true,
				"diskType": "SSD",
			},
		},
	}

	var mapA map[string]interface{}
	mapA = make(map[string]interface{})
	mapA["key1"] = "value1"
	mapA["key2"] = "value2"
	Pool.Extras.Advanced["mapA"] = mapA

	result, err := GetPoolCapabilityMap(&Pool)
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, result)
	}

	CreatedAt, ok := result["createdAt"].(string)
	if (!ok) || (Pool.CreatedAt != CreatedAt) {
		t.Errorf("Expected %v/%v, get %v/%v", true, Pool.CreatedAt, ok, CreatedAt)
	}

	FreeCapacity, ok := result["freeCapacity"].(float64)
	if (!ok) || (!utils.IsFloatEqual(FreeCapacity, float64(Pool.FreeCapacity))) {
		t.Errorf("Expected %v/%v, get %v/%v", true, float64(Pool.FreeCapacity), ok, FreeCapacity)
	}

	thin, ok := result["extras.advanced.thin"].(bool)
	if (!ok) || (Pool.Extras.Advanced["thin"] != thin) {
		t.Errorf("Expected %v/%v, get %v/%v", true, Pool.Extras.Advanced["thin"], ok, thin)
	}

	value1, ok := result["extras.advanced.mapA.key1"].(string)
	if (!ok) || ("value1" != value1) {
		t.Errorf("Expected %v/%v, get %v/%v", true, "value1", ok, value1)
	}

	RecoveryTimeObjective, ok := result["extras.dataStorage.recoveryTimeObjective"].(float64)
	if (!ok) || (!utils.IsFloatEqual(RecoveryTimeObjective, float64(Pool.Extras.DataStorage.RecoveryTimeObjective))) {
		t.Errorf("Expected %v/%v, get %v/%v", true, float64(Pool.Extras.DataStorage.RecoveryTimeObjective), ok, RecoveryTimeObjective)
	}

	ProvisioningPolicy, ok := result["extras.dataStorage.provisioningPolicy"].(string)
	if (!ok) || (Pool.Extras.DataStorage.ProvisioningPolicy != ProvisioningPolicy) {
		t.Errorf("Expected %v/%v, get %v/%v", true, Pool.Extras.DataStorage.ProvisioningPolicy, ok, ProvisioningPolicy)
	}

	IsSpaceEfficient, ok := result["extras.dataStorage.isSpaceEfficient"].(bool)
	if (!ok) || (Pool.Extras.DataStorage.IsSpaceEfficient != IsSpaceEfficient) {
		t.Errorf("Expected %v/%v, get %v/%v", true, Pool.Extras.DataStorage.IsSpaceEfficient, ok, IsSpaceEfficient)
	}
}

var (
	FakePools = []*model.StoragePoolSpec{
		&model.StoragePoolSpec{},
		&model.StoragePoolSpec{},
		&model.StoragePoolSpec{},
	}

	TestCases = []struct {
		request  map[string]interface{}
		pools    []*model.StoragePoolSpec
		expected []*model.StoragePoolSpec
	}{
		{
			pools: FakePools,
		},
		{
			pools: FakePools,
		},
	}
)

func TestCreatedAtFilter(t *testing.T) {
	FakePools[0].BaseModel = &model.BaseModel{
		CreatedAt: "2017-10-24T15:04:05",
	}
	FakePools[1].BaseModel = &model.BaseModel{
		CreatedAt: "2017-10-24T15:04:06",
	}
	FakePools[2].BaseModel = &model.BaseModel{
		CreatedAt: "2017-10-24T15:04:07",
	}

	TestCases[0].request = map[string]interface{}{
		"createdAt": "s== 2017-10-24T15:04:05",
	}
	TestCases[0].expected = []*model.StoragePoolSpec{
		FakePools[0],
	}

	TestCases[1].request = map[string]interface{}{
		"createdAt": "s>= 2017-10-24T15:04:06",
	}
	TestCases[1].expected = []*model.StoragePoolSpec{
		FakePools[1],
		FakePools[2],
	}

	for _, testCase := range TestCases {
		result, _ := SelectSupportedPools(len(testCase.pools), testCase.request,
			testCase.pools)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestFreeCapacityFilter(t *testing.T) {
	FakePools[0].FreeCapacity = 100
	FakePools[1].FreeCapacity = 50
	FakePools[2].FreeCapacity = 66

	TestCases[0].request = map[string]interface{}{
		"freeCapacity": ">= 66",
	}
	TestCases[0].expected = []*model.StoragePoolSpec{
		FakePools[0],
		FakePools[2],
	}

	TestCases[1].request = map[string]interface{}{
		"freeCapacity": "> 100",
	}
	TestCases[1].expected = nil

	for _, testCase := range TestCases {
		result, _ := SelectSupportedPools(len(testCase.pools), testCase.request,
			testCase.pools)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestAccessProtocolFilter(t *testing.T) {
	FakePools[0].Extras.IOConnectivity = model.IOConnectivityLoS{
		AccessProtocol: "dbr",
	}
	FakePools[1].Extras.IOConnectivity = model.IOConnectivityLoS{
		AccessProtocol: "rbd",
	}
	FakePools[2].Extras.IOConnectivity = model.IOConnectivityLoS{
		AccessProtocol: "brd",
	}

	TestCases[0].request = map[string]interface{}{
		"extras.ioConnectivity.accessProtocol": "rbd",
	}
	TestCases[0].expected = []*model.StoragePoolSpec{
		FakePools[1],
	}

	TestCases[1].request = map[string]interface{}{
		"extras.ioConnectivity.accessProtocol": "s!= rbd",
	}

	TestCases[1].expected = []*model.StoragePoolSpec{
		FakePools[0],
		FakePools[2],
	}

	for _, testCase := range TestCases {
		result, _ := SelectSupportedPools(len(testCase.pools), testCase.request,
			testCase.pools)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestAdvancedFilter(t *testing.T) {
	FakePools[0].Extras.Advanced = model.ExtraSpec{
		"compress": true,
	}
	FakePools[1].Extras.Advanced = model.ExtraSpec{
		"compress": true,
	}
	FakePools[2].Extras.Advanced = model.ExtraSpec{
		"compress": false,
	}

	TestCases[0].request = map[string]interface{}{
		"extras.advanced.compress": true,
	}
	TestCases[0].expected = []*model.StoragePoolSpec{
		FakePools[0],
		FakePools[1],
	}

	TestCases[1].request = map[string]interface{}{
		"extras.advanced.compress": false,
	}
	TestCases[1].expected = []*model.StoragePoolSpec{
		FakePools[2],
	}

	for _, testCase := range TestCases {
		result, _ := SelectSupportedPools(len(testCase.pools), testCase.request,
			testCase.pools)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}
