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

package policy

import (
	"reflect"
	"testing"
)

func TestIsStorageTagSupported(t *testing.T) {
	var tags = map[string]string{
		"intervalSnapshot":     "operation",
		"deleteSnapshotPolicy": "operation",
	}

	if !IsStorageTagSupported(tags) {
		t.Errorf("tags %v are not supported by %v\n", tags, PolicyTypeMappingTable)
	}
}

func TestFindPolicyType(t *testing.T) {
	var policys = []string{"thinProvision", "intervalSnapshot"}
	var expectedTypes = []string{"feature", "operation"}

	for i, policy := range policys {
		pType, err := FindPolicyType(policy)
		if err != nil {
			t.Errorf("Failed to find the type of policy %v\n", policy)
		}

		if !reflect.DeepEqual(pType, expectedTypes[i]) {
			t.Errorf("Expected %v, got %v\n", expectedTypes[i], pType)
		}
	}
}

func TestNewStorageTag(t *testing.T) {
	var tags = map[string]interface{}{
		"thinProvision":        true,
		"highAvailability":     false,
		"intervalSnapshot":     "1d",
		"deleteSnapshotPolicy": true,
	}
	var expectedSt = &StorageTag{
		syncTag: map[string]interface{}{
			"thinProvision":    true,
			"highAvailability": false,
		},
		asyncTag: map[string]string{
			"intervalSnapshot":     "1d",
			"deleteSnapshotPolicy": "true",
		},
	}

	st := NewStorageTag(tags, 1)
	if reflect.DeepEqual(st, expectedSt) {
		t.Errorf("Expected %v, got %v\n", expectedSt, st)
	}
}
