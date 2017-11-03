// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package client

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/opensds/opensds/pkg/model"
)

func NewFakePoolReceiver() Receiver {
	return &fakePoolReceiver{}
}

type fakePoolReceiver struct{}

func (*fakePoolReceiver) Recv(
	f reqFunc,
	string,
	method string,
	in interface{},
	out interface{},
) error {
	if strings.ToUpper(method) != "GET" {
		return errors.New("method not supported!")
	}

	switch out.(type) {
	case *model.StoragePoolSpec:
		if err := json.Unmarshal([]byte(samplePool), out); err != nil {
			return err
		}
		break
	case *[]*model.StoragePoolSpec:
		if err := json.Unmarshal([]byte(samplePools), out); err != nil {
			return err
		}
		break
	default:
		return errors.New("output format not supported!")
	}

	return nil
}

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
		TotalCapacity: int64(100),
		FreeCapacity:  int64(90),
		DockId:        "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		Parameters: map[string]interface{}{
			"diskType":  "SSD",
			"iops":      float64(1000),
			"bandwidth": float64(1000),
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
			TotalCapacity: int64(100),
			FreeCapacity:  int64(90),
			DockId:        "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			Parameters: map[string]interface{}{
				"diskType":  "SSD",
				"iops":      float64(1000),
				"bandwidth": float64(1000),
			},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "a594b8ac-a103-11e7-985f-d723bcf01b5f",
			},
			Name:          "sample-pool-02",
			Description:   "This is the second sample storage pool for testing",
			TotalCapacity: int64(200),
			FreeCapacity:  int64(170),
			DockId:        "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			Parameters: map[string]interface{}{
				"diskType":  "SAS",
				"iops":      float64(800),
				"bandwidth": float64(800),
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

var (
	samplePool = `{
		"id": "084bf71e-a102-11e7-88a8-e31fe6d52248",
		"name": "sample-pool-01",
		"description": "This is the first sample storage pool for testing",
		"totalCapacity": 100,
		"freeCapacity": 90,
		"dockId": "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		"extras": {
			"diskType": "SSD",
			"iops": 1000,
			"bandwidth": 1000
		}
	}`

	samplePools = `[
		{
			"id": "084bf71e-a102-11e7-88a8-e31fe6d52248",
			"name": "sample-pool-01",
			"description": "This is the first sample storage pool for testing",
			"totalCapacity": 100,
			"freeCapacity": 90,
			"dockId": "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			"extras": {
				"diskType": "SSD",
				"iops": 1000,
				"bandwidth": 1000
			}
		},
		{
			"id": "a594b8ac-a103-11e7-985f-d723bcf01b5f",
			"name": "sample-pool-02",
			"description": "This is the second sample storage pool for testing",
			"totalCapacity": 200,
			"freeCapacity": 170,
			"dockId": "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			"extras": {
				"diskType": "SAS",
				"iops": 800,
				"bandwidth": 800
			}
		}
	]`
)
