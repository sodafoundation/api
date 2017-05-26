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
This module implements a standard SouthBound interface of volume resource to
storage plugins.

*/

package volume

import (
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
	storagePlugins "github.com/opensds/opensds/pkg/dock/plugins"
	"github.com/satori/go.uuid"
)

type VolumePolicy struct {
	//The standard volume policy configuration will be defined here.
}

type VolumeDriver interface {
	//Any initialization the volume driver does while starting.
	Setup()
	//Any operation the volume driver does while stoping.
	Unset()

	CreateVolume(name string, size int32) (*api.VolumeResponse, error)

	GetVolume(volID string) (*api.VolumeResponse, error)

	DeleteVolume(volID string) error

	InitializeConnection(volID string, doLocalAttach, multiPath bool, hostInfo *api.HostInfo) (*api.ConnectionInfo, error)

	AttachVolume(volID, host, mountpoint string) error

	DetachVolume(volID string) error

	CreateSnapshot(name, volID, description string) (*api.VolumeSnapshot, error)

	GetSnapshot(snapID string) (*api.VolumeSnapshot, error)

	DeleteSnapshot(snapID string) error
}

func CreateVolume(resourceType, name string, size int32) (*api.VolumeResponse, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitVP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return &api.VolumeResponse{}, err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.CreateVolume(name, size)
	if err != nil {
		log.Println("Call plugin to create volume failed:", err)
		return &api.VolumeResponse{}, err
	} else {
		return result, nil
	}
}

func GetVolume(resourceType, volID string) (*api.VolumeResponse, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitVP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return &api.VolumeResponse{}, err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.GetVolume(volID)
	if err != nil {
		log.Println("Call plugin to get volume failed:", err)
		return &api.VolumeResponse{}, err
	} else {
		return result, nil
	}
}

func DeleteVolume(resourceType, volID string) error {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitVP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	if err = volumeDriver.DeleteVolume(volID); err != nil {
		log.Println("Call plugin to delete volume failed:", err)
		return err
	}
	return nil
}

func CreateVolumeAttachment(resourceType, volID string, doLocalAttach, multiPath bool, hostInfo *api.HostInfo) (*api.VolumeAttachment, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitVP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return &api.VolumeAttachment{}, err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	connInfo, err := volumeDriver.InitializeConnection(volID, doLocalAttach, multiPath, hostInfo)
	if err != nil {
		log.Println("Call plugin to initialize volume connection failed:", err)
		return &api.VolumeAttachment{}, err
	}

	return &api.VolumeAttachment{
		Id:             uuid.NewV4().String(),
		HostInfo:       *hostInfo,
		ConnectionInfo: *connInfo,
	}, nil
}

func UpdateVolumeAttachment(resourceType, volID, host, mountpoint string) error {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitVP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	if err = volumeDriver.AttachVolume(volID, host, mountpoint); err != nil {
		log.Println("Call plugin to update volume attachment failed:", err)
		return err
	}
	return nil
}

func DeleteVolumeAttachment(resourceType, volID string) error {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitVP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	if err = volumeDriver.DetachVolume(volID); err != nil {
		log.Println("Call plugin to delete volume attachment failed:", err)
		return err
	}
	return nil
}

func CreateSnapshot(resourceType, name, volID, description string) (*api.VolumeSnapshot, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitVP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return &api.VolumeSnapshot{}, err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.CreateSnapshot(name, volID, description)
	if err != nil {
		log.Println("Call plugin to create snapshot failed:", err)
		return &api.VolumeSnapshot{}, err
	} else {
		return result, nil
	}
}

func GetSnapshot(resourceType, snapID string) (*api.VolumeSnapshot, error) {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitVP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return &api.VolumeSnapshot{}, err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	result, err := volumeDriver.GetSnapshot(snapID)
	if err != nil {
		log.Println("Call plugin to get snapshot failed:", err)
		return &api.VolumeSnapshot{}, err
	} else {
		return result, nil
	}
}

func DeleteSnapshot(resourceType, snapID string) error {
	//Get the storage plugins and do some initializations.
	plugins, err := storagePlugins.InitVP(resourceType)
	if err != nil {
		log.Printf("Find %s failed: %v\n", resourceType, err)
		return err
	}

	//Call function of StoragePlugins configured by storage plugins.
	var volumeDriver VolumeDriver = plugins
	if err = volumeDriver.DeleteSnapshot(snapID); err != nil {
		log.Println("Call plugin to delete snapshot failed:", err)
		return err
	}
	return nil
}
