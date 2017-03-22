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

package util

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"git.openstack.org/openstack/golang-client/testUtil"
)

var token = "2350971-5716-8165"

func TestDelete(t *testing.T) {
	var apiServer = testUtil.CreateDeleteTestRequestServer(t, token, "/other")
	defer apiServer.Close()

	err := Delete(apiServer.URL+"/other", token, *http.DefaultClient)
	testUtil.IsNil(t, err)
}

func TestPostJsonWithValidResponse(t *testing.T) {
	var apiServer = testUtil.CreatePostJSONTestRequestServer(t, token, `{"id":"id1","name":"Chris"}`, "", `{"id":"id1","name":"name"}`)
	defer apiServer.Close()
	actual := TestStruct{}
	ti := TestStruct{ID: "id1", Name: "name"}

	err := PostJSON(apiServer.URL, token, *http.DefaultClient, ti, &actual)
	testUtil.IsNil(t, err)
	expected := TestStruct{ID: "id1", Name: "Chris"}

	testUtil.Equals(t, expected, actual)
}

func TestCallAPI(t *testing.T) {
	tokn := "eaaafd18-0fed-4b3a-81b4-663c99ec1cbb"
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Auth-Token") != tokn {
				t.Error(errors.New("Token failed"))
			}
			w.WriteHeader(200) //ok
		}))
	zeroByte := &([]byte{})
	if _, err := CallAPI("HEAD", apiServer.URL, zeroByte, "X-Auth-Token", tokn); err != nil {
		t.Error(err)
	}
	if _, err := CallAPI("DELETE", apiServer.URL, zeroByte, "X-Auth-Token", tokn); err != nil {
		t.Error(err)
	}
	if _, err := CallAPI("POST", apiServer.URL, zeroByte, "X-Auth-Token", tokn); err != nil {
		t.Error(err)
	}
}

func TestCallAPIGetContent(t *testing.T) {
	tokn := "eaaafd18-0fed-4b3a-81b4-663c99ec1cbb"
	fContent, err := ioutil.ReadFile("./util.go")
	if err != nil {
		t.Error(err)
	}
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
			}
			if r.Header.Get("X-Auth-Token") != tokn {
				t.Error(errors.New("Token failed"))
			}
			w.Header().Set("Content-Length", r.Header.Get("Content-Length"))
			w.Write(body)
		}))
	var resp *http.Response
	if resp, err = CallAPI("GET", apiServer.URL, &fContent, "X-Auth-Token", tokn,
		"Etag", "md5hash-blahblah"); err != nil {
		t.Error(err)
	}
	if strconv.Itoa(len(fContent)) != resp.Header.Get("Content-Length") {
		t.Error(errors.New("Failed: Content-Length comparison"))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(fContent, body) {
		t.Error(errors.New("Failed: Content body comparison"))
	}
}

func TestCallAPIPutContent(t *testing.T) {
	tokn := "eaaafd18-0fed-4b3a-81b4-663c99ec1cbb"
	fContent, err := ioutil.ReadFile("./util.go")
	if err != nil {
		t.Error(err)
	}
	var apiServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Auth-Token") != tokn {
				t.Error(errors.New("Token failed"))
			}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
			}
			if strconv.Itoa(len(fContent)) != r.Header.Get("Content-Length") {
				t.Error(errors.New("Failed: Content-Length comparison"))
			}
			if !bytes.Equal(fContent, body) {
				t.Error(errors.New("Failed: Content body comparison"))
			}
			w.WriteHeader(200)
		}))
	if _, err = CallAPI("PUT", apiServer.URL, &fContent, "X-Auth-Token", tokn); err != nil {
		t.Error(err)
	}
}

type TestStruct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
