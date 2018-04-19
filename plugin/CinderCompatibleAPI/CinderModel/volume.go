// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package CinderModel

// *******************Create*******************

// CreateVolumeReqSpec ...
type CreateVolumeReqSpec struct {
	Volume         VolumeOfCreateVolumeReq `json:"volume"`
	SchedulerHints SchedulerHints          `json:"OS-SCH-HNT:scheduler_hints,omitempty"`
}

// VolumeOfCreateVolumeReq ...
type VolumeOfCreateVolumeReq struct {
	Size               int64             `json:"size"`
	AvailabilityZone   string            `json:"availability_zone,omitempty"`
	SourceVolID        string            `json:"source_volid,omitempty"`
	Description        string            `json:"description,omitempty"`
	Multiattach        bool              `json:"multiattach,omitempty"`
	SnapshotID         string            `json:"snapshot_id,omitempty"`
	BackupID           string            `json:"backup_id,omitempty"`
	Name               string            `json:"name"`
	ImageRef           string            `json:"imageRef,omitempty"`
	VolumeType         string            `json:"volume_type,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	ConsistencygroupID string            `json:"consistencygroup_id,omitempty"`
}

// SchedulerHints ...
type SchedulerHints struct {
	SameHost []string `json:"same_host,omitempty"`
}

// CreateVolumeRespSpec ...
type CreateVolumeRespSpec struct {
	Volume VolumeOfCreateVolumeResp `json:"volume,omitempty"`
}

// VolumeOfCreateVolumeResp ...
type VolumeOfCreateVolumeResp struct {
	MigrationStatus    string                   `json:"migration_status,omitempty"`
	Attachments        []AttachmentOfVolumeResp `json:"attachments"`
	Links              []Link                   `json:"links,omitempty"`
	AvailabilityZone   string                   `json:"availability_zone,omitempty"`
	Encrypted          bool                     `json:"encrypted ,omitempty"`
	UpdatedAt          string                   `json:"updated_at,omitempty"`
	ReplicationStatus  string                   `json:"replication_status,omitempty"`
	SnapshotID         string                   `json:"snapshot_id,omitempty"`
	ID                 string                   `json:"id,omitempty"`
	Size               int64                    `json:"size,omitempty"`
	UserID             string                   `json:"user_id,omitempty"`
	Metadata           map[string]string        `json:"metadata"`
	Status             string                   `json:"status,omitempty"`
	Description        string                   `json:"description,omitempty"`
	Multiattach        bool                     `json:"multiattach,omitempty"`
	SourceVolID        string                   `json:"source_volid,omitempty"`
	ConsistencygroupID string                   `json:"consistencygroup_id,omitempty"`
	Name               string                   `json:"name,omitempty"`
	Bootable           bool                     `json:"bootable,omitempty"`
	CreatedAt          string                   `json:"created_at,omitempty"`
	VolumeType         string                   `json:"volume_type,omitempty"`
}

// Link ...
type Link struct {
	Href string `json:"href,omitempty"`
	Rel  string `json:"rel,omitempty"`
}

// *******************List*******************

// ListVolumeRespSpec ...
type ListVolumeRespSpec struct {
	Volumes []VolumeForListResp `json:"volumes"`
	Count   int64               `json:"count,omitempty"`
}

// VolumeForListResp ...
type VolumeForListResp struct {
	ID    string `json:"id"`
	Links []Link `json:"links,omitempty"`
	Name  string `json:"name"`
}

// *******************List Detail*******************

// ListVolumeDetailRespSpec ...
type ListVolumeDetailRespSpec struct {
	Volumes []VolumeForListDetailResp `json:"volumes"`
	Count   int64                     `json:"count,omitempty"`
}

// VolumeForListDetailResp ...
type VolumeForListDetailResp struct {
	MigrationStatus     string                   `json:"migration_status,omitempty"`
	Attachments         []AttachmentOfVolumeResp `json:"attachments"`
	Links               []Link                   `json:"links,omitempty"`
	AvailabilityZone    string                   `json:"availability_zone,omitempty"`
	Host                string                   `json:"os-vol-host-attr:host,omitempty"`
	Encrypted           bool                     `json:"encrypted ,omitempty"`
	UpdatedAt           string                   `json:"updated_at"`
	ReplicationStatus   string                   `json:"replication_status,omitempty"`
	SnapshotID          string                   `json:"snapshot_id,omitempty"`
	ID                  string                   `json:"id"`
	Size                int64                    `json:"size"`
	UserID              string                   `json:"user_id"`
	TenantID            string                   `json:"os-vol-tenant-attr:tenant_id,omitempty"`
	Migstat             string                   `json:"os-vol-mig-status-attr:migstat,omitempty"`
	Metadata            map[string]string        `json:"metadata"`
	Status              string                   `json:"status"`
	VolumeImageMetadata map[string]string        `json:"volume_image_metadata ,omitempty"`
	Description         string                   `json:"description"`
	Multiattach         bool                     `json:"multiattach,omitempty"`
	SourceVolID         string                   `json:"source_volid,omitempty"`
	ConsistencygroupID  string                   `json:"consistencygroup_id,omitempty"`
	NameID              string                   `json:"os-vol-mig-status-attr:name_id,omitempty"`
	Name                string                   `json:"name"`
	Bootable            bool                     `json:"bootable,omitempty"`
	CreatedAt           string                   `json:"created_at"`
	VolumeType          string                   `json:"volume_type,omitempty"`
}

// *******************Show*******************

// ShowVolumeRespSpec ...
type ShowVolumeRespSpec struct {
	Volume VolumeOfShowVolumeResp `json:"volume"`
}

// VolumeOfShowVolumeResp ...
type VolumeOfShowVolumeResp struct {
	MigrationStatus     string                   `json:"migration_status,omitempty"`
	Attachments         []AttachmentOfVolumeResp `json:"attachments"`
	Links               []Link                   `json:"links,omitempty"`
	AvailabilityZone    string                   `json:"availability_zone,omitempty"`
	Host                string                   `json:"os-vol-host-attr:host,omitempty"`
	Encrypted           bool                     `json:"encrypted ,omitempty"`
	UpdatedAt           string                   `json:"updated_at"`
	ReplicationStatus   string                   `json:"replication_status,omitempty"`
	SnapshotID          string                   `json:"snapshot_id,omitempty"`
	ID                  string                   `json:"id"`
	Size                int64                    `json:"size"`
	UserID              string                   `json:"user_id"`
	TenantID            string                   `json:"os-vol-tenant-attr:tenant_id,omitempty"`
	Migstat             string                   `json:"os-vol-mig-status-attr:migstat,omitempty"`
	Metadata            map[string]string        `json:"metadata"`
	Status              string                   `json:"status"`
	VolumeImageMetadata map[string]string        `json:"volume_image_metadata,omitempty"`
	Description         string                   `json:"description"`
	Multiattach         bool                     `json:"multiattach,omitempty"`
	SourceVolID         string                   `json:"source_volid,omitempty"`
	ConsistencygroupID  string                   `json:"consistencygroup_id,omitempty"`
	NameID              string                   `json:"os-vol-mig-status-attr:name_id,omitempty"`
	Name                string                   `json:"name"`
	Bootable            bool                     `json:"bootable,omitempty"`
	CreatedAt           string                   `json:"created_at"`
	VolumeType          string                   `json:"volume_type,omitempty"`
	ServiceUuID         string                   `json:"service_uuid,omitempty"`
	SharedTargets       bool                     `json:"shared_targets,omitempty"`
}

// *******************Update*******************

// UpdateVolumeReqSpec ...
type UpdateVolumeReqSpec struct {
	Volume VolumeOfUpdateVolumeReq `json:"volume"`
}

// VolumeOfUpdateVolumeReq ...
type VolumeOfUpdateVolumeReq struct {
	Description string            `json:"description,omitempty"`
	Name        string            `json:"name,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// UpdateVolumeRespSpec ...
type UpdateVolumeRespSpec struct {
	Volume VolumeOfUpdateVolumeResp `json:"volume"`
}

// VolumeOfUpdateVolumeResp ...
type VolumeOfUpdateVolumeResp struct {
	MigrationStatus    string                   `json:"migration_status,omitempty"`
	Attachments        []AttachmentOfVolumeResp `json:"attachments"`
	Links              []Link                   `json:"links,omitempty"`
	AvailabilityZone   string                   `json:"availability_zone,omitempty"`
	Encrypted          bool                     `json:"encrypted ,omitempty"`
	UpdatedAt          string                   `json:"updated_at"`
	ReplicationStatus  string                   `json:"replication_status,omitempty"`
	SnapshotID         string                   `json:"snapshot_id,omitempty"`
	ID                 string                   `json:"id"`
	Size               int64                    `json:"size"`
	UserID             string                   `json:"user_id"`
	Metadata           map[string]string        `json:"metadata"`
	Status             string                   `json:"status"`
	Description        string                   `json:"description"`
	Multiattach        bool                     `json:"multiattach,omitempty"`
	SourceVolID        string                   `json:"source_volid,omitempty"`
	ConsistencygroupID string                   `json:"consistencygroup_id,omitempty"`
	Name               string                   `json:"name"`
	Bootable           bool                     `json:"bootable,omitempty"`
	CreatedAt          string                   `json:"created_at"`
	VolumeType         string                   `json:"volume_type,omitempty"`
}

// AttachmentOfVolumeResp ...
type AttachmentOfVolumeResp struct {
	ID string `json:"id,omitempty"`
}
