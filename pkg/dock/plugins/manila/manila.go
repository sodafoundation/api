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
This module implements manila plugin for OpenSDS. Manila plugin will pass these
operation requests about file shares to OpenStack go-client module.

*/

package manila

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"openstack/golang-client/share"

	"git.openstack.org/openstack/golang-client.git/openstack"
)

type ManilaPlugin struct {
	Host        string
	Username    string
	Password    string
	ProjectName string
}

func (plugin *ManilaPlugin) Setup() {

}

func (plugin *ManilaPlugin) Unset() {

}

func (plugin *ManilaPlugin) CreateShare(name string, shrType string, shrProto string, size int32) (string, error) {
	//Get the certified file share service.
	shareService, err := plugin.getShareService()
	if err != nil {
		log.Println("Cannot access file share service:", err)
		return "", err
	}

	//Configure HTTP request body, the body is defined in share package.
	requestBody := share.RequestBody{}
	requestBody.Name = name
	requestBody.Size = size
	requestBody.Share_proto = shrProto
	requestBody.Share_type = shrType
	body := share.CreateBody{requestBody}
	share, err := shareService.Create(&body)
	if err != nil {
		log.Println("Cannot create file share:", err)
		return "", err
	}

	a, _ := json.Marshal(share)
	result := string(a)
	log.Println("Create file share success, dls =", result)
	return result, nil
}

func (plugin *ManilaPlugin) GetShare(shrID string) (string, error) {
	shareService, err := plugin.getShareService()
	if err != nil {
		log.Println("Cannot access file share service:", err)
		return "", err
	}

	share, err := shareService.Show(shrID)
	if err != nil {
		log.Println("Cannot show file share:", err)
		return "", err
	}

	a, _ := json.Marshal(share)
	result := string(a)
	log.Println("Get file share success, dls =", result)
	return result, nil
}

func (plugin *ManilaPlugin) GetAllShares(allowDetails bool) (string, error) {
	shareService, err := plugin.getShareService()
	if err != nil {
		log.Println("Cannot access file share service:", err)
		return "", err
	}

	var shares interface{}
	if allowDetails {
		shares, err = shareService.Detail()
		if err != nil {
			log.Println("Cannot detail file shares:", err)
			return "", err
		}
	} else {
		shares, err = shareService.List()
		if err != nil {
			log.Println("Cannot list file shares:", err)
			return "", err
		}
	}

	a, _ := json.Marshal(shares)
	result := string(a)
	log.Println("Get all file shares success, dls =", result)
	return result, nil
}

func (plugin *ManilaPlugin) DeleteShare(shrID string) (string, error) {
	shareService, err := plugin.getShareService()
	if err != nil {
		log.Println("Cannot access file share service:", err)
		return "", err
	}

	_, err = shareService.Show(shrID)
	if err != nil {
		log.Println("Cannot get file share:", err)
		return "", err
	}

	err = shareService.Delete(shrID)
	if err != nil {
		log.Println("Cannot delete file share:", err)
		return "", err
	}

	result := "Delete file share success!"
	return result, nil
}

func (plugin *ManilaPlugin) getShareService() (share.Service, error) {
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

	// Find the endpoint for the file share service.
	url, err := auth.GetEndpoint("sharev2", "")
	if url == "" || err != nil {
		log.Fatalln("Share service url not found during authentication.")
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

	shareService, _ := share.NewService(*sess, *http.DefaultClient, url)
	return shareService, nil
}
