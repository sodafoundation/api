// Copyright (c) 2019 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package oceanstor

const (
	BasicURI                   = "/dsware/service/"
	CallTimeout                = 50
	LoginSocketTimeout         = 32
	UnitGiShiftBit             = 10
	DefaultAZ                  = "default"
	NamePrefix                 = "opensds"
	LunId                      = "lunId"
	ConnectorType              = "connectorType"
	FusionstorageIscsi         = "fusionstorage_iscsi"
	InitiatorNotExistErrorCode = "32155103"
	VolumeNotExistErrorCode    = "32150005"
	CmdBin                     = "fsc_cli"
	MaxRetryNode               = 3
	DefaultConfPath            = "/etc/opensds/driver/fusionstorage.yaml"
	PwdExpired                 = 3
	PwdReset                   = 4
	NFS                        = "NFS"
	CIFS                       = "CIFS"
	UnitGi                     = 1024 * 1024 * 1024
	defaultAZ                  = "default"
	defaultFileSystem          = "opensds_file_system"
	StatusFSHealth             = "1"
	StatusFSRunning            = "27"
)
