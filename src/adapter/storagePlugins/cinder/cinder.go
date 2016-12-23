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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"git.openstack.org/openstack/golang-client.git/openstack"
	"git.openstack.org/openstack/golang-client.git/volume/v3"
)

type CinderPlugin struct {
	Host        string
	Username    string
	Password    string
	ProjectName string
}

func (plugin *CinderPlugin) Setup() {

}

func (plugin *CinderPlugin) Unset() {

}

func (plugin *CinderPlugin) CreateVolume(name string, size int) (string, error) {
	//Get the certified volume service.
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		panicString := fmt.Sprint("Cannot access volume service:", err)
		panic(panicString)
	}

	//Configure HTTP request body, the body is defined in v3 package.
	requestBody := v3.RequestBody{name, size}
	body := v3.Body{requestBody}
	volumes, err := volumeService.Create(&body)
	if err != nil {
		panicString := fmt.Sprint("Cannot create volume:", err)
		panic(panicString)
	}

	a, _ := json.Marshal(volumes)
	result := fmt.Sprint("Create volume success!\n", string(a))
	return result, nil
}

func (plugin *CinderPlugin) GetVolume(volID string) (string, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		panicString := fmt.Sprint("Cannot access volume service:", err)
		panic(panicString)
	}

	volumes, err := volumeService.Show(volID)
	if err != nil {
		panicString := fmt.Sprint("Cannot show volume:", err)
		panic(panicString)
	}

	a, _ := json.Marshal(volumes)
	return string(a), nil
}

func (plugin *CinderPlugin) GetAllVolumes(allowDetails bool) (string, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		panicString := fmt.Sprint("Cannot access volume service:", err)
		panic(panicString)
	}

	var volumes interface{}
	if allowDetails {
		volumes, err = volumeService.Detail()
		if err != nil {
			panicString := fmt.Sprint("Cannot detail volumes:", err)
			panic(panicString)
		}
	} else {
		volumes, err = volumeService.List()
		if err != nil {
			panicString := fmt.Sprint("Cannot list volumes:", err)
			panic(panicString)
		}
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

	requestBody := v3.RequestBody{name, 0}
	body := v3.Body{requestBody}
	volumes, err := volumeService.Update(volID, &body)
	if err != nil {
		panicString := fmt.Sprint("Cannot update volume:", err)
		panic(panicString)
	}

	a, _ := json.Marshal(volumes)
	result := fmt.Sprint("Update volume success!\n", string(a))
	return result, nil
}

func (plugin *CinderPlugin) DeleteVolume(volID string) (string, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		panicString := fmt.Sprint("Cannot access volume service:", err)
		panic(panicString)
	}

	err = volumeService.Delete(volID)
	if err != nil {
		panicString := fmt.Sprint("Cannot delete volume:", err)
		panic(panicString)
	}

	resp := "Delete volume success!"
	return resp, nil
}

func (plugin *CinderPlugin) getVolumeService() (v3.Service, error) {
	creds := openstack.AuthOpts{
		AuthUrl:     plugin.Host,
		Username:    plugin.Username,
		Password:    plugin.Password,
		ProjectName: plugin.ProjectName,
	}
	auth, err := openstack.DoAuthRequest(creds)
	if err != nil {
		panicString := fmt.Sprint("There was an error authenticating:", err)
		panic(panicString)
	}
	if !auth.GetExpiration().After(time.Now()) {
		panic("There was an error. The auth token has an invalid expiration.")
	}

	// Find the endpoint for the volume service.
	url, err := auth.GetEndpoint("volumev2", "")
	if url == "" || err != nil {
		panic("Volume service url not found during authentication")
	}

	// Make a new client with these creds, here configure InsecureSkipVerify
	// in tls.Config to skip the certificate verification.
	tls := &tls.Config{}
	tls.InsecureSkipVerify = true
	sess, err := openstack.NewSession(nil, auth, tls)
	if err != nil {
		panicString := fmt.Sprint("Error crating new Session:", err)
		panic(panicString)
	}

	volumeService := v3.Service{
		Session: *sess,
		Client:  *http.DefaultClient,
		URL:     url,
	}
	return volumeService, nil
}

func (plugin *CinderPlugin) Mount(host string, volID string) {

}

func (plugin *CinderPlugin) Unmount(host string, volID string) {

}
