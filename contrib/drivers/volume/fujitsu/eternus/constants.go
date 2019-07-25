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

package eternus

// define config parameter
const (
	defaultConfPath    = "/etc/opensds/driver/fujitsu_eternus.yaml"
	defaultCliConfPath = "/etc/opensds/driver/cli_response.yml"
	defaultAZ          = "default"
	SSHPort            = "22"
)

// define eternus specific parameter
const (
	LBASize = 2097152
)

// define command
const (
	CreateVolume = "create volume"
)

// define error code
const (
	NotFound = "0110"
)

// response key name
const (
	KLunId        = "eternusVolId"
	KLunGrpId     = "eternusLunGrpId"
	KLunMaskGrpNo = "eternusLunMaskGrpNo"
	KHostId       = "eternusHostId"
	KSnapId       = "eternusSnapshotId"
	KSnapLunId    = "eternusSnapshotVolId"
)
