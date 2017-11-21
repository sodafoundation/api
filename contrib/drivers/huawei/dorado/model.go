// Copyright (c) 2017 OpenSDS Authors.
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

package dorado

type Error struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
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
	CompressedCapacity              string `json:"COMPRESSEDCAPACITY"`
	CompressInvolvedCapacity        string `json:"COMPRESSINVOLVEDCAPACITY"`
	CompressionRate                 string `json:"COMPRESSIONRATE"`
	DataSpace                       string `json:"DATASPACE"`
	DedupedCapacity                 string `json:"DEDUPEDCAPACITY"`
	DedupInvolvedCapacity           string `json:"DEDUPINVOLVEDCAPACITY"`
	DeduplicationRate               string `json:"DEDUPLICATIONRATE"`
	Description                     string `json:"DESCRIPTION"`
	EndingUpThreshold               string `json:"ENDINGUPTHRESHOLD"`
	HealthStatus                    string `json:"HEALTHSTATUS"`
	Id                              string `json:"ID"`
	LunConfigedCapacity             string `json:"LUNCONFIGEDCAPACITY"`
	Name                            string `json:"NAME"`
	ParentId                        string `json:"PARENTID"`
	ParentName                      string `json:"PARENTNAME"`
	ParentType                      int    `json:"PARENTTYPE"`
	ReductionInvolvedCapacity       string `json:"REDUCTIONINVOLVEDCAPACITY"`
	ReplicationCapacity             string `json:"REPLICATIONCAPACITY"`
	RunningStatus                   string `json:"RUNNINGSTATUS"`
	SaveCapacityRate                string `json:"SAVECAPACITYRATE"`
	SpaceReductionRate              string `json:"SPACEREDUCTIONRATE"`
	ThinProvisionSavePercentage     string `json:"THINPROVISIONSAVEPERCENTAGE"`
	Tier0Capacity                   string `json:"TIER0CAPACITY"`
	Tier0Disktype                   string `json:"TIER0DISKTYPE"`
	Tier0RaidLv                     string `json:"TIER0RAIDLV"`
	TotalLunWriteCapacity           string `json:"TOTALLUNWRITECAPACITY"`
	Type                            int    `json:"TYPE"`
	UsageType                       string `json:"USAGETYPE"`
	UserConsumedCapacity            string `json:"USERCONSUMEDCAPACITY"`
	UserConsumedCapacityPercentage  string `json:"USERCONSUMEDCAPACITYPERCENTAGE"`
	UserConsumedCapacityThreshold   string `json:"USERCONSUMEDCAPACITYTHRESHOLD"`
	UserConsumedCapacityWithoutMeta string `json:"USERCONSUMEDCAPACITYWITHOUTMETA"`
	UserFreeCapacity                string `json:"USERFREECAPACITY"`
	UserTotalCapacity               string `json:"USERTOTALCAPACITY"`
	UserWriteAllocCapacity          string `json:"USERWRITEALLOCCAPACITY"`
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

