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
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego/context"
	c "github.com/opensds/opensds/pkg/context"

	"github.com/astaxie/beego"
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

var (
	ByteHostReq = []byte(`
        {
            "accessMode": "agentless",
            "hostName": "sap1",
            "ip": "192.168.56.12",
            "availabilityZones": [
                "az1",
                "az2"
            ],
            "initiators": [
                {
                    "portName": "20000024ff5bb888",
                    "protocol": "iscsi"
                },
                {
                    "portName": "20000024ff5bc999",
                    "protocol": "iscsi"
                }
            ]
        }`)

	hostReq = model.HostSpec{
		BaseModel:         &model.BaseModel{},
		AccessMode:        "agentless",
		HostName:          "sap1",
		IP:                "192.168.56.12",
		AvailabilityZones: []string{"az1", "az2"},
		Initiators: []*model.Initiator{
			&model.Initiator{
				PortName: "20000024ff5bb888",
				Protocol: "iscsi",
			},
			&model.Initiator{
				PortName: "20000024ff5bc999",
				Protocol: "iscsi",
			},
		},
	}
)

func TestCreateHost(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		fakeHost := &SampleHosts[0]

		mockClient := new(dbtest.Client)
		mockClient.On("CreateHost", c.NewAdminContext(), &hostReq).Return(fakeHost, nil)
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/host/hosts", bytes.NewBuffer(ByteHostReq))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		var output model.HostSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, fakeHost)

	})
}

func TestListHosts(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		fakeHosts := []*model.HostSpec{&SampleHosts[0], &SampleHosts[1]}
		mockClient := new(dbtest.Client)
		mockClient.On("ListHosts", c.NewAdminContext()).Return(fakeHosts, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/host/hosts", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output []*model.HostSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, fakeHosts)
	})
}

func TestGetHost(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		fakeHost := &SampleHosts[0]
		mockClient := new(dbtest.Client)
		mockClient.On("GetHost", c.NewAdminContext(), fakeHost.Id).Return(fakeHost, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/host/hosts/"+SampleHosts[0].Id, nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.HostSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, fakeHost)
	})
}

func TestUpdateHost(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		fakeHost := &SampleHosts[0]

		var fakeHostUpdateReq model.HostSpec
		tmp, _ := json.Marshal(&hostReq)
		json.Unmarshal(tmp, &fakeHostUpdateReq)
		fakeHostUpdateReq.Id = fakeHost.Id

		mockClient := new(dbtest.Client)
		mockClient.On("UpdateHost", c.NewAdminContext(), &fakeHostUpdateReq).Return(fakeHost, nil)
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/host/hosts/"+fakeHost.Id, bytes.NewBuffer(ByteHostReq))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		var output model.HostSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, fakeHost)

	})
}

func TestDeleteHost(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		fakeHost := &SampleHosts[0]
		mockClient := new(dbtest.Client)
		mockClient.On("DeleteHost", c.NewAdminContext(), fakeHost.Id).Return(nil)
		db.C = mockClient

		r, _ := http.NewRequest("DELETE", "/v1beta/host/hosts/"+fakeHost.Id, nil)
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		assertTestResult(t, w.Code, 200)

	})
}
