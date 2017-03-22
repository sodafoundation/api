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
	// "git.openstack.org/openstack/golang-client/identity/v2"
	// "time"

	"git.openstack.org/openstack/golang-client/openstack/v3"
)

// Authentication examples.
func main() {
	config := getConfig()

	// Authenticate with a project id, username, password.
	creds := v3.AuthOpts{
		AuthUrl:   config.Host,
		Methods:   config.Methods,
		ProjectId: config.ProjectID,
		Username:  config.Username,
		Password:  config.Password,
	}
	auth, err := v3.DoAuthRequestV3(creds)
	if err != nil {
		fmt.Println("Error authenticating project/username/password:", err)
		return
	}

	// Get the endpoint with service type and region name.
	ept, err := auth.GetEndpoint("volumev3", "RegionOne")
	if err != nil {
		fmt.Println("No volume endpoint found:", err)
		return
	} else {
		fmt.Printf("Endpoint is: %+v\n", ept)
	}
}
