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

package db

import (
	"errors"

	c "github.com/opensds/opensds/pkg/context"
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
func (fc *FakeDbClient) CreateDock(ctx *c.Context, dck *model.DockSpec) (*model.DockSpec, error) {
	return &SampleDocks[0], nil
}

// GetDock
func (fc *FakeDbClient) GetDock(ctx *c.Context, dckID string) (*model.DockSpec, error) {
	for _, dock := range SampleDocks {
		if dock.Id == dckID {
			return &dock, nil
		}
	}

	return nil, errors.New("Can't find this dock resource!")
}

// GetDockByPoolId
func (fc *FakeDbClient) GetDockByPoolId(ctx *c.Context, poolId string) (*model.DockSpec, error) {
	pool, err := fc.GetPool(ctx, poolId)
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
func (fc *FakeDbClient) ListDocksWithFilter(ctx *c.Context, m map[string][]string) ([]*model.DockSpec, error) {
	var dcks []*model.DockSpec

	for i := range SampleDocks {
		dcks = append(dcks, &SampleDocks[i])
	}
	return dcks, nil
}
func (fc *FakeDbClient) ListDocks(ctx *c.Context) ([]*model.DockSpec, error) {
	var dcks []*model.DockSpec

	for i := range SampleDocks {
		dcks = append(dcks, &SampleDocks[i])
	}
	return dcks, nil
}

//ListAvailabilityZones
func (fc *FakeDbClient) ListAvailabilityZones(ctx *c.Context) ([]string, error) {
	var azs []string
	for i := range SamplePools {
		az := SamplePools[i].AvailabilityZone
		azs = append(azs, az)
	}
	return azs, nil
}

// UpdateDock
func (fc *FakeDbClient) UpdateDock(ctx *c.Context, dckID, name, desp string) (*model.DockSpec, error) {
	return nil, nil
}

// DeleteDock
func (fc *FakeDbClient) DeleteDock(ctx *c.Context, dckID string) error {
	return nil
}

func (fc *FakeDbClient) CreatePool(ctx *c.Context, pol *model.StoragePoolSpec) (*model.StoragePoolSpec, error) {
	return &SamplePools[0], nil
}

// GetPool
func (fc *FakeDbClient) GetPool(ctx *c.Context, polID string) (*model.StoragePoolSpec, error) {
	for _, pool := range SamplePools {
		if pool.Id == polID {
			return &pool, nil
		}
	}

	return nil, errors.New("Can't find this pool resource!")
}

// ListPools
func (fc *FakeDbClient) ListPoolsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.StoragePoolSpec, error) {
	var pols []*model.StoragePoolSpec

	for i := range SamplePools {
		pols = append(pols, &SamplePools[i])
	}
	return pols, nil
}
func (fc *FakeDbClient) ListPools(ctx *c.Context) ([]*model.StoragePoolSpec, error) {
	var pols []*model.StoragePoolSpec

	for i := range SamplePools {
		pols = append(pols, &SamplePools[i])
	}
	return pols, nil
}

// UpdatePool
func (fc *FakeDbClient) UpdatePool(ctx *c.Context, polID, name, desp string, usedCapacity int64, used bool) (*model.StoragePoolSpec, error) {
	return nil, nil
}

// DeletePool
func (fc *FakeDbClient) DeletePool(ctx *c.Context, polID string) error {
	return nil
}

// CreateProfile
func (fc *FakeDbClient) CreateProfile(ctx *c.Context, prf *model.ProfileSpec) (*model.ProfileSpec, error) {
	return &SampleProfiles[0], nil
}

// GetProfile
func (fc *FakeDbClient) GetProfile(ctx *c.Context, prfID string) (*model.ProfileSpec, error) {
	for _, profile := range SampleProfiles {
		if profile.Id == prfID {
			return &profile, nil
		}
	}

	return nil, errors.New("Can't find this profile resource!")
}

// GetDefaultProfile
func (fc *FakeDbClient) GetDefaultProfile(ctx *c.Context) (*model.ProfileSpec, error) {
	for _, profile := range SampleProfiles {
		if profile.Name == "default" {
			return &profile, nil
		}
	}

	return nil, errors.New("Can't find default profile resource!")
}

// ListProfiles
func (fc *FakeDbClient) ListProfilesWithFilter(ctx *c.Context, m map[string][]string) ([]*model.ProfileSpec, error) {
	var prfs []*model.ProfileSpec

	for i := range SampleProfiles {
		prfs = append(prfs, &SampleProfiles[i])
	}
	return prfs, nil
}
func (fc *FakeDbClient) ListProfiles(ctx *c.Context) ([]*model.ProfileSpec, error) {
	var prfs []*model.ProfileSpec

	for i := range SampleProfiles {
		prfs = append(prfs, &SampleProfiles[i])
	}
	return prfs, nil
}

// UpdateProfile
func (fc *FakeDbClient) UpdateProfile(ctx *c.Context, prfID string, input *model.ProfileSpec) (*model.ProfileSpec, error) {
	return nil, nil
}

// DeleteProfile
func (fc *FakeDbClient) DeleteProfile(ctx *c.Context, prfID string) error {
	return nil
}

// AddExtraProperty
func (fc *FakeDbClient) AddExtraProperty(ctx *c.Context, prfID string, ext model.ExtraSpec) (*model.ExtraSpec, error) {
	extra := SampleProfiles[0].Extras
	return &extra, nil
}

// ListExtraProperties
func (fc *FakeDbClient) ListExtraProperties(ctx *c.Context, prfID string) (*model.ExtraSpec, error) {
	extra := SampleProfiles[0].Extras
	return &extra, nil
}

// RemoveExtraProperty
func (fc *FakeDbClient) RemoveExtraProperty(ctx *c.Context, prfID, extraKey string) error {
	return nil
}

// CreateVolume
func (fc *FakeDbClient) CreateVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	return vol, nil
}

