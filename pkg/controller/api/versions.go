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
This module implements the common data structure.

*/

package api

import (
	"encoding/json"
	"errors"
	"log"
)

const (
	KnownVersions = `{
		"versions": [
			{
				"id": "v1",
				"status": "SUPPORTED"
			}
		]
	}`
)

func ListVersions() (AvailableVersions, error) {
	var aVersions AvailableVersions
	err := json.Unmarshal([]byte(KnownVersions), &aVersions)
	if err != nil {
		log.Println(err)
		return AvailableVersions{}, err
	}
	return aVersions, nil
}

func GetVersionv1() (VersionInfo, error) {
	aVersions, err := ListVersions()
	if err != nil {
		log.Println(err)
		return VersionInfo{}, err
	}

	versions := aVersions.Versions
	for _, version := range versions {
		if version.Id == "v1" {
			return version, nil
		}
	}

	err = errors.New("Can't find v1 in available versions!")
	log.Println(err)
	return VersionInfo{}, err
}
