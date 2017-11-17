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
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/model"
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
				"size": int64(66),
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
				"size": 101,
			},
			pools:    fakePools,
			expected: nil,
		},
	}
	filter := &CapacityFilter{}
	for _, testCase := range testCases {
		result, _ := filter.Handle(testCase.request, testCase.pools)
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
	filter := &AZFilter{}
	for _, testCase := range testCases {
		result, _ := filter.Handle(testCase.request, testCase.pools)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestThinFilter(t *testing.T) {
	fakePools := []*model.StoragePoolSpec{
		&model.StoragePoolSpec{
			Parameters: map[string]interface{}{
				"thin": true,
			},
		},
		&model.StoragePoolSpec{
			Parameters: map[string]interface{}{
				"thin": true,
			},
		},
		&model.StoragePoolSpec{
			Parameters: map[string]interface{}{
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
				"thin": true,
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					Parameters: map[string]interface{}{
						"thin": true,
					},
				},
				&model.StoragePoolSpec{
					Parameters: map[string]interface{}{
						"thin": true,
					},
				},
			},
		},
		{
			request: map[string]interface{}{
				"thin": false,
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					Parameters: map[string]interface{}{
						"thin": false,
					},
				},
			},
		},
	}
	filter := &ThinFilter{}
	for _, testCase := range testCases {
		result, _ := filter.Handle(testCase.request, testCase.pools)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestDedupeFilter(t *testing.T) {
	fakePools := []*model.StoragePoolSpec{
		&model.StoragePoolSpec{
			Parameters: map[string]interface{}{
				"dedupe": true,
			},
		},
		&model.StoragePoolSpec{
			Parameters: map[string]interface{}{
				"dedupe": true,
			},
		},
		&model.StoragePoolSpec{
			Parameters: map[string]interface{}{
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
				"dedupe": true,
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					Parameters: map[string]interface{}{
						"dedupe": true,
					},
				},
				&model.StoragePoolSpec{
					Parameters: map[string]interface{}{
						"dedupe": true,
					},
				},
			},
		},
		{
			request: map[string]interface{}{
				"dedupe": false,
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					Parameters: map[string]interface{}{
						"dedupe": false,
					},
				},
			},
		},
	}
	filter := &DedupeFilter{}
	for _, testCase := range testCases {
		result, _ := filter.Handle(testCase.request, testCase.pools)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestCompressFilter(t *testing.T) {
	fakePools := []*model.StoragePoolSpec{
		&model.StoragePoolSpec{
			Parameters: map[string]interface{}{
				"compress": true,
			},
		},
		&model.StoragePoolSpec{
			Parameters: map[string]interface{}{
				"compress": true,
			},
		},
		&model.StoragePoolSpec{
			Parameters: map[string]interface{}{
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
				"compress": true,
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					Parameters: map[string]interface{}{
						"compress": true,
					},
				},
				&model.StoragePoolSpec{
					Parameters: map[string]interface{}{
						"compress": true,
					},
				},
			},
		},
		{
			request: map[string]interface{}{
				"compress": false,
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					Parameters: map[string]interface{}{
						"compress": false,
					},
				},
			},
		},
	}
	filter := &CompressFilter{}
	for _, testCase := range testCases {
		result, _ := filter.Handle(testCase.request, testCase.pools)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestDiskTypeFilter(t *testing.T) {
	fakePools := []*model.StoragePoolSpec{
		&model.StoragePoolSpec{
			Parameters: map[string]interface{}{
				"diskType": "SSD",
			},
		},
		&model.StoragePoolSpec{
			Parameters: map[string]interface{}{
				"diskType": "SAS",
			},
		},
		&model.StoragePoolSpec{
			Parameters: map[string]interface{}{
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
				"diskType": "SSD",
			},
			pools: fakePools,
			expected: []*model.StoragePoolSpec{
				&model.StoragePoolSpec{
					Parameters: map[string]interface{}{
						"diskType": "SSD",
					},
				},
			},
		},
		{
			request: map[string]interface{}{
				"diskType": "NVMe SSD",
			},
			pools:    fakePools,
			expected: nil,
		},
	}
	filter := &DiskTypeFilter{}
	for _, testCase := range testCases {
		result, _ := filter.Handle(testCase.request, testCase.pools)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}
