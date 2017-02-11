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
This module implements the enry into the operations of orchestration module.

Request about volume operation will be passed to the grpc client and requests
about other resources (database, fileSystem, etc) will be passed to metaData
service module.

*/

package api

import (
	"log"

	"github.com/opensds/opensds/pkg/orchestration/grpcapi"
)

func CreateShare(resourceType, name, shrType, shrProto string, size int) (string, error) {
	result, err := grpcapi.CreateShare(resourceType, name, shrType, shrProto, size)

	if err != nil {
		log.Println("Error occured in orchestration module when create share!")
		return "", err
	} else {
		return result, nil
	}
}

func GetShare(resourceType string, shrID string) (string, error) {
	result, err := grpcapi.GetShare(resourceType, shrID)

	if err != nil {
		log.Println("Error occured in orchestration module when get share!")
		return "", err
	} else {
		return result, nil
	}
}

func GetAllShares(resourceType string, allowDetails bool) (string, error) {
	result, err := grpcapi.GetAllShares(resourceType, allowDetails)

	if err != nil {
		log.Println("Error occured in orchestration module when get all shares!")
		return "", err
	} else {
		return result, nil
	}
}

func UpdateShare(resourceType string, shrID string, name string) (string, error) {
	result, err := grpcapi.UpdateShare(resourceType, shrID, name)

	if err != nil {
		log.Println("Error occured in orchestration module when update share!")
		return "", err
	} else {
		return result, nil
	}
}

func DeleteShare(resourceType string, shrID string) (string, error) {
	result, err := grpcapi.DeleteShare(resourceType, shrID)

	if err != nil {
		log.Println("Error occured in orchestration module when delete share!")
		return "", err
	} else {
		return result, nil
	}
}
