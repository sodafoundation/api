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

package db

import (
	"errors"

	"github.com/opensds/opensds/pkg/model"
)

type FakeDbClient struct{}

func NewFakeDbClient() Client {
	return &FakeDbClient{}
}

func (fc *FakeDbClient) CreateDock(dck *model.DockSpec) error {
	return nil
}

func (fc *FakeDbClient) GetDock(dckID string) (*model.DockSpec, error) {
	for i := range sampleDocks {
		if sampleDocks[i].GetId() == dckID {
			return &sampleDocks[i], nil
		}
	}

	return nil, errors.New("Can't find this dock resource!")
}

func (fc *FakeDbClient) ListDocks() ([]*model.DockSpec, error) {
	var dcks []*model.DockSpec

	for i := range sampleDocks {
		dcks = append(dcks, &sampleDocks[i])
	}
	return dcks, nil
}

func (fc *FakeDbClient) UpdateDock(dckID, name, desp string) (*model.DockSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) DeleteDock(dckID string) error {
	return nil
}

func (fc *FakeDbClient) CreatePool(pol *model.StoragePoolSpec) error {
	return nil
}

func (fc *FakeDbClient) GetPool(polID string) (*model.StoragePoolSpec, error) {
	for i := range samplePools {
		if samplePools[i].GetId() == polID {
			return &samplePools[i], nil
		}
	}

	return nil, errors.New("Can't find this pool resource!")
}

func (fc *FakeDbClient) ListPools() ([]*model.StoragePoolSpec, error) {
	var pols []*model.StoragePoolSpec

	for i := range samplePools {
		pols = append(pols, &samplePools[i])
	}
	return pols, nil
}

func (fc *FakeDbClient) UpdatePool(polID, name, desp string, usedCapacity int64, used bool) (*model.StoragePoolSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) DeletePool(polID string) error {
	return nil
}

func (fc *FakeDbClient) CreateProfile(prf *model.ProfileSpec) error {
	return nil
}

func (fc *FakeDbClient) GetProfile(prfID string) (*model.ProfileSpec, error) {
	for i := range sampleProfiles {
		if sampleProfiles[i].GetId() == prfID {
			return &sampleProfiles[i], nil
		}
	}

	return nil, errors.New("Can't find this profile resource!")
}

func (fc *FakeDbClient) ListProfiles() ([]*model.ProfileSpec, error) {
	var prfs []*model.ProfileSpec

	for i := range sampleProfiles {
		prfs = append(prfs, &sampleProfiles[i])
	}
	return prfs, nil
}

func (fc *FakeDbClient) UpdateProfile(prfID string, input *model.ProfileSpec) (*model.ProfileSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) DeleteProfile(prfID string) error {
	return nil
}

func (fc *FakeDbClient) AddExtraProperty(prfID string, ext model.ExtraSpec) (*model.ExtraSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) ListExtraProperties(prfID string) (*model.ExtraSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) RemoveExtraProperty(prfID, extraKey string) error {
	return nil
}

func (fc *FakeDbClient) CreateVolume(vol *model.VolumeSpec) error {
	return nil
}

func (fc *FakeDbClient) GetVolume(volID string) (*model.VolumeSpec, error) {
	return &sampleVolume, nil
}

func (fc *FakeDbClient) ListVolumes() ([]*model.VolumeSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) DeleteVolume(volID string) error {
	return nil
}

func (fc *FakeDbClient) CreateVolumeAttachment(volID string, atc *model.VolumeAttachmentSpec) error {
	return nil
}

func (fc *FakeDbClient) GetVolumeAttachment(volID, attachmentID string) (*model.VolumeAttachmentSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) ListVolumeAttachments(volID string) ([]*model.VolumeAttachmentSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) UpdateVolumeAttachment(volID, attachmentID, mountpoint string, hostInfo *model.HostInfo) (*model.VolumeAttachmentSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) DeleteVolumeAttachment(volID, attachmentID string) error {
	return nil
}

func (fc *FakeDbClient) CreateVolumeSnapshot(vs *model.VolumeSnapshotSpec) error {
	return nil
}

func (fc *FakeDbClient) GetVolumeSnapshot(snapshotID string) (*model.VolumeSnapshotSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) ListVolumeSnapshots() ([]*model.VolumeSnapshotSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) DeleteVolumeSnapshot(snapshotID string) error {
	return nil
}

var (
	sampleProfiles = []model.ProfileSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
			},
			Name:        "default",
			Description: "default policy",
			Extra:       model.ExtraSpec{},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "2f9c0a04-66ef-11e7-ade2-43158893e017",
			},
			Name:        "ceph",
			Description: "ceph policy",
			Extra: model.ExtraSpec{
				"highAvailability":     "true",
				"intervalSnapshot":     "10s",
				"deleteSnapshotPolicy": "true",
			},
		},
	}

	sampleDocks = []model.DockSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			Name:        "cinder",
			Description: "cinder backend service",
			Endpoint:    "localhost:50050",
			DriverName:  "cinder",
			Parameters: map[string]interface{}{
				"thinProvision":    "true",
				"highAvailability": "false",
			},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "076454a8-65da-11e7-9a65-5f5d9b935b9f",
			},
			Name:        "ceph",
			Description: "ceph backend service",
			Endpoint:    "localhost:50050",
			DriverName:  "ceph",
			Parameters: map[string]interface{}{
				"thinProvision":    "false",
				"highAvailability": "true",
			},
		},
	}

	samplePools = []model.StoragePoolSpec{
		{
			BaseModel: &model.BaseModel{
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
			BaseModel: &model.BaseModel{
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

	sampleVolume = model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id:        "9193c3ec-771f-11e7-8ca3-d32c0a8b2725",
			CreatedAt: "2017-08-02T09:17:05",
		},
		Name:        "fake-volume",
		Description: "fake volume for testing",
		Size:        1,
		PoolId:      "80287bf8-66de-11e7-b031-f3b0af1675ba",
	}
)
