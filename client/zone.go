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

// ZoneBuilder contains request body of handling a zone request.
// Currently it's assigned as the pointer of ZoneSpec struct, but it
// could be discussed if it's better to define an interface.
type ZoneBuilder *model.ZoneSpec

// NewZoneMgr
func NewZoneMgr(r Receiver, edp string, tenantId string) *ZoneMgr {
	return &ZoneMgr{
		Receiver: r,
		Endpoint: edp,
		TenantId: tenantId,
	}
}

// ZoneMgr
type ZoneMgr struct {
	Receiver
	Endpoint string
	TenantId string
}

// CreateZone
func (p *ZoneMgr) CreateZone(body ZoneBuilder) (*model.ZoneSpec, error) {
	var res model.ZoneSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateZoneURL(urls.Client, p.TenantId)}, "/")

	if err := p.Recv(url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetZone
func (p *ZoneMgr) GetZone(zoneID string) (*model.ZoneSpec, error) {
	var res model.ZoneSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateZoneURL(urls.Client, p.TenantId, zoneID)}, "/")

	if err := p.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// UpdateZone ...
func (p *ZoneMgr) UpdateZone(zoneID string, body ZoneBuilder) (*model.ZoneSpec, error) {
	var res model.ZoneSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateZoneURL(urls.Client, p.TenantId, zoneID)}, "/")

	if err := p.Recv(url, "PUT", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListZones
func (p *ZoneMgr) ListZones(args ...interface{}) ([]*model.ZoneSpec, error) {
	var res []*model.ZoneSpec

	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateZoneURL(urls.Client, p.TenantId)}, "/")

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

// DeleteZone
func (p *ZoneMgr) DeleteZone(zoneID string) error {
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateZoneURL(urls.Client, p.TenantId, zoneID)}, "/")

	return p.Recv(url, "DELETE", nil, nil)
}
