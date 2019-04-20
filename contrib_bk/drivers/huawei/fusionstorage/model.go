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

package fusionstorage

type requesData struct {
	Timeout int         `json:"timeout"`
	Data    interface{} `json:"data`
}

type responseResult struct {
	RespCode int      `json:"result"`
	Details  []detail `json:"detail"`
}

type detail struct {
	Description string `json:"description,omitempty"`
	ErrorCode   int    `json:"errorCode,omitempty"`
}

func (r *responseResult) GetDescription() string {
	for _, v := range r.Details {
		if v.Description != "" {
			return v.Description
		}
	}

	return ""
}

func (r *responseResult) GetErrorCode() int {
	for _, v := range r.Details {
		if v.ErrorCode != 0 {
			return v.ErrorCode
		}
	}

	return 0
}

type version struct {
	CurrentVersion string `json:"currentVersion"`
}

type poolResp struct {
	Pools []pool `json:"storagePools"`
}

type pool struct {
	PoolId        int   `json:"poolId"`
	TotalCapacity int64 `json:"totalCapacity"`
	AllocCapacity int64 `json:"allocatedCapacity"`
	UsedCapacity  int64 `json:"usedCapacity"`
}

type hostList struct {
	HostList []host `json:"hostList"`
}

type host struct {
	HostName string `json:"hostName"`
}

type portHostMap struct {
	PortHostMap map[string][]string `json:"portHostMap"`
}

type hostLunList struct {
	LunList []lunList `json:"hostLunList"`
}

type lunList struct {
	Id   int    `json:"lunId"`
	Name string `json:"lunName"`
}
