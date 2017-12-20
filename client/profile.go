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

// ExtraBuilder contains request body of handling a profile extra request.
// Currently it's assigned as the pointer of Extra struct, but it
// could be discussed if it's better to define an interface.
type ExtraBuilder *model.ExtraSpec

func NewProfileMgr(edp string) *ProfileMgr {
	return &ProfileMgr{
		Receiver: NewReceiver(),
		Endpoint: edp,
	}
}

type ProfileMgr struct {
	Receiver

	Endpoint string
}

func (p *ProfileMgr) CreateProfile(body ProfileBuilder) (*model.ProfileSpec, error) {
	var res model.ProfileSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL()}, "/")

	if err := p.Recv(request, url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (p *ProfileMgr) GetProfile(prfID string) (*model.ProfileSpec, error) {
	var res model.ProfileSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(prfID)}, "/")

	if err := p.Recv(request, url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (p *ProfileMgr) ListProfiles() ([]*model.ProfileSpec, error) {
	var res []*model.ProfileSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL()}, "/")

	if err := p.Recv(request, url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (p *ProfileMgr) DeleteProfile(prfID string) error {
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(prfID)}, "/")

	return p.Recv(request, url, "DELETE", nil, nil)
}

func (p *ProfileMgr) AddExtraProperty(prfID string, body ExtraBuilder) (*model.ExtraSpec, error) {
	var res model.ExtraSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(prfID),
		"extras"}, "/")

	if err := p.Recv(request, url, "POST", body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (p *ProfileMgr) ListExtraProperties(prfID string) (*model.ExtraSpec, error) {
	var res model.ExtraSpec
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(prfID),
		"extras"}, "/")

	if err := p.Recv(request, url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (p *ProfileMgr) RemoveExtraProperty(prfID, extraKey string) error {
	url := strings.Join([]string{
		p.Endpoint,
		urls.GenerateProfileURL(prfID),
		"extras", extraKey}, "/")

	return p.Recv(request, url, "DELETE", nil, nil)
}
