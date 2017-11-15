// Copyright 2017 The OpenSDS Authors.
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
This module implements a entry into the OpenSDS northbound REST service.

*/

package api

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/utils"
)

var KnownVersions = []map[string]string{
	{
		"name":        "v1alpha",
		"description": "v1alpha version",
		"status":      "CURRENT",
		"updatedAt":   "2017-07-10T14:36:58.014Z",
	},
}

type VersionPortal struct {
	beego.Controller
}

func (this *VersionPortal) ListVersions() {
	body, err := json.Marshal(KnownVersions)
	if err != nil {
		reason := fmt.Sprintf("Marshal versions failed: %s", err.Error())
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}

func (this *VersionPortal) GetVersion() {
	apiVersion := this.Ctx.Input.Param(":apiVersion")

	// Find version by specified api version
	var result map[string]string
	for _, version := range KnownVersions {
		if version["name"] == apiVersion {
			result = version
			break
		}
	}
	if result == nil {
		reason := fmt.Sprintf("Can't find the version: %s", apiVersion)
		this.Ctx.Output.SetStatus(StatusNotFound)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		reason := fmt.Sprintf("Marshal version failed: %s", err.Error())
		this.Ctx.Output.SetStatus(StatusInternalServerError)
		this.Ctx.Output.Body(utils.ErrorStatus(this.Ctx.Output.Status, reason))
		log.Error(reason)
		return
	}

	this.Ctx.Output.SetStatus(StatusOK)
	this.Ctx.Output.Body(body)
	return
}
