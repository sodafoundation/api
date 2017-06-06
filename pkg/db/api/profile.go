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

func CreateProfile(prf *api.StorageProfile) (*api.StorageProfile, error) {
	prfBody, err := json.Marshal(prf)
	if err != nil {
		return &api.StorageProfile{}, err
	}

	dbReq := &db.DbRequest{
		Url:     URL_PREFIX + "profiles/" + prf.Id,
		Content: string(prfBody),
	}
	dbRes := dbCli.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create profile in db:", dbRes.Error)
		return &api.StorageProfile{}, errors.New(dbRes.Error)
	}

	return prf, nil
}

func GetProfile(prfID string) (*api.StorageProfile, error) {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "profiles/" + prfID,
	}
	dbRes := dbCli.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get profile in db:", dbRes.Error)
		return &api.StorageProfile{}, errors.New(dbRes.Error)
	}

	var prf = &api.StorageProfile{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), prf); err != nil {
		log.Println("[Error] When parsing profile in db:", dbRes.Error)
		return &api.StorageProfile{}, errors.New(dbRes.Error)
	}
	return prf, nil
}

func ListProfiles() (*[]api.StorageProfile, error) {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "profiles",
	}
	dbRes := dbCli.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list profiles in db:", dbRes.Error)
		return &[]api.StorageProfile{}, errors.New(dbRes.Error)
	}

	var prfs = []api.StorageProfile{}
	if len(dbRes.Message) == 0 {
		return &prfs, nil
	}
	for _, msg := range dbRes.Message {
		var prf = api.StorageProfile{}
		if err := json.Unmarshal([]byte(msg), &prf); err != nil {
			log.Println("[Error] When parsing profile in db:", dbRes.Error)
			return &[]api.StorageProfile{}, errors.New(dbRes.Error)
		}
		prfs = append(prfs, prf)
	}
	return &prfs, nil
}

func UpdateProfile(prfID, name, desp string, sTag map[string]string) (*api.StorageProfile, error) {
	prf, err := GetProfile(prfID)
	if err != nil {
		return &api.StorageProfile{}, err
	}
	if name != "" {
		prf.Name = name
	}
	if desp != "" {
		prf.Description = desp
	}
	if len(sTag) != 0 {
		prf.StorageTags = sTag
	}
	prfBody, err := json.Marshal(prf)
	if err != nil {
		return &api.StorageProfile{}, err
	}

	dbReq := &db.DbRequest{
		Url:        URL_PREFIX + "profiles/" + prfID,
		NewContent: string(prfBody),
	}
	dbRes := dbCli.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When update profile in db:", dbRes.Error)
		return &api.StorageProfile{}, errors.New(dbRes.Error)
	}
	return prf, nil
}

func DeleteProfile(prfID string) error {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "profiles/" + prfID,
	}
	dbRes := dbCli.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete profile in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}
