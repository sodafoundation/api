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

package client

import (
	"strings"

	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/urls"
)

func NewDockMgr(r Receiver, edp string, tenantId string) *DockMgr {
	return &DockMgr{
		Receiver: r,
		Endpoint: edp,
		TenantId: tenantId,
	}
}

type DockMgr struct {
	Receiver
	Endpoint string
	TenantId string
}

func (d *DockMgr) GetDock(dckID string) (*model.DockSpec, error) {
	var res model.DockSpec
	url := strings.Join([]string{
		d.Endpoint,
		urls.GenerateDockURL(urls.Client, d.TenantId, dckID)}, "/")

	if err := d.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (d *DockMgr) ListDocks(v []string, dock *model.DockSpec) ([]*model.DockSpec, error) {
	var res []*model.DockSpec
	var u string

	url := strings.Join([]string{
		d.Endpoint,
		urls.GenerateDockURL(urls.Client, d.TenantId)}, "/")

	var limit, offset, sortDir, sortKey, createdAt, description, driverName, endpoint, id, name, status, storageType, updatedAt string
	var urlpara []string

	if len(v) > 0 {
		if v[0] != "" {
			limit = "limit=" + v[0]
			urlpara = append(urlpara, limit)
		}
		if v[1] != "" {
			offset = "offset=" + v[1]
			urlpara = append(urlpara, offset)
		}
		if v[2] != "" {
			sortDir = "sortDir=" + v[2]
			urlpara = append(urlpara, sortDir)
		}
		if v[3] != "" {
			sortKey = "sortKey=" + v[3]
			urlpara = append(urlpara, sortKey)
		}
	}
	if dock.CreatedAt != "" {
		createdAt = "CreatedAt=" + dock.CreatedAt
		urlpara = append(urlpara, createdAt)
	}
	if dock.Description != "" {
		description = "Description=" + dock.Description
		urlpara = append(urlpara, description)
	}
	if dock.DriverName != "" {
		driverName = "DriverName=" + dock.DriverName
		urlpara = append(urlpara, driverName)
	}
	if dock.Endpoint != "" {
		endpoint = "Endpoint=" + dock.Endpoint
		urlpara = append(urlpara, endpoint)
	}
	if dock.Id != "" {
		id = "id=" + dock.Id
		urlpara = append(urlpara, id)
	}
	if dock.Name != "" {
		name = "Name=" + dock.Name
		urlpara = append(urlpara, name)
	}
	if dock.Status != "" {
		status = "Status=" + dock.Status
		urlpara = append(urlpara, status)
	}
	if dock.StorageType != "" {
		storageType = "StorageType=" + dock.StorageType
		urlpara = append(urlpara, storageType)
	}
	if dock.UpdatedAt != "" {
		updatedAt = "UpdatedAt=" + dock.UpdatedAt
		urlpara = append(urlpara, updatedAt)
	}
	if len(urlpara) > 0 {
		u = strings.Join(urlpara, "&")
		url += "?" + u
	}
	if err := d.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return res, nil
}
