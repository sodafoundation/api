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

// ProfileBuilder contains request body of handling a profile request.
// Currently it's assigned as the pointer of ProfileSpec struct, but it
// could be discussed if it's better to define an interface.
type ProfileBuilder *model.ProfileSpec

// CustomBuilder contains request body of handling a profile customized
// properties request.
// Currently it's assigned as the pointer of CustomPropertiesSpec struct, but it
// could be discussed if it's better to define an interface.
type CustomBuilder *model.CustomPropertiesSpec

// NewProfileMgr
func NewProfileMgr(r Receiver, edp string, tenantId string) *ProfileMgr {
	return &ProfileMgr{
		Receiver: r,
		Endpoint: edp,
		TenantId: tenantId,
	}
}

// ProfileMgr
type ProfileMgr struct {
	Receiver
	Endpoint string
	TenantId string
}

// CreateProfile
func (p *ProfileMgr) CreateProfile(body ProfileBuilder) (*model.ProfileSpec, error) {
	var res model.ProfileSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(urls.Client, p.TenantId)}, "/")

	if err := p.Recv(url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetProfile
func (p *ProfileMgr) GetProfile(prfID string) (*model.ProfileSpec, error) {
	var res model.ProfileSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(urls.Client, p.TenantId, prfID)}, "/")

	if err := p.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// UpdateProfile ...
func (p *ProfileMgr) UpdateProfile(prfID string, body ProfileBuilder) (*model.ProfileSpec, error) {
	var res model.ProfileSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(urls.Client, p.TenantId, prfID)}, "/")

	if err := p.Recv(url, "PUT", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListProfiles
func (p *ProfileMgr) ListProfiles(args ...interface{}) ([]*model.ProfileSpec, error) {
	var res []*model.ProfileSpec

	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(urls.Client, p.TenantId)}, "/")

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

// DeleteProfile
func (p *ProfileMgr) DeleteProfile(prfID string) error {
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(urls.Client, p.TenantId, prfID)}, "/")

	return p.Recv(url, "DELETE", nil, nil)
}

// AddCustomProperty
func (p *ProfileMgr) AddCustomProperty(prfID string, body CustomBuilder) (*model.CustomPropertiesSpec, error) {
	var res model.CustomPropertiesSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(urls.Client, p.TenantId, prfID),
		"customProperties"}, "/")

	if err := p.Recv(url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListCustomProperties
func (p *ProfileMgr) ListCustomProperties(prfID string) (*model.CustomPropertiesSpec, error) {
	var res model.CustomPropertiesSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(urls.Client, p.TenantId, prfID),
		"customProperties"}, "/")

	if err := p.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// RemoveCustomProperty
func (p *ProfileMgr) RemoveCustomProperty(prfID, customKey string) error {
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(urls.Client, p.TenantId, prfID),
		"customProperties", customKey}, "/")

	return p.Recv(url, "DELETE", nil, nil)
}
