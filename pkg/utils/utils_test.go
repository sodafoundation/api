// Copyright 2019 The OpenSDS Authors.
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

package utils

import (
	"os"
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/constants"
)

func TestRvRepElement(t *testing.T) {
	var strs = []string{"default", "default"}
	str := RvRepElement(strs)
	res := str[0]
	var expect = "default"
	if len(str) != 1 || res != expect {
		t.Errorf("%v remove redundant elements fail,expect:%v,result:%v\n", str, expect, res)
	}
}

func TestContains(t *testing.T) {
	var permissions = []string{"Read", "Write", "Execute"}
	var testkey = "Read"
	if !Contains(permissions, testkey) {
		t.Errorf("%v is not contains in %v\n", testkey, permissions)
	}
}

func TestContained(t *testing.T) {
	var targets = []interface{}{
		[]interface{}{"key01", 123, true},
		map[interface{}]string{
			"key01": "value01",
			true:    "value02",
			123:     "value03",
		},
	}
	var objs = []interface{}{"key01", 123, true}

	for _, obj := range objs {
		for _, target := range targets {
			if !Contained(obj, target) {
				t.Errorf("%v is not contained in %v\n", obj, target)
			}
		}
	}
}

func TestMergeGeneralMaps(t *testing.T) {
	input := []map[string]interface{}{
		{
			"Lee": 100,
			"fat": false,
		},
		{
			"Ming": 50,
			"tall": true,
		},
	}
	var expected = map[string]interface{}{
		"Lee":  100,
		"fat":  false,
		"Ming": 50,
		"tall": true,
	}

	result := MergeGeneralMaps(input...)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, get %v\n", expected, result)
	}
}

func TestMergeStringMaps(t *testing.T) {
	input := []map[string]string{
		{
			"Lee": "fat",
		},
		{
			"Ming": "thin",
		},
	}
	var expected = map[string]string{
		"Lee":  "fat",
		"Ming": "thin",
	}

	result := MergeStringMaps(input...)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, get %v\n", expected, result)
	}
}

func TestPathExists(t *testing.T) {
	testDir := "./testDir"
	isExist, _ := PathExists(testDir)
	if isExist {
		t.Errorf("Expected false, get %v\n", isExist)
	}
	os.MkdirAll(testDir, 0755)
	isExist, _ = PathExists(testDir)
	if !isExist {
		t.Errorf("Expected true, get %v\n", isExist)
	}
	os.RemoveAll(testDir)
}

func TestStructToMap(t *testing.T) {
	PoolA := model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4486139-78d5-462d-a7b9-fdaf6c797e11",
			CreatedAt: "2017-10-24T15:04:05",
		},
		FreeCapacity:     int64(50),
		AvailabilityZone: "az1",
		Extras: model.StoragePoolExtraSpec{
			DataStorage:    model.DataStorageLoS{ProvisioningPolicy: "Thin", Compression: false, Deduplication: false},
			IOConnectivity: model.IOConnectivityLoS{AccessProtocol: "iscsi", MaxIOPS: 700000, MaxBWS: 600, MinIOPS: 100000, MinBWS: 100, Latency: 100},
			DataProtection: model.DataProtectionLoS{IsIsolated: true, ReplicaType: "Mirror"},
			Advanced:       map[string]interface{}{"diskType": "SSD", "latency": 5},
		},
	}

	poolMap, err := StructToMap(PoolA)

	if nil != err {
		t.Errorf(err.Error())
	}

	value, ok := poolMap["freeCapacity"]
	if !ok {
		t.Errorf("Expected %v, get %v", true, ok)
	}

	freeCapacity, ok := value.(float64)
	if !ok || !IsFloatEqual(freeCapacity, 50) {
		t.Errorf("Expected %v/%v, get %v/%v", true, 50.0, ok, freeCapacity)
	}

	value, ok = poolMap["extras"]
	if !ok {
		t.Errorf("Expected %v, get %v", true, ok)
	}
	_, ok = value.(map[string]interface{})
	if !ok {
		t.Errorf("Expected %v, get %v", true, ok)
	}
}

