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
This module implements the operation of dock resources.

*/

package api

import (
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
	dockRoute "github.com/opensds/opensds/pkg/controller/metadata/dock_route"
)

type DockRequestDeliver interface {
	registerDock(dockIp string) (*api.DockRoute, error)

	deregisterDock(dockIp string) (string, error)

	getDock(dockId string) (*api.DockRoute, error)

	listDocks() (*api.DockRoutes, error)
}

// DockRequest is a structure for all properties of
// a dock request
type DockRequest struct{}

func (dr DockRequest) registerDock(dockIp string) (*api.DockRoute, error) {
	return dockRoute.RegisterDockRoute(dockIp)
}

func (dr DockRequest) deregisterDock(dockIp string) (string, error) {
	return dockRoute.DeregisterDockRoute(dockIp)
}

func (dr DockRequest) getDock(dockId string) (*api.DockRoute, error) {
	return dockRoute.GetDockRoute(dockId)
}

func (dr DockRequest) listDocks() (*api.DockRoutes, error) {
	return dockRoute.ListDockRoutes()
}

func RegisterDock(drd DockRequestDeliver, dockIp string) (api.DockRoute, error) {
	dock, err := drd.registerDock(dockIp)
	if err != nil {
		log.Printf("Register dock ip %s error: %v\n", dockIp, err)
		return api.DockRoute{}, err
	}

	return *dock, nil
}

func DeregisterDock(drd DockRequestDeliver, dockIp string) (string, error) {
	res, err := drd.deregisterDock(dockIp)
	if err != nil {
		log.Printf("Deregister dock ip %s error: %v\n", dockIp, err)
		return "", err
	}

	return res, nil
}

func GetDock(drd DockRequestDeliver, dockId string) (api.DockRoute, error) {
	dock, err := drd.getDock(dockId)
	if err != nil {
		log.Printf("Get dock %s error: %v\n", dockId, err)
		return api.DockRoute{}, err
	}

	return *dock, nil
}

func ListDocks(drd DockRequestDeliver) (api.DockRoutes, error) {
	docks, err := drd.listDocks()
	if err != nil {
		log.Printf("List docks error: %v\n", err)
		return api.DockRoutes{}, err
	}

	return *docks, nil
}
