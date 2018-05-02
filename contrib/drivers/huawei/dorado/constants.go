// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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
//    License for the specific l

package dorado

import "time"

const (
	defaultConfPath = "/etc/opensds/driver/huawei_dorado.yaml"
	defaultAZ       = "default"
)

const (
	KLunId  = "huaweiLunId"
	KSnapId = "huaweiSnapId"
)

const (
	KPairId          = "huaweiReplicaPairId"   // replication pair
	KSecondaryLunId  = "huaweiSecondaryLunId"  // secondary lun id
	KSecondaryLunWwn = "huaweiSecondaryLunWwn" // secondary lun wwn
)

const (
	StatusHealth       = "1"
	StatusActive       = "43"
	StatusRunning      = "10"
	StatusVolumeReady  = "27"
	StatusLuncopyReady = "40"
	StatusQosActive    = "2"
	StatusQosInactive  = "45"
)

// Array type
const (
	ReplicationArrayType   = "1"
	HeterogeneityArrayType = "2"
	UnknownArrayType       = "3"
)

// Health status
const (
	NormalHealthStatus  = "1"
	FaultHealthStatus   = "2"
	InvalidHealthStatus = "14"
)

// Running status
const (
	LinkUpRunningStatus     = "10"
	LinkDownRunningStatus   = "11"
	DisabledRunningStatus   = "31"
	ConnectingRunningStatus = "101"
)

const (
	ThickLuntype         = 0
	ThinLuntype          = 1
	MaxNameLength        = 31
	MaxDescriptionLength = 170
	PortNumPerContr      = 2
	PwdExpired           = 3
	PwdReset             = 4
)

const (
	ErrorConnectToServer      = -403
	ErrorUnauthorizedToServer = -401
)

const (
	MappingViewPrefix = "OpenSDS_MappingView_"
	LunGroupPrefix    = "OpenSDS_LunGroup_"
	HostGroupPrefix   = "OpenSDS_HostGroup_"
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

	ReplicaRunningStatusKey         = "RUNNINGSTATUS"
	ReplicaRunningStatusInitialSync = "21"
	ReplicaRunningStatusSync        = "23"
	ReplicaRunningStatusSynced      = "24"
	ReplicaRunningStatusNormal      = "1"
	ReplicaRunningStatusSplit       = "26"
	ReplicaRunningStatusErrupted    = "34"
	ReplicaRunningStatusInvalid     = "35"

	ReplicaHealthStatusKey    = "HEALTHSTATUS"
	ReplicaHealthStatusNormal = "1"

	ReplicaLocalDataStatusKey   = "PRIRESDATASTATUS"
	ReplicaRemoteDataStatusKey  = "SECRESDATASTATUS"
	ReplicaDataSyncKey          = "ISDATASYNC"
	ReplicaDataStatusSynced     = "1"
	ReplicaDataStatusComplete   = "2"
	ReplicaDataStatusIncomplete = "3"
)
