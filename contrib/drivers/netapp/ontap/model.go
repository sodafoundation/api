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

package ontap

import (
	"github.com/netapp/trident/storage_drivers/ontap"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
)

type BackendOptions struct {
	Version           int    `yaml:"version"`
	StorageDriverName string `yaml:"storageDriverName"`
	ManagementLIF     string `yaml:"managementLIF"`
	DataLIF           string `yaml:"dataLIF"`
	Svm               string `yaml:"svm"`
	IgroupName        string `yaml:"igroupName"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
}

type ONTAPConfig struct {
	BackendOptions `yaml:"backendOptions"`
	Pool           map[string]PoolProperties `yaml:"pool,flow"`
}

type Driver struct {
	sanStorageDriver *ontap.SANStorageDriver
	conf             *ONTAPConfig
}

type Pool struct {
	PoolId        int   `json:"poolId"`
	TotalCapacity int64 `json:"totalCapacity"`
	AllocCapacity int64 `json:"allocatedCapacity"`
	UsedCapacity  int64 `json:"usedCapacity"`
}
