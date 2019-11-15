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
	// "strings"

	// log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/policy"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	// "github.com/opensds/opensds/pkg/utils/constants"
)

type ZonePortal struct {
	BasePortal
}

func (p *ZonePortal) ListAvailabilityZones() {
	if !policy.Authorize(p.Ctx, "availability_zone:list") {
		return
	}
	azs, err := db.C.ListZones(c.GetContext(p.Ctx))
	if err != nil {
		errMsg := fmt.Sprintf("get AvailabilityZones failed: %s", err.Error())
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	body, err := json.Marshal(azs)
	if err != nil {
		errMsg := fmt.Sprintf("marshal AvailabilityZones failed: %s", err.Error())
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, body)
	return
}

func (p *ZonePortal) CreateZone() {
	if !policy.Authorize(p.Ctx, "availability_zone:create") {
		return
	}

	var zone = model.ZoneSpec{
		BaseModel: &model.BaseModel{},
	}

	// Unmarshal the request body
	if err := json.NewDecoder(p.Ctx.Request.Body).Decode(&zone); err != nil {
		errMsg := fmt.Sprintf("parse zone request body failed: %v", err)
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// Call db api module to handle create zone request.
	result, err := db.C.CreateZone(c.GetContext(p.Ctx), &zone)
	if err != nil {
		errMsg := fmt.Sprintf("create zone failed: %v", err)
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal zone created result failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, body)
	return
}

func (p *ZonePortal) ListZones() {
	if !policy.Authorize(p.Ctx, "availability_zone:list") {
		return
	}

	m, err := p.GetParameters()
	if err != nil {
		errMsg := fmt.Sprintf("list zones failed: %v", err)
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	result, err := db.C.ListZonesWithFilter(c.GetContext(p.Ctx), m)
	if err != nil {
		errMsg := fmt.Sprintf("list zones failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal zones listed result failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, body)
	return
}

func (p *ZonePortal) GetZone() {
	if !policy.Authorize(p.Ctx, "availability_zone:get") {
		return
	}
	id := p.Ctx.Input.Param(":zoneId")

	result, err := db.C.GetZone(c.GetContext(p.Ctx), id)
	if err != nil {
		errMsg := fmt.Sprintf("zone %s not found: %v", id, err)
		p.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal zone got result failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, body)
	return
}

func (p *ZonePortal) UpdateZone() {

	if !policy.Authorize(p.Ctx, "availability_zone:update") {
		return
	}

	var zone = model.ZoneSpec{
		BaseModel: &model.BaseModel{},
	}
	id := p.Ctx.Input.Param(":zoneId")

	if err := json.NewDecoder(p.Ctx.Request.Body).Decode(&zone); err != nil {
		errMsg := fmt.Sprintf("parse zone request body failed: %v", err)
		p.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	result, err := db.C.UpdateZone(c.GetContext(p.Ctx), id, &zone)
	if err != nil {
		errMsg := fmt.Sprintf("update zones failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal zone updated result failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, body)
	return
}

func (p *ZonePortal) DeleteZone() {

	if !policy.Authorize(p.Ctx, "availability_zone:delete") {
		return
	}
	id := p.Ctx.Input.Param(":zoneId")
	ctx := c.GetContext(p.Ctx)
	zone, err := db.C.GetZone(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("zone %s not found: %v", id, err)
		p.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	err = db.C.DeleteZone(ctx, zone.Id)
	if err != nil {
		errMsg := fmt.Sprintf("delete zones failed: %v", err)
		p.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	p.SuccessHandle(StatusOK, nil)
	return
}
