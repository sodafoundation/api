// Copyright 2017 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package dorado

type Error struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type GenericResult struct {
	Data  interface{} `json:"data"`
	Error Error       `json:"error"`
}

type Auth struct {
	AccountState  int    `json:"accountstate"`
	DeviceId      string `json:"deviceid"`
	IBaseToken    string `json:"iBaseToken"`
	LastLoginIp   string `json:"lastloginip"`
	LastLoginTime int    `json:"lastlogintime"`
	Level         int    `json:"level"`
	PwdChanGeTime int    `json:"pwdchangetime"`
	UserGroup     string `json:"usergroup"`
	UserId        string `json:"userid"`
	UserName      string `json:"username"`
	UserScope     string `json:"userscope"`
}
type AuthResp struct {
	Data  Auth  `json:"data"`
	Error Error `json:"error"`
}

type StoragePool struct {
	Description       string `json:"DESCRIPTION"`
	Id                string `json:"ID"`
	Name              string `json:"NAME"`
	UserFreeCapacity  string `json:"USERFREECAPACITY"`
	UserTotalCapacity string `json:"USERTOTALCAPACITY"`
}
type StoragePoolsResp struct {
	Data  []StoragePool `json:"data"`
	Error Error         `json:"error"`
}

type Lun struct {
	AllocCapacity               string `json:"ALLOCCAPACITY"`
	AllocType                   string `json:"ALLOCTYPE"`
	Capability                  string `json:"CAPABILITY"`
	Capacity                    string `json:"CAPACITY"`
	CapacityAlarmLevel          string `json:"CAPACITYALARMLEVEL"`
	Description                 string `json:"DESCRIPTION"`
	DrsEnable                   string `json:"DRS_ENABLE"`
	EnableCompression           string `json:"ENABLECOMPRESSION"`
	EnableIscsiThinLunThreshold string `json:"ENABLEISCSITHINLUNTHRESHOLD"`
	EnableSmartDedup            string `json:"ENABLESMARTDEDUP"`
	ExposedToInitiator          string `json:"EXPOSEDTOINITIATOR"`
	ExtendIfSwitch              string `json:"EXTENDIFSWITCH"`
	HealthStatus                string `json:"HEALTHSTATUS"`
	Id                          string `json:"ID"`
	IsAdd2LunGroup              string `json:"ISADD2LUNGROUP"`
	IsCheckZeroPage             string `json:"ISCHECKZEROPAGE"`
	IscsiThinLunThreshold       string `json:"ISCSITHINLUNTHRESHOLD"`
	LunMigrationOrigin          string `json:"LUNMigrationOrigin"`
	MirrorPolicy                string `json:"MIRRORPOLICY"`
	MirrorType                  string `json:"MIRRORTYPE"`
	Name                        string `json:"NAME"`
	OwningController            string `json:"OWNINGCONTROLLER"`
	ParentId                    string `json:"PARENTID"`
	ParentName                  string `json:"PARENTNAME"`
	PrefetChPolicy              string `json:"PREFETCHPOLICY"`
	PrefetChValue               string `json:"PREFETCHVALUE"`
	RemoteLunId                 string `json:"REMOTELUNID"`
	RemoteReplicationIds        string `json:"REMOTEREPLICATIONIDS"`
	ReplicationCapacity         string `json:"REPLICATION_CAPACITY"`
	RunningStatus               string `json:"RUNNINGSTATUS"`
	RunningWritePolicy          string `json:"RUNNINGWRITEPOLICY"`
	SectorSize                  string `json:"SECTORSIZE"`
	SnapShotIds                 string `json:"SNAPSHOTIDS"`
	SubType                     string `json:"SUBTYPE"`
	ThinCapacityUsage           string `json:"THINCAPACITYUSAGE"`
	Type                        int    `json:"TYPE"`
	UsageType                   string `json:"USAGETYPE"`
	WorkingController           string `json:"WORKINGCONTROLLER"`
	WritePolicy                 string `json:"WRITEPOLICY"`
	Wwn                         string `json:"WWN"`
	RemoteLunWwn                string `json:"remoteLunWwn"`
}

type LunResp struct {
	Data  Lun   `json:"data"`
	Error Error `json:"error"`
}

type LunsResp struct {
	Data  []Lun `json:"data"`
	Error Error `json:"error"`
}

