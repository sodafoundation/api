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

package fusionstorage

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/opensds/opensds/contrib/drivers/utils/config"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
	c "github.com/opensds/opensds/pkg/utils/config"
	"github.com/opensds/opensds/pkg/utils/exec"
)

var startServer = `[2018-10-30 17:25:50] [INFO] [fsc_cli.start_dsware_api_daemon:59] INFO - Start FSCTool service successfully.
`

var queryAllPoolInfo = `
pool_id=0,total_capacity=7205408,alloc_capacity=11264,used_capacity=10330,pool_model=0
pool_id=1,total_capacity=7205408,alloc_capacity=0,used_capacity=1113,pool_model=0
result=0

`

var createVolume = `
wwn=688860300000000180d44b2538981926
result=0

`

var queryVolume = `
vol_name=vol002,father_name=test,status=0,vol_size=2048,real_size=-1,pool_id=0,create_time=1540543278,encrypt_flag=0,lun_id=5,lld_progress=-1,rw_right=true,wwn=688860300000000580d44b2538981926
result=0

`

var expandVolume = `result=0

`

var deleteVolume = `result=0

`

var createVolumeFromSnap = `
wwn=688860300000000580d44b2538981926
result=0

`
var createSnapshot = `
snap_name=snapshot,father_name=,status=0,snap_size=1024,real_size=-1,pool_id=0,delete_priority=0,create_time=1540471203,encrypt_flag=0,smartCacheFlag=0,tree_id=0,branch_id=0,snap_id=0
result=0

`
var querySnapshot = `
snap_name=snapshot,father_name=,status=0,snap_size=1024,real_size=-1,pool_id=0,delete_priority=0,create_time=1540471203,encrypt_flag=0,smartCacheFlag=0,tree_id=0,branch_id=0,snap_id=0
result=0

`

var deleteSnapshot = `
result=0

`

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
	v, ok := f.RespMap[args[1]]
	if !ok {
		fmt.Println(args)
		return "", fmt.Errorf("can find specified op: %s", args[1])
	}
	return v.out, v.err
}

