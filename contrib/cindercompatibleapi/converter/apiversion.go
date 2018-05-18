// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package converter

import (
	"github.com/opensds/opensds/pkg/model"
)

// *******************List All Api Versions*******************

// ListAllApiVersionsRespSpec ...
type ListAllApiVersionsRespSpec struct {
	Versions []ListAllApiVersions `json:"versions"`
}

// ListAllApiVersions ...
type ListAllApiVersions struct {
	Status     string            `json:"status"`
	Updated    string            `json:"updated"`
	Links      []VersionLink     `json:"links"`
	MinVersion string            `json:"min_version,"`
	Version    string            `json:"version"`
	MediaTypes map[string]string `json:"media-types"`
	ID         string            `json:"id"`
}

// VersionLink ...
type VersionLink struct {
	Href string `json:"href"`
	Type string `json:"type"`
	Rel  string `json:"rel"`
}

// ListAllApiVersionsResp ...
func ListAllApiVersionsResp(versions []*model.VersionSpec) *ListAllApiVersionsRespSpec {
	var resp ListAllApiVersionsRespSpec
	var cinderVersion ListAllApiVersions

	if 0 == len(versions) {
		resp.Versions = make([]ListAllApiVersions, 0, 0)
	} else {
		for _, version := range versions {

			cinderVersion.Status = version.Status
			cinderVersion.Updated = version.UpdatedAt
			cinderVersion.MinVersion = "3.0"
			cinderVersion.ID = "v3.0"

			resp.Versions = append(resp.Versions, cinderVersion)

		}
	}

	return &resp
}
