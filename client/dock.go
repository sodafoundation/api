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

type DockMgr struct {
	Receiver

	Endpoint string
	Opt      map[string]string
	lock     sync.Mutex
}

func NewDockMgr(edp string) *DockMgr {
	return &DockMgr{
		Receiver: NewReceiver(),
		Endpoint: edp,
	}
}

func (d *DockMgr) GetDock(dckID string) (*model.DockSpec, error) {
	var res model.DockSpec
	url := d.Endpoint + "/api/v1alpha/docks/" + dckID

	if err := d.Recv(request, url, "GET", d.Opt, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (d *DockMgr) ListDocks() ([]*model.DockSpec, error) {
	var res []*model.DockSpec
	url := d.Endpoint + "/api/v1alpha/docks"

	if err := d.Recv(request, url, "GET", d.Opt, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func (d *DockMgr) ResetAndUpdateDockRequestContent(in interface{}) error {
	var err error

	d.lock.Lock()
	defer d.lock.Unlock()
	// Clear all content stored in Opt field.
	d.Opt = make(map[string]string)
	// Valid the input data.
	switch in.(type) {
	case map[string]string:
		d.Opt = in.(map[string]string)
		break
	default:
		err = errors.New("Request content type not supported")
	}

	return err
}
