// Copyright 2019 The OpenSDS Authors.
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

package ontap

import (
	"github.com/netapp/trident/storage"
	"github.com/netapp/trident/storage_drivers/ontap/api"
	"github.com/netapp/trident/storage_drivers/ontap/api/azgo"
	"github.com/netapp/trident/utils"
	"github.com/opensds/opensds/pkg/model"
	"reflect"
	"testing"

	tridentconfig "github.com/netapp/trident/config"
	sa "github.com/netapp/trident/storage_attribute"
	drivers "github.com/netapp/trident/storage_drivers"
	odu "github.com/opensds/opensds/contrib/drivers/utils"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils/config"
)

var md *SANStorageDriverMock

func SetUpSANDriver(t *testing.T) {
	config.CONF.OsdsDock.Backends.NetappOntapSan.ConfigPath = "testdata/netapp_ontap_san.yaml"
	// Save current function and restore at the end:
	old := initialize
	defer func() { initialize = old }()

	initialize = func(context tridentconfig.DriverContext, configJSON string, commonConfig *drivers.CommonStorageDriverConfig) error {
		// This will be called, do whatever you want to,
		// return whatever you want to
		md = GetSANStorageDriverMock()
		return nil
	}
	if err := d.Setup(); err != nil {
		t.Errorf("Setup netapp ontap driver failed: %+v\n", err)
	}
}

func GetSANDriver(t *testing.T) *SANDriver {
	var d = &SANDriver{}
	SetUpSANDriver(t)
	return d
}

func TestSetup(t *testing.T) {
	var d = &SANDriver{}
	config.CONF.OsdsDock.Backends.NetappOntapSan.ConfigPath = "testdata/netapp_ontap_san.yaml"

	expectedPool := map[string]PoolProperties{
		"pool-0": {
			StorageType:      "block",
			AvailabilityZone: "default",
			MultiAttach:      true,
			Extras: model.StoragePoolExtraSpec{
				DataStorage: model.DataStorageLoS{
					ProvisioningPolicy: "Thin",
					Compression:        false,
					Deduplication:      false,
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "iscsi",
				},
			},
		},
	}
	expectedBackend := BackendOptions{
		Version:           1,
		StorageDriverName: "ontap-san",
		ManagementLIF:     "127.0.0.1",
		DataLIF:           "127.0.0.1",
		IgroupName:        "opensds",
		Svm:               "vserver",
		Username:          "admin",
		Password:          "password",
	}
	expectedDriver := &SANDriver{
		conf: &ONTAPConfig{
			BackendOptions: expectedBackend,
			Pool:           expectedPool,
		},
	}

	old := initialize
	defer func() { initialize = old }()

	initialize = func(context tridentconfig.DriverContext, configJSON string, commonConfig *drivers.CommonStorageDriverConfig) error {
		md, err := NewSANStorageDriverMock()
		if err != nil {
			t.Fatalf("Unable to create mock driver.")
		}
		t.Logf("Mock Driver {%+v} created successfully ", md)
		return nil
	}

	if err := d.Setup(); err != nil {
		t.Errorf("Setup netapp ontap driver failed: %+v\n", err)
	}

	if !reflect.DeepEqual(d.conf, expectedDriver.conf) {
		t.Errorf("Expected %+v, got %+v", expectedDriver.conf, d.conf)
	}
	t.Logf("Expected %+v, got %+v", expectedDriver.conf, d.conf)
}

