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
This module implements the common data structure.

*/

package v1

type StoragePool struct {
	Id               string            `json:"id"`
	Name             string            `json:"name"`
	Status           string            `json:"status"`
	Description      string            `json:"description"`
	StorageType      string            `json:"storageType"`
	DockName         string            `json:"dockName"`
	AvailabilityZone string            `json:"availabilityZone"`
	TotalCapacity    int64             `json:"totalCapacity"`
	FreeCapacity     int64             `json:"freeCapacity"`
	StorageTag       map[string]string `json:"storageTag"`
}
