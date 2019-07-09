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
	"github.com/opensds/opensds/pkg/utils"
	. "github.com/opensds/opensds/testutils/collection"
)

type FilterCaseSpec struct {
	request  map[string]interface{}
	expected []*model.StoragePoolSpec
}

var (
	FakePools = []*model.StoragePoolSpec{
		&SamplePools[0],
		&SamplePools[1],
	}
)

func TestSimplifyPoolCapabilityMap(t *testing.T) {
	input := make(map[string]interface{})
	input["name"] = "sample-pool-02"
	input["extras.advanced"] = map[string]interface{}{
		"diskType": "SAS",
	}
	simpleMap, unSimpleMap := simplifyPoolCapabilityMap(input)

	if 2 != len(simpleMap) {
		t.Errorf("Expected %v, get %v", 2, len(simpleMap))
	}

	if 0 != len(unSimpleMap) {
		t.Errorf("Expected %v, get %v", 0, len(simpleMap))
	}

	name, ok := simpleMap["name"].(string)
	if (!ok) || (input["name"] != name) {
		t.Errorf("Expected %v/%v, get %v/%v", true, input["name"], ok, name)
	}

	diskType, ok := simpleMap["extras.advanced.diskType"].(string)
	if (!ok) || ("SAS" != diskType) {
		t.Errorf("Expected %v/%v, get %v/%v", true, "SAS", ok, diskType)
	}
}

func TestGetPoolCapabilityMap(t *testing.T) {
	Pool := SamplePools[0]
	result, err := GetPoolCapabilityMap(&Pool)
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	id, ok := result["id"].(string)
	if (!ok) || (Pool.Id != id) {
		t.Errorf("Expected %v/%v, get %v/%v", true, Pool.Id, ok, id)
	}

	FreeCapacity, ok := result["freeCapacity"].(float64)
	if (!ok) || (!utils.IsFloatEqual(FreeCapacity, float64(Pool.FreeCapacity))) {
		t.Errorf("Expected %v/%v, get %v/%v", true, float64(Pool.FreeCapacity), ok, FreeCapacity)
	}

	IsSpaceEfficient, ok := result["extras.dataStorage.isSpaceEfficient"].(bool)
	if (!ok) || (Pool.Extras.DataStorage.IsSpaceEfficient != IsSpaceEfficient) {
		t.Errorf("Expected %v/%v, get %v/%v", true, Pool.Extras.DataStorage.IsSpaceEfficient, ok, IsSpaceEfficient)
	}

	latency, ok := result["extras.advanced.latency"].(string)
	if (!ok) || (Pool.Extras.Advanced["latency"] != latency) {
		t.Errorf("Expected %v/%v, get %v/%v", true, Pool.Extras.Advanced["latency"], ok, latency)
	}
}

func TestIsAvailablePool(t *testing.T) {
	filterReq := make(map[string]interface{})
	filterReq["totalCapacity"] = "<= 100"
	isAvailable, err := IsAvailablePool(filterReq, &SamplePools[0])
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if true != isAvailable {
		t.Errorf("Expected %v, get %v", true, isAvailable)
	}

	filterReq["totalCapacity"] = "!= 200"
	isAvailable, err = IsAvailablePool(filterReq, &SamplePools[1])
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if false != isAvailable {
		t.Errorf("Expected %v, get %v", false, isAvailable)
	}

	delete(filterReq, "totalCapacity")
	filterReq[":totalCapacity"] = "!= 200"
	isAvailable, err = IsAvailablePool(filterReq, &SamplePools[1])
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if true != isAvailable {
		t.Errorf("Expected %v, get %v", false, isAvailable)
	}

}

func TestMatch(t *testing.T) {

	isMatch, err := match("availabilityZone", "default", "<in> defau")
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if true != isMatch {
		t.Errorf("Expected %v, get %v", true, isMatch)
	}

	isMatch, err = match("availabilityZone", "default", "<in> defau1")
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if false != isMatch {
		t.Errorf("Expected %v, get %v", false, isMatch)
	}
}

func TestInOperator(t *testing.T) {

	result, err := InOperator("availabilityZone", "fau", "default")
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if true != result {
		t.Errorf("Expected %v, get %v", true, result)
	}

	result, err = InOperator("availabilityZone", "default1", "default")
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if false != result {
		t.Errorf("Expected %v, get %v", false, result)
	}
}

