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
	api "github.com/opensds/opensds/pkg/model"
)

type Plugin struct{}

func (p *Plugin) Setup() {}

func (p *Plugin) Unset() {}

func (p *Plugin) CreateVolume(name string, size int64) (*api.VolumeSpec, error) {
	return &api.VolumeSpec{}, nil
}

func (p *Plugin) GetVolume(volID string) (*api.VolumeSpec, error) {
	return &api.VolumeSpec{}, nil
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
	return &api.VolumeSnapshotSpec{}, nil
}

func (p *Plugin) GetSnapshot(snapID string) (*api.VolumeSnapshotSpec, error) {
	return &api.VolumeSnapshotSpec{}, nil
}

func (p *Plugin) DeleteSnapshot(snapID string) error {
	return nil
}
