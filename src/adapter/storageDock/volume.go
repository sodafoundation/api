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
	"adapter/storagePlugins"
)

type Policy struct {
	//The standard policy configuration will be defined here.
}

type VolumeDriver interface {
	//Any initialization the volume driver does while starting.
	Setup()
	//Any operation the volume driver does while stoping.
	Unset()

	CreateVolume(name string, size int) (string, error)

	GetVolume(volID string) (string, error)

	GetAllVolumes(allowDetails bool) (string, error)

	UpdateVolume(volID string, name string) (string, error)

	DeleteVolume(volID string) (string, error)

	Mount(host string, volID string)

	Unmount(host string, volID string)
}

func CreateVolume(resourceType string, name string, size int) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins := storagePlugins.Init(resourceType)

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.CreateVolume(name, size)
	if err != nil {
		return "Error", err
	} else {
		return result, nil
	}
}

func GetVolume(resourceType string, volID string) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins := storagePlugins.Init(resourceType)

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.GetVolume(volID)
	if err != nil {
		return "Error", err
	} else {
		return result, nil
	}
}

func GetAllVolumes(resourceType string, allowDetails bool) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins := storagePlugins.Init(resourceType)

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.GetAllVolumes(allowDetails)
	if err != nil {
		return "Error", err
	} else {
		return result, nil
	}
}

func UpdateVolume(resourceType string, volID string, name string) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins := storagePlugins.Init(resourceType)

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.UpdateVolume(volID, name)
	if err != nil {
		return "Error", err
	} else {
		return result, nil
	}
}

func DeleteVolume(resourceType string, volID string) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins := storagePlugins.Init(resourceType)

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.DeleteVolume(volID)
	if err != nil {
		return "Error", err
	} else {
		return result, nil
	}
}
