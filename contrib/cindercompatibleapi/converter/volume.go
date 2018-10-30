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
This module implements a entry into the OpenSDS northbound service.
*/

package converter

import (
	"errors"

	"github.com/opensds/opensds/pkg/model"
)

var (
	// APIVersion ...
	APIVersion = "v3"
	// Endpoint ...
	Endpoint = "http://127.0.0.1:8777/v3"
)

// *******************List accessible volumes with details*******************

// ListVolumesDetailsRespSpec ...
type ListVolumesDetailsRespSpec struct {
	Volumes []ListRespVolumeDetails `json:"volumes"`
	Count   int64                   `json:"count,omitempty"`
}

// ListRespVolumeDetails ...
type ListRespVolumeDetails struct {
	MigrationStatus     string            `json:"migration_status,omitempty"`
	Attachments         []RespAttachment  `json:"attachments"`
	Links               []Link            `json:"links,omitempty"`
	AvailabilityZone    string            `json:"availability_zone,omitempty"`
	Host                string            `json:"os-vol-host-attr:host,omitempty"`
	Encrypted           bool              `json:"encrypted,omitempty"`
	UpdatedAt           string            `json:"updated_at"`
	ReplicationStatus   string            `json:"replication_status,omitempty"`
	SnapshotID          string            `json:"snapshot_id,omitempty"`
	ID                  string            `json:"id"`
	Size                int64             `json:"size"`
	UserID              string            `json:"user_id"`
	TenantID            string            `json:"os-vol-tenant-attr:tenant_id,omitempty"`
	Migstat             string            `json:"os-vol-mig-status-attr:migstat,omitempty"`
	Metadata            map[string]string `json:"metadata"`
	Status              string            `json:"status"`
	VolumeImageMetadata map[string]string `json:"volume_image_metadata,omitempty"`
	Description         string            `json:"description"`
	Multiattach         bool              `json:"multiattach,omitempty"`
	SourceVolID         string            `json:"source_volid,omitempty"`
	ConsistencygroupID  string            `json:"consistencygroup_id,omitempty"`
	NameID              string            `json:"os-vol-mig-status-attr:name_id,omitempty"`
	Name                string            `json:"name"`
	Bootable            bool              `json:"bootable,omitempty"`
	CreatedAt           string            `json:"created_at"`
	VolumeType          string            `json:"volume_type,omitempty"`
}

// ListVolumesDetailsResp ...
func ListVolumesDetailsResp(volumes []*model.VolumeSpec) *ListVolumesDetailsRespSpec {
	var resp ListVolumesDetailsRespSpec
	var cinderVolume ListRespVolumeDetails

	if 0 == len(volumes) {
		resp.Volumes = make([]ListRespVolumeDetails, 0, 0)
	} else {
		for _, volume := range volumes {

			cinderVolume.Attachments = make([]RespAttachment, 0, 0)
			cinderVolume.AvailabilityZone = volume.AvailabilityZone
			cinderVolume.UpdatedAt = volume.BaseModel.UpdatedAt
			cinderVolume.ID = volume.BaseModel.Id
			cinderVolume.Size = volume.Size
			cinderVolume.UserID = volume.UserId
			cinderVolume.Metadata = make(map[string]string)
			//cinderVolume.TenantID = volume.TenantId
			cinderVolume.Status = volume.Status
			cinderVolume.Description = volume.Description
			cinderVolume.Name = volume.Name
			cinderVolume.CreatedAt = volume.BaseModel.CreatedAt
			cinderVolume.VolumeType = volume.ProfileId

			resp.Volumes = append(resp.Volumes, cinderVolume)
		}
	}

	return &resp
}

// *******************Create a volume*******************

// CreateVolumeReqSpec ...
type CreateVolumeReqSpec struct {
	Volume         CreateReqVolume `json:"volume"`
	SchedulerHints SchedulerHints  `json:"OS-SCH-HNT:scheduler_hints,omitempty"`
}

// CreateReqVolume ...
type CreateReqVolume struct {
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
	Volume CreateRespVolume `json:"volume,omitempty"`
}

