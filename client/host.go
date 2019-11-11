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

package client

import (
	"strings"

	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/urls"
)

// HostBuilder contains request body of handling a host request.
type HostBuilder *model.HostSpec

// NewHostMgr implementation
func NewHostMgr(r Receiver, edp string, tenantID string) *HostMgr {
	return &HostMgr{
		Receiver: r,
		Endpoint: edp,
		TenantID: tenantID,
	}
}

// HostMgr implementation
type HostMgr struct {
	Receiver
	Endpoint string
	TenantID string
}

// CreateHost implementation
func (h *HostMgr) CreateHost(body HostBuilder) (*model.HostSpec, error) {
	var res model.HostSpec

	url := strings.Join([]string{
		h.Endpoint,
		urls.GenerateHostURL(urls.Client, h.TenantID)}, "/")

	if err := h.Recv(url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetHost implementation
func (h *HostMgr) GetHost(ID string) (*model.HostSpec, error) {
	var res model.HostSpec
	url := strings.Join([]string{
		h.Endpoint,
		urls.GenerateHostURL(urls.Client, h.TenantID, ID)}, "/")

	if err := h.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListHosts implementation
func (h *HostMgr) ListHosts(args ...interface{}) ([]*model.HostSpec, error) {
	url := strings.Join([]string{
		h.Endpoint,
		urls.GenerateHostURL(urls.Client, h.TenantID)}, "/")

	param, err := processListParam(args)
	if err != nil {
		return nil, err
	}

	if param != "" {
		url += "?" + param
	}

	var res []*model.HostSpec
	if err := h.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateHost implementation
func (h *HostMgr) UpdateHost(ID string, body HostBuilder) (*model.HostSpec, error) {
	var res model.HostSpec
	url := strings.Join([]string{
		h.Endpoint,
		urls.GenerateHostURL(urls.Client, h.TenantID, ID)}, "/")

	if err := h.Recv(url, "PUT", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// DeleteHost implementation
func (h *HostMgr) DeleteHost(ID string) error {
	url := strings.Join([]string{
		h.Endpoint,
		urls.GenerateHostURL(urls.Client, h.TenantID, ID)}, "/")

	return h.Recv(url, "DELETE", nil, nil)
}
