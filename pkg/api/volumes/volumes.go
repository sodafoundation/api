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
	"encoding/json"
	"log"

	"github.com/opensds/opensds/pkg/api"
	"github.com/opensds/opensds/pkg/api/rpcapi"
)

func Create(resourceType string, name string, size int) (api.VolumeResponse, error) {
	var nullResponse api.VolumeResponse

	result, err := rpcapi.CreateVolume(resourceType, name, size)
	if err != nil {
		log.Println("Create volume error: ", err)
		return nullResponse, err
	}

	var volumeResponse api.VolumeResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &volumeResponse); err != nil {
		return nullResponse, err
	}
	return volumeResponse, nil
}

func Show(resourceType string, shrID string) (api.VolumeDetailResponse, error) {
	var nullResponse api.VolumeDetailResponse

	result, err := rpcapi.GetVolume(resourceType, shrID)
	if err != nil {
		log.Println("Show volume error: ", err)
		return nullResponse, err
	}

	var volumeDetailResponse api.VolumeDetailResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &volumeDetailResponse); err != nil {
		return nullResponse, err
	}
	return volumeDetailResponse, nil
}

func List(resourceType string, allowDetails bool) ([]api.VolumeResponse, error) {
	var nullResponses []api.VolumeResponse

	result, err := rpcapi.GetAllVolumes(resourceType, allowDetails)
	if err != nil {
		log.Println("List volumes error: ", err)
		return nullResponses, err
	}

	var volumesResponse []api.VolumeResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &volumesResponse); err != nil {
		return nullResponses, err
	}
	return volumesResponse, nil
}

func Update(resourceType string, shrID string, name string) (api.VolumeResponse, error) {
	var nullResponse api.VolumeResponse

	result, err := rpcapi.UpdateVolume(resourceType, shrID, name)
	if err != nil {
		log.Println("Update volume error: ", err)
		return nullResponse, err
	}

	var volumeResponse api.VolumeResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &volumeResponse); err != nil {
		return nullResponse, err
	}
	return volumeResponse, nil
}

func Delete(resourceType string, volID string) (string, error) {
	result, err := rpcapi.DeleteVolume(resourceType, volID)

	if err != nil {
		log.Println("Delete volume error: ", err)
		return "", err
	} else {
		return result, nil
	}
}

func Mount(resourceType, volID, host, mountpoint string) (string, error) {
	result, err := rpcapi.MountVolume(resourceType, volID, host, mountpoint)

	if err != nil {
		log.Println("Mount volume error: ", err)
		return "", err
	} else {
		return result, nil
	}
}

func Unmount(resourceType, volID, attachment string) (string, error) {
	result, err := rpcapi.UnmountVolume(resourceType, volID, attachment)

	if err != nil {
		log.Println("Unmount volume error: ", err)
		return "", err
	} else {
		return result, nil
	}
}
