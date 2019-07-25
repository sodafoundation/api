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
// License for the specific language governing permissions and limitations
// under the License.

package fusionstorage

const (
	BasicURI                           = "/dsware/service/"
	UnitGiShiftBit                     = 10
	DefaultAZ                          = "default"
	NamePrefix                         = "opensds"
	LunId                              = "lunId"
	FusionstorageIscsi                 = "fusionstorage_iscsi"
	InitiatorNotExistErrorCodeVersion6 = "32155103"
	InitiatorNotExistErrorCodeVersion8 = "155103"
	VolumeAlreadyInHostErrorCode       = "157001"
	CmdBin                             = "fsc_cli"
	DefaultConfPath                    = "/etc/opensds/driver/fusionstorage.yaml"
	ClientVersion6_3                   = "6.3"
	ClientVersion8_0                   = "8.0"
	MaxRetry                           = 3
)
