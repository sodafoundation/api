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

type PoolMgr struct {
	Receiver

	Endpoint string
	Opt      map[string]string
	lock     sync.Mutex
}

func NewPoolMgr(edp string) *PoolMgr {
	return &PoolMgr{
		Receiver: NewReceiver(),
		Endpoint: edp,
	}
}

func (p *PoolMgr) GetPool(polID string) (*model.StoragePoolSpec, error) {
	var res model.StoragePoolSpec
	url := p.Endpoint + "/api/v1alpha/block/pools/" + polID

	if err := p.Recv(request, url, "GET", p.Opt, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (p *PoolMgr) ListPools() ([]*model.StoragePoolSpec, error) {
	var res []*model.StoragePoolSpec
	url := p.Endpoint + "/api/v1alpha/block/pools"

	if err := p.Recv(request, url, "GET", p.Opt, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func (p *PoolMgr) ResetAndUpdatePoolRequestContent(in interface{}) error {
	var err error

	p.lock.Lock()
	defer p.lock.Unlock()
	// Clear all content stored in Opt field.
	p.Opt = make(map[string]string)
	// Valid the input data.
	switch in.(type) {
	case map[string]string:
		p.Opt = in.(map[string]string)
		break
	default:
		err = errors.New("Request content type not supported")
	}

	return err
}
