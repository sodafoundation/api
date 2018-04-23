// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package dorado

import (
	"github.com/opensds/opensds/pkg/dock/proto"
	"testing"
)

func TestLoadConf(t *testing.T) {
}

func TestCreateReplication(t *testing.T) {
	d := &ReplicationDriver{}
	err := d.Setup()
	if err != nil {
		t.Error("set up ...", err)
	}
	opt := &proto.CreateReplicationOpts{
		Id: "29693780-260a-47b7-91bc-bd03f0bf4f68",
		PrimaryReplicationDriverData:   map[string]string{KLunId: "362"},
		SecondaryReplicationDriverData: map[string]string{KLunId: "4761"},
		ReplicationMode:                "ASYNC",
	}
	d.CreateReplication(opt)
}

func TestEnableReplication(t *testing.T) {
	d := &ReplicationDriver{}
	err := d.Setup()
	if err != nil {
		t.Error("set up ...", err)
	}
	opt := &proto.EnableReplicationOpts{
		Id: "29693780-260a-47b7-91bc-bd03f0bf4f68",
		PrimaryReplicationDriverData:   map[string]string{KLunId: "362"},
		SecondaryReplicationDriverData: map[string]string{KLunId: "4761"},
		Metadata:                       map[string]string{KPairId: "7079902c95e8000c"},
	}
	d.EnableReplication(opt)
}

func TestDisableReplication(t *testing.T) {
	d := &ReplicationDriver{}
	err := d.Setup()
	if err != nil {
		t.Error("set up ...", err)
	}
	opt := &proto.DisableReplicationOpts{
		Id: "29693780-260a-47b7-91bc-bd03f0bf4f68",
		PrimaryReplicationDriverData:   map[string]string{KLunId: "362"},
		SecondaryReplicationDriverData: map[string]string{KLunId: "4761"},
		Metadata:                       map[string]string{KPairId: "7079902c95e8000c"},
	}
	d.DisableReplication(opt)
}

func TestFailoverReplication(t *testing.T) {
	d := &ReplicationDriver{}
	err := d.Setup()
	if err != nil {
		t.Error("set up ...", err)
	}
	opt := &proto.FailoverReplicationOpts{
		Id: "29693780-260a-47b7-91bc-bd03f0bf4f68",
		PrimaryReplicationDriverData:   map[string]string{KLunId: "362"},
		SecondaryReplicationDriverData: map[string]string{KLunId: "4761"},
		Metadata:                       map[string]string{KPairId: "7079902c95e8000c"},
	}
	d.FailoverReplication(opt)
}

func TestDeleteReplication(t *testing.T) {
	d := &ReplicationDriver{}
	err := d.Setup()
	if err != nil {
		t.Error("set up ...", err)
	}
	opt := &proto.DeleteReplicationOpts{
		Id:       "29693780-260a-47b7-91bc-bd03f0bf4f68",
		Metadata: map[string]string{KPairId: "7079902c95e8000c"},
	}
	d.DeleteReplication(opt)
}
