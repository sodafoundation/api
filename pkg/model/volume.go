// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
)

// VolumeSpec is an block device created by storage service, it can be attached
// to physical machine or virtual machine instance.
type VolumeSpec struct {
	*BaseModel

	// The uuid of the project that the volume belongs to.
	TenantId string `json:"tenantId,omitempty"`

	// The uuid of the user that the volume belongs to.
	// +optional
	UserId string `json:"userId,omitempty"`

	// The name of the volume.
	Name string `json:"name,omitempty"`

	// The description of the volume.
	// +optional
	Description string `json:"description,omitempty"`

	// The group id of the volume.
	GroupId string `json:"groupId,omitempty"`

	// The size of the volume requested by the user.
	// Default unit of volume Size is GB.
	Size int64 `json:"size,omitempty"`

	// The locality that volume belongs to.
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// The status of the volume.
	// One of: "available", "error", "in-use", etc.
	Status string `json:"status,omitempty"`

	// The uuid of the pool which the volume belongs to.
	// +readOnly
	PoolId string `json:"poolId,omitempty"`

	// The uuid of the profile which the volume belongs to.
	ProfileId string `json:"profileId,omitempty"`

	// Metadata should be kept until the scemantics between opensds volume
	// and backend storage resouce description are clear.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`

	// Attach status of the volume.
	AttachStatus string
}

// VolumeAttachmentSpec is a description of volume attached resource.
type VolumeAttachmentSpec struct {
	*BaseModel

	// The uuid of the project that the volume belongs to.
	TenantId string `json:"tenantId,omitempty"`

	// The uuid of the user that the volume belongs to.
	// +optional
	UserId string `json:"userId,omitempty"`

	// The uuid of the volume which the attachment belongs to.
	VolumeId string `json:"volumeId,omitempty"`

	// The locaility when the volume was attached to a host.
	Mountpoint string `json:"mountpoint,omitempty"`

	// The status of the attachment.
	// One of: "attaching", "attached", "error", etc.
	Status string `json:"status,omitempty"`

	// Metadata should be kept until the scemantics between opensds volume
	// attachment and backend attached storage resouce description are clear.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`

	// See details in `HostInfo`
	HostInfo `json:"hostInfo,omitempty"`

	// See details in `ConnectionInfo`
	ConnectionInfo `json:"connectionInfo,omitempty"`
}

// HostInfo is a structure for all properties of host when create a volume
// attachment.
type HostInfo struct {
	Platform  string `json:"platform,omitempty"`
	OsType    string `json:"osType,omitempty"`
	Ip        string `json:"ip,omitempty"`
	Host      string `json:"host,omitempty"`
	Initiator string `json:"initiator,omitempty"`
}

// ConnectionInfo is a structure for all properties of connection when
// create a volume attachment.
type ConnectionInfo struct {
	DriverVolumeType     string                 `json:"driverVolumeType,omitempty"`
	ConnectionData       map[string]interface{} `json:"data,omitempty"`
	AdditionalProperties map[string]interface{} `json:"additionalProperties,omitempty"`
}

func (con *ConnectionInfo) EncodeConnectionData() []byte {
	conBody, _ := json.Marshal(&con.ConnectionData)
	return conBody
}

