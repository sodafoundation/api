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

package dock

import (
	"errors"
	"log"

	api "github.com/opensds/opensds/pkg/model"
)

/*
func RegisterDock(edp string, backends []string) (*api.DockSpec, error) {
	docks, err := readDocksFromFile()
	if err != nil {
		log.Println("Could not get dock resource:", err)
		return &api.DockSpec{}, err
	}

	docks = append(docks, dock)

	if !writeDocksToFile(docks) {
		err = errors.New("Register dock endpoint " + edp + " failed!")
		return &api.DockSpec{}, err
	} else {
		return &dock, nil
	}
}

func DeregisterDock(edp string) (string, error) {
	docks, err := readDocksFromFile()
	if err != nil {
		log.Println("Could not get dock resource:", err)
		return "", err
	}

	var dockFound bool
	var newDocks []api.DockSpec

	for i, dock := range docks {
		if dock.Endpoint == edp {
			dockFound = true
			newDocks = append(docks[:i], docks[i+1:]...)
			break
		}
	}
	if !dockFound {
		err = errors.New("Couldn't find dock endpoint " + edp + " in dock resource!")
		return "", err
	}

	if !writeDocksToFile(newDocks) {
		err = errors.New("Deregister dock endpoint " + edp + " failed!")
		return "", err
	} else {
		return "Deregister dock success!", nil
	}
}
*/

func GetDock(dockId string) (*api.DockSpec, error) {
	docks, err := readDocksFromFile()
	if err != nil {
		log.Println("Could not read dock resource:", err)
		return &api.DockSpec{}, err
	}

	for _, dock := range docks {
		if dockId == dock.Id {
			return &dock, nil
		}
	}

	err = errors.New("Could not find this dock service!")
	return &api.DockSpec{}, err
}

func ListDocks() (*[]api.DockSpec, error) {
	docks, err := readDocksFromFile()
	if err != nil {
		log.Println("Could not read dock resource:", err)
		return &[]api.DockSpec{}, err
	}

	return &docks, nil
}
