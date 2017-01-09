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
This module implements the orchestation module client of grpc service.

*/

package grpcapi

import (
	"strconv"
	"strings"

	"grpc"
)

func CreateVolume(resourceType string, name string, size int) (string, error) {
	var client grpc.Client
	url := "opensds/adapter"
	action := []string{"CreateVolume", resourceType, name, strconv.Itoa(size)}
	return client.Run(url, strings.Join(action[:], ","))
}

func GetVolume(resourceType string, volID string) (string, error) {
	var client grpc.Client
	url := "opensds/adapter"
	action := []string{"GetVolume", resourceType, volID}
	return client.Run(url, strings.Join(action[:], ","))
}

func GetAllVolumes(resourceType string, allowDetails bool) (string, error) {
	var client grpc.Client
	url := "opensds/adapter"
	action := []string{"GetAllVolumes", resourceType,
		strconv.FormatBool(allowDetails)}
	return client.Run(url, strings.Join(action[:], ","))
}

func UpdateVolume(resourceType string, volID string, name string) (string, error) {
	var client grpc.Client
	url := "opensds/adapter"
	action := []string{"UpdateVolume", resourceType, volID, name}
	return client.Run(url, strings.Join(action[:], ","))
}

func DeleteVolume(resourceType string, volID string) (string, error) {
	var client grpc.Client
	url := "opensds/adapter"
	action := []string{"DeleteVolume", resourceType, volID}
	return client.Run(url, strings.Join(action[:], ","))
}

func MountVolume(resourceType, volID, host, mountpoint string) (string, error) {
	var client grpc.Client
	url := "opensds/adapter"
	action := []string{"MountVolume", resourceType, volID, host, mountpoint}
	return client.Run(url, strings.Join(action[:], ","))
}

func UnmountVolume(resourceType, volID, attchment string) (string, error) {
	var client grpc.Client
	url := "opensds/adapter"
	action := []string{"UnmountVolume", resourceType, volID, attchment}
	return client.Run(url, strings.Join(action[:], ","))
}
