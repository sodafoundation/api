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
	"errors"
	"log"
	"math/rand"
	"time"

	api "github.com/opensds/opensds/pkg/api/v1"
)

func RegisterDockRoute(dockIp string) (*api.DockRoute, error) {
	route := generateDockRoute(dockIp)

	routes, err := readDockRoutesFromFile()
	if err != nil {
		log.Println("Could not get dock routes:", err)
		return &api.DockRoute{}, err
	}

	routes = append(routes, route)

	if !writeDockRoutesToFile(routes) {
		err = errors.New("Register dock ip " + dockIp + " failed!")
		return &api.DockRoute{}, err
	} else {
		return &route, nil
	}
}

func DeregisterDockRoute(dockIp string) (string, error) {
	routes, err := readDockRoutesFromFile()
	if err != nil {
		log.Println("Could not get dock routes:", err)
		return "", err
	}

	var dockFound bool
	var newRoutes []api.DockRoute

	for i, route := range routes {
		if route.Address == dockIp {
			dockFound = true
			newRoutes = append(routes[:i], routes[i+1:]...)
			break
		}
	}
	if !dockFound {
		err = errors.New("Couldn't find dock ip " + dockIp + " in dock route tables!")
		return "", err
	}

	if !writeDockRoutesToFile(newRoutes) {
		err = errors.New("Deregister dock ip " + dockIp + " failed!")
		return "", err
	} else {
		return "Deregister dock success!", nil
	}
}

func GetDockRoute(dockId string) (*api.DockRoute, error) {
	routes, err := readDockRoutesFromFile()
	if err != nil {
		log.Println("Could not read dock routes:", err)
		return &api.DockRoute{}, err
	}

	if dockId == "" {
		log.Println("Dock id not provided, arrange the address randomly!")
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		n := r.Intn(len(routes))
		return &routes[n], nil
	}

	for _, route := range routes {
		if dockId == route.Id {
			return &route, nil
		}
	}

	err = errors.New("Could not find this dock service!")
	return &api.DockRoute{}, err
}

func ListDockRoutes() (*api.DockRoutes, error) {
	routes, err := readDockRoutesFromFile()
	if err != nil {
		log.Println("Could not read dock routes:", err)
		return &api.DockRoutes{}, err
	}

	return &api.DockRoutes{
		Routes: routes,
	}, nil
}
