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
func (v *VolumeMgr) ListVolumes(args ...interface{}) ([]*model.VolumeSpec, error) {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeURL(urls.Client, v.TenantId)}, "/")

	param, err := processListParam(args)
	if err != nil {
		return nil, err
	}

	if param != "" {
		url += "?" + param
	}

	var res []*model.VolumeSpec
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
		urls.GenerateVolumeURL(urls.Client, v.TenantId, volID, "resize")}, "/")

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
func (v *VolumeMgr) ListVolumeAttachments(args ...interface{}) ([]*model.VolumeAttachmentSpec, error) {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateAttachmentURL(urls.Client, v.TenantId)}, "/")

	param, err := processListParam(args)
	if err != nil {
		return nil, err
	}

	if param != "" {
		url += "?" + param
	}
	var res []*model.VolumeAttachmentSpec
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
func (v *VolumeMgr) ListVolumeSnapshots(args ...interface{}) ([]*model.VolumeSnapshotSpec, error) {
	var res []*model.VolumeSnapshotSpec

	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateSnapshotURL(urls.Client, v.TenantId)}, "/")

	param, err := processListParam(args)
	if err != nil {
		return nil, err
	}

	if param != "" {
		url += "?" + param
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
func (v *VolumeMgr) ListVolumeGroups(args ...interface{}) ([]*model.VolumeGroupSpec, error) {
	var res []*model.VolumeGroupSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeGroupURL(urls.Client, v.TenantId)}, "/")

	param, err := processListParam(args)
	if err != nil {
		return nil, err
	}

	if param != "" {
		url += "?" + param
	}

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
