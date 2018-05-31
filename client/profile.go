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

// ProfileBuilder contains request body of handling a profile request.
// Currently it's assigned as the pointer of ProfileSpec struct, but it
// could be discussed if it's better to define an interface.
type ProfileBuilder *model.ProfileSpec

// ExtraBuilder contains request body of handling a profile extra request.
// Currently it's assigned as the pointer of Extra struct, but it
// could be discussed if it's better to define an interface.
type ExtraBuilder *model.ExtraSpec

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
func (p *ProfileMgr) ListProfiles(v []string, prof *model.ProfileSpec) ([]*model.ProfileSpec, error) {
	var res []*model.ProfileSpec
	var u string

	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(urls.Client, p.TenantId)}, "/")

	var limit, offset, sortDir, sortKey, createdAt, description, name, storageType, updatedAt, id string
	var urlpara []string

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
	if prof.CreatedAt != "" {
		createdAt = "CreatedAt=" + prof.CreatedAt
		urlpara = append(urlpara, createdAt)
	}
	if prof.Description != "" {
		description = "Description=" + prof.Description
		urlpara = append(urlpara, description)
	}
	if prof.Name != "" {
		name = "Name=" + prof.Name
		urlpara = append(urlpara, name)
	}
	if prof.StorageType != "" {
		storageType = "StorageType=" + prof.StorageType
		urlpara = append(urlpara, storageType)
	}
	if prof.UpdatedAt != "" {
		updatedAt = "UpdatedAt=" + prof.UpdatedAt
		urlpara = append(urlpara, updatedAt)
	}
	if prof.Id != "" {
		id = "Id=" + prof.Id
		urlpara = append(urlpara, id)
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

// DeleteProfile
func (p *ProfileMgr) DeleteProfile(prfID string) error {
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(urls.Client, p.TenantId, prfID)}, "/")

	return p.Recv(url, "DELETE", nil, nil)
}

// AddExtraProperty
func (p *ProfileMgr) AddExtraProperty(prfID string, body ExtraBuilder) (*model.ExtraSpec, error) {
	var res model.ExtraSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(urls.Client, p.TenantId, prfID),
		"extras"}, "/")

	if err := p.Recv(url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListExtraProperties
func (p *ProfileMgr) ListExtraProperties(prfID string) (*model.ExtraSpec, error) {
	var res model.ExtraSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(urls.Client, p.TenantId, prfID),
		"extras"}, "/")

	if err := p.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// RemoveExtraProperty
func (p *ProfileMgr) RemoveExtraProperty(prfID, extraKey string) error {
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(urls.Client, p.TenantId, prfID),
		"extras", extraKey}, "/")

	return p.Recv(url, "DELETE", nil, nil)
}
