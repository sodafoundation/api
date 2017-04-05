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

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
)

const (
	URL_PREFIX string = "http://10.2.0.115:50048"
)

type OpenSDSOptions struct {
	DefaultOptions
	VolumeId     string `json:"volumeId"`
	ResourceType string `json:"resourceType"`
	VolumeType   string `json:"volumeType"`
}

type OpenSDSPlugin struct{}

func (OpenSDSPlugin) Init() Result {
	return Succeed()
}

func (OpenSDSPlugin) NewOptions() interface{} {
	var option = &OpenSDSOptions{}
	return option
}

func (OpenSDSPlugin) Attach(opts interface{}) Result {
	opt := opts.(*OpenSDSOptions)

	url := URL_PREFIX + "/api/v1/volumes/attach"

	dockId, err := GetDockId()
	if err != nil {
		return Fail(err.Error())
	}

	vr := &VolumeRequest{
		DockId:       dockId,
		ResourceType: opt.ResourceType,
		Id:           opt.VolumeId,
		VolumeType:   opt.VolumeType,
	}

	// fmt.Println("Start PUT request to attach volume, url =", url)
	req := httplib.Put(url).SetTimeout(100*time.Second, 50*time.Second)
	req.JSONBody(vr)

	resp, err := req.Response()
	if err != nil {
		return Fail(err.Error())
	}

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return Fail(err.Error())
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Fail(err.Error())
	}

	var volumeResponse = &VolumeResponse{}
	err = json.Unmarshal(rbody, volumeResponse)
	if err != nil {
		return Fail(err.Error())
	} else {
		if volumeResponse.Status == "Success" {
			return Result{
				Status: "Success",
				Device: volumeResponse.Message,
			}
		} else {
			err = errors.New("Detach volume failed!")
			return Fail(err.Error())
		}
	}
}

func (OpenSDSPlugin) Detach(device string) Result {
	url := URL_PREFIX + "/api/v1/volumes/attach"

	dockId, err := GetDockId()
	if err != nil {
		return Fail(err.Error())
	}

	vr := &VolumeRequest{
		DockId:       dockId,
		ResourceType: "cinder",
	}

	// fmt.Println("Start DELETE request to detach volume, url =", url)
	req := httplib.Delete(url).SetTimeout(100*time.Second, 50*time.Second)
	req.JSONBody(vr)

	resp, err := req.Response()
	if err != nil {
		return Fail(err.Error())
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Fail(err.Error())
	}

	var volumeResponse = &VolumeResponse{}
	err = json.Unmarshal(rbody, volumeResponse)
	if err != nil {
		return Fail(err.Error())
	}
	if strings.Contains(string(rbody), "Success") {
		return Succeed()
	} else {
		err = errors.New("Detach volume failed!")
		return Fail(err.Error())
	}
}

func (OpenSDSPlugin) Mount(mountDir string, device string, opts interface{}) Result {
	opt := opts.(*OpenSDSOptions)

	url := URL_PREFIX + "/api/v1/volumes/action/mount"

	dockId, err := GetDockId()
	if err != nil {
		return Fail(err.Error())
	}

	vr := &VolumeRequest{
		DockId:       dockId,
		ResourceType: opt.ResourceType,
		MountDir:     mountDir,
		Device:       device,
		FsType:       opt.FsType,
	}

	// fmt.Println("Start PUT request to mount volume, url =", url)
	req := httplib.Put(url).SetTimeout(100*time.Second, 50*time.Second)
	req.JSONBody(vr)

	resp, err := req.Response()
	if err != nil {
		return Fail(err.Error())
	}

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return Fail(err.Error())
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Fail(err.Error())
	}
	if strings.Contains(string(rbody), "Success") {
		return Succeed()
	} else {
		err = errors.New("Mount volume failed!")
		return Fail(err.Error())
	}
}

func (OpenSDSPlugin) Unmount(mountDir string) Result {
	url := URL_PREFIX + "/api/v1/volumes/mount"

	dockId, err := GetDockId()
	if err != nil {
		return Fail(err.Error())
	}

	vr := &VolumeRequest{
		DockId:       dockId,
		ResourceType: "cinder",
		MountDir:     mountDir,
	}

	// fmt.Println("Start DELETE request to unmount volume, url =", url)
	req := httplib.Delete(url).SetTimeout(100*time.Second, 50*time.Second)
	req.JSONBody(vr)

	resp, err := req.Response()
	if err != nil {
		return Fail(err.Error())
	}

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return Fail(err.Error())
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Fail(err.Error())
	}
	if strings.Contains(string(rbody), "Success") {
		return Succeed()
	} else {
		err = errors.New("Unmount volume failed!")
		return Fail(err.Error())
	}
}

func main() {
	RunPlugin(&OpenSDSPlugin{})
}
