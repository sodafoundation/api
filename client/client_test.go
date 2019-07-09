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

package client

import (
	"reflect"
	"testing"
)

var assertTestResult = func(t *testing.T, got, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v, got %v\n", expected, got)
	}
}

func TestProcessListParam(t *testing.T) {

	t.Run("The args should only support one parameter", func(t *testing.T) {
		argA := map[string]string{"limit": "3"}
		argB := map[string]string{"offset": "4"}
		args := []interface{}{argA, argB}
		_, err := processListParam(args)
		expectedError := "args should only support one parameter"
		assertTestResult(t, err.Error(), expectedError)
	})

	t.Run("The args type should only be map[string]string", func(t *testing.T) {
		args := []interface{}{"limit=3&offset=4"}
		_, err := processListParam(args)
		expectedError := "args element type should be map[string]string"
		assertTestResult(t, err.Error(), expectedError)
	})

	t.Run("Test the output format if everything works well", func(t *testing.T) {
		arg := map[string]string{"limit": "3", "offset": "4"}
		expectedA, expectedB := "limit=3&offset=4", "offset=4&limit=3"
		params, err := processListParam([]interface{}{arg})
		assertTestResult(t, err, nil)
		if !reflect.DeepEqual(params, expectedA) && !reflect.DeepEqual(params, expectedB) {
			t.Errorf("expected %v or %v, got %v\n", expectedA, expectedB, params)
		}
	})
}
