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

func CreatePool(pol *api.StoragePool) (*api.StoragePool, error) {
	polBody, err := json.Marshal(pol)
	if err != nil {
		return &api.StoragePool{}, err
	}

	dbReq := &db.DbRequest{
		Url:     URL_PREFIX + "pools/" + pol.Id,
		Content: string(polBody),
	}
	dbRes := dbCli.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When create pol in db:", dbRes.Error)
		return &api.StoragePool{}, errors.New(dbRes.Error)
	}

	return pol, nil
}

func GetPool(polID string) (*api.StoragePool, error) {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "pools/" + polID,
	}
	dbRes := dbCli.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When get pool in db:", dbRes.Error)
		return &api.StoragePool{}, errors.New(dbRes.Error)
	}

	var pol = &api.StoragePool{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), pol); err != nil {
		log.Println("[Error] When parsing pool in db:", dbRes.Error)
		return &api.StoragePool{}, errors.New(dbRes.Error)
	}
	return pol, nil
}

func ListPools() (*[]api.StoragePool, error) {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "pools",
	}
	dbRes := dbCli.List(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When list pools in db:", dbRes.Error)
		return &[]api.StoragePool{}, errors.New(dbRes.Error)
	}

	var pols = []api.StoragePool{}
	if len(dbRes.Message) == 0 {
		return &pols, nil
	}
	for _, msg := range dbRes.Message {
		var pol = api.StoragePool{}
		if err := json.Unmarshal([]byte(msg), &pol); err != nil {
			log.Println("[Error] When parsing pool in db:", dbRes.Error)
			return &[]api.StoragePool{}, errors.New(dbRes.Error)
		}
		pols = append(pols, pol)
	}
	return &pols, nil
}

func UpdatePool(polID, name, status, desp string, usedCapacity int64, used bool) (*api.StoragePool, error) {
	pol, err := GetPool(polID)
	if err != nil {
		return &api.StoragePool{}, err
	}
	if name != "" {
		pol.Name = name
	}
	if status != "" {
		pol.Status = status
	}
	if desp != "" {
		pol.Description = desp
	}
	if usedCapacity != 0 {
		if used {
			pol.FreeCapacity = pol.FreeCapacity - usedCapacity
		} else {
			pol.FreeCapacity = pol.FreeCapacity + usedCapacity
		}
	}
	polBody, err := json.Marshal(pol)
	if err != nil {
		return &api.StoragePool{}, err
	}

	dbReq := &db.DbRequest{
		Url:        URL_PREFIX + "pools/" + polID,
		NewContent: string(polBody),
	}
	dbRes := dbCli.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When update pool in db:", dbRes.Error)
		return &api.StoragePool{}, errors.New(dbRes.Error)
	}
	return pol, nil
}

func DeletePool(polID string) error {
	dbReq := &db.DbRequest{
		Url: URL_PREFIX + "pools/" + polID,
	}
	dbRes := dbCli.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Println("[Error] When delete pool in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}
