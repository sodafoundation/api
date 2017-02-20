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
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
)

const DEVICE_PREFIX string = "/dev/mapper/vg01-volume--"

type OpenSDSOptions struct {
	DefaultOptions
	VolumeId string `json:"volumeId"`
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

	volId := opt.VolumeId
	url := "http://162.3.140.36:8080/api/v1/volumes/action/cinder/" + volId

	fmt.Println("Start POST request to attach volume, url =", url)

	req := httplib.Post(url).SetTimeout(100*time.Second, 50*time.Second)

	var volumeRequest VolumeRequest
	volumeRequest.ActionType = "attach"
	volumeRequest.Device = DEVICE_PREFIX + volId

	req.JSONBody(volumeRequest)
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
	if strings.Contains(string(rbody), "success") {
		volumeId := strings.Replace(opt.VolumeId, "-", "--", 4)
		return Result{
			Status: "Success",
			Device: DEVICE_PREFIX + volumeId,
		}
	} else {
		err = errors.New("Attach volume failed!")
		return Fail(err.Error())
	}
}

func (OpenSDSPlugin) Detach(device string) Result {
	volumeId := device[len(device)-40 : len(device)]
	volId := strings.Replace(string(volumeId), "--", "-", 4)
	url := "http://162.3.140.36:8080/api/v1/volumes/action/cinder/" + volId

	fmt.Println("Start GET request to get volume, url =", url)

	req := httplib.Get(url).SetTimeout(100*time.Second, 50*time.Second)

	resp, err := req.Response()
	if err != nil {
		return Fail(err.Error())
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Fail(err.Error())
	}
	volumeResponse := new(VolumeResponse)
	err = json.Unmarshal(rbody, volumeResponse)
	if err != nil {
		return Fail(err.Error())
	}

	if volumeResponse.Status != "in-use" {
		err = errors.New("The status of volume is not in-use!")
		return Fail(err.Error())
	}

	url = url + "/action"

	fmt.Println("Start POST request to detach volume, url =", url)

	req = httplib.Post(url).SetTimeout(10*time.Second, 5*time.Second)

	var volumeRequest VolumeRequest
	volumeRequest.ActionType = "detach"
	volumeRequest.Attachment = volumeResponse.Attachments[0]["attachment_id"]

	req.JSONBody(volumeRequest)
	resp, err = req.Response()
	if err != nil {
		return Fail(err.Error())
	}

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return Fail(err.Error())
	}

	rbody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return Fail(err.Error())
	}
	if strings.Contains(string(rbody), "success") {
		return Succeed()
	} else {
		err = errors.New("Detach volume failed!")
		return Fail(err.Error())
	}
}

func (OpenSDSPlugin) Mount(mountDir string, device string, opts interface{}) Result {
	opt := opts.(*OpenSDSOptions)

	volId := opt.VolumeId
	url := "http://162.3.140.36:8080/api/v1/volumes/action/cinder/" + volId

	fmt.Println("Start POST request to mount volume, url =", url)

	req := httplib.Post(url).SetTimeout(100*time.Second, 50*time.Second)

	var volumeRequest VolumeRequest
	volumeRequest.ActionType = "mount"
	volumeRequest.MountDir = mountDir
	volumeRequest.Device = device
	volumeRequest.FsType = opt.FsType

	req.JSONBody(volumeRequest)
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
	if strings.Contains(string(rbody), "success") {
		return Succeed()
	} else {
		err = errors.New("Mount volume failed!")
		return Fail(err.Error())
	}
}

func (OpenSDSPlugin) Unmount(mountDir string) Result {
	volId := "null"
	url := "http://162.3.140.36:8080/api/v1/volumes/action/cinder/" + volId

	fmt.Println("Start POST request to unmount volume, url =", url)

	req := httplib.Post(url).SetTimeout(100*time.Second, 50*time.Second)

	var volumeRequest VolumeRequest
	volumeRequest.ActionType = "unmount"
	volumeRequest.MountDir = mountDir

	req.JSONBody(volumeRequest)
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
	if strings.Contains(string(rbody), "success") {
		return Succeed()
	} else {
		err = errors.New("Unmount volume failed!")
		return Fail(err.Error())
	}
}

func main() {
	RunPlugin(&OpenSDSPlugin{})
}
