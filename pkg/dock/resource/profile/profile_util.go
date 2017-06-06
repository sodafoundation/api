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
This module defines the configuration work of profile resource, please DO NOT
modify this file.

*/

package profile

import (
	"encoding/json"
	"io/ioutil"
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
)

func readProfilesFromFile() ([]api.StorageProfile, error) {
	var profiles []api.StorageProfile

	userJSON, err := ioutil.ReadFile("/etc/opensds/profile.json")
	if err != nil {
		log.Println("ReadFile json failed:", err)
		return profiles, err
	}

	// If the profile table is empty, consider it as a normal condition
	if len(userJSON) == 0 {
		return profiles, nil
	} else if err = json.Unmarshal(userJSON, &profiles); err != nil {
		log.Println("Unmarshal json failed:", err)
		return profiles, err
	}
	return profiles, nil
}

func writeProfilesToFile(profiles []api.StorageProfile) bool {
	body, err := json.Marshal(&profiles)
	if err != nil {
		log.Println("Marshal json failed:", err)
		return false
	}

	if err = ioutil.WriteFile("/etc/opensds/profile.json", body, 0666); err != nil {
		log.Println("WriteFile json failed:", err)
		return false
	}
	return true
}
