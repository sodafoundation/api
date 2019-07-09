// Copyright 2018 The OpenSDS Authors.
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

import "github.com/opensds/opensds/pkg/model/proto"

const (
	ReplicationModeSync         = "sync"
	ReplicationModeAsync        = "async"
	ReplicationDefaultBackendId = "default"
	ReplicationDefaultPeriod    = 60
)

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
	//ReplicationStatus string `json:"status,omitempty"`

	// The uuid of the volume on the primary site.
	PrimaryVolumeId string `json:"primaryVolumeId,omitempty"`

	// The uuid of the volume on the secondary site.
	SecondaryVolumeId string `json:"secondaryVolumeId,omitempty"`

	// NOTE: Need to figure out how to represent the relationship
	// when there are more than 2 sites. May need to use array.
	AvailabilityZone string `json:"availabilityZone,omitempty"`
	// region
	Region string `json:"region,omitempty"`
	// group id
	GroupId string `json:"groupId,omitempty"`
	// primary replication driver data
	PrimaryReplicationDriverData map[string]string `json:"primaryReplicationDriverData,omitempty"`
	// secondary replication driver data
	SecondaryReplicationDriverData map[string]string `json:"secondaryReplicationDriverData,omitempty"`
	// replication status
	ReplicationStatus string `json:"replicationStatus,omitempty"`
	// supports "async" or "sync" now
	ReplicationMode string `json:"replicationMode,omitempty"`
	// 0 means sync replication.
	ReplicationPeriod int64 `json:"replicationPeriod,omitempty"`
	// replication period
	ReplicationBandwidth int64 `json:"replicationBandwidth,omitempty"`
	// profile id
	ProfileId string `json:"profileId,omitempty"`
	// pool id
	PoolId string `json:"poolId,omitempty"`
	// metadata
	Metadata map[string]string `json:"metadata,omitempty"`
	// volume data list
	VolumeDataList []*proto.VolumeData `json:"volumeDataList,omitempty"`
}

type FailoverReplicationSpec struct {
	AllowAttachedVolume bool   `json:"allowAttachedVolume,omitempty"`
	SecondaryBackendId  string `json:"secondaryBackendId,omitempty"`
}
