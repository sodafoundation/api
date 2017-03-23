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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"git.openstack.org/openstack/golang-client/objectstorage/v1"
	"git.openstack.org/openstack/golang-client/openstack"
)

func main() {
	config := getConfig()

	// Before working with object storage we need to authenticate with a project
	// that has active object storage.
	// Authenticate with a project name, username, password.
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

	// Find the endpoint for object storage.
	url, err := auth.GetEndpoint("object-store", "")
	if url == "" || err != nil {
		panic("object-store url not found during authentication")
	}

	// Make a new client with these creds
	sess, err := openstack.NewSession(nil, auth, nil)
	if err != nil {
		panicString := fmt.Sprint("Error crating new Session:", err)
		panic(panicString)
	}

	hdr, err := objectstorage.GetAccountMeta(sess, url)
	if err != nil {
		panicString := fmt.Sprint("There was an error getting account metadata:", err)
		panic(panicString)
	}
	_ = hdr

	// Create a new container.
	var headers http.Header = http.Header{}
	headers.Add("X-Log-Retention", "true")
	if err = objectstorage.PutContainer(sess, url+"/"+config.Container, headers); err != nil {
		panicString := fmt.Sprint("PutContainer Error:", err)
		panic(panicString)
	}

	// Get a list of all the containers at the selected endoint.
	containersJson, err := objectstorage.ListContainers(sess, 0, "", url)
	if err != nil {
		panic(err)
	}

	type containerType struct {
		Name         string
		Bytes, Count int
	}
	containersList := []containerType{}

	if err = json.Unmarshal(containersJson, &containersList); err != nil {
		panic(err)
	}

	found := false
	for i := 0; i < len(containersList); i++ {
		if containersList[i].Name == config.Container {
			found = true
		}
	}
	if !found {
		panic("Created container is missing from downloaded containersList")
	}

	// Set and Get container metadata.
	headers = http.Header{}
	headers.Add("X-Container-Meta-fubar", "false")
	if err = objectstorage.SetContainerMeta(sess, url+"/"+config.Container, headers); err != nil {
		panic(err)
	}

	hdr, err = objectstorage.GetContainerMeta(sess, url+"/"+config.Container)
	if err != nil {
		panicString := fmt.Sprint("GetContainerMeta Error:", err)
		panic(panicString)
	}
	if hdr.Get("X-Container-Meta-fubar") != "false" {
		panic("container meta does not match")
	}

	// Create an object in a container.
	var fContent []byte
	srcFile := "10-objectstore.go"
	fContent, err = ioutil.ReadFile(srcFile)
	if err != nil {
		panic(err)
	}

	headers = http.Header{}
	headers.Add("X-Container-Meta-fubar", "false")
	object := config.Container + "/" + srcFile
	if err = objectstorage.PutObject(sess, &fContent, url+"/"+object, headers); err != nil {
		panic(err)
	}
	objectsJson, err := objectstorage.ListObjects(sess, 0, "", "", "", "",
		url+"/"+config.Container)

	type objectType struct {
		Name, Hash, Content_type, Last_modified string
		Bytes                                   int
	}
	objectsList := []objectType{}

	if err = json.Unmarshal(objectsJson, &objectsList); err != nil {
		panic(err)
	}
	found = false
	for i := 0; i < len(objectsList); i++ {
		if objectsList[i].Name == srcFile {
			found = true
		}
	}
	if !found {
		panic("created object is missing from the objectsList")
	}

	// Manage object metadata
	headers = http.Header{}
	headers.Add("X-Object-Meta-fubar", "true")
	if err = objectstorage.SetObjectMeta(sess, url+"/"+object, headers); err != nil {
		panicString := fmt.Sprint("SetObjectMeta Error:", err)
		panic(panicString)
	}
	hdr, err = objectstorage.GetObjectMeta(sess, url+"/"+object)
	if err != nil {
		panicString := fmt.Sprint("GetObjectMeta Error:", err)
		panic(panicString)
	}
	if hdr.Get("X-Object-Meta-fubar") != "true" {
		panicString := fmt.Sprint("SetObjectMeta Error:", err)
		panic(panicString)
	}

	// Retrieve an object and check that it is the same as what as uploaded.
	_, body, err := objectstorage.GetObject(sess, url+"/"+object)
	if err != nil {
		panicString := fmt.Sprint("GetObject Error:", err)
		panic(panicString)
	}
	if !bytes.Equal(fContent, body) {
		panicString := fmt.Sprint("GetObject Error:", "byte comparison of uploaded != downloaded")
		panic(panicString)
	}

	// Duplication (Copy) an existing object.
	if err = objectstorage.CopyObject(sess, url+"/"+object, "/"+object+".dup"); err != nil {
		panicString := fmt.Sprint("CopyObject Error:", err)
		panic(panicString)
	}

	// Delete the objects.
	if err = objectstorage.DeleteObject(sess, url+"/"+object); err != nil {
		panicString := fmt.Sprint("DeleteObject Error:", err)
		panic(panicString)
	}
	if err = objectstorage.DeleteObject(sess, url+"/"+object+".dup"); err != nil {
		panicString := fmt.Sprint("DeleteObject Error:", err)
		panic(panicString)
	}

	// Delete the container that was previously created.
	if err = objectstorage.DeleteContainer(sess, url+"/"+config.Container); err != nil {
		panicString := fmt.Sprint("DeleteContainer Error:", err)
		panic(panicString)
	}
}