// CreateRespVolume ...
type CreateRespVolume struct {
	MigrationStatus    string            `json:"migration_status,omitempty"`
	Attachments        []RespAttachment  `json:"attachments"`
	Links              []Link            `json:"links,omitempty"`
	AvailabilityZone   string            `json:"availability_zone,omitempty"`
	Encrypted          bool              `json:"encrypted,omitempty"`
	UpdatedAt          string            `json:"updated_at,omitempty"`
	ReplicationStatus  string            `json:"replication_status,omitempty"`
	SnapshotID         string            `json:"snapshot_id,omitempty"`
	ID                 string            `json:"id,omitempty"`
	Size               int64             `json:"size,omitempty"`
	UserID             string            `json:"user_id,omitempty"`
	Metadata           map[string]string `json:"metadata"`
	Status             string            `json:"status,omitempty"`
	Description        string            `json:"description,omitempty"`
	Multiattach        bool              `json:"multiattach,omitempty"`
	SourceVolID        string            `json:"source_volid,omitempty"`
	ConsistencygroupID string            `json:"consistencygroup_id,omitempty"`
	Name               string            `json:"name,omitempty"`
	Bootable           bool              `json:"bootable,omitempty"`
	CreatedAt          string            `json:"created_at,omitempty"`
	VolumeType         string            `json:"volume_type,omitempty"`
}

// Link ...
type Link struct {
	Href string `json:"href,omitempty"`
	Rel  string `json:"rel,omitempty"`
}

// CreateVolumeReq ...
func CreateVolumeReq(cinderReq *CreateVolumeReqSpec) (*model.VolumeSpec, error) {
	volume := model.VolumeSpec{}
	volume.BaseModel = &model.BaseModel{}
	volume.Name = cinderReq.Volume.Name
	volume.Description = cinderReq.Volume.Description
	volume.Size = cinderReq.Volume.Size
	volume.AvailabilityZone = cinderReq.Volume.AvailabilityZone
	volume.ProfileId = cinderReq.Volume.VolumeType

	if ("" != cinderReq.Volume.SourceVolID) || (false != cinderReq.Volume.Multiattach) ||
		("" != cinderReq.Volume.SnapshotID) || ("" != cinderReq.Volume.BackupID) ||
		("" != cinderReq.Volume.ImageRef) || (0 != len(cinderReq.Volume.Metadata)) ||
		("" != cinderReq.Volume.ConsistencygroupID) {

		return nil, errors.New("OpenSDS does not support the parameter: " +
			"id/source_volid/multiattach/snapshot_id/backup_id/imageRef/metadata/consistencygroup_id")
	}

	return &volume, nil
}

// CreateVolumeResp ...
func CreateVolumeResp(volume *model.VolumeSpec) *CreateVolumeRespSpec {
	resp := CreateVolumeRespSpec{}

	resp.Volume.Attachments = make([]RespAttachment, 0, 0)
	resp.Volume.AvailabilityZone = volume.AvailabilityZone
	resp.Volume.UpdatedAt = volume.BaseModel.UpdatedAt
	resp.Volume.ID = volume.BaseModel.Id
	resp.Volume.Size = volume.Size
	resp.Volume.UserID = volume.UserId
	resp.Volume.Metadata = make(map[string]string)
	resp.Volume.Status = volume.Status
	resp.Volume.Description = volume.Description
	resp.Volume.Name = volume.Name
	resp.Volume.CreatedAt = volume.BaseModel.CreatedAt
	resp.Volume.VolumeType = volume.ProfileId

	return &resp
}

// *******************List accessible volumes*******************

// ListVolumesRespSpec ...
type ListVolumesRespSpec struct {
	Volumes []ListRespVolume `json:"volumes"`
	Count   int64            `json:"count,omitempty"`
}

// ListRespVolume ...
type ListRespVolume struct {
	ID    string `json:"id"`
	Links []Link `json:"links,omitempty"`
	Name  string `json:"name"`
}

// ListVolumesResp ...
func ListVolumesResp(volumes []*model.VolumeSpec) *ListVolumesRespSpec {
	var resp ListVolumesRespSpec
	var cinderVolume ListRespVolume

	if 0 == len(volumes) {
		resp.Volumes = make([]ListRespVolume, 0, 0)
	} else {
		for _, volume := range volumes {
			cinderVolume.ID = volume.Id
			cinderVolume.Name = volume.Name

			resp.Volumes = append(resp.Volumes, cinderVolume)
		}
	}

	return &resp
}

// *******************Show a volume's details*******************

// ShowVolumeRespSpec ...
type ShowVolumeRespSpec struct {
	Volume ShowRespVolume `json:"volume"`
}

