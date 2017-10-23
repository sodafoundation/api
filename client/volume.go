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

package client

import (
	"errors"
	"fmt"
	"sync"

	"github.com/opensds/opensds/pkg/model"
)

type VolumeMgr struct {
	Receiver

	Endpoint string
	Opt      map[string]string
	Body     interface{}
	lock     sync.Mutex
}

func NewVolumeMgr(edp string) *VolumeMgr {
	return &VolumeMgr{
		Receiver: NewReceiver(),
		Endpoint: edp,
	}
}

func (v *VolumeMgr) CreateVolume() (*model.VolumeSpec, error) {
	var res model.VolumeSpec
	url := v.Endpoint + "/api/v1alpha/block/volumes"

	if err := v.Recv(request, url, "POST", v.Body, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (v *VolumeMgr) GetVolume(volID string) (*model.VolumeSpec, error) {
	var res model.VolumeSpec
	url := v.Endpoint + "/api/v1alpha/block/volumes/" + volID

	if err := v.Recv(request, url, "GET", v.Opt, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (v *VolumeMgr) ListVolumes() ([]*model.VolumeSpec, error) {
	var res []*model.VolumeSpec
	url := v.Endpoint + "/api/v1alpha/block/volumes"

	if err := v.Recv(request, url, "GET", v.Opt, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func (v *VolumeMgr) DeleteVolume(volID string) *model.Response {
	var res model.Response
	url := v.Endpoint + "/api/v1alpha/block/volumes" + volID

	if err := v.Recv(request, url, "DELETE", v.Body, &res); err != nil {
		res.Status, res.Error = "Failure", fmt.Sprint(err)
	}

	return &res
}

func (v *VolumeMgr) CreateVolumeSnapshot() (*model.VolumeSnapshotSpec, error) {
	var res model.VolumeSnapshotSpec
	url := v.Endpoint + "/api/v1alpha/block/snapshots"

	if err := v.Recv(request, url, "POST", v.Body, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (v *VolumeMgr) GetVolumeSnapshot(snpID string) (*model.VolumeSnapshotSpec, error) {
	var res model.VolumeSnapshotSpec
	url := v.Endpoint + "/api/v1alpha/block/snapshots/" + snpID

	if err := v.Recv(request, url, "GET", v.Opt, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (v *VolumeMgr) ListVolumeSnapshots() ([]*model.VolumeSnapshotSpec, error) {
	var res []*model.VolumeSnapshotSpec
	url := v.Endpoint + "/api/v1alpha/block/snapshots"

	if err := v.Recv(request, url, "GET", v.Opt, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func (v *VolumeMgr) DeleteVolumeSnapshot(snpID string) *model.Response {
	var res model.Response
	url := v.Endpoint + "/api/v1alpha/block/snapshots" + snpID

	if err := v.Recv(request, url, "DELETE", v.Body, &res); err != nil {
		res.Status, res.Error = "Failure", fmt.Sprint(err)
	}

	return &res
}

func (v *VolumeMgr) ResetAndUpdateVolumeRequestContent(in interface{}) error {
	var err error

	v.lock.Lock()
	defer v.lock.Unlock()
	// Clear all content stored in Opt field.
	v.Opt, v.Body = make(map[string]string), nil
	// Valid the input data.
	switch in.(type) {
	case map[string]string:
		v.Opt = in.(map[string]string)
		break
	case model.VolumeSpec, *model.VolumeSpec, model.VolumeSnapshotSpec, *model.VolumeSnapshotSpec:
		v.Body = in
		break
	default:
		err = errors.New("Request content type not supported")
	}

	return err
}
