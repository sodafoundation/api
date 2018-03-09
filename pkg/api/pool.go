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

	"github.com/astaxie/beego"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/policy"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
)

type PoolPortal struct {
	beego.Controller
}

func (this *PoolPortal) ListPools() {
	if !policy.Authorize(this.Ctx, "pool:list") {
		return
	}
	// Call db api module to handle list pools request.
	result, err := db.C.ListPools(c.GetContext(this.Ctx))
	if err != nil {
		reason := fmt.Sprintf("List pools failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal pools failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *PoolPortal) GetPool() {
	if !policy.Authorize(this.Ctx, "pool:get") {
		return
	}
	id := this.Ctx.Input.Param(":poolId")
	result, err := db.C.GetPool(c.GetContext(this.Ctx), id)
	if err != nil {
		reason := fmt.Sprintf("Get pool failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorBadRequest)
		this.Ctx.Output.Body(model.ErrorBadRequestStatus(reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal pool failed: %s", err.Error())
		this.Ctx.Output.SetStatus(model.ErrorInternalServer)
		this.Ctx.Output.Body(model.ErrorInternalServerStatus(reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}
