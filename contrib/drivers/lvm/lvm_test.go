// Copyright 2017 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package lvm

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils/config"
	"github.com/opensds/opensds/pkg/utils/exec"
)

var fp = map[string]PoolProperties{
	"vg001": {
		StorageType:      "block",
		AvailabilityZone: "default",
		MultiAttach:      true,
		Extras: model.StoragePoolExtraSpec{
			DataStorage: model.DataStorageLoS{
				ProvisioningPolicy: "Thin",
				IsSpaceEfficient:   false,
			},
			IOConnectivity: model.IOConnectivityLoS{
				AccessProtocol: "iscsi",
				MaxIOPS:        7000000,
				MaxBWS:         600,
			},
			Advanced: map[string]interface{}{
				"diskType": "SSD",
				"latency":  "5ms",
			},
		},
	},
}

func TestSetup(t *testing.T) {
	var d = &Driver{}
	config.CONF.OsdsDock.Backends.LVM.ConfigPath = "testdata/lvm.yaml"
	var expectedDriver = &Driver{
		conf: &LVMConfig{
			Pool:           fp,
			TgtBindIp:      "192.168.56.105",
			TgtConfDir:     "/etc/tgt/conf.d",
			EnableChapAuth: true,
		},
	}

	if err := d.Setup(); err != nil {
		t.Errorf("Setup lvm driver failed: %+v\n", err)
	}
	if !reflect.DeepEqual(d.conf, expectedDriver.conf) {
		t.Errorf("Expected %+v, got %+v", expectedDriver.conf, d.conf)
	}
}

type FakeResp struct {
	out string
	err error
}

func NewFakeExecuter(respMap map[string]*FakeResp) exec.Executer {
	return &FakeExecuter{RespMap: respMap}
}

type FakeExecuter struct {
	RespMap map[string]*FakeResp
}

func (f *FakeExecuter) Run(name string, args ...string) (string, error) {
	var cmd = name
	if name == "env" {
		cmd = args[1]
	}
	v, ok := f.RespMap[cmd]
	if !ok {
		return "", fmt.Errorf("can find specified op: %s", args[1])
	}
	return v.out, v.err
}

func TestCreateVolume(t *testing.T) {
	var fd = &Driver{}
	config.CONF.OsdsDock.Backends.LVM.ConfigPath = "testdata/lvm.yaml"
	fd.Setup()

	respMap := map[string]*FakeResp{
		"lvcreate": {"", nil},
	}
	fd.cli.RootExecuter = NewFakeExecuter(respMap)
	fd.cli.BaseExecuter = NewFakeExecuter(respMap)

	opt := &pb.CreateVolumeOpts{
		Id:          "e1bb066c-5ce7-46eb-9336-25508cee9f71",
		Name:        "test001",
		Description: "volume for testing",
		Size:        int64(1),
		PoolName:    "vg001",
	}
	var expected = &model.VolumeSpec{
		BaseModel:   &model.BaseModel{},
		Name:        "test001",
		Description: "volume for testing",
		Size:        int64(1),
		Metadata: map[string]string{
			"lvPath": "/dev/vg001/volume-e1bb066c-5ce7-46eb-9336-25508cee9f71",
		},
	}
	vol, err := fd.CreateVolume(opt)
	if err != nil {
		t.Error("Failed to create volume:", err)
	}
	vol.Id = ""
	if !reflect.DeepEqual(vol, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, vol)
	}
}

func TestCreateVolumeFromSnapshot(t *testing.T) {
	var fd = &Driver{}
	config.CONF.OsdsDock.Backends.LVM.ConfigPath = "testdata/lvm.yaml"
	fd.Setup()

	respMap := map[string]*FakeResp{
		"lvcreate": {"", nil},
		"dd":       {"", nil},
	}
	fd.cli.RootExecuter = NewFakeExecuter(respMap)
	fd.cli.BaseExecuter = NewFakeExecuter(respMap)

	opt := &pb.CreateVolumeOpts{
		Id:           "e1bb066c-5ce7-46eb-9336-25508cee9f71",
		Name:         "test001",
		Description:  "volume for testing",
		Size:         int64(1),
		PoolName:     "vg001",
		SnapshotId:   "3769855c-a102-11e7-b772-17b880d2f537",
		SnapshotSize: int64(1),
	}
	var expected = &model.VolumeSpec{
		BaseModel:   &model.BaseModel{},
		Name:        "test001",
		Description: "volume for testing",
		Size:        int64(1),
		Metadata: map[string]string{
			"lvPath": "/dev/vg001/volume-e1bb066c-5ce7-46eb-9336-25508cee9f71",
		},
	}
	vol, err := fd.CreateVolume(opt)
	if err != nil {
		t.Error("Failed to create volume:", err)
	}
	vol.Id = ""
	if !reflect.DeepEqual(vol, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, vol)
	}
}

func TestDeleteVolume(t *testing.T) {
	var fd = &Driver{}
	config.CONF.OsdsDock.Backends.LVM.ConfigPath = "testdata/lvm.yaml"
	fd.Setup()

	respMap := map[string]*FakeResp{
		"lvdisplay": {"-wi-a-----", nil},
		"lvremove":  {"", nil},
	}
	fd.cli.RootExecuter = NewFakeExecuter(respMap)
	fd.cli.BaseExecuter = NewFakeExecuter(respMap)

	opt := &pb.DeleteVolumeOpts{
		Metadata: map[string]string{
			"lvPath": "/dev/vg001/test001",
		},
	}
	if err := fd.DeleteVolume(opt); err != nil {
		t.Error("Failed to delete volume:", err)
	}
}

