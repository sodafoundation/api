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

package db

import (
	"errors"

	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/testutils/collection"
)

type FakeDbClient struct{}

func NewFakeDbClient() *FakeDbClient {
	return &FakeDbClient{}
}

func (fc *FakeDbClient) CreateDock(dck *model.DockSpec) (*model.DockSpec, error) {
	return &SampleDocks[0], nil
}

func (fc *FakeDbClient) GetDock(dckID string) (*model.DockSpec, error) {
	for _, dock := range SampleDocks {
		if dock.Id == dckID {
			return &dock, nil
		}
	}

	return nil, errors.New("Can't find this dock resource!")
}
func (fc *FakeDbClient) GetDockByPoolId(poolId string) (*model.DockSpec, error) {
	pool, err := fc.GetPool(poolId)
	if err != nil {
		return nil, err
	}
	for _, dock := range SampleDocks {
		if dock.Id == pool.DockId {
			return &dock, nil
		}
	}
	return nil, errors.New("Can't find this dock resource by pool id!")
}

func (fc *FakeDbClient) ListDocks() ([]*model.DockSpec, error) {
	var dcks []*model.DockSpec

	for i := range SampleDocks {
		dcks = append(dcks, &SampleDocks[i])
	}
	return dcks, nil
}

func (fc *FakeDbClient) UpdateDock(dckID, name, desp string) (*model.DockSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) DeleteDock(dckID string) error {
	return nil
}

func (fc *FakeDbClient) CreatePool(pol *model.StoragePoolSpec) (*model.StoragePoolSpec, error) {
	return &SamplePools[0], nil
}

func (fc *FakeDbClient) GetPool(polID string) (*model.StoragePoolSpec, error) {
	for _, pool := range SamplePools {
		if pool.Id == polID {
			return &pool, nil
		}
	}

	return nil, errors.New("Can't find this pool resource!")
}

func (fc *FakeDbClient) ListPools() ([]*model.StoragePoolSpec, error) {
	var pols []*model.StoragePoolSpec

	for i := range SamplePools {
		pols = append(pols, &SamplePools[i])
	}
	return pols, nil
}

func (fc *FakeDbClient) UpdatePool(polID, name, desp string, usedCapacity int64, used bool) (*model.StoragePoolSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) DeletePool(polID string) error {
	return nil
}

func (fc *FakeDbClient) CreateProfile(prf *model.ProfileSpec) (*model.ProfileSpec, error) {
	return &SampleProfiles[0], nil
}

func (fc *FakeDbClient) GetProfile(prfID string) (*model.ProfileSpec, error) {
	for _, profile := range SampleProfiles {
		if profile.Id == prfID {
			return &profile, nil
		}
	}

	return nil, errors.New("Can't find this profile resource!")
}

func (fc *FakeDbClient) GetDefaultProfile() (*model.ProfileSpec, error) {
	for _, profile := range SampleProfiles {
		if profile.Name == "default" {
			return &profile, nil
		}
	}

	return nil, errors.New("Can't find default profile resource!")
}

func (fc *FakeDbClient) ListProfiles() ([]*model.ProfileSpec, error) {
	var prfs []*model.ProfileSpec

	for i := range SampleProfiles {
		prfs = append(prfs, &SampleProfiles[i])
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
	extra := SampleProfiles[0].Extras
	return &extra, nil
}

func (fc *FakeDbClient) ListExtraProperties(prfID string) (*model.ExtraSpec, error) {
	extra := SampleProfiles[0].Extras
	return &extra, nil
}

func (fc *FakeDbClient) RemoveExtraProperty(prfID, extraKey string) error {
	return nil
}

func (fc *FakeDbClient) CreateVolume(vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	return &SampleVolumes[0], nil
}

func (fc *FakeDbClient) GetVolume(volID string) (*model.VolumeSpec, error) {
	return &SampleVolumes[0], nil
}

func (fc *FakeDbClient) ListVolumes() ([]*model.VolumeSpec, error) {
	var vols []*model.VolumeSpec

	vols = append(vols, &SampleVolumes[0])
	return vols, nil
}

func (fc *FakeDbClient) UpdateVolume(volID string, vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	return &SampleVolumes[0], nil
}

func (fc *FakeDbClient) DeleteVolume(volID string) error {
	return nil
}

func (fc *FakeDbClient) CreateVolumeAttachment(attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	return &SampleAttachments[0], nil
}

func (fc *FakeDbClient) GetVolumeAttachment(attachmentId string) (*model.VolumeAttachmentSpec, error) {
	return &SampleAttachments[0], nil
}

func (fc *FakeDbClient) ListVolumeAttachments(volumeId string) ([]*model.VolumeAttachmentSpec, error) {
	var atcs []*model.VolumeAttachmentSpec

	atcs = append(atcs, &SampleAttachments[0])
	return atcs, nil
}

func (fc *FakeDbClient) UpdateVolumeAttachment(attachmentId string, attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) DeleteVolumeAttachment(attachmentId string) error {
	return nil
}

func (fc *FakeDbClient) CreateVolumeSnapshot(vs *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	return &SampleSnapshots[0], nil
}

func (fc *FakeDbClient) GetVolumeSnapshot(snapshotID string) (*model.VolumeSnapshotSpec, error) {
	return &SampleSnapshots[0], nil
}

func (fc *FakeDbClient) ListVolumeSnapshots() ([]*model.VolumeSnapshotSpec, error) {
	var snps []*model.VolumeSnapshotSpec

	snps = append(snps, &SampleSnapshots[0], &SampleSnapshots[1])
	return snps, nil
}

func (fc *FakeDbClient) UpdateVolumeSnapshot(snapshotID string, vs *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	return &SampleSnapshots[0], nil
}

func (fc *FakeDbClient) DeleteVolumeSnapshot(snapshotID string) error {
	return nil
}
