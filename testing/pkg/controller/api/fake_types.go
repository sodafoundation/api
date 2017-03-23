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

type VersionInfo struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type AvailableVersions struct {
	Versions []VersionInfo `json:"versions"`
}

type DefaultResponse struct {
	Error  string `json:"error"`
	Status string `json:"status"`
}

// VolumeResponse is a structure for all properties of
// an volume for a non detailed query
type VolumeResponse struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Status      string              `json:"status"`
	Size        int                 `json:"size"`
	VolumeType  string              `json:"volume_type"`
	Attachments []map[string]string `json:"attachments"`
}

// VolumeDetailResponse is a structure for all properties of
// an volume for a detailed query
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

// ShareResponse is a structure for all properties of
// an share for a non detailed query
type ShareResponse struct {
	Id    string              `json:"id,omitempty"`
	Name  string              `json:"name,omitempty"`
	Links []map[string]string `json:"links"`
}

// ShareDetailResponse is a structure for all properties of
// an share for a detailed query
type ShareDetailResponse struct {
	Links                    []map[string]string `json:"links"`
	AvailabilityZone         string              `json:"availability_zone,omitempty"`
	ShareNetworkId           string              `json:"share_network_id,omitempty"`
	ExportLocations          []string            `json:"export_locations"`
	ShareServerId            string              `json:"share_server_id,omitempty"`
	SnapshotId               string              `json:"snapshot_id,omitempty"`
	Id                       string              `json:"id,omitempty"`
	Size                     int                 `json:"size"`
	ShareType                string              `json:"share_type,omitempty"`
	ShareTypeName            string              `json:"share_type_name,omitempty"`
	ExportLocation           string              `json:"export_location,omitempty"`
	ConsistencyGroupId       string              `json:"consistency_group_id,omitempty"`
	ProjectId                string              `json:"project_id,omitempty"`
	Metadata                 map[string]string   `json:"metadata"`
	Status                   string              `json:"status,omitempty"`
	AccessRulesStatus        string              `json:"access_rules_status,omitempty"`
	Description              string              `json:"description,omitempty"`
	Host                     string              `json:"host,omitempty"`
	TaskState                string              `json:"task_state,omitempty"`
	IsPublic                 bool                `json:"is_public"`
	SnapshotSupport          bool                `json:"snapshot_support"`
	Name                     string              `json:"name,omitempty"`
	HasReplicas              bool                `json:"has_replicas"`
	ReplicationType          string              `json:"replication_type,omitempty"`
	ShareProto               string              `json:"share_proto,omitempty"`
	VolumeType               string              `json:"volume_type,omitempty"`
	SourceCgsnapshotMemberId string              `json:"source_cgsnapshot_member_id,omitempty"`
}
