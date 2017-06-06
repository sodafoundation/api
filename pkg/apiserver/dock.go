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
	db "github.com/opensds/opensds/pkg/db/api"
)

type DockRequestDeliver interface {
	getDock() (*api.Dock, error)

	listDocks() (*[]api.Dock, error)
}

// DockRequest is a structure for all properties of
// a storage dock request
type DockRequest struct {
	Id string `json:"id"`
}

func (dr *DockRequest) getDock() (*api.Dock, error) {
	return db.GetDock(dr.Id)
}

func (dr *DockRequest) listDocks() (*[]api.Dock, error) {
	return db.ListDocks()
}

func GetDock(drd DockRequestDeliver) (api.Dock, error) {
	var nullResponse = api.Dock{}

	dck, err := drd.getDock()
	if err != nil {
		log.Println("Get dock error:", err)
		return nullResponse, err
	}
	return *dck, nil
}

func ListDocks(drd DockRequestDeliver) ([]api.Dock, error) {
	var nullResponses = []api.Dock{}

	dcks, err := drd.listDocks()
	if err != nil {
		log.Println("List all docks error:", err)
		return nullResponses, err
	}
	return *dcks, nil
}
