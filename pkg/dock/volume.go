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

package dock

import (
	"log"

	storagePlugins "github.com/opensds/opensds/pkg/dock/plugins"
)

type VolumePolicy struct {
	//The standard volume policy configuration will be defined here.
}

type VolumeDriver interface {
	//Any initialization the volume driver does while starting.
	Setup()
	//Any operation the volume driver does while stoping.
	Unset()

	CreateVolume(name string, volType string, size int32) (string, error)

	GetVolume(volID string) (string, error)

	GetAllVolumes(allowDetails bool) (string, error)

	DeleteVolume(volID string) (string, error)

	AttachVolume(volID, host, device string) (string, error)

	DetachVolume(volID string, attachment string) (string, error)
}

func CreateVolume(resourceType, name, volType string, size int32) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitVP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return "", err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.CreateVolume(name, volType, size)
	if err != nil {
		log.Println("Call plugin to create volume failed:", err)
		return "", err
	} else {
		return result, nil
	}
}

func GetVolume(resourceType, volID string) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitVP(resourceType)
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
	plugins, err := storagePlugins.InitVP(resourceType)
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

func DeleteVolume(resourceType, volID string) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitVP(resourceType)
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

func AttachVolume(resourceType, volID, host, device string) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitVP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return "", err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.AttachVolume(volID, host, device)
	if err != nil {
		log.Println("Call plugin to attach volume failed:", err)
		return "", err
	} else {
		return result, nil
	}
}

func DetachVolume(resourceType, volID, attachment string) (string, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitVP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return "", err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.DetachVolume(volID, attachment)
	if err != nil {
		log.Println("Call plugin to detach volume failed:", err)
		return "", err
	} else {
		return result, nil
	}
}

func MountVolume(mountDir, device, fsType string) (string, error) {
	if err := storagePlugins.Mount(mountDir, device, fsType); err != nil {
		log.Println("Call plugin to mount volume failed:", err)
		return "Mount volume failed!", err
	}
	return "Mount volume success!", nil
}

func UnmountVolume(mountDir string) (string, error) {
	if err := storagePlugins.Unmount(mountDir); err != nil {
		log.Println("Call plugin to unmount volume failed:", err)
		return "Unmount volume failed!", err
	}
	return "Unmount volume success!", nil
}
