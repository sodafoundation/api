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
This module implements a standard SouthBound interface to storage plugins.

*/

package storageDock

import (
	"log"

	"adapter/storagePlugins"
)

type SharePolicy struct {
	//The standard share policy configuration will be defined here.
}

type ShareDriver interface {
	//Any initialization the volume driver does while starting.
	Setup()
	//Any operation the volume driver does while stoping.
	Unset()

	CreateShare(name string, size int) (string, error)

	GetShare(shrID string) (string, error)

	GetAllShares(allowDetails bool) (string, error)

	UpdateShare(shrID string, name string) (string, error)

	DeleteShare(shrID string) (string, error)
}

func CreateShare(resourceType string, name string, size int) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitSP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return "", err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var shareDriver ShareDriver = plugins
	result, err := shareDriver.CreateShare(name, size)
	if err != nil {
		log.Println("Call plugin to create volume failed:", err)
		return "", err
	} else {
		return result, nil
	}
}

func GetShare(resourceType, shrID string) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitSP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return "", err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var shareDriver ShareDriver = plugins
	result, err := shareDriver.GetShare(shrID)
	if err != nil {
		log.Println("Call plugin to get volume failed:", err)
		return "", err
	} else {
		return result, nil
	}
}

func GetAllShares(resourceType string, allowDetails bool) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitSP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return "", err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var shareDriver ShareDriver = plugins
	result, err := shareDriver.GetAllShares(allowDetails)
	if err != nil {
		log.Println("Call plugin to get all volumes failed:", err)
		return "", err
	} else {
		return result, nil
	}
}

func UpdateShare(resourceType, shrID, name string) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitSP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return "", err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var shareDriver ShareDriver = plugins
	result, err := shareDriver.UpdateShare(shrID, name)
	if err != nil {
		log.Println("Call plugin to update volume failed:", err)
		return "", err
	} else {
		return result, nil
	}
}

func DeleteShare(resourceType, shrID string) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitSP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return "", err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var shareDriver ShareDriver = plugins
	result, err := shareDriver.DeleteShare(shrID)
	if err != nil {
		log.Println("Call plugin to delete share failed:", err)
		return "", err
	} else {
		return result, nil
	}
}
