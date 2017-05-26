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
This module implements the common data structure.

*/

package api

// StorageProfile is a structure for all properties of
// profile configured by admin
type StorageProfile struct {
	Id            string            `json:"id"`
	Name          string            `json:"name"`
	BackendDriver string            `json:"backend"`
	StorageTags   map[string]string `json:"tags"`
}

// VolumeOperationSchema is a structure for all properties of
// volume operation
type VolumeOperationSchema struct {
	// Some properties related to basic operation of volumes
	DockId       string `json:"dockId,omitempty"`
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	VolumeType   string `json:"volumeType"`
	Size         int32  `json:"size"`
	AllowDetails bool   `json:"allowDetails"`

	// Some properties related to basic operation of volume attachments
	DoLocalAttach bool `json:"doLocalAttach"`
	MultiPath     bool `json:"multipath"`
	HostInfo      `json:"hostInfo"`
	Mountpoint    string `json:"mountpoint"`

	// Some properties related to basic operation of volume snapshots
	SnapshotId      string `json:"snapshotId,omitempty"`
	SnapshotName    string `json:"snapshotName,omitempty"`
	Description     string `json:"description,omitempty"`
	ForceSnapshoted bool   `json:"forceSnapshoted,omitempty"`
}

// VolumeRequest is a structure for all properties of
// a volume request
type VolumeRequest struct {
	Schema  *VolumeOperationSchema `json:"schema"`
	Profile *StorageProfile        `json:"profile"`
}

// VolumeResponse is a structure for all properties of
// an volume for a non detailed query
type VolumeResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ConnectorProperties struct {
	DoLocalAttach bool   `json:"do_local_attach"`
	Platform      string `json:"platform"`
	OsType        string `json:"os_type"`
	Ip            string `json:"ip"`
	Host          string `json:"host"`
	MultiPath     bool   `json:"multipath"`
	Initiator     string `json:"initiator"`
}

// ConnectionInfo is a structure for all properties of
// connection when create a volume attachment
type ConnectionInfo struct {
	DriverVolumeType        string      `json:"driver_volume_type"`
	ConnectionDataContainer interface{} `json:"data"`
}

// HostInfo is a structure for all properties of host
// when create a volume attachment
type HostInfo struct {
	Platform  string `json:"platform"`
	OsType    string `json:"osType"`
	Ip        string `json:"ip"`
	Host      string `json:"host"`
	Initiator string `json:"initiator"`
}

// VolumeSnapshotResponse is a structure for all properties of
// a volume attachment for a non detailed query
type VolumeAttachment struct {
	Id             string `json:"id"`
	Mountpoint     string `json:"mountpoint"`
	Status         string `json:"status"`
	HostInfo       `json:"hostInfo"`
	ConnectionInfo `json:"connectionInfo"`
}