// VolumeSnapshotSpec is a description of volume snapshot resource.
type VolumeSnapshotSpec struct {
	*BaseModel

	// The uuid of the project that the volume snapshot belongs to.
	TenantId string `json:"tenantId,omitempty"`

	// The uuid of the user that the volume snapshot belongs to.
	// +optional
	UserId string `json:"userId,omitempty"`

	// The name of the volume snapshot.
	Name string `json:"name,omitempty"`

	// The description of the volume snapshot.
	// +optional
	Description string `json:"description,omitempty"`

	// The size of the volume which the snapshot belongs to.
	// Default unit of volume Size is GB.
	Size int64 `json:"size,omitempty"`

	// The status of the volume snapshot.
	// One of: "available", "error", etc.
	Status string `json:"status,omitempty"`

	// The uuid of the volume which the snapshot belongs to.
	VolumeId string `json:"volumeId,omitempty"`

	// Metadata should be kept until the scemantics between opensds volume
	// snapshot and backend storage resouce snapshot description are clear.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ExtendSpec ...
type ExtendSpec struct {
	NewSize int64 `json:"newSize,omitempty"`
}

// ExtendVolumeSpec ...
type ExtendVolumeSpec struct {
	Extend ExtendSpec `json:"extend,omitempty"`
}

type VolumeGroupSpec struct {
	*BaseModel
	// The name of the volume group.
	Name string `json:"name,omitempty"`

	Status string `json:"status,omitempty"`

	// The uuid of the project that the volume snapshot belongs to.
	TenantId string `json:"tenantId,omitempty"`

	// The uuid of the user that the volume snapshot belongs to.
	// +optional
	UserId string `json:"userId,omitempty"`

	// The description of the volume group.
	// +optional
	Description string `json:"description,omitempty"`

	// The uuid of the profile which the volume group belongs to.
	Profiles []string `json:"profileId,omitempty"`

	// The locality that volume group belongs to.
	// +optional
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// The addVolumes contain UUIDs of volumes to be added to the group.
	AddVolumes []string `json:"addVolumes,omitempty"`

	// The removeVolumes contains the volumes to be removed from the group.
	RemoveVolumes []string `json:"removeVolumes,omitempty"`

	// The uuid of the pool which the volume belongs to.
	// +readOnly
	PoolId string `json:"poolId,omitempty"`

	GroupSnapshots []string `json:"groupSnapshots,omitempty"`
}

var volume_sortKey string

type VolumeSlice []*VolumeSpec

func (volume VolumeSlice) Len() int { return len(volume) }

func (volume VolumeSlice) Swap(i, j int) { volume[i], volume[j] = volume[j], volume[i] }

func (volume VolumeSlice) Less(i, j int) bool {
	switch volume_sortKey {

	case "ID":
		return volume[i].Id < volume[j].Id
	case "NAME":
		return volume[i].Name < volume[j].Name
	case "STATUS":
		return volume[i].Status < volume[j].Status
	case "AVAILABILITYZONE":
		return volume[i].AvailabilityZone < volume[j].AvailabilityZone
	case "PROFILEID":
		return volume[i].ProfileId < volume[j].ProfileId
	case "TENANTID":
		return volume[i].TenantId < volume[j].TenantId
	case "SIZE":
		return volume[i].Size < volume[j].Size
	case "POOLID":
		return volume[i].PoolId < volume[j].PoolId
	case "DESCRIPTION":
		return volume[i].Description < volume[j].Description
		// TODO:case "lun_id" (admin_only)
		// TODO:case "GroupId"
	}
	return false
}

func (c *VolumeSpec) FindValue(k string, p *VolumeSpec) string {
	switch k {
	case "Id":
		return p.Id
	case "CreatedAt":
		return p.CreatedAt
	case "UpdatedAt":
		return p.UpdatedAt
	case "TenantId":
		return p.TenantId
	case "UserId":
		return p.UserId
	case "Name":
		return p.Name
	case "Description":
		return p.Description
	case "AvailabilityZone":
		return p.AvailabilityZone
	case "Size":
		return strconv.FormatInt(p.Size, 10)
	case "Status":
		return p.Status
	case "PoolId":
		return p.PoolId
	case "ProfileId":
		return p.ProfileId
	}
	return ""
}

func (c *VolumeSpec) SortList(volumes []*VolumeSpec, sortKey, sortDir string) []*VolumeSpec {

	volume_sortKey = sortKey

	if strings.EqualFold(sortDir, "asc") {
		sort.Sort(VolumeSlice(volumes))

	} else {
		sort.Sort(sort.Reverse(VolumeSlice(volumes)))
	}
	return volumes
}

var volumeAttachment_sortKey string

type VolumeAttachmentSlice []*VolumeAttachmentSpec

func (volumeAttachment VolumeAttachmentSlice) Len() int { return len(volumeAttachment) }

func (volumeAttachment VolumeAttachmentSlice) Swap(i, j int) {

	volumeAttachment[i], volumeAttachment[j] = volumeAttachment[j], volumeAttachment[i]
}

func (volumeAttachment VolumeAttachmentSlice) Less(i, j int) bool {
	switch volumeAttachment_sortKey {

	case "ID":
		return volumeAttachment[i].Id < volumeAttachment[j].Id
	case "VOLUMEID":
		return volumeAttachment[i].VolumeId < volumeAttachment[j].VolumeId
	case "STATUS":
		return volumeAttachment[i].Status < volumeAttachment[j].Status
	case "USERID":
		return volumeAttachment[i].UserId < volumeAttachment[j].UserId
	case "TENANTID":
		return volumeAttachment[i].TenantId < volumeAttachment[j].TenantId
	}
	return false
}

func (c *VolumeAttachmentSpec) FindValue(k string, p *VolumeAttachmentSpec) string {
	switch k {
	case "Id":
		return p.Id
	case "CreatedAt":
		return p.CreatedAt
	case "UpdatedAte":
		return p.UpdatedAt
	case "TenantId":
		return p.TenantId
	case "UserId":
		return p.UserId
	case "VolumeId":
		return p.VolumeId
	case "Mountpoint":
		return p.Mountpoint
	case "Status":
		return p.Status
	}
	return ""
}

func (c *VolumeAttachmentSpec) SortList(attachments []*VolumeAttachmentSpec, sortKey, sortDir string) []*VolumeAttachmentSpec {

	volumeAttachment_sortKey = sortKey

	if strings.EqualFold(sortDir, "asc") {
		sort.Sort(VolumeAttachmentSlice(attachments))
	} else {
		sort.Sort(sort.Reverse(VolumeAttachmentSlice(attachments)))
	}
	return attachments
}

var volumeSnapshot_sortKey string

type VolumeSnapshotSlice []*VolumeSnapshotSpec

func (volumeSnapshot VolumeSnapshotSlice) Len() int { return len(volumeSnapshot) }

func (volumeSnapshot VolumeSnapshotSlice) Swap(i, j int) {

	volumeSnapshot[i], volumeSnapshot[j] = volumeSnapshot[j], volumeSnapshot[i]
}

func (volumeSnapshot VolumeSnapshotSlice) Less(i, j int) bool {
	switch volumeSnapshot_sortKey {

	case "ID":
		return volumeSnapshot[i].Id < volumeSnapshot[j].Id
	case "VOLUMEID":
		return volumeSnapshot[i].VolumeId < volumeSnapshot[j].VolumeId
	case "STATUS":
		return volumeSnapshot[i].Status < volumeSnapshot[j].Status
	case "USERID":
		return volumeSnapshot[i].UserId < volumeSnapshot[j].UserId
	case "TENANTID":
		return volumeSnapshot[i].TenantId < volumeSnapshot[j].TenantId
	case "SIZE":
		return volumeSnapshot[i].Size < volumeSnapshot[j].Size
		//TODO:case "GroupSnapshotId"
	}
	return false
}

func (c *VolumeSnapshotSpec) FindValue(k string, p *VolumeSnapshotSpec) string {
	switch k {
	case "Id":
		return p.Id
	case "CreatedAt":
		return p.CreatedAt
	case "UpdatedAte":
		return p.UpdatedAt
	case "TenantId":
		return p.TenantId
	case "UserId":
		return p.UserId
	case "Name":
		return p.Name
	case "Description":
		return p.Description
	case "Status":
		return p.Status
	case "Size":
		return strconv.FormatInt(p.Size, 10)
	case "VolumeId":
		return p.VolumeId
	}
	return ""
}

func (c *VolumeSnapshotSpec) SortList(snapshots []*VolumeSnapshotSpec, sortKey, sortDir string) []*VolumeSnapshotSpec {

	volumeSnapshot_sortKey = sortKey

	if strings.EqualFold(sortDir, "asc") {
		sort.Sort(VolumeSnapshotSlice(snapshots))
	} else {
		sort.Sort(sort.Reverse(VolumeSnapshotSlice(snapshots)))
	}
	return snapshots
}

var volumeGroup_sortKey string

type VolumeGroupSlice []*VolumeGroupSpec

func (volumeGroup VolumeGroupSlice) Len() int { return len(volumeGroup) }

func (volumeGroup VolumeGroupSlice) Swap(i, j int) {

	volumeGroup[i], volumeGroup[j] = volumeGroup[j], volumeGroup[i]
}

func (volumeGroup VolumeGroupSlice) Less(i, j int) bool {
	switch volumeGroup_sortKey {

	case "ID":
		return volumeGroup[i].Id < volumeGroup[j].Id
	case "CREATEDAT":
		return volumeGroup[i].CreatedAt < volumeGroup[j].CreatedAt
	case "NAME":
		return volumeGroup[i].Name < volumeGroup[j].Name
	case "USERID":
		return volumeGroup[i].UserId < volumeGroup[j].UserId
	case "TENANTID":
		return volumeGroup[i].TenantId < volumeGroup[j].TenantId
	case "STATUS":
		return volumeGroup[i].Status < volumeGroup[j].Status
	case "POOLID":
		return volumeGroup[i].PoolId < volumeGroup[j].PoolId
	case "AVAILABILITYZONE":
		return volumeGroup[i].AvailabilityZone < volumeGroup[i].AvailabilityZone
	}
	return false
}

func (c *VolumeGroupSpec) FindValue(k string, p *VolumeGroupSpec) string {
	switch k {
	case "Id":
		return p.Id
	case "CreatedAt":
		return p.CreatedAt
	case "UpdatedAte":
		return p.UpdatedAt
	case "TenantId":
		return p.TenantId
	case "UserId":
		return p.UserId
	case "Name":
		return p.Name
	case "Description":
		return p.Description
	case "Status":
		return p.Status
	case "AvailabilityZone":
		return p.AvailabilityZone
	case "PoolId":
		return p.PoolId
	}
	return ""
}

func (c *VolumeGroupSpec) SortList(groups []*VolumeGroupSpec, sortKey, sortDir string) []*VolumeGroupSpec {

	volumeGroup_sortKey = sortKey

	if strings.EqualFold(sortDir, "asc") {
		sort.Sort(VolumeGroupSlice(groups))
	} else {
		sort.Sort(sort.Reverse(VolumeGroupSlice(groups)))
	}
	return groups

}
