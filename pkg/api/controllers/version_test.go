// Copyright 2019 The OpenSDS Authors.
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

package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	c "github.com/opensds/opensds/pkg/context"
)

func init() {
	var versionPortal VersionPortal
	beego.Router("/", &versionPortal, "get:ListVersions")
	beego.Router("/:apiVersion", &versionPortal, "get:GetVersion")
}

func TestListVersions(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		var output []map[string]string
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, KnownVersions)
	})
}

func TestGetVersion(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "/v1beta", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		var output map[string]string
		json.Unmarshal(w.Body.Bytes(), &output)
		var expected = map[string]string{
			"name":        "v1beta",
			"description": "v1beta version",
			"status":      "CURRENT",
			"updatedAt":   "2017-07-10T14:36:58.014Z",
		}
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, expected)
	})

	t.Run("Should return 404 if get version with invalid API version", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "/InvalidAPIVersion", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}
