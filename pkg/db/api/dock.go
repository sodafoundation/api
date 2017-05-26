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

func CreateDock(dck *api.Dock) (*api.Dock, error) {
	dckBody, err := json.Marshal(dck)
	if err != nil {
		return &api.Dock{}, err
	}

	dbReq := &db.DbRequest{
		Url:     URL_PREFIX + "docks/" + dck.Id,
		Content: string(dckBody),
	}
	dbRes := dbCli.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create dock in db:", dbRes.Error)
		return &api.Dock{}, errors.New(dbRes.Error)
	}

	return dck, nil
}

func GetDock(dckID string) (*api.Dock, error) {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "docks/" + dckID,
	}
	dbRes := dbCli.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get dock in db:", dbRes.Error)
		return &api.Dock{}, errors.New(dbRes.Error)
	}

	var dck = &api.Dock{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), dck); err != nil {
		log.Println("[Error] When parsing dock in db:", dbRes.Error)
		return &api.Dock{}, errors.New(dbRes.Error)
	}
	return dck, nil
}

func ListDocks() (*[]api.Dock, error) {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "docks",
	}
	dbRes := dbCli.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list docks in db:", dbRes.Error)
		return &[]api.Dock{}, errors.New(dbRes.Error)
	}

	var dcks = []api.Dock{}
	if len(dbRes.Message) == 0 {
		return &dcks, nil
	}
	for _, msg := range dbRes.Message {
		var dck = api.Dock{}
		if err := json.Unmarshal([]byte(msg), &dck); err != nil {
			log.Println("[Error] When parsing dock in db:", dbRes.Error)
			return &[]api.Dock{}, errors.New(dbRes.Error)
		}
		dcks = append(dcks, dck)
	}
	return &dcks, nil
}

func UpdateDock(dckID, status, edp string) (*api.Dock, error) {
	dck, err := GetDock(dckID)
	if err != nil {
		return &api.Dock{}, err
	}
	if status != "" {
		dck.Status = status
	}
	if edp != "" {
		dck.Endpoint = edp
	}
	dckBody, err := json.Marshal(dck)
	if err != nil {
		return &api.Dock{}, err
	}

	dbReq := &db.DbRequest{
		Url:        URL_PREFIX + "docks/" + dckID,
		NewContent: string(dckBody),
	}
	dbRes := dbCli.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When update dock in db:", dbRes.Error)
		return &api.Dock{}, errors.New(dbRes.Error)
	}
	return dck, nil
}

func DeleteDock(dckID string) error {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "docks/" + dckID,
	}
	dbRes := dbCli.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete dock in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}
