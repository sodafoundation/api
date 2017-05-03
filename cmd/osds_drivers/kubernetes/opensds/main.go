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
	api "github.com/opensds/opensds/pkg/api/v1"
)

const (
	URL_PREFIX string = "http://10.169.149.191:50040"
)

type OpenSDSOptions struct {
	DefaultOptions
	VolumeId      string `json:"volumeId"`
	BackendDriver string `json:"backendDriver"`
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
		Schema: &api.VolumeOperationSchema{
			DockId: dockId,
			Id:     opt.VolumeId,
		},
		Profile: &api.StorageProfile{
			BackendDriver: opt.BackendDriver,
		},
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
	linkPath, err := FindLinkPath(device)
	if err != nil {
		return Fail(err.Error())
	}
	backendDriver, err := FindBackendDriver(device)
	if err != nil {
		return Fail(err.Error())
	}

	vr := &VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			DockId: dockId,
			Device: linkPath,
		},
		Profile: &api.StorageProfile{
			BackendDriver: backendDriver,
		},
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

	url := URL_PREFIX + "/api/v1/volumes/mount"

	dockId, err := GetDockId()
	if err != nil {
		return Fail(err.Error())
	}

	vr := &VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			DockId:   dockId,
			FsType:   opt.FsType,
			Device:   device,
			MountDir: mountDir,
		},
		Profile: &api.StorageProfile{},
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
		Schema: &api.VolumeOperationSchema{
			DockId:   dockId,
			MountDir: mountDir,
		},
		Profile: &api.StorageProfile{},
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
