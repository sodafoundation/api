// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
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

/*
This module implements a sample driver for OpenSDS. This driver will handle all
operations of volume and return a fake value.

*/

package upsplugin

import (
	"encoding/json"
	"io/ioutil"
	"log"

	api "github.com/opensds/opensds/pkg/model"
)

const (
	upspluginPoolConfig = "/etc/opensds/pool.json"
)

type Plugin struct{}

func (p *Plugin) Setup() {}

func (p *Plugin) Unset() {}

func (p *Plugin) CreateVolume(name string, size int64) (*api.VolumeSpec, error) {
	return &api.VolumeSpec{BaseModel: &api.BaseModel{}}, nil
}

func (p *Plugin) GetVolume(volID string) (*api.VolumeSpec, error) {
	return &api.VolumeSpec{BaseModel: &api.BaseModel{}}, nil
}

func (p *Plugin) DeleteVolume(volID string) error {
	return nil
}

func (p *Plugin) InitializeConnection(volID string, doLocalAttach, multiPath bool, hostInfo *api.HostInfo) (*api.ConnectionInfo, error) {
	return &api.ConnectionInfo{}, nil
}

func (p *Plugin) AttachVolume(volID, host, mountpoint string) error {
	return nil
}

func (p *Plugin) DetachVolume(volID string) error {
	return nil
}

func (p *Plugin) CreateSnapshot(name, volID, description string) (*api.VolumeSnapshotSpec, error) {
	return &api.VolumeSnapshotSpec{
		BaseModel: &api.BaseModel{},
		VolumeId:  volID,
	}, nil
}

func (p *Plugin) GetSnapshot(snapID string) (*api.VolumeSnapshotSpec, error) {
	return &api.VolumeSnapshotSpec{BaseModel: &api.BaseModel{}}, nil
}

func (p *Plugin) DeleteSnapshot(snapID string) error {
	return nil
}

func (p *Plugin) ListPools() (*[]api.StoragePoolSpec, error) {
	pools, err := readPoolsFromFile()
	if err != nil {
		log.Println("Could not read pool resource:", err)
		return &[]api.StoragePoolSpec{}, err
	}

	return &pools, nil
}

func readPoolsFromFile() ([]api.StoragePoolSpec, error) {
	var pools []api.StoragePoolSpec

	userJSON, err := ioutil.ReadFile(upspluginPoolConfig)
	if err != nil {
		log.Println("ReadFile json failed:", err)
		return pools, err
	}

	// If the pool resource is empty, consider it as a normal condition
	if len(userJSON) == 0 {
		return pools, nil
	}

	// Unmarshal the result
	if err = json.Unmarshal(userJSON, &pools); err != nil {
		log.Println("Unmarshal json failed:", err)
		return pools, err
	}
	return pools, nil
}
