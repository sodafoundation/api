// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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
This module implements the entry into operations of open policy agent module.

*/

package opa

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
)

const URLPrefix = "http://localhost:8181/v1/data/"

var httpClient = &http.Client{}

type Data interface{}

func RegisterData(input Data) error {
	var url = URLPrefix

	switch input.(type) {
	case *[]*model.ProfileSpec:
		url += "profiles"
	case *[]*model.StoragePoolSpec:
		url += "pools"
	case *[]*model.DockSpec:
		url += "docks"
	}

	log.Info("Start PUT request to register policy data, url =", url)

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewReader(inputJSON))
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return err
	}
	return nil
}

func PatchData(input Data, op, path string) error {
	var url = URLPrefix

	switch input.(type) {
	case *model.ProfileSpec:
		url += "profiles"
	case *model.StoragePoolSpec:
		url += "pools"
	case *model.DockSpec:
		url += "docks"
	}

	log.Info("Start PATCH request to update policy data, url =", url)

	var data = []struct {
		Op    string      `json:"op"`
		Path  string      `json:"path"`
		Value interface{} `json:"value"`
	}{
		{
			Op:    op,
			Path:  path,
			Value: input,
		},
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(dataJSON))
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return err
	}
	return nil
}

func GetPoolData(input string) (*model.StoragePoolSpec, error) {
	var url = URLPrefix + "opa/policies/find_supported_pools/" + input

	log.Info("Start GET request to get pool data, url =", url)

	req, err := http.NewRequest("GET", url, nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nil, err
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out = struct {
		Result *model.StoragePoolSpec `json:"result"`
	}{
		Result: &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{},
		},
	}
	if err = json.Unmarshal(rbody, &out); err != nil {
		return nil, err
	}

	return out.Result, nil
}

func GetDockData(input string) (*model.DockSpec, error) {
	var url = URLPrefix + "opa/policies/find_dock_by_pool_id/" + input

	log.Info("Start GET request to get dock data, url =", url)

	req, err := http.NewRequest("GET", url, nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nil, err
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out = struct {
		Result *model.DockSpec `json:"result"`
	}{
		Result: &model.DockSpec{
			BaseModel: &model.BaseModel{},
		},
	}
	if err = json.Unmarshal(rbody, &out); err != nil {
		return nil, err
	}

	return out.Result, nil
}

// CheckHTTPResponseStatusCode compares http response header StatusCode against expected
// statuses. Primary function is to ensure StatusCode is in the 20x (return nil).
// Ok: 200. Created: 201. Accepted: 202. No Content: 204. Partial Content: 206.
// Otherwise return error message.
func CheckHTTPResponseStatusCode(resp *http.Response) error {
	switch resp.StatusCode {
	case 200, 201, 202, 204, 206:
		return nil
	case 400:
		return errors.New("Error: response == 400 bad request")
	case 401:
		return errors.New("Error: response == 401 unauthorised")
	case 403:
		return errors.New("Error: response == 403 forbidden")
	case 404:
		return errors.New("Error: response == 404 not found")
	case 405:
		return errors.New("Error: response == 405 method not allowed")
	case 409:
		return errors.New("Error: response == 409 conflict")
	case 413:
		return errors.New("Error: response == 413 over limit")
	case 415:
		return errors.New("Error: response == 415 bad media type")
	case 422:
		return errors.New("Error: response == 422 unprocessable")
	case 429:
		return errors.New("Error: response == 429 too many request")
	case 500:
		return errors.New("Error: response == 500 instance fault / server err")
	case 501:
		return errors.New("Error: response == 501 not implemented")
	case 503:
		return errors.New("Error: response == 503 service unavailable")
	}
	return errors.New("Error: unexpected response status code")
}
