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
This module implements the entry into CRUD operation of volumes.

*/

package volumes

import (
	"log"

	"api/grpcapi"
)

func Create(resourceType string, name string, size int) (string, error) {
	result, err := grpcapi.CreateVolume(resourceType, name, size)

	if err != nil {
		log.Println("Create volume error: ", err)
		return "", err
	} else {
		return result, nil
	}
}

func Show(resourceType string, volID string) (string, error) {
	result, err := grpcapi.GetVolume(resourceType, volID)

	if err != nil {
		log.Println("Show volume error: ", err)
		return "", err
	} else {
		return result, nil
	}
}

func List(resourceType string, allowDetails bool) (string, error) {
	result, err := grpcapi.GetAllVolumes(resourceType, allowDetails)

	if err != nil {
		log.Println("List volume error: ", err)
		return "", err
	} else {
		return result, nil
	}
}

func Update(resourceType string, volID string, name string) (string, error) {
	result, err := grpcapi.UpdateVolume(resourceType, volID, name)

	if err != nil {
		log.Println("Update volume error: ", err)
		return "", err
	} else {
		return result, nil
	}
}

func Delete(resourceType string, volID string) (string, error) {
	result, err := grpcapi.DeleteVolume(resourceType, volID)

	if err != nil {
		log.Println("Delete volume error: ", err)
		return "", err
	} else {
		return result, nil
	}
}

func Mount(resourceType, volID, host, mountpoint string) (string, error) {
	result, err := grpcapi.MountVolume(resourceType, volID, host, mountpoint)

	if err != nil {
		log.Println("Mount volume error: ", err)
		return "", err
	} else {
		return result, nil
	}
}

func Unmount(resourceType, volID, attachment string) (string, error) {
	result, err := grpcapi.UnmountVolume(resourceType, volID, attachment)

	if err != nil {
		log.Println("Unmount volume error: ", err)
		return "", err
	} else {
		return result, nil
	}
}
