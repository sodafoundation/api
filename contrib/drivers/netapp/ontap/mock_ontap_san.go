// Copyright 2020 The OpenSDS Authors.
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
	"fmt"
	log "github.com/golang/glog"
	tridentconfig "github.com/netapp/trident/config"
	"github.com/netapp/trident/storage"
	fakestorage "github.com/netapp/trident/storage/fake"
	sa "github.com/netapp/trident/storage_attribute"
	drivers "github.com/netapp/trident/storage_drivers"
	"github.com/netapp/trident/storage_drivers/fake"
	tu "github.com/netapp/trident/storage_drivers/fake/test_utils"
	"github.com/netapp/trident/utils"
	"github.com/stretchr/testify/mock"
	"sync"
)

// SANStorageDriverMock
type SANStorageDriverMock struct {
	mock.Mock
	storageDriver *fake.StorageDriver
}

var instance *SANStorageDriverMock
var once sync.Once

func GetSANStorageDriverMock() *SANStorageDriverMock {
	once.Do(func() {
		var err error
		instance, err = NewSANStorageDriverMock()
		if err != nil {
			log.Fatalf("Unable to create mock driver.")
		}
		log.Infof("Mock Driver {%+v} created successfully ", instance)

	})
	return instance
}

func NewSANStorageDriverMock() (*SANStorageDriverMock, error) {
	volumes := make([]fakestorage.Volume, 0)
	fakeConfig, err := fake.NewFakeStorageDriverConfigJSON("test", tridentconfig.Block,
		tu.GenerateFakePools(1), volumes)

	// Parse the common config struct from JSON
	commonConfig, err := drivers.ValidateCommonSettings(fakeConfig)
	if err != nil {
		err = fmt.Errorf("input failed validation: %v", err)
		return nil, err
	}

	var md = &SANStorageDriverMock{}

	if initializeErr := md.Initialize(
		tridentconfig.CurrentDriverContext, fakeConfig, commonConfig); initializeErr != nil {
		err = fmt.Errorf("problem initializing storage driver '%s': %v",
			commonConfig.StorageDriverName, initializeErr)
		return nil, err
	}
	log.Infof("storage driver (%s) initialized successfully.", commonConfig.StorageDriverName)
	return md, nil
}

func (md *SANStorageDriverMock) Initialize(context tridentconfig.DriverContext, configJSON string,
	commonConfig *drivers.CommonStorageDriverConfig) error {

	storageDriver := &fake.StorageDriver{}

	if err := storageDriver.Initialize(
		tridentconfig.CurrentDriverContext, configJSON, commonConfig); err != nil {
		err = fmt.Errorf("problem initializing storage driver '%s': %v",
			commonConfig.StorageDriverName, err)
		return err
	}
	md.storageDriver = storageDriver
	return nil
}

func (md *SANStorageDriverMock) Create(volConfig *storage.VolumeConfig, storagePool *storage.Pool,
	volAttributes map[string]sa.Request) error {
	if err := md.storageDriver.Create(volConfig, storagePool, volAttributes); err != nil {
		err = fmt.Errorf("problem creating volume '%s': %v",
			volConfig.Name, err)
		return err
	}
	return nil
}

func (md *SANStorageDriverMock) CreateClone(volConfig *storage.VolumeConfig,
	storagePool *storage.Pool) error {

	md.storageDriver.CreatePrepare(volConfig)

	md.storageDriver.Volumes[volConfig.CloneSourceVolume] = fakestorage.Volume{
		Name:         volConfig.CloneSourceVolume,
		PhysicalPool: "pool-0",
	}
	if err := md.storageDriver.CreateClone(volConfig, storagePool); err != nil {
		err = fmt.Errorf("problem creating clone volume '%s': %v",
			volConfig.Name, err)
		return err
	}
	return nil
}

func (md *SANStorageDriverMock) Destroy(name string) error {

	if err := md.storageDriver.Destroy(name); err != nil {
		err = fmt.Errorf("problem deleting volume '%s': %v",
			name, err)
		return err
	}
	return nil
}

func (md *SANStorageDriverMock) Resize(volConfig *storage.VolumeConfig, sizeBytes uint64) error {

	if err := md.storageDriver.Resize(volConfig, sizeBytes); err != nil {
		err = fmt.Errorf("problem resizing volume '%s': %v",
			volConfig.Name, err)
		return err
	}
	return nil
}

func (md *SANStorageDriverMock) Publish(name string, publishInfo *utils.VolumePublishInfo) error {

	if err := md.storageDriver.Publish(name, publishInfo); err != nil {
		err = fmt.Errorf("problem publishing connection '%s': %v",
			name, err)
		return err
	}
	return nil
}

func (md *SANStorageDriverMock) CreateSnapshot(snapConfig *storage.SnapshotConfig) (*storage.Snapshot, error) {

	md.storageDriver.Volumes[snapConfig.VolumeInternalName] = fakestorage.Volume{
		Name: snapConfig.VolumeInternalName,
	}
	snapshot, err := md.storageDriver.CreateSnapshot(snapConfig)
	if err != nil {
		err = fmt.Errorf("problem creating volume snapshot '%s': %v",
			snapConfig.Name, err)
		return nil, err
	}
	return snapshot, nil
}

func (md *SANStorageDriverMock) DeleteSnapshot(snapConfig *storage.SnapshotConfig) error {

	if err := md.storageDriver.DeleteSnapshot(snapConfig); err != nil {
		err = fmt.Errorf("problem deleting volume snapshot '%s': %v",
			snapConfig.Name, err)
		return err
	}
	return nil
}

func (md *SANStorageDriverMock) isTerminated() bool {
	return !md.storageDriver.Initialized()
}
