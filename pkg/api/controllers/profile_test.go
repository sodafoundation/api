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
	"strings"
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/testutils/collection"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
)

func init() {
	var profilePortal ProfilePortal
	beego.Router("/v1beta/profiles", &profilePortal, "post:CreateProfile;get:ListProfiles")
	beego.Router("/v1beta/profiles/:profileId", &profilePortal, "get:GetProfile;put:UpdateProfile;delete:DeleteProfile")
	beego.Router("/v1beta/profiles/:profileId/customProperties", &profilePortal, "post:AddCustomProperty;get:ListCustomProperties")
	beego.Router("/v1beta/profiles/:profileId/customProperties/:customKey", &profilePortal, "delete:RemoveCustomProperty")
}

////////////////////////////////////////////////////////////////////////////////
//                            Tests for profile                               //
////////////////////////////////////////////////////////////////////////////////

func TestCreateProfile(t *testing.T) {
	var fakeBody = `{
		"name": "default",
		"description": "default policy",
		"storageType": "block",
		"customProperties": {
			"dataStorage": {
				"provisioningPolicy": "Thin",
				"isSpaceEfficient": true
			},
			"ioConnectivity": {
				"accessProtocol": "rbd",
				"maxIOPS": 5000000,
				"maxBWS": 500
			}
		}
	}`

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("CreateProfile", c.NewAdminContext(), &model.ProfileSpec{
			BaseModel:   &model.BaseModel{},
			Name:        "default",
			Description: "default policy",
			StorageType: "block",
			CustomProperties: model.CustomPropertiesSpec{
				"dataStorage": map[string]interface{}{
					"provisioningPolicy": "Thin",
					"isSpaceEfficient":   true,
				},
				"ioConnectivity": map[string]interface{}{
					"accessProtocol": "rbd",
					"maxIOPS":        float64(5000000),
					"maxBWS":         float64(500),
				},
			}}).Return(&SampleProfiles[1], nil)
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/profiles", strings.NewReader(fakeBody))
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.ProfileSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SampleProfiles[1])
	})
}

