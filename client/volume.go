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

package client

import (
	"strconv"
	"strings"

	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/urls"
)

// VolumeBuilder contains request body of handling a volume request.
// Currently it's assigned as the pointer of VolumeSpec struct, but it
// could be discussed if it's better to define an interface.
type VolumeBuilder *model.VolumeSpec

// ExtendVolumeBuilder contains request body of handling a extend volume request.
// Currently it's assigned as the pointer of ExtendVolumeSpec struct, but it
// could be discussed if it's better to define an interface.
type ExtendVolumeBuilder *model.ExtendVolumeSpec

// VolumeAttachmentBuilder contains request body of handling a volume request.
// Currently it's assigned as the pointer of VolumeSpec struct, but it
// could be discussed if it's better to define an interface.
type VolumeAttachmentBuilder *model.VolumeAttachmentSpec

// VolumeSnapshotBuilder contains request body of handling a volume snapshot
// request. Currently it's assigned as the pointer of VolumeSnapshotSpec
// struct, but it could be discussed if it's better to define an interface.
type VolumeSnapshotBuilder *model.VolumeSnapshotSpec

// VolumeGroupBuilder contains request body of handling a volume group
// request. Currently it's assigned as the pointer of VolumeGroupSpec
// struct, but it could be discussed if it's better to define an interface.
type VolumeGroupBuilder *model.VolumeGroupSpec

// NewVolumeMgr
func NewVolumeMgr(r Receiver, edp string, tenantId string) *VolumeMgr {
	return &VolumeMgr{
		Receiver: r,
		Endpoint: edp,
		TenantId: tenantId,
	}
}

// VolumeMgr
type VolumeMgr struct {
	Receiver
	Endpoint string
	TenantId string
}

// CreateVolume
func (v *VolumeMgr) CreateVolume(body VolumeBuilder) (*model.VolumeSpec, error) {
	var res model.VolumeSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeURL(urls.Client, v.TenantId)}, "/")

	if err := v.Recv(url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetVolume
func (v *VolumeMgr) GetVolume(volID string) (*model.VolumeSpec, error) {
	var res model.VolumeSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeURL(urls.Client, v.TenantId, volID)}, "/")

	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListVolumes
func (v *VolumeMgr) ListVolumes(p []string, vol *model.VolumeSpec) ([]*model.VolumeSpec, error) {
	var res []*model.VolumeSpec
	var u string

	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeURL(urls.Client, v.TenantId)}, "/")

	var limit, offset, sortDir, sortKey, id, createdAt, updatedAt, name, description, userId, availabilityZone, size, status, poolId, profileId string
	var urlpara []string
	if p[0] != "" {
		limit = "limit=" + p[0]
		urlpara = append(urlpara, limit)
	}
	if p[1] != "" {
		offset = "offset=" + p[1]
		urlpara = append(urlpara, offset)
	}
	if p[2] != "" {
		sortDir = "sortDir=" + p[2]
		urlpara = append(urlpara, sortDir)
	}
	if p[3] != "" {
		sortKey = "sortKey=" + p[3]
		urlpara = append(urlpara, sortKey)
	}
	if vol.Id != "" {
		id = "Id=" + vol.Id
		urlpara = append(urlpara, id)
	}
	if vol.CreatedAt != "" {
		createdAt = "CreatedAt=" + vol.CreatedAt
		urlpara = append(urlpara, createdAt)
	}
	if vol.UpdatedAt != "" {
		updatedAt = "UpdatedAt=" + vol.UpdatedAt
		urlpara = append(urlpara, updatedAt)
	}
	if vol.Name != "" {
		name = "Name=" + vol.Name
		urlpara = append(urlpara, name)
	}
	if vol.Description != "" {
		description = "Description=" + vol.Description
		urlpara = append(urlpara, description)
	}
	if vol.UserId != "" {
		userId = "UserId=" + vol.UserId
		urlpara = append(urlpara, userId)
	}
	if vol.AvailabilityZone != "" {
		availabilityZone = "AvailabilityZone=" + vol.AvailabilityZone
		urlpara = append(urlpara, availabilityZone)
	}
	if vol.Size != 0 {
		size = "Size=" + strconv.FormatInt(vol.Size, 10)
		urlpara = append(urlpara, size)
	}
	if vol.Status != "" {
		status = "Status=" + vol.Status
		urlpara = append(urlpara, status)
	}
	if vol.PoolId != "" {
		poolId = "PoolId=" + vol.PoolId
		urlpara = append(urlpara, poolId)
	}
	if vol.ProfileId != "" {
		profileId = "ProfileId=" + vol.ProfileId
		urlpara = append(urlpara, profileId)
	}
	if len(urlpara) > 0 {
		u = strings.Join(urlpara, "&")
		url += "?" + u
	}

	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteVolume
