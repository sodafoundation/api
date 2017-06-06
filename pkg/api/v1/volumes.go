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

package v1

type DefaultResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// VolumeOperationSchema is a structure for all properties of
// volume operation
type VolumeOperationSchema struct {
	// Some properties related to basic operation of volumes
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Size int32  `json:"size"`

	// Some properties related to basic operation of volume attachments
	AttachmentId  string `json:"attachmentId"`
	DoLocalAttach bool   `json:"doLocalAttach"`
	MultiPath     bool   `json:"multiPath"`
	HostInfo      `json:"hostInfo"`
	Mountpoint    string `json:"mountpoint"`

	// Some properties related to basic operation of volume snapshots
	SnapshotId      string `json:"snapshotId,omitempty"`
	SnapshotName    string `json:"snapshotName,omitempty"`
	Description     string `json:"description,omitempty"`
	ForceSnapshoted bool   `json:"forceSnapshoted,omitempty"`
}

// VolumeResponse is a structure for all properties of
// a volume for a non detailed query
type VolumeResponse struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	Size             int    `json:"size"`
	PoolName         string `json:"pool_name"`
	AvailabilityZone string `json:"availability_zone"`
	Iops             int32  `json:"iops"`
}

// VolumeDetailResponse is a structure for all properties of
// a volume for a detailed query
type VolumeDetailResponse struct {
	Id              string              `json:"id,omitempty"`
	Attachments     []map[string]string `json:"attachments"`
	Links           []map[string]string `json:"links"`
	Metadata        map[string]string   `json:"metadata"`
	Protected       bool                `json:"protected"`
	Status          string              `json:"status,omitempty"`
	MigrationStatus string              `json:"migration_status,omitempty"`
	UserID          string              `json:"user_id,omitempty"`
	Encrypted       bool                `json:"encrypted"`
	Multiattach     bool                `json:"multiattach"`
	Description     string              `json:"description,omitempty"`
	VolumeType      string              `json:"volume_type,omitempty"`
	Name            string              `json:"name,omitempty"`
	SourceVolid     string              `json:"source_volid,omitempty"`
	SnapshotId      string              `json:"snapshot_id,omitempty"`
	Size            int                 `json:"size"`

	AvailabilityZone   string `json:"availability_zone,omitempty"`
	ReplicationStatus  string `json:"replication_status,omitempty"`
	ConsistencygroupId string `json:"consistencygroup_id,omitempty"`
}

// VolumeSnapshot is a structure for all properties of
// a volume snapshot for a non detailed query
type VolumeSnapshot struct {
	Status   string `json:"status,omitempty"`
	Name     string `json:"name,omitempty"`
	VolumeId string `json:"volume_id,omitempty"`
	Size     int    `json:"size"`
	Id       string `json:"id,omitempty"`
}

// HostInfo is a structure for all properties of host
// when create a volume attachment
type HostInfo struct {
	Platform  string `json:"platform"`
	OsType    string `json:"os_type"`
	Ip        string `json:"ip"`
	Host      string `json:"host"`
	Initiator string `json:"initiator"`
}

// ConnectionInfo is a structure for all properties of
// connection when create a volume attachment
type ConnectionInfo struct {
	DriverVolumeType string                 `json:"driver_volume_type"`
	ConnectionData   map[string]interface{} `json:"data"`
}

// VolumeAttachment is a structure for all properties of
// a volume attachment for a non detailed query
type VolumeAttachment struct {
	Id             string `json:"id"`
	Mountpoint     string `json:"mountpoint"`
	Status         string `json:"status"`
	HostInfo       `json:"hostInfo"`
	ConnectionInfo `json:"connectionInfo"`
}
