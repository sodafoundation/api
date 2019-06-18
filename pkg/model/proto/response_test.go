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

package proto

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

func TestGenericResponseResult(t *testing.T) {

	t.Run("Should marshal string type", func(t *testing.T) {
		resMsg := "hello, world"
		expected := &GenericResponse{
			Reply: &GenericResponse_Result_{
				Result: &GenericResponse_Result{
					Message: "hello, world",
				},
			},
		}

		result := GenericResponseResult(resMsg)
		assertTestResult(t, result, expected)
	})

	t.Run("Should marshal []string type", func(t *testing.T) {
		resMsg := []string{"hello, world"}
		expected := &GenericResponse{
			Reply: &GenericResponse_Result_{
				Result: &GenericResponse_Result{
					Message: `["hello, world"]`,
				},
			},
		}

		result := GenericResponseResult(resMsg)
		assertTestResult(t, result, expected)
	})

	t.Run("Should marshal map[string]interface{} type", func(t *testing.T) {
		resMsg := map[string]interface{}{"key": "hello, world"}
		expected := &GenericResponse{
			Reply: &GenericResponse_Result_{
				Result: &GenericResponse_Result{
					Message: `{"key":"hello, world"}`,
				},
			},
		}

		result := GenericResponseResult(resMsg)
		assertTestResult(t, result, expected)
	})
}
