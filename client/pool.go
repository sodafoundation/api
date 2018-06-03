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
	"strconv"
	"strings"

	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/urls"
)

// NewPoolMgr
func NewPoolMgr(r Receiver, edp string, tenantId string) *PoolMgr {
	return &PoolMgr{
		Receiver: r,
		Endpoint: edp,
		TenantId: tenantId,
	}
}

// PoolMgr
type PoolMgr struct {
	Receiver
	Endpoint string
	TenantId string
}

// GetPool
func (p *PoolMgr) GetPool(polID string) (*model.StoragePoolSpec, error) {
	var res model.StoragePoolSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GeneratePoolURL(urls.Client, p.TenantId, polID)}, "/")

	if err := p.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListPools
func (p *PoolMgr) ListPools(v []string, pool *model.StoragePoolSpec) ([]*model.StoragePoolSpec, error) {
	var res []*model.StoragePoolSpec
	var u string

	url := strings.Join([]string{
		p.Endpoint,
		urls.GeneratePoolURL(urls.Client, p.TenantId)}, "/")

	var limit, offset, sortDir, sortKey, availabilityZone, createdAt, description, dockId, freeCapacity, id, name, status,
		storageType, totalCapacity, updatedAt string
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
	if pool.AvailabilityZone != "" {
		availabilityZone = "AvailabilityZone=" + pool.AvailabilityZone
		urlpara = append(urlpara, availabilityZone)
	}
	if pool.CreatedAt != "" {
		createdAt = "CreatedAt=" + pool.CreatedAt
		urlpara = append(urlpara, createdAt)
	}
	if pool.Description != "" {
		description = "Description=" + pool.Description
		urlpara = append(urlpara, description)
	}
	if pool.DockId != "" {
		dockId = "DockId=" + pool.DockId
		urlpara = append(urlpara, dockId)
	}
	if pool.FreeCapacity != 0 {
		freeCapacity = "FreeCapacity=" + strconv.FormatInt(pool.FreeCapacity, 10)
		urlpara = append(urlpara, freeCapacity)
	}
	if pool.Id != "" {
		id = "Id=" + pool.Id
		urlpara = append(urlpara, id)
	}
	if pool.Name != "" {
		name = "Name=" + pool.Name
		urlpara = append(urlpara, name)
	}
	if pool.Status != "" {
		status = "Status=" + pool.Status
		urlpara = append(urlpara, status)
	}
	if pool.StorageType != "" {
		storageType = "StorageType=" + pool.StorageType
		urlpara = append(urlpara, storageType)
	}
	if pool.TotalCapacity != 0 {
		totalCapacity = "TotalCapacity=" + strconv.FormatInt(pool.TotalCapacity, 10)
		urlpara = append(urlpara, totalCapacity)
	}
	if pool.UpdatedAt != "" {
		updatedAt = "UpdatedAt=" + pool.UpdatedAt
		urlpara = append(urlpara, updatedAt)
	}
	if len(urlpara) > 0 {
		u = strings.Join(urlpara, "&")
		url += "?" + u
	}

	if err := p.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return res, nil
}
