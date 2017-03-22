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

package objectstorage

import (
	"errors"
	"git.openstack.org/openstack/golang-client/openstack"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var znHome = "./"
var objFile = "objectstorage_test.go"
var srcFile = znHome + objFile
var tokn = "eaaafd18-0fed-4b3a-81b4-663c99ec1cbb"
var containerName = "John's container"
var containerPrefix = "/" + containerName
var objPrefix = containerPrefix + "/" + objFile

func TestGetAccountMeta(t *testing.T) {
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "HEAD" {
				w.Header().Set("X-Account-Container-Count", "7")
				w.Header().Set("X-Account-Object-Count", "413")
				w.Header().Set("X-Account-Bytes-Used", "987654321000")
				w.WriteHeader(204)
				return
			}
			t.Error(errors.New("Failed: r.Method == HEAD"))
		}))
	defer apiServer.Close()

	sess := testGetSession()
	meta, err := GetAccountMeta(sess, apiServer.URL)
	if err != nil {
		t.Error(err)
	}
	if meta.Get("X-Account-Container-Count") != "7" ||
		meta.Get("X-Account-Object-Count") != "413" ||
		meta.Get("X-Account-Bytes-Used") != "987654321000" {
		t.Error("Failed: meta not matching")
	}
}

func TestListContainers(t *testing.T) {
	var containerList = `[
	{"name":"container 1",
	"count":2, "bytes":78},
	{"name":"container 2",
	"count":1,
	"bytes":17}]`
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(200)
				w.Write([]byte(containerList))
				return
			}
			t.Error(errors.New("Failed: r.Method == GET"))
		}))
	defer apiServer.Close()

	sess := testGetSession()
	myList, err := ListContainers(sess, 0, "", apiServer.URL)
	if err != nil {
		t.Error(err)
	}
	if string(myList) != containerList {
		t.Error(errors.New("Failed: input != output"))
	}
}

func TestListObjects(t *testing.T) {
	var objList = `[
		{"name":"test obj 1", 
		"hash":"4281c348eaf83e70ddce0e07221c3d28",
		"bytes":14,
		"content_type":"application\/octet-stream",
		"last_modified":"2009-02-03T05:26:32.612278"},
		{"name":"test obj 2",
		"hash":"b039efe731ad111bc1b0ef221c3849d0",
		"bytes":64,
		"content_type":"application\/octet-stream",
		"last_modified":"2009-02-03T05:26:32.612278"}
		]`
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(200)
				w.Write([]byte(objList))
				return
			}
			t.Error(errors.New("Failed: r.Method == GET"))
		}))
	defer apiServer.Close()

	sess := testGetSession()
	myList, err := ListObjects(
		sess, 0, "", "", "", "", apiServer.URL+containerPrefix)
	if err != nil {
		t.Error(err)
	}
	if string(myList) != objList {
		t.Error(errors.New("Failed: input != output"))
	}
}

func TestDeleteContainer(t *testing.T) {
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "DELETE" {
				w.WriteHeader(204)
				return
			}
			t.Error(errors.New("Failed: r.Method == DELETE"))
		}))
	defer apiServer.Close()

	sess := testGetSession()
	if err := DeleteContainer(sess, apiServer.URL+containerPrefix); err != nil {
		t.Error(err)
	}
}

func TestGetContainerMeta(t *testing.T) {
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "HEAD" {
				w.Header().Set("X-Container-Object-Count", "7")
				w.Header().Set("X-Container-Bytes-Used", "413")
				w.Header().Set("X-Container-Meta-InspectedBy", "Jack Wolf")
				w.WriteHeader(204)
				return
			}
			t.Error(errors.New("Failed: r.Method == HEAD"))
		}))
	defer apiServer.Close()

	sess := testGetSession()
	meta, err := GetContainerMeta(sess, apiServer.URL+containerPrefix)
	if err != nil {
		t.Error(err)
	}
	if meta.Get("X-Container-Object-Count") != "7" ||
		meta.Get("X-Container-Bytes-Used") != "413" ||
		meta.Get("X-Container-Meta-InspectedBy") != "Jack Wolf" {
		t.Error("Failed: meta not matching")
	}
}

func TestSetContainerMeta(t *testing.T) {
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" && r.Header.Get("X-Container-Meta-Fruit") == "Apple" {
				w.WriteHeader(204)
				return
			}
			t.Error(errors.New(
				"Failed: r.Method == POST && X-Container-Meta-Fruit == Apple"))
		}))
	defer apiServer.Close()

	sess := testGetSession()
	headers := http.Header{}
	headers.Add("X-Container-Meta-Fruit", "Apple")
	if err := SetContainerMeta(
		sess, apiServer.URL+containerPrefix,
		headers); err != nil {
		t.Error(err)
	}
}

