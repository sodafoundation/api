// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
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

/*
This module implements the policy-based scheduling by parsing storage
profiles configured by admin.

*/

package controller

import (
	"errors"
	"reflect"
	"testing"

	api "github.com/opensds/opensds/pkg/model"
)

type fakeDbClient struct{}

func (fc *fakeDbClient) CreateDock(dck *api.DockSpec) (*api.DockSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) GetDock(dckID string) (*api.DockSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) ListDocks() (*[]api.DockSpec, error) {
	return &sampleDocks, nil
}

func (fc *fakeDbClient) UpdateDock(dckID, name, desp string) (*api.DockSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) DeleteDock(dckID string) error {
	return errors.New("Not implemented!")
}

func (fc *fakeDbClient) CreatePool(pol *api.StoragePoolSpec) (*api.StoragePoolSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) GetPool(polID string) (*api.StoragePoolSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) ListPools() (*[]api.StoragePoolSpec, error) {
	return &samplePools, nil
}

func (fc *fakeDbClient) UpdatePool(polID, name, desp string, usedCapacity int64, used bool) (*api.StoragePoolSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) DeletePool(polID string) error {
	return errors.New("Not implemented!")
}

func (fc *fakeDbClient) CreateProfile(prf *api.ProfileSpec) (*api.ProfileSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) GetProfile(prfID string) (*api.ProfileSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) ListProfiles() (*[]api.ProfileSpec, error) {
	return &sampleProfiles, nil
}

func (fc *fakeDbClient) UpdateProfile(prfID string, input *api.ProfileSpec) (*api.ProfileSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) DeleteProfile(prfID string) error {
	return errors.New("Not implemented!")
}

func (fc *fakeDbClient) AddExtraProperty(prfID string, ext api.ExtraSpec) (*api.ExtraSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) ListExtraProperties(prfID string) (*api.ExtraSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) RemoveExtraProperty(prfID, extraKey string) error {
	return errors.New("Not implemented!")
}

func (fc *fakeDbClient) CreateVolume(vol *api.VolumeSpec) (*api.VolumeSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) GetVolume(volID string) (*api.VolumeSpec, error) {
	return &sampleVolume, nil
}

func (fc *fakeDbClient) ListVolumes() (*[]api.VolumeSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) DeleteVolume(volID string) error {
	return errors.New("Not implemented!")
}

func (fc *fakeDbClient) CreateVolumeAttachment(volID string, atc *api.VolumeAttachmentSpec) (*api.VolumeAttachmentSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) GetVolumeAttachment(volID, attachmentID string) (*api.VolumeAttachmentSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) ListVolumeAttachments(volID string) (*[]api.VolumeAttachmentSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) UpdateVolumeAttachment(volID, attachmentID, mountpoint string, hostInfo *api.HostInfo) (*api.VolumeAttachmentSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) DeleteVolumeAttachment(volID, attachmentID string) error {
	return errors.New("Not implemented!")
}

func (fc *fakeDbClient) CreateVolumeSnapshot(vs *api.VolumeSnapshotSpec) (*api.VolumeSnapshotSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) GetVolumeSnapshot(snapshotID string) (*api.VolumeSnapshotSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) ListVolumeSnapshots() (*[]api.VolumeSnapshotSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (fc *fakeDbClient) DeleteVolumeSnapshot(snapshotID string) error {
	return errors.New("Not implemented!")
}

func NewFakeDbSearcher() Searcher {
	var fc *fakeDbClient

	return &DbSearcher{
		Client: fc,
	}
}

func TestSearchProfile(t *testing.T) {
	s := NewFakeDbSearcher()
	var expectedDefaultProfile = api.ProfileSpec{
		BaseModel: &api.BaseModel{
			Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		Name:        "default",
		Description: "default policy",
		Extra:       api.ExtraSpec{},
	}
	var expectedAssignedProfile = api.ProfileSpec{
		BaseModel: &api.BaseModel{
			Id: "2f9c0a04-66ef-11e7-ade2-43158893e017",
		},
		Name:        "ceph",
		Description: "ceph policy",
		Extra: api.ExtraSpec{
			"highAvailability":     "true",
			"intervalSnapshot":     "10s",
			"deleteSnapshotPolicy": "true",
		},
	}

	// Test if the method would return default profile when no profile id assigned.
	prf, err := s.SearchProfile("")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&expectedDefaultProfile, prf) {
		t.Fatalf("Expected %v, get %v", &expectedDefaultProfile, prf)
	}

	// Test if the method would return specified profile when profile id assigned.
	prf, err = s.SearchProfile("2f9c0a04-66ef-11e7-ade2-43158893e017")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&expectedAssignedProfile, prf) {
		t.Fatalf("Expected %v, get %v", &expectedAssignedProfile, prf)
	}
}

