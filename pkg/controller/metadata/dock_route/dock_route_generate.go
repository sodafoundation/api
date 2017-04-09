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
This module defines the initialization work about dock service: generate dock
node information and record it to dock route service.

*/

package dock_route

import (
	"log"
	"os/exec"

	"github.com/opensds/opensds/pkg/controller/api"
)

func GenerateDockRoute(dockIp string) (api.DockRoute, error) {
	dockId, err := generateDockId()
	if err != nil {
		log.Println("Could not generate dock id:", err)
		return api.DockRoute{}, err
	}

	return api.DockRoute{
		Id:      dockId,
		Status:  "available",
		Address: dockIp,
	}, nil
}

func generateDockId() (string, error) {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", err
	}

	dockId := string(out)
	return dockId[0 : len(dockId)-1], nil
}
