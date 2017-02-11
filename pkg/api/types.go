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

import (
	"openstack/golang-client/util"
)

// VolumeResponse is a structure for all properties of
// an volume for a non detailed query
type VolumeResponse struct {
	ID    string              `json:"id,omitempty"`
	Name  string              `json:"name,omitempty"`
	Links []map[string]string `json:"links"`
	// Consistencygroup_id string `json:"consistencygroup_id"`
}

// VolumeDetailResponse is a structure for all properties of
// an volume for a detailed query
type VolumeDetailResponse struct {
	ID              string               `json:"id,omitempty"`
	Attachments     []map[string]string  `json:"attachments"`
	Links           []map[string]string  `json:"links"`
	Metadata        map[string]string    `json:"metadata"`
	Protected       bool                 `json:"protected"`
	Status          string               `json:"status,omitempty"`
	MigrationStatus string               `json:"migration_status,omitempty"`
	UserID          string               `json:"user_id,omitempty"`
	Encrypted       bool                 `json:"encrypted"`
	Multiattach     bool                 `json:"multiattach"`
	CreatedAt       util.RFC8601DateTime `json:"created_at"`
	Description     string               `json:"description,omitempty"`
	Volume_type     string               `json:"volume_type,omitempty"`
	Name            string               `json:"name,omitempty"`
	Source_volid    string               `json:"source_volid,omitempty"`
	Snapshot_id     string               `json:"snapshot_id,omitempty"`
	Size            int64                `json:"size"`

	Aavailability_zone  string `json:"availability_zone,omitempty"`
	Rreplication_status string `json:"replication_status,omitempty"`
	Consistencygroup_id string `json:"consistencygroup_id,omitempty"`
}

// ShareResponse is a structure for all properties of
// an share for a non detailed query
type ShareResponse struct {
	ID    string              `json:"id,omitempty"`
	Name  string              `json:"name,omitempty"`
	Links []map[string]string `json:"links"`
}

// ShareDetailResponse is a structure for all properties of
// an share for a detailed query
type ShareDetailResponse struct {
	Links                       []map[string]string  `json:"links"`
	Availability_zone           string               `json:"availability_zone,omitempty"`
	Share_network_id            string               `json:"share_network_id,omitempty"`
	Export_locations            []string             `json:"export_locations"`
	Share_server_id             string               `json:"share_server_id,omitempty"`
	Snapshot_id                 string               `json:"snapshot_id,omitempty"`
	ID                          string               `json:"id,omitempty"`
	Size                        int64                `json:"size"`
	Share_type                  string               `json:"share_type,omitempty"`
	Share_type_name             string               `json:"share_type_name,omitempty"`
	Export_location             string               `json:"export_location,omitempty"`
	Consistency_group_id        string               `json:"consistency_group_id,omitempty"`
	Project_id                  string               `json:"project_id,omitempty"`
	Metadata                    map[string]string    `json:"metadata"`
	Status                      string               `json:"status,omitempty"`
	Access_rules_status         string               `json:"access_rules_status,omitempty"`
	Description                 string               `json:"description,omitempty"`
	Host                        string               `json:"host,omitempty"`
	Task_state                  string               `json:"task_state,omitempty"`
	Is_public                   bool                 `json:"is_public"`
	Snapshot_support            bool                 `json:"snapshot_support"`
	Name                        string               `json:"name,omitempty"`
	Has_replicas                bool                 `json:"has_replicas"`
	Replication_type            string               `json:"replication_type,omitempty"`
	CreatedAt                   util.RFC8601DateTime `json:"created_at"`
	Share_proto                 string               `json:"share_proto,omitempty"`
	Volume_type                 string               `json:"volume_type,omitempty"`
	Source_cgsnapshot_member_id string               `json:"source_cgsnapshot_member_id,omitempty"`
}
