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

package utils

import (
	"os"
	"reflect"
	"testing"
)

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
		map[string]interface{}{
			"Lee": 100,
			"fat": false,
		},
		map[string]interface{}{
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
		map[string]string{
			"Lee": "fat",
		},
		map[string]string{
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
