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
	"strings"
)

// DockSpec is initialized by specific driver configuration. Each backend
// can be regarded as a docking service between SDS controller and storage
// service.
type DockSpec struct {
	*BaseModel

	// The name of the dock.
	Name string `json:"name,omitempty"`

	// The description of the dock.
	// +optional
	Description string `json:"description,omitempty"`

	// The status of the dock.
	// One of: "available" or "unavailable".
	Status string `json:"status,omitempty"`

	// The storage type of the dock.
	// One of: "block", "file" or "object".
	StorageType string `json:"storageType,omitempty"`

	// Endpoint represents the dock server's access address.
	Endpoint string `json:"endpoint,omitempty"`

	// NodeId represents the identification of the host, it can be considered
	// as instance id or hostname.
	NodeId string `json:"nodeId,omitempty"`

	// DriverName represents the dock provider.
	// Currently One of: "cinder", "ceph", "lvm", "default".
	DriverName string `json:"driverName,omitempty"`
}

var dockSortKey string

type DockSlice []*DockSpec

func (dock DockSlice) Len() int { return len(dock) }

func (dock DockSlice) Swap(i, j int) { dock[i], dock[j] = dock[j], dock[i] }

func (dock DockSlice) Less(i, j int) bool {
	switch dockSortKey {

	case "ID":
		return dock[i].Id < dock[j].Id
	case "NAME":
		return dock[i].Name < dock[j].Name
	case "STATUS":
		return dock[i].Status < dock[j].Status
	case "ENDPOINT":
		return dock[i].Endpoint < dock[j].Endpoint
	case "DRIVERNAME":
		return dock[i].DriverName < dock[j].DriverName
	case "DESCRIPTION":
		return dock[i].Description < dock[j].Description
	}
	return false
}

func (c *DockSpec) FindValue(k string, d *DockSpec) string {
	switch k {
	case "Id":
		return d.Id
	case "CreatedAt":
		return d.CreatedAt
	case "Name":
		return d.Name
	case "UpdatedAt":
		return d.UpdatedAt
	case "Description":
		return d.Description
	case "Status":
		return d.Status
	case "StorageType":
		return d.StorageType
	case "Endpoint":
		return d.Endpoint
	case "DriverName":
		return d.DriverName
	}
	return ""
}

func (c *DockSpec) SortList(dcks []*DockSpec, sortKey, sortDir string) []*DockSpec {
	dockSortKey = sortKey
	if strings.EqualFold(sortDir, "asc") {
		sort.Sort(DockSlice(dcks))
	} else {
		sort.Sort(sort.Reverse(DockSlice(dcks)))
	}
	return dcks
}