func TestSetup(t *testing.T) {
	f := Driver{}
	c.CONF.OsdsDock.Backends.HuaweiFusionStorage.ConfigPath = "./testdata/fusionstorage.yaml"
	respMap := map[string]*FakeResp{
		"startServer": &FakeResp{startServer, nil},
	}
	baseExecuter = NewFakeExecuter(respMap)
	rootExecuter = NewFakeExecuter(respMap)
	f.Setup()
	expect := &Config{
		AuthOptions: AuthOptions{
			FmIp: "192.168.0.100",
			FsaIp: []string{
				"192.168.0.1",
				"192.168.0.2",
				"192.168.0.3",
			},
		},
		Pool: map[string]config.PoolProperties{
			"0": config.PoolProperties{
				StorageType:      "block",
				AvailabilityZone: "nova-01",
				Extras: model.StoragePoolExtraSpec{
					DataStorage: model.DataStorageLoS{
						ProvisioningPolicy: "Thin",
						IsSpaceEfficient:   false,
					},
					IOConnectivity: model.IOConnectivityLoS{
						AccessProtocol: "DSWARE",
						MaxIOPS:        7000000,
						MaxBWS:         600,
					},
					Advanced: map[string]interface{}{
						"diskType": "SSD",
						"latency":  "3ms",
					},
				},
			},
			"1": {
				StorageType:      "block",
				AvailabilityZone: "nova-01",
				Extras: model.StoragePoolExtraSpec{
					DataStorage: model.DataStorageLoS{
						ProvisioningPolicy: "Thin",
						IsSpaceEfficient:   false,
					},
					IOConnectivity: model.IOConnectivityLoS{
						AccessProtocol: "DSWARE",
						MaxIOPS:        3000000,
						MaxBWS:         300,
					},
					Advanced: map[string]interface{}{
						"diskType": "SSD",
						"latency":  "500ms",
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(expect, f.conf) {
		t.Errorf("Test driver setup failed:\n expect:%v\n got:\t%v", expect, f.conf)
	}
}

func TestCreateVolume(t *testing.T) {
	f := Driver{}
	c.CONF.OsdsDock.Backends.HuaweiFusionStorage.ConfigPath = "./testdata/fusionstorage.yaml"
	respMap := map[string]*FakeResp{
		"startServer":  &FakeResp{startServer, nil},
		"createVolume": &FakeResp{createVolume, nil},
	}
	baseExecuter = NewFakeExecuter(respMap)
	rootExecuter = NewFakeExecuter(respMap)
	f.Setup()

	opt := &pb.CreateVolumeOpts{
		Id:       "b4c29f4b-6ab8-40ed-be3e-0111bcea7b14",
		Name:     "FakeVolumeName",
		Size:     10,
		PoolId:   "49cc1071-18af-49f3-913e-e8b36370f32c",
		PoolName: "0",
	}
	resp, err := f.CreateVolume(opt)
	if err != nil {
		t.Errorf("Test CreateVolume failed: %v", resp)
	}
	expect := &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: "b4c29f4b-6ab8-40ed-be3e-0111bcea7b14",
		},
		Name:   "FakeVolumeName",
		Size:   10,
		PoolId: "49cc1071-18af-49f3-913e-e8b36370f32c",
	}
	if !reflect.DeepEqual(expect, resp) {
		t.Errorf("Test create volume failed:\n expect:%v \n got:\t%v", &expect, resp)
	}
}

func TestCreateVolumeFromSnapshot(t *testing.T) {
	f := Driver{}
	c.CONF.OsdsDock.Backends.HuaweiFusionStorage.ConfigPath = "./testdata/fusionstorage.yaml"
	respMap := map[string]*FakeResp{
		"startServer":          &FakeResp{startServer, nil},
		"createVolumeFromSnap": &FakeResp{createVolumeFromSnap, nil},
	}
	baseExecuter = NewFakeExecuter(respMap)
	rootExecuter = NewFakeExecuter(respMap)
	f.Setup()

	opt := &pb.CreateVolumeOpts{
		Id:         "b4c29f4b-6ab8-40ed-be3e-0111bcea7b14",
		Name:       "FakeVolumeName",
		Size:       10,
		PoolId:     "49cc1071-18af-49f3-913e-e8b36370f32c",
		PoolName:   "0",
		SnapshotId: "b84e7fc5-4feb-42d2-86e4-1a0a287e3fad",
	}
	resp, err := f.CreateVolume(opt)
	if err != nil {
		t.Errorf("Test CreateVolume failed: %v", resp)
	}
	expect := model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: "b4c29f4b-6ab8-40ed-be3e-0111bcea7b14",
		},
		Name:       "FakeVolumeName",
		Size:       10,
		PoolId:     "49cc1071-18af-49f3-913e-e8b36370f32c",
		SnapshotId: "b84e7fc5-4feb-42d2-86e4-1a0a287e3fad",
	}
	if !reflect.DeepEqual(&expect, resp) {
		t.Errorf("Test create volume failed:\n expect:%v \n got:\t%v", &expect, resp)
	}
}

func TestDeleteVolume(t *testing.T) {
	f := Driver{}
	c.CONF.OsdsDock.Backends.HuaweiFusionStorage.ConfigPath = "./testdata/fusionstorage.yaml"
	respMap := map[string]*FakeResp{
		"startServer":  &FakeResp{startServer, nil},
		"deleteVolume": &FakeResp{deleteVolume, nil},
	}
	baseExecuter = NewFakeExecuter(respMap)
	rootExecuter = NewFakeExecuter(respMap)
	f.Setup()

	opt := &pb.DeleteVolumeOpts{
		Id: "b4c29f4b-6ab8-40ed-be3e-0111bcea7b14",
	}
	if err := f.DeleteVolume(opt); err != nil {
		t.Errorf("Test Delete volume failed, %v", err)
	}

	respMap["deleteVolume"] = &FakeResp{deleteVolume, VolumeNotExist}
	baseExecuter = NewFakeExecuter(respMap)
	rootExecuter = NewFakeExecuter(respMap)
	if err := f.DeleteVolume(opt); err != nil {
		t.Errorf("Test Delete volume failed")
	}

	respMap["deleteVolume"] = &FakeResp{deleteVolume, fmt.Errorf("fake error")}
	baseExecuter = NewFakeExecuter(respMap)
	rootExecuter = NewFakeExecuter(respMap)
	if err := f.DeleteVolume(opt); err == nil {
		t.Errorf("Test Delete volume failed")
	}
}

func TestCreateSnapshot(t *testing.T) {
	f := Driver{}
	c.CONF.OsdsDock.Backends.HuaweiFusionStorage.ConfigPath = "./testdata/fusionstorage.yaml"
	respMap := map[string]*FakeResp{
		"startServer":    &FakeResp{startServer, nil},
		"createSnapshot": &FakeResp{createSnapshot, nil},
	}
	baseExecuter = NewFakeExecuter(respMap)
	rootExecuter = NewFakeExecuter(respMap)
	f.Setup()
	opt := &pb.CreateVolumeSnapshotOpts{
		Id:       "b4c29f4b-6ab8-40ed-be3e-0111bcea7b14",
		VolumeId: "b84e7fc5-4feb-42d2-86e4-1a0a287e3fad",
		Name:     "FakeVolumeName",
		Size:     10,
	}
	resp, err := f.CreateSnapshot(opt)
	if err != nil {
		t.Errorf("Test Create volume snapshot failed: %v", err)
	}
	expect := &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: "b4c29f4b-6ab8-40ed-be3e-0111bcea7b14",
		},
		VolumeId: "b84e7fc5-4feb-42d2-86e4-1a0a287e3fad",
		Name:     "FakeVolumeName",
		Size:     10,
	}

	if !reflect.DeepEqual(expect, resp) {
		t.Errorf("Test create volume failed:\n expect:%v \n got:\t%v", &expect, resp)
	}
}

