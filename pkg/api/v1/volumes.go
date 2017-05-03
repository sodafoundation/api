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
	DockId       string `json:"dockId,omitempty"`
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	VolumeType   string `json:"volumeType"`
	Size         int32  `json:"size"`
	AllowDetails bool   `json:"allowDetails"`

	// Some properties related to attach and mount operation of volumes
	Device   string `json:"device,omitempty"`
	MountDir string `json:"mountDir,omitempty"`
	FsType   string `json:"fsType,omitempty"`

	// Some properties related to basic operation of volume snapshots
	SnapshotId      string `json:"snapshotId,omitempty"`
	SnapshotName    string `json:"snapshotName,omitempty"`
	Description     string `json:"description,omitempty"`
	ForceSnapshoted bool   `json:"forceSnapshoted,omitempty"`
}

// VolumeResponse is a structure for all properties of
// a volume for a non detailed query
type VolumeResponse struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Status      string              `json:"status"`
	Size        int                 `json:"size"`
	VolumeType  string              `json:"volume_type"`
	Attachments []map[string]string `json:"attachments"`
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

// VolumeSnapshotResponse is a structure for all properties of
// a volume snapshot for a non detailed query
type VolumeSnapshotResponse struct {
	Status   string `json:"status,omitempty"`
	Name     string `json:"name,omitempty"`
	VolumeId string `json:"volume_id,omitempty"`
	Size     int    `json:"size"`
	Id       string `json:"id,omitempty"`
}
