// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package client

import (
	"errors"
	"fmt"
	"sync"

	"github.com/opensds/opensds/pkg/model"
)

type ProfileMgr struct {
	Receiver

	Endpoint string
	Opt      map[string]string
	Body     interface{}
	lock     sync.Mutex
}

func NewProfileMgr(edp string) *ProfileMgr {
	return &ProfileMgr{
		Receiver: NewReceiver(),
		Endpoint: edp,
	}
}

func (p *ProfileMgr) CreateProfile() (*model.ProfileSpec, error) {
	var res model.ProfileSpec
	url := p.Endpoint + "/api/v1alpha/profiles"

	if err := p.Recv(request, url, "POST", p.Body, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (p *ProfileMgr) GetProfile(prfID string) (*model.ProfileSpec, error) {
	var res model.ProfileSpec
	url := p.Endpoint + "/api/v1alpha/profiles/" + prfID

	if err := p.Recv(request, url, "GET", p.Opt, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (p *ProfileMgr) ListProfiles() ([]*model.ProfileSpec, error) {
	var res []*model.ProfileSpec
	url := p.Endpoint + "/api/v1alpha/profiles"

	if err := p.Recv(request, url, "GET", p.Opt, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func (p *ProfileMgr) DeleteProfile(prfID string) *model.Response {
	var res model.Response
	url := p.Endpoint + "/api/v1alpha/profiles/" + prfID

	if err := p.Recv(request, url, "DELETE", p.Body, &res); err != nil {
		res.Status, res.Error = "Failure", fmt.Sprint(err)
	}

	return &res
}

func (p *ProfileMgr) AddExtraProperty(prfID string, ext *model.ExtraSpec) (*model.ExtraSpec, error) {
	var res model.ExtraSpec
	url := p.Endpoint + "/api/v1alpha/profiles/" + prfID + "/extras"

	if err := p.Recv(request, url, "POST", p.Body, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (p *ProfileMgr) ListExtraProperties(prfID string) (*model.ExtraSpec, error) {
	var res model.ExtraSpec
	url := p.Endpoint + "/api/v1alpha/profiles/" + prfID + "/extras"

	if err := p.Recv(request, url, "GET", p.Opt, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (p *ProfileMgr) RemoveExtraProperty(prfID, extraKey string) *model.Response {
	var res model.Response
	url := p.Endpoint + "/api/v1alpha/profiles/" + prfID + "/extras/" + extraKey

	if err := p.Recv(request, url, "DELETE", p.Body, &res); err != nil {
		res.Status, res.Error = "Failure", fmt.Sprint(err)
	}

	return &res
}

func (p *ProfileMgr) ResetAndUpdateProfileRequestContent(in interface{}) error {
	var err error

	p.lock.Lock()
	defer p.lock.Unlock()
	// Clear all content stored in Opt field.
	p.Opt, p.Body = make(map[string]string), nil
	// Valid the input data.
	switch in.(type) {
	case map[string]string:
		p.Opt = in.(map[string]string)
		break
	case model.ProfileSpec, *model.ProfileSpec:
		p.Body = in
		break
	default:
		err = errors.New("Request content type not supported")
	}

	return err
}