func TestPutContainer(t *testing.T) {
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "PUT" {
				w.WriteHeader(201)
				return
			}
			t.Error(errors.New("Failed: r.Method == PUT"))
		}))
	defer apiServer.Close()

	sess := testGetSession()
	headers := http.Header{}
	headers.Add("X-TTL", "259200")
	headers.Add("X-Log-Retention", "true")
	if err := PutContainer(sess, apiServer.URL+containerPrefix,
		headers); err != nil {
		t.Error(err)
	}
}

func TestPutObject(t *testing.T) {
	var fContent []byte
	f, err := os.Open(srcFile)
	defer f.Close()
	if err != nil {
		t.Error(err)
	}
	fContent, err = ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	f.Close()
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			rBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
			}
			if r.Method == "PUT" && len(fContent) == len(rBody) {
				w.WriteHeader(201)
				return
			}
			t.Error(errors.New("Failed: Not 201"))
		}))
	defer apiServer.Close()

	sess := testGetSession()
	headers := http.Header{}
	if err = PutObject(sess, &fContent, apiServer.URL+objPrefix, headers); err != nil {
		t.Error(err)
	}
}

func TestCopyObject(t *testing.T) {
	destURL := "/destContainer/dest/Obj"
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "COPY" && r.Header.Get("Destination") == destURL {
				w.WriteHeader(200)
				return
			}
			t.Error(errors.New(
				"Failed: r.Method == COPY && r.Header.Get(Destination) == destURL"))
		}))
	defer apiServer.Close()

	sess := testGetSession()
	if err := CopyObject(sess, apiServer.URL+objPrefix, destURL); err != nil {
		t.Error(err)
	}
}

func TestGetObjectMeta(t *testing.T) {
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "HEAD" {
				w.Header().Set("X-Object-Meta-Fruit", "Apple")
				w.Header().Set("X-Object-Meta-Veggie", "Carrot")
				w.WriteHeader(200)
				return
			}
			t.Error(errors.New(
				"Failed: r.Method == HEAD && r.Header.Get(X-Auth-Token) == tokn"))
		}))
	defer apiServer.Close()

	sess := testGetSession()
	meta, err := GetObjectMeta(sess, apiServer.URL+objPrefix)
	if err != nil {
		t.Error(err)
	}
	if meta.Get("X-Object-Meta-Fruit") != "Apple" ||
		meta.Get("X-Object-Meta-Veggie") != "Carrot" {
		t.Error("Failed: meta not matching")
	}
}

func TestSetObjectMeta(t *testing.T) {
	var apiServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		if r.Method == "POST" && r.Header.Get("X-Object-Meta-Fruit") == "Apple" {
			w.WriteHeader(202)
			return
		}
		t.Error(errors.New("Failed: r.Method == POST && X-Object-Meta-Fruit == Apple"))
	}))
	defer apiServer.Close()

	sess := testGetSession()
	headers := http.Header{}
	headers.Add("X-Object-Meta-Fruit", "Apple")
	if err := SetObjectMeta(sess, apiServer.URL+objPrefix,
		headers); err != nil {
		t.Error(err)
	}
}

func TestGetObject(t *testing.T) {
	var unCompressedLen int
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				fContent, err := ioutil.ReadFile(srcFile)
				if err != nil {
					t.Error(err)
				}
				unCompressedLen = len(fContent)
				w.Header().Set("Content-Length", strconv.Itoa(unCompressedLen))
				w.Header().Set("X-Object-ModTime", "93000299")
				w.Header().Set("X-Object-Mode", "rwxrwxrwx")
				w.Write(fContent)
				return
			}
			t.Error(errors.New("Failed: r.Method == GET"))
		}))
	defer apiServer.Close()

	sess := testGetSession()
	hdr, body, err := GetObject(sess, apiServer.URL+objPrefix)
	if err != nil {
		t.Error(err)
	}
	if unCompressedLen != len(body) {
		t.Error(errors.New("GET: incorrect uncompressed len"))
	}
	if hdr.Get("X-Object-ModTime") != "93000299" ||
		hdr.Get("Content-Length") != strconv.Itoa(len(body)) ||
		hdr.Get("X-Object-Mode") != "rwxrwxrwx" {
		//
		t.Error(errors.New("GET: incorrect hdr"))
	}
}

func TestDeleteObject(t *testing.T) {
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "DELETE" {
				w.WriteHeader(204)
				return
			}
			t.Error(errors.New("Failed: r.Method == DELETE"))
		}))
	defer apiServer.Close()

	sess := testGetSession()
	if err := DeleteObject(sess, apiServer.URL+objPrefix); err != nil {
		t.Error(err)
	}
}

func testGetSession() *openstack.Session {
	auth := openstack.AuthToken{
		Access: openstack.AccessType{
			Token: openstack.Token{
				ID: tokn,
			},
		},
	}
	sess, _ := openstack.NewSession(http.DefaultClient, auth, nil)
	return sess
}
