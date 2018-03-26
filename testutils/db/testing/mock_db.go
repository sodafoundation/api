// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package testing

import (
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/model"
	"github.com/stretchr/testify/mock"
)

// MockClient is a fake database client for testing.
type MockClient struct {
	mock.Mock
}

// AddExtraProperty
func (_m *MockClient) AddExtraProperty(ctx *c.Context, prfID string, ext model.ExtraSpec) (*model.ExtraSpec, error) {
	ret := _m.Called(prfID, ext)

	var r0 *model.ExtraSpec
	if rf, ok := ret.Get(0).(func(string, model.ExtraSpec) *model.ExtraSpec); ok {
		r0 = rf(prfID, ext)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ExtraSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, model.ExtraSpec) error); ok {
		r1 = rf(prfID, ext)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateDock
func (_m *MockClient) CreateDock(ctx *c.Context, dck *model.DockSpec) (*model.DockSpec, error) {
	ret := _m.Called(dck)

	var r0 *model.DockSpec
	if rf, ok := ret.Get(0).(func(*model.DockSpec) *model.DockSpec); ok {
		r0 = rf(dck)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DockSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.DockSpec) error); ok {
		r1 = rf(dck)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreatePool
func (_m *MockClient) CreatePool(ctx *c.Context, pol *model.StoragePoolSpec) (*model.StoragePoolSpec, error) {
	ret := _m.Called(pol)

	var r0 *model.StoragePoolSpec
	if rf, ok := ret.Get(0).(func(*model.StoragePoolSpec) *model.StoragePoolSpec); ok {
		r0 = rf(pol)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.StoragePoolSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.StoragePoolSpec) error); ok {
		r1 = rf(pol)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateProfile
func (_m *MockClient) CreateProfile(ctx *c.Context, prf *model.ProfileSpec) (*model.ProfileSpec, error) {
	ret := _m.Called(prf)

	var r0 *model.ProfileSpec
	if rf, ok := ret.Get(0).(func(*model.ProfileSpec) *model.ProfileSpec); ok {
		r0 = rf(prf)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ProfileSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.ProfileSpec) error); ok {
		r1 = rf(prf)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateVolume
func (_m *MockClient) CreateVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	ret := _m.Called(vol)

	var r0 *model.VolumeSpec
	if rf, ok := ret.Get(0).(func(*model.VolumeSpec) *model.VolumeSpec); ok {
		r0 = rf(vol)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.VolumeSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.VolumeSpec) error); ok {
		r1 = rf(vol)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateVolumeAttachment
func (_m *MockClient) CreateVolumeAttachment(ctx *c.Context, attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	ret := _m.Called(attachment)

	var r0 *model.VolumeAttachmentSpec
	if rf, ok := ret.Get(0).(func(*model.VolumeAttachmentSpec) *model.VolumeAttachmentSpec); ok {
		r0 = rf(attachment)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.VolumeAttachmentSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.VolumeAttachmentSpec) error); ok {
		r1 = rf(attachment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateVolumeSnapshot
func (_m *MockClient) CreateVolumeSnapshot(ctx *c.Context, vs *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	ret := _m.Called(vs)

	var r0 *model.VolumeSnapshotSpec
	if rf, ok := ret.Get(0).(func(*model.VolumeSnapshotSpec) *model.VolumeSnapshotSpec); ok {
		r0 = rf(vs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.VolumeSnapshotSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.VolumeSnapshotSpec) error); ok {
		r1 = rf(vs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteDock
func (_m *MockClient) DeleteDock(ctx *c.Context, dckID string) error {
	ret := _m.Called(dckID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(dckID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePool
func (_m *MockClient) DeletePool(ctx *c.Context, polID string) error {
	ret := _m.Called(polID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(polID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteProfile
func (_m *MockClient) DeleteProfile(ctx *c.Context, prfID string) error {
	ret := _m.Called(prfID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(prfID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteVolume
func (_m *MockClient) DeleteVolume(ctx *c.Context, volID string) error {
	ret := _m.Called(volID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(volID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteVolumeAttachment
func (_m *MockClient) DeleteVolumeAttachment(ctx *c.Context, attachmentId string) error {
	ret := _m.Called(attachmentId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(attachmentId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteVolumeSnapshot
func (_m *MockClient) DeleteVolumeSnapshot(ctx *c.Context, snapshotID string) error {
	ret := _m.Called(snapshotID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(snapshotID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDock
func (_m *MockClient) GetDock(ctx *c.Context, dckID string) (*model.DockSpec, error) {
	ret := _m.Called(dckID)

	var r0 *model.DockSpec
	if rf, ok := ret.Get(0).(func(string) *model.DockSpec); ok {
		r0 = rf(dckID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DockSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(dckID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDockByPoolId
func (_m *MockClient) GetDockByPoolId(ctx *c.Context, poolId string) (*model.DockSpec, error) {
	ret := _m.Called(poolId)

	var r0 *model.DockSpec
	if rf, ok := ret.Get(0).(func(string) *model.DockSpec); ok {
		r0 = rf(poolId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DockSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(poolId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPool
func (_m *MockClient) GetPool(ctx *c.Context, polID string) (*model.StoragePoolSpec, error) {
	ret := _m.Called(polID)

	var r0 *model.StoragePoolSpec
	if rf, ok := ret.Get(0).(func(string) *model.StoragePoolSpec); ok {
		r0 = rf(polID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.StoragePoolSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(polID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProfile
func (_m *MockClient) GetProfile(ctx *c.Context, prfID string) (*model.ProfileSpec, error) {
	ret := _m.Called(prfID)

	var r0 *model.ProfileSpec
	if rf, ok := ret.Get(0).(func(string) *model.ProfileSpec); ok {
		r0 = rf(prfID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ProfileSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(prfID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDefaultProfile
func (_m *MockClient) GetDefaultProfile(ctx *c.Context) (*model.ProfileSpec, error) {
	ret := _m.Called()

	var r0 *model.ProfileSpec
	if rf, ok := ret.Get(0).(func() *model.ProfileSpec); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ProfileSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVolume
func (_m *MockClient) GetVolume(ctx *c.Context, volID string) (*model.VolumeSpec, error) {
	ret := _m.Called(volID)

	var r0 *model.VolumeSpec
	if rf, ok := ret.Get(0).(func(string) *model.VolumeSpec); ok {
		r0 = rf(volID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.VolumeSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(volID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVolumeAttachment
func (_m *MockClient) GetVolumeAttachment(ctx *c.Context, attachmentId string) (*model.VolumeAttachmentSpec, error) {
	ret := _m.Called(attachmentId)

	var r0 *model.VolumeAttachmentSpec
	if rf, ok := ret.Get(0).(func(string) *model.VolumeAttachmentSpec); ok {
		r0 = rf(attachmentId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.VolumeAttachmentSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(attachmentId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVolumeSnapshot
func (_m *MockClient) GetVolumeSnapshot(ctx *c.Context, snapshotID string) (*model.VolumeSnapshotSpec, error) {
	ret := _m.Called(snapshotID)

	var r0 *model.VolumeSnapshotSpec
	if rf, ok := ret.Get(0).(func(string) *model.VolumeSnapshotSpec); ok {
		r0 = rf(snapshotID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.VolumeSnapshotSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(snapshotID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDocks
func (_m *MockClient) ListDocks(ctx *c.Context) ([]*model.DockSpec, error) {
	ret := _m.Called()

	var r0 []*model.DockSpec
	if rf, ok := ret.Get(0).(func() []*model.DockSpec); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.DockSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDocks
func (_m *MockClient) ListDocksWithFilter(ctx *c.Context, m map[string][]string) ([]*model.DockSpec, error) {
	ret := _m.Called(m)
	var r0 []*model.DockSpec
	if rf, ok := ret.Get(0).(func() []*model.DockSpec); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.DockSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListExtraProperties
func (_m *MockClient) ListExtraProperties(ctx *c.Context, prfID string) (*model.ExtraSpec, error) {
	ret := _m.Called(prfID)

	var r0 *model.ExtraSpec
	if rf, ok := ret.Get(0).(func(string) *model.ExtraSpec); ok {
		r0 = rf(prfID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ExtraSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(prfID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPools
func (_m *MockClient) ListPools(ctx *c.Context) ([]*model.StoragePoolSpec, error) {
	ret := _m.Called()

	var r0 []*model.StoragePoolSpec
	if rf, ok := ret.Get(0).(func() []*model.StoragePoolSpec); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.StoragePoolSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPools
func (_m *MockClient) ListPoolsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.StoragePoolSpec, error) {
	ret := _m.Called(m)

	var r0 []*model.StoragePoolSpec
	if rf, ok := ret.Get(0).(func() []*model.StoragePoolSpec); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.StoragePoolSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListProfiles
func (_m *MockClient) ListProfiles(ctx *c.Context) ([]*model.ProfileSpec, error) {
	ret := _m.Called()

	var r0 []*model.ProfileSpec
	if rf, ok := ret.Get(0).(func() []*model.ProfileSpec); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ProfileSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockClient) ListProfilesWithFilter(ctx *c.Context, m map[string][]string) ([]*model.ProfileSpec, error) {
	ret := _m.Called(m)

	var r0 []*model.ProfileSpec
	if rf, ok := ret.Get(0).(func() []*model.ProfileSpec); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ProfileSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListVolumeAttachments
func (_m *MockClient) ListVolumeAttachments(ctx *c.Context, volumeId string) ([]*model.VolumeAttachmentSpec, error) {
	ret := _m.Called(volumeId)

	var r0 []*model.VolumeAttachmentSpec
	if rf, ok := ret.Get(0).(func(string) []*model.VolumeAttachmentSpec); ok {
		r0 = rf(volumeId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.VolumeAttachmentSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(volumeId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockClient) ListVolumeAttachmentsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.VolumeAttachmentSpec, error) {
	ret := _m.Called(m)

	var r0 []*model.VolumeAttachmentSpec
	if rf, ok := ret.Get(0).(func(string) []*model.VolumeAttachmentSpec); ok {
		r0 = rf(m["volumeId"][0])
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.VolumeAttachmentSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(m["volumeId"][0])
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListVolumeSnapshots
func (_m *MockClient) ListVolumeSnapshots(ctx *c.Context) ([]*model.VolumeSnapshotSpec, error) {
	ret := _m.Called()

	var r0 []*model.VolumeSnapshotSpec
	if rf, ok := ret.Get(0).(func() []*model.VolumeSnapshotSpec); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.VolumeSnapshotSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockClient) ListVolumeSnapshotsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.VolumeSnapshotSpec, error) {
	ret := _m.Called(m)

	var r0 []*model.VolumeSnapshotSpec
	if rf, ok := ret.Get(0).(func() []*model.VolumeSnapshotSpec); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.VolumeSnapshotSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListVolumes
func (_m *MockClient) ListVolumes(ctx *c.Context) ([]*model.VolumeSpec, error) {
	ret := _m.Called()

	var r0 []*model.VolumeSpec
	if rf, ok := ret.Get(0).(func() []*model.VolumeSpec); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.VolumeSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockClient) ListVolumesWithFilter(ctx *c.Context, m map[string][]string) ([]*model.VolumeSpec, error) {
	ret := _m.Called(m)

	var r0 []*model.VolumeSpec
	if rf, ok := ret.Get(0).(func() []*model.VolumeSpec); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.VolumeSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveExtraProperty
func (_m *MockClient) RemoveExtraProperty(ctx *c.Context, prfID string, extraKey string) error {
	ret := _m.Called(prfID, extraKey)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(prfID, extraKey)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateDock
func (_m *MockClient) UpdateDock(ctx *c.Context, dckID string, name string, desp string) (*model.DockSpec, error) {
	ret := _m.Called(dckID, name, desp)

	var r0 *model.DockSpec
	if rf, ok := ret.Get(0).(func(string, string, string) *model.DockSpec); ok {
		r0 = rf(dckID, name, desp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DockSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(dckID, name, desp)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePool
func (_m *MockClient) UpdatePool(ctx *c.Context, polID string, name string, desp string, usedCapacity int64, used bool) (*model.StoragePoolSpec, error) {
	ret := _m.Called(polID, name, desp, usedCapacity, used)

	var r0 *model.StoragePoolSpec
	if rf, ok := ret.Get(0).(func(string, string, string, int64, bool) *model.StoragePoolSpec); ok {
		r0 = rf(polID, name, desp, usedCapacity, used)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.StoragePoolSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, int64, bool) error); ok {
		r1 = rf(polID, name, desp, usedCapacity, used)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProfile
func (_m *MockClient) UpdateProfile(ctx *c.Context, prfID string, input *model.ProfileSpec) (*model.ProfileSpec, error) {
	ret := _m.Called(prfID, input)

	var r0 *model.ProfileSpec
	if rf, ok := ret.Get(0).(func(string, *model.ProfileSpec) *model.ProfileSpec); ok {
		r0 = rf(prfID, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ProfileSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *model.ProfileSpec) error); ok {
		r1 = rf(prfID, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateVolume ...
func (_m *MockClient) UpdateVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	ret := _m.Called(vol.Id, vol)

	var r0 *model.VolumeSpec
	if rf, ok := ret.Get(0).(func(string, *model.VolumeSpec) *model.VolumeSpec); ok {
		r0 = rf(vol.Id, vol)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.VolumeSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *model.VolumeSpec) error); ok {
		r1 = rf(vol.Id, vol)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateVolumeAttachment ...
func (_m *MockClient) UpdateVolumeAttachment(ctx *c.Context, attachmentID string, attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	ret := _m.Called(attachmentID, attachment)

	var r0 *model.VolumeAttachmentSpec
	if rf, ok := ret.Get(0).(func(string, *model.VolumeAttachmentSpec) *model.VolumeAttachmentSpec); ok {
		r0 = rf(attachmentID, attachment)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.VolumeAttachmentSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *model.VolumeAttachmentSpec) error); ok {
		r1 = rf(attachmentID, attachment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateVolumeSnapshot ...
func (_m *MockClient) UpdateVolumeSnapshot(ctx *c.Context, snapshotID string, vs *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	ret := _m.Called(snapshotID, vs)

	var r0 *model.VolumeSnapshotSpec
	if rf, ok := ret.Get(0).(func(string, *model.VolumeSnapshotSpec) *model.VolumeSnapshotSpec); ok {
		r0 = rf(snapshotID, vs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.VolumeSnapshotSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *model.VolumeSnapshotSpec) error); ok {
		r1 = rf(snapshotID, vs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExtendVolume ...
func (_m *MockClient) ExtendVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	ret := _m.Called(vol)

	var r0 *model.VolumeSpec
	if rf, ok := ret.Get(0).(func(string, *model.VolumeSpec) *model.VolumeSpec); ok {
		r0 = rf(vol.Id, vol)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.VolumeSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *model.VolumeSpec) error); ok {
		r1 = rf(vol.Id, vol)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
