// Copyright 2019 The OpenSDS Authors.
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

package ontap

//  default value for driver
const (
	defaultConfPath   = "/etc/opensds/driver/netapp_ontap_san.yaml"
	DefaultAZ         = "default"
	volumePrefix      = "opensds_"
	snapshotPrefix    = "opensds_snapshot_"
	naaPrefix         = "60a98000" //Applicable to Clustered Data ONTAP and Data ONTAP 7-Mode
	KLvPath           = "lunPath"
	KLvIdFormat       = "NAA"
	StorageDriverName = "ontap-san"
	driverContext      = "csi"
	VolumeVersion     = "1"
	SnapshotVersion   = "1"
	accessMode        = ""
	volumeMode        = "Block"
	bytesGiB           = 1073741824
)
