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
package client

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/model"
)

var fp = &PoolMgr{
	Receiver: NewFakePoolReceiver(),
}

func TestGetPool(t *testing.T) {
	var polID = "084bf71e-a102-11e7-88a8-e31fe6d52248"
	expected := &model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
		},
		Name:          "sample-pool-01",
		Description:   "This is the first sample storage pool for testing",
		StorageType:   "block",
		TotalCapacity: int64(100),
		FreeCapacity:  int64(90),
		DockId:        "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		Extras: model.StoragePoolExtraSpec{
			DataStorage: model.DataStorageLoS{
				ProvisioningPolicy: "Thin",
				IsSpaceEfficient:   true,
			},
			IOConnectivity: model.IOConnectivityLoS{
				AccessProtocol: "rbd",
				MaxIOPS:        1000,
			},
			Advanced: map[string]interface{}{
				"diskType":   "SSD",
				"throughput": float64(1000),
			},
		},
	}

	pol, err := fp.GetPool(polID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(pol, expected) {
		t.Errorf("Expected %v, got %v", expected, pol)
		return
	}
}

func TestListPools(t *testing.T) {
	expected := []*model.StoragePoolSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
			},
			Name:          "sample-pool-01",
			Description:   "This is the first sample storage pool for testing",
			StorageType:   "block",
			TotalCapacity: int64(100),
			FreeCapacity:  int64(90),
			DockId:        "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			Extras: model.StoragePoolExtraSpec{
				DataStorage: model.DataStorageLoS{
					ProvisioningPolicy: "Thin",
					IsSpaceEfficient:   true,
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "rbd",
					MaxIOPS:        8000000,
					MaxBWS:         700,
				},
				Advanced: map[string]interface{}{
					"diskType": "SSD",
					"latency":  "3ms",
				},
			},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "a594b8ac-a103-11e7-985f-d723bcf01b5f",
			},
			Name:          "sample-pool-02",
			Description:   "This is the second sample storage pool for testing",
			StorageType:   "block",
			TotalCapacity: int64(200),
			FreeCapacity:  int64(170),
			DockId:        "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			Extras: model.StoragePoolExtraSpec{
				DataStorage: model.DataStorageLoS{
					ProvisioningPolicy: "Thin",
					IsSpaceEfficient:   true,
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "rbd",
					MaxIOPS:        3000000,
					MaxBWS:         350,
				},
				Advanced: map[string]interface{}{
					"diskType": "SAS",
					"latency":  "500ms",
				},
			},
		},
	}

	pols, err := fp.ListPools()
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(pols, expected) {
		t.Errorf("Expected %v, got %v", expected, pols)
		return
	}
}
