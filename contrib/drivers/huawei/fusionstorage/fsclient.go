// Copyright (c) 2019 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

type RestCommon struct {
	username string
	password string
	version  string
	addess   string
	headers  map[string]string
}

func newRestCommon(username, password, url string) *RestCommon {

	return &RestCommon{
		addess:   url,
		username: username,
		password: password,
		headers:  map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}
}

func (r *RestCommon) getVersion() error {
	url := "rest/version"
	r.headers["Referer"] = r.addess + BasicURI
	content, err := r.request(url, "GET", true, nil)
	if err != nil {
		return fmt.Errorf("Failed to get version, %v", err)
	}

	var v version
	err = json.Unmarshal(content, &v)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal the result, %v", err)
	}

	r.version = v.CurrentVersion

	return nil
}

func (r *RestCommon) login() error {
	r.getVersion()
	url := "/sec/login"
	data := map[string]string{"userName": r.username, "password": r.password}
	_, err := r.request(url, "POST", false, data)
	if err != nil {
		return err
	}

	return nil
}

func (r *RestCommon) queryPoolInfo() (*poolResp, error) {
	url := "/storagePool"
	result, err := r.request(url, "GET", false, nil)
	if err != nil {
		return nil, err
	}

	var pools *poolResp
	if err := json.Unmarshal(result, &pools); err != nil {
		return nil, err
	}
	return pools, nil
}

func (r *RestCommon) createVolume(volName, poolId string, volSize int64) error {
	url := "/volume/create"
	polID, _ := strconv.Atoi(poolId)
	params := map[string]interface{}{"volName": volName, "volSize": volSize, "poolId": polID}
	if _, err := r.request(url, "POST", false, params); err != nil {
		return err
	}
	return nil
}

func (r *RestCommon) deleteVolume(volName string) error {
	url := "/volume/delete"
	params := map[string]interface{}{"volNames": []string{volName}}
	_, err := r.request(url, "POST", false, params)
	if err != nil {
		return err
	}

	return nil
}

func (r *RestCommon) request(url, method string, isGetVersion bool, reqParams interface{}) ([]byte, error) {
	var callUrl string
	if !isGetVersion {
		callUrl = r.addess + BasicURI + r.version + url
	} else {
		callUrl = r.addess + BasicURI + url
	}

	fmt.Println(callUrl)
	// No verify by SSL
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// initialize http client
	client := &http.Client{Transport: tr}

	var body []byte
	var err error
	if reqParams != nil {
		body, err = json.Marshal(reqParams)
		if err != nil {
			return nil, fmt.Errorf("Failed to marshal the request parameters, url is %s, error is %v", callUrl, err)
		}
	}

	req, err := http.NewRequest(strings.ToUpper(method), callUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("Failed to initiate the request, url is %s, error is %v", callUrl, err)
	}

	// initiate the header
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	// do the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Process request failed: %v, url is %s", err, callUrl)
	}
	defer resp.Body.Close()

	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read from response body failed: %v, url is %s", err, callUrl)
	}

	if 400 <= resp.StatusCode && resp.StatusCode <= 599 {
		pc, _, line, _ := runtime.Caller(2)
		return nil, fmt.Errorf("return status code is: %s, return content is: %s, error function is: %s, error line is: %s, url is %s",
			strconv.Itoa(resp.StatusCode), string(respContent), runtime.FuncForPC(pc).Name(), strconv.Itoa(line), callUrl)
	}

	// Check the error code in the returned content
	var respResult *responseResult
	if err := json.Unmarshal(respContent, &respResult); err != nil {
		return nil, err
	}

	if respResult.RespCode != 0 {
		return nil, fmt.Errorf("Request failed, url is %s, %s", callUrl, string(respContent))
	}

	if resp.Header != nil && len(resp.Header["X-Auth-Token"]) > 0 {
		token := resp.Header["X-Auth-Token"][0]
		r.headers["x-auth-token"] = token
	}

	return respContent, nil
}
