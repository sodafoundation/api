// Copyright 2017 The OpenSDS Authors.
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
func (p *PoolMgr) ListPools(args ...interface{}) ([]*model.StoragePoolSpec, error) {
	var res []*model.StoragePoolSpec

	url := strings.Join([]string{
		p.Endpoint,
		urls.GeneratePoolURL(urls.Client, p.TenantId)}, "/")

	param, err := processListParam(args)
	if err != nil {
		return nil, err
	}

	if param != "" {
		url += "?" + param
	}

	if err := p.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return res, nil
}
