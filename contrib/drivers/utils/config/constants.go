// Copyright 2019 The OpenSDS Authors.
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

/*
This module defines some essential configuration infos for all storage drivers.

*/

package config

// These constants below represent the vendor name of all storage drivers which
// can be supported by now.
const (
	CinderDriverType              = "cinder"
	CephDriverType                = "ceph"
	LVMDriverType                 = "lvm"
	HuaweiDoradoDriverType        = "huawei_dorado"
	HuaweiFusionStorageDriverType = "huawei_fusionstorage"
	HuaweiOceanFileDriverType     = "huawei_oceanstor"
	HpeNimbleDriverType           = "hpe_nimble"
	DRBDDriverType                = "drbd"
	ManilaDriverType              = "manila"
)

const (
	NFSDriverType = "nfs"
)

// These constants below represent the access protocol type of all storage
// drivers which can be supported by now. Please NOTICE that currently these
// constants can NOT be used by all methods except InitializeConnection().
const (
	ISCSIProtocol  = "iscsi"
	RBDProtocol    = "rbd"
	FCProtocol     = "fibre_channel"
	NVMEOFProtocol = "nvmeof"
)

// Telemetry metric resource type
const (
	MetricResourceTypeController = "controller"
	MetricResourceTypePool       = "pool"
	MetricResourceTypeFilesystem = "filesystem"
	MetricResourceTypeDisk       = "disk"
	MetricResourceTypeVolume     = "volume"
)

const (
	KMetricIOPS               = "iops"
	KMetricBandwidth          = "bandwidth"
	KMetricLatency            = "latency"
	KMetricServiceTime        = "service_time"
	KMetricUtilizationPercent = "utilization_prcnt"
	KMetricCacheHitRatio      = "cache_hit_ratio"
	KMetricCpuUsage           = "cpu_usage"
)
