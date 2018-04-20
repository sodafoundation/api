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
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/astaxie/beego"
	c "github.com/opensds/opensds/client"
	"github.com/opensds/opensds/plugin/CinderCompatibleAPI/cindermodel"
)

func init() {
	beego.Router("/V3/snapshots/:snapshotId", &SnapshotPortal{},
		"get:GetSnapshot;delete:DeleteSnapshot;put:UpdateSnapshot")
	beego.Router("/V3/snapshots", &SnapshotPortal{},
		"post:CreateSnapshot;get:ListSnapshot")
	beego.Router("/V3/snapshots/detail", &SnapshotPortal{},
		"get:ListSnapshotDetail")
	if false == IsFakeClient {
		client = NewFakeClient(&c.Config{Endpoint: TestEp})
	}
}

////////////////////////////////////////////////////////////////////////////////
//                            Tests for Snapshot                              //
////////////////////////////////////////////////////////////////////////////////
func TestCreateSnapshot(t *testing.T) {
	RequestBodyStr := `{
    	"snapshot": {
        "name": "sample-snapshot-01",
		"description": "This is the first sample snapshot for testing",
		"volume_id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
		"metadata": null
		}
	}`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("POST", "/V3/snapshots", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output cindermodel.CreateSnapshotRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	var expected cindermodel.CreateSnapshotRespSpec
	json.Unmarshal([]byte(RequestBodyStr), &expected)

	if w.Code != StatusAccepted {
		t.Errorf("Expected %v, actual %v", StatusAccepted, w.Code)
	}

	expected.Snapshot.ID = "3769855c-a102-11e7-b772-17b880d2f537"
	expected.Snapshot.VolumeID = "bd5b12a8-a101-11e7-941e-d77981b584d8"
	expected.Snapshot.Status = "available"
	expected.Snapshot.Size = 1
	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestGetSnapshot(t *testing.T) {
	r, _ := http.NewRequest("GET", "/V3/snapshots/f2dda3d2-bf79-11e7-8665-f750b088f63e", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output cindermodel.ShowSnapshotDetailsResp
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJSON := `{
    	"snapshot": {
		"id": "3769855c-a102-11e7-b772-17b880d2f537",
        "name": "sample-snapshot-01",
		"description": "This is the first sample snapshot for testing",
		"volume_id": "bd5b12a8-a101-11e7-941e-d77981b584d8"
		}
	}`

	var expected cindermodel.ShowSnapshotDetailsResp
	json.Unmarshal([]byte(expectedJSON), &expected)

	if w.Code != StatusOK {
		t.Errorf("Expected %v, actual %v", StatusOK, w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListSnapshot(t *testing.T) {
	r, _ := http.NewRequest("GET", "/V3/snapshots", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output []cindermodel.ListSnapshotRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJSON :=
		`{"snapshots":[
		{"status":"created",
		 "description":"This is the first sample snapshot for testing",
		 "name":"sample-snapshot-01",
		 "volume_id":"bd5b12a8-a101-11e7-941e-d77981b584d8",
		 "id":"3769855c-a102-11e7-b772-17b880d2f537",
		 "size":1},
		{"status":"created",
		 "description":"This is the second sample snapshot for testing",
		 "name":"sample-snapshot-02",
		 "volume_id":"bd5b12a8-a101-11e7-941e-d77981b584d8",
		 "id":"3bfaf2cc-a102-11e7-8ecb-63aea739d755","size":1
		}
		]}`

	var expected []cindermodel.ListSnapshotRespSpec
	json.Unmarshal([]byte(expectedJSON), &expected)

	if w.Code != StatusOK {
		t.Errorf("Expected %v, actual %v", StatusOK, w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListSnapshotDetail(t *testing.T) {
	r, _ := http.NewRequest("GET", "/V3/snapshots/detail", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output []cindermodel.ListSnapshotDetailRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJSON :=
		`{"snapshots":[
		{"status":"created",
		 "description":"This is the first sample snapshot for testing",
		 "name":"sample-snapshot-01",
		 "volume_id":"bd5b12a8-a101-11e7-941e-d77981b584d8",
		 "id":"3769855c-a102-11e7-b772-17b880d2f537",
		 "size":1},
		{"status":"created",
		 "description":"This is the second sample snapshot for testing",
		 "name":"sample-snapshot-02",
		 "volume_id":"bd5b12a8-a101-11e7-941e-d77981b584d8",
		 "id":"3bfaf2cc-a102-11e7-8ecb-63aea739d755","size":1
		}
		]}`

	var expected []cindermodel.ListSnapshotDetailRespSpec
	json.Unmarshal([]byte(expectedJSON), &expected)

	if w.Code != StatusOK {
		t.Errorf("Expected %v, actual %v", StatusOK, w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestDeleteSnapshot(t *testing.T) {
	r, _ := http.NewRequest("DELETE", "/V3/snapshots/3769855c-a102-11e7-b772-17b880d2f537", nil)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != StatusAccepted {
		t.Errorf("Expected %v, actual %v", StatusAccepted, w.Code)
	}
}

func TestUpdateSnapshot(t *testing.T) {
	RequestBodyStr := `{
    	"snapshot": {
        "name": "sample-snapshot-01",
		"description": "This is the first sample snapshot for testing"
		}
	}`

	var jsonStr = []byte(RequestBodyStr)
	r, _ := http.NewRequest("PUT", "/V3/snapshots/3769855c-a102-11e7-b772-17b880d2f537", bytes.NewBuffer(jsonStr))

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output cindermodel.UpdateSnapshotRespSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	var expected cindermodel.UpdateSnapshotRespSpec
	json.Unmarshal([]byte(RequestBodyStr), &expected)

	if w.Code != StatusOK {
		t.Errorf("Expected %v, actual %v", StatusOK, w.Code)
	}

	expected.Snapshot.ID = "3769855c-a102-11e7-b772-17b880d2f537"
	expected.Snapshot.VolumeID = "bd5b12a8-a101-11e7-941e-d77981b584d8"
	expected.Snapshot.Status = "available"
	expected.Snapshot.Size = 1

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}