type Snapshot struct {
	CascadedLevel         string `json:"CASCADEDLEVEL"`
	CascadedNum           string `json:"CASCADEDNUM"`
	ConsumedCapacity      string `json:"CONSUMEDCAPACITY"`
	Description           string `json:"DESCRIPTION"`
	ExposedToInitiator    string `json:"EXPOSEDTOINITIATOR"`
	HealthStatus          string `json:"HEALTHSTATUS"`
	Id                    string `json:"ID"`
	IoClassId             string `json:"IOCLASSID"`
	IoPriority            string `json:"IOPRIORITY"`
	SourceLunCapacity     string `json:"SOURCELUNCAPACITY"`
	Name                  string `json:"NAME"`
	ParentId              string `json:"PARENTID"`
	ParentName            string `json:"PARENTNAME"`
	ParentType            int    `json:"PARENTTYPE"`
	RollBackendTime       string `json:"ROLLBACKENDTIME"`
	RollbackRate          string `json:"ROLLBACKRATE"`
	RollbackSpeed         string `json:"ROLLBACKSPEED"`
	RollbackStartTime     string `json:"ROLLBACKSTARTTIME"`
	RollbackTargetObjId   string `json:"ROLLBACKTARGETOBJID"`
	RollbackTargetObjName string `json:"ROLLBACKTARGETOBJNAME"`
	RunningStatus         string `json:"RUNNINGSTATUS"`
	SourceLunId           string `json:"SOURCELUNID"`
	SourceLunName         string `json:"SOURCELUNNAME"`
	SubType               string `json:"SUBTYPE"`
	TimeStamp             string `json:"TIMESTAMP"`
	Type                  int    `json:"TYPE"`
	UserCapacity          string `json:"USERCAPACITY"`
	WorkingController     string `json:"WORKINGCONTROLLER"`
	Wwn                   string `json:"WWN"`
	ReplicationCapacity   string `json:"replicationCapacity"`
}

type SnapshotResp struct {
	Data  Snapshot `json:"data"`
	Error Error    `json:"error"`
}

type SnapshotsResp struct {
	Data  []Snapshot `json:"data"`
	Error Error      `json:"error"`
}

type Initiator struct {
	Id         string `json:"ID"`
	Name       string `json:"NAME"`
	ParentId   string `json:"PARENTID"`
	ParentType string `json:"PARENTTYPE"`
	ParentName string `json:"PARENTNAME"`
}

type InitiatorResp struct {
	Data  Initiator `json:"data"`
	Error Error     `json:"error"`
}

type InitiatorsResp struct {
	Data  []Initiator `json:"data"`
	Error Error       `json:"error"`
}

type Host struct {
	Id               string `json:"ID"`
	Name             string `json:"NAME"`
	OsType           string `json:"OPERATIONSYSTEM"`
	Ip               string `json:"IP"`
	IsAddToHostGroup bool   `json:"ISADD2HOSTGROUP"`
}

type HostResp struct {
	Data  Host  `json:"data"`
	Error Error `json:"error"`
}

type HostsResp struct {
	Data  []Host `json:"data"`
	Error Error  `json:"error"`
}

type HostGroup struct {
	Id                string `json:"ID"`
	Name              string `json:"NAME"`
	Description       string `json:"DESCRIPTION"`
	IsAdd2MappingView string `json:"ISADD2MAPPINGVIEW"`
}

type HostGroupResp struct {
	Data  HostGroup `json:"data"`
	Error Error     `json:"error"`
}

type HostGroupsResp struct {
	Data  []HostGroup `json:"data"`
	Error Error       `json:"error"`
}

type LunGroup struct {
	Id                string `json:"ID"`
	Name              string `json:"NAME"`
	Description       string `json:"DESCRIPTION"`
	IsAdd2MappingView string `json:"ISADD2MAPPINGVIEW"`
}

type LunGroupResp struct {
	Data  LunGroup `json:"data"`
	Error Error    `json:"error"`
}

type LunGroupsResp struct {
	Data  []LunGroup `json:"data"`
	Error Error      `json:"error"`
}

type MappingView struct {
	Id          string `json:"ID"`
	Name        string `json:"NAME"`
	Description string `json:"DESCRIPTION"`
}

type MappingViewResp struct {
	Data  MappingView `json:"data"`
	Error Error       `json:"error"`
}

type MappingViewsResp struct {
	Data  []MappingView `json:"data"`
	Error Error         `json:"error"`
}

type IscsiTgtPort struct {
	EthPortId string `json:"ETHPORTID"`
	Id        string `json:"ID"`
	Tpgt      string `json:"TPGT"`
	Type      int    `json:"TYPE"`
}

type IscsiTgtPortsResp struct {
	Data  []IscsiTgtPort `json:"data"`
	Error Error          `json:"error"`
}

type HostAssociateLun struct {
	Id                string `json:"ID"`
	AssociateMetadata string `json:"ASSOCIATEMETADATA"`
}

