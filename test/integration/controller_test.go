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

// +build integration

package integration

import (
	"reflect"
	"testing"

	"github.com/opensds/opensds/pkg/controller/volume"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	. "github.com/opensds/opensds/testutils/collection"
)

var (
	vc      = volume.NewController()
	dckInfo = &model.DockSpec{
		Endpoint:   "localhost:50050",
		DriverName: "default",
	}
)

func TestControllerCreateVolume(t *testing.T) {
	vc.SetDock(dckInfo)

	vol, err := vc.CreateVolume(&pb.CreateVolumeOpts{})
	if err != nil {
		t.Error("create volume in controller failed:", err)
		return
	}

	var expected = &SampleVolumes[0]
	if !reflect.DeepEqual(vol, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, vol)
	}
}

func TestControllerDeleteVolume(t *testing.T) {
	vc.SetDock(dckInfo)

	err := vc.DeleteVolume(&pb.DeleteVolumeOpts{})
	if err != nil {
		t.Error("delete volume in controller failed:", err)
	}
}

func TestControllerExtendVolume(t *testing.T) {
	vc.SetDock(dckInfo)

	vol, err := vc.ExtendVolume(&pb.ExtendVolumeOpts{})
	if err != nil {
		t.Error("extend volume in controller failed:", err)
		return
	}

	var expected = &SampleVolumes[0]
	if !reflect.DeepEqual(vol, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, vol)
	}
}

func TestControllerCreateVolumeAttachment(t *testing.T) {
	vc.SetDock(dckInfo)

	atc, err := vc.CreateVolumeAttachment(&pb.CreateVolumeAttachmentOpts{})
	if err != nil {
		t.Error("create volume attachment in controller failed:", err)
		return
	}

	var expected = &model.VolumeAttachmentSpec{
		BaseModel:      &model.BaseModel{},
		ConnectionInfo: SampleConnection,
	}
	if !reflect.DeepEqual(atc, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, atc)
	}
}

func TestControllerDeleteVolumeAttachment(t *testing.T) {
	vc.SetDock(dckInfo)

	err := vc.DeleteVolumeAttachment(&pb.DeleteVolumeAttachmentOpts{})
	if err != nil {
		t.Error("delete volume attachment in controller failed:", err)
	}
}

func TestControllerCreateVolumeSnapshot(t *testing.T) {
	vc.SetDock(dckInfo)

	snp, err := vc.CreateVolumeSnapshot(&pb.CreateVolumeSnapshotOpts{})
	if err != nil {
		t.Error("create volume snapshot in controller failed:", err)
		return
	}

	var expected = &SampleSnapshots[0]
	if !reflect.DeepEqual(snp, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, snp)
	}
}

func TestControllerDeleteVolumeSnapshot(t *testing.T) {
	vc.SetDock(dckInfo)

	err := vc.DeleteVolumeSnapshot(&pb.DeleteVolumeSnapshotOpts{})
	if err != nil {
		t.Error("delete volume snapshot in controller failed:", err)
	}
}

func TestControllerCreateVolumeGroup(t *testing.T) {
	vc.SetDock(dckInfo)

	vg, err := vc.CreateVolumeGroup(&pb.CreateVolumeGroupOpts{})
	if err != nil {
		t.Error("create volume group in controller failed:", err)
		return
	}

	var expected = &SampleVolumeGroups[0]
	if !reflect.DeepEqual(vg, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, vg)
	}
}

func TestControllerUpdateVolumeGroup(t *testing.T) {
	vc.SetDock(dckInfo)

	vg, err := vc.UpdateVolumeGroup(&pb.UpdateVolumeGroupOpts{})
	if err != nil {
		t.Error("update volume group in controller failed:", err)
		return
	}

	var expected = &SampleVolumeGroups[0]
	if !reflect.DeepEqual(vg, expected) {
		t.Errorf("expected %+v, got %+v\n", expected, vg)
	}
}

func TestControllerDeleteVolumeGroup(t *testing.T) {
	vc.SetDock(dckInfo)

	err := vc.DeleteVolumeGroup(&pb.DeleteVolumeGroupOpts{})
	if err != nil {
		t.Error("delete volume group in controller failed:", err)
	}
}
