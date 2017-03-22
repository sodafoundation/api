// session_test - REST client session tests
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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"git.openstack.org/openstack/golang-client/testUtil"
)

type TestStruct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestSessionGet(t *testing.T) {
	tokn := "eaaafd18-0fed-4b3a-81b4-663c99ec1cbb"
	var apiServer = testUtil.CreateGetJsonTestServer(
		t,
		tokn,
		`{"id":"id1","name":"Chris"}`,
		nil,
	)
	expected := TestStruct{ID: "id1", Name: "Chris"}
	actual := TestStruct{}

	s, _ := NewSession(nil, nil, nil)
	var headers http.Header = http.Header{}
	headers.Set("X-Auth-Token", tokn)
	headers.Set("Accept", "application/json")
	headers.Set("Etag", "md5hash-blahblah")
	resp, err := s.Get(apiServer.URL, nil, &headers)
	if err != nil {
		t.Error(err)
	}
	testUtil.IsNil(t, err)

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if err = json.Unmarshal(rbody, &actual); err != nil {
		t.Error(err)
	}

	testUtil.Equals(t, expected, actual)
}
