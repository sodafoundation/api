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

package mocks

import (
	"github.com/opensds/opensds/pkg/model"
)

type MockSetter struct {
	Uuid        string
	CreatedTime string
	UpdatedTime string
}

func (_m *MockSetter) SetCreatedTimeStamp(m model.Modeler) error {
	m.SetCreatedTime(_m.CreatedTime)
	return nil
}

func (_m *MockSetter) SetUpdatedTimeStamp(m model.Modeler) error {
	m.SetUpdatedTime(_m.UpdatedTime)
	return nil
}

func (_m *MockSetter) SetUuid(m model.Modeler) error {
	m.SetId(_m.Uuid)
	return nil
}
