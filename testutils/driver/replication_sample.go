// Copyright 2018 The OpenSDS Authors.
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
This module implements a sample driver for OpenSDS. This driver will handle all
operations of volume and return a fake value.

*/

package sample

import (
	//"errors"

	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	. "github.com/opensds/opensds/testutils/collection"
)

// ReplicationDriver
type ReplicationDriver struct{}

// Setup
func (r *ReplicationDriver) Setup() error { return nil }

// Unset
func (r *ReplicationDriver) Unset() error { return nil }

// CreateReplication
func (r *ReplicationDriver) CreateReplication(opt *pb.CreateReplicationOpts) (*model.ReplicationSpec, error) {
	return &SampleReplications[0], nil
}

func (r *ReplicationDriver) DeleteReplication(opt *pb.DeleteReplicationOpts) error {
	return nil
}

func (r *ReplicationDriver) EnableReplication(opt *pb.EnableReplicationOpts) error {
	return nil
}

func (r *ReplicationDriver) DisableReplication(opt *pb.DisableReplicationOpts) error {
	return nil
}

func (r *ReplicationDriver) FailoverReplication(opt *pb.FailoverReplicationOpts) error {
	return nil
}
