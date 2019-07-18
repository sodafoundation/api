// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package testhelper

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	// Mux is a multiplexer that can be used to register handlers.
	Mux *http.ServeMux

	// Server is an in-memory HTTP server for testing.
	Server *httptest.Server
)

// SetupHTTP prepares the Mux and Server.
func SetupHTTP() {
	Mux = http.NewServeMux()
	Server = httptest.NewServer(Mux)
}

// TeardownHTTP releases HTTP-related resources.
func TeardownHTTP() {
	Server.Close()
}

// Endpoint returns a fake endpoint that will actually target the Mux.
func Endpoint() string {
	return Server.URL + "/"
}

// TestMethod checks that the Request has the expected method (e.g. GET, POST).
func TestMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

// TestHeader checks that the header on the http.Request matches the expected value.
func TestHeader(t *testing.T, r *http.Request, header string, expected string) {
	if actual := r.Header.Get(header); expected != actual {
		t.Errorf("Header %s = %s, expected %s", header, actual, expected)
	}
}

// TestBody verifies that the request body matches an expected body.
func TestBody(t *testing.T, r *http.Request, expected string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Unable to read body: %v", err)
	}
	str := string(b)
	if expected != str {
		t.Errorf("Body = %s, expected %s", str, expected)
	}
}
