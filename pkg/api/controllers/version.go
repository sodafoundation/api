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
This module implements a entry into the OpenSDS northbound REST service.

*/

package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/opensds/opensds/pkg/model"
)

// KnownVersions
var KnownVersions = []map[string]string{
	{
		"name":        "v1beta",
		"description": "v1beta version",
		"status":      "CURRENT",
		"updatedAt":   "2017-07-10T14:36:58.014Z",
	},
}

// VersionPortal
type VersionPortal struct {
	BasePortal
}

// ListVersions
func (v *VersionPortal) ListVersions() {
	body, err := json.Marshal(KnownVersions)
	if err != nil {
		errMsg := fmt.Sprintf("marshal versions failed: %s", err.Error())
		v.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	v.SuccessHandle(StatusOK, body)
	return
}

// GetVersion
func (v *VersionPortal) GetVersion() {
	apiVersion := v.Ctx.Input.Param(":apiVersion")

	// Find version by specified api version
	var result map[string]string
	for _, version := range KnownVersions {
		if version["name"] == apiVersion {
			result = version
			break
		}
	}
	if result == nil {
		errMsg := fmt.Sprintf("can't find the version: %s", apiVersion)
		v.ErrorHandle(model.ErrorNotFound, errMsg)
		return
	}

	// Marshal the result.
	body, err := json.Marshal(result)
	if err != nil {
		errMsg := fmt.Sprintf("marshal version failed: %s", err.Error())
		v.ErrorHandle(model.ErrorInternalServer, errMsg)
		return
	}

	v.SuccessHandle(StatusOK, body)
	return
}
