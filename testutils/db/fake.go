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

// FakeDbClient
type FakeDbClient struct{}

// NewFakeDbClient
func NewFakeDbClient() *FakeDbClient {
	return &FakeDbClient{}
}

// CreateDock
func (fc *FakeDbClient) CreateDock(dck *model.DockSpec) (*model.DockSpec, error) {
	return &SampleDocks[0], nil
}

// GetDock
func (fc *FakeDbClient) GetDock(dckID string) (*model.DockSpec, error) {
	for _, dock := range SampleDocks {
		if dock.Id == dckID {
			return &dock, nil
		}
	}

	return nil, errors.New("Can't find this dock resource!")
}

// GetDockByPoolId
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

// ListDocks
func (fc *FakeDbClient) ListDocks() ([]*model.DockSpec, error) {
	var dcks []*model.DockSpec

	for i := range SampleDocks {
		dcks = append(dcks, &SampleDocks[i])
	}
	return dcks, nil
}

// UpdateDock
func (fc *FakeDbClient) UpdateDock(dckID, name, desp string) (*model.DockSpec, error) {
	return nil, nil
}

// DeleteDock
func (fc *FakeDbClient) DeleteDock(dckID string) error {
	return nil
}

func (fc *FakeDbClient) CreatePool(pol *model.StoragePoolSpec) (*model.StoragePoolSpec, error) {
	return &SamplePools[0], nil
}

// GetPool
func (fc *FakeDbClient) GetPool(polID string) (*model.StoragePoolSpec, error) {
	for _, pool := range SamplePools {
		if pool.Id == polID {
			return &pool, nil
		}
	}

	return nil, errors.New("Can't find this pool resource!")
}

// ListPools
func (fc *FakeDbClient) ListPools() ([]*model.StoragePoolSpec, error) {
	var pols []*model.StoragePoolSpec

	for i := range SamplePools {
		pols = append(pols, &SamplePools[i])
	}
	return pols, nil
}

// UpdatePool
func (fc *FakeDbClient) UpdatePool(polID, name, desp string, usedCapacity int64, used bool) (*model.StoragePoolSpec, error) {
	return nil, nil
}

// DeletePool
func (fc *FakeDbClient) DeletePool(polID string) error {
	return nil
}

// CreateProfile
func (fc *FakeDbClient) CreateProfile(prf *model.ProfileSpec) (*model.ProfileSpec, error) {
	return &SampleProfiles[0], nil
}

// GetProfile
func (fc *FakeDbClient) GetProfile(prfID string) (*model.ProfileSpec, error) {
	for _, profile := range SampleProfiles {
		if profile.Id == prfID {
			return &profile, nil
		}
	}

	return nil, errors.New("Can't find this profile resource!")
}

// GetDefaultProfile
func (fc *FakeDbClient) GetDefaultProfile() (*model.ProfileSpec, error) {
	for _, profile := range SampleProfiles {
		if profile.Name == "default" {
			return &profile, nil
		}
	}

	return nil, errors.New("Can't find default profile resource!")
}

// ListProfiles
func (fc *FakeDbClient) ListProfiles() ([]*model.ProfileSpec, error) {
	var prfs []*model.ProfileSpec

	for i := range SampleProfiles {
		prfs = append(prfs, &SampleProfiles[i])
	}
	return prfs, nil
}

// UpdateProfile
func (fc *FakeDbClient) UpdateProfile(prfID string, input *model.ProfileSpec) (*model.ProfileSpec, error) {
	return nil, nil
}

// DeleteProfile
func (fc *FakeDbClient) DeleteProfile(prfID string) error {
	return nil
}

// AddExtraProperty
func (fc *FakeDbClient) AddExtraProperty(prfID string, ext model.ExtraSpec) (*model.ExtraSpec, error) {
	extra := SampleProfiles[0].Extras
	return &extra, nil
}

// ListExtraProperties
func (fc *FakeDbClient) ListExtraProperties(prfID string) (*model.ExtraSpec, error) {
	extra := SampleProfiles[0].Extras
	return &extra, nil
}

// RemoveExtraProperty
func (fc *FakeDbClient) RemoveExtraProperty(prfID, extraKey string) error {
	return nil
}

// CreateVolume
func (fc *FakeDbClient) CreateVolume(vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	return &SampleVolumes[0], nil
}

// GetVolume
func (fc *FakeDbClient) GetVolume(volID string) (*model.VolumeSpec, error) {
	return &SampleVolumes[0], nil
}

// ListVolumes
func (fc *FakeDbClient) ListVolumes() ([]*model.VolumeSpec, error) {
	var vols []*model.VolumeSpec

	vols = append(vols, &SampleVolumes[0])
	return vols, nil
}

// UpdateVolume
func (fc *FakeDbClient) UpdateVolume(volID string, vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	return &SampleVolumes[0], nil
}

// DeleteVolume
func (fc *FakeDbClient) DeleteVolume(volID string) error {
	return nil
}

// CreateVolumeAttachment
func (fc *FakeDbClient) CreateVolumeAttachment(attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	return &SampleAttachments[0], nil
}

// GetVolumeAttachment
func (fc *FakeDbClient) GetVolumeAttachment(attachmentId string) (*model.VolumeAttachmentSpec, error) {
	return &SampleAttachments[0], nil
}

// ListVolumeAttachments
func (fc *FakeDbClient) ListVolumeAttachments(volumeId string) ([]*model.VolumeAttachmentSpec, error) {
	var atcs []*model.VolumeAttachmentSpec

	atcs = append(atcs, &SampleAttachments[0])
	return atcs, nil
}

// UpdateVolumeAttachment
func (fc *FakeDbClient) UpdateVolumeAttachment(attachmentId string, attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	return nil, nil
}

// DeleteVolumeAttachment
func (fc *FakeDbClient) DeleteVolumeAttachment(attachmentId string) error {
	return nil
}

// CreateVolumeSnapshot
func (fc *FakeDbClient) CreateVolumeSnapshot(vs *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	return &SampleSnapshots[0], nil
}

// GetVolumeSnapshot
func (fc *FakeDbClient) GetVolumeSnapshot(snapshotID string) (*model.VolumeSnapshotSpec, error) {
	return &SampleSnapshots[0], nil
}

// ListVolumeSnapshots
func (fc *FakeDbClient) ListVolumeSnapshots() ([]*model.VolumeSnapshotSpec, error) {
	var snps []*model.VolumeSnapshotSpec

	snps = append(snps, &SampleSnapshots[0], &SampleSnapshots[1])
	return snps, nil
}

// UpdateVolumeSnapshot
func (fc *FakeDbClient) UpdateVolumeSnapshot(snapshotID string, vs *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	return &SampleSnapshots[0], nil
}

// DeleteVolumeSnapshot
func (fc *FakeDbClient) DeleteVolumeSnapshot(snapshotID string) error {
	return nil
}