func TestIsFloatEqual(t *testing.T) {
	isEqual := IsFloatEqual(0.0, 0.00)
	if true != isEqual {
		t.Errorf("Expected %v, get %v", true, isEqual)
	}

	isEqual = IsFloatEqual(1.00, 1)
	if true != isEqual {
		t.Errorf("Expected %v, get %v", true, isEqual)
	}

	isEqual = IsFloatEqual(-1.00, -1)
	if true != isEqual {
		t.Errorf("Expected %v, get %v", true, isEqual)
	}

	isEqual = IsFloatEqual(2.00, 1)
	if false != isEqual {
		t.Errorf("Expected %v, get %v", false, isEqual)
	}

	isEqual = IsFloatEqual(-1.00, -2)
	if false != isEqual {
		t.Errorf("Expected %v, get %v", false, isEqual)
	}
}

func TestIsEqual(t *testing.T) {
	isEqual, err := IsEqual("", true, true)
	if true != isEqual {
		t.Errorf("Expected %v, get %v", true, isEqual)
	}

	isEqual, err = IsEqual("", false, true)
	if false != isEqual {
		t.Errorf("Expected %v, get %v", false, isEqual)
	}

	isEqual, err = IsEqual("", -1.00, -1.000)
	if true != isEqual {
		t.Errorf("Expected %v, get %v", true, isEqual)
	}

	isEqual, err = IsEqual("", 2.00, 1)
	if false != isEqual {
		t.Errorf("Expected %v, get %v", false, isEqual)
	}

	isEqual, err = IsEqual("", "abc", "abc")
	if true != isEqual {
		t.Errorf("Expected %v, get %v", true, isEqual)
	}

	isEqual, err = IsEqual("", "abc", "ABC")
	if false != isEqual {
		t.Errorf("Expected %v, get %v", false, isEqual)
	}

	isEqual, err = IsEqual("keyA", "abc", true)
	expectedErr := "the type of keyA must be string"
	if expectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", expectedErr, err.Error())
	}
}

type animal struct {
	Name string
	Kind string
	Age  int
}

func TestFilter(t *testing.T) {

	testCases := []struct {
		animals  []*animal
		params   map[string][]string
		expected []*animal
	}{
		{
			animals: []*animal{
				&animal{
					Name: "Tony",
					Kind: "dog",
					Age:  5,
				},
				&animal{
					Name: "Jack",
					Kind: "horse",
					Age:  10,
				},
				&animal{
					Name: "Huhu",
					Kind: "tiger",
					Age:  3,
				},
			},
			params: map[string][]string{"kind": []string{"tiger"}},
			expected: []*animal{
				&animal{
					Name: "Huhu",
					Kind: "tiger",
					Age:  3,
				},
			},
		},
		{
			animals: []*animal{
				&animal{
					Name: "Tony",
					Kind: "dog",
					Age:  5,
				},
				&animal{
					Name: "Jack",
					Kind: "horse",
					Age:  10,
				},
				&animal{
					Name: "Huhu",
					Kind: "tiger",
					Age:  3,
				},
			},
			params: map[string][]string{"name": []string{"Tony", "Jack"}},
			expected: []*animal{
				&animal{
					Name: "Tony",
					Kind: "dog",
					Age:  5,
				},
				&animal{
					Name: "Jack",
					Kind: "horse",
					Age:  10,
				},
			},
		},
		{
			animals: []*animal{
				&animal{
					Name: "Tony",
					Kind: "dog",
					Age:  5,
				},
			},
			params:   map[string][]string{"name": []string{"Kio"}},
			expected: nil,
		},
	}
	for _, c := range testCases {
		results := Filter(c.animals, c.params).([]interface{})
		if len(results) != len(c.expected) {
			t.Errorf("Expected %v, get %v", len(c.expected), len(results))
		}
		for i, actual := range results {
			if !reflect.DeepEqual(actual, c.expected[i]) {
				t.Errorf("Expected %v, get %v", c.expected[i], actual)
			}
		}

	}
}