// ShowRespVolume ...
type ShowRespVolume struct {
	MigrationStatus     string            `json:"migration_status,omitempty"`
	Attachments         []RespAttachment  `json:"attachments"`
	Links               []Link            `json:"links,omitempty"`
	AvailabilityZone    string            `json:"availability_zone,omitempty"`
	Host                string            `json:"os-vol-host-attr:host,omitempty"`
	Encrypted           bool              `json:"encrypted,omitempty"`
	UpdatedAt           string            `json:"updated_at"`
	ReplicationStatus   string            `json:"replication_status,omitempty"`
	SnapshotID          string            `json:"snapshot_id"`
	ID                  string            `json:"id"`
	Size                int64             `json:"size"`
	UserID              string            `json:"user_id"`
	TenantID            string            `json:"os-vol-tenant-attr:tenant_id,omitempty"`
	Migstat             string            `json:"os-vol-mig-status-attr:migstat,omitempty"`
	Metadata            map[string]string `json:"metadata"`
	Status              string            `json:"status"`
	VolumeImageMetadata map[string]string `json:"volume_image_metadata"`
	Description         string            `json:"description"`
	Multiattach         bool              `json:"multiattach,omitempty"`
	SourceVolID         string            `json:"source_volid,omitempty"`
	ConsistencygroupID  string            `json:"consistencygroup_id,omitempty"`
	NameID              string            `json:"os-vol-mig-status-attr:name_id,omitempty"`
	Name                string            `json:"name"`
	Bootable            bool              `json:"bootable"`
	CreatedAt           string            `json:"created_at"`
	VolumeType          string            `json:"volume_type,omitempty"`
	ServiceUuID         string            `json:"service_uuid,omitempty"`
	SharedTargets       bool              `json:"shared_targets,omitempty"`
}

// ShowVolumeResp ...
func ShowVolumeResp(volume *model.VolumeSpec) *ShowVolumeRespSpec {
	resp := ShowVolumeRespSpec{}

	resp.Volume.Attachments = make([]RespAttachment, 0, 0)
	resp.Volume.AvailabilityZone = volume.AvailabilityZone
	resp.Volume.UpdatedAt = volume.BaseModel.UpdatedAt
	resp.Volume.ID = volume.BaseModel.Id
	resp.Volume.Size = volume.Size
	resp.Volume.UserID = volume.UserId
	resp.Volume.Metadata = make(map[string]string)
	resp.Volume.Status = volume.Status
	resp.Volume.Description = volume.Description
	resp.Volume.Name = volume.Name
	resp.Volume.CreatedAt = volume.BaseModel.CreatedAt
	resp.Volume.VolumeType = volume.ProfileId
	//resp.Volume.TenantID = volume.TenantId

	return &resp
}

// *******************Update a volume*******************

// UpdateVolumeReqSpec ...
type UpdateVolumeReqSpec struct {
	Volume UpdateReqVolume `json:"volume"`
}

