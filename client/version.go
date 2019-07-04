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

package client

import (
	"strings"

	"github.com/opensds/opensds/pkg/model"
)

// VersionBuilder contains request body of handling a version request.
// Currently it's assigned as the pointer of VersionSpec struct, but it
// could be discussed if it's better to define an interface.
type VersionBuilder *model.VersionSpec

// NewVersionMgr ...
func NewVersionMgr(r Receiver, edp string, tenantId string) *VersionMgr {
	return &VersionMgr{
		Receiver: r,
		Endpoint: edp,
	}
}

// VersionMgr ...
type VersionMgr struct {
	Receiver
	Endpoint string
	tenantId string
}

// GetVersion ...
func (v *VersionMgr) GetVersion(apiVersion string) (*model.VersionSpec, error) {
	var res model.VersionSpec
	url := strings.Join([]string{
		v.Endpoint, apiVersion}, "/")

	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListVersions ...
func (v *VersionMgr) ListVersions() ([]*model.VersionSpec, error) {
	var res []*model.VersionSpec
	url := v.Endpoint

	if err := v.Recv(url, "GET", nil, &res); err != nil {
		return nil, err
	}

	return res, nil
}