func TestUpdateProfile(t *testing.T) {
	var jsonStr = []byte(`{
		"id": "2f9c0a04-66ef-11e7-ade2-43158893e017",
		"name": "silver",
		"description": "silver policy"
	}`)
	var expectedJson = []byte(`{
		"id": "2f9c0a04-66ef-11e7-ade2-43158893e017",
		"name": "silver",
		"description": "silver policy",
		"customProperties": {
			"dataStorage": {
				"provisioningPolicy": "Thin",
				"isSpaceEfficient":   true
			},
			"ioConnectivity": {
				"accessProtocol": "rbd",
				"maxIOPS":        5000000,
				"maxBWS":         500
			}
		}
	}`)
	var expected model.ProfileSpec
	json.Unmarshal(expectedJson, &expected)

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		profile := model.ProfileSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&profile)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateProfile", c.NewAdminContext(), profile.Id, &profile).
			Return(&expected, nil)
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.ProfileSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &expected)
	})

	t.Run("Should return 500 if update profile with bad request", func(t *testing.T) {
		profile := model.ProfileSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&profile)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateProfile", c.NewAdminContext(), profile.Id, &profile).
			Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestListProfiles(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		var sampleProfiles = []*model.ProfileSpec{&SampleProfiles[1]}
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListProfilesWithFilter", c.NewAdminContext(), m).Return(
			sampleProfiles, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/profiles?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output []*model.ProfileSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, sampleProfiles)
	})

	t.Run("Should return 500 if list profiles with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListProfilesWithFilter", c.NewAdminContext(), m).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/profiles?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestGetProfile(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetProfile", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").
			Return(&SampleProfiles[1], nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.ProfileSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SampleProfiles[1])
	})

	t.Run("Should return 404 if get profile with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetProfile", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(
			nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET",
			"/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}

func TestDeleteProfile(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetProfile", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(
			&SampleProfiles[1], nil)
		mockClient.On("ListVolumesByProfileId", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(
			SampleVolumeNames, nil)
		mockClient.On("DeleteProfile", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(nil)
		db.C = mockClient

		r, _ := http.NewRequest("DELETE",
			"/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 200)
	})

	t.Run("Should return 404 if delete profile with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetProfile", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(
			nil, errors.New("Invalid resource uuid"))
		mockClient.On("ListVolumesByProfileId", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(
			nil, errors.New("Depency volumes"))
		db.C = mockClient

		r, _ := http.NewRequest("DELETE",
			"/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}

////////////////////////////////////////////////////////////////////////////////
//               Tests for file share profile                                 //
////////////////////////////////////////////////////////////////////////////////

func TestFileShareCreateProfile(t *testing.T) {
	var fakeBody = `{
                "name": "silver",
                "description": "silver policy",
                "storageType": "file",
                "provisioningProperties":{
                        "dataStorage":{
                                "storageAccessCapability": ["Read","Write","Execute"],
                                "provisioningPolicy": "Thin",
                                "isSpaceEfficient": true
                        },
                        "ioConnectivity": {
                                "accessProtocol": "NFS",
                                "maxIOPS": 5000000,
                                "maxBWS": 500
                        }
                }
        }`

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("CreateProfile", c.NewAdminContext(), &model.ProfileSpec{
			BaseModel:   &model.BaseModel{},
			Name:        "silver",
			Description: "silver policy",
			StorageType: "file",
			ProvisioningProperties: model.ProvisioningPropertiesSpec{
				DataStorage: model.DataStorageLoS{
					StorageAccessCapability: []string{"Read", "Write", "Execute"},
					ProvisioningPolicy:      "Thin",
					IsSpaceEfficient:        true,
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "NFS",
					MaxIOPS:        5000000,
					MaxBWS:         500,
				},
			}}).Return(&SampleFileShareProfiles[1], nil)
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/profiles", strings.NewReader(fakeBody))
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.ProfileSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SampleFileShareProfiles[1])
	})
}
func TestFileShareUpdateProfile(t *testing.T) {
	var jsonStr = []byte(`{
                "id": "2f9c0a04-66ef-11e7-ade2-43158893e017",
                "name": "silver",
                "description": "silver policy"
        }`)
	var expectedJson = []byte(`{
                "id": "2f9c0a04-66ef-11e7-ade2-43158893e017",
                "name": "silver",
                "description": "silver policy",
                "storageType": "file",
                "provisioningProperties":{
                        "dataStorage":{
                                "storageAccessCapability": ["Read","Write","Execute"],
                                "provisioningPolicy": "Thin",
                                "isSpaceEfficient": true
                        },
                        "ioConnectivity": {
                                "accessProtocol": "NFS",
                                "maxIOPS":        5000000,
                                "maxBWS":         500
                        }
                }
        }`)
	var expected model.ProfileSpec
	json.Unmarshal(expectedJson, &expected)

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		profile := model.ProfileSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&profile)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateProfile", c.NewAdminContext(), profile.Id, &profile).
			Return(&expected, nil)
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.ProfileSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &expected)
	})

	t.Run("Should return 500 if update profile with bad request", func(t *testing.T) {
		profile := model.ProfileSpec{BaseModel: &model.BaseModel{}}
		json.NewDecoder(bytes.NewBuffer(jsonStr)).Decode(&profile)
		mockClient := new(dbtest.Client)
		mockClient.On("UpdateProfile", c.NewAdminContext(), profile.Id, &profile).
			Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("PUT", "/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017", bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestListFileShareProfiles(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		var sampleProfiles = []*model.ProfileSpec{&SampleFileShareProfiles[1]}
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListProfilesWithFilter", c.NewAdminContext(), m).Return(
			sampleProfiles, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/profiles?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output []*model.ProfileSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, output, sampleProfiles)
	})

	t.Run("Should return 500 if list profiles with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		m := map[string][]string{
			"offset":  {"0"},
			"limit":   {"1"},
			"sortDir": {"asc"},
			"sortKey": {"name"},
		}
		mockClient.On("ListProfilesWithFilter", c.NewAdminContext(), m).Return(nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/profiles?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestGetFileShareProfile(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetProfile", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").
			Return(&SampleFileShareProfiles[1], nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.ProfileSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SampleFileShareProfiles[1])
	})

	t.Run("Should return 404 if get profile with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetProfile", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(
			nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET",
			"/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}

func TestDeleteFileShareProfile(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetProfile", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(
			&SampleFileShareProfiles[1], nil)
		mockClient.On("ListFileSharesByProfileId", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(
			SampleShareNames, nil)
		mockClient.On("DeleteProfile", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(nil)
		db.C = mockClient

		r, _ := http.NewRequest("DELETE",
			"/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 200)
	})

	t.Run("Should return 404 if delete profile with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("GetProfile", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(
			nil, errors.New("Invalid resource uuid"))
		mockClient.On("ListFileSharesByProfileId", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(
			nil, errors.New("Depency FileShares"))
		db.C = mockClient

		r, _ := http.NewRequest("DELETE",
			"/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 404)
	})
}

////////////////////////////////////////////////////////////////////////////////
//               Tests for profile custom properties spec                     //
////////////////////////////////////////////////////////////////////////////////

func TestAddCustomProperty(t *testing.T) {
	var fakeBody = `{
		"dataStorage": {
			"provisioningPolicy": "Thin",
			"isSpaceEfficient": true
		}
	}`

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("AddCustomProperty", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017", model.CustomPropertiesSpec{
			"dataStorage": map[string]interface{}{
				"provisioningPolicy": "Thin",
				"isSpaceEfficient":   true}}).Return(&SampleCustomProperties, nil)
		db.C = mockClient

		r, _ := http.NewRequest("POST", "/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017/customProperties", strings.NewReader(fakeBody))
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.CustomPropertiesSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SampleCustomProperties)
	})
}

func TestListCustomProperties(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("ListCustomProperties", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(
			&SampleCustomProperties, nil)
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017/customProperties", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		var output model.CustomPropertiesSpec
		json.Unmarshal(w.Body.Bytes(), &output)
		assertTestResult(t, w.Code, 200)
		assertTestResult(t, &output, &SampleCustomProperties)
	})

	t.Run("Should return 500 if list custom properties with bad request", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("ListCustomProperties", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(
			nil, errors.New("db error"))
		db.C = mockClient

		r, _ := http.NewRequest("GET", "/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017/customProperties", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		assertTestResult(t, w.Code, 500)
	})
}

func TestRemoveCustomProperty(t *testing.T) {

	t.Run("Should return 200 if everything works well", func(t *testing.T) {
		mockClient := new(dbtest.Client)
		mockClient.On("RemoveCustomProperty", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017", "key1").Return(nil)
		db.C = mockClient

		r, _ := http.NewRequest("DELETE",
			"/v1beta/profiles/2f9c0a04-66ef-11e7-ade2-43158893e017/customProperties/key1", nil)
		w := httptest.NewRecorder()
		beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
			httpCtx.Input.SetData("context", c.NewAdminContext())
		})
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		assertTestResult(t, w.Code, 200)
	})
}
