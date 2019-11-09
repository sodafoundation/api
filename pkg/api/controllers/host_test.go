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
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/testutils/collection"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
)

func init() {
	var hostPortal HostPortal
	beego.Router("/v1beta/host/hosts", &hostPortal, "get:ListHosts;post:CreateHost")
	beego.Router("/v1beta/host/hosts/:hostId", &hostPortal, "get:GetHost;put:UpdateHost;delete:DeleteHost")
}

func TestListHosts(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("ListHosts", c.NewAdminContext()).Return(SampleHosts, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/host/hosts", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output []*model.HostSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, SampleHosts)
	})
}

func TestGetHost(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetHost", c.NewAdminContext(), SampleHosts[0].Id).Return(SampleHosts[0], nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/host/hosts/"+SampleHosts[0].Id, nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.HostSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SampleHosts[0])
	})
}