func TestCreateVolume(t *testing.T) {
	var d = GetSANDriver(t)

	opt := &pb.CreateVolumeOpts{
		Id:          "e1bb066c-5ce7-46eb-9336-25508cee9f71",
		Name:        "testOntapVol1",
		Description: "volume for testing netapp ontap",
		Size:        int64(1),
		PoolName:    "pool-0",
	}
	var expected = &model.VolumeSpec{
		BaseModel:   &model.BaseModel{Id: "e1bb066c-5ce7-46eb-9336-25508cee9f71"},
		Name:        "testOntapVol1",
		Description: "volume for testing netapp ontap",
		Size:        int64(1),
		Identifier:  &model.Identifier{DurableName: "60a98000486e542d4f5a2f47694d684b", DurableNameFormat: "NAA"},
		Metadata: map[string]string{
			"lunPath": "/vol/opensds_e1bb066c5ce746eb933625508cee9f71/lun0",
		},
	}

	old := create
	old1 := LunGetSerialNumber
	defer func() { create = old; LunGetSerialNumber = old1 }()

	var serialNumber = "HnT-OZ/GiMhK"
	LunGetSerialNumber = func(lunPath string) (*azgo.LunGetSerialNumberResponse, error) {
		return &azgo.LunGetSerialNumberResponse{
			Result: azgo.LunGetSerialNumberResponseResult{
				SerialNumberPtr: &serialNumber,
			},
		}, nil
	}

	create = func(volConfig *storage.VolumeConfig, storagePool *storage.Pool, volAttributes map[string]sa.Request) error {
		var name = getVolumeName(opt.GetId())
		volConfig = d.GetVolumeConfig(name, opt.GetSize())
		storagePool = &storage.Pool{
			Name:               opt.GetPoolName(),
			StorageClasses:     make([]string, 0),
			Attributes:         make(map[string]sa.Offer),
			InternalAttributes: make(map[string]string),
		}
		err := md.Create(volConfig, storagePool, make(map[string]sa.Request))
		if err != nil {
			t.Error(err)
		}
		return nil
	}

	vol, err := d.CreateVolume(opt)
	if err != nil {
		t.Error("Failed to create volume:", err)
	}

	if !reflect.DeepEqual(vol, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, vol)
	}
	t.Logf("Expected %+v, got %+v\n", expected, vol)
}

func TestCreateVolumeFromSnapshot(t *testing.T) {
	var d = GetSANDriver(t)

	opt := &pb.CreateVolumeOpts{
		Id:          "e1bb066c-5ce7-46eb-9336-25508cee9f72",
		Name:        "testOntapVol1",
		Description: "volume for testing netapp ontap",
		Size:        int64(1),
		PoolName:    "pool-0",
		SnapshotId:  "3769855c-a102-11e7-b772-17b880d2f537",
	}
	var expected = &model.VolumeSpec{
		BaseModel:   &model.BaseModel{Id: "e1bb066c-5ce7-46eb-9336-25508cee9f72"},
		Name:        "testOntapVol1",
		Description: "volume for testing netapp ontap",
		Size:        int64(1),
		SnapshotId:  "3769855c-a102-11e7-b772-17b880d2f537",
		Identifier:  &model.Identifier{DurableName: "60a98000486e542d4f5a2f47694d684b", DurableNameFormat: "NAA"},
		Metadata: map[string]string{
			"lunPath": "/vol/opensds_e1bb066c5ce746eb933625508cee9f72/lun0",
		},
	}

	old := createClone
	old1 := LunGetSerialNumber
	defer func() { createClone = old; LunGetSerialNumber = old1 }()

	var serialNumber = "HnT-OZ/GiMhK"
	LunGetSerialNumber = func(lunPath string) (*azgo.LunGetSerialNumberResponse, error) {
		return &azgo.LunGetSerialNumberResponse{
			Result: azgo.LunGetSerialNumberResponseResult{
				SerialNumberPtr: &serialNumber,
			},
		}, nil
	}

	createClone = func(volConfig *storage.VolumeConfig, storagePool *storage.Pool) error {
		var snapName = getSnapshotName(opt.GetSnapshotId())
		var volName = opt.GetMetadata()["volume"]
		var name = getVolumeName(opt.GetId())

		volConfig = d.GetVolumeConfig(name, opt.GetSize())
		volConfig.CloneSourceVolumeInternal = volName
		volConfig.CloneSourceVolume = volName
		volConfig.CloneSourceSnapshot = snapName

		storagePool = &storage.Pool{
			Name:               opt.GetPoolName(),
			StorageClasses:     make([]string, 0),
			Attributes:         make(map[string]sa.Offer),
			InternalAttributes: make(map[string]string),
		}
		if err := md.CreateClone(volConfig, storagePool); err != nil {
			t.Error(err)
		}
		return nil
	}

	vol, err := d.CreateVolume(opt)
	if err != nil {
		t.Error("Failed to create volume:", err)
	}

	if !reflect.DeepEqual(vol, expected) {
		t.Errorf("Expected %+v, Got %+v\n", expected, vol)
	}
	t.Logf("Expected %+v, Got %+v\n", expected, vol)
}

