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
	"encoding/json"
	"io/ioutil"
	"log"

	api "github.com/opensds/opensds/pkg/model"
)

func readPoolsFromFile() ([]api.StoragePoolSpec, error) {
	var pools []api.StoragePoolSpec

	userJSON, err := ioutil.ReadFile("/etc/opensds/pool.json")
	if err != nil {
		log.Println("ReadFile json failed:", err)
		return pools, err
	}

	// If the pool resource is empty, consider it as a normal condition
	if len(userJSON) == 0 {
		return pools, nil
	}

	// Unmarshal the result
	if err = json.Unmarshal(userJSON, &pools); err != nil {
		log.Println("Unmarshal json failed:", err)
		return pools, err
	}
	return pools, nil
}

func writeDocksToFile(pools []api.StoragePoolSpec) bool {
	body, err := json.Marshal(&pools)
	if err != nil {
		log.Println("Marshal json failed:", err)
		return false
	}

	if err = ioutil.WriteFile("/etc/opensds/pool.json", body, 0666); err != nil {
		log.Println("WriteFile json failed:", err)
		return false
	}
	return true
}
