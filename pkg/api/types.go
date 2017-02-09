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

import "time"

// VolumeResponse is a structure for all properties of
// an volume for a non detailed query
type VolumeResponse struct {
	ID    string              `json:"id"`
	Name  string              `json:"name"`
	Links []map[string]string `json:"links"`
	// Consistencygroup_id string `json:"consistencygroup_id"`
}

// VolumeDetailResponse is a structure for all properties of
// an volume for a detailed query
type VolumeDetailResponse struct {
	ID              string              `json:"id"`
	Attachments     []map[string]string `json:"attachments"`
	Links           []map[string]string `json:"links"`
	Metadata        map[string]string   `json:"metadata"`
	Protected       bool                `json:"protected"`
	Status          string              `json:"status"`
	MigrationStatus string              `json:"migration_status"`
	UserID          string              `json:"user_id"`
	Encrypted       bool                `json:"encrypted"`
	Multiattach     bool                `json:"multiattach"`
	CreatedAt       time.Time           `json:"created_at"`
	Description     string              `json:"description"`
	Volume_type     string              `json:"volume_type"`
	Name            string              `json:"name"`
	Source_volid    string              `json:"source_volid"`
	Snapshot_id     string              `json:"snapshot_id"`
	Size            int64               `json:"size"`

	Aavailability_zone  string `json:"availability_zone"`
	Rreplication_status string `json:"replication_status"`
	Consistencygroup_id string `json:"consistencygroup_id"`
}

// ShareResponse is a structure for all properties of
// an share for a non detailed query
type ShareResponse struct {
	ID    string              `json:"id"`
	Name  string              `json:"name"`
	Links []map[string]string `json:"links"`
}

// ShareDetailResponse is a structure for all properties of
// an share for a detailed query
type ShareDetailResponse struct {
	Links                       []map[string]string `json:"links"`
	Availability_zone           string              `json:"availability_zone"`
	Share_network_id            string              `json:"share_network_id"`
	Export_locations            []string            `json:"export_locations"`
	Share_server_id             string              `json:"share_server_id"`
	Snapshot_id                 string              `json:"snapshot_id"`
	ID                          string              `json:"id"`
	Size                        int64               `json:"size"`
	Share_type                  string              `json:"share_type"`
	Share_type_name             string              `json:"share_type_name"`
	Export_location             string              `json:"export_location"`
	Consistency_group_id        string              `json:"consistency_group_id"`
	Project_id                  string              `json:"project_id"`
	Metadata                    map[string]string   `json:"metadata"`
	Status                      string              `json:"status"`
	Access_rules_status         string              `json:"access_rules_status"`
	Description                 string              `json:"description"`
	Host                        string              `json:"host"`
	Task_state                  string              `json:"task_state"`
	Is_public                   bool                `json:"is_public"`
	Snapshot_support            bool                `json:"snapshot_support"`
	Name                        string              `json:"name"`
	Has_replicas                bool                `json:"has_replicas"`
	Replication_type            string              `json:"replication_type"`
	CreatedAt                   time.Time           `json:"created_at"`
	Share_proto                 string              `json:"share_proto"`
	Volume_type                 string              `json:"volume_type"`
	Source_cgsnapshot_member_id string              `json:"source_cgsnapshot_member_id"`
}
