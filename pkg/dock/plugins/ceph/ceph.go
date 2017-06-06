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
This module implements ceph driver for OpenSDS. Ceph driver will pass these
operation requests about volume to go-ceph module.

*/

package ceph

import (
	api "github.com/opensds/opensds/pkg/api/v1"
)

type CephPlugin struct{}

func (plugin *CephPlugin) Setup() {}

func (plugin *CephPlugin) Unset() {}

func (plugin *CephPlugin) CreateVolume(name string, size int32) (*api.VolumeResponse, error) {
	return &api.VolumeResponse{}, nil
}

func (plugin *CephPlugin) GetVolume(volID string) (*api.VolumeResponse, error) {
	return &api.VolumeResponse{}, nil
}

func (plugin *CephPlugin) DeleteVolume(volID string) error {
	return nil
}

func (plugin *CephPlugin) InitializeConnection(volID string, doLocalAttach, multiPath bool, hostInfo *api.HostInfo) (*api.ConnectionInfo, error) {
	return &api.ConnectionInfo{}, nil
}

func (plugin *CephPlugin) AttachVolume(volID, host, mountpoint string) error {
	return nil
}

func (plugin *CephPlugin) DetachVolume(volID string) error {
	return nil
}

func (plugin *CephPlugin) CreateSnapshot(name, volID, description string) (*api.VolumeSnapshot, error) {
	return &api.VolumeSnapshot{}, nil
}

func (plugin *CephPlugin) GetSnapshot(snapID string) (*api.VolumeSnapshot, error) {
	return &api.VolumeSnapshot{}, nil
}

func (plugin *CephPlugin) DeleteSnapshot(snapID string) error {
	return nil
}
