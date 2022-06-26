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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego/context"
	c "github.com/sodafoundation/api/pkg/context"

	"github.com/astaxie/beego"
	"github.com/sodafoundation/api/pkg/db"
	"github.com/sodafoundation/api/pkg/model"
	. "github.com/sodafoundation/api/testutils/collection"
	dbtest "github.com/sodafoundation/api/testutils/db/testing"
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
                "default",
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
		AvailabilityZones: []string{"default", "az2"},
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
		mockClient.On("ListHostsByName", c.NewAdminContext(), hostReq.HostName).Return(nil, nil)
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
	t.Run("Should return 400 for invalid request - lists not found db error", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("ListHostsByName", c.NewAdminContext(), hostReq.HostName).Return(nil, errors.New("List hosts failed: "))
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/host/hosts", bytes.NewBuffer(ByteHostReq))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		assertTestResult(t, w.Code, 400)

	})
	t.Run("Should return 400 - for some db error while creating host", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("ListHostsByName", c.NewAdminContext(), hostReq.HostName).Return(nil, nil)
		mockClient.On("CreateHost", c.NewAdminContext(), &hostReq).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/host/hosts", bytes.NewBuffer(ByteHostReq))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		assertTestResult(t, w.Code, 400)
	})
	t.Run("Should return 400 - the host with same name already exists in the system", func(t *testing.T) {
		fakeHost := []*model.HostSpec{&SampleHosts[0], &SampleHosts[1]}

		mockClient := new(dbtest.Client)
		mockClient.On("ListHostsByName", c.NewAdminContext(), hostReq.HostName).Return(fakeHost, nil)
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/host/hosts", bytes.NewBuffer(ByteHostReq))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		assertTestResult(t, w.Code, 400)

	})
}

func TestListHosts(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		fakeHosts := []*model.HostSpec{&SampleHosts[0], &SampleHosts[1]}
		mockClient := new(dbtest.Client)
		mockClient.On("ListHosts", c.NewAdminContext(), map[string][]string{}).Return(fakeHosts, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/host/hosts", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output []*model.HostSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, fakeHosts)
	})

	t.Run("Should return 400 - when lists fails with db error", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("ListHosts", c.NewAdminContext(), map[string][]string{}).Return(nil, errors.New("When list hosts in db:"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/host/hosts", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 400)
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
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.HostSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, fakeHost)
	})

	t.Run("Should return 404 - specified host id doesn't exists in db", func(t *testing.T) {
		fakeHost := &SampleHosts[0]
		mockClient := new(dbtest.Client)
		mockClient.On("GetHost", c.NewAdminContext(), fakeHost.Id).Return(nil, errors.New("specified host(202964b5-8e73-46fd-b41b-a8e403f3c30b) can't find"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/host/hosts/"+SampleHosts[0].Id, nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
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

	t.Run("Should return 400 - update host failed with db error", func(t *testing.T) {
		fakeHost := &SampleHosts[0]
		var fakeHostUpdateReq model.HostSpec
		tmp, _ := json.Marshal(&hostReq)
		json.Unmarshal(tmp, &fakeHostUpdateReq)
		fakeHostUpdateReq.Id = fakeHost.Id

		mockClient := new(dbtest.Client)
		mockClient.On("UpdateHost", c.NewAdminContext(), &fakeHostUpdateReq).Return(nil, errors.New("update host failed:"))
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/host/hosts/"+fakeHost.Id, bytes.NewBuffer(ByteHostReq))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/JSON")
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		assertTestResult(t, w.Code, 400)
	})
}

func TestDeleteHost(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		fakeHost := &SampleHosts[0]
		mockClient := new(dbtest.Client)
		mockClient.On("DeleteHost", c.NewAdminContext(), fakeHost.Id).Return(nil)
		mockClient.On("GetHost", c.NewAdminContext(), fakeHost.Id).Return(fakeHost, nil)
		mockClient.On("ListVolumeAttachmentsWithFilter", c.NewAdminContext(),
			map[string][]string{"hostId": []string{fakeHost.Id}}).Return(nil, nil)
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
