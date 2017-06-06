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
	"errors"
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
)

/*
func CreateProfile(name, backend string, tags map[string]string) (*api.StorageProfile, error) {
	profiles, err := readProfilesFromFile()
	if err != nil {
		log.Println("Could not get profiles:", err)
		return &api.StorageProfile{}, err
	}

	profiles = append(profiles, profile)

	if !writeProfilesToFile(profiles) {
		err = errors.New("Create profile " + name + " failed!")
		return &api.StorageProfile{}, err
	} else {
		return &profile, nil
	}
}
*/

func GetProfile(name string) (*api.StorageProfile, error) {
	profiles, err := readProfilesFromFile()
	if err != nil {
		log.Println("Could not read profile table:", err)
		return &api.StorageProfile{}, err
	}

	for _, profile := range profiles {
		if name == profile.Name {
			return &profile, nil
		}
	}

	err = errors.New("Could not find this profile!")
	return &api.StorageProfile{}, err
}

func ListProfiles() (*[]api.StorageProfile, error) {
	profiles, err := readProfilesFromFile()
	if err != nil {
		log.Println("Could not read dock routes:", err)
		return &[]api.StorageProfile{}, err
	}

	return &profiles, nil
}

/*
func DeleteProfile(name string) (string, error) {
	profiles, err := readProfilesFromFile()
	if err != nil {
		log.Println("Could not get profiles:", err)
		return "", err
	}

	var profileFound bool
	var newProfiles []api.StorageProfile

	for i, profile := range profiles {
		if profile.Name == name {
			profileFound = true
			newProfiles = append(profiles[:i], profiles[i+1:]...)
			break
		}
	}
	if !profileFound {
		err = errors.New("Couldn't find profile " + name + " in profile table!")
		return "", err
	}

	if !writeProfilesToFile(newProfiles) {
		err = errors.New("Delete profile " + name + " failed!")
		return "", err
	} else {
		return "Delete profile success!", nil
	}
}
*/
