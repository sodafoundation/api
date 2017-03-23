// +build !unit

// Copyright (c) 2014 Hewlett-Packard Development Company, L.P.
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
	"fmt"
	"net/http"
	"time"

	"git.openstack.org/openstack/golang-client/image/v1"
	"git.openstack.org/openstack/golang-client/openstack"
)

// Image examples.
func main() {
	config := getConfig()

	// Authenticate with a username, password, tenant id.
	creds := openstack.AuthOpts{
		AuthUrl:     config.Host,
		ProjectName: config.ProjectName,
		Username:    config.Username,
		Password:    config.Password,
	}
	auth, err := openstack.DoAuthRequest(creds)
	if err != nil {
		panicString := fmt.Sprint("There was an error authenticating:", err)
		panic(panicString)
	}
	if !auth.GetExpiration().After(time.Now()) {
		panic("There was an error. The auth token has an invalid expiration.")
	}

	// Find the endpoint for the image service.
	url, err := auth.GetEndpoint("image", "")
	if url == "" || err != nil {
		panic("v1 image service url not found during authentication")
	}

	// Make a new client with these creds
	sess, err := openstack.NewSession(nil, auth, nil)
	if err != nil {
		panicString := fmt.Sprint("Error crating new Session:", err)
		panic(panicString)
	}

	imageService := image.Service{
		Session: *sess,
		Client:  *http.DefaultClient,
		URL:     url + "/v1", // We're forcing Image v1 for now
	}
	imagesDetails, err := imageService.ImagesDetail()
	if err != nil {
		panicString := fmt.Sprint("Cannot access images:", err)
		panic(panicString)
	}

	var imageIDs = make([]string, 0)
	for _, element := range imagesDetails {
		imageIDs = append(imageIDs, element.ID)
	}

	if len(imageIDs) == 0 {
		panicString := fmt.Sprint("No images found, check to make sure access is correct")
		panic(panicString)
	}
}
