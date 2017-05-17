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
This module implements the policy-based scheduling by parsing storage
profiles configured by admin.

*/

package scheduler

import (
	dockRoute "github.com/opensds/opensds/pkg/controller/metadata/dock_route"
)

func getDockId(backend string) string {
	// Get all docks information.
	docks, err := dockRoute.ListDocks()
	if err != nil {
		return ""
	}

	for _, dock := range docks.DockList {
		for _, value := range dock.Backends {
			if value == backend {
				return dock.Id
			}
		}
	}
	return ""
}
