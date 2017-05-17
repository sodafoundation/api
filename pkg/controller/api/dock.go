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
	registerDock(edp string, backends []string) (*api.Dock, error)

	deregisterDock(edp string) (string, error)

	getDock(dockId string) (*api.Dock, error)

	listDocks() (*api.Docks, error)
}

// DockRequest is a structure for all properties of
// a dock request
type DockRequest struct{}

func (dr DockRequest) registerDock(edp string, backends []string) (*api.Dock, error) {
	return dockRoute.RegisterDock(edp, backends)
}

func (dr DockRequest) deregisterDock(edp string) (string, error) {
	return dockRoute.DeregisterDock(edp)
}

func (dr DockRequest) getDock(dockId string) (*api.Dock, error) {
	return dockRoute.GetDock(dockId)
}

func (dr DockRequest) listDocks() (*api.Docks, error) {
	return dockRoute.ListDocks()
}

func RegisterDock(drd DockRequestDeliver, edp string, backends []string) (api.Dock, error) {
	dock, err := drd.registerDock(edp, backends)
	if err != nil {
		log.Printf("Register dock endpoint %s error: %v\n", edp, err)
		return api.Dock{}, err
	}

	return *dock, nil
}

func DeregisterDock(drd DockRequestDeliver, edp string) (string, error) {
	res, err := drd.deregisterDock(edp)
	if err != nil {
		log.Printf("Deregister dock endpoint %s error: %v\n", edp, err)
		return "", err
	}

	return res, nil
}

func GetDock(drd DockRequestDeliver, dockId string) (api.Dock, error) {
	dock, err := drd.getDock(dockId)
	if err != nil {
		log.Printf("Get dock %s error: %v\n", dockId, err)
		return api.Dock{}, err
	}

	return *dock, nil
}

func ListDocks(drd DockRequestDeliver) (api.Docks, error) {
	docks, err := drd.listDocks()
	if err != nil {
		log.Printf("List docks error: %v\n", err)
		return api.Docks{}, err
	}

	return *docks, nil
}