// GetVolume
func (fc *FakeDbClient) GetVolume(ctx *c.Context, volID string) (*model.VolumeSpec, error) {
	vol := SampleVolumes[0]
	return &vol, nil
}

// ListVolumes
func (fc *FakeDbClient) ListVolumesWithFilter(ctx *c.Context, m map[string][]string) ([]*model.VolumeSpec, error) {
	var vols []*model.VolumeSpec

	for i := range SampleVolumes {
		vols = append(vols, &SampleVolumes[i])
	}
	return vols, nil
}
func (fc *FakeDbClient) ListVolumes(ctx *c.Context) ([]*model.VolumeSpec, error) {
	var vols []*model.VolumeSpec

	for i := range SampleVolumes {
		vols = append(vols, &SampleVolumes[i])
	}
	return vols, nil
}

// UpdateVolume
func (fc *FakeDbClient) UpdateVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	return &SampleVolumes[0], nil
}

// DeleteVolume
func (fc *FakeDbClient) DeleteVolume(ctx *c.Context, volID string) error {
	return nil
}

// ExtendVolume ...
func (fc *FakeDbClient) ExtendVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	return &SampleVolumes[0], nil
}

// CreateVolumeAttachment
func (fc *FakeDbClient) CreateVolumeAttachment(ctx *c.Context, attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	return attachment, nil
}

// GetVolumeAttachment
func (fc *FakeDbClient) GetVolumeAttachment(ctx *c.Context, attachmentId string) (*model.VolumeAttachmentSpec, error) {
	attach := SampleAttachments[0]
	return &attach, nil
}

// ListVolumeAttachments
func (fc *FakeDbClient) ListVolumeAttachmentsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.VolumeAttachmentSpec, error) {
	var atcs []*model.VolumeAttachmentSpec

	for i := range SampleAttachments {
		atcs = append(atcs, &SampleAttachments[i])
	}
	return atcs, nil
}
func (fc *FakeDbClient) ListVolumeAttachments(ctx *c.Context, volumeId string) ([]*model.VolumeAttachmentSpec, error) {
	var atcs []*model.VolumeAttachmentSpec

	for i := range SampleAttachments {
		atcs = append(atcs, &SampleAttachments[i])
	}
	return atcs, nil
}

// UpdateVolumeAttachment
func (fc *FakeDbClient) UpdateVolumeAttachment(ctx *c.Context, attachmentId string, attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	return nil, nil
}