func TestPullVolume(t *testing.T) {
	var d = GetSANDriver(t)
	volId := "e1bb066c-5ce7-46eb-9336-25508cee9f71"
	vol, err := d.PullVolume(volId)
	t.Logf("Expected Error %+v\n", err)
	expectedErr := &model.NotImplementError{"method PullVolume has not been implemented yet"}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Error("Failed to pull volume:", err)
	}
	if vol != nil {
		t.Errorf("Expected %+v, Got %+v\n", nil, vol)
	}
}

func TestDeleteVolume(t *testing.T) {
	var d = GetSANDriver(t)
	opt := &pb.DeleteVolumeOpts{
		Id: "e1bb066c-5ce7-46eb-9336-25508cee9f71",
	}

	old := destroy
	defer func() { destroy = old }()
	destroy = func(name string) error {
		name = getVolumeName(opt.GetId())
		err := md.Destroy(name)
		if err != nil {
			t.Error(err)
		}
		return nil
	}
	if err := d.DeleteVolume(opt); err != nil {
		t.Error("Failed to delete volume:", err)
	}
	t.Logf("volume (%s) deleted successfully.", opt.GetId())
}

func TestExtendVolume(t *testing.T) {
	var d = GetSANDriver(t)

	opt := &pb.ExtendVolumeOpts{
		Id:   "591c43e6-1156-42f5-9fbc-161153da185c",
		Size: int64(2),
	}

	old := resize
	defer func() { resize = old }()
	resize = func(volConfig *storage.VolumeConfig, newSize uint64) error {
		var name = getVolumeName(opt.GetId())
		volConfig = d.GetVolumeConfig(name, opt.GetSize())

		newSize = uint64(opt.GetSize() * bytesGiB)
		if err := md.Resize(volConfig, newSize); err != nil {
			t.Error(err)
		}
		return nil
	}

	vol, err := d.ExtendVolume(opt)
	if err != nil {
		t.Error("Failed to extend volume:", err)
	}

	if vol.Size != 2 {
		t.Errorf("Expected %+v, Got %+v\n", 2, vol.Size)
	}
	t.Logf("Got Extended Volume %+v\n", vol)
}

func TestInitializeConnection(t *testing.T) {
	var d = GetSANDriver(t)
	initiators := []*pb.Initiator{}
	initiator := &pb.Initiator{
		PortName: "iqn.2020-01.io.opensds:example",
		Protocol: "iscsi",
	}
	initiators = append(initiators, initiator)

	opt := &pb.CreateVolumeAttachmentOpts{
		Id:       "591c43e6-1156-42f5-9fbc-161153da185c",
		VolumeId: "e1bb066c-5ce7-46eb-9336-25508cee9f71",
		HostInfo: &pb.HostInfo{
			OsType:     "linux",
			Host:       "localhost",
			Ip:         "127.0.0.1",
			Initiators: initiators,
		},
		Metadata:       nil,
		DriverName:     "netapp_ontap_san",
		AccessProtocol: "iscsi",
	}

	expected := &model.ConnectionInfo{
		DriverVolumeType: "iscsi",
		ConnectionData: map[string]interface{}{
			"target_discovered": true,
			"volumeId":          "e1bb066c-5ce7-46eb-9336-25508cee9f71",
			"volume":            "opensds_e1bb066c5ce746eb933625508cee9f71",
			"description":       "NetApp ONTAP Attachment",
			"hostName":          "localhost",
			"initiator":         "iqn.2020-01.io.opensds:example",
			"targetIQN":         []string{""},
			"targetPortal":      []string{"127.0.0.1" + ":3260"},
			"targetLun":         int32(0),
			"igroup":            "",
		},
	}

	old := publish
	defer func() { publish = old }()
	publish = func(name string, publishInfo *utils.VolumePublishInfo) error {
		name = getVolumeName(opt.GetVolumeId())
		hostInfo := opt.GetHostInfo()
		initiator := odu.GetInitiatorName(hostInfo.GetInitiators(), opt.GetAccessProtocol())
		hostName := hostInfo.GetHost()
		publishInfo = &utils.VolumePublishInfo{
			HostIQN:  []string{initiator},
			HostIP:   []string{hostInfo.GetIp()},
			HostName: hostName,
		}

		md.Publish(name, publishInfo)
		return nil
	}

	connectionInfo, err := d.InitializeConnection(opt)
	if err != nil {
		t.Error("Failed to initialize connection:", err)
	}

	if !reflect.DeepEqual(connectionInfo, expected) {
		t.Errorf("Expected %+v, Got %+v\n", expected, connectionInfo)
	}
	t.Logf("Expected %+v, Got %+v\n", expected, connectionInfo)
}

