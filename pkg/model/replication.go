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

package model

//import (
//	"encoding/json"
//)

// ReplicationSpec represents a replication relationship between the volumes
// on the primary and secondary sites.
type ReplicationSpec struct {
	*BaseModel

	// The uuid of the tenant that the replication belongs to.
	TenantId string `json:"tenantId,omitempty"`

	// The uuid of the user that the replication belongs to.
	// +optional
	UserId string `json:"userId,omitempty"`

	// The name of the replication.
	Name string `json:"name,omitempty"`

	// The description of the replication.
	// +optional
	Description string `json:"description,omitempty"`

	// The status of the replication.
	Status string `json:"status,omitempty"`

	// The uuid of the volume on the primary site.
	PrimaryVolumeId string `json:"primaryVolumeId,omitempty"`

	// The uuid of the volume on the secondary site.
	SecondaryVolumeId string `json:"secondaryVolumeId,omitempty"`

	// NOTE: Need to figure out how to represent the relationship
	// when there are more than 2 sites. May need to use array.
}
