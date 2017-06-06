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

package pool

import (
	"errors"
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
)

func GetPool(poolId string) (*api.StoragePool, error) {
	pools, err := readPoolsFromFile()
	if err != nil {
		log.Println("Could not read pool resource:", err)
		return &api.StoragePool{}, err
	}

	for _, pool := range pools {
		if poolId == pool.Id {
			return &pool, nil
		}
	}

	err = errors.New("Could not find any pool resource!")
	return &api.StoragePool{}, err
}

func ListPools() (*[]api.StoragePool, error) {
	pools, err := readPoolsFromFile()
	if err != nil {
		log.Println("Could not read pool resource:", err)
		return &[]api.StoragePool{}, err
	}

	return &pools, nil
}