func TestPullSnapshot(t *testing.T) {
	var d = GetSANDriver(t)
	snapshotId := "d1916c49-3088-4a40-b6fb-0fda18d074c3"
	snapshot, err := d.PullSnapshot(snapshotId)
	t.Logf("Expected Error %+v\n", err)
	expectedErr := &model.NotImplementError{"method PullSnapshot has not been implemented yet"}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Error("Failed to pull snapshot:", err)
	}
	if snapshot != nil {
		t.Errorf("Expected %+v, Got %+v\n", nil, snapshot)
	}
}

func TestInitializeSnapshotConnection(t *testing.T) {
	var d = GetSANDriver(t)
	opt := &pb.CreateSnapshotAttachmentOpts{}
	connectionInfo, err := d.InitializeSnapshotConnection(opt)
	t.Logf("Expected Error %+v\n", err)
	expectedErr := &model.NotImplementError{"method InitializeSnapshotConnection has not been implemented yet."}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Error("Failed to initialize snapshot connection:", err)
	}
	if connectionInfo != nil {
		t.Errorf("Expected %+v, Got %+v\n", nil, connectionInfo)
	}
}

func TestTerminateSnapshotConnection(t *testing.T) {
	var d = GetSANDriver(t)
	opt := &pb.DeleteSnapshotAttachmentOpts{}
	err := d.TerminateSnapshotConnection(opt)
	t.Logf("Expected Error %+v\n", err)
	expectedErr := &model.NotImplementError{"method TerminateSnapshotConnection has not been implemented yet."}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Error("Failed to terminate snapshot connection:", err)
	}
}

func TestCreateVolumeGroup(t *testing.T) {
	var d = GetSANDriver(t)
	opt := &pb.CreateVolumeGroupOpts{}
	vg, err := d.CreateVolumeGroup(opt)
	t.Logf("Expected Error %+v\n", err)
	expectedErr := &model.NotImplementError{"method CreateVolumeGroup has not been implemented yet"}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Error("Failed to create volume group:", err)
	}
	if vg != nil {
		t.Errorf("Expected %+v, Got %+v\n", nil, vg)
	}
}

func TestUpdateVolumeGroup(t *testing.T) {
	var d = GetSANDriver(t)
	opt := &pb.UpdateVolumeGroupOpts{}
	vg, err := d.UpdateVolumeGroup(opt)
	t.Logf("Expected Error %+v\n", err)
	expectedErr := &model.NotImplementError{"method UpdateVolumeGroup has not been implemented yet"}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Error("Failed to update volume group:", err)
	}
	if vg != nil {
		t.Errorf("Expected %+v, Got %+v\n", nil, vg)
	}
}

func TestDeleteVolumeGroup(t *testing.T) {
	var d = GetSANDriver(t)
	opt := &pb.DeleteVolumeGroupOpts{}
	err := d.DeleteVolumeGroup(opt)
	t.Logf("Expected Error %+v\n", err)
	expectedErr := &model.NotImplementError{"method DeleteVolumeGroup has not been implemented yet"}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Error("Failed to delete volume group:", err)
	}
}

func TestCreateSnapshot(t *testing.T) {
	var d = GetSANDriver(t)

	opt := &pb.CreateVolumeSnapshotOpts{
		Id:          "d1916c49-3088-4a40-b6fb-0fda18d074c3",
		Name:        "snap001",
		Description: "volume snapshot for Netapp ontap testing",
		Size:        int64(1),
		VolumeId:    "e1bb066c-5ce7-46eb-9336-25508cee9f71",
	}
	var expected = &model.VolumeSnapshotSpec{
		BaseModel:   &model.BaseModel{Id: "d1916c49-3088-4a40-b6fb-0fda18d074c3"},
		Name:        "snap001",
		Description: "volume snapshot for Netapp ontap testing",
		Size:        int64(1),
		VolumeId:    "e1bb066c-5ce7-46eb-9336-25508cee9f71",
		Metadata: map[string]string{
			"name":         "opensds_snapshot_d1916c4930884a40b6fb0fda18d074c3",
			"volume":       "opensds_e1bb066c5ce746eb933625508cee9f71",
			"creationTime": "2020-01-29T09:05:18Z",
			"size":         "0K",
		},
	}

	old := createSnapshot
	defer func() { createSnapshot = old }()
	createSnapshot = func(snapConfig *storage.SnapshotConfig) (*storage.Snapshot, error) {
		var snapName = getSnapshotName(opt.GetId())
		var volName = getVolumeName(opt.GetVolumeId())

		snapConfig = d.GetSnapshotConfig(snapName, volName)

		snapshot, err := md.CreateSnapshot(snapConfig)
		if err != nil {
			t.Error(err)
		}
		snapshot.Created = "2020-01-29T09:05:18Z"
		return snapshot, nil
	}

	snapshot, err := d.CreateSnapshot(opt)
	if err != nil {
		t.Error("Failed to create volume snapshot:", err)
	}

	if !reflect.DeepEqual(snapshot, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected, snapshot)
	}
	t.Logf("Expected %+v, Got %+v\n", expected, snapshot)
}

