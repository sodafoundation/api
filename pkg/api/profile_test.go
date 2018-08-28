// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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
	"errors"

	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
)

func init() {
	var profilePortal ProfilePortal
	beego.Router("/v1beta/profiles", &profilePortal, "post:CreateProfile;get:ListProfiles")
	beego.Router("/v1beta/profiles/:profileId", &profilePortal, "get:GetProfile;put:UpdateProfile;delete:DeleteProfile")
	beego.Router("/v1beta/profiles/:profileId/extras", &profilePortal, "post:AddExtraProperty;get:ListExtraProperties")
	beego.Router("/v1beta/profiles/:profileId/extras/:extraKey", &profilePortal, "delete:RemoveExtraProperty")
}

var (
	fakeExtras = model.ExtraSpec{
		"key1": "val1",
		"key2": "val2",
		"key3": map[string]interface{}{
			"subKey1": "subVal1",
			"subKey2": "subVal2",
		},
	}
	fakeProfile = &model.ProfileSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			CreatedAt: "2017-10-24T16:21:32",
		},
		Name:        "Gold",
		Description: "Gold service",
		Extras:      fakeExtras,
	}
	fakeProfiles = []*model.ProfileSpec{fakeProfile}
)

////////////////////////////////////////////////////////////////////////////////
//                            Tests for profile                               //
////////////////////////////////////////////////////////////////////////////////

