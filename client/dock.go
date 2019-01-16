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
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/urls"
)

func NewDockMgr(r Receiver, edp string) *DockMgr {
	return &DockMgr{
		Receiver: r,
		Endpoint: edp,
	}
}

type DockMgr struct {
	Receiver
	Endpoint string
}

func (d *DockMgr) GetDock(dckID string) (*model.DockSpec, error) {
	var res model.DockSpec
	if err := d.Receiver.Recv(urls.DockResource, d.Endpoint, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (d *DockMgr) ListDocks(args ...interface{}) ([]*model.DockSpec, error) {
	var res []*model.DockSpec
	var filter string
	param, err := processListParam(args)
	if err != nil {
		return nil, err
	}

	if param != "" {
		filter += "?" + param
	}

	if err := d.Receiver.Recv(urls.DockResource, d.Endpoint, "GET", nil, &res, filter); err != nil {
		return nil, err
	}

	return res, nil
}
