// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

/*
This module implements the common data structure.

*/

package model

import (
	"sort"
	"strconv"
	"strings"
)

// A pool is discoveried and updated by a dock service. Each pool can be regarded
// as a physical storage pool or a virtual storage pool. It's a logical and
// atomic pool and can be abstracted from any storage platform.
type StoragePoolSpec struct {
	*BaseModel
	// The uuid of project
	// + readOnly
	TenantId string `json:"tenantId"`

	// The name of the pool.
	Name string `json:"name,omitempty"`

	// The description of the pool.
	// +optional
	Description string `json:"description,omitempty"`

	// The status of the pool.
	// One of: "available" or "unavailable".
	Status string `json:"status,omitempty"`

	// The uuid of the dock which the pool belongs to.
	DockId string `json:"dockId,omitempty"`

	// The locality that pool belongs to.
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// The total capacity of the pool.
	// Default unit of TotalCapacity is GB.
	TotalCapacity int64 `json:"totalCapacity,omitempty"`

	// The free capaicty of the pool.
	// Default unit of FreeCapacity is GB.
	FreeCapacity int64 `json:"freeCapacity,omitempty"`

	// The storage type of the storage pool.
	// One of: "block", "file" or "object".
	StorageType string `json:"storageType,omitempty"`

	// Map of keys and StoragePoolExtraSpec object that represents the properties
	// of the pool, such as supported capabilities.
	// +optional
	Extras StoragePoolExtraSpec `json:"extras,omitempty"`
}

type StoragePoolExtraSpec struct {
	// DataStorage represents suggested some data storage capabilities.
	DataStorage DataStorageLoS `json:"dataStorage,omitempty" yaml:"dataStorage,omitempty"`
	// IOConnectivity represents some suggested IO connectivity capabilities.
	IOConnectivity IOConnectivityLoS `json:"ioConnectivity,omitempty" yaml:"ioConnectivity,omitempty"`
	// DataProtection represents some suggested data protection capabilities.
	DataProtection DataProtectionLos `json:"dataProtection,omitempty" yaml:"dataProtection,omitempty"`

	// Besides those basic suggested pool properties above, vendors can configure
	// some advanced features (diskType, IOPS, throughout, latency, etc)
	// themselves, all these properties can be exposed to controller scheduler
	// and filtered by selector in a extensible way.
	Advanced map[string]interface{} `json:"advanced,omitempty" yaml:"advanced,omitempty"`
}

var poolSortKey string

type StoragePoolSlice []*StoragePoolSpec

func (pool StoragePoolSlice) Len() int { return len(pool) }

func (pool StoragePoolSlice) Swap(i, j int) { pool[i], pool[j] = pool[j], pool[i] }

func (pool StoragePoolSlice) Less(i, j int) bool {
	switch poolSortKey {

	case "ID":
		return pool[i].Id < pool[j].Id
	case "NAME":
		return pool[i].Name < pool[j].Name
	case "STATUS":
		return pool[i].Status < pool[j].Status
	case "AVAILABILITYZONE":
		return pool[i].AvailabilityZone < pool[j].AvailabilityZone
	case "DOCKID":
		return pool[i].DockId < pool[j].DockId
	case "DESCRIPTION":
		return pool[i].Description < pool[j].Description
	}
	return false
}

func (c *StoragePoolSpec) FindValue(k string, p *StoragePoolSpec) string {
	switch k {
	case "Id":
		return p.Id
	case "CreatedAt":
		return p.CreatedAt
	case "UpdatedAt":
		return p.UpdatedAt
	case "Name":
		return p.Name
	case "Description":
		return p.Description
	case "Status":
		return p.Status
	case "DockId":
		return p.DockId
	case "AvailabilityZone":
		return p.AvailabilityZone
	case "TotalCapacity":
		return strconv.FormatInt(p.TotalCapacity, 10)
	case "FreeCapacity":
		return strconv.FormatInt(p.FreeCapacity, 10)
	case "StorageType":
		return p.StorageType
	}
	return ""
}

func (c *StoragePoolSpec) SortList(pools []*StoragePoolSpec, sortKey, sortDir string) []*StoragePoolSpec {

	poolSortKey = sortKey

	if strings.EqualFold(sortDir, "asc") {
		sort.Sort(StoragePoolSlice(pools))
	} else {
		sort.Sort(sort.Reverse(StoragePoolSlice(pools)))
	}
	return pools
}
