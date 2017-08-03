// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

/*
This module implements the common data structure.

*/

package model

type QuotaSpec struct {
	*BaseModel
	Name         string           `json:"name,omitempty"`
	Description  string           `json:"description,omitempty"`
	ResourceList map[string]int64 `json:"resourceList,omitempty"`
}

func (quota *QuotaSpec) GetName() string {
	return quota.Name
}

func (quota *QuotaSpec) GetDescription() string {
	return quota.Description
}

func (quota *QuotaSpec) GetResourceList() map[string]int64 {
	return quota.ResourceList
}