func TestDeleteSnapshot(t *testing.T) {
	var d = GetSANDriver(t)
	opt := &pb.DeleteVolumeSnapshotOpts{
		Id:       "d1916c49-3088-4a40-b6fb-0fda18d074c3",
		VolumeId: "e1bb066c-5ce7-46eb-9336-25508cee9f71",
	}

	old := deleteSnapshot
	defer func() { deleteSnapshot = old }()
	deleteSnapshot = func(snapConfig *storage.SnapshotConfig) error {
		var snapName = getSnapshotName(opt.GetId())
		var volName = getVolumeName(opt.GetVolumeId())

		snapConfig = d.GetSnapshotConfig(snapName, volName)
		if err := md.DeleteSnapshot(snapConfig); err != nil {
			t.Error(err)
		}
		return nil
	}
	if err := d.DeleteSnapshot(opt); err != nil {
		t.Error("Failed to delete snapshot:", err)
	}
	t.Logf("volume snapshot (%s) deleted successfully.", opt.GetId())
}

func TestListPools(t *testing.T) {
	var d = GetSANDriver(t)

	var expected = []*model.StoragePoolSpec{
		{
			BaseModel:        &model.BaseModel{},
			Name:             "pool-0",
			TotalCapacity:    int64(10),
			FreeCapacity:     int64(10),
			AvailabilityZone: "default",
			StorageType:      "block",
			MultiAttach:      true,
			Extras: model.StoragePoolExtraSpec{
				DataStorage: model.DataStorageLoS{
					ProvisioningPolicy: "Thin",
					Compression:        false,
					Deduplication:      false,
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "iscsi",
				},
			},
		},
	}

	old := VserverGetAggregateNames
	old1 := AggregateCommitment
	defer func() { VserverGetAggregateNames = old; AggregateCommitment = old1 }()
	VserverGetAggregateNames = func() ([]string, error) {
		return []string{"pool-0"}, nil
	}
	AggregateCommitment = func(aggregate string) (*api.AggregateCommitment, error) {
		return &api.AggregateCommitment{
			TotalAllocated: 0.0,
			AggregateSize:  10737418240.0,
		}, nil
	}
	fakepool := map[string]PoolProperties{
		"pool-0": {
			StorageType:      "block",
			AvailabilityZone: "default",
			MultiAttach:      true,
			Extras: model.StoragePoolExtraSpec{
				DataStorage: model.DataStorageLoS{
					ProvisioningPolicy: "Thin",
					Compression:        false,
					Deduplication:      false,
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "iscsi",
				},
			},
		},
	}
	d.conf = &ONTAPConfig{
		Pool: fakepool,
	}
	pools, err := d.ListPools()
	if err != nil {
		t.Error("Failed to list pools:", err)
	}
	for i := range pools {
		pools[i].Id = ""
	}

	if !reflect.DeepEqual(pools, expected) {
		t.Errorf("Expected %+v, got %+v\n", expected[0], pools[0])
	}
}

func TestUnSet(t *testing.T) {
	var d = GetSANDriver(t)

	old := Terminate
	defer func() { Terminate = old }()
	Terminate = func() {
		md.storageDriver.Terminate()
	}
	if err := d.Unset(); err != nil {
		t.Errorf("Unset netapp ontap driver failed: %+v\n", err)
	}
	isIntialized := md.isTerminated()
	if isIntialized == false {
		t.Errorf("Expected %+v, got %+v\n", false, isIntialized)
	}
}
