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
This module implements a entry into the OpenSDS northbound service.

*/

package api

import (
	"encoding/json"
	"fmt"
	c "github.com/opensds/opensds/pkg/context"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/policy"
	"github.com/opensds/opensds/pkg/controller/client"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	. "github.com/opensds/opensds/pkg/utils/config"
	"golang.org/x/net/context"
)

func NewMetricsPortal() *MetricsPortal {
	return &MetricsPortal{
		CtrClient: client.NewClient(),
	}
}

type MetricsPortal struct {
	BasePortal

	CtrClient client.Client
}

func (m *MetricsPortal) CollectMetrics() {
	if !policy.Authorize(m.Ctx, "metrics:collect") {
		return
	}
	ctx := c.GetContext(m.Ctx)
	var collMetricSpec = model.CollectMetricSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(m.Ctx.Request.Body).Decode(&collMetricSpec); err != nil {
		errMsg := fmt.Sprintf("parse collect metric request body failed: %s", err.Error())
		m.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	m.SuccessHandle(StatusOK, nil)

	// NOTE:The real volume creation process.
	// Volume creation request is sent to the Dock. Dock will update volume status to "available"
	// after volume creation is completed.
	if err := m.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer m.CtrClient.Close()

	opt := &pb.CollectMetricsOpts{
		InstanceId: 	collMetricSpec.InstanceId,
		MetricValues:   collMetricSpec.Metrics,
		Context: 		ctx.ToJson(),
	}
	if _, err := m.CtrClient.CollectMetrics(context.Background(), opt); err != nil {
		log.Error("collect metrics failed in controller service:", err)
		return
	}

	return
}


