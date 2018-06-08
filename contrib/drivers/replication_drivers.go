// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package drivers

import (
	"reflect"

	"github.com/opensds/opensds/contrib/drivers/drbd"
	"github.com/opensds/opensds/contrib/drivers/huawei/dorado"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/config"
	replication_sample "github.com/opensds/opensds/testutils/driver"
)

// ReplicationDriver is an interface for exposing some operations of different
// replication drivers, currently supporting DRBD.
type ReplicationDriver interface {
	// Any initialization the replication driver does while starting.
	Setup() error
	// Any operation the replication driver does while stopping.
	Unset() error

	CreateReplication(opt *pb.CreateReplicationOpts) (*model.ReplicationSpec, error)
	DeleteReplication(opt *pb.DeleteReplicationOpts) error
	EnableReplication(opt *pb.EnableReplicationOpts) error
	DisableReplication(opt *pb.DisableReplicationOpts) error
	FailoverReplication(opt *pb.FailoverReplicationOpts) error
}

func IsSupportHostBasedReplication(resourceType string) bool {
	v := reflect.ValueOf(config.CONF.Backends)
	t := reflect.TypeOf(config.CONF.Backends)
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("conf")
		if resourceType == tag && field.Interface().(config.BackendProperties).SupportReplication {
			return true
		}
	}
	return false
}

// Init
func InitReplicationDriver(resourceType string) (ReplicationDriver, error) {
	var d ReplicationDriver
	switch resourceType {
	case "drbd":
		d = &drbd.ReplicationDriver{}
		break
	case "huawei_dorado":
		d = &dorado.ReplicationDriver{}
		break
	default:
		d = &replication_sample.ReplicationDriver{}
		break
	}
	err := d.Setup()
	return d, err
}

// Clean
func CleanReplicationDriver(d ReplicationDriver) ReplicationDriver {
	// Execute different clean operations according to the ReplicationDriver type.
	switch d.(type) {
	case *drbd.ReplicationDriver:
		break
	case *dorado.ReplicationDriver:
		d = &dorado.ReplicationDriver{}
	default:
		break
	}
	d.Unset()
	d = nil

	return d
}