func TestSort(t *testing.T) {

	testCases := []struct {
		animals  []*animal
		sortKey  string
		sortDir  string
		expected []*animal
	}{
		{
			animals: []*animal{
				&animal{
					Name: "Tony",
					Kind: "dog",
					Age:  5,
				},
				&animal{
					Name: "Huhu",
					Kind: "tiger",
					Age:  3,
				},
				&animal{
					Name: "Jack",
					Kind: "horse",
					Age:  10,
				},
			},
			sortKey: "name",
			sortDir: "asc",
			expected: []*animal{
				&animal{
					Name: "Huhu",
					Kind: "tiger",
					Age:  3,
				},
				&animal{
					Name: "Jack",
					Kind: "horse",
					Age:  10,
				},
				&animal{
					Name: "Tony",
					Kind: "dog",
					Age:  5,
				},
			},
		},
		{
			animals: []*animal{
				&animal{
					Name: "Tony",
					Kind: "dog",
					Age:  5,
				},
				&animal{
					Name: "Jack",
					Kind: "horse",
					Age:  10,
				},
				&animal{
					Name: "Huhu",
					Kind: "tiger",
					Age:  3,
				},
			},
			sortKey: "age",
			sortDir: constants.DefaultSortDir,
			expected: []*animal{
				&animal{
					Name: "Jack",
					Kind: "horse",
					Age:  10,
				},
				&animal{
					Name: "Tony",
					Kind: "dog",
					Age:  5,
				},
				&animal{
					Name: "Huhu",
					Kind: "tiger",
					Age:  3,
				},
			},
		},
	}
	for _, c := range testCases {
		results := Sort(c.animals, c.sortKey, c.sortDir).([]*animal)
		if len(results) != len(c.expected) {
			t.Errorf("Expected %v, get %v", len(c.expected), len(results))
		}
		for i, actual := range results {
			if !reflect.DeepEqual(actual, c.expected[i]) {
				t.Errorf("Expected %v, get %v", c.expected[i], actual)
			}
		}

	}
}

func TestSlice(t *testing.T) {

	testCases := []struct {
		animals  []*animal
		offset   int
		limit    int
		expected []*animal
	}{
		{
			animals: []*animal{
				&animal{
					Name: "Tony",
					Kind: "dog",
					Age:  5,
				},
				&animal{
					Name: "Huhu",
					Kind: "tiger",
					Age:  3,
				},
				&animal{
					Name: "Jack",
					Kind: "horse",
					Age:  10,
				},
			},
			offset: 1,
			limit:  2,
			expected: []*animal{
				&animal{
					Name: "Huhu",
					Kind: "tiger",
					Age:  3,
				},
				&animal{
					Name: "Jack",
					Kind: "horse",
					Age:  10,
				},
			},
		},
		{
			animals: []*animal{
				&animal{
					Name: "Tony",
					Kind: "dog",
					Age:  5,
				},
				&animal{
					Name: "Huhu",
					Kind: "tiger",
					Age:  3,
				},
				&animal{
					Name: "Jack",
					Kind: "horse",
					Age:  10,
				},
			},
			offset:   3,
			limit:    3,
			expected: nil,
		},
	}
	for _, c := range testCases {
		results := Slice(c.animals, c.offset, c.limit).([]interface{})
		if len(results) != len(c.expected) {
			t.Errorf("Expected %v, get %v", len(c.expected), len(results))
		}
		for i, actual := range results {
			if !reflect.DeepEqual(actual, c.expected[i]) {
				t.Errorf("Expected %v, get %v", c.expected[i], actual)
			}
		}

	}
}

func TestContainsIgnoreCase(t *testing.T) {
	var permissions = []string{"Read", "Write", "Execute"}
	var testkey1 = "READ"
	var testkey2 = "Read"
	var testkey3 = "read"
	var testkey4 = "Raed"
	contains := ContainsIgnoreCase(permissions, testkey1)
	if contains != true {
		t.Errorf("%v should contains %v", permissions, testkey1)
	}
	contains = ContainsIgnoreCase(permissions, testkey2)
	if contains != true {
		t.Errorf("%v should contains %v", permissions, testkey2)
	}
	contains = ContainsIgnoreCase(permissions, testkey3)
	if contains != true {
		t.Errorf("%v should contains %v", permissions, testkey3)
	}
	contains = ContainsIgnoreCase(permissions, testkey4)
	if contains != false {
		t.Errorf("%v shouldn't contains %v", permissions, testkey4)
	}
}
