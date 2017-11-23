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
	"github.com/opensds/opensds/pkg/utils"
	. "github.com/opensds/opensds/pkg/utils/config"
	fakedb "github.com/opensds/opensds/testutils/db"
	fakedriver "github.com/opensds/opensds/testutils/driver"
	fakesetter "github.com/opensds/opensds/testutils/utils/testing"
)

var expectedSetter = fakesetter.NewFakeSetter()

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
	utils.S = expectedSetter
}

func NewFakeDiscoverer() *DockDiscoverer {
	return &DockDiscoverer{
		c: fakedb.NewFakeDbClient(),
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
			Parameters: map[string]interface{}{
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
			Parameters: map[string]interface{}{
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
			Parameters: map[string]interface{}{
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
			Parameters: map[string]interface{}{
				"diskType": "SAS",
				"thin":     true,
			},
		},
	}

	if err := dd.Store(); err != nil {
		t.Errorf("Failed to store docks and pools into database: %v\n", err)
	}
	for _, dck := range dd.dcks {
		if !reflect.DeepEqual(dck.Id, expectedSetter.Uuid) {
			t.Errorf("Assert dock id: expected %v, got %v\n", expectedSetter.Uuid, dck.Id)
		}
		if !reflect.DeepEqual(dck.CreatedAt, expectedSetter.CreatedTime) {
			t.Errorf("Assert dock create time: expected %v, got %v\n",
				expectedSetter.CreatedTime, dck.CreatedAt)
		}
		if !reflect.DeepEqual(dck.Id, expectedSetter.Uuid) {
			t.Errorf("Assert dock update time: expected %v, got %v\n", expectedSetter.UpdatedTime, dck.UpdatedAt)
		}
	}
	for _, pol := range dd.pols {
		if !reflect.DeepEqual(pol.Id, expectedSetter.Uuid) {
			t.Errorf("Assert pool id: expected %v, got %v\n", expectedSetter.Uuid, pol.Id)
		}
		if !reflect.DeepEqual(pol.CreatedAt, expectedSetter.CreatedTime) {
			t.Errorf("Assert pool create time: expected %v, got %v\n",
				expectedSetter.CreatedTime, pol.CreatedAt)
		}
		if !reflect.DeepEqual(pol.Id, expectedSetter.Uuid) {
			t.Errorf("Assert pool update time: expected %v, got %v\n", expectedSetter.UpdatedTime, pol.UpdatedAt)
		}
	}
}
