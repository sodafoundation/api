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

package discovery

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/pkg/utils/config"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
	fakedriver "github.com/opensds/opensds/testutils/driver"
)

const (
	expectedUuid      = "0e9c3c68-8a0b-11e7-94a7-67f755e235cb"
	expectedCreatedAt = "2017-08-26T11:01:09"
	expectedUpdatedAt = "2017-08-26T11:01:55"
)

func init() {
	CONF.OsdsDock = OsdsDock{
		ApiEndpoint:     "localhost:50050",
		EnabledBackends: []string{"sample"},
		Backends: Backends{
			Sample: BackendProperties{
				Name:        "sample",
				Description: "Sample Test",
				DriverName:  "default",
			},
		},
	}
}

func NewFakeDiscoverer() *DockDiscoverer {
	return &DockDiscoverer{
	//		c: fakedb.NewFakeDbClient(),
	}
}

func TestInit(t *testing.T) {
	var dd = NewFakeDiscoverer()
	var expectedDocks = []*model.DockSpec{
		{
			BaseModel:   &model.BaseModel{},
			Name:        "sample",
			Description: "Sample Test",
			Endpoint:    "localhost:50050",
			DriverName:  "default",
		},
	}

	if err := dd.Init(); err != nil {
		t.Errorf("Failed to init discoverer struct: %v\n", err)
	}
	for _, dck := range dd.dcks {
		dck.Id = ""
	}
	if !reflect.DeepEqual(dd.dcks, expectedDocks) {
		t.Errorf("Expected %+v, got %+v\n", expectedDocks, dd.dcks)
	}
}

func TestDiscover(t *testing.T) {
	var dd = NewFakeDiscoverer()
	dd.dcks = []*model.DockSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			Name:        "sample",
			Description: "Sample Test",
			Endpoint:    "localhost:50050",
			DriverName:  "default",
		},
	}
	var expectedPools = []*model.StoragePoolSpec{
		{
			BaseModel:        &model.BaseModel{},
			Name:             "sample-pool-01",
			Description:      "This is the first sample storage pool for testing",
			TotalCapacity:    int64(100),
			FreeCapacity:     int64(90),
			AvailabilityZone: "default",
			DockId:           "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			Extras: model.ExtraSpec{
				"diskType": "SSD",
				"thin":     true,
			},
		},
		{
			BaseModel:        &model.BaseModel{},
			Name:             "sample-pool-02",
			Description:      "This is the second sample storage pool for testing",
			TotalCapacity:    int64(200),
			FreeCapacity:     int64(170),
			AvailabilityZone: "default",
			DockId:           "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			Extras: model.ExtraSpec{
				"diskType": "SAS",
				"thin":     true,
			},
		},
	}

	if err := dd.Discover(&fakedriver.Driver{}); err != nil {
		t.Errorf("Failed to discoverer pools: %v\n", err)
	}
	for _, pol := range dd.pols {
		pol.Id = ""
	}
	if !reflect.DeepEqual(dd.pols, expectedPools) {
		t.Errorf("Expected %+v, got %+v\n", expectedPools, dd.pols)
	}
}

func TestStore(t *testing.T) {
	var dd = NewFakeDiscoverer()
	dd.dcks = []*model.DockSpec{
		{
			BaseModel:   &model.BaseModel{},
			Name:        "sample",
			Description: "Sample Test",
			Endpoint:    "localhost:50050",
			DriverName:  "default",
		},
	}
	dd.pols = []*model.StoragePoolSpec{
		{
			BaseModel:        &model.BaseModel{},
			Name:             "sample-pool-01",
			Description:      "This is the first sample storage pool for testing",
			TotalCapacity:    int64(100),
			FreeCapacity:     int64(90),
			AvailabilityZone: "default",
			DockId:           "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			Extras: model.ExtraSpec{
				"diskType": "SSD",
				"thin":     true,
			},
		},
		{
			BaseModel:        &model.BaseModel{},
			Name:             "sample-pool-02",
			Description:      "This is the second sample storage pool for testing",
			TotalCapacity:    int64(200),
			FreeCapacity:     int64(170),
			AvailabilityZone: "default",
			DockId:           "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			Extras: model.ExtraSpec{
				"diskType": "SAS",
				"thin":     true,
			},
		},
	}

	mockClient := new(dbtest.MockClient)
	mockClient.On("CreateDock", dd.dcks[0]).Return(nil, nil)
	mockClient.On("CreatePool", dd.pols[0]).Return(nil, nil)
	mockClient.On("CreatePool", dd.pols[1]).Return(nil, nil)
	dd.c = mockClient

	if err := dd.Store(); err != nil {
		t.Errorf("Failed to store docks and pools into database: %v\n", err)
	}
}
