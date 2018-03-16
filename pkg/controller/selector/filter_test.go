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

func TestCapacityFilter(t *testing.T) {
	fakePools := []*model.StoragePoolSpec{
		&model.StoragePoolSpec{
			FreeCapacity: int64(100),
		},
		&model.StoragePoolSpec{
			FreeCapacity: int64(50),
		},
		&model.StoragePoolSpec{
			FreeCapacity: int64(66),
		},
	}
	testCases := []struct {
		request  map[string]interface{}
		pools    []*model.StoragePoolSpec
		expected []*model.StoragePoolSpec
	}{
		{
			request: map[string]interface{}{
				"freeCapacity": ">= 66",
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					FreeCapacity: int64(100),
				},
				&model.StoragePoolSpec{
					FreeCapacity: int64(66),
				},
			},
		},
		{
			request: map[string]interface{}{
				"freeCapacity": ">= 101",
			},
			pools:    fakePools,
			expected: nil,
		},
	}

	for _, testCase := range testCases {
		result, _ := SelectSupportedPools(len(testCase.pools), testCase.request,
			testCase.pools)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestAZFilter(t *testing.T) {
	fakePools := []*model.StoragePoolSpec{
		&model.StoragePoolSpec{
			AvailabilityZone: "az1",
		},
		&model.StoragePoolSpec{
			AvailabilityZone: "az2",
		},
		&model.StoragePoolSpec{
			AvailabilityZone: "az1",
		},
	}
	testCases := []struct {
		request  map[string]interface{}
		pools    []*model.StoragePoolSpec
		expected []*model.StoragePoolSpec
	}{
		{
			request: map[string]interface{}{
				"availabilityZone": "az1",
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					AvailabilityZone: "az1",
				},
				&model.StoragePoolSpec{
					AvailabilityZone: "az1",
				},
			},
		},
		{
			request: map[string]interface{}{
				"availabilityZone": "az3",
			},
			pools:    fakePools,
			expected: nil,
		},
	}

	for _, testCase := range testCases {
		result, _ := SelectSupportedPools(len(testCase.pools), testCase.request,
			testCase.pools)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestThinFilter(t *testing.T) {
	fakePools := []*model.StoragePoolSpec{
		&model.StoragePoolSpec{
			Extras: model.ExtraSpec{
				"thin": true,
			},
		},
		&model.StoragePoolSpec{
			Extras: model.ExtraSpec{
				"thin": true,
			},
		},
		&model.StoragePoolSpec{
			Extras: model.ExtraSpec{
				"thin": false,
			},
		},
	}
	testCases := []struct {
		request  map[string]interface{}
		pools    []*model.StoragePoolSpec
		expected []*model.StoragePoolSpec
	}{
		{
			request: map[string]interface{}{
				"extras.thin": true,
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					Extras: model.ExtraSpec{
						"thin": true,
					},
				},
				&model.StoragePoolSpec{
					Extras: model.ExtraSpec{
						"thin": true,
					},
				},
			},
		},
		{
			request: map[string]interface{}{
				"extras.thin": false,
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					Extras: model.ExtraSpec{
						"thin": false,
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		result, _ := SelectSupportedPools(len(testCase.pools), testCase.request,
			testCase.pools)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestDedupeFilter(t *testing.T) {
	fakePools := []*model.StoragePoolSpec{
		&model.StoragePoolSpec{
			Extras: model.ExtraSpec{
				"dedupe": true,
			},
		},
		&model.StoragePoolSpec{
			Extras: model.ExtraSpec{
				"dedupe": true,
			},
		},
		&model.StoragePoolSpec{
			Extras: model.ExtraSpec{
				"dedupe": false,
			},
		},
	}
	testCases := []struct {
		request  map[string]interface{}
		pools    []*model.StoragePoolSpec
		expected []*model.StoragePoolSpec
	}{
		{
			request: map[string]interface{}{
				"extras.dedupe": true,
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					Extras: model.ExtraSpec{
						"dedupe": true,
					},
				},
				&model.StoragePoolSpec{
					Extras: model.ExtraSpec{
						"dedupe": true,
					},
				},
			},
		},
		{
			request: map[string]interface{}{
				"extras.dedupe": false,
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					Extras: model.ExtraSpec{
						"dedupe": false,
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		result, _ := SelectSupportedPools(len(testCase.pools), testCase.request,
			testCase.pools)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestCompressFilter(t *testing.T) {
	fakePools := []*model.StoragePoolSpec{
		&model.StoragePoolSpec{
			Extras: model.ExtraSpec{
				"compress": true,
			},
		},
		&model.StoragePoolSpec{
			Extras: model.ExtraSpec{
				"compress": true,
			},
		},
		&model.StoragePoolSpec{
			Extras: model.ExtraSpec{
				"compress": false,
			},
		},
	}
	testCases := []struct {
		request  map[string]interface{}
		pools    []*model.StoragePoolSpec
		expected []*model.StoragePoolSpec
	}{
		{
			request: map[string]interface{}{
				"extras.compress": true,
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					Extras: model.ExtraSpec{
						"compress": true,
					},
				},
				&model.StoragePoolSpec{
					Extras: model.ExtraSpec{
						"compress": true,
					},
				},
			},
		},
		{
			request: map[string]interface{}{
				"extras.compress": false,
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					Extras: model.ExtraSpec{
						"compress": false,
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		result, _ := SelectSupportedPools(len(testCase.pools), testCase.request,
			testCase.pools)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestDiskTypeFilter(t *testing.T) {
	fakePools := []*model.StoragePoolSpec{
		&model.StoragePoolSpec{
			Extras: model.ExtraSpec{
				"diskType": "SSD",
			},
		},
		&model.StoragePoolSpec{
			Extras: model.ExtraSpec{
				"diskType": "SAS",
			},
		},
		&model.StoragePoolSpec{
			Extras: model.ExtraSpec{
				"diskType": "SATA",
			},
		},
	}
	testCases := []struct {
		request  map[string]interface{}
		pools    []*model.StoragePoolSpec
		expected []*model.StoragePoolSpec
	}{
		{
			request: map[string]interface{}{
				"extras.diskType": "SSD",
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					Extras: model.ExtraSpec{
						"diskType": "SSD",
					},
				},
			},
		},
		{
			request: map[string]interface{}{
				"extras.diskType": "NVMe SSD",
			},
			pools:    fakePools,
			expected: nil,
		},
	}

	for _, testCase := range testCases {
		result, _ := SelectSupportedPools(len(testCase.pools), testCase.request,
			testCase.pools)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestGetPoolCapabilityMap(t *testing.T) {
	Pool := model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4486139-78d5-462d-a7b9-fdaf6c797e11",
			CreatedAt: "2017-10-24T15:04:05",
		},
		FreeCapacity:     int64(50),
		AvailabilityZone: "az1",
		Extras: model.ExtraSpec{
			"thin":     true,
			"dedupe":   true,
			"compress": true,
			"diskType": "SSD",
		},
	}

	var mapA map[string]interface{}
	mapA = make(map[string]interface{})
	mapA["key1"] = "value1"
	mapA["key2"] = "value2"
	Pool.Extras["mapA"] = mapA

	result, err := GetPoolCapabilityMap(&Pool)
	if nil != err {
		t.Errorf("Expected nil, get %v", result)
	}

	CreatedAt, ok := result["createdAt"].(string)
	if (!ok) || (Pool.CreatedAt != CreatedAt) {
		t.Errorf("Expected true/%v, get %v/%v ", Pool.CreatedAt, ok, CreatedAt)
	}

	FreeCapacity, ok := result["freeCapacity"].(float64)
	if (!ok) || (!utils.IsFloatEqual(FreeCapacity, 50)) {
		t.Errorf("Expected true/%v, get %v/%v ", Pool.FreeCapacity, ok, FreeCapacity)
	}

	thin, ok := result["extras.thin"].(bool)
	if (!ok) || (Pool.Extras["thin"] != thin) {
		t.Errorf("Expected true/%v, get %v/%v ", Pool.Extras["thin"], ok, thin)
	}

	value1, ok := result["extras.mapA.key1"].(string)
	if (!ok) || ("value1" != value1) {
		t.Errorf("Expected true/'value1', get %v/%v ", ok, value1)
	}
}
