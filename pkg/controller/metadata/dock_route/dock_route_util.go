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
node information and record it to dock route service. Please DO NOT
modify this file.

*/

package dock_route

import (
	"encoding/json"
	"io/ioutil"
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
	"github.com/satori/go.uuid"
)

func generateDockRoute(dockIp string) api.DockRoute {
	return api.DockRoute{
		Id:      uuid.NewV4().String(),
		Status:  "available",
		Address: dockIp,
	}
}

func readDockRoutesFromFile() ([]api.DockRoute, error) {
	var routes []api.DockRoute

	userJSON, err := ioutil.ReadFile("/etc/opensds/dock_route.json")
	if err != nil {
		log.Println("ReadFile json failed:", err)
		return routes, err
	}

	// If the dock route table is empty, consider it as a normal condition
	if len(userJSON) == 0 {
		return routes, nil
	} else if err = json.Unmarshal(userJSON, &routes); err != nil {
		log.Println("Unmarshal json failed:", err)
		return routes, err
	}
	return routes, nil
}

func writeDockRoutesToFile(routes []api.DockRoute) bool {
	body, err := json.Marshal(&routes)
	if err != nil {
		log.Println("Marshal json failed:", err)
		return false
	}

	if err = ioutil.WriteFile("/etc/opensds/dock_route.json", body, 0666); err != nil {
		log.Println("WriteFile json failed:", err)
		return false
	}
	return true
}
