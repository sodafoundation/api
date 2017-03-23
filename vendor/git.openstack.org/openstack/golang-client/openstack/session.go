// session - REST client session
// Copyright 2015 Dean Troyer
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openstack

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"git.openstack.org/openstack/golang-client/util"
)

var Debug = new(bool)

type Response struct {
	Resp *http.Response
	Body []byte
}

// Generic callback to get a token from the auth plugin
type AuthFunc func(s *Session, opts interface{}) (AuthRef, error)

type Session struct {
	httpClient *http.Client
	AuthToken  AuthRef
	Headers    http.Header
}

func NewSession(hclient *http.Client, auth AuthRef, tls *tls.Config) (session *Session, err error) {
	if hclient == nil {
		// Only build a transport if we're also building the client
		tr := &http.Transport{
			TLSClientConfig:    tls,
			DisableCompression: true,
		}
		hclient = &http.Client{Transport: tr}
	}
	session = &Session{
		httpClient: hclient,
		AuthToken:  auth,
		Headers:    http.Header{},
	}
	return session, nil
}

func (s *Session) NewRequest(method, url string, headers *http.Header, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		req.Header = *headers
	}
	if s.AuthToken != nil {
		req.Header.Add("X-Auth-Token", s.AuthToken.GetToken())
	}
	return
}

func (s *Session) Do(req *http.Request) (*http.Response, error) {
	// Add session headers
	for k := range s.Headers {
		req.Header.Set(k, s.Headers.Get(k))
	}

	if *Debug {
		d, _ := httputil.DumpRequestOut(req, true)
		log.Printf(">>>>>>>>>> REQUEST:\n", string(d))
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if *Debug {
		dr, _ := httputil.DumpResponse(resp, true)
		log.Printf("<<<<<<<<<< RESULT:\n", string(dr))
	}

	return resp, nil
}

// Perform a simple get to an endpoint
func (s *Session) Request(
	method string,
	url string,
	params *url.Values,
	headers *http.Header,
	body *[]byte,
) (resp *http.Response, err error) {
	// add params to url here
	if params != nil {
		url = url + "?" + params.Encode()
	}

	// Get the body if one is present
	var buf io.Reader
	if body != nil {
		buf = bytes.NewReader(*body)
	}

	req, err := s.NewRequest(method, url, headers, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err = s.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Perform a simple get to an endpoint and unmarshall returned JSON
func (s *Session) RequestJSON(
	method string,
	url string,
	params *url.Values,
	headers *http.Header,
	body interface{},
	responseContainer interface{},
) (resp *http.Response, err error) {
	bodyjson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	if headers == nil {
		headers = &http.Header{}
		headers.Add("Content-Type", "application/json")
	}

	resp, err = s.Request(method, url, params, headers, &bodyjson)
	if err != nil {
		return nil, err
	}

	err = util.CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nil, err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading response body")
	}
	if err = json.Unmarshal(rbody, &responseContainer); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Session) Delete(
	url string,
	params *url.Values,
	headers *http.Header,
) (resp *http.Response, err error) {
	return s.Request("DELETE", url, params, headers, nil)
}

func (s *Session) Get(
	url string,
	params *url.Values,
	headers *http.Header,
) (resp *http.Response, err error) {
	return s.Request("GET", url, params, headers, nil)
}

func (s *Session) GetJSON(
	url string,
	params *url.Values,
	headers *http.Header,
	responseContainer interface{},
) (resp *http.Response, err error) {
	return s.RequestJSON("GET", url, params, headers, nil, responseContainer)
}

func (s *Session) Head(
	url string,
	params *url.Values,
	headers *http.Header,
) (resp *http.Response, err error) {
	return s.Request("HEAD", url, params, headers, nil)
}

func (s *Session) Post(
	url string,
	params *url.Values,
	headers *http.Header,
	body *[]byte,
) (resp *http.Response, err error) {
	return s.Request("POST", url, params, headers, body)
}

func (s *Session) PostJSON(
	url string,
	params *url.Values,
	headers *http.Header,
	body interface{},
	responseContainer interface{},
) (resp *http.Response, err error) {
	return s.RequestJSON("POST", url, params, headers, body, responseContainer)
}

func (s *Session) Put(
	url string,
	params *url.Values,
	headers *http.Header,
	body *[]byte,
) (resp *http.Response, err error) {
	return s.Request("PUT", url, params, headers, body)
}

// Delete sends a DELETE request.
func Delete(
	url string,
	params *url.Values,
	headers *http.Header,
) (resp *http.Response, err error) {
	s, _ := NewSession(nil, nil, nil)
	return s.Delete(url, params, headers)
}

// Get sends a GET request.
func Get(
	url string,
	params *url.Values,
	headers *http.Header,
) (resp *http.Response, err error) {
	s, _ := NewSession(nil, nil, nil)
	return s.Get(url, params, headers)
}

// GetJSON sends a GET request and unmarshalls returned JSON.
func GetJSON(
	url string,
	params *url.Values,
	headers *http.Header,
	responseContainer interface{},
) (resp *http.Response, err error) {
	s, _ := NewSession(nil, nil, nil)
	return s.RequestJSON("GET", url, params, headers, nil, responseContainer)
}

// Post sends a POST request.
func Post(
	url string,
	params *url.Values,
	headers *http.Header,
	body *[]byte,
) (resp *http.Response, err error) {
	s, _ := NewSession(nil, nil, nil)
	return s.Post(url, params, headers, body)
}

// PostJSON sends a POST request and unmarshalls returned JSON.
func PostJSON(
	url string,
	params *url.Values,
	headers *http.Header,
	body interface{},
	responseContainer interface{},
) (resp *http.Response, err error) {
	s, _ := NewSession(nil, nil, nil)
	return s.RequestJSON("POST", url, params, headers, body, responseContainer)
}

// Put sends a PUT request.
func Put(
	url string,
	params *url.Values,
	headers *http.Header,
	body *[]byte,
) (resp *http.Response, err error) {
	s, _ := NewSession(nil, nil, nil)
	return s.Put(url, params, headers, body)
}
