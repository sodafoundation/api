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

package client

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/model"
)

var fpr = &ProfileMgr{
	Receiver: NewFakeProfileReceiver(),
}

func TestCreateProfile(t *testing.T) {
	expected := &model.ProfileSpec{
		BaseModel: &model.BaseModel{
			Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		Name:        "default",
		Description: "default policy",
		StorageType: "block",
	}

	prf, err := fpr.CreateProfile(&model.ProfileSpec{})
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(prf, expected) {
		t.Errorf("Expected %v, got %v", expected, prf)
		return
	}
}

func TestGetProfile(t *testing.T) {
	var prfID = "1106b972-66ef-11e7-b172-db03f3689c9c"
	expected := &model.ProfileSpec{
		BaseModel: &model.BaseModel{
			Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		Name:        "default",
		Description: "default policy",
		StorageType: "block",
	}

	prf, err := fpr.GetProfile(prfID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(prf, expected) {
		t.Errorf("Expected %v, got %v", expected, prf)
		return
	}
}

func TestListProfiles(t *testing.T) {
	expected := []*model.ProfileSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
			},
			Name:        "default",
			Description: "default policy",
			StorageType: "block",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "2f9c0a04-66ef-11e7-ade2-43158893e017",
			},
			Name:        "silver",
			Description: "silver policy",
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
			},
		},
	}

	prfs, err := fpr.ListProfiles()
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(prfs, expected) {
		t.Errorf("Expected %v, got %v", expected[1], prfs[1])
		return
	}
}

func TestDeleteProfile(t *testing.T) {
	var prfID = "1106b972-66ef-11e7-b172-db03f3689c9c"

	if err := fpr.DeleteProfile(prfID); err != nil {
		t.Error(err)
		return
	}
}

func TestAddCustomProperty(t *testing.T) {
	var prfID = "2f9c0a04-66ef-11e7-ade2-43158893e017"
	expected := &model.CustomPropertiesSpec{
		"dataStorage": map[string]interface{}{
			"provisioningPolicy": "Thin",
			"isSpaceEfficient":   true,
		},
		"ioConnectivity": map[string]interface{}{
			"accessProtocol": "rbd",
			"maxIOPS":        float64(5000000),
			"maxBWS":         float64(500),
		},
	}

	cps, err := fpr.AddCustomProperty(prfID, &model.CustomPropertiesSpec{})
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(cps, expected) {
		t.Errorf("Expected %v, got %v", expected, cps)
		return
	}
}

func TestListCustomProperties(t *testing.T) {
	var prfID = "2f9c0a04-66ef-11e7-ade2-43158893e017"
	expected := &model.CustomPropertiesSpec{
		"dataStorage": map[string]interface{}{
			"provisioningPolicy": "Thin",
			"isSpaceEfficient":   true,
		},
		"ioConnectivity": map[string]interface{}{
			"accessProtocol": "rbd",
			"maxIOPS":        float64(5000000),
			"maxBWS":         float64(500),
		},
	}

	cps, err := fpr.ListCustomProperties(prfID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(cps, expected) {
		t.Errorf("Expected %v, got %v", expected, cps)
		return
	}
}

func TestRemoveCustomProperty(t *testing.T) {
	var prfID, customKey = "2f9c0a04-66ef-11e7-ade2-43158893e017", "diskType"

	if err := fpr.RemoveCustomProperty(prfID, customKey); err != nil {
		t.Error(err)
		return
	}
}