func TestSearchSupportedPool(t *testing.T) {
	s := NewFakeDbSearcher()
	var expectedPool = api.StoragePoolSpec{
		BaseModel: &api.BaseModel{
			Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		},
		Name:          "rbd-pool",
		Description:   "ceph pool1",
		StorageType:   "block",
		DockId:        "076454a8-65da-11e7-9a65-5f5d9b935b9f",
		TotalCapacity: 200,
		FreeCapacity:  200,
		Parameters: map[string]interface{}{
			"thinProvision":    "false",
			"highAvailability": "true",
		},
	}
	var inputTag = map[string]string{
		"highAvailability": "true",
	}

	// Test if the method would return correct pool when storage tag assigned.
	pol, err := s.SearchSupportedPool(inputTag)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&expectedPool, pol) {
		t.Fatalf("Expected %v, get %v", &expectedPool, pol)
	}
}

func TestSearchDockByPool(t *testing.T) {
	s := NewFakeDbSearcher()
	var expectedDock = api.DockSpec{
		BaseModel: &api.BaseModel{
			Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		},
		Name:        "cinder",
		Description: "cinder backend service",
		Endpoint:    "127.0.0.1",
		DriverName:  "cinder",
		Parameters: map[string]interface{}{
			"thinProvision":    "true",
			"highAvailability": "false",
		},
	}
	var inputPool = api.StoragePoolSpec{
		BaseModel: &api.BaseModel{
			Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
		},
		Name:          "cinder-pool",
		Description:   "cinder pool1",
		StorageType:   "block",
		DockId:        "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		TotalCapacity: 100,
		FreeCapacity:  100,
		Parameters: map[string]interface{}{
			"thinProvision":    "true",
			"highAvailability": "false",
		},
	}

	// Test if the method would return correct dock when storage pool assigned.
	dck, err := s.SearchDockByPool(&inputPool)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&expectedDock, dck) {
		t.Fatalf("Expected %v, get %v", &expectedDock, dck)
	}
}

func TestSearchDockByVolume(t *testing.T) {
	s := NewFakeDbSearcher()
	var expectedDock = api.DockSpec{
		BaseModel: &api.BaseModel{
			Id: "076454a8-65da-11e7-9a65-5f5d9b935b9f",
		},
		Name:        "ceph",
		Description: "ceph backend service",
		Endpoint:    "127.0.0.1",
		DriverName:  "ceph",
		Parameters: map[string]interface{}{
			"thinProvision":    "false",
			"highAvailability": "true",
		},
	}
	var inputVolID = "9193c3ec-771f-11e7-8ca3-d32c0a8b2725"

	// Test if the method would return correct dock when volume id assigned.
	dck, err := s.SearchDockByVolume(inputVolID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&expectedDock, dck) {
		t.Fatalf("Expected %v, get %v", &expectedDock, dck)
	}
}

var (
	sampleProfiles = []api.ProfileSpec{
		{
			BaseModel: &api.BaseModel{
				Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
			},
			Name:        "default",
			Description: "default policy",
			Extra:       api.ExtraSpec{},
		},
		{
			BaseModel: &api.BaseModel{
				Id: "2f9c0a04-66ef-11e7-ade2-43158893e017",
			},
			Name:        "ceph",
			Description: "ceph policy",
			Extra: api.ExtraSpec{
				"highAvailability":     "true",
				"intervalSnapshot":     "10s",
				"deleteSnapshotPolicy": "true",
			},
		},
	}

	sampleDocks = []api.DockSpec{
		{
			BaseModel: &api.BaseModel{
				Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			Name:        "cinder",
			Description: "cinder backend service",
			Endpoint:    "127.0.0.1",
			DriverName:  "cinder",
			Parameters: map[string]interface{}{
				"thinProvision":    "true",
				"highAvailability": "false",
			},
		},
		{
			BaseModel: &api.BaseModel{
				Id: "076454a8-65da-11e7-9a65-5f5d9b935b9f",
			},
			Name:        "ceph",
			Description: "ceph backend service",
			Endpoint:    "127.0.0.1",
			DriverName:  "ceph",
			Parameters: map[string]interface{}{
				"thinProvision":    "false",
				"highAvailability": "true",
			},
		},
	}

	samplePools = []api.StoragePoolSpec{
		{
			BaseModel: &api.BaseModel{
				Id: "6edc7604-7725-11e7-b2b1-1335d0254e7c",
			},
			Name:          "cinder-pool",
			Description:   "cinder pool1",
			StorageType:   "block",
			DockId:        "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			TotalCapacity: 100,
			FreeCapacity:  100,
			Parameters: map[string]interface{}{
				"thinProvision":    "true",
				"highAvailability": "false",
			},
		},
		{
			BaseModel: &api.BaseModel{
				Id: "80287bf8-66de-11e7-b031-f3b0af1675ba",
			},
			Name:          "rbd-pool",
			Description:   "ceph pool1",
			StorageType:   "block",
			DockId:        "076454a8-65da-11e7-9a65-5f5d9b935b9f",
			TotalCapacity: 200,
			FreeCapacity:  200,
			Parameters: map[string]interface{}{
				"thinProvision":    "false",
				"highAvailability": "true",
			},
		},
	}

	sampleVolume = api.VolumeSpec{
		BaseModel: &api.BaseModel{
			Id:        "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
			CreatedAt: "2017-08-02T09:17:05",
		},
		Name:        "fake-volume",
		Description: "fake volume for testing",
		Size:        1,
		PoolId:      "80287bf8-66de-11e7-b031-f3b0af1675ba",
	}
)
