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
This module implements CoprHD plugin for OpenSDS. CoprHD plugin will pass these
operation requests about volume to REST API.

*/

package coprhd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"gopkg.in/jmcvetta/napping.v3"
)

const (
	// LoginUri path to create a authentication token
	loginUri = "/login.json"
	// CreateVolumeUri path to create volume
	createVolumeUri = "/block/volumes.json"

	projectId = "urn:storageos:Project:dff66e0c-6af1-4bd1-9d22-7dbb355f5c62:global"
	varrayId  = "urn:storageos:VirtualPool:09cbc520-cd92-47a7-83d6-0d25c9cbd053:vdc1"
	vpoolId   = "urn:storageos:VirtualArray:24604053-968f-444e-9797-470b384aaa2e:vdc1"
)

// VolumeArgs represents the json parameters for the volume REST call
type VolumeArgs struct {
	Name    string `json:"name"`
	Size    string `json:"size"`
	Project string `json:"project"`
	VArray  string `json:"varray"`
	VPool   string `json:"vpool"`
}

// CreateAndDeleteVolumeReply is the reply from the volume REST call
type CreateAndDeleteVolumeReply struct {
	Task []struct {
		Inactive bool `json:"inactive"`
		Resource struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"resource"`
	} `json:"task"`
}

type GetVolumeReply struct {
	Name       string `json:"name"`
	Id         string `json:"id"`
	Size       string `json:"provisioned_capacity_gb"`
	Inactive   bool   `json:"inactive"`
	SystemType string `json:"system_type"`
}

type VolumeResponse struct {
	Name       string `json:"name"`
	Id         string `json:"id"`
	Size       int32  `json:"size"`
	Status     string `json:"status"`
	VolumeType string `json:"volume_type"`
}

type BulkReply struct {
	Id []string `json:"id"`
}

type Driver struct {
	Url        string
	Creds      *url.Userinfo
	HttpClient *http.Client
}

func (d *Driver) Setup() {

}

func (d *Driver) Unset() {

}

func (d *Driver) CreateVolume(name, volType string, size int32) (string, error) {

	s, err := d.getAuthSession()

	if err != nil {
		log.Fatal("Failed to create session: %s", err.Error())
	}

	res := &CreateAndDeleteVolumeReply{}

	payload := VolumeArgs{
		Name:    name,                      // Name
		Size:    fmt.Sprintf("%dGB", size), // Volume Size
		Project: projectId,
		VArray:  varrayId,
		VPool:   vpoolId,
	}

	url := d.Url + createVolumeUri

	resp, err := s.Post(url, &payload, res, nil)

	if resp.Status() != http.StatusAccepted {
		return "", fmt.Errorf("Failed to create volume: %s", resp.Result)
	}

	log.Println("Create volume success, dls =", res)

	vres := &VolumeResponse{
		Id:   res.Task[0].Resource.Id,
		Name: res.Task[0].Resource.Name,
		Size: size,
	}
	if res.Task[0].Inactive {
		vres.Status = "inactive"
	} else {
		vres.Status = "active"
	}

	rbody, _ := json.Marshal(vres)
	return string(rbody), nil
}

func (d *Driver) GetVolume(volID string) (string, error) {
	s, err := d.getAuthSession()

	if err != nil {
		log.Fatal("Failed to create session: %s", err.Error())
	}

	res := &GetVolumeReply{}

	url := d.Url + "/block/volumes/" + volID + ".json"

	resp, err := s.Get(url, nil, res, nil)

	if resp.Status() != http.StatusOK {
		return "", fmt.Errorf("Failed to get volume: %s", resp.Result)
	}

	log.Println("Get volume success, dls =", res)

	size, err := strconv.ParseFloat(res.Size, 32)
	if err != nil {
		return "", err
	}

	vres := &VolumeResponse{
		Id:         res.Id,
		Name:       res.Name,
		VolumeType: res.SystemType,
		Size:       int32(size),
	}
	if res.Inactive {
		vres.Status = "inactive"
	} else {
		vres.Status = "active"
	}

	rbody, _ := json.Marshal(vres)
	return string(rbody), nil
}

func (d *Driver) GetAllVolumes(allowDetails bool) (string, error) {
	allowDetails = true
	s, err := d.getAuthSession()

	if err != nil {
		log.Fatal("Failed to create session: %s", err.Error())
	}

	bres := &BulkReply{}

	url := d.Url + "/block/volumes/bulk.json"

	resp, err := s.Get(url, nil, bres, nil)

	if resp.Status() != http.StatusOK {
		return "", fmt.Errorf("Failed to get all volumes: %s", resp.Result)
	}

	log.Println("List volume bulks success, dls =", bres)

	var vress []VolumeResponse

	for _, id := range bres.Id {
		res := &GetVolumeReply{}

		url := d.Url + "/block/volumes/" + id + ".json"

		resp, err = s.Get(url, nil, res, nil)

		if resp.Status() != http.StatusOK {
			return "", fmt.Errorf("Failed to get volume: %s", resp.Result)
		}

		log.Println("Get volume success, dls =", res)

		size, err := strconv.ParseFloat(res.Size, 32)
		if err != nil {
			return "", err
		}

		vres := &VolumeResponse{
			Id:         res.Id,
			Name:       res.Name,
			VolumeType: res.SystemType,
			Size:       int32(size),
		}
		if res.Inactive {
			vres.Status = "inactive"
		} else {
			vres.Status = "active"
		}

		vress = append(vress, *vres)
	}

	log.Println("List volumes success, dls =", vress)

	rbody, _ := json.Marshal(vress)
	return string(rbody), nil
}

func (d *Driver) DeleteVolume(volID string) (string, error) {
	s, err := d.getAuthSession()
	if err != nil {
		log.Fatal("Failed to create session: %s", err.Error())
	}

	res := &CreateAndDeleteVolumeReply{}

	url := d.Url + "/block/volumes/" + volID + "/deactivate.json"

	resp, err := s.Post(url, nil, res, nil)
	if resp.Status() != http.StatusAccepted {
		return "", fmt.Errorf("Failed to delete volume: %s", resp.Result)
	}

	log.Println("Delete success, dls =", res)

	return "Delete volume success!", nil
}

func (d *Driver) AttachVolume(volId string) (string, error) {
	return "", nil
}

func (d *Driver) DetachVolume(device string) (string, error) {
	return "", nil
}

func (d *Driver) CreateSnapshot(name, volID, description string, forced bool) (string, error) {
	return "", nil
}

func (d *Driver) GetSnapshot(snapID string) (string, error) {
	return "", nil
}

func (d *Driver) GetAllSnapshots() (string, error) {
	return "", nil
}

func (d *Driver) DeleteSnapshot(snapID string) (string, error) {
	return "", nil
}

// getAuthSession returns an authenticated API Session
func (d *Driver) getAuthSession() (session *napping.Session, err error) {
	s := napping.Session{
		Userinfo: d.Creds,
		Client:   d.HttpClient,
	}

	url := d.Url + loginUri

	resp, err := s.Get(url, nil, nil, nil)

	if err != nil {
		return
	}

	token := resp.HttpResponse().Header.Get("X-SDS-AUTH-TOKEN")

	h := http.Header{}

	h.Set("X-SDS-AUTH-TOKEN", token)

	session = &napping.Session{
		Client: d.HttpClient,
		Header: &h,
	}

	return
}
