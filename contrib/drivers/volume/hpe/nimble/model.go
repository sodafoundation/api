// Copyright 2019 The OpenSDS Authors.
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

package nimble

type AuthRespBody struct {
	Data AuthRespData `json:"data"`
}

type AuthRespData struct {
	AppName      string `json:"app_name"`
	CreationTime int    `json:"creation_time"`
	Id           int    `json:"id"`
	LastModified int    `json:"last_modified"`
	SessionToken string `json:"session_token"`
	SourceIp     string `json:"source_ip"`
	UserName     string `json:"username"`
}

type ClinetErrors struct {
	Errs []error
}

type ArrayInnerErrorBody struct {
	Errs []ArrayInnerErrorResp `json:"messages"`
}

type ArrayInnerErrorResp struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Text     string `json:"text"`
}

type StoragePoolsRespBody struct {
	StartRow  int                   `json:"startRow"`
	EndRow    int                   `json:"endRow"`
	TotalRows int                   `json:"totalRows"`
	Data      []StoragePoolRespData `json:"data"`
}

type StoragePoolRespData struct {
	Id            string      `json:"id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	TotalCapacity int64       `json:"capacity"`
	FreeCapacity  int64       `json:"free_space"`
	ArrayList     []ArrayList `json:"array_list"`
	Endpoint      string
	Token         string
}

type ArrayList struct {
	Id                       string `json:"id"`
	ArrayId                  string `json:"array_id"`
	Name                     string `json:"name"`
	ArrayName                string `json:"array_name"`
	Usage                    int64  `json:"usage"`
	UsageValid               bool   `json:"usage_valid"`
	Migrate                  string `json:"migrate"`
	EvacUsage                int64  `json:"evac_usage"`
	EvacTime                 int64  `json:"evac_time"`
	SnapUsageCompressedBytes int64  `json:"snap_usage_compressed_bytes"`
	UsableCapacity           int64  `json:"usable_capacity"`
	VolUsageCompressedBytes  int64  `json:"vol_usage_compressed_bytes"`
}

type AuthOptions struct {
	Username  string `yaml:"username,omitempty"`
	Password  string `yaml:"password,omitempty"`
	Endpoints string `yaml:"endpoints,omitempty"`
	Insecure  bool   `yaml:"insecure,omitempty"`
}

type NimbleClient struct {
	user      string
	passwd    string
	endpoints []string
	tokens    []string
	insecure  bool
}

type VolumeRespBody struct {
	Data VolumeRespData `json:"data"`
}

type AllVolumeRespBody struct {
	StartRow  int              `json:"startRow"`
	EndRow    int              `json:"endRow"`
	TotalRows int              `json:"totalRows"`
	Data      []VolumeRespData `json:"data"`
}

type VolumeRespData struct {
	AgentType                  string      `json:"agent_type"`
	AppCategory                string      `json:"app_category"`
	AppUuid                    string      `json:"app_uuid"`
	AvgStatsLast5mins          interface{} `json:"avg_stats_last_5mins"`
	BaseSnapId                 string      `json:"base_snap_id"`
	BaseSnapName               string      `json:"base_snap_name"`
	BlockSize                  int64       `json:"block_size"`
	CacheNeededForPin          int64       `json:"cache_needed_for_pin"`
	CachePinned                bool        `json:"cache_pinned"`
	CachePolicy                string      `json:"cache_policy"`
	CachingEnabled             bool        `json:"caching_enabled"`
	CksumLastVerified          int64       `json:"cksum_last_verified"`
	Clone                      bool        `json:"clone"`
	ContentReplErrorsFound     bool        `json:"content_repl_errors_found"`
	CreationTime               int64       `json:"creation_time"`
	DedupeEnabled              bool        `json:"dedupe_enabled"`
	Description                string      `json:"description"`
	DestPoolId                 string      `json:"dest_pool_id"`
	DestPoolName               string      `json:"dest_pool_name"`
	EncryptionCipher           string      `json:"encryption_cipher"`
	FolderId                   string      `json:"folder_id"`
	FolderName                 string      `json:"folder_name"`
	FullName                   string      `json:"full_name"`
	Id                         string      `json:"id"`
	LastContentSnapBrCgUid     int64       `json:"last_content_snap_br_cg_uid"`
	LastContentSnapBrGid       int64       `json:"last_content_snap_br_gid"`
	LastContentSnapId          int64       `json:"last_content_snap_id"`
	LastModified               int64       `json:"last_modified"`
	LastReplicatedSnap         interface{} `json:"last_replicated_snap"`
	LastSnap                   interface{} `json:"last_snap"`
	Limit                      int64       `json:"limit"`
	LimitIops                  int64       `json:"limit_iops"`
	LimitMbps                  int64       `json:"limit_mbps"`
	Metadata                   interface{} `json:"metadata"`
	MoveAborting               bool        `json:"move_aborting"`
	MoveBytesMigrated          int64       `json:"move_bytes_migrated"`
	MoveBytesRemaining         int64       `json:"move_bytes_remaining"`
	MoveEstComplTime           int64       `json:"move_est_compl_time"`
	MoveStartTime              int64       `json:"move_start_time"`
	MultiInitiator             bool        `json:"multi_initiator"`
	Name                       string      `json:"name"`
	NeedsContentRepl           bool        `json:"needs_content_repl"`
	NumConnections             int64       `json:"num_connections"`
	NumFcConnections           int64       `json:"num_fc_connections"`
	NumIscsiConnections        int64       `json:"num_iscsi_connections"`
	NumSnaps                   int64       `json:"num_snaps"`
	OfflineReason              interface{} `json:"offline_reason"`
	Online                     bool        `json:"online"`
	OnlineSnaps                interface{} `json:"online_snaps"`
	OwnedByGroup               string      `json:"owned_by_group"`
	OwnedByGroupId             string      `json:"owned_by_group_id"`
	ParentVolId                string      `json:"parent_vol_id"`
	ParentVolName              string      `json:"parent_vol_name"`
	PerfpolicyId               string      `json:"perfpolicy_id"`
	PerfpolicyName             string      `json:"perfpolicy_name"`
	PinnedCacheSize            int64       `json:"pinned_cache_size"`
	PoolId                     int64       `json:"pool_id"`
	PoolName                   string      `json:"pool_name"`
	PreviouslyDeduped          bool        `json:"previously_deduped"`
	ProjectedNumSnaps          int64       `json:"projected_num_snaps"`
	ProtectionType             string      `json:"protection_type"`
	ReadOnly                   bool        `json:"read_only"`
	Reserve                    int64       `json:"reserve"`
	SearchName                 string      `json:"search_name"`
	SerialNumber               string      `json:"serial_number"`
	Size                       int64       `json:"size"`
	SnapLimit                  int64       `json:"snap_limit"`
	SnapLimitPercent           int64       `json:"snap_limit_percent"`
	SnapReserve                int64       `json:"snap_reserve"`
	SnapUsageCompressedBytes   int64       `json:"snap_usage_compressed_bytes"`
	SnapUsagePopulatedBytes    int64       `json:"snap_usage_populated_bytes"`
	SnapUsageUncompressedBytes int64       `json:"snap_usage_uncompressed_bytes"`
	SnapWarnLevel              int64       `json:"snap_warn_level"`
	SpaceUsageLevel            string      `json:"space_usage_level"`
	TargetName                 string      `json:"target_name"`
	ThinlyProvisioned          bool        `json:"thinly_provisioned"`
	TotalUsageBytes            int64       `json:"total_usage_bytes"`
	UpstreamCachePinned        bool        `json:"upstream_cache_pinned"`
	UsageValid                 bool        `json:"usage_valid"`
	VolState                   string      `json:"vol_state"`
	VolUsageCompressedBytese   int64       `json:"vol_usage_compressed_bytes"`
	VolUsageUncompressedBytes  int64       `json:"vol_usage_uncompressed_bytes"`
	VolcollId                  string      `json:"volcoll_id"`
	VolcollName                string      `json:"volcoll_name"`
	VpdIeee0                   string      `json:"vpd_ieee0"`
	VpdIeee1                   string      `json:"vpd_ieee1"`
	VpdT10                     string      `json:"vpd_t10"`
	WarnLevel                  int64       `json:"warn_level"`
	IscsiSessions              interface{} `json:"iscsi_sessions"`
	FcSessions                 interface{} `json:"fc_sessions"`
	AccessControlRecords       interface{} `json:"access_control_records"`
}

type SnapshotRespBody struct {
	Data SnapshotRespData `json:"data"`
}

type AllSnapshotRespBody struct {
	StartRow  int                `json:"startRow"`
	EndRow    int                `json:"endRow"`
	TotalRows int                `json:"totalRows"`
	Data      []SnapshotRespData `json:"data"`
}

type SnapshotRespData struct {
	AccessControlRecords     interface{} `json:"access_control_records"`
	AgentType                string      `json:"agent_type"`
	AppUuid                  string      `json:"app_uuid"`
	CreationTime             int64       `json:"creation_time"`
	Description              string      `json:"description"`
	Id                       string      `json:"id"`
	IsReplica                bool        `json:"is_replica"`
	IsUnmanaged              bool        `json:"is_unmanaged"`
	LastModified             int64       `json:"last_modified"`
	Metadata                 interface{} `json:"metadata"`
	Name                     string      `json:"name"`
	NewDataCompressedBytes   int64       `json:"new_data_compressed_bytes"`
	NewDataUncompressedBytes string      `json:"new_data_uncompressed_bytes"`
	NewDataValid             bool        `json:"new_data_valid"`
	OfflineReason            string      `json:"offline_reason"`
	Online                   bool        `json:"online"`
	OriginName               string      `json:"origin_name"`
	ReplicationStatus        interface{} `json:"replication_status"`
	ScheduleId               string      `json:"schedule_id"`
	ScheduleName             string      `json:"schedule_name"`
	SerialNumber             string      `json:"serial_number"`
	Size                     int64       `json:"size"`
	SnapCollectionId         string      `json:"snap_collection_id"`
	SnapCollectionName       string      `json:"snap_collection_name"`
	TargetName               string      `json:"target_name"`
	VolId                    string      `json:"vol_id"`
	VolName                  string      `json:"vol_name"`
	VpdIeee0                 string      `json:"vpd_ieee0"`
	VpdIeee1                 string      `json:"vpd_ieee1"`
	VpdT10                   string      `json:"vpd_t10"`
	Writable                 bool        `json:"writable"`
}

type LoginReqBody struct {
	Data LoginReqData `json:"data"`
}
type LoginReqData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateVolumeReqBody struct {
	Data CreateVolumeReqData `json:"data"`
}

type CreateVolumeReqData struct {
	Name             string      `json:"name"`                  //not meta
	Size             int64       `json:"size"`                  //not meta
	Description      string      `json:"description,omitempty"` //not meta
	PerfpolicyId     string      `json:"perfpolicy_id,omitempty"`
	Reserve          int64       `json:"reserve,omitempty"`
	WarnLevel        int64       `json:"warn_level,omitempty"`
	Limit            int64       `json:"limit,omitempty"`
	SnapReserve      int64       `json:"snap_reserve,omitempty"` //not meta
	SnapWarnLevel    int64       `json:"snap_warn_level,omitempty"`
	SnapLimit        int64       `json:"snap_limit,omitempty"`
	SnapLimitPercent int64       `json:"snap_limit_percent,omitempty"`
	Online           *bool       `json:"online,omitempty"`
	OwnedByGroupId   string      `json:"owned_by_group_id,omitempty"`
	MultiInitiator   bool        `json:"multi_initiator,omitempty"`
	PoolId           string      `json:"pool_id"` //not meta
	ReadOnly         bool        `json:"read_only,omitempty"`
	BlockSize        int64       `json:"block_size,omitempty"`
	Clone            bool        `json:"clone,omitempty"`
	BaseSnapId       string      `json:"base_snap_id,omitempty"` //not meta
	AgentType        string      `json:"agent_type,omitempty"`
	DestPoolId       string      `json:"dest_pool_id,omitempty"`
	CachePinned      *bool       `json:"cache_pinned,omitempty"`
	EncryptionCipher string      `json:"encryption_cipher,omitempty"`
	AppUuid          string      `json:"app_uuid,omitempty"`
	FolderId         string      `json:"folder_id,omitempty"`
	Metadata         interface{} `json:"metadata,omitempty"`
	DedupeEnabled    *bool       `json:"dedupe_enabled,omitempty"`
	LimitIops        int64       `json:"limit_iops,omitempty"`
	LimitMbps        int64       `json:"limit_mbps,omitempty"`
}

type ExtendVolumeReqBody struct {
	Data ExtendVolumeReqData `json:"data"`
}
type ExtendVolumeReqData struct {
	Name             string      `json:"name,omitempty"`
	Size             int64       `json:"size"` //not meta
	Description      string      `json:"description,omitempty"`
	PerfpolicyId     string      `json:"perfpolicy_id,omitempty"`
	Reserve          int64       `json:"reserve,omitempty"`
	WarnLevel        int64       `json:"warn_level,omitempty"`
	Limit            int64       `json:"limit,omitempty"`
	SnapReserve      int64       `json:"snap_reserve,omitempty"` //not meta
	SnapWarnLevel    int64       `json:"snap_warn_level,omitempty"`
	SnapLimit        int64       `json:"snap_limit,omitempty"`
	SnapLimitPercent int64       `json:"snap_limit_percent,omitempty"`
	Online           *bool       `json:"online,omitempty"`
	OwnedByGroupId   string      `json:"owned_by_group_id,omitempty"`
	MultiInitiator   bool        `json:"multi_initiator,omitempty"`
	ReadOnly         bool        `json:"read_only,omitempty"`
	BlockSize        int64       `json:"block_size,omitempty"`
	VolcollId        string      `json:"volcoll_id,omitempty"`
	AgentType        string      `json:"agent_type,omitempty"`
	Force            *bool       `json:"force,omitempty"`
	CachePinned      *bool       `json:"cache_pinned,omitempty"`
	AppUuid          string      `json:"app_uuid,omitempty"`
	FolderId         string      `json:"folder_id,omitempty"`
	Metadata         interface{} `json:"metadata,omitempty"`
	CachingEnabled   *bool       `json:"caching_enabled,omitempty"`
	DedupeEnabled    *bool       `json:"dedupe_enabled,omitempty"`
	LimitIops        int64       `json:"limit_iops,omitempty"`
	LimitMbps        int64       `json:"limit_mbps,omitempty"`
}
type ExtendVolumeRespBody struct {
	Data VolumeRespData `json:"data"`
}

type CreateSnapshotReqBody struct {
	Data CreateSnapshotReqData `json:"data"`
}

type CreateSnapshotReqData struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	VolId       string      `json:"vol_id"`
	Online      *bool       `json:"online,omitempty"`
	Writable    *bool       `json:"writable,omitempty"`
	AppUuid     string      `json:"app_uuid,omitempty"`
	Metadata    interface{} `json:"metadata,omitempty"`
	AgentType   string      `json:"agent_type,omitempty"`
}

type OfflineVolumeReqBody struct {
	Data OfflineVolumeReqData `json:"data"`
}

type OfflineVolumeReqData struct {
	Online bool `json:"online"`
	Force  bool `json:"force"`
}

type AllInitiatorRespBody struct {
	StartRow  int                 `json:"startRow"`
	EndRow    int                 `json:"endRow"`
	TotalRows int                 `json:"totalRows"`
	Data      []InitiatorRespData `json:"data"`
}

type InitiatorRespBody struct {
	Data InitiatorRespData `json:"data"`
}

type InitiatorRespData struct {
	Id                 string `json:"id"`
	AccessProtocol     string `json:"access_protocol"`
	InitiatorGroupId   string `json:"initiator_group_id"`
	InitiatorGroupName string `json:"initiator_group_name"`
	Label              string `json:"label"`
	Iqn                string `json:"iqn"`
	IpAddress          string `json:"ip_address"`
	Alias              string `json:"alias"`
	Wwpn               string `json:"wwpn"`
	CreationTime       int64  `json:"creation_time"`
	LastModified       int64  `json:"last_modified"`
}

type CreateInitiatorReqBody struct {
	Data CreateInitiatorReqData `json:"data"`
}

type CreateInitiatorReqData struct {
	AccessProtocol   string `json:"access_protocol,omitempty"`
	InitiatorGroupId string `json:"initiator_group_id,omitempty"`
	Label            string `json:"label,omitempty"`
	Iqn              string `json:"iqn,omitempty"`
	IpAddress        string `json:"ip_address,omitempty"`
	Alias            string `json:"alias,omitempty"`
	Wwpn             string `json:"wwpn,omitempty"`
}

type AllInitiatorGrpRespBody struct {
	StartRow  int                    `json:"startRow"`
	EndRow    int                    `json:"endRow"`
	TotalRows int                    `json:"totalRows"`
	Data      []InitiatorGrpRespData `json:"data"`
}

type InitiatorGrpRespBody struct {
	Data InitiatorGrpRespData `json:"data"`
}

type InitiatorGrpRespData struct {
	Id              string              `json:"id"`
	Name            string              `json:"name"`
	FullName        string              `json:"full_name"`
	SearchName      string              `json:"search_name"`
	Description     string              `json:"description"`
	AccessProtocol  string              `json:"access_protocol"`
	HostType        string              `json:"host_type"`
	TargetSubnets   []map[string]string `json:"target_subnets"`
	IscsiInitiators []map[string]string `json:"iscsi_initiators"`
	FcInitiators    []map[string]string `json:"fc_initiators"`
	CreationTime    int64               `json:"creation_time"`
	LastModified    int64               `json:"last_modified"`
	AppUuid         string              `json:"app_uuid"`
	VolumeCount     string              `json:"volume_count"`
	VolumeList      []map[string]string `json:"volume_list"`
	NumConnections  int64               `json:"num_connections"`
}

type CreateInitiatorGrpReqBody struct {
	Data CreateInitiatorGrpReqData `json:"data"`
}

type CreateInitiatorGrpReqData struct {
	Name            string              `json:"name,omitempty"`
	Description     string              `json:"description,omitempty"`
	AccessProtocol  string              `json:"access_protocol,omitempty"`
	HostType        string              `json:"host_type,omitempty"`
	TargetSubnets   []map[string]string `json:"target_subnets,omitempty"`
	IscsiInitiators []map[string]string `json:"iscsi_initiators,omitempty"`
	FcInitiators    []map[string]string `json:"fc_initiators,omitempty"`
	AppUuid         string              `json:"app_uuid,omitempty"`
}

type CreateAccessControlReqBody struct {
	Data CreateAccessControlReqData `json:"data"`
}

type CreateAccessControlReqData struct {
	ApplyTo          string   `json:"apply_to,omitempty"`
	ChapUserId       string   `json:"chap_user_id,omitempty"`
	InitiatorGroupId string   `json:"initiator_group_id,omitempty"`
	Lun              string   `json:"lun,omitempty"`
	VolId            string   `json:"vol_id,omitempty"`
	PeId             string   `json:"pe_id,omitempty"`
	SnapId           string   `json:"snap_id,omitempty"`
	PeIds            []string `json:"pe_ids,omitempty"`
}

type AllAccessControlRespBody struct {
	StartRow  int                     `json:"startRow"`
	EndRow    int                     `json:"endRow"`
	TotalRows int                     `json:"totalRows"`
	Data      []AccessControlRespData `json:"data"`
}

type AccessControlRespBody struct {
	Data AccessControlRespData `json:"data"`
}

type AccessControlRespData struct {
	Id                 string   `json:"id"`
	ApplyTo            string   `json:"apply_to"`
	ChapUserId         string   `json:"chap_user_id"`
	ChapUserName       string   `json:"chap_user_name"`
	InitiatorGroupId   string   `json:"initiator_group_id"`
	InitiatorGroupName string   `json:"initiator_group_name"`
	Lun                int64    `json:"lun"`
	VolId              string   `json:"vol_id"`
	VolName            string   `json:"vol_name"`
	VolAgentType       string   `json:"vol_agent_type"`
	PeId               string   `json:"pe_id"`
	PeName             string   `json:"pe_name"`
	PeLun              string   `json:"pe_lun"`
	SnapId             string   `json:"snap_id"`
	SnapName           string   `json:"snap_name"`
	PeIds              []string `json:"pe_ids"`
	SnapLuns           []string `json:"snapluns"`
	CreationTime       int64    `json:"creation_time"`
	LastModified       int64    `json:"last_modified"`
	AccessProtocol     string   `json:"access_protocol"`
}
