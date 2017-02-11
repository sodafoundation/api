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
This module implements the api module client of grpc service.

*/

package grpcapi

import (
	"strconv"
	"strings"

	"github.com/opensds/opensds/pkg/grpc"
)

func CreateVolume(resourceType string, name string, size int) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"CreateVolume", resourceType, name, strconv.Itoa(size)}
	return client.Run(url, strings.Join(action[:], ","))
}

func GetVolume(resourceType string, volID string) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"GetVolume", resourceType, volID}
	return client.Run(url, strings.Join(action[:], ","))
}

func GetAllVolumes(resourceType string, allowDetails bool) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"GetAllVolumes", resourceType,
		strconv.FormatBool(allowDetails)}
	return client.Run(url, strings.Join(action[:], ","))
}

func UpdateVolume(resourceType string, volID string, name string) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"UpdateVolume", resourceType, volID, name}
	return client.Run(url, strings.Join(action[:], ","))
}

func DeleteVolume(resourceType string, volID string) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"DeleteVolume", resourceType, volID}
	return client.Run(url, strings.Join(action[:], ","))
}

func MountVolume(resourceType, volID, host, mountpoint string) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"MountVolume", resourceType, volID, host, mountpoint}
	return client.Run(url, strings.Join(action[:], ","))
}

func UnmountVolume(resourceType, volID, attchment string) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"UnmountVolume", resourceType, volID, attchment}
	return client.Run(url, strings.Join(action[:], ","))
}

func CreateShare(resourceType, name, shrType, shrProto string, size int) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"CreateShare", resourceType, name, shrType, shrProto, strconv.Itoa(size)}
	return client.Run(url, strings.Join(action[:], ","))
}

func GetShare(resourceType string, shrID string) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"GetShare", resourceType, shrID}
	return client.Run(url, strings.Join(action[:], ","))
}

func GetAllShares(resourceType string, allowDetails bool) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"GetAllShares", resourceType,
		strconv.FormatBool(allowDetails)}
	return client.Run(url, strings.Join(action[:], ","))
}

func UpdateShare(resourceType string, shrID string, name string) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"UpdateShare", resourceType, shrID, name}
	return client.Run(url, strings.Join(action[:], ","))
}

func DeleteShare(resourceType string, shrID string) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"DeleteShare", resourceType, shrID}
	return client.Run(url, strings.Join(action[:], ","))
}

func CreateDatabase(name string, size int) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"CreateDatabase", name, strconv.Itoa(size)}
	return client.Run(url, strings.Join(action[:], ","))
}

func GetDatabase(id int, name string) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"GetDatabase", strconv.Itoa(id), name}
	return client.Run(url, strings.Join(action[:], ","))
}

func GetAllDatabases() (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"GetAllDatabases"}
	return client.Run(url, strings.Join(action[:], ","))
}

func UpdateDatabase(id int, size int, name string) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"UpdateDatabase", strconv.Itoa(id),
		strconv.Itoa(size), name}
	return client.Run(url, strings.Join(action[:], ","))
}

func DeleteDatabase(id int, name string, cascade bool) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"DeleteDatabase", strconv.Itoa(id),
		name, strconv.FormatBool(cascade)}
	return client.Run(url, strings.Join(action[:], ","))
}

func CreateFileSystem(name string, size int) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"CreateFileSystem", name, strconv.Itoa(size)}
	return client.Run(url, strings.Join(action[:], ","))
}

func GetFileSystem(id int, name string) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"GetFileSystem", strconv.Itoa(id), name}
	return client.Run(url, strings.Join(action[:], ","))
}

func GetAllFileSystems() (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"GetAllFileSystems"}
	return client.Run(url, strings.Join(action[:], ","))
}

func UpdateFileSystem(id int, size int, name string) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"UpdateFileSystem", strconv.Itoa(id),
		strconv.Itoa(size), name}
	return client.Run(url, strings.Join(action[:], ","))
}

func DeleteFileSystem(id int, name string, cascade bool) (string, error) {
	var client grpc.Client
	url := "opensds/orchestration"
	action := []string{"DeleteFileSystem", strconv.Itoa(id),
		name, strconv.FormatBool(cascade)}
	return client.Run(url, strings.Join(action[:], ","))
}
