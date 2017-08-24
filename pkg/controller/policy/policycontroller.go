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

/*
This module implements the policy-based scheduling by parsing storage
profiles configured by admin.

*/

package policy

import (
	"github.com/opensds/opensds/pkg/controller/policy/executor"
	pb "github.com/opensds/opensds/pkg/grpc/opensds"
	api "github.com/opensds/opensds/pkg/model"
)

type Controller interface {
	Setup(flag int)

	StorageTag() *StorageTag

	ExecuteSyncPolicy(req *pb.DockRequest) error

	ExecuteAsyncPolicy(req *pb.DockRequest, in string, errChan chan error)
}

func NewController(profileSpec *api.ProfileSpec) Controller {
	return &controller{
		Profile: profileSpec,
	}
}

type controller struct {
	Profile *api.ProfileSpec
	Tag     *StorageTag
}

func (c *controller) Setup(flag int) {
	c.Tag = NewStorageTag(c.Profile.Extra, flag)
}

func (c *controller) StorageTag() *StorageTag {
	return c.Tag
}

func (c *controller) ExecuteSyncPolicy(vr *pb.DockRequest) error {
	swf, err := executor.RegisterSynchronizedWorkflow(vr, c.Tag.syncTag)
	if err != nil {
		return err
	}

	if err = executor.ExecuteSynchronizedWorkflow(swf); err != nil {
		return err
	}
	return nil
}

func (c *controller) ExecuteAsyncPolicy(vr *pb.DockRequest, in string, errChan chan error) {
	awf, err := executor.RegisterAsynchronizedWorkflow(vr, c.Tag.asyncTag, in)
	if err != nil {
		errChan <- err
	}

	defer close(errChan)

	if err = executor.ExecuteAsynchronizedWorkflow(awf); err != nil {
		errChan <- err
	}
	errChan <- nil
}
