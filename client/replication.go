// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

type ReplicationBuilder *model.ReplicationSpec
type FailoverReplicationBuilder *model.FailoverReplicationSpec

// NewReplicationMgr
func NewReplicationMgr(r Receiver, edp string, tenantId string) *ReplicationMgr {
	return &ReplicationMgr{
		Receiver: r,
		Endpoint: edp,
		TenantId: tenantId,
	}
}

// ReplicationMgr
type ReplicationMgr struct {
	Receiver
	Endpoint string
	TenantId string
}

// CreateReplication
func (v *ReplicationMgr) CreateReplication(body ReplicationBuilder) (*model.ReplicationSpec, error) {
	var res model.ReplicationSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateReplicationURL(urls.Client, v.TenantId)}, "/")

	if err := v.Recv(url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetReplication
func (v *ReplicationMgr) GetReplication(replicaId string) (*model.ReplicationSpec, error) {
	var res model.ReplicationSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateReplicationURL(urls.Client, v.TenantId, replicaId)}, "/")

	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListReplications
func (v *ReplicationMgr) ListReplications(p []string, rep *model.ReplicationSpec) ([]*model.ReplicationSpec, error) {
	var res []*model.ReplicationSpec
	var u string

	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateReplicationURL(urls.Client, v.TenantId, "detail")}, "/")

	var limit, offset, sortDir, sortKey, id, createdAt, updatedAt, name, description, primaryVolumeId, secondaryVolumeId string
	var urlpara []string
	if len(p) > 0 {
		if p[0] != "" {
			limit = "limit=" + p[0]
			urlpara = append(urlpara, limit)
		}
		if p[1] != "" {
			offset = "offset=" + p[1]
			urlpara = append(urlpara, offset)
		}
		if p[2] != "" {
			sortDir = "sortDir=" + p[2]
			urlpara = append(urlpara, sortDir)
		}
		if p[3] != "" {
			sortKey = "sortKey=" + p[3]
			urlpara = append(urlpara, sortKey)
		}
	}
	if rep.Id != "" {
		id = "Id=" + rep.Id
		urlpara = append(urlpara, id)
	}
	if rep.CreatedAt != "" {
		createdAt = "CreatedAt=" + rep.CreatedAt
		urlpara = append(urlpara, createdAt)
	}
	if rep.UpdatedAt != "" {
		updatedAt = "UpdatedAt=" + rep.UpdatedAt
		urlpara = append(urlpara, updatedAt)
	}
	if rep.Name != "" {
		name = "Name=" + rep.Name
		urlpara = append(urlpara, name)
	}
	if rep.Description != "" {
		description = "Description=" + rep.Description
		urlpara = append(urlpara, description)
	}
	if rep.PrimaryVolumeId != "" {
		primaryVolumeId = "Description=" + rep.PrimaryVolumeId
		urlpara = append(urlpara, primaryVolumeId)
	}
	if rep.SecondaryVolumeId != "" {
		secondaryVolumeId = "Description=" + rep.SecondaryVolumeId
		urlpara = append(urlpara, secondaryVolumeId)
	}

	if len(urlpara) > 0 {
		u = strings.Join(urlpara, "&")
		url += "?" + u
	}

	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteReplication
func (v *ReplicationMgr) DeleteReplication(replicaId string, body ReplicationBuilder) error {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateReplicationURL(urls.Client, v.TenantId, replicaId)}, "/")
	return v.Recv(url, "DELETE", body, nil)
}

// UpdateReplication
func (v *ReplicationMgr) UpdateReplication(replicaId string, body ReplicationBuilder) (*model.ReplicationSpec, error) {
	var res model.ReplicationSpec
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateReplicationURL(urls.Client, v.TenantId, replicaId)}, "/")

	if err := v.Recv(url, "PUT", body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// EnableReplication
func (v *ReplicationMgr) EnableReplication(replicaId string) error {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateReplicationURL(urls.Client, v.TenantId, replicaId, "enable")}, "/")
	return v.Recv(url, "POST", nil, nil)
}

// EnableReplication
func (v *ReplicationMgr) DisableReplication(replicaId string) error {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateReplicationURL(urls.Client, v.TenantId, replicaId, "disable")}, "/")
	return v.Recv(url, "POST", nil, nil)
}

// EnableReplication
func (v *ReplicationMgr) FailoverReplication(replicaId string, body FailoverReplicationBuilder) error {
	url := strings.Join([]string{
		v.Endpoint,
		urls.GenerateReplicationURL(urls.Client, v.TenantId, replicaId, "failover")}, "/")
	return v.Recv(url, "POST", body, nil)
}
