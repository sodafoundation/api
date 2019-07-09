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

package oceanstor

const (
	DefaultConfPath       = "/etc/opensds/driver/oceanstor_fileshare.yaml"
	PwdExpired            = 3
	PwdReset              = 4
	NFSProto              = "nfs"
	CIFSProto             = "cifs"
	UnitGi                = 1024 * 1024 * 1024
	defaultAZ             = "default"
	defaultFileSystem     = "opensds_file_system"
	StatusFSHealth        = "1"
	StatusFSRunning       = "27"
	AccessLevelRW         = "rw"
	AccessLevelRO         = "ro"
	AccessNFSRw           = "1"
	AccessNFSRo           = "0"
	AccessCIFSRo          = "0"
	AccessCIFSFullControl = "1"
	MaxRetry              = 3
	FileShareName         = "fileshareName"
	FileShareID           = "shareId"
	NamePrefix            = "opensds"
	FileShareSnapshotID   = "fileshareSnapId"
	AccessTypeUser        = "user"
	AccessTypeIp          = "ip"
	AccessLevelRead       = "read"
	AccessLevelWrite      = "write"
	AccessLevelExecute    = "execute"
)
