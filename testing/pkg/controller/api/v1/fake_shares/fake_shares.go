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
	"errors"
	"log"

	"github.com/opensds/opensds/testing/pkg/controller/api"
	"github.com/opensds/opensds/testing/pkg/controller/api/grpcapi"
	pb "github.com/opensds/opensds/testing/pkg/grpc/fake_opensds"
)

type ShareRequestDeliver interface {
	createShare() *pb.Response

	getShare() *pb.Response

	listShares() *pb.Response

	deleteShare() *pb.Response
}

// ShareRequest is a structure for all properties of
// a share request
type ShareRequest struct {
	ResourceType string `json:"resourceType,omitempty"`
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Size         int32  `json:"size"`
	ShareType    string `json:"shareType,omitempty"`
	ShareProto   string `json:"shareProto,omitempty"`
	AllowDetails bool   `json:"allowDetails"`
}

func (sr ShareRequest) createShare() *pb.Response {
	return grpcapi.CreateShare(sr.ResourceType, sr.Name, sr.ShareType, sr.ShareProto, sr.Size)
}

func (sr ShareRequest) getShare() *pb.Response {
	return grpcapi.GetShare(sr.ResourceType, sr.Id)
}

func (sr ShareRequest) listShares() *pb.Response {
	return grpcapi.ListShares(sr.ResourceType, sr.AllowDetails)
}

func (sr ShareRequest) deleteShare() *pb.Response {
	return grpcapi.DeleteShare(sr.ResourceType, sr.Id)
}

func CreateShare(srd ShareRequestDeliver) (api.ShareResponse, error) {
	var nullResponse api.ShareResponse

	result := srd.createShare()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("Create file share error:", err)
		return nullResponse, err
	}

	var shareResponse api.ShareResponse
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &shareResponse); err != nil {
		return nullResponse, err
	}
	return shareResponse, nil
}

func GetShare(srd ShareRequestDeliver) (api.ShareDetailResponse, error) {
	var nullResponse api.ShareDetailResponse

	result := srd.getShare()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("Get file share error:", err)
		return nullResponse, err
	}

	var shareDetailResponse api.ShareDetailResponse
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &shareDetailResponse); err != nil {
		return nullResponse, err
	}
	return shareDetailResponse, nil
}

func ListShares(srd ShareRequestDeliver) ([]api.ShareResponse, error) {
	var nullResponses []api.ShareResponse

	result := srd.listShares()
	if result.GetStatus() == "Failure" {
		err := errors.New(result.GetError())
		log.Println("List all file shares error:", err)
		return nullResponses, err
	}

	var sharesResponse []api.ShareResponse
	rbody := []byte(result.GetMessage())
	if err := json.Unmarshal(rbody, &sharesResponse); err != nil {
		return nullResponses, err
	}
	return sharesResponse, nil
}

func DeleteShare(srd ShareRequestDeliver) api.DefaultResponse {
	var defaultResponse api.DefaultResponse

	result := srd.deleteShare()
	if result.GetStatus() == "Failure" {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = result.GetError()
		log.Println("Delete file share error:", defaultResponse.Error)
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}
