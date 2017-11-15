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
	"fmt"

	"github.com/opensds/opensds/pkg/model"
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

func NewVolumeMgr(edp string) *VolumeMgr {
	return &VolumeMgr{
		Receiver: NewReceiver(),
		Endpoint: edp,
	}
}

type VolumeMgr struct {
	Receiver

	Endpoint string
}

func (v *VolumeMgr) CreateVolume(body VolumeBuilder) (*model.VolumeSpec, error) {
	var res model.VolumeSpec
	url := v.Endpoint + "/v1alpha/block/volumes"

	if err := v.Recv(request, url, "POST", body, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (v *VolumeMgr) GetVolume(volID string) (*model.VolumeSpec, error) {
	var res model.VolumeSpec
	url := v.Endpoint + "/v1alpha/block/volumes/" + volID

	if err := v.Recv(request, url, "GET", nil, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (v *VolumeMgr) ListVolumes() ([]*model.VolumeSpec, error) {
	var res []*model.VolumeSpec
	url := v.Endpoint + "/v1alpha/block/volumes"

	if err := v.Recv(request, url, "GET", nil, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func (v *VolumeMgr) DeleteVolume(volID string, body VolumeBuilder) error {
	url := v.Endpoint + "/v1alpha/block/volumes/" + volID

	return v.Recv(request, url, "DELETE", body, nil)
}

func (v *VolumeMgr) CreateVolumeAttachment(body VolumeAttachmentBuilder) (*model.VolumeAttachmentSpec, error) {
	var res model.VolumeAttachmentSpec
	url := v.Endpoint + "/v1alpha/block/attachments"

	if err := v.Recv(request, url, "POST", body, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (v *VolumeMgr) UpdateVolumeAttachment(atcID string, body VolumeAttachmentBuilder) (*model.VolumeAttachmentSpec, error) {
	var res model.VolumeAttachmentSpec
	url := v.Endpoint + "/v1alpha/block/attachments/" + atcID

	if err := v.Recv(request, url, "PUT", body, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (v *VolumeMgr) GetVolumeAttachment(atcID string) (*model.VolumeAttachmentSpec, error) {
	var res model.VolumeAttachmentSpec
	url := v.Endpoint + "/v1alpha/block/attachments/" + atcID

	if err := v.Recv(request, url, "GET", nil, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (v *VolumeMgr) ListVolumeAttachments() ([]*model.VolumeAttachmentSpec, error) {
	var res []*model.VolumeAttachmentSpec
	url := v.Endpoint + "/v1alpha/block/attachments"

	if err := v.Recv(request, url, "GET", nil, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func (v *VolumeMgr) DeleteVolumeAttachment(atcID string, body VolumeAttachmentBuilder) error {
	url := v.Endpoint + "/v1alpha/block/attachments/" + atcID

	return v.Recv(request, url, "DELETE", body, nil)
}

func (v *VolumeMgr) CreateVolumeSnapshot(body VolumeSnapshotBuilder) (*model.VolumeSnapshotSpec, error) {
	var res model.VolumeSnapshotSpec
	url := v.Endpoint + "/v1alpha/block/snapshots"

	if err := v.Recv(request, url, "POST", body, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (v *VolumeMgr) GetVolumeSnapshot(snpID string) (*model.VolumeSnapshotSpec, error) {
	var res model.VolumeSnapshotSpec
	url := v.Endpoint + "/v1alpha/block/snapshots/" + snpID

	if err := v.Recv(request, url, "GET", nil, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (v *VolumeMgr) ListVolumeSnapshots() ([]*model.VolumeSnapshotSpec, error) {
	var res []*model.VolumeSnapshotSpec
	url := v.Endpoint + "/v1alpha/block/snapshots"

	if err := v.Recv(request, url, "GET", nil, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func (v *VolumeMgr) DeleteVolumeSnapshot(snpID string, body VolumeSnapshotBuilder) error {
	url := v.Endpoint + "/v1alpha/block/snapshots/" + snpID

	return v.Recv(request, url, "DELETE", body, nil)
}
