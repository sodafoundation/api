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
This module implements osbrick plugin for OpenSDS. It will pass these
operation requests about connector to OpenStack osbrick module.

*/

package connector

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/opensds/opensds/cmd/osds_drivers/kubernetes/opensds/api"
)

const (
	URL_PREFIX = "http://localhost:7879"
)

type InitializeRequest struct {
	Multipath bool `json:"multipath"`
}

func GetConnectorProperties(isMultipath bool) (*api.ConnectorProperties, error) {
	url := URL_PREFIX + "/Volume/Initialize"

	initReq := &InitializeRequest{
		Multipath: isMultipath,
	}

	log.Printf("Start POST request to initialize volume, url = %s, body = %+v\n",
		url, initReq)

	req := httplib.Post(url).SetTimeout(100*time.Second, 50*time.Second)
	req.JSONBody(initReq)

	resp, err := req.Response()
	if err != nil {
		return nil, err
	}

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nil, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var prop = &api.ConnectorProperties{}
	err = json.Unmarshal(rbody, prop)
	if err != nil {
		log.Println("The format of response is not supported, resp:", string(rbody))
		return nil, err
	}

	prop.DoLocalAttach = true
	return prop, nil
}

type Connector struct {
	ConnInfo api.ConnectionInfo `json:"connection_info"`
}

func (conn *Connector) ConnectVolume() (string, error) {
	url := URL_PREFIX + "/Volume/Connect"

	log.Printf("Start POST request to connect volume, url = %s, body = %+v\n",
		url, conn)

	req := httplib.Post(url).SetTimeout(100*time.Second, 50*time.Second)
	req.JSONBody(conn)

	resp, err := req.Response()
	if err != nil {
		return "", err
	}

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return "", err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	devPath, err := parseDevicePath(rbody)
	if err != nil {
		return "", err
	}
	return devPath, nil
}

func (conn *Connector) DisconnectVolume() (string, error) {
	url := URL_PREFIX + "/Volume/Disconnect"

	log.Printf("Start POST request to disconnect volume, url = %s, body = %+v\n",
		url, conn)

	req := httplib.Post(url).SetTimeout(100*time.Second, 50*time.Second)
	req.JSONBody(conn)

	resp, err := req.Response()
	if err != nil {
		return "", err
	}

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return "", err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(rbody), nil
}

func parseDevicePath(body []byte) (string, error) {
	var deviceData map[string]string

	if err := json.Unmarshal(body, &deviceData); err != nil {
		log.Println("Unable to parse device path:", err)
		return "", err
	}
	return deviceData["path"], nil
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
