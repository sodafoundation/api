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
	"errors"
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
	beego.Router("/v1beta/block/volumeGroups", &VolumeGroupPortal{}, "post:CreateVolumeGroup;get:ListVolumeGroups")
	beego.Router("/v1beta/block/volumeGroups/:groupId", &VolumeGroupPortal{}, "put:UpdateVolumeGroup;get:GetVolumeGroup;delete:DeleteVolumeGroup")
}

func TestListVolumeGroups(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		var sampleVGs = []*model.VolumeGroupSpec{&SampleVolumeGroups[0]}
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListVolumeGroupsWithFilter", c.NewAdminContext(), m).Return(sampleVGs, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/volumeGroups?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output []*model.VolumeGroupSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, sampleVGs)
	})

	t.Run("Should return 500 if list volume groups with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListVolumeGroupsWithFilter", c.NewAdminContext(), m).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/volumeGroups?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestGetVolumeGroup(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolumeGroup", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f555").Return(&SampleVolumeGroups[0], nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/volumeGroups/3769855c-a102-11e7-b772-17b880d2f555", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.VolumeGroupSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SampleVolumeGroups[0])
	})

	t.Run("Should return 404 if get volume group with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetVolumeGroup", c.NewAdminContext(), "3769855c-a102-11e7-b772-17b880d2f555").Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/block/volumeGroups/3769855c-a102-11e7-b772-17b880d2f555", nil)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}
