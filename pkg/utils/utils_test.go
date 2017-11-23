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

package utils

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/model"
	fake "github.com/opensds/opensds/testutils/utils/testing"
)

var fset = fake.NewFakeSetter()

func InitializeModelers() []model.Modeler {
	var (
		modelers   []model.Modeler
		volume     = model.VolumeSpec{BaseModel: &model.BaseModel{}}
		attachment = model.VolumeAttachmentSpec{BaseModel: &model.BaseModel{}}
		snapshot   = model.VolumeSnapshotSpec{BaseModel: &model.BaseModel{}}
		profile    = model.ProfileSpec{BaseModel: &model.BaseModel{}}
		pool       = model.StoragePoolSpec{BaseModel: &model.BaseModel{}}
		dock       = model.DockSpec{BaseModel: &model.BaseModel{}}
	)

	modelers = append(modelers, volume, attachment, snapshot, profile, pool, dock)
	modelers = append(modelers, &volume, &attachment, &snapshot, &profile, &pool, &dock)

	return modelers
}

func InitializeModelersWithSomething(uuid, createdAt, updatedAt string) []model.Modeler {
	var (
		modelers   []model.Modeler
		volume     = model.VolumeSpec{BaseModel: &model.BaseModel{Id: uuid, CreatedAt: createdAt, UpdatedAt: updatedAt}}
		attachment = model.VolumeAttachmentSpec{BaseModel: &model.BaseModel{Id: uuid, CreatedAt: createdAt, UpdatedAt: updatedAt}}
		snapshot   = model.VolumeSnapshotSpec{BaseModel: &model.BaseModel{Id: uuid, CreatedAt: createdAt, UpdatedAt: updatedAt}}
		profile    = model.ProfileSpec{BaseModel: &model.BaseModel{Id: uuid, CreatedAt: createdAt, UpdatedAt: updatedAt}}
		pool       = model.StoragePoolSpec{BaseModel: &model.BaseModel{Id: uuid, CreatedAt: createdAt, UpdatedAt: updatedAt}}
		dock       = model.DockSpec{BaseModel: &model.BaseModel{Id: uuid, CreatedAt: createdAt, UpdatedAt: updatedAt}}
	)

	modelers = append(modelers, volume, attachment, snapshot, profile, pool, dock)
	modelers = append(modelers, &volume, &attachment, &snapshot, &profile, &pool, &dock)

	return modelers
}

func TestSetUuid(t *testing.T) {
	modelers := InitializeModelers()
	expectedModelers := InitializeModelersWithSomething(fset.Uuid, "", "")

	for i, model := range modelers {
		if ok := fset.SetUuid(model); ok != nil {
			t.Errorf("Failed to set uuid to model %v\n", model)
		}

		if !reflect.DeepEqual(model, expectedModelers[i]) {
			t.Errorf("Expected %v, got %v\n", expectedModelers[i], model)
		}
	}
}

func TestSetCreatedTimeStamp(t *testing.T) {
	modelers := InitializeModelers()
	expectedModelers := InitializeModelersWithSomething("", fset.CreatedTime, "")

	for i, model := range modelers {
		if ok := fset.SetCreatedTimeStamp(model); ok != nil {
			t.Errorf("Failed to set created time to model %v\n", model)
		}

		if !reflect.DeepEqual(model, expectedModelers[i]) {
			t.Errorf("Expected %v, got %v\n", expectedModelers[i], model)
		}
	}
}

func TestSetUpdatedTimeStamp(t *testing.T) {
	modelers := InitializeModelers()
	expectedModelers := InitializeModelersWithSomething("", "", fset.UpdatedTime)

	for i, model := range modelers {
		if ok := fset.SetUpdatedTimeStamp(model); ok != nil {
			t.Errorf("Failed to set updated time to model %v\n", model)
		}

		if !reflect.DeepEqual(model, expectedModelers[i]) {
			t.Errorf("Expected %v, got %v\n", expectedModelers[i], model)
		}
	}
}

func TestErrorStatus(t *testing.T) {
	var fakeErrorRes = ErrorRes{
		Code:    http.StatusAccepted,
		Message: "this is a test",
	}

	expected, err := json.Marshal(fakeErrorRes)
	if err != nil {
		t.Fatal(err)
	}

	result := ErrorStatus(fakeErrorRes.Code, fakeErrorRes.Message)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %v, get %v\n", expected, result)
	}
}

func TestValidateData(t *testing.T) {
	// First test.
	var data1 = &model.StoragePoolSpec{BaseModel: &model.BaseModel{}}
	var expected1 = &model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id:        fset.Uuid,
			CreatedAt: fset.CreatedTime,
		},
	}

	if err := ValidateData(data1, fset); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected1, data1) {
		t.Fatalf("Expected %v, get %v\n", expected1, data1)
	}

	// Second test.
	var data2 = &model.DockSpec{BaseModel: &model.BaseModel{}}
	var expected2 = &model.DockSpec{
		BaseModel: &model.BaseModel{
			Id:        fset.Uuid,
			CreatedAt: fset.CreatedTime,
		},
	}

	if err := ValidateData(data2, fset); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected2, data2) {
		t.Fatalf("Expected %v, get %v\n", expected2, data2)
	}
}

func TestContained(t *testing.T) {
	var targets = []interface{}{
		[]interface{}{"key01", 123, true},
		map[interface{}]string{
			"key01": "value01",
			true:    "value02",
			123:     "value03",
		},
	}
	var objs = []interface{}{"key01", 123, true}

	for _, obj := range objs {
		for _, target := range targets {
			if !Contained(obj, target) {
				t.Errorf("%v is not contained in %v\n", obj, target)
			}
		}
	}
}
