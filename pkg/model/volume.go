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

/*
This module implements the common data structure.

*/

package model

import (
	"encoding/json"
)

type VolumeSpec struct {
	*BaseModel
	Name             string            `json:"name,omitempty"`
	Description      string            `json:"description,omitempty"`
	Size             int64             `json:"size,omitempty"`
	AvailabilityZone string            `json:"availabilityZone,omitempty"`
	Status           string            `json:"status,omitempty"`
	PoolId           string            `json:"poolId,omitempty"`
	ProfileId        string            `json:"profileId,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

func (vol *VolumeSpec) GetName() string {
	return vol.Name
}

func (vol *VolumeSpec) GetDescription() string {
	return vol.Description
}

func (vol *VolumeSpec) GetSize() int64 {
	return vol.Size
}

func (vol *VolumeSpec) GetAvailabilityZone() string {
	return vol.AvailabilityZone
}

func (vol *VolumeSpec) GetPoolId() string {
	return vol.PoolId
}

func (vol *VolumeSpec) GetProfileId() string {
	return vol.ProfileId
}

func (vol *VolumeSpec) GetMetadata() map[string]string {
	return vol.Metadata
}

type VolumeAttachmentSpec struct {
	*BaseModel
	VolumeId        string            `json:"volumeId,omitempty"`
	Mountpoint      string            `json:"mountpoint,omitempty"`
	Status          string            `json:"status,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	*HostInfo       `json:"hostInfo,omitempty"`
	*ConnectionInfo `json:"connectionInfo,omitempty"`
}

func (atc *VolumeAttachmentSpec) GetVolumeId() string {
	return atc.VolumeId
}

func (atc *VolumeAttachmentSpec) GetMountpoint() string {
	return atc.Mountpoint
}

func (atc *VolumeAttachmentSpec) GetMetadata() map[string]string {
	return atc.Metadata
}

// HostInfo is a structure for all properties of host
// when create a volume attachment
type HostInfo struct {
	Platform  string `json:"platform,omitempty"`
	OsType    string `json:"osType,omitempty"`
	Ip        string `json:"ip,omitempty"`
	Host      string `json:"host,omitempty"`
	Initiator string `json:"initiator,omitempty"`
}

func (host *HostInfo) GetPlatform() string {
	return host.Platform
}

func (host *HostInfo) GetOsType() string {
	return host.OsType
}

func (host *HostInfo) GetIp() string {
	return host.Ip
}

func (host *HostInfo) GetHost() string {
	return host.Host
}

func (host *HostInfo) GetInitiator() string {
	return host.Initiator
}

// ConnectionInfo is a structure for all properties of
// connection when create a volume attachment
type ConnectionInfo struct {
	DriverVolumeType string                 `json:"driverVolumeType,omitempty"`
	ConnectionData   map[string]interface{} `json:"data,omitempty"`
}

func (con *ConnectionInfo) GetDriverVolumeType() string {
	return con.DriverVolumeType
}

func (con *ConnectionInfo) GetConnectionData() map[string]interface{} {
	return con.ConnectionData
}

func (con *ConnectionInfo) EncodeConnectionData() []byte {
	conBody, _ := json.Marshal(&con.ConnectionData)
	return conBody
}

type VolumeSnapshotSpec struct {
	*BaseModel
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Size        int64             `json:"size,omitempty"`
	Status      string            `json:"status,omitempty"`
	VolumeId    string            `json:"volumeId,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

func (snp *VolumeSnapshotSpec) GetName() string {
	return snp.Name
}

func (snp *VolumeSnapshotSpec) GetDescription() string {
	return snp.Description
}

func (snp *VolumeSnapshotSpec) GetSize() int64 {
	return snp.Size
}

func (snp *VolumeSnapshotSpec) GetVolumeId() string {
	return snp.VolumeId
}

func (snp *VolumeSnapshotSpec) GetMetadata() map[string]string {
	return snp.Metadata
}
