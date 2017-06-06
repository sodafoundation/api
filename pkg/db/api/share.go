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
This module implements the database operation interface of data structure
defined in api module.

*/

package api

import (
	"encoding/json"
	"errors"
	"log"

	api "github.com/opensds/opensds/pkg/api/v1"
	"github.com/opensds/opensds/pkg/db"
)

func CreateShare(shr *api.ShareResponse) (*api.ShareResponse, error) {
	shrBody, err := json.Marshal(shr)
	if err != nil {
		return &api.ShareResponse{}, err
	}

	dbReq := &db.DbRequest{
		Url:     URL_PREFIX + "shares/" + shr.Id,
		Content: string(shrBody),
	}
	dbRes := dbCli.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create share in db:", dbRes.Error)
		return &api.ShareResponse{}, errors.New(dbRes.Error)
	}

	return shr, nil
}

func GetShare(shrID string) (*api.ShareResponse, error) {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "shares/" + shrID,
	}
	dbRes := dbCli.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get share in db:", dbRes.Error)
		return &api.ShareResponse{}, errors.New(dbRes.Error)
	}

	var shr = &api.ShareResponse{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), shr); err != nil {
		log.Println("[Error] When parsing share in db:", dbRes.Error)
		return &api.ShareResponse{}, errors.New(dbRes.Error)
	}
	return shr, nil
}

func ListShares() (*[]api.ShareResponse, error) {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "shares",
	}
	dbRes := dbCli.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list shares in db:", dbRes.Error)
		return &[]api.ShareResponse{}, errors.New(dbRes.Error)
	}

	var shrs = []api.ShareResponse{}
	if len(dbRes.Message) == 0 {
		return &shrs, nil
	}
	for _, msg := range dbRes.Message {
		var shr = api.ShareResponse{}
		if err := json.Unmarshal([]byte(msg), &shr); err != nil {
			log.Println("[Error] When parsing share in db:", dbRes.Error)
			return &[]api.ShareResponse{}, errors.New(dbRes.Error)
		}
		shrs = append(shrs, shr)
	}
	return &shrs, nil
}

func DeleteShare(shrID string) error {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "shares/" + shrID,
	}
	dbRes := dbCli.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete share in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}
