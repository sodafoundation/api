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

// The acceptance package is a set of acceptance tests showcasing how the
// contents of the package are meant to be used. This is setup in a similar
// manner to a consuming application.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"git.openstack.org/openstack/golang-client/openstack"
)

// testconfig contains the user information needed by the acceptance and
// integration tests.
type testconfig struct {
	Host        string
	Username    string
	Password    string
	ProjectID   string
	ProjectName string
	Container   string
	ImageRegion string
	Debug       bool
}

// getConfig provides access to credentials in other tests and examples.
func getConfig() *testconfig {
	config := &testconfig{}
	userJSON, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("ReadFile json failed")
	}
	if err = json.Unmarshal(userJSON, &config); err != nil {
		panic("Unmarshal json failed")
	}

	// Propagate debug setting to packages
	openstack.Debug = &config.Debug
	fmt.Printf("config: %+v\n", config)
	return config
}
