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

package selector

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
)

func TestSelectSupportedPool(t *testing.T) {
	mockClient := new(dbtest.MockClient)
	mockClient.On("ListPools").Return(fakePools, nil)
	db.C = mockClient

	testCases := []struct {
		request  map[string]interface{}
		expected *model.StoragePoolSpec
	}{
		{
			request: map[string]interface{}{
				"size":             int64(5001),
				"availabilityZone": "az1",
				"thin":             true,
			},
			expected: fakePools[1],
		},
		{
			request: map[string]interface{}{
				"size":             int64(400),
				"availabilityZone": "default",
				"diskType":         "SSD",
			},
			expected: nil,
		},
	}

	s := NewSelector()
	for _, testCase := range testCases {
		result, _ := s.SelectSupportedPool(testCase.request)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

var (
	fakePools = []*model.StoragePoolSpec{
		{
			BaseModel: &model.BaseModel{
				Id:        "f4486139-78d5-462d-a7b9-fdaf6c797e1b",
				CreatedAt: "2017-10-24T15:04:05",
			},
			Name:             "fakePool",
			Description:      "fake pool for testing",
			Status:           "available",
			AvailabilityZone: "az1",
			TotalCapacity:    99999,
			FreeCapacity:     5000,
			DockId:           "ccac4f33-e603-425a-8813-371bbe10566e",
			Extras: model.ExtraSpec{
				"thin":     true,
				"dedupe":   false,
				"compress": false,
				"diskType": "SSD",
			},
		},
		{
			BaseModel: &model.BaseModel{
				Id:        "f4486139-78d5-462d-a7b9-fdaf6c797e1b",
				CreatedAt: "2017-10-24T15:04:05",
			},
			Name:             "fakePool",
			Description:      "fake pool for testing",
			Status:           "available",
			AvailabilityZone: "az1",
			TotalCapacity:    99999,
			FreeCapacity:     6999,
			DockId:           "ccac4f33-e603-425a-8813-371bbe10566e",
			Extras: model.ExtraSpec{
				"thin":     true,
				"dedupe":   true,
				"compress": true,
				"diskType": "SATA",
			},
		},
	}
)
