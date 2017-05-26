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
This module implements the entry into CRUD operation of volumes.

*/

package api

import (
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
	db "github.com/opensds/opensds/pkg/db/api"
)

type PoolRequestDeliver interface {
	getPool() (*api.StoragePool, error)

	listPools() (*[]api.StoragePool, error)
}

// PoolRequest is a structure for all properties of
// a storage pool request
type PoolRequest struct {
	Id string `json:"id"`
}

func (pr *PoolRequest) getPool() (*api.StoragePool, error) {
	return db.GetPool(pr.Id)
}

func (pr *PoolRequest) listPools() (*[]api.StoragePool, error) {
	return db.ListPools()
}

func GetPool(prd PoolRequestDeliver) (api.StoragePool, error) {
	var nullResponse = api.StoragePool{}

	pol, err := prd.getPool()
	if err != nil {
		log.Println("Get pool error:", err)
		return nullResponse, err
	}
	return *pol, nil
}

func ListPools(prd PoolRequestDeliver) ([]api.StoragePool, error) {
	var nullResponses = []api.StoragePool{}

	pols, err := prd.listPools()
	if err != nil {
		log.Println("List all pools error:", err)
		return nullResponses, err
	}
	return *pols, nil
}
