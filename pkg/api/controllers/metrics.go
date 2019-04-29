// Copyright (c) 2019 The OpenSDS Authors All Rights Reserved.
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

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/policy"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/controller/client"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	. "github.com/opensds/opensds/pkg/utils/config"
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

func (m *MetricsPortal) GetMetrics() {
	if !policy.Authorize(m.Ctx, "metrics:get") {
		return
	}
	ctx := c.GetContext(m.Ctx)
	var getMetricSpec = model.GetMetricSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(m.Ctx.Request.Body).Decode(&getMetricSpec); err != nil {
		errMsg := fmt.Sprintf("parse get metric request body failed: %s", err.Error())
		m.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	if err := m.CtrClient.Connect(CONF.OsdsLet.ApiEndpoint); err != nil {
		log.Error("when connecting controller client:", err)
		return
	}
	defer m.CtrClient.Close()

	opt := &pb.GetMetricsOpts{
		InstanceId: getMetricSpec.InstanceId,
		MetricName: getMetricSpec.MetricName,
		StartTime:  getMetricSpec.StartTime,
		EndTime:    getMetricSpec.EndTime,
		Context:    ctx.ToJson(),
	}
	res, err := m.CtrClient.GetMetrics(context.Background(), opt)

	if err != nil {
		log.Error("collect metrics failed in controller service:", err)
		return
	}

	m.SuccessHandle(StatusOK, []byte(res.GetResult().GetMessage()))

	return
}
