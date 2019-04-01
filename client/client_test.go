// Copyright (c) 2019 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package client

import (
	"reflect"
	"testing"
)

func TestProcessListParam(t *testing.T) {
	// Test case 1: The args should only support one parameter.
	argA := map[string]string{"limit": "3"}
	argB := map[string]string{"offset": "4"}
	args := []interface{}{argA, argB}
	_, err := processListParam(args)
	expectedError := "args should only support one parameter"
	if err == nil {
		t.Errorf("expected Non-%v, got %v\n", nil, err)
	} else {
		if expectedError != err.Error() {
			t.Errorf("expected Non-%v, got %v\n", expectedError, err.Error())
		}
	}

	// Test case 2: The args type should only be map[string]string.
	args = []interface{}{"limit=3&offset=4"}
	_, err = processListParam(args)
	expectedError = "args element type should be map[string]string"
	if err == nil {
		t.Errorf("expected Non-%v, got %v\n", nil, err)
	} else {
		if expectedError != err.Error() {
			t.Errorf("expected Non-%v, got %v\n", expectedError, err.Error())
		}
	}

	// Test case 3: Test the output format if everything works well.
	arg := map[string]string{"limit": "3", "offset": "4"}
	expectedA, expectedB := "limit=3&offset=4", "offset=4&limit=3"
	params, err := processListParam([]interface{}{arg})
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(params, expectedA) && !reflect.DeepEqual(params, expectedB) {
		t.Errorf("expected %v or %v, got %v\n", expectedA, expectedB, params)
	}
}