// DeleteVolumeAttachment
func (fc *FakeDbClient) DeleteVolumeAttachment(ctx *c.Context, attachmentId string) error {
	SampleAttachments = []model.VolumeAttachmentSpec{}
	return nil
}

// CreateVolumeSnapshot
func (fc *FakeDbClient) CreateVolumeSnapshot(ctx *c.Context, vs *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	return vs, nil
}

// GetVolumeSnapshot
func (fc *FakeDbClient) GetVolumeSnapshot(ctx *c.Context, snapshotID string) (*model.VolumeSnapshotSpec, error) {
	snap := SampleSnapshots[0]
	return &snap, nil
}

// ListVolumeSnapshots
func (fc *FakeDbClient) ListVolumeSnapshotsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.VolumeSnapshotSpec, error) {
	var snps []*model.VolumeSnapshotSpec

	for i := range SampleSnapshots {
		snps = append(snps, &SampleSnapshots[i])
	}
	return snps, nil
}
func (fc *FakeDbClient) ListVolumeSnapshots(ctx *c.Context) ([]*model.VolumeSnapshotSpec, error) {
	var snps []*model.VolumeSnapshotSpec

	for i := range SampleSnapshots {
		snps = append(snps, &SampleSnapshots[i])
	}
	return snps, nil
}

// UpdateVolumeSnapshot
func (fc *FakeDbClient) UpdateVolumeSnapshot(ctx *c.Context, snapshotID string, vs *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	return &SampleSnapshots[0], nil
}

// DeleteVolumeSnapshot
func (fc *FakeDbClient) DeleteVolumeSnapshot(ctx *c.Context, snapshotID string) error {
	return nil
}

func (fc *FakeDbClient) CreateReplication(ctx *c.Context, replication *model.ReplicationSpec) (*model.ReplicationSpec, error) {
	return &SampleReplications[0], nil
}

func (fc *FakeDbClient) GetReplication(ctx *c.Context, replicationId string) (*model.ReplicationSpec, error) {
	return &SampleReplications[0], nil
}

func (fc *FakeDbClient) ListReplication(ctx *c.Context) ([]*model.ReplicationSpec, error) {
	var replications = []*model.ReplicationSpec{
		&SampleReplications[0], &SampleReplications[1],
	}
	return replications, nil
}

func (fc *FakeDbClient) GetReplicationByVolumeId(ctx *c.Context, volumeId string) (*model.ReplicationSpec, error) {
	return &SampleReplications[0], nil
}

func (fc *FakeDbClient) ListReplicationWithFilter(ctx *c.Context, m map[string][]string) ([]*model.ReplicationSpec, error) {
	var replications = []*model.ReplicationSpec{
		&SampleReplications[0], &SampleReplications[1],
	}
	return replications, nil
}

func (fc *FakeDbClient) DeleteReplication(ctx *c.Context, replicationId string) error {
	return nil
}

func (fc *FakeDbClient) UpdateReplication(ctx *c.Context, replicationId string, input *model.ReplicationSpec) (*model.ReplicationSpec, error) {
	return nil, nil
}

// CreateVolumeGroup
func (fc *FakeDbClient) CreateVolumeGroup(ctx *c.Context, vg *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) UpdateVolumeGroup(ctx *c.Context, vg *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) GetVolumeGroup(ctx *c.Context, vgId string) (*model.VolumeGroupSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) UpdateStatus(ctx *c.Context, in interface{}, status string) error {
	return nil
}

func (fc *FakeDbClient) ListVolumesByIds(ctx *c.Context, ids []string) ([]*model.VolumeSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) ListVolumesByGroupId(ctx *c.Context, vgId string) ([]*model.VolumeSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) ListSnapshotsByVolumeId(ctx *c.Context, volumeId string) ([]*model.VolumeSnapshotSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) DeleteVolumeGroup(ctx *c.Context, volumeId string) error {
	return nil
}

func (fc *FakeDbClient) ListVolumeGroupsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.VolumeGroupSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) ListVolumeGroups(ctx *c.Context) ([]*model.VolumeGroupSpec, error) {
	return nil, nil
}

func (fc *FakeDbClient) VolumesToUpdate(ctx *c.Context, volumeList []*model.VolumeSpec) ([]*model.VolumeSpec, error) {
	return nil, nil
}
