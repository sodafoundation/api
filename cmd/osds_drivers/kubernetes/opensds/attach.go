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
This module implements a standard SouthBound interface of volume resource to
storage plugins.

*/

package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/astaxie/beego/httplib"

	"github.com/opensds/opensds/cmd/osds_drivers/kubernetes/opensds/api"
)

func CreateVolumeAttachment(volID string, prop *api.ConnectorProperties) (*api.VolumeAttachment, error) {
	url := URL_PREFIX + "/api/v1/volumes/" + volID + "/attachments"
	vr := &api.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			DoLocalAttach: prop.DoLocalAttach,
			MultiPath:     prop.MultiPath,
			HostInfo: api.HostInfo{
				Platform:  prop.Platform,
				OsType:    prop.OsType,
				Ip:        prop.Ip,
				Host:      prop.Host,
				Initiator: prop.Initiator,
			},
		},
	}

	// fmt.Println("Start POST request to create volume attachment, url =", url)
	req := httplib.Post(url).SetTimeout(100*time.Second, 50*time.Second)
	req.JSONBody(vr)

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

	var atc = &api.VolumeAttachment{}
	if err = json.Unmarshal(rbody, atc); err != nil {
		return nil, err
	}
	return atc, nil
}

func GetVolumeAttachment(id, volID string) (*api.VolumeAttachment, error) {
	url := URL_PREFIX + "/api/v1/volumes/" + volID + "/attachments/" + id

	// fmt.Println("Start GET request to get volume attachment, url =", url)
	req := httplib.Get(url).SetTimeout(100*time.Second, 50*time.Second)

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

	var atc = &api.VolumeAttachment{}
	if err = json.Unmarshal(rbody, atc); err != nil {
		return nil, err
	}
	return atc, nil
}

func ListVolumeAttachments(volID string) (*[]api.VolumeAttachment, error) {
	url := URL_PREFIX + "/api/v1/volumes/" + volID + "/attachments"

	// fmt.Println("Start GET request to list volume attachments, url =", url)
	req := httplib.Get(url).SetTimeout(100*time.Second, 50*time.Second)

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

	var atcs = &[]api.VolumeAttachment{}
	if err = json.Unmarshal(rbody, atcs); err != nil {
		return nil, err
	}
	return atcs, nil
}

func UpdateVolumeAttachment(id, volID, mountpoint string, hostInfo api.HostInfo) (*api.VolumeAttachment, error) {
	url := URL_PREFIX + "/api/v1/volumes/" + volID + "/attachments/" + id
	vr := &api.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			HostInfo:   hostInfo,
			Mountpoint: mountpoint,
		},
	}

	// fmt.Println("Start PUT request to update volume attachment, url =", url)
	req := httplib.Put(url).SetTimeout(100*time.Second, 50*time.Second)
	req.JSONBody(vr)

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

	var atc = &api.VolumeAttachment{}
	if err = json.Unmarshal(rbody, atc); err != nil {
		return nil, err
	}
	return atc, nil
}

func DeleteVolumeAttachment(id, volID string) (*api.VolumeResponse, error) {
	url := URL_PREFIX + "/api/v1/volumes/" + volID + "/attachments/" + id

	// fmt.Println("Start DELETE request to delete volume attachment, url =", url)
	req := httplib.Delete(url).SetTimeout(100*time.Second, 50*time.Second)

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

	var volumeResponse = &api.VolumeResponse{}
	err = json.Unmarshal(rbody, volumeResponse)
	if err != nil {
		return nil, err
	}
	return volumeResponse, nil
}
