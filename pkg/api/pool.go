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

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/policy"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
)

type PoolPortal struct {
	BasePortal
}

func (p *PoolPortal) ListAvailabilityZones() {
	if !policy.Authorize(p.Ctx, "availability_zone:list") {
		return
	}
	azs, err := db.C.ListAvailabilityZones(c.GetContext(p.Ctx))
	if err != nil {
		reason := fmt.Sprintf("Get AvailabilityZones for pools failed: %s", err.Error())
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	body, err := json.Marshal(azs)
	if err != nil {
		reason := fmt.Sprintf("Marshal AvailabilityZones failed: %s", err.Error())
		p.Ctx.Output.SetStatus(model.ErrorInternalServer)
		p.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	p.Ctx.Output.SetStatus(StatusOK)
	p.Ctx.Output.Body(body)
	return
}

func (p *PoolPortal) ListPools() {
	if !policy.Authorize(p.Ctx, "pool:list") {
		return
	}
	// Call db api module to handle list pools request.
	m, err := p.GetParameters()
	if err != nil {
		reason := fmt.Sprintf("List pools failed: %s", err.Error())
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	result, err := db.C.ListPoolsWithFilter(c.GetContext(p.Ctx), m)
	if err != nil {
		reason := fmt.Sprintf("List pools failed: %s", err.Error())
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal pools failed: %s", err.Error())
		p.Ctx.Output.SetStatus(model.ErrorInternalServer)
		p.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	p.Ctx.Output.SetStatus(StatusOK)
	p.Ctx.Output.Body(body)
	return
}

func (p *PoolPortal) GetPool() {
	if !policy.Authorize(p.Ctx, "pool:get") {
		return
	}
	id := p.Ctx.Input.Param(":poolId")
	result, err := db.C.GetPool(c.GetContext(p.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("Get pool failed: %s", err.Error())
		p.Ctx.Output.SetStatus(model.ErrorBadRequest)
		p.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal pool failed: %s", err.Error())
		p.Ctx.Output.SetStatus(model.ErrorInternalServer)
		p.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	p.Ctx.Output.SetStatus(StatusOK)
	p.Ctx.Output.Body(body)
	return
}
