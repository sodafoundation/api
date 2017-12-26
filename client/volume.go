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

// VolumeAttachmentBuilder contains request body of handling a volume request.
// Currently it's assigned as the pointer of VolumeSpec struct, but it
// could be discussed if it's better to define an interface.
type VolumeAttachmentBuilder *model.VolumeAttachmentSpec

// VolumeSnapshotBuilder contains request body of handling a volume snapshot
// request. Currently it's assigned as the pointer of VolumeSnapshotSpec
// struct, but it could be discussed if it's better to define an interface.
type VolumeSnapshotBuilder *model.VolumeSnapshotSpec

// NewVolumeMgr
func NewVolumeMgr(edp string) *VolumeMgr {
	return &VolumeMgr{
		Receiver: NewReceiver(),
		Endpoint: edp,
	}
}

// VolumeMgr
type VolumeMgr struct {
	Receiver

	Endpoint string
}

// CreateVolume
func (v *VolumeMgr) CreateVolume(body VolumeBuilder) (*model.VolumeSpec, error) {
	var res model.VolumeSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeURL()}, "/")

	if err := v.Recv(request, url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetVolume
func (v *VolumeMgr) GetVolume(volID string) (*model.VolumeSpec, error) {
	var res model.VolumeSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeURL(volID)}, "/")

	if err := v.Recv(request, url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListVolumes
func (v *VolumeMgr) ListVolumes() ([]*model.VolumeSpec, error) {
	var res []*model.VolumeSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeURL()}, "/")

	if err := v.Recv(request, url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteVolume
func (v *VolumeMgr) DeleteVolume(volID string, body VolumeBuilder) error {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeURL(volID)}, "/")

	return v.Recv(request, url, "DELETE", body, nil)
}

// UpdateVolume
func (v *VolumeMgr) UpdateVolume(volID string, body VolumeBuilder) (*model.VolumeSpec, error) {
	var res model.VolumeSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateVolumeURL(volID)}, "/")

	if err := v.Recv(request, url, "PUT", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// CreateVolumeAttachment
func (v *VolumeMgr) CreateVolumeAttachment(body VolumeAttachmentBuilder) (*model.VolumeAttachmentSpec, error) {
	var res model.VolumeAttachmentSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateAttachmentURL()}, "/")

	if err := v.Recv(request, url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// UpdateVolumeAttachment
func (v *VolumeMgr) UpdateVolumeAttachment(atcID string, body VolumeAttachmentBuilder) (*model.VolumeAttachmentSpec, error) {
	var res model.VolumeAttachmentSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateAttachmentURL(atcID)}, "/")

	if err := v.Recv(request, url, "PUT", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetVolumeAttachment
func (v *VolumeMgr) GetVolumeAttachment(atcID string) (*model.VolumeAttachmentSpec, error) {
	var res model.VolumeAttachmentSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateAttachmentURL(atcID)}, "/")

	if err := v.Recv(request, url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListVolumeAttachments
func (v *VolumeMgr) ListVolumeAttachments() ([]*model.VolumeAttachmentSpec, error) {
	var res []*model.VolumeAttachmentSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateAttachmentURL()}, "/")

	if err := v.Recv(request, url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteVolumeAttachment
func (v *VolumeMgr) DeleteVolumeAttachment(atcID string, body VolumeAttachmentBuilder) error {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateAttachmentURL(atcID)}, "/")

	return v.Recv(request, url, "DELETE", body, nil)
}

// CreateVolumeSnapshot
func (v *VolumeMgr) CreateVolumeSnapshot(body VolumeSnapshotBuilder) (*model.VolumeSnapshotSpec, error) {
	var res model.VolumeSnapshotSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateSnapshotURL()}, "/")

	if err := v.Recv(request, url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetVolumeSnapshot
func (v *VolumeMgr) GetVolumeSnapshot(snpID string) (*model.VolumeSnapshotSpec, error) {
	var res model.VolumeSnapshotSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateSnapshotURL(snpID)}, "/")

	if err := v.Recv(request, url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListVolumeSnapshots
func (v *VolumeMgr) ListVolumeSnapshots() ([]*model.VolumeSnapshotSpec, error) {
	var res []*model.VolumeSnapshotSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateSnapshotURL()}, "/")

	if err := v.Recv(request, url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteVolumeSnapshot
func (v *VolumeMgr) DeleteVolumeSnapshot(snpID string, body VolumeSnapshotBuilder) error {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateSnapshotURL(snpID)}, "/")

	return v.Recv(request, url, "DELETE", body, nil)
}

// UpdateVolumeSnapshot
func (v *VolumeMgr) UpdateVolumeSnapshot(snpID string, body VolumeSnapshotBuilder) (*model.VolumeSnapshotSpec, error) {
	var res model.VolumeSnapshotSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateSnapshotURL(snpID)}, "/")

	if err := v.Recv(request, url, "PUT", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
