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
	"log"
	"net/http"
	// "time"

	"openstack/golang-client/auth"
	"openstack/golang-client/volume"

	"git.openstack.org/openstack/golang-client/openstack"
)

type CinderPlugin struct {
	Host        string
	Methods     []string
	Username    string
	Password    string
	ProjectId   string
	ProjectName string
}

func (plugin *CinderPlugin) Setup() {

}

func (plugin *CinderPlugin) Unset() {

}

func (plugin *CinderPlugin) CreateVolume(name, volType string, size int32) (string, error) {
	//Get the certified volume service.
	volumeService, err := plugin.getVolumeService()
	if err != nil {
		log.Println("Cannot access volume service:", err)
		return "", err
	}

	//Configure create request body, the body is defined in volume package.
	body := &volume.CreateBody{
		VolumeBody: volume.RequestBody{
			Name:       name,
			VolumeType: volType,
			Size:       size,
		},
	}

	volume, err := volumeService.Create(body)
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

func (plugin *CinderPlugin) AttachVolume(volID string) (string, error) {
	return AttachVolumeToHost(plugin, volID)
}

func (plugin *CinderPlugin) DetachVolume(device string) (string, error) {
	return DetachVolumeFromHost(plugin, device)
}

/*
There is some touble now in getVolumeService(). After setting up OpenSDS

service, this process would dump if any credential works don't work. And

we thought it could be solved by make this function a goroutine.

*/
func (plugin *CinderPlugin) getVolumeService() (volume.Service, error) {
	creds := auth.AuthOpts{
		AuthUrl:     plugin.Host,
		Methods:     plugin.Methods,
		Username:    plugin.Username,
		Password:    plugin.Password,
		ProjectId:   plugin.ProjectId,
		ProjectName: plugin.ProjectName,
	}
	auth, err := auth.DoAuthRequestV3(creds)
	if err != nil {
		log.Fatalln("There was an error authenticating:", err)
	}
	/*
		if !auth.GetExpiration().After(time.Now()) {
			log.Fatalln("There was an error. The auth token has an invalid expiration.")
		}
	*/

	// Find the endpoint for the volume v2 service.
	url, err := auth.GetEndpoint("volumev2", "")
	if url == "" || err != nil {
		log.Fatalln("Volume service url not found during authentication.")
	}

	// Make a new client with these creds, here configure InsecureSkipVerify
	// in tls.Config to skip the certificate verification.
	tls := &tls.Config{
		InsecureSkipVerify: true,
	}

	sess, err := openstack.NewSession(nil, auth, tls)
	if err != nil {
		log.Fatalln("Error creating new Session:", err)
	}

	volumeService := volume.NewService(sess, http.DefaultClient, url)
	return volumeService, nil
}
