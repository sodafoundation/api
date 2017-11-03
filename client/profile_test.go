// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package client

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/opensds/opensds/pkg/model"
)

var fpr = &ProfileMgr{
	Receiver: NewFakeProfileReceiver(),
}

func NewFakeProfileReceiver() Receiver {
	return &fakeProfileReceiver{}
}

type fakeProfileReceiver struct{}

func (*fakeProfileReceiver) Recv(
	f reqFunc,
	string,
	method string,
	in interface{},
	out interface{},
) error {
	switch strings.ToUpper(method) {
	case "POST":
		switch out.(type) {
		case *model.ProfileSpec:
			if err := json.Unmarshal([]byte(sampleProfile), out); err != nil {
				return err
			}
			break
		case *model.ExtraSpec:
			if err := json.Unmarshal([]byte(sampleExtras), out); err != nil {
				return err
			}
			break
		default:
			return errors.New("output format not supported!")
		}
		break
	case "GET":
		switch out.(type) {
		case *model.ProfileSpec:
			if err := json.Unmarshal([]byte(sampleProfile), out); err != nil {
				return err
			}
			break
		case *[]*model.ProfileSpec:
			if err := json.Unmarshal([]byte(sampleProfiles), out); err != nil {
				return err
			}
			break
		case *model.ExtraSpec:
			if err := json.Unmarshal([]byte(sampleExtras), out); err != nil {
				return err
			}
			break
		default:
			return errors.New("output format not supported!")
		}
		break
	case "DELETE":
		break
	default:
		return errors.New("inputed method format not supported!")
	}

	return nil
}

func TestCreateProfile(t *testing.T) {
	expected := &model.ProfileSpec{
		BaseModel: &model.BaseModel{
			Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		Name:        "default",
		Description: "default policy",
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
		},
		{
			BaseModel: &model.BaseModel{
				Id: "2f9c0a04-66ef-11e7-ade2-43158893e017",
			},
			Name:        "silver",
			Description: "silver policy",
			Extra: model.ExtraSpec{
				"diskType":  "SAS",
				"iops":      float64(300),
				"bandwidth": float64(500),
			},
		},
	}

	prfs, err := fpr.ListProfiles()
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(prfs, expected) {
		t.Errorf("Expected %v, got %v", expected, prfs)
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

func TestAddExtraProperty(t *testing.T) {
	var prfID = "2f9c0a04-66ef-11e7-ade2-43158893e017"
	expected := &model.ExtraSpec{
		"diskType":  "SAS",
		"iops":      float64(300),
		"bandwidth": float64(500),
	}

	exts, err := fpr.AddExtraProperty(prfID, &model.ExtraSpec{})
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(exts, expected) {
		t.Errorf("Expected %v, got %v", expected, exts)
		return
	}
}

func TestListExtraProperties(t *testing.T) {
	var prfID = "2f9c0a04-66ef-11e7-ade2-43158893e017"
	expected := &model.ExtraSpec{
		"diskType":  "SAS",
		"iops":      float64(300),
		"bandwidth": float64(500),
	}

	exts, err := fpr.ListExtraProperties(prfID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(exts, expected) {
		t.Errorf("Expected %v, got %v", expected, exts)
		return
	}
}

func TestRemoveExtraProperty(t *testing.T) {
	var prfID, extraKey = "2f9c0a04-66ef-11e7-ade2-43158893e017", "diskType"

	if err := fpr.RemoveExtraProperty(prfID, extraKey); err != nil {
		t.Error(err)
		return
	}
}

var (
	sampleProfile = `{
		"id": "1106b972-66ef-11e7-b172-db03f3689c9c",
		"name": "default",
		"description": "default policy"
	}`

	sampleProfiles = `[
		{
			"id": "1106b972-66ef-11e7-b172-db03f3689c9c",
			"name": "default",
			"description": "default policy"
		},
		{
			"id": "2f9c0a04-66ef-11e7-ade2-43158893e017",
			"name": "silver",
			"description": "silver policy",
			"extras": {
				"diskType":"SAS",
				"iops": 300,
				"bandwidth": 500
			}
		}
	]`

	sampleExtras = `{
		"diskType":"SAS",
		"iops": 300,
		"bandwidth": 500
	}`
)