type HostAssociateLunsResp struct {
	Data  []HostAssociateLun `json:"data"`
	Error Error              `json:"error"`
}

type System struct {
	Id          string `json:"ID"`
	Name        string `json:"NAME"`
	Location    string `json:"LOCATION"`
	ProductMode string `json:"PRODUCTMODE"`
	Wwn         string `json:"wwn"`
}

type SystemResp struct {
	Data  System `json:"data"`
	Error Error  `json:"error"`
}

type RemoteDevice struct {
	Id            string `json:"ID"`
	Name          string `json:"NAME"`
	ArrayType     string `json:"ARRAYTYPE"`
	HealthStatus  string `json:"HEALTHSTATUS"`
	RunningStatus string `json:"RUNNINGSTATUS"`
	Wwn           string `json:"WWN"`
}

type RemoteDevicesResp struct {
	Data  []RemoteDevice `json:"data"`
	Error Error          `json:"error"`
}

type ReplicationPair struct {
	Capacity            string `json:"CAPACITY"`
	CompressValid       string `json:"COMPRESSVALID"`
	EnableCompress      string `json:"ENABLECOMPRESS"`
	HealthStatus        string `json:"HEALTHSTATUS"`
	Id                  string `json:"ID"`
	IsDataSync          string `json:"ISDATASYNC"`
	IsInCg              string `json:"ISINCG"`
	IsPrimary           string `json:"ISPRIMARY"`
	IsRollback          string `json:"ISROLLBACK"`
	LocalResId          string `json:"LOCALRESID"`
	LocalResName        string `json:"LOCALRESNAME"`
	LocalResType        string `json:"LOCALRESTYPE"`
	PriResDataStatus    string `json:"PRIRESDATASTATUS"`
	RecoveryPolicy      string `json:"RECOVERYPOLICY"`
	RemoteDeviceId      string `json:"REMOTEDEVICEID"`
	RemoteDeviceName    string `json:"REMOTEDEVICENAME"`
	RemoteDeviceSn      string `json:"REMOTEDEVICESN"`
	RemoteResId         string `json:"REMOTERESID"`
	RemoteResName       string `json:"REMOTERESNAME"`
	ReplicationMode     string `json:"REPLICATIONMODEL"`
	ReplicationProgress string `json:"REPLICATIONPROGRESS"`
	RunningStatus       string `json:"RUNNINGSTATUS"`
	SecResAccess        string `json:"SECRESACCESS"`
	SecResDataStatus    string `json:"SECRESDATASTATUS"`
	Speed               string `json:"SPEED"`
	SynchronizeType     string `json:"SYNCHRONIZETYPE"`
	SyncLeftTime        string `json:"SYNCLEFTTIME"`
	TimeDifference      string `json:"TIMEDIFFERENCE"`
	RemTimeoutPeriod    string `json:"REMTIMEOUTPERIOD"`
	Type                string `json:"TYPE"`
}

type ReplicationPairResp struct {
	Data  ReplicationPair `json:"data"`
	Error Error           `json:"error"`
}

type SimpleStruct struct {
	Id   string `json:"ID"`
	Name string `json:"NAME"`
}

type SimpleResp struct {
	Data  []SimpleStruct `json:"data"`
	Error Error          `json:"error"`
}

type FCInitiatorsResp struct {
	Data  []FCInitiator `json:"data"`
	Error Error         `json:"error"`
}

type FCInitiator struct {
	Isfree        bool   `json:"ISFREE"`
	Id            string `json:"ID"`
	Type          int    `json:"TYPE"`
	RunningStatus string `json:"RUNNINGSTATUS"`
	ParentId      string `json:"PARENTID"`
	ParentType    int    `json:"PARENTTYPE"`
}

type FCTargWWPNResp struct {
	Data  []FCTargWWPN `json:"data"`
	Error Error        `json:"error"`
}

type FCTargWWPN struct {
	IniPortWWN  string `json:"INITIATOR_PORT_WWN"`
	TargPortWWN string `json:"TARGET_PORT_WWN"`
}

type ObjCountResp struct {
	Data  Count `json:"data"`
	Error Error `json:"error"`
}

type Count struct {
	Count string `json:"COUNT"`
}

type Performance struct {
	Uuid       string `json:"CMO_STATISTIC_UUID"`
	DataIdList string `json:"CMO_STATISTIC_DATA_ID_LIST"`
	DataList   string `json:"CMO_STATISTIC_DATA_LIST"`
	TimeStamp  string `json:"CMO_STATISTIC_TIMESTAMP"`
}

type PerformancesResp struct {
	Data  []Performance `json:"data"`
	Error Error         `json:"error"`
}