func TestCreateProfile(t *testing.T) {
	var fakeBody = `{
			"name": "Gold",
			"description": "Gold service"
		}`

	mockClient := new(dbtest.Client)
	mockClient.On("CreateProfile", c.NewAdminContext(), &model.ProfileSpec{
		BaseModel:   &model.BaseModel{},
		Name:        "Gold",
		Description: "Gold service"}).Return(&model.ProfileSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			CreatedAt: "2017-10-24T16:21:32",
		},
		Name:        "Gold",
		Description: "Gold service"}, nil)
	db.C = mockClient

	r, _ := http.NewRequest("POST", "/v1beta/profiles", strings.NewReader(fakeBody))
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.ProfileSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `
		{
			"id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"name": "Gold",
			"description": "Gold service",
			"createdAt": "2017-10-24T16:21:32",
			"updatedAt": ""
		}`

	var expected model.ProfileSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestUpdateProfile(t *testing.T) {

	mockClient := new(dbtest.Client)
	mockClient.On("UpdateProfile", c.NewAdminContext(), "f4a5e666-c669-4c64-a2a1-8f9ecd560c78", fakeProfile).Return(fakeProfile, nil)
	db.C = mockClient

	var fakeBody = `
		{
			"id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"name": "Gold",
			"description": "Gold service",
			"createdAt": "2017-10-24T16:21:32",
			"updatedAt": "",
			"extras": {
				"key1": "val1",
				"key2": "val2",
				"key3": {
					"subKey1": "subVal1",
					"subKey2": "subVal2"
				}
			}	
		}`
	r, _ := http.NewRequest("PUT", "/v1beta/profiles/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", strings.NewReader(fakeBody))
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.ProfileSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `
		{
			"id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"name": "Gold",
			"description": "Gold service",
			"createdAt": "2017-10-24T16:21:32",
			"updatedAt": "",
			"extras": {
				"key1": "val1",
				"key2": "val2",
				"key3": {
					"subKey1": "subVal1",
					"subKey2": "subVal2"
				}
			}	
		}`

	var expected model.ProfileSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListProfiles(t *testing.T) {

	mockClient := new(dbtest.Client)
	m := map[string][]string{
		"offset":  []string{"0"},
		"limit":   []string{"1"},
		"sortDir": []string{"asc"},
		"sortKey": []string{"name"},
	}
	mockClient.On("ListProfilesWithFilter", c.NewAdminContext(), m).Return(fakeProfiles, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/profiles?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output []model.ProfileSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `[
		{
			"id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"name": "Gold",
			"description": "Gold service",
			"createdAt": "2017-10-24T16:21:32",
			"updatedAt": "",
			"extras": {
				"key1": "val1",
				"key2": "val2",
				"key3": {
					"subKey1": "subVal1",
					"subKey2": "subVal2"
				}
			}	
		}
	]`

	var expected []model.ProfileSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListProfilesWithBadRequest(t *testing.T) {

	mockClient := new(dbtest.Client)
	m := map[string][]string{
		"offset":  []string{"0"},
		"limit":   []string{"1"},
		"sortDir": []string{"asc"},
		"sortKey": []string{"name"},
	}
	mockClient.On("ListProfilesWithFilter", c.NewAdminContext(), m).Return(nil, errors.New("db error"))
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/profiles?offset=0&limit=1&sortDir=asc&sortKey=name", nil)
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}

func TestGetProfile(t *testing.T) {

	mockClient := new(dbtest.Client)
	mockClient.On("GetProfile", c.NewAdminContext(), "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(fakeProfile, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/profiles/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.ProfileSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `
		{
			"id": "f4a5e666-c669-4c64-a2a1-8f9ecd560c78",
			"name": "Gold",
			"description": "Gold service",
			"createdAt": "2017-10-24T16:21:32",
			"updatedAt": "",
			"extras": {
				"key1": "val1",
				"key2": "val2",
				"key3": {
					"subKey1": "subVal1",
					"subKey2": "subVal2"
				}
			}	
		}`

	var expected model.ProfileSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestGetProfileWithBadRequest(t *testing.T) {

	mockClient := new(dbtest.Client)
	mockClient.On("GetProfile", c.NewAdminContext(), "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(
		nil, errors.New("db error"))
	db.C = mockClient

	r, _ := http.NewRequest("GET",
		"/v1beta/profiles/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}

func TestDeleteProfile(t *testing.T) {

	mockClient := new(dbtest.Client)
	mockClient.On("GetProfile", c.NewAdminContext(), "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(
		fakeProfile, nil)
	mockClient.On("DeleteProfile", c.NewAdminContext(), "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(nil)
	db.C = mockClient

	r, _ := http.NewRequest("DELETE",
		"/v1beta/profiles/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}
}

func TestDeleteProfileWithBadrequest(t *testing.T) {

	mockClient := new(dbtest.Client)
	mockClient.On("GetProfile", c.NewAdminContext(), "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(
		nil, errors.New("Invalid resource uuid"))
	db.C = mockClient

	r, _ := http.NewRequest("DELETE",
		"/v1beta/profiles/f4a5e666-c669-4c64-a2a1-8f9ecd560c78", nil)
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}

////////////////////////////////////////////////////////////////////////////////
//                          Tests for profile spec                            //
////////////////////////////////////////////////////////////////////////////////

func TestListExtraProperties(t *testing.T) {

	mockClient := new(dbtest.Client)
	mockClient.On("ListExtraProperties", c.NewAdminContext(), "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(&fakeExtras, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/profiles/f4a5e666-c669-4c64-a2a1-8f9ecd560c78/extras", nil)
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.ExtraSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `{
		"key1": "val1",
		"key2": "val2",
		"key3": {
			"subKey1": "subVal1",
			"subKey2": "subVal2"
		}
	}`

	var expected model.ExtraSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListExtraPropertiesWithBadRequest(t *testing.T) {

	mockClient := new(dbtest.Client)
	mockClient.On("ListExtraProperties", c.NewAdminContext(), "f4a5e666-c669-4c64-a2a1-8f9ecd560c78").Return(nil, errors.New("db error"))
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1beta/profiles/f4a5e666-c669-4c64-a2a1-8f9ecd560c78/extras", nil)
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}

func TestAddExtraProperty(t *testing.T) {

	mockClient := new(dbtest.Client)
	mockClient.On("AddExtraProperty", c.NewAdminContext(), "f4a5e666-c669-4c64-a2a1-8f9ecd560c78", fakeExtras).Return(&fakeExtras, nil)
	db.C = mockClient

	var fakeBody = `
		{
				"key1": "val1",
				"key2": "val2",
				"key3": {
					"subKey1": "subVal1",
					"subKey2": "subVal2"
				}
		}`
	r, _ := http.NewRequest("POST", "/v1beta/profiles/f4a5e666-c669-4c64-a2a1-8f9ecd560c78/extras", strings.NewReader(fakeBody))
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.ExtraSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `
		{
				"key1": "val1",
				"key2": "val2",
				"key3": {
					"subKey1": "subVal1",
					"subKey2": "subVal2"
				}	
		}`

	var expected model.ExtraSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestRemoveExtraProperty(t *testing.T) {

	mockClient := new(dbtest.Client)
	mockClient.On("RemoveExtraProperty", c.NewAdminContext(), "f4a5e666-c669-4c64-a2a1-8f9ecd560c78", "key1").Return(nil)
	db.C = mockClient

	r, _ := http.NewRequest("DELETE",
		"/v1beta/profiles/f4a5e666-c669-4c64-a2a1-8f9ecd560c78/extras/key1", nil)
	w := httptest.NewRecorder()
	beego.InsertFilter("*", beego.BeforeExec, func(httpCtx *context.Context) {
		httpCtx.Input.SetData("context", c.NewAdminContext())
	})
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}
}
