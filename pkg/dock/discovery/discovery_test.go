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

	c "github.com/opensds/opensds/pkg/context"
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

func NewFakeDockDiscoverer() *provisionDockDiscoverer {
	return &provisionDockDiscoverer{
		DockRegister: &DockRegister{},
	}
}

func TestInit(t *testing.T) {
	var fdd = NewFakeDockDiscoverer()
	var expected []*model.DockSpec

	for i := range SampleDocks {
		expected = append(expected, &SampleDocks[i])
	}
	name := map[string][]string{"Name": {SampleDocks[0].Name}}
	mockClient := new(dbtest.Client)
	mockClient.On("ListDocksWithFilter", c.NewAdminContext(), name).Return(expected, nil)
	fdd.c = mockClient
	if err := fdd.Init(); err != nil {
		t.Errorf("Failed to init discoverer struct: %v\n", err)
	}
	for i := range fdd.dcks {
		fdd.dcks[i].Id = ""
		fdd.dcks[i].NodeId = ""
		fdd.dcks[i].Metadata = nil
		expected[i].Id = ""
	}
	if !reflect.DeepEqual(fdd.dcks, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, fdd.dcks)
	}
}

func TestDiscover(t *testing.T) {
	var fdd = NewFakeDockDiscoverer()
	var expected []*model.StoragePoolSpec

	for i := range SampleDocks {
		fdd.dcks = append(fdd.dcks, &SampleDocks[i])
	}
	for i := range SamplePools {
		fdd.pols = append(fdd.pols, &SamplePools[i])
		expected = append(expected, &SamplePools[i])
	}
	m1 := map[string][]string{
		"Name":   {SamplePools[0].Name},
		"DockId": {""},
	}
	m2 := map[string][]string{
		"Name":   {SamplePools[1].Name},
		"DockId": {""},
	}
	m3 := map[string][]string{
		"Name":   {SamplePools[2].Name},
		"DockId": {""},
	}

	mockClient := new(dbtest.Client)
	mockClient.On("ListPools", c.NewAdminContext()).Return(fdd.pols, nil)
	mockClient.On("ListPoolsWithFilter", c.NewAdminContext(), m1).Return(expected, nil)
	mockClient.On("ListPoolsWithFilter", c.NewAdminContext(), m2).Return(expected, nil)
	mockClient.On("ListPoolsWithFilter", c.NewAdminContext(), m3).Return(expected, nil)
	fdd.c = mockClient

	if err := fdd.Discover(); err != nil {
		t.Errorf("Failed to discoverer pools: %v\n", err)
	}
	for _, pol := range fdd.pols {
		pol.Id = ""
	}
	if !reflect.DeepEqual(fdd.pols, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, fdd.pols)
	}
}

func TestReport(t *testing.T) {
	var fdd = NewFakeDockDiscoverer()
	var azs []*model.AvailabilityZoneSpec
	var az *model.AvailabilityZoneSpec

	az = &model.AvailabilityZoneSpec {
		BaseModel: &model.BaseModel {
			Id: "",
		},
		Name: "default",
	}

	for i := range SampleDocks {
		fdd.dcks = append(fdd.dcks, &SampleDocks[i])
	}
	for i := range SamplePools {
		fdd.pols = append(fdd.pols, &SamplePools[i])
	}
	for i := range SampleAvailabilityZones {
		azs = append(azs, &SampleAvailabilityZones[i])
	}

	mockClient := new(dbtest.Client)
	mockClient.On("CreateDock", c.NewAdminContext(), fdd.dcks[0]).Return(nil, nil)
	mockClient.On("CreatePool", c.NewAdminContext(), fdd.pols[0]).Return(nil, nil)
	mockClient.On("CreatePool", c.NewAdminContext(), fdd.pols[1]).Return(nil, nil)
	mockClient.On("CreatePool", c.NewAdminContext(), fdd.pols[2]).Return(nil, nil)
	mockClient.On("ListAvailabilityZones", c.NewAdminContext()).Return(azs, nil)
	mockClient.On("CreateAvailabilityZone", c.NewAdminContext(), az).Return(&SampleAvailabilityZones[0], nil)
	fdd.c = mockClient

	if err := fdd.Report(); err != nil {
		t.Errorf("Failed to store docks and pools into database: %v\n", err)
	}
}
