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
	"time"

	"git.openstack.org/openstack/golang-client/openstack"
)

// Authentication examples.
func main() {
	config := getConfig()

	// Authenticate with just a username and password. The returned token is
	// unscoped to a tenant.
	creds := openstack.AuthOpts{
		AuthUrl:  config.Host,
		Username: config.Username,
		Password: config.Password,
	}
	auth, err := openstack.DoAuthRequest(creds)
	if err != nil {
		fmt.Println("Error authenticating username/password:", err)
		return
	}
	if !auth.GetExpiration().After(time.Now()) {
		fmt.Println("There was an error. The auth token has an invalid expiration.")
		return
	}

	// Authenticate with a project name, username, password.
	creds = openstack.AuthOpts{
		AuthUrl:     config.Host,
		ProjectName: config.ProjectName,
		Username:    config.Username,
		Password:    config.Password,
	}
	auth, err = openstack.DoAuthRequest(creds)
	if err != nil {
		fmt.Println("Error authenticating project/username/password:", err)
		return
	}
	if !auth.GetExpiration().After(time.Now()) {
		fmt.Println("There was an error. The auth token has an invalid expiration.")
		return
	}

	// Authenticate with a project id, username, password.
	creds = openstack.AuthOpts{
		AuthUrl:   config.Host,
		ProjectId: config.ProjectID,
		Username:  config.Username,
		Password:  config.Password,
	}
	auth, err = openstack.DoAuthRequest(creds)
	if err != nil {
		fmt.Println("Error authenticating project/username/password:", err)
		return
	}
	if !auth.GetExpiration().After(time.Now()) {
		fmt.Println("There was an error. The auth token has an invalid expiration.")
		return
	}

	// Get the first endpoint
	_, err = auth.GetEndpoint("compute", "")
	if err != nil {
		fmt.Println("No compute endpoint found:", err)
		return
	}
}
