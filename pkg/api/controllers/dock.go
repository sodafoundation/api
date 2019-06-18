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

// DockPortal
type DockPortal struct {
	BasePortal
}

// ListDocks
func (d *DockPortal) ListDocks() {
	if !policy.Authorize(d.Ctx, "dock:list") {
		return
	}
	// Call db api module to handle list docks request.
	m, err := d.GetParameters()
	if err != nil {
		errMsg := fmt.Sprintf("list docks failed: %s", err.Error())
		d.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}
	result, err := db.C.ListDocksWithFilter(c.GetContext(d.Ctx), m)
	if err != nil {
		errMsg := fmt.Sprintf("list docks failed: %s", err.Error())
		d.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal docks failed: %s", err.Error())
		d.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	d.SuccessHandle(StatusOK, body)
	return
}

// GetDock
func (d *DockPortal) GetDock() {
	if !policy.Authorize(d.Ctx, "dock:get") {
		return
	}
	id := d.Ctx.Input.Param(":dockId")
	result, err := db.C.GetDock(c.GetContext(d.Ctx), id)
	if err != nil {
		errMsg := fmt.Sprintf("dock %s not found: %s", id, err.Error())
		d.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal dock failed: %s", err.Error())
		d.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	d.SuccessHandle(StatusOK, body)
	return
}
