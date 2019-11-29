// Copyright 2019 The OpenSDS Authors.
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
	"encoding/json"
	"fmt"

	"github.com/opensds/opensds/pkg/api/policy"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
)

type HostPortal struct {
	BasePortal
}

func NewHostPortal() *HostPortal {
	return &HostPortal{}
}

func (p *HostPortal) ListHosts() {
	if !policy.Authorize(p.Ctx, "host:list") {
		return
	}
	// TODO: handle query parameters
	hosts, err := db.C.ListHosts(c.GetContext(p.Ctx))
	if err != nil {
		errMsg := fmt.Sprintf("list hosts failed: %s", err.Error())
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	body, err := json.Marshal(hosts)
	if err != nil {
		errMsg := fmt.Sprintf("marshal hosts failed: %s", err.Error())
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, body)
	return
}

func (p *HostPortal) CreateHost() {
	if !policy.Authorize(p.Ctx, "host:create") {
		return
	}

	var host = model.HostSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(p.Ctx.Request.Body).Decode(&host); err != nil {
		errMsg := fmt.Sprintf("parse host request body failed: %s", err.Error())
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// TODO: Add prameter validation

	result, err := db.C.CreateHost(c.GetContext(p.Ctx), &host)
	if err != nil {
		errMsg := fmt.Sprintf("create host failed: %v", err)
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal host created result failed: %s", err.Error())
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, body)
	return
}

func (p *HostPortal) GetHost() {
	if !policy.Authorize(p.Ctx, "host:get") {
		return
	}
	id := p.Ctx.Input.Param(":hostId")
	result, err := db.C.GetHost(c.GetContext(p.Ctx), id)
	if err != nil {
		errMsg := fmt.Sprintf("host %s not found: %s", id, err.Error())
		p.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal host failed: %s", err.Error())
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, body)
	return
}

func (p *HostPortal) UpdateHost() {
	if !policy.Authorize(p.Ctx, "host:update") {
		return
	}

	id := p.Ctx.Input.Param(":hostId")
	var host = model.HostSpec{
		BaseModel: &model.BaseModel{
			Id: id,
		},
	}
	if err := json.NewDecoder(p.Ctx.Request.Body).Decode(&host); err != nil {
		errMsg := fmt.Sprintf("parse host request body failed: %v", err)
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// TODO: Add parameter validation

	result, err := db.C.UpdateHost(c.GetContext(p.Ctx), &host)
	if err != nil {
		errMsg := fmt.Sprintf("update host failed: %v", err)
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal host updated result failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, body)
	return
}

func (p *HostPortal) DeleteHost() {
	if !policy.Authorize(p.Ctx, "host:delete") {
		return
	}
	id := p.Ctx.Input.Param(":hostId")

	// TODO: Add deletion constraition check

	err := db.C.DeleteHost(c.GetContext(p.Ctx), id)
	if err != nil {
		errMsg := fmt.Sprintf("delete host failed: %v", err)
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, nil)
	return
}
