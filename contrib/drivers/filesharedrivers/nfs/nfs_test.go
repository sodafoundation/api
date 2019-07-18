// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package nfs

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/opensds/opensds/contrib/drivers/filesharedrivers/nfs"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils/config"
	"github.com/opensds/opensds/pkg/utils/exec"
)

var fp = map[string]PoolProperties{
	"vg001": {
		StorageType:      "file",
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
	var d = &nfs.Driver{}
	config.CONF.OsdsDock.Backends.NFS.ConfigPath = "testdata/nfs.yaml"
	var expectedDriver = &Driver{
		conf: &NFSConfig{
			Pool:           fp,
			TgtBindIp:      "11.242.178.20",
			TgtConfDir:     "/etc/tgt/conf.d",
			EnableChapAuth: true,
		},
	}

	if err := d.Setup(); err != nil {
		t.Errorf("Setup nfs driver failed: %+v\n", err)
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

func TestCreateFileShare(t *testing.T) {
	var fd = &Driver{}
	config.CONF.OsdsDock.Backends.NFS.ConfigPath = "testdata/nfs.yaml"
	fd.Setup()

	respMap := map[string]*FakeResp{
		"lvcreate": {"", nil},
	}
	fd.cli.RootExecuter = NewFakeExecuter(respMap)
	fd.cli.BaseExecuter = NewFakeExecuter(respMap)

	opt := &pb.CreateFileShareOpts{
		Id:          "e1bb066c-5ce7-46eb-9336-25508cee9f71",
		Name:        "test001",
		Description: "fileshare for testing",
		Size:        int64(1),
		PoolName:    "vg001",
	}
	var expected = &model.FileShareSpec{
		BaseModel:   &model.BaseModel{},
		Name:        "test001",
		Description: "fileshare for testing",
		Size:        int64(1),
		Metadata: map[string]string{
			"lvPath": "/dev/vg001/test001",
		},
	}
	fileshare, err := fd.CreateFileShare(opt)
	if err != nil {
		t.Error("Failed to create fileshare:", err)
	}
	fileshare.Id = ""
	if !reflect.DeepEqual(fileshare, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, fileshare)
	}
}

func TestListPools(t *testing.T) {
	var fd = &Driver{}
	config.CONF.OsdsDock.Backends.NFS.ConfigPath = "testdata/nfs.yaml"
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
			StorageType:      "file",
			MultiAttach:      true,
			Extras: model.StoragePoolExtraSpec{
				DataStorage: model.DataStorageLoS{
					ProvisioningPolicy:      "Thin",
					IsSpaceEfficient:        false,
					StorageAccessCapability: []string{"Read", "Write", "Execute"},
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "nfs",
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
