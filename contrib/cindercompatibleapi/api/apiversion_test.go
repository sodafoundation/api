// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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
	c "github.com/opensds/opensds/client"
	"github.com/opensds/opensds/contrib/cindercompatibleapi/converter"
)

func init() {
	beego.Router("/", &VersionPortal{},
		"get:ListAllAPIVersions")

	if false == IsFakeClient {
		client = NewFakeClient(&c.Config{Endpoint: TestEp})
	}
}

////////////////////////////////////////////////////////////////////////////////
//                            Tests for Version                              //
////////////////////////////////////////////////////////////////////////////////
func TestListAllAPIVersions(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output converter.ListAllAPIVersionsRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)
	//fmt.Println(string(w.Body.Bytes()))
	expectedJSON := `
    {
		"versions": [{
			"status": "CURRENT",
			"updated": "2017-07-10T14:36:58.014Z",
			"min_version": "3.0",
			"id": "v3.0"
		}]
	}`

	var expected converter.ListAllAPIVersionsRespSpec
	json.Unmarshal([]byte(expectedJSON), &expected)

	if w.Code != http.StatusMultipleChoices {
		t.Errorf("Expected %v, actual %v", http.StatusMultipleChoices, w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}
