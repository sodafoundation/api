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
This module implements cinder plugin for OpenSDS. Cinder plugin will pass these
operation requests about volume to OpenStack go-client module.

*/

package cinder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"git.openstack.org/openstack/golang-client.git/openstack"
	"git.openstack.org/openstack/golang-client.git/volume/v3"
)

type CinderPlugin struct {
	Host     string "localhost"
	Username string "cloud_admin"
	Password string "CloudService@123!"
}

func (plugin *CinderPlugin) Setup() {

}

func (plugin *CinderPlugin) Unset() {

}

func (plugin *CinderPlugin) CreateVolume(name string, size int) (string, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		panicString := fmt.Sprint("Cannot access volume service:", err)
		panic(panicString)
	}

	requestBody := v3.RequestBody{name, size}
	volumes, err := volumeService.Create(&requestBody)
	if err != nil {
		panicString := fmt.Sprint("Cannot access volumes:", err)
		panic(panicString)
	}

	a, _ := json.Marshal(volumes)
	return string(a), nil
}

func (plugin *CinderPlugin) GetVolume(volID string) (string, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		panicString := fmt.Sprint("Cannot access volume service:", err)
		panic(panicString)
	}

	volumes, err := volumeService.Show(volID)
	if err != nil {
		panicString := fmt.Sprint("Cannot access volumes:", err)
		panic(panicString)
	}

	a, _ := json.Marshal(volumes)
	return string(a), nil
}

func (plugin *CinderPlugin) GetAllVolumes() (string, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		panicString := fmt.Sprint("Cannot access volume service:", err)
		panic(panicString)
	}

	volumes, err := volumeService.List()
	if err != nil {
		panicString := fmt.Sprint("Cannot access volumes:", err)
		panic(panicString)
	}

	a, _ := json.Marshal(volumes)
	return string(a), nil
}

func (plugin *CinderPlugin) UpdateVolume(volID string, name string) (string, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		panicString := fmt.Sprint("Cannot access volume service:", err)
		panic(panicString)
	}

	requestBody := v3.RequestBody{name, 100}
	volumes, err := volumeService.Update(volID, &requestBody)
	if err != nil {
		panicString := fmt.Sprint("Cannot access volumes:", err)
		panic(panicString)
	}

	a, _ := json.Marshal(volumes)
	return string(a), nil
}

func (plugin *CinderPlugin) DeleteVolume(volID string) (string, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		panicString := fmt.Sprint("Cannot access volume service:", err)
		panic(panicString)
	}

	volumes, err := volumeService.Delete(volID)
	if err != nil {
		panicString := fmt.Sprint("Cannot access volumes:", err)
		panic(panicString)
	}

	a, _ := json.Marshal(volumes)
	return string(a), nil
}

func (plugin *CinderPlugin) getVolumeService() (v3.Service, error) {
	creds := openstack.AuthOpts{
		AuthUrl:  plugin.Host,
		Username: plugin.Username,
		Password: plugin.Password,
	}
	auth, err := openstack.DoAuthRequest(creds)
	if err != nil {
		panicString := fmt.Sprint("There was an error authenticating:", err)
		panic(panicString)
	}
	if !auth.GetExpiration().After(time.Now()) {
		panic("There was an error. The auth token has an invalid expiration.")
	}

	// Find the endpoint for the volume v2 service.
	url, err := auth.GetEndpoint("volumev2", "")
	if url == "" || err != nil {
		panic("v2 volume service url not found during authentication")
	}

	// Make a new client with these creds
	sess, err := openstack.NewSession(nil, auth, nil)
	if err != nil {
		panicString := fmt.Sprint("Error crating new Session:", err)
		panic(panicString)
	}

	volumeService := v3.Service{
		Session: *sess,
		Client:  *http.DefaultClient,
		URL:     url, // We're forcing Volume v2 for now
	}
	return volumeService, nil
}

func (plugin *CinderPlugin) Mount(host string, volID string) {

}

func (plugin *CinderPlugin) Unmount(host string, volID string) {

}
