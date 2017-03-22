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
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
)

const (
	URL_PREFIX         string = "http://10.2.0.115:50048"
	CEPH_POOL_NAME     string = "volumes"
	CEPH_LINK_PREFIX   string = "volumes/volume-"
	CINDER_LINK_PREFIX string = "/dev/mapper/cinder--volumes-volume--"
	MANILA_LINK_PREFIX string = "/var/lib/manila/mnt/share-"
	DEVICE_PREFIX      string = "/dev/"
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

	switch opt.ResourceType {
	case "cinder":
		return cinderAttach(opt)
	case "manila":
		return manilaAttach(opt)
	default:
		err := errors.New("Backend resource not supported!")
		return Fail(err.Error())
	}
}

func cinderAttach(opt *OpenSDSOptions) Result {
	volId := opt.VolumeId
	volType := opt.VolumeType
	url := URL_PREFIX + "/api/v1/volumes/action/cinder/" + volId

	var path string
	switch volType {
	case "lvm":
		device, err := getLvmDevicePath(volId)
		if err != nil {
			return Fail(err.Error())
		}

		path = string(device)
	case "ceph":
		device, err := getCephDevicePath(volId)
		if err != nil {
			return Fail(err.Error())
		}

		path = device
	}

	dockId, err := GetDockId()
	if err != nil {
		return Fail(err.Error())
	}

	var volumeRequest VolumeRequest
	volumeRequest.DockId = dockId
	volumeRequest.ActionType = "attach"
	volumeRequest.Device = path

	// fmt.Println("Start POST request to attach volume, url =", url)

	req := httplib.Post(url).SetTimeout(100*time.Second, 50*time.Second)

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
	if strings.Contains(string(rbody), "Success") {
		return Result{
			Status: "Success",
			Device: path,
		}
	} else {
		err = errors.New("Attach volume failed!")
		return Fail(err.Error())
	}
}

func manilaAttach(opt *OpenSDSOptions) Result {
	shrId := opt.VolumeId
	url := URL_PREFIX + "/api/v1/shares/manila/" + shrId

	req := httplib.Get(url).SetTimeout(100*time.Second, 50*time.Second)

	resp, err := req.Response()
	if err != nil {
		return Fail(err.Error())
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Fail(err.Error())
	}
	sdr := new(ShareDetailResponse)
	err = json.Unmarshal(rbody, sdr)
	if err != nil {
		return Fail(err.Error())
	}

	if sdr.ExportLocation == "" {
		err = errors.New("Share not exported!")
		return Fail(err.Error())
	} else {
		return Result{
			Status: "Success",
			Device: sdr.ExportLocation,
		}
	}
}

func (OpenSDSPlugin) Detach(device string) Result {
	if strings.Contains(device, MANILA_LINK_PREFIX) {
		return Succeed()
	}

	var volId string

	if strings.HasPrefix(device, CINDER_LINK_PREFIX) {
		volumeId := device[len(CINDER_LINK_PREFIX):len(device)]
		volId = strings.Replace(volumeId, "--", "-", 4)
	} else {
		if strings.HasPrefix(device, "/dev/rbd") {
			image, err := parseDevicePath(device)
			if err != nil {
				return Fail(err.Error())
			}

			volId = image[7:len(image)]

			unmapCmd := exec.Command("rbd", "unmap", device)
			_, err = unmapCmd.CombinedOutput()
			if err != nil {
				return Fail(err.Error())
			}
		} else {
			err := errors.New("Unexpect device prefix: " + device)
			return Fail(err.Error())
		}
	}

	url := URL_PREFIX + "/api/v1/volumes/cinder/" + volId

	// fmt.Println("Start GET request to get volume, url =", url)

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

	url = URL_PREFIX + "/api/v1/volumes/action/cinder/" + volId
	dockId, err := GetDockId()
	if err != nil {
		return Fail(err.Error())
	}

	var volumeRequest VolumeRequest
	volumeRequest.DockId = dockId
	volumeRequest.ActionType = "detach"
	volumeRequest.Attachment = volumeResponse.Attachments[0]["attachment_id"]

	// fmt.Println("Start POST request to detach volume, url =", url)

	req = httplib.Post(url).SetTimeout(10*time.Second, 5*time.Second)

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
	if strings.Contains(string(rbody), "Success") {
		return Succeed()
	} else {
		err = errors.New("Detach volume failed!")
		return Fail(err.Error())
	}
}

func (OpenSDSPlugin) Mount(mountDir string, device string, opts interface{}) Result {
	opt := opts.(*OpenSDSOptions)

	volId := opt.VolumeId
	url := URL_PREFIX + "/api/v1/volumes/action/cinder/" + volId

	// fmt.Println("Start POST request to mount volume, url =", url)

	req := httplib.Post(url).SetTimeout(100*time.Second, 50*time.Second)

	dockId, err := GetDockId()
	if err != nil {
		return Fail(err.Error())
	}

	var volumeRequest VolumeRequest
	volumeRequest.DockId = dockId
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
	if strings.Contains(string(rbody), "Success") {
		return Succeed()
	} else {
		err = errors.New("Mount volume failed!")
		return Fail(err.Error())
	}
}

func (OpenSDSPlugin) Unmount(mountDir string) Result {
	volId := "null"
	url := URL_PREFIX + "/api/v1/volumes/action/cinder/" + volId

	// fmt.Println("Start POST request to unmount volume, url =", url)

	req := httplib.Post(url).SetTimeout(100*time.Second, 50*time.Second)

	dockId, err := GetDockId()
	if err != nil {
		return Fail(err.Error())
	}

	var volumeRequest VolumeRequest
	volumeRequest.DockId = dockId
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

func getLvmDevicePath(volId string) (string, error) {
	link := CINDER_LINK_PREFIX + strings.Replace(volId, "-", "--", 4)

	path, err := os.Readlink(link)
	if err != nil {
		err = errors.New("Can't find device path!")
		return "", err
	}

	slice := strings.Split(path, "/")
	device := "/dev/" + slice[1]
	return device, nil
}

func getCephDevicePath(volId string) (string, error) {
	imagesCmd := exec.Command("rbd", "showmapped")
	images, err := imagesCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	if strings.Contains(string(images), volId) {
		err := errors.New("This volume have been attached!")
		return "", err
	}

	pathCmd := exec.Command("rbd", "map", CEPH_LINK_PREFIX+volId)
	device, err := pathCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	devSlice := strings.Split(string(device), "\n")
	return devSlice[0], nil
}

func parseDevicePath(device string) (string, error) {
	linksCmd := exec.Command("rbd", "showmapped")
	links, err := linksCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	linkSlice := strings.Split(string(links), "\n")
	for _, i := range linkSlice {
		if strings.Contains(i, device) {
			imageSlice := strings.Fields(i)
			return imageSlice[2], nil
		}
	}

	err = errors.New("Can't parse device path!")
	return "", err
}