func TestCompareOperator(t *testing.T) {

	result, err := CompareOperator("s<", "dockId", "123", "122")
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if true != result {
		t.Errorf("Expected %v, get %v", true, result)
	}

	result, err = CompareOperator("s<", "dockId", "123", "124")
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if false != result {
		t.Errorf("Expected %v, get %v", false, result)
	}
}

func TestStringCompare(t *testing.T) {

	result, err := StringCompare("s<", "dockId", "123", "122")
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if false != result {
		t.Errorf("Expected %v, get %v", false, result)
	}

	result, err = StringCompare("s<", "dockId", "123", "124")
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if true != result {
		t.Errorf("Expected %v, get %v", true, result)
	}
}

func TestParseBoolAndCompare(t *testing.T) {

	result, err := ParseBoolAndCompare("isSpaceEfficient", true, "T")
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if true != result {
		t.Errorf("Expected %v, get %v", true, result)
	}

	result, err = ParseBoolAndCompare("isSpaceEfficient", true, "0")
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if false != result {
		t.Errorf("Expected %v, get %v", false, result)
	}
}

func TestParseFloat64AndCompare(t *testing.T) {

	result, err := ParseFloat64AndCompare("==", "freeCapacity", 60.0, "60")
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if true != result {
		t.Errorf("Expected %v, get %v", true, result)
	}

	result, err = ParseFloat64AndCompare("!=", "freeCapacity", 61.0, "60")
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if true != result {
		t.Errorf("Expected %v, get %v", true, result)
	}
}

func TestOrOperator(t *testing.T) {
	words := []string{"<or>", "10", "<or>", "20"}
	result, err := OrOperator("freeCapacity", words, 20.0)
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if true != result {
		t.Errorf("Expected %v, get %v", true, result)
	}

	result, err = OrOperator("freeCapacity", words, 30.0)
	if nil != err {
		t.Errorf("Expected %v, get %v", nil, err)
	}

	if false != result {
		t.Errorf("Expected %v, get %v", false, result)
	}
}

func TestIdFilter(t *testing.T) {
	testCases := []FilterCaseSpec{
		{
			request: map[string]interface{}{
				"id": "084bf71e-a102-11e7-88a8-e31fe6d52248",
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[0],
			},
		},
		{
			request: map[string]interface{}{
				"id": "s== a594b8ac-a103-11e7-985f-d723bcf01b5f",
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[1],
			},
		},
	}

	for _, testCase := range testCases {
		result, _ := SelectSupportedPools(len(FakePools), testCase.request,
			FakePools)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestFreeCapacityFilter(t *testing.T) {
	testCases := []FilterCaseSpec{
		{
			request: map[string]interface{}{
				"freeCapacity": ">= 170",
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[1],
			},
		},
		{
			request: map[string]interface{}{
				"freeCapacity": "== 90",
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[0],
			},
		},
	}

	for _, testCase := range testCases {
		result, _ := SelectSupportedPools(len(FakePools), testCase.request,
			FakePools)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestIsSpaceEfficientFilter(t *testing.T) {
	testCases := []FilterCaseSpec{
		{
			request: map[string]interface{}{
				"extras.dataStorage.isSpaceEfficient": "<is> true",
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[0],
				&SamplePools[1],
			},
		},
		{
			request: map[string]interface{}{
				"extras.dataStorage.isSpaceEfficient": "<is> false",
			},
			expected: nil,
		},
	}

	for _, testCase := range testCases {
		result, _ := SelectSupportedPools(len(FakePools), testCase.request,
			FakePools)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestAdvancedFilter(t *testing.T) {
	testCases := []FilterCaseSpec{
		{
			request: map[string]interface{}{
				"extras.advanced.diskType": "SAS",
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[1],
			},
		},
		{
			request: map[string]interface{}{
				"extras.advanced.diskType": "s>= SSD",
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[0],
			},
		},
	}

	for _, testCase := range testCases {
		result, _ := SelectSupportedPools(len(FakePools), testCase.request,
			FakePools)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

func TestIsMultiAttachFilter(t *testing.T) {
	testCases := []FilterCaseSpec{
		{
			request: map[string]interface{}{
				"multiAttach": "<is> true",
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[0],
			},
		},
		{
			request: map[string]interface{}{
				"multiAttach": "<is> false",
			},
			expected: []*model.StoragePoolSpec{
				&SamplePools[1],
			},
		},
	}

	for _, testCase := range testCases {
		result, _ := SelectSupportedPools(len(FakePools), testCase.request,
			FakePools)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected[0], result[0])
		}
	}
}
