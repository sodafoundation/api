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
	"github.com/opensds/opensds/pkg/api/rpcapi"
)

type ShareRequestDeliver interface {
	createShare() (string, error)

	getShare() (string, error)

	getAllShares() (string, error)

	updateShare() (string, error)

	deleteShare() (string, error)
}

// ShareRequest is a structure for all properties of
// a share request
type ShareRequest struct {
	ResourceType string `json:"resourceType,omitempty"`
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Size         int    `json:"size"`
	ShareType    string `json:"shareType,omitempty"`
	ShareProto   string `json:"shareProto,omitempty"`
	AllowDetails bool   `json:"allowDetails"`
}

func (sr ShareRequest) createShare() (string, error) {
	return rpcapi.CreateShare(sr.ResourceType, sr.Name, sr.ShareType, sr.ShareProto, sr.Size)
}

func (sr ShareRequest) getShare() (string, error) {
	return rpcapi.GetShare(sr.ResourceType, sr.Id)
}

func (sr ShareRequest) getAllShares() (string, error) {
	return rpcapi.GetAllShares(sr.ResourceType, sr.AllowDetails)
}

func (sr ShareRequest) updateShare() (string, error) {
	return rpcapi.UpdateShare(sr.ResourceType, sr.Id, sr.Name)
}

func (sr ShareRequest) deleteShare() (string, error) {
	return rpcapi.DeleteShare(sr.ResourceType, sr.Id)
}

func CreateShare(srd ShareRequestDeliver) (api.ShareResponse, error) {
	var nullResponse api.ShareResponse

	result, err := srd.createShare()
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

func GetShare(srd ShareRequestDeliver) (api.ShareDetailResponse, error) {
	var nullResponse api.ShareDetailResponse

	result, err := srd.getShare()
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

func ListShares(srd ShareRequestDeliver) ([]api.ShareResponse, error) {
	var nullResponses []api.ShareResponse

	result, err := srd.getAllShares()
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

func UpdateShare(srd ShareRequestDeliver) (api.ShareResponse, error) {
	var nullResponse api.ShareResponse

	result, err := srd.updateShare()
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

func DeleteShare(srd ShareRequestDeliver) (string, error) {
	result, err := srd.deleteShare()

	if err != nil {
		log.Println("Delete file share error: ", err)
		return "", err
	}
	return result, nil
}
