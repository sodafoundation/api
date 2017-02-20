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
	"errors"
	"log"
	"net/http"
	"time"

	"openstack/golang-client/volume"

	"git.openstack.org/openstack/golang-client.git/openstack"
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
		log.Println("Cannot access volume service:", err)
		return "", err
	}

	//Configure HTTP request body, the body is defined in volume package.
	requestBody := volume.RequestBody{}
	requestBody.Name = name
	requestBody.Size = size
	body := volume.CreateBody{requestBody}
	volume, err := volumeService.Create(&body)
	if err != nil {
		log.Println("Cannot create volume:", err)
		return "", err
	}

	a, _ := json.Marshal(volume)
	result := string(a)
	log.Println("Create volume success, dls =", result)
	return result, nil
}

func (plugin *CinderPlugin) GetVolume(volID string) (string, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return "", err
	}

	volume, err := volumeService.Show(volID)
	if err != nil {
		log.Println("Cannot show volume:", err)
		return "", err
	}

	a, _ := json.Marshal(volume)
	result := string(a)
	log.Println("Get volume success, dls =", result)
	return result, nil
}

func (plugin *CinderPlugin) GetAllVolumes(allowDetails bool) (string, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return "", err
	}

	var volumes interface{}
	if allowDetails {
		volumes, err = volumeService.Detail()
		if err != nil {
			log.Println("Cannot detail volumes:", err)
			return "", err
		}
	} else {
		volumes, err = volumeService.List()
		if err != nil {
			log.Println("Cannot list volumes:", err)
			return "", err
		}
	}

	a, _ := json.Marshal(volumes)
	result := string(a)
	log.Println("Get all volumes success, dls =", result)
	return result, nil
}

func (plugin *CinderPlugin) UpdateVolume(volID string, name string) (string, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return "", err
	}

	requestBody := volume.RequestBody{}
	requestBody.Name = name
	body := volume.CreateBody{requestBody}
	volume, err := volumeService.Update(volID, &body)
	if err != nil {
		log.Println("Cannot update volume:", err)
		return "", err
	}

	a, _ := json.Marshal(volume)
	result := string(a)
	log.Println("Update volume success, dls =", result)
	return result, nil
}

func (plugin *CinderPlugin) DeleteVolume(volID string) (string, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return "", err
	}

	_, err = volumeService.Show(volID)
	if err != nil {
		log.Println("Cannot get volume:", err)
		return "", err
	}

	err = volumeService.Delete(volID)
	if err != nil {
		log.Println("Cannot delete volume:", err)
		return "", err
	}

	result := "Delete volume success!"
	return result, nil
}

func (plugin *CinderPlugin) AttachVolume(volID, host, device string) (string, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return "", err
	}

	vol, err := volumeService.Show(volID)
	if err != nil {
		log.Println("Cannot get volume:", err)
		return "", err
	}
	if vol.Status != "available" {
		err = errors.New("The status of volume is not available!")
		log.Println("Cannot attach volume:", err)
		return "", err
	}

	requestBody := volume.RequestBody{}
	requestBody.HostName = host
	requestBody.Device = device
	body := volume.AttachBody{requestBody}
	err = volumeService.Attach(volID, &body)
	if err != nil {
		log.Println("Cannot attach volume:", err)
		return "", err
	}

	result := "Attach volume success!"
	return result, nil
}

func (plugin *CinderPlugin) DetachVolume(volID, attachment string) (string, error) {
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return "", err
	}

	vol, err := volumeService.Show(volID)
	if err != nil {
		log.Println("Cannot get volume:", err)
		return "", err
	}
	if vol.Status != "in-use" {
		err = errors.New("The status of volume is not in-use!")
		log.Println("Cannot attach volume:", err)
		return "", err
	}

	requestBody := volume.RequestBody{}
	requestBody.AttachmentID = attachment
	body := volume.DetachBody{requestBody}
	err = volumeService.Detach(volID, &body)
	if err != nil {
		log.Println("Cannot detach volume:", err)
		return "", err
	}

	result := "Detach volume success!"
	return result, nil
}

/*
There is some touble now in getVolumeService(). After setting up OpenSDS

service, this process would dump if any credential works don't work. And

we thought it could be solved by make this function a goroutine.

*/
func (plugin *CinderPlugin) getVolumeService() (volume.Service, error) {
	creds := openstack.AuthOpts{
		AuthUrl:     plugin.Host,
		Username:    plugin.Username,
		Password:    plugin.Password,
		ProjectName: plugin.ProjectName,
	}
	auth, err := openstack.DoAuthRequest(creds)
	if err != nil {
		log.Fatalln("There was an error authenticating:", err)
	}
	if !auth.GetExpiration().After(time.Now()) {
		log.Fatalln("There was an error. The auth token has an invalid expiration.")
	}

	// Find the endpoint for the volume service.
	url, err := auth.GetEndpoint("volumev2", "")
	if url == "" || err != nil {
		log.Fatalln("Volume service url not found during authentication.")
	}

	// Make a new client with these creds, here configure InsecureSkipVerify
	// in tls.Config to skip the certificate verification.
	tls := &tls.Config{}
	tls.InsecureSkipVerify = true
	sess, err := openstack.NewSession(nil, auth, tls)
	if err != nil {
		log.Fatalln("Error creating new Session:", err)
	}

	volumeService, _ := volume.NewService(*sess, *http.DefaultClient, url)
	return volumeService, nil
}
