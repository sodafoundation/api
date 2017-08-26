// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/opensds/opensds/pkg/model"
	"github.com/satori/go.uuid"
)

type Setter interface {
	SetUuid(m model.Modeler) error

	SetCreatedTimeStamp(m model.Modeler) error

	SetUpdatedTimeStamp(m model.Modeler) error
}

type setter struct{}

func NewSetter() Setter {
	return &setter{}
}

func (s *setter) SetUuid(m model.Modeler) error {
	switch m.(type) {
	case model.VolumeSpec, model.VolumeSnapshotSpec, model.VolumeAttachmentSpec,
		model.ProfileSpec, model.DockSpec, model.StoragePoolSpec:
		// Set uuid.
		m.SetId(uuid.NewV4().String())

		return nil
	case *model.VolumeSpec, *model.VolumeSnapshotSpec, *model.VolumeAttachmentSpec,
		*model.ProfileSpec, *model.DockSpec, *model.StoragePoolSpec:
		// Set uuid.
		m.SetId(uuid.NewV4().String())

		return nil
	default:
		return errors.New("Unexpected input object format!")
	}
}

func (s *setter) SetCreatedTimeStamp(m model.Modeler) error {
	switch m.(type) {
	case model.VolumeSpec, model.VolumeSnapshotSpec, model.VolumeAttachmentSpec,
		model.ProfileSpec, model.DockSpec, model.StoragePoolSpec:
		// Set created time.
		m.SetCreatedTime(time.Now().Format(TimeFormat))

		return nil
	case *model.VolumeSpec, *model.VolumeSnapshotSpec, *model.VolumeAttachmentSpec,
		*model.ProfileSpec, *model.DockSpec, *model.StoragePoolSpec:
		// Set created time.
		m.SetCreatedTime(time.Now().Format(TimeFormat))

		return nil
	default:
		return errors.New("Unexpected input object format!")
	}
}

func (s *setter) SetUpdatedTimeStamp(m model.Modeler) error {
	switch m.(type) {
	case model.VolumeSpec, model.VolumeSnapshotSpec, model.VolumeAttachmentSpec,
		model.ProfileSpec, model.DockSpec, model.StoragePoolSpec:
		// Set updated time.
		m.SetUpdatedTime(time.Now().Format(TimeFormat))

		return nil
	case *model.VolumeSpec, *model.VolumeSnapshotSpec, *model.VolumeAttachmentSpec,
		*model.ProfileSpec, *model.DockSpec, *model.StoragePoolSpec:
		// Set updated time.
		m.SetUpdatedTime(time.Now().Format(TimeFormat))

		return nil
	default:
		return errors.New("Unexpected input object format!")
	}
}

type ErrorRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Title   string `json:"title"`
}

func ErrorStatus(code int, message string) []byte {
	status := &ErrorRes{
		Code:    code,
		Message: message,
		Title:   http.StatusText(code),
	}

	// Mashal the status.
	body, err := json.Marshal(status)
	if err != nil {
		log.Println("Failed to mashal error response:", err.Error())
		return []byte("Failed to mashal error response: " + err.Error())
	}
	return body
}

func Contained(obj, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	default:
		return false
	}
	return false
}
