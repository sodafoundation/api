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

package discovery

import (
	"io/ioutil"

	log "github.com/golang/glog"
)

const (
	defaultConfigFile = "/etc/opensds/dock.json"
)

func loadFromFile(path string) (interface{}, error) {
	if path == "" {
		path = defaultConfigFile
	}

	userJSON, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("ReadFile json failed:", err)
		return nil, err
	}

	return userJSON, nil
}