func (v *VolumeMgr) DeleteVolume(volID string, body VolumeBuilder) error {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeURL(urls.Client, v.TenantId, volID)}, "/")

	return v.Recv(url, "DELETE", body, nil)
}

// UpdateVolume
func (v *VolumeMgr) UpdateVolume(volID string, body VolumeBuilder) (*model.VolumeSpec, error) {
	var res model.VolumeSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeURL(urls.Client, v.TenantId, volID)}, "/")

	if err := v.Recv(url, "PUT", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ExtendVolume ...
func (v *VolumeMgr) ExtendVolume(volID string, body ExtendVolumeBuilder) (*model.VolumeSpec, error) {
	var res model.VolumeSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateNewVolumeURL(urls.Client, v.TenantId, volID, "resize")}, "/")

	if err := v.Recv(url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// CreateVolumeAttachment
func (v *VolumeMgr) CreateVolumeAttachment(body VolumeAttachmentBuilder) (*model.VolumeAttachmentSpec, error) {
	var res model.VolumeAttachmentSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateAttachmentURL(urls.Client, v.TenantId)}, "/")

	if err := v.Recv(url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// UpdateVolumeAttachment
func (v *VolumeMgr) UpdateVolumeAttachment(atcID string, body VolumeAttachmentBuilder) (*model.VolumeAttachmentSpec, error) {
	var res model.VolumeAttachmentSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateAttachmentURL(urls.Client, v.TenantId, atcID)}, "/")

	if err := v.Recv(url, "PUT", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetVolumeAttachment
func (v *VolumeMgr) GetVolumeAttachment(atcID string) (*model.VolumeAttachmentSpec, error) {
	var res model.VolumeAttachmentSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateAttachmentURL(urls.Client, v.TenantId, atcID)}, "/")

	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListVolumeAttachments
func (v *VolumeMgr) ListVolumeAttachments(p []string, volAtm *model.VolumeAttachmentSpec) ([]*model.VolumeAttachmentSpec, error) {
	var res []*model.VolumeAttachmentSpec
	var u string

	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateAttachmentURL(urls.Client, v.TenantId)}, "/")

	var limit, offset, sortDir, sortKey, id, createdAt, updatedAt, userId, tenantId, volumeId, status, mountpoint string
	var urlpara []string

	if p[0] != "" {
		limit = "limit=" + p[0]
		urlpara = append(urlpara, limit)
	}
	if p[1] != "" {
		offset = "offset=" + p[1]
		urlpara = append(urlpara, offset)
	}
	if p[2] != "" {
		sortDir = "sortDir=" + p[2]
		urlpara = append(urlpara, sortDir)
	}
	if p[3] != "" {
		sortKey = "sortKey=" + p[3]
		urlpara = append(urlpara, sortKey)
	}
	if volAtm.Id != "" {
		id = "Id" + volAtm.Id
		urlpara = append(urlpara, id)
	}
	if volAtm.CreatedAt != "" {
		createdAt = "CreatedAt=" + volAtm.CreatedAt
		urlpara = append(urlpara, createdAt)
	}
	if volAtm.UpdatedAt != "" {
		updatedAt = "UpdatedAt=" + volAtm.UpdatedAt
		urlpara = append(urlpara, updatedAt)
	}
	if volAtm.UserId != "" {
		userId = "UserId=" + volAtm.UserId
		urlpara = append(urlpara, userId)
	}
	if volAtm.TenantId != "" {
		tenantId = "TenantId=" + volAtm.TenantId
		urlpara = append(urlpara, tenantId)
	}
	if volAtm.VolumeId != "" {
		volumeId = "VolumeId=" + volAtm.VolumeId
		urlpara = append(urlpara, volumeId)
	}
	if volAtm.Status != "" {
		status = "Status=" + volAtm.Status
		urlpara = append(urlpara, status)
	}
	if volAtm.Mountpoint != "" {
		mountpoint = "Mountpoint=" + volAtm.Mountpoint
		urlpara = append(urlpara, mountpoint)
	}
	if len(urlpara) > 0 {
		u = strings.Join(urlpara, "&")
		url += "?" + u
	}

	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteVolumeAttachment
func (v *VolumeMgr) DeleteVolumeAttachment(atcID string, body VolumeAttachmentBuilder) error {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateAttachmentURL(urls.Client, v.TenantId, atcID)}, "/")

	return v.Recv(url, "DELETE", body, nil)
}

// CreateVolumeSnapshot
func (v *VolumeMgr) CreateVolumeSnapshot(body VolumeSnapshotBuilder) (*model.VolumeSnapshotSpec, error) {
	var res model.VolumeSnapshotSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateSnapshotURL(urls.Client, v.TenantId)}, "/")

	if err := v.Recv(url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetVolumeSnapshot
func (v *VolumeMgr) GetVolumeSnapshot(snpID string) (*model.VolumeSnapshotSpec, error) {
	var res model.VolumeSnapshotSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateSnapshotURL(urls.Client, v.TenantId, snpID)}, "/")

	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListVolumeSnapshots
