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

package app

import (
	"encoding/json"
	"io/ioutil"
	log "github.com/golang/glog"

	api "github.com/opensds/opensds/pkg/model"
)

func ListDocks() (*[]api.DockSpec, error) {
	docks, err := readDocksFromFile()
	if err != nil {
		log.Error("Could not read dock resource:", err)
		return &[]api.DockSpec{}, err
	}

	return &docks, nil
}

func readDocksFromFile() ([]api.DockSpec, error) {
	var docks []api.DockSpec

	userJSON, err := ioutil.ReadFile("/etc/opensds/dock.json")
	if err != nil {
		log.Error("ReadFile json failed:", err)
		return docks, err
	}

	// If the dock resource is empty, consider it as a normal condition
	if len(userJSON) == 0 {
		return docks, nil
	}

	// Unmarshal the result
	if err = json.Unmarshal(userJSON, &docks); err != nil {
		log.Error("Unmarshal json failed:", err)
		return docks, err
	}
	return docks, nil
}
