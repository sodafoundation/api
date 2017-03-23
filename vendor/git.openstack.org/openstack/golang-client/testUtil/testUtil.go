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

// Package testUtil has helpers to be used with unit tests
package testUtil

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// Equals fails the test if exp is not equal to act.
// Code was copied from https://github.com/benbjohnson/testing MIT license
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// Assert fails the test if the condition is false.
// Code was copied from https://github.com/benbjohnson/testing MIT license
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// IsNil ensures that the act interface is nil
// otherwise an error is raised.
func IsNil(tb testing.TB, act interface{}) {
	if act != nil {
		tb.Error("expected nil", act)
		tb.FailNow()
	}
}

// CreateGetJSONTestServer creates a httptest.Server that can be used to test
// JSON Get requests. Takes a token, JSON payload, and a verification function
// to do additional validation
func CreateGetJsonTestServer(
	t *testing.T,
	expectedAuthToken string,
	jsonResponsePayload string,
	verifyRequest func(*http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			headerValuesEqual(t, r, "X-Auth-Token", expectedAuthToken)
			headerValuesEqual(t, r, "Accept", "application/json")
			// verifyRequest(r)
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(jsonResponsePayload))
				w.WriteHeader(http.StatusOK)
				return
			}

			t.Error(errors.New("Failed: r.Method == GET"))
		}))
}

// CreateGetJSONTestRequestServer creates a httptest.Server that can be used to test GetJson requests. Just specify the token,
// json payload that is to be read by the response, and a verification func that can be used
// to do additional validation of the request that is built
func CreateGetJSONTestRequestServer(t *testing.T, expectedAuthTokenValue string, jsonResponsePayload string, verifyRequest func(*http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			headerValuesEqual(t, r, "X-Auth-Token", expectedAuthTokenValue)
			headerValuesEqual(t, r, "Accept", "application/json")
			verifyRequest(r)
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(jsonResponsePayload))
				w.WriteHeader(http.StatusOK)
				return
			}

			t.Error(errors.New("Failed: r.Method == GET"))
		}))
}

// CreatePostJSONTestRequestServer creates a http.Server that can be used to test PostJson requests. Specify the token,
// response json payload and the url and request body that is expected.
func CreatePostJSONTestRequestServer(t *testing.T, expectedAuthTokenValue string, outputResponseJSONPayload string, expectedRequestURLEndsWith string, expectedRequestBody string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			headerValuesEqual(t, r, "X-Auth-Token", expectedAuthTokenValue)
			headerValuesEqual(t, r, "Accept", "application/json")
			headerValuesEqual(t, r, "Content-Type", "application/json")
			reqURL := r.URL.String()
			if !strings.HasSuffix(reqURL, expectedRequestURLEndsWith) {
				t.Error(errors.New("Incorrect url created, expected:" + expectedRequestURLEndsWith + " at the end, actual url:" + reqURL))
			}
			actualRequestBody := dumpRequestBody(r)
			if actualRequestBody != expectedRequestBody {
				t.Error(errors.New("Incorrect payload created, expected:'" + expectedRequestBody + "', actual '" + actualRequestBody + "'"))
			}
			if r.Method == "POST" {
				w.Header().Set("Content-Type", "application/json")
				// Status Code has to be written before writing the payload or
				// else it defaults to 200 OK instead.
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(outputResponseJSONPayload))
				return
			}

			t.Error(errors.New("Failed: r.Method == POST"))
		}))
}

// CreateDeleteTestRequestServer creates a http.Server that can be used to test Delete requests.
func CreateDeleteTestRequestServer(t *testing.T, expectedAuthTokenValue string, urlEndsWith string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			headerValuesEqual(t, r, "X-Auth-Token", expectedAuthTokenValue)
			reqURL := r.URL.String()
			if !strings.HasSuffix(reqURL, urlEndsWith) {
				t.Error(errors.New("Incorrect url created, expected '" + urlEndsWith + "' at the end, actual url:" + reqURL))
			}
			if r.Method == "DELETE" {
				w.WriteHeader(204)
				return
			}

			t.Error(errors.New("Failed: r.Method == DELETE"))
		}))
}

func dumpRequestBody(request *http.Request) (body string) {
	requestWithBody, _ := httputil.DumpRequest(request, true)
	requestWithoutBody, _ := httputil.DumpRequest(request, false)
	body = strings.Replace(string(requestWithBody), string(requestWithoutBody), "", 1)
	return
}

func headerValuesEqual(t *testing.T, req *http.Request, name string, expectedValue string) {
	actualValue := req.Header.Get(name)
	if actualValue != expectedValue {
		t.Error(fmt.Errorf("Expected Header {Name:'%s', Value:'%s', actual value '%s'", name, expectedValue, actualValue))
	}
}
