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

type Dock struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Endpoint         string `json:"endpoint"`
	Status           string `json:"status"`
	AvailabilityZone string `json:"availabilityZone"`
	DriverName       string `json:"driverName"`
	StorageName      string `json:"storageName"`
	StorageType      string `json:"storageType"`
	StorageFqdn      string `json:"storageFqdn"`
}
