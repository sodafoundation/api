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

package model

// Fileshare status
const (
	FileShareCreating      = "creating"
	FileShareAvailable     = "available"
	FileShareInUse         = "in_Use"
	FileShareDeleting      = "deleting"
	FileShareError         = "error"
	FileShareErrorDeleting = "errorDeleting"
)

// fileshare snapshot status
const (
	FileShareSnapCreating      = "creating"
	FileShareSnapAvailable     = "available"
	FileShareSnapDeleting      = "deleting"
	FileShareSnapError         = "error"
	FileShareSnapErrorDeleting = "errorDeleting"
)

// volume status
const (
	VolumeCreating       = "creating"
	VolumeAvailable      = "available"
	VolumeInUse          = "inUse"
	VolumeDeleting       = "deleting"
	VolumeError          = "error"
	VolumeErrorDeleting  = "errorDeleting"
	VolumeErrorExtending = "errorExtending"
	VolumeExtending      = "extending"
)

// volume attach status
const (
	VolumeAttaching      = "attaching"
	VolumeAttached       = "attached"
	VolumeDetached       = "detached"
	VolumeReserved       = "reserved"
	VolumeErrorAttaching = "errorAttaching"
	VolumeErrorDetaching = "errorDetaching"
)

// volume snapshot status
const (
	VolumeSnapCreating      = "creating"
	VolumeSnapAvailable     = "available"
	VolumeSnapDeleting      = "deleting"
	VolumeSnapError         = "error"
	VolumeSnapErrorDeleting = "errorDeleting"
)

// volume attachment status
const (
	VolumeAttachCreating      = "creating"
	VolumeAttachAvailable     = "available"
	VolumeAttachErrorDeleting = "errorDeleting"
	VolumeAttachError         = "error"
)

//volume replication status
const (
	ReplicationDeleted        = "deleted"
	ReplicationCreating       = "creating"
	ReplicationDeleting       = "deleting"
	ReplicationEnabling       = "enabling"
	ReplicationDisabling      = "disabling"
	ReplicationFailingOver    = "failing_over"
	ReplicationFailingBack    = "failing_back"
	ReplicationAvailable      = "available"
	ReplicationError          = "error"
	ReplicationErrorDeleting  = "error_deleting"
	ReplicationErrorEnabling  = "error_enabling"
	ReplicationErrorDisabling = "error_disabling"
	ReplicationErrorFailover  = "error_failover"
	ReplicationErrorFailback  = "error_failback"
	ReplicationEnabled        = "enabled"
	ReplicationDisabled       = "disabled"
	ReplicationFailover       = "failed_over"
)

// volume group status
const (
	VolumeGroupCreating      = "creating"
	VolumeGroupAvailable     = "available"
	VolumeGroupErrorDeleting = "errorDeleting"
	VolumeGroupError         = "error"
	VolumeGroupDeleting      = "deleting"
	VolumeGroupUpdating      = "updating"
	VolumeGroupInUse         = "inUse"
)
