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

// AvailabilityZoneBuilder contains request body of handling a zone request.
// Currently it's assigned as the pointer of AvailabilityZoneSpec struct, but it
// could be discussed if it's better to define an interface.
type AvailabilityZoneBuilder *model.AvailabilityZoneSpec

// NewAvailabilityZoneMgr
func NewAvailabilityZoneMgr(r Receiver, edp string, tenantId string) *AvailabilityZoneMgr {
	return &AvailabilityZoneMgr{
		Receiver: r,
		Endpoint: edp,
		TenantId: tenantId,
	}
}

// AvailabilityZoneMgr
type AvailabilityZoneMgr struct {
	Receiver
	Endpoint string
	TenantId string
}

// CreateAvailabilityZone
func (p *AvailabilityZoneMgr) CreateAvailabilityZone(body AvailabilityZoneBuilder) (*model.AvailabilityZoneSpec, error) {
	var res model.AvailabilityZoneSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateZoneURL(urls.Client, p.TenantId)}, "/")

	if err := p.Recv(url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetAvailabilityZone
func (p *AvailabilityZoneMgr) GetAvailabilityZone(zoneID string) (*model.AvailabilityZoneSpec, error) {
	var res model.AvailabilityZoneSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateZoneURL(urls.Client, p.TenantId, zoneID)}, "/")

	if err := p.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// UpdateAvailabilityZone ...
func (p *AvailabilityZoneMgr) UpdateAvailabilityZone(zoneID string, body AvailabilityZoneBuilder) (*model.AvailabilityZoneSpec, error) {
	var res model.AvailabilityZoneSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateZoneURL(urls.Client, p.TenantId, zoneID)}, "/")

	if err := p.Recv(url, "PUT", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListAvailabilityZones
func (p *AvailabilityZoneMgr) ListAvailabilityZones(args ...interface{}) ([]*model.AvailabilityZoneSpec, error) {
	var res []*model.AvailabilityZoneSpec

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

// DeleteAvailabilityZone
func (p *AvailabilityZoneMgr) DeleteAvailabilityZone(zoneID string) error {
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateZoneURL(urls.Client, p.TenantId, zoneID)}, "/")

	return p.Recv(url, "DELETE", nil, nil)
}
