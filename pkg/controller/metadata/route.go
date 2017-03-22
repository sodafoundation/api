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
This module defines the route configuation about dock service.

*/

package metadata

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os/exec"
)

type DockRoute struct {
	Id      string   `json:"id"`
	Drivers []string `json:"drivers"`
	Address string   `json:"address"`
}

func getDockRoutes() ([]DockRoute, error) {
	var routes []DockRoute

	userJSON, err := ioutil.ReadFile("/etc/opensds/dock_route.json")
	if err != nil {
		log.Println("ReadFile json failed:", err)
		return routes, err
	}

	if err = json.Unmarshal(userJSON, &routes); err != nil {
		log.Println("Unmarshal json failed:", err)
		return routes, err
	}
	return routes, nil
}

func GetDockAddress(dockId string) (string, error) {
	routes, err := getDockRoutes()
	if err != nil {
		log.Println("Could not get dock routes:", err)
		return "", err
	}

	if dockId == "" {
		log.Println("Dock id not provided, arrange the address randomly!")
		return routes[0].Address, nil
	}

	for _, i := range routes {
		if dockId == i.Id {
			return i.Address, nil
		}
	}

	err = errors.New("Could not find this dock service!")
	return "", err
}

func AddDockRoute(address string, drivers []string) (bool, error) {
	routes, err := getDockRoutes()
	if err != nil {
		log.Println("Could not get dock routes:", err)
		return false, err
	}

	uuid, err := generateDockId()
	if err != nil {
		log.Println("Generate new dock id failed:", err)
		return false, nil
	}

	newRoute := DockRoute{
		Id:      uuid,
		Address: address,
	}
	for _, driver := range drivers {
		newRoute.Drivers = append(newRoute.Drivers, driver)
	}

	routes = append(routes, newRoute)

	if !setDockRoutes(&routes) {
		err = errors.New("Add dock route failed!")
		return false, err
	} else {
		return true, nil
	}
}

func generateDockId() (string, error) {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func setDockRoutes(routes *[]DockRoute) bool {
	body, err := json.Marshal(routes)
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
