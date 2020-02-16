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
This module implements the common data structure.

*/

package model

import (
	"encoding/json"

	"github.com/golang/glog"
)

// An OpenSDS availability zone is identified by a unique name and ID.
type AvailabilityZoneSpec struct {
	*BaseModel

	// The name of the availability zone.
	Name string `json:"name,omitempty"`

	// The description of the availability zone.
	// +optional
	Description string `json:"description,omitempty"`
}

func NewZoneFromJson(s string) *AvailabilityZoneSpec {
	p := &AvailabilityZoneSpec{}
	err := json.Unmarshal([]byte(s), p)
	if err != nil {
		glog.Errorf("Unmarshal json to AvailabilityZoneSpec failed, %v", err)
	}
	return p
}

func (p *AvailabilityZoneSpec) ToJson() string {
	b, err := json.Marshal(p)
	if err != nil {
		glog.Errorf("AvailabilityZoneSpec convert to json failed, %v", err)
	}
	return string(b)
}
