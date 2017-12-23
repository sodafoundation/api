// Copyright 2017 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/astaxie/beego"
)

func init() {
	var versionPortal VersionPortal
	beego.Router("/", &versionPortal, "get:ListVersions")
	beego.Router("/:apiVersion", &versionPortal, "get:GetVersion")
}

func TestListVersions(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output []map[string]string
	json.Unmarshal(w.Body.Bytes(), &output)

	if !reflect.DeepEqual(KnownVersions, output) {
		t.Errorf("Expected %v, actual %v", KnownVersions, output)
	}
}

func TestGetVersion(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1beta", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output map[string]string
	json.Unmarshal(w.Body.Bytes(), &output)

	var expected = map[string]string{
		"name":        "v1beta",
		"description": "v1beta version",
		"status":      "CURRENT",
		"updatedAt":   "2017-07-10T14:36:58.014Z",
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestGetVersionWithInvalidAPIVersion(t *testing.T) {
	r, _ := http.NewRequest("GET", "/InvalidAPIVersion", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("Expected 404, actual %v", w.Code)
	}
}
