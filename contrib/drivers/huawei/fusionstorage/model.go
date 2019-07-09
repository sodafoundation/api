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

import (
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
)

type Driver struct {
	Client *FsClient
	Conf   *Config
}

type AuthOptions struct {
	Username        string   `yaml:"username"`
	Password        string   `yaml:"password"`
	Url             string   `yaml:"url"`
	FmIp            string   `yaml:"fmIp,omitempty"`
	FsaIp           []string `yaml:"fsaIp,flow"`
	PwdEncrypter    string   `yaml:"PwdEncrypter,omitempty"`
	EnableEncrypted bool     `yaml:"EnableEncrypted,omitempty"`
	Version         string   `json:"version"`
}

type Config struct {
	AuthOptions `yaml:"authOptions"`
	Pool        map[string]PoolProperties `yaml:"pool,flow"`
}

type RequesData struct {
	Timeout int         `json:"timeout"`
	Data    interface{} `json:"data`
}

type ResponseResult struct {
	RespCode int      `json:"result"`
	Details  []Detail `json:"detail"`
}

type Detail struct {
	Description string `json:"description,omitempty"`
	ErrorCode   int    `json:"errorCode,omitempty"`
}

type Version struct {
	CurrentVersion string `json:"currentVersion"`
}

type PoolResp struct {
	Pools []Pool `json:"storagePools"`
}

type Pool struct {
	PoolId        int   `json:"poolId"`
	TotalCapacity int64 `json:"totalCapacity"`
	AllocCapacity int64 `json:"allocatedCapacity"`
	UsedCapacity  int64 `json:"usedCapacity"`
}

type HostList struct {
	HostList []Host `json:"hostList"`
}

type Host struct {
	HostName string `json:"hostName"`
}

type PortHostMap struct {
	PortHostMap map[string][]string `json:"portHostMap"`
}

type HostLunList struct {
	LunList []LunList `json:"hostLunList"`
}

type LunList struct {
	Id   int    `json:"lunId"`
	Name string `json:"lunName"`
}

type IscsiPortal struct {
	NodeResultList []NodeResult `json:"nodeResultList"`
}

type NodeResult struct {
	PortalList []Portal `json:"iscsiPortalList"`
}

type Portal struct {
	IscsiPortal string `json:"iscsiPortal"`
	Status      string `json:"iscsiStatus"`
}

type DeviceVersion struct {
	Version string `json:"version"`
}
