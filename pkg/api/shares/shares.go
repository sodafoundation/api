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
This module implements the entry into CRUD operation of shares.

*/

package shares

import (
	"encoding/json"
	"log"

	"github.com/opensds/opensds/pkg/api"
	"github.com/opensds/opensds/pkg/api/grpcapi"
)

func Create(resourceType, name, shrType, shrProto string, size int) (api.ShareResponse, error) {
	var nullResponse api.ShareResponse

	result, err := grpcapi.CreateShare(resourceType, name, shrType, shrProto, size)
	if err != nil {
		log.Println("Create file share error: ", err)
		return nullResponse, err
	}

	var shareResponse api.ShareResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &shareResponse); err != nil {
		return nullResponse, err
	}
	return shareResponse, nil
}

func Show(resourceType string, shrID string) (api.ShareDetailResponse, error) {
	var nullResponse api.ShareDetailResponse

	result, err := grpcapi.GetShare(resourceType, shrID)
	if err != nil {
		log.Println("Show file share error: ", err)
		return nullResponse, err
	}

	var shareDetailResponse api.ShareDetailResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &shareDetailResponse); err != nil {
		return nullResponse, err
	}
	return shareDetailResponse, nil
}

func List(resourceType string, allowDetails bool) ([]api.ShareResponse, error) {
	var nullResponses []api.ShareResponse

	result, err := grpcapi.GetAllShares(resourceType, allowDetails)
	if err != nil {
		log.Println("List file shares error: ", err)
		return nullResponses, err
	}

	var sharesResponse []api.ShareResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &sharesResponse); err != nil {
		return nullResponses, err
	}
	return sharesResponse, nil
}

func Update(resourceType string, shrID string, name string) (api.ShareResponse, error) {
	var nullResponse api.ShareResponse

	result, err := grpcapi.UpdateShare(resourceType, shrID, name)
	if err != nil {
		log.Println("Update file share error: ", err)
		return nullResponse, err
	}

	var shareResponse api.ShareResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &shareResponse); err != nil {
		return nullResponse, err
	}
	return shareResponse, nil
}

func Delete(resourceType string, shrID string) (string, error) {
	result, err := grpcapi.DeleteShare(resourceType, shrID)

	if err != nil {
		log.Println("Delete file share error: ", err)
		return "", err
	} else {
		return result, nil
	}
}
