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

package integration

import (
	"fmt"
	"testing"

	"github.com/opensds/opensds/pkg/controller/volume"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
)

var vc = volume.NewController(
	&pb.CreateVolumeOpts{},
	&pb.DeleteVolumeOpts{},
	&pb.CreateVolumeSnapshotOpts{},
	&pb.DeleteVolumeSnapshotOpts{},
	&pb.CreateAttachmentOpts{},
)

var dckInfo = &model.DockSpec{
	Endpoint:   "localhost:50050",
	DriverName: "default",
}

func TestCreateVolume(t *testing.T) {
	vc.SetDock(dckInfo)

	vol, err := vc.CreateVolume()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(vol)
}

func TestDeleteVolume(t *testing.T) {
	vc.SetDock(dckInfo)

	if res := vc.DeleteVolume(); res.GetStatus() == "Failure" {
		t.Error(res.GetError())
	}
}

func TestCreateVolumeAttachment(t *testing.T) {
	vc.SetDock(dckInfo)

	atc, err := vc.CreateVolumeAttachment()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(atc)
}

func TestCreateVolumeSnapshot(t *testing.T) {
	vc.SetDock(dckInfo)

	snp, err := vc.CreateVolumeSnapshot()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(snp)
}

func TestDeleteVolumeSnapshot(t *testing.T) {
	vc.SetDock(dckInfo)

	if res := vc.DeleteVolumeSnapshot(); res.GetStatus() == "Failure" {
		t.Error(res.GetError())
	}
}
