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
	"fmt"
	"log"
	"net/http"
	"net/url"

	"gopkg.in/jmcvetta/napping.v3"
)

const (
	// LoginUri path to create a authentication token
	loginUri = "login.json"
	// CreateVolumeUri path to create volume
	createVolumeUri = "block/volumes.json"
)

// VolumeArgs represents the json parameters for the volume REST call
type VolumeArgs struct {
	Name string `json:"name"`
	Size string `json:"size"`
	// Project string `json:"project"`
	// VArray  string `json:"varray"`
	// VPool   string `json:"vpool"`
}

// VolumeReply is the reply from the volume REST call
type VolumeReply struct {
	Task []struct {
		Resource struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"resource"`
	} `json:"task"`
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

func (d *Driver) CreateVolume(name string, size int) (string, error) {

	s, err := d.getAuthSession()

	if err != nil {
		log.Fatal("Failed to create session: %s", err.Error())
	}

	res := &VolumeReply{}

	payload := VolumeArgs{
		name, // Name
		fmt.Sprintf("%.6fGB", size), // Volume Size
	}

	url := d.Url + createVolumeUri

	resp, err := s.Post(url, &payload, res, nil)

	if resp.Status() != http.StatusAccepted {

		return "", fmt.Errorf("Failed to create volume: %s", resp.Status())
	}

	return res.Task[0].Resource.Id, err
}

func (d *Driver) GetVolume(volID string) (string, error) {
	s, err := d.getAuthSession()

	if err != nil {
		log.Fatal("Failed to create session: %s", err.Error())
	}

	res := &VolumeReply{}

	url := d.Url + "block/volumes" + volID + ".json"

	resp, err := s.Get(url, nil, res, nil)

	if resp.Status() != http.StatusAccepted {

		return "", fmt.Errorf("Failed to get volume: %s", resp.Status())
	}

	return res.Task[0].Resource.Id, err
}

func (d *Driver) GetAllVolumes(allowDetails bool) (string, error) {
	allowDetails = true
	s, err := d.getAuthSession()

	if err != nil {
		log.Fatal("Failed to create session: %s", err.Error())
	}

	res := &VolumeReply{}

	url := d.Url + "block/volumes.json"

	resp, err := s.Get(url, nil, res, nil)

	if resp.Status() != http.StatusAccepted {

		return "", fmt.Errorf("Failed to get all volumes: %s", resp.Status())
	}

	return res.Task[0].Resource.Id, err
	return "", nil
}

func (d *Driver) UpdateVolume(volID string, name string) (string, error) {
	return "", nil
}

func (d *Driver) DeleteVolume(volumeID string) (string, error) {
	return "", nil
}

func (d *Driver) Mount(host string, volID string) {

}

func (d *Driver) Unmount(host string, volID string) {

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