// UpdateReqVolume ...
type UpdateReqVolume struct {
	Description string            `json:"description,omitempty"`
	Name        string            `json:"name,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// UpdateVolumeRespSpec ...
type UpdateVolumeRespSpec struct {
	Volume UpdateRespVolume `json:"volume"`
}

// UpdateRespVolume ...
type UpdateRespVolume struct {
	MigrationStatus    string            `json:"migration_status,omitempty"`
	Attachments        []RespAttachment  `json:"attachments"`
	Links              []Link            `json:"links,omitempty"`
	AvailabilityZone   string            `json:"availability_zone,omitempty"`
	Encrypted          bool              `json:"encrypted,omitempty"`
	UpdatedAt          string            `json:"updated_at"`
	ReplicationStatus  string            `json:"replication_status,omitempty"`
	SnapshotID         string            `json:"snapshot_id,omitempty"`
	ID                 string            `json:"id"`
	Size               int64             `json:"size"`
	UserID             string            `json:"user_id"`
	Metadata           map[string]string `json:"metadata"`
	Status             string            `json:"status"`
	Description        string            `json:"description"`
	Multiattach        bool              `json:"multiattach,omitempty"`
	SourceVolID        string            `json:"source_volid,omitempty"`
	ConsistencygroupID string            `json:"consistencygroup_id,omitempty"`
	Name               string            `json:"name"`
	Bootable           bool              `json:"bootable,omitempty"`
	CreatedAt          string            `json:"created_at"`
	VolumeType         string            `json:"volume_type,omitempty"`
}

// RespAttachment ...
type RespAttachment struct {
	ID string `json:"id,omitempty"`
}

// UpdateVolumeReq ...
func UpdateVolumeReq(cinderReq *UpdateVolumeReqSpec) (*model.VolumeSpec, error) {
	volume := model.VolumeSpec{}
	volume.BaseModel = &model.BaseModel{}
	volume.Description = cinderReq.Volume.Description
	volume.Name = cinderReq.Volume.Name

	if 0 != len(cinderReq.Volume.Metadata) {

		return nil, errors.New("OpenSDS does not support the parameter: metadata")
	}

	return &volume, nil
}

// UpdateVolumeResp ...
func UpdateVolumeResp(volume *model.VolumeSpec) *UpdateVolumeRespSpec {
	resp := UpdateVolumeRespSpec{}
	resp.Volume.Attachments = make([]RespAttachment, 0, 0)
	resp.Volume.AvailabilityZone = volume.AvailabilityZone
	resp.Volume.UpdatedAt = volume.BaseModel.UpdatedAt
	resp.Volume.ID = volume.BaseModel.Id
	resp.Volume.Size = volume.Size
	resp.Volume.UserID = volume.UserId
	resp.Volume.Metadata = make(map[string]string)
	resp.Volume.Status = volume.Status
	resp.Volume.Description = volume.Description
	resp.Volume.Name = volume.Name
	resp.Volume.CreatedAt = volume.BaseModel.CreatedAt
	resp.Volume.VolumeType = volume.ProfileId

	return &resp
}

// *******************Volume actions*******************

// InitializeConnectionReqSpec ...
type InitializeConnectionReqSpec struct {
	InitializeConnection InitializeConnection `json:"os-initialize_connection"`
}

// InitializeConnection ...
type InitializeConnection struct {
	Connector InitializeConnector `json:"connector"`
}

// InitializeConnector ...
type InitializeConnector struct {
	Platform      string `json:"platform"`
	Host          string `json:"host"`
	DoLocalAttach bool   `json:"do_local_attach"`
	IP            string `json:"ip"`
	OsType        string `json:"os_type"`
	Multipath     bool   `json:"multipath"`
	Initiator     string `json:"initiator"`
}

// InitializeConnectionRespSpec ...
type InitializeConnectionRespSpec struct {
	ConnectionInfo InitializeConnectionInfo `json:"connection_info"`
}

// InitializeConnectionInfo ...
type InitializeConnectionInfo struct {
	DriverVolumeType string                 `json:"driver_volume_type"`
	Data             map[string]interface{} `json:"data"`
}

// InitializeConnectionReq ...
func InitializeConnectionReq(initializeConnectionReq *InitializeConnectionReqSpec, volumeID string) *model.VolumeAttachmentSpec {
	attachment := model.VolumeAttachmentSpec{}
	attachment.Metadata = make(map[string]string)
	//attachment.Metadata["instance_uuid"] = cinderReq.Attachment.InstanceUuID
	attachment.HostInfo.Initiator = initializeConnectionReq.InitializeConnection.Connector.Initiator
	attachment.HostInfo.Ip = initializeConnectionReq.InitializeConnection.Connector.IP
	attachment.HostInfo.Platform = initializeConnectionReq.InitializeConnection.Connector.Platform
	attachment.HostInfo.Host = initializeConnectionReq.InitializeConnection.Connector.Host
	attachment.HostInfo.OsType = initializeConnectionReq.InitializeConnection.Connector.OsType
	//attachment.Mountpoint = cinderReq.Attachment.Connector.Mountpoint
	attachment.VolumeId = volumeID

	return &attachment
}

// InitializeConnectionResp ...
func InitializeConnectionResp(attachment *model.VolumeAttachmentSpec) *InitializeConnectionRespSpec {
	resp := InitializeConnectionRespSpec{}
	resp.ConnectionInfo.DriverVolumeType = attachment.ConnectionInfo.DriverVolumeType
	resp.ConnectionInfo.Data = make(map[string]interface{})

	resp.ConnectionInfo.Data["auth_password"] = attachment.ConnectionInfo.ConnectionData["authPassword"]
	resp.ConnectionInfo.Data["target_discovered"] = attachment.ConnectionInfo.ConnectionData["targetDiscovered"]
	resp.ConnectionInfo.Data["encrypted"] = attachment.ConnectionInfo.ConnectionData["encrypted"]
	//resp.ConnectionInfo.Data["qos_specs"]
	resp.ConnectionInfo.Data["target_iqn"] = attachment.ConnectionInfo.ConnectionData["targetIQN"]
	resp.ConnectionInfo.Data["target_portal"] = attachment.ConnectionInfo.ConnectionData["targetPortal"]
	resp.ConnectionInfo.Data["volume_id"] = attachment.ConnectionInfo.ConnectionData["volumeId"]
	resp.ConnectionInfo.Data["target_lun"] = attachment.ConnectionInfo.ConnectionData["targetLun"]
	resp.ConnectionInfo.Data["access_mode"] = attachment.ConnectionInfo.ConnectionData["accessMode"]
	resp.ConnectionInfo.Data["auth_username"] = attachment.ConnectionInfo.ConnectionData["authUserName"]
	resp.ConnectionInfo.Data["auth_method"] = attachment.ConnectionInfo.ConnectionData["authMethod"]

	return &resp
}