func TestExtendVolume(t *testing.T) {
	var fd = &Driver{}
	config.CONF.OsdsDock.Backends.LVM.ConfigPath = "testdata/lvm.yaml"
	fd.Setup()

	respMap := map[string]*FakeResp{
		"lvdisplay": {"-wi-a-----", nil},
		"lvchange":  {"", nil},
		"lvextend":  {"", nil},
	}
	fd.cli.RootExecuter = NewFakeExecuter(respMap)
	fd.cli.BaseExecuter = NewFakeExecuter(respMap)

	opt := &pb.ExtendVolumeOpts{
		Id: "591c43e6-1156-42f5-9fbc-161153da185c",
		Metadata: map[string]string{
			"lvPath": "/dev/vg001/test001",
		},
		Size: int64(1),
	}

	vol, err := fd.ExtendVolume(opt)
	if err != nil {
		t.Error("Failed to extend volume:", err)
	}

	if vol.Size != 1 {
		t.Errorf("Expected %+v, got %+v\n", 1, vol.Size)
	}
}

func TestCreateSnapshot(t *testing.T) {
	var fd = &Driver{}
	config.CONF.OsdsDock.Backends.LVM.ConfigPath = "testdata/lvm.yaml"
	fd.Setup()

	respMap := map[string]*FakeResp{
		"lvcreate": {"-wi-a-----", nil},
	}
	fd.cli.RootExecuter = NewFakeExecuter(respMap)
	fd.cli.BaseExecuter = NewFakeExecuter(respMap)

	opt := &pb.CreateVolumeSnapshotOpts{
		Id:          "d1916c49-3088-4a40-b6fb-0fda18d074c3",
		Name:        "snap001",
		Description: "volume snapshot for testing",
		Size:        int64(1),
		VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Metadata: map[string]string{
			"lvPath": "/dev/vg001/test001",
		},
	}
	var expected = &model.VolumeSnapshotSpec{
		BaseModel:   &model.BaseModel{},
		Name:        "snap001",
		Description: "volume snapshot for testing",
		Size:        int64(1),
		VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
		Metadata: map[string]string{
			"lvsPath": "/dev/vg001/_snapshot-d1916c49-3088-4a40-b6fb-0fda18d074c3",
		},
	}
	snp, err := fd.CreateSnapshot(opt)
	if err != nil {
		t.Error("Failed to create volume snapshot:", err)
	}
	snp.Id = ""
	snp.Metadata["lvsPath"] = "/dev/vg001/_snapshot-d1916c49-3088-4a40-b6fb-0fda18d074c3"
	if !reflect.DeepEqual(snp, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, snp)
	}
}

func TestDeleteSnapshot(t *testing.T) {
	var fd = &Driver{}
	config.CONF.OsdsDock.Backends.LVM.ConfigPath = "testdata/lvm.yaml"
	fd.Setup()

	lvsResp := `  _snapshot-f0594d2b-ffdf-4947-8380-089f0bc17389
  volume-0e2f4a9e-4a94-4d27-b1b4-83464811605c
  volume-591c43e6-1156-42f5-9fbc-161153da185c
  root
  swap_1
`
	respMap := map[string]*FakeResp{
		"lvs":       {lvsResp, nil},
		"lvdisplay": {"-wi-a-----", nil},
		"lvremove":  {"", nil},
	}
	fd.cli.RootExecuter = NewFakeExecuter(respMap)
	fd.cli.BaseExecuter = NewFakeExecuter(respMap)

	opt := &pb.DeleteVolumeSnapshotOpts{
		Metadata: map[string]string{
			"lvsPath": "/dev/vg001/snap001",
		},
	}

	if err := fd.DeleteSnapshot(opt); err != nil {
		t.Error("Failed to delete volume snapshot:", err)
	}
}

func TestListPools(t *testing.T) {
	var fd = &Driver{}
	config.CONF.OsdsDock.Backends.LVM.ConfigPath = "testdata/lvm.yaml"
	fd.Setup()

	var vgsResp = `  vg001  18.00 18.00 ahF6kS-QNOH-X63K-avat-6Kag-XLTo-c9ghQ6
  ubuntu-vg               127.52  0.03 fQbqtg-3vDQ-vk3U-gfsT-50kJ-30pq-OZVSJH
`
	respMap := map[string]*FakeResp{
		"vgs": {vgsResp, nil},
	}
	fd.cli.RootExecuter = NewFakeExecuter(respMap)
	fd.cli.BaseExecuter = NewFakeExecuter(respMap)

	var expected = []*model.StoragePoolSpec{
		{
			BaseModel:        &model.BaseModel{},
			Name:             "vg001",
			TotalCapacity:    int64(18),
			FreeCapacity:     int64(18),
			AvailabilityZone: "default",
			StorageType:      "block",
			MultiAttach:      true,
			Extras: model.StoragePoolExtraSpec{
				DataStorage: model.DataStorageLoS{
					ProvisioningPolicy: "Thin",
					IsSpaceEfficient:   false,
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "iscsi",
					MaxIOPS:        7000000,
					MaxBWS:         600,
				},
				Advanced: map[string]interface{}{
					"diskType": "SSD",
					"latency":  "5ms",
				},
			},
		},
	}

	pols, err := fd.ListPools()
	if err != nil {
		t.Error("Failed to list pools:", err)
	}
	for i := range pols {
		pols[i].Id = ""
	}
	if !reflect.DeepEqual(pols, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected[0], pols[0])
	}
}