func TestDeleteSnapshot(t *testing.T) {
	f := Driver{}
	c.CONF.OsdsDock.Backends.HuaweiFusionStorage.ConfigPath = "./testdata/fusionstorage.yaml"
	respMap := map[string]*FakeResp{
		"startServer":    &FakeResp{startServer, nil},
		"deleteSnapshot": &FakeResp{deleteSnapshot, nil},
	}
	baseExecuter = NewFakeExecuter(respMap)
	rootExecuter = NewFakeExecuter(respMap)
	f.Setup()

	opt := &pb.DeleteVolumeSnapshotOpts{
		Id:       "b4c29f4b-6ab8-40ed-be3e-0111bcea7b14",
		VolumeId: "b84e7fc5-4feb-42d2-86e4-1a0a287e3fad",
	}
	err := f.DeleteSnapshot(opt)
	if err != nil {
		t.Errorf("Test delete snapshot failed, %v", err)
	}

	respMap["deleteSnapshot"] = &FakeResp{deleteSnapshot, SnapshotNotExist}
	baseExecuter = NewFakeExecuter(respMap)
	rootExecuter = NewFakeExecuter(respMap)
	if err := f.DeleteSnapshot(opt); err != nil {
		t.Errorf("Test Delete snapshot failed")
	}

	respMap["deleteSnapshot"] = &FakeResp{deleteSnapshot, fmt.Errorf("fake error")}
	baseExecuter = NewFakeExecuter(respMap)
	rootExecuter = NewFakeExecuter(respMap)
	if err := f.DeleteSnapshot(opt); err == nil {
		t.Errorf("Test Delete snapshot failed")
	}
}

func TestListPool(t *testing.T) {
	f := Driver{}
	c.CONF.OsdsDock.Backends.HuaweiFusionStorage.ConfigPath = "./testdata/fusionstorage.yaml"
	respMap := map[string]*FakeResp{
		"startServer":      &FakeResp{startServer, nil},
		"queryAllPoolInfo": &FakeResp{queryAllPoolInfo, nil},
	}
	baseExecuter = NewFakeExecuter(respMap)
	rootExecuter = NewFakeExecuter(respMap)
	f.Setup()
	resp, err := f.ListPools()
	if err != nil {
		t.Errorf("List pool failed: %v", err)
	}
	expect := []*model.StoragePoolSpec{
		&model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: resp[0].Id,
			},
			Name:             "0",
			FreeCapacity:     7026,
			TotalCapacity:    7036,
			StorageType:      "block",
			AvailabilityZone: "nova-01",
			Extras: model.StoragePoolExtraSpec{
				DataStorage: model.DataStorageLoS{
					ProvisioningPolicy: "Thin",
					IsSpaceEfficient:   false,
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "DSWARE",
					MaxIOPS:        7000000,
					MaxBWS:         600,
				},
				Advanced: map[string]interface{}{
					"diskType": "SSD",
					"latency":  "3ms",
				},
			},
		},
		&model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: resp[1].Id,
			},
			Name:             "1",
			FreeCapacity:     7035,
			TotalCapacity:    7036,
			StorageType:      "block",
			AvailabilityZone: "nova-01",
			Extras: model.StoragePoolExtraSpec{
				DataStorage: model.DataStorageLoS{
					ProvisioningPolicy: "Thin",
					IsSpaceEfficient:   false,
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "DSWARE",
					MaxIOPS:        3000000,
					MaxBWS:         300,
				},
				Advanced: map[string]interface{}{
					"diskType": "SSD",
					"latency":  "500ms",
				},
			},
		},
	}
	if !reflect.DeepEqual(expect, resp) {
		t.Errorf("Test create volume failed:\n expect:%v %v \n got:\t%v %v",
			expect[0], expect[1], resp[0], resp[1])
	}
}
