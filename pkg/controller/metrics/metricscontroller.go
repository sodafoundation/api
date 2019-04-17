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

/*
This module implements a entry into the OpenSDS volume controller service.

*/

package metrics

import (
	"encoding/json"
	//"errors"
	"fmt"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/dock/client"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"golang.org/x/net/context"
)

// Controller is an interface for exposing some operations of different volume
// controllers.
type Controller interface {
	CollectMetrics(opt *pb.CollectMetricsOpts) (*model.CollectMetricSpec, error)
	SetDock(dockInfo *model.DockSpec)
}

// NewController method creates a controller structure and expose its pointer.
func NewController() Controller {
	return &controller{
		Client: client.NewClient(),
	}
}

type controller struct {
	client.Client
	DockInfo *model.DockSpec
}

func (c *controller) CollectMetrics(opt *pb.CollectMetricsOpts) (*model.CollectMetricSpec, error) {
	if err := c.Client.Connect(c.DockInfo.Endpoint); err != nil {
		log.Error("when connecting dock client:", err)
		return nil, err
	}

	response, err := c.Client.CollectMetrics(context.Background(), opt)
	if err != nil {
		log.Error("create volume failed in volume controller:", err)
		return nil, err
	}
	defer c.Client.Close()

	if errorMsg := response.GetError(); errorMsg != nil {
		return nil,
			fmt.Errorf("failed to create volume in volume controller, code: %v, message: %v",
				errorMsg.GetCode(), errorMsg.GetDescription())
	}

	var vol = &model.CollectMetricSpec{}
	if err = json.Unmarshal([]byte(response.GetResult().GetMessage()), vol); err != nil {
		log.Error("create volume failed in volume controller:", err)
		return nil, err
	}

	return vol, nil

}

func (c *controller) SetDock(dockInfo *model.DockSpec) {
	c.DockInfo = dockInfo
}
