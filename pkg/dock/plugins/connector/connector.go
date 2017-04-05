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
	"io/ioutil"
	"log"
	"time"

	"openstack/golang-client/volume"

	"github.com/astaxie/beego/httplib"
)

const (
	URL_PREFIX = "http://localhost:7879"
)

type Connector struct {
	ConnInfo volume.ConnectionInfo `json:"connection_info"`
}

type InitializeRequest struct {
	Multipath bool `json:"multipath"`
}

func GetConnectorProperties(isMultipath bool) (*volume.ConnectorProperties, error) {
	url := URL_PREFIX + "Volume/Initialize"

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

	var prop = &volume.ConnectorProperties{}
	err = json.Unmarshal(rbody, prop)
	if err != nil {
		log.Println("The format of response is not supported, resp:", string(rbody))
		return nil, err
	}

	prop.DoLocalAttach = true
	return prop, nil
}

func (conn *Connector) ConnectVolume() (string, error) {
	url := URL_PREFIX + "Volume/Connect"

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

	return string(rbody), nil
}

func (conn *Connector) DisconnectVolume() (string, error) {
	url := URL_PREFIX + "Volume/Disconnect"

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
