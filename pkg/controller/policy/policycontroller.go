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

/*
This module implements the policy-based scheduling by parsing storage
profiles configured by admin.

*/

package policy

import (
	"github.com/opensds/opensds/pkg/controller/policy/executor"
	"github.com/opensds/opensds/pkg/model"
)

type Controller interface {
	Setup(flag int)

	StorageTag() *StorageTag

	ExecuteSyncPolicy(req interface{}) error

	ExecuteAsyncPolicy(req interface{}, in string, errChan chan error)

	SetDock(dockInfo *model.DockSpec)
}

func NewController(profileSpec *model.ProfileSpec) Controller {
	return &controller{
		Profile: profileSpec,
	}
}

type controller struct {
	Profile  *model.ProfileSpec
	DockInfo *model.DockSpec
	Tag      *StorageTag
}

func (c *controller) Setup(flag int) {
	c.Tag = NewStorageTag(c.Profile.Extras, flag)
}

func (c *controller) StorageTag() *StorageTag {
	return c.Tag
}

func (c *controller) ExecuteSyncPolicy(req interface{}) error {
	swf, err := executor.RegisterSynchronizedWorkflow(req, c.Tag.syncTag)
	if err != nil {
		return err
	}

	if err = executor.ExecuteSynchronizedWorkflow(swf); err != nil {
		return err
	}
	return nil
}

func (c *controller) ExecuteAsyncPolicy(req interface{}, in string, errChan chan error) {
	awf, err := executor.RegisterAsynchronizedWorkflow(req, c.Tag.asyncTag, c.DockInfo, in)
	if err != nil {
		errChan <- err
	}

	defer close(errChan)

	if err = executor.ExecuteAsynchronizedWorkflow(awf); err != nil {
		errChan <- err
	}
	errChan <- nil
}

func (c *controller) SetDock(dockInfo *model.DockSpec) {
	c.DockInfo = dockInfo
}
