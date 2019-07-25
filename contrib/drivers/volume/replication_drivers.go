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
This module defines an standard table of storage driver. The default storage
driver is sample driver used for testing. If you want to use other storage
plugin, just modify Init() and Clean() method.

*/

package volume

import (
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
)

// VolumeReplicationDriver is an interface for exposing some operations of different
// replication drivers, currently supporting DRBD.
type VolumeReplicationDriver interface {
	// Any initialization the replication driver does while starting.
	Setup() error
	// Any operation the replication driver does while stopping.
	Teardown() error

	CreateReplication(opt *pb.CreateReplicationOpts) (*model.ReplicationSpec, error)
	DeleteReplication(opt *pb.DeleteReplicationOpts) error
	EnableReplication(opt *pb.EnableReplicationOpts) error
	DisableReplication(opt *pb.DisableReplicationOpts) error
	FailoverReplication(opt *pb.FailoverReplicationOpts) error
}
