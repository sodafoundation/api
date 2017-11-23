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

package testing

import (
	"errors"

	"github.com/opensds/opensds/pkg/model"
)

func NewFakeSetter() *MockSetter {
	return &MockSetter{
		Uuid:        "0e9c3c68-8a0b-11e7-94a7-67f755e235cb",
		CreatedTime: "2017-08-26T11:01:09",
		UpdatedTime: "2017-08-26T11:01:55",
	}
}

type MockSetter struct {
	Uuid        string
	CreatedTime string
	UpdatedTime string
}

func (_m *MockSetter) SetCreatedTimeStamp(m model.Modeler) error {
	switch m.(type) {
	case model.VolumeSpec, model.VolumeSnapshotSpec, model.VolumeAttachmentSpec,
		model.ProfileSpec, model.DockSpec, model.StoragePoolSpec:
		break
	case *model.VolumeSpec, *model.VolumeSnapshotSpec, *model.VolumeAttachmentSpec,
		*model.ProfileSpec, *model.DockSpec, *model.StoragePoolSpec:
		break
	default:
		return errors.New("Unexpected input object format!")
	}

	m.SetCreatedTime(_m.CreatedTime)
	return nil
}

func (_m *MockSetter) SetUpdatedTimeStamp(m model.Modeler) error {
	switch m.(type) {
	case model.VolumeSpec, model.VolumeSnapshotSpec, model.VolumeAttachmentSpec,
		model.ProfileSpec, model.DockSpec, model.StoragePoolSpec:
		break
	case *model.VolumeSpec, *model.VolumeSnapshotSpec, *model.VolumeAttachmentSpec,
		*model.ProfileSpec, *model.DockSpec, *model.StoragePoolSpec:
		break
	default:
		return errors.New("Unexpected input object format!")
	}

	m.SetUpdatedTime(_m.UpdatedTime)
	return nil
}

func (_m *MockSetter) SetUuid(m model.Modeler) error {
	switch m.(type) {
	case model.VolumeSpec, model.VolumeSnapshotSpec, model.VolumeAttachmentSpec,
		model.ProfileSpec, model.DockSpec, model.StoragePoolSpec:
		break
	case *model.VolumeSpec, *model.VolumeSnapshotSpec, *model.VolumeAttachmentSpec,
		*model.ProfileSpec, *model.DockSpec, *model.StoragePoolSpec:
		break
	default:
		return errors.New("Unexpected input object format!")
	}

	m.SetId(_m.Uuid)
	return nil
}