func (v *VolumeMgr) ListVolumeSnapshots(p []string, volSnap *model.VolumeSnapshotSpec) ([]*model.VolumeSnapshotSpec, error) {
	var res []*model.VolumeSnapshotSpec
	var u string

	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateSnapshotURL(urls.Client, v.TenantId)}, "/")

	var limit, offset, sortDir, sortKey, id, createdAt string
	var updatedAt, userId, tenantId, name, description, status, volumeId, size string
	var urlpara []string

	if p[0] != "" {
		limit = "limit=" + p[0]
		urlpara = append(urlpara, limit)
	}
	if p[1] != "" {
		offset = "offset=" + p[1]
		urlpara = append(urlpara, offset)
	}
	if p[2] != "" {
		sortDir = "sortDir=" + p[2]
		urlpara = append(urlpara, sortDir)
	}
	if p[3] != "" {
		sortKey = "sortKey=" + p[3]
		urlpara = append(urlpara, sortKey)
	}
	if volSnap.Id != "" {
		id = "Id=" + volSnap.Id
		urlpara = append(urlpara, id)
	}
	if volSnap.CreatedAt != "" {
		createdAt = "CreatedAt=" + volSnap.CreatedAt
		urlpara = append(urlpara, createdAt)
	}
	if volSnap.UpdatedAt != "" {
		updatedAt = "UpdatedAt=" + volSnap.UpdatedAt
		urlpara = append(urlpara, updatedAt)
	}

	if volSnap.UserId != "" {
		userId = "UserId=" + volSnap.UserId
		urlpara = append(urlpara, userId)
	}
	if volSnap.TenantId != "" {
		tenantId = "TenantId=" + volSnap.TenantId
		urlpara = append(urlpara, tenantId)
	}
	if volSnap.Name != "" {
		name = "Name=" + volSnap.Name
		urlpara = append(urlpara, name)
	}
	if volSnap.Description != "" {
		description = "Description=" + volSnap.Description
		urlpara = append(urlpara, description)
	}
	if volSnap.Status != "" {
		status = "Status=" + volSnap.Status
		urlpara = append(urlpara, status)
	}
	if volSnap.Size != 0 {
		size = "Size=" + strconv.FormatInt(volSnap.Size, 10)
		urlpara = append(urlpara, size)
	}
	if volSnap.VolumeId != "" {
		volumeId = "VolumeId=" + volSnap.VolumeId
		urlpara = append(urlpara, volumeId)
	}
	if len(urlpara) > 0 {
		u = strings.Join(urlpara, "&")
		url += "?" + u
	}

	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteVolumeSnapshot
func (v *VolumeMgr) DeleteVolumeSnapshot(snpID string, body VolumeSnapshotBuilder) error {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateSnapshotURL(urls.Client, v.TenantId, snpID)}, "/")

	return v.Recv(url, "DELETE", body, nil)
}

// UpdateVolumeSnapshot
func (v *VolumeMgr) UpdateVolumeSnapshot(snpID string, body VolumeSnapshotBuilder) (*model.VolumeSnapshotSpec, error) {
	var res model.VolumeSnapshotSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateSnapshotURL(urls.Client, v.TenantId, snpID)}, "/")

	if err := v.Recv(url, "PUT", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// CreateVolumeGroup
func (v *VolumeMgr) CreateVolumeGroup(body VolumeGroupBuilder) (*model.VolumeGroupSpec, error) {
	var res model.VolumeGroupSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeGroupURL(urls.Client, v.TenantId)}, "/")

	if err := v.Recv(url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetVolumeGroup
func (v *VolumeMgr) GetVolumeGroup(vgId string) (*model.VolumeGroupSpec, error) {
	var res model.VolumeGroupSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeGroupURL(urls.Client, v.TenantId, vgId)}, "/")

	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListVolumeGroups
func (v *VolumeMgr) ListVolumeGroups() ([]*model.VolumeGroupSpec, error) {
	var res []*model.VolumeGroupSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeGroupURL(urls.Client, v.TenantId)}, "/")

	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteVolumeGroup
func (v *VolumeMgr) DeleteVolumeGroup(vgId string, body VolumeGroupBuilder) error {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeGroupURL(urls.Client, v.TenantId, vgId)}, "/")

	return v.Recv(url, "DELETE", body, nil)
}

// UpdateVolumeSnapshot
func (v *VolumeMgr) UpdateVolumeGroup(vgId string, body VolumeGroupBuilder) (*model.VolumeGroupSpec, error) {
	var res model.VolumeGroupSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeGroupURL(urls.Client, v.TenantId, vgId)}, "/")

	if err := v.Recv(url, "PUT", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
