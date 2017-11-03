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
	extra := sampleProfiles[0].Extra
	return &extra, nil
}

func (fc *FakeDbClient) ListExtraProperties(prfID string) (*model.ExtraSpec, error) {
	extra := sampleProfiles[0].Extra
	return &extra, nil
}

func (fc *FakeDbClient) RemoveExtraProperty(prfID, extraKey string) error {
	return nil
}

func (fc *FakeDbClient) CreateVolume(vol *model.VolumeSpec) error {
	return nil
}

func (fc *FakeDbClient) GetVolume(volID string) (*model.VolumeSpec, error) {
	return &sampleVolumes[0], nil
}

func (fc *FakeDbClient) ListVolumes() ([]*model.VolumeSpec, error) {
	var vols []*model.VolumeSpec

	vols = append(vols, &sampleVolumes[0])
	return vols, nil
}

func (fc *FakeDbClient) DeleteVolume(volID string) error {
	return nil
}

func (fc *FakeDbClient) CreateVolumeAttachment(volID string, atc *model.VolumeAttachmentSpec) error {
	return nil
}

func (fc *FakeDbClient) GetVolumeAttachment(volID, attachmentID string) (*model.VolumeAttachmentSpec, error) {
	return &sampleAttachments[0], nil
}

func (fc *FakeDbClient) ListVolumeAttachments(volID string) ([]*model.VolumeAttachmentSpec, error) {
	var atcs []*model.VolumeAttachmentSpec

	atcs = append(atcs, &sampleAttachments[0])
	return atcs, nil
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
	return &sampleSnapshots[0], nil
}

func (fc *FakeDbClient) ListVolumeSnapshots() ([]*model.VolumeSnapshotSpec, error) {
	var snps []*model.VolumeSnapshotSpec

	snps = append(snps, &sampleSnapshots[0], &sampleSnapshots[1])
	return snps, nil
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
			Name:        "silver",
			Description: "silver policy",
			Extra: model.ExtraSpec{
				"diskType":  "SAS",
				"iops":      300,
				"bandwidth": 500,
			},
		},
	}

	sampleDocks = []model.DockSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			Name:        "sample",
			Description: "sample backend service",
			Endpoint:    "localhost:50050",
			DriverName:  "sample",
		},
	}

	samplePools = []model.StoragePoolSpec{
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
				"iops":      1000,
				"bandwidth": 1000,
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
				"iops":      800,
				"bandwidth": 800,
			},
		},
	}

	sampleVolumes = []model.VolumeSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			},
			Name:        "sample-volume",
			Description: "This is a sample volume for testing",
			Size:        int64(1),
			Status:      "available",
			PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
	}

	sampleAttachments = []model.VolumeAttachmentSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "f2dda3d2-bf79-11e7-8665-f750b088f63e",
			},
			Name:        "sample-volume-attachment",
			Description: "This is a sample volume attachment for testing",
			Status:      "available",
			VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
			HostInfo:    &model.HostInfo{},
			ConnectionInfo: &model.ConnectionInfo{
				DriverVolumeType: "iscsi",
				ConnectionData: map[string]interface{}{
					"targetDiscovered": true,
					"targetIqn":        "iqn.2017-10.io.opensds:volume:00000001",
					"targetPortal":     "127.0.0.0.1:3260",
					"discard":          false,
				},
			},
		},
	}

	sampleSnapshots = []model.VolumeSnapshotSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "3769855c-a102-11e7-b772-17b880d2f537",
			},
			Name:        "sample-snapshot-01",
			Description: "This is the first sample snapshot for testing",
			Size:        int64(1),
			Status:      "created",
			VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "3bfaf2cc-a102-11e7-8ecb-63aea739d755",
			},
			Name:        "sample-snapshot-02",
			Description: "This is the second sample snapshot for testing",
			Size:        int64(1),
			Status:      "created",
			VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
		},
	}
)
