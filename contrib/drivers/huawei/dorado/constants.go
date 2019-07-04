// Copyright 2018 The OpenSDS Authors.
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
// License for the specific l

package dorado

import "time"

//  default value for driver
const (
	defaultConfPath = "/etc/opensds/driver/huawei_dorado.yaml"
	defaultAZ       = "default"
)

// metadata keys
const (
	KLunId  = "huaweiLunId"
	KSnapId = "huaweiSnapId"
	KPairId = "huaweiReplicaPairId" // replication pair
)

// name prefix
const (
	PrefixMappingView = "OpenSDS_MappingView_"
	PrefixLunGroup    = "OpenSDS_LunGroup_"
	PrefixHostGroup   = "OpenSDS_HostGroup_"
)

// object type id
const (
	ObjectTypeLun             = "11"
	ObjectTypeHost            = "21"
	ObjectTypeSnapshot        = "27"
	ObjectTypeHostGroup       = "14"
	ObjectTypeController      = "207"
	ObjectTypePool            = "216"
	ObjectTypeLunCopy         = 219 // not string, should be integer
	ObjectTypeIscsiInitiator  = "222"
	ObjectTypeFcInitiator     = "223"
	ObjectTypeMappingView     = "245"
	ObjectTypeLunGroup        = "256"
	ObjectTypeReplicationPair = "263"
)

// Error Code
const (
	ErrorConnectToServer      = -403
	ErrorUnauthorizedToServer = -401
	ErrorObjectUnavailable    = 1077948996
	ErrorHostGroupNotExist    = 1077937500
)

// misc
const (
	ThickLunType         = 0
	ThinLunType          = 1
	MaxNameLength        = 31
	MaxDescriptionLength = 170
	PortNumPerContr      = 2
	PwdExpired           = 3
	PwdReset             = 4
)

// lun copy
const (
	LunCopySpeedLow     = "1"
	LunCopySpeedMedium  = "2"
	LunCopySpeedHigh    = "3"
	LunCopySpeedHighest = "4"
)

var LunCopySpeedTypes = []string{LunCopySpeedLow, LunCopySpeedMedium, LunCopySpeedHigh, LunCopySpeedHighest}

const (
	LunReadyWaitInterval = 2 * time.Second
	LunReadyWaitTimeout  = 20 * time.Second
	LunCopyWaitInterval  = 2 * time.Second
	LunCopyWaitTimeout   = 200 * time.Second
)

// Object status key id
const (
	StatusHealth          = "1"
	StatusQosActive       = "2"
	StatusRunning         = "10"
	StatusVolumeReady     = "27"
	StatusLunCoping       = "39"
	StatusLunCopyStop     = "38"
	StatusLunCopyQueue    = "37"
	StatusLunCopyNotStart = "36"
	StatusLunCopyReady    = "40"
	StatusActive          = "43"
	StatusQosInactive     = "45"
)

// Array type
const (
	ArrayTypeReplication   = "1"
	ArrayTypeHeterogeneity = "2"
	ArrayTypeUnknown       = "3"
)

// Health status
const (
	HealthStatusNormal          = "1"
	HealthStatusFault           = "2"
	HealthStatusPreFail         = "3"
	HealthStatusPartiallyBroken = "4"
	HealthStatusDegraded        = "5"
	HealthStatusBadSectorsFound = "6"
	HealthStatusBitErrorsFound  = "7"
	HealthStatusConsistent      = "8"
	HealthStatusInconsistent    = "9"
	HealthStatusBusy            = "10"
	HealthStatusNoInput         = "11"
	HealthStatusLowBattery      = "12"
	HealthStatusSingleLinkFault = "13"
	HealthStatusInvalid         = "14"
	HealthStatusWriteProtect    = "15"
)

// Running status
const (
	RunningStatusNormal      = "1"
	RunningStatusLinkUp      = "10"
	RunningStatusLinkDown    = "11"
	RunningStatusOnline      = "27"
	RunningStatusDisabled    = "31"
	RunningStatusInitialSync = "21"
	RunningStatusSync        = "23"
	RunningStatusSynced      = "24"
	RunningStatusSplit       = "26"
	RunningStatusInterrupted = "34"
	RunningStatusInvalid     = "35"
	RunningStatusConnecting  = "101"
)

const (
	DefaultReplicaWaitInterval = 1 * time.Second
	DefaultReplicaWaitTimeout  = 20 * time.Second

	ReplicaSyncMode  = "1"
	ReplicaAsyncMode = "2"
	ReplicaSpeed     = "2"
	ReplicaPeriod    = "3600"
	ReplicaSecondRo  = "2"
	ReplicaSecondRw  = "3"

	ReplicaRunningStatusKey   = "RUNNINGSTATUS"
	ReplicaHealthStatusKey    = "HEALTHSTATUS"
	ReplicaHealthStatusNormal = "1"

	ReplicaLocalDataStatusKey   = "PRIRESDATASTATUS"
	ReplicaRemoteDataStatusKey  = "SECRESDATASTATUS"
	ReplicaDataSyncKey          = "ISDATASYNC"
	ReplicaDataStatusSynced     = "1"
	ReplicaDataStatusComplete   = "2"
	ReplicaDataStatusIncomplete = "3"
)

// performance key ids
const (
	PerfUtilizationPercent = "18"  // usage ratioPerf
	PerfBandwidth          = "21"  // mbs
	PerfIOPS               = "22"  // tps
	PerfServiceTime        = "29"  // excluding queue time(ms)
	PerfCpuUsage           = "68"  // %
	PerfCacheHitRatio      = "303" // %
	PerfLatency            = "370" // ms
)
