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

	MountVolume(volID, host, mountpoint string) (string, error)

	UnmountVolume(volID string, attachment string) (string, error)
}

func CreateVolume(resourceType string, name string, size int) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.Init(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return "", err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.CreateVolume(name, size)
	if err != nil {
		log.Println("Call plugin to create volume failed:", err)
		return "", err
	} else {
		return result, nil
	}
}

func GetVolume(resourceType, volID string) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.Init(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return "", err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.GetVolume(volID)
	if err != nil {
		log.Println("Call plugin to get volume failed:", err)
		return "", err
	} else {
		return result, nil
	}
}

func GetAllVolumes(resourceType string, allowDetails bool) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.Init(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return "", err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.GetAllVolumes(allowDetails)
	if err != nil {
		log.Println("Call plugin to get all volumes failed:", err)
		return "", err
	} else {
		return result, nil
	}
}

func UpdateVolume(resourceType, volID, name string) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.Init(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return "", err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.UpdateVolume(volID, name)
	if err != nil {
		log.Println("Call plugin to update volume failed:", err)
		return "", err
	} else {
		return result, nil
	}
}

func DeleteVolume(resourceType, volID string) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.Init(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return "", err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.DeleteVolume(volID)
	if err != nil {
		log.Println("Call plugin to delete volume failed:", err)
		return "", err
	} else {
		return result, nil
	}
}

func MountVolume(resourceType, volID, host, mountpoint string) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.Init(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return "", err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.MountVolume(volID, host, mountpoint)
	if err != nil {
		log.Println("Call plugin to mount volume failed:", err)
		return "", err
	} else {
		return result, nil
	}
}

func UnmountVolume(resourceType, volID, attachment string) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.Init(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return "", err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.UnmountVolume(volID, attachment)
	if err != nil {
		log.Println("Call plugin to unmount volume failed:", err)
		return "", err
	} else {
		return result, nil
	}
}
