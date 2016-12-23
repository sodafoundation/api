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
This module implements the entry into operations of storageDock module.

*/

package api

import (
	"adapter/storageDock"
)

func CreateVolume(resourceType string, name string, size int) (string, error) {
	result, err := storageDock.CreateVolume(resourceType, name, size)

	if err != nil {
		return "Error", err
	} else {
		return result, nil
	}
}

func GetVolume(resourceType string, volID string) (string, error) {
	result, err := storageDock.GetVolume(resourceType, volID)

	if err != nil {
		return "Error", err
	} else {
		return result, nil
	}
}

func GetAllVolumes(resourceType string, allowDetails bool) (string, error) {
	result, err := storageDock.GetAllVolumes(resourceType, allowDetails)

	if err != nil {
		return "Error", err
	} else {
		return result, nil
	}
}

func UpdateVolume(resourceType string, volID string, name string) (string, error) {
	result, err := storageDock.UpdateVolume(resourceType, volID, name)

	if err != nil {
		return "Error", err
	} else {
		return result, nil
	}
}

func DeleteVolume(resourceType string, volID string) (string, error) {
	result, err := storageDock.DeleteVolume(resourceType, volID)

	if err != nil {
		return "Error", err
	} else {
		return result, nil
	}
}
