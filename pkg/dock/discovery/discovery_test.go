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

package discovery

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/pkg/utils/config"
	. "github.com/opensds/opensds/testutils/collection"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
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
				Description: "sample backend service",
				DriverName:  "sample",
			},
		},
	}
}

func NewFakeDiscoverer() *DockDiscoverer {
	return &DockDiscoverer{}
}

func TestInit(t *testing.T) {
	var dd = NewFakeDiscoverer()
	var expected []*model.DockSpec

	for i := range SampleDocks {
		expected = append(expected, &SampleDocks[i])
	}
	if err := dd.Init(); err != nil {
		t.Errorf("Failed to init discoverer struct: %v\n", err)
	}
	for i := range dd.dcks {
		dd.dcks[i].Id = ""
		expected[i].Id = ""
	}
	if !reflect.DeepEqual(dd.dcks, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, dd.dcks)
	}
}

func TestDiscover(t *testing.T) {
	var dd = NewFakeDiscoverer()
	var expected []*model.StoragePoolSpec

	for i := range SampleDocks {
		dd.dcks = append(dd.dcks, &SampleDocks[i])
	}
	for i := range SamplePools {
		expected = append(expected, &SamplePools[i])
	}
	if err := dd.Discover(); err != nil {
		t.Errorf("Failed to discoverer pools: %v\n", err)
	}
	for _, pol := range dd.pols {
		pol.Id = ""
	}
	if !reflect.DeepEqual(dd.pols, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, dd.pols)
	}
}

func TestStore(t *testing.T) {
	var dd = NewFakeDiscoverer()

	for i := range SampleDocks {
		dd.dcks = append(dd.dcks, &SampleDocks[i])
	}
	for i := range SamplePools {
		dd.pols = append(dd.pols, &SamplePools[i])
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
