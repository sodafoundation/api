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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	log "github.com/golang/glog"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"

	"github.com/netapp/trident/storage"
	sa "github.com/netapp/trident/storage_attribute"
	drivers "github.com/netapp/trident/storage_drivers"
	"github.com/netapp/trident/storage_drivers/ontap"
	"github.com/netapp/trident/storage_drivers/ontap/api"
	"github.com/netapp/trident/utils"

	odu "github.com/opensds/opensds/contrib/drivers/utils"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils/config"
)

func lunPath(name string) string {
	return fmt.Sprintf("/vol/%v/lun0", name)
}

func getVolumeName(id string) string {
	r := strings.NewReplacer("-", "")
	return volumePrefix + r.Replace(id)
}

func getSnapshotName(id string) string {
	r := strings.NewReplacer("-", "")
	return snapshotPrefix + r.Replace(id)
}

// Get LUN Serial Number
func (d *Driver) getLunSerialNumber(lunPath string) (string, error) {

	lunSrNumber, err := d.sanStorageDriver.API.LunGetSerialNumber(lunPath)
	if err != nil {
		return "", fmt.Errorf("problem reading maps for LUN %s: %v", lunPath, err)
	}

	return naaPrefix + hex.EncodeToString([]byte(lunSrNumber.Result.SerialNumber())), nil
}

func (d *Driver) GetVolumeConfig(name string, size int64) (volConfig *storage.VolumeConfig) {
	volConfig = &storage.VolumeConfig{
		Version:      VolumeVersion,
		Name:         name,
		InternalName: name,
		Size:         strconv.FormatInt(size*bytesGB, 10),
		Protocol:     d.sanStorageDriver.GetProtocol(),
		AccessMode:   accessMode,
		VolumeMode:   volumeMode,
		AccessInfo:   utils.VolumeAccessInfo{},
	}
	return volConfig
}

func (d *Driver) GetSnapshotConfig(snapName string, volName string) (snapConfig *storage.SnapshotConfig) {
	snapConfig = &storage.SnapshotConfig{
		Version:            SnapshotVersion,
		Name:               snapName,
		InternalName:       snapName,
		VolumeName:         volName,
		VolumeInternalName: volName,
	}
	return snapConfig
}

func (d *Driver) Setup() error {
	// Read NetApp ONTAP config file
	d.conf = &ONTAPConfig{}

	p := config.CONF.OsdsDock.Backends.NetappOntapSan.ConfigPath
	if "" == p {
		p = defaultConfPath
	}
	if _, err := Parse(d.conf, p); err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("unable to instantiate ontap backend.")
		}
	}()

	empty := ""
	config := &drivers.OntapStorageDriverConfig{
		CommonStorageDriverConfig: &drivers.CommonStorageDriverConfig{
			Version:           d.conf.Version,
			StorageDriverName: StorageDriverName,
			StoragePrefixRaw:  json.RawMessage("{}"),
			StoragePrefix:     &empty,
		},
		ManagementLIF: d.conf.ManagementLIF,
		DataLIF:       d.conf.DataLIF,
		IgroupName:    d.conf.IgroupName,
		SVM:           d.conf.Svm,
		Username:      d.conf.Username,
		Password:      d.conf.Password,
	}
	marshaledJSON, err := json.Marshal(config)
	if err != nil {
		log.Fatal("unable to marshal ONTAP config:  ", err)
	}
	configJSON := string(marshaledJSON)

	// Convert config (JSON or YAML) to JSON
	configJSONBytes, err := yaml.YAMLToJSON([]byte(configJSON))
	if err != nil {
		err = fmt.Errorf("invalid config format: %v", err)
		return err
	}
	configJSON = string(configJSONBytes)

	// Parse the common config struct from JSON
	commonConfig, err := drivers.ValidateCommonSettings(configJSON)
	if err != nil {
		err = fmt.Errorf("input failed validation: %v", err)
		return err
	}

	d.sanStorageDriver = &ontap.SANStorageDriver{
		Config: *config,
	}

	// Initialize the driver.
	if err = d.sanStorageDriver.Initialize(driverContext, configJSON, commonConfig); err != nil {
		log.Errorf("could not initialize storage driver (%s). failed: %v", commonConfig.StorageDriverName, err)
		return err
	}
	log.Infof("storage driver (%s) initialized successfully.", commonConfig.StorageDriverName)

	return nil
}

func (d *Driver) Unset() error {
	//driver to clean up and stop any ongoing operations.
	d.sanStorageDriver.Terminate()
	return nil
}

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (vol *model.VolumeSpec, err error) {

	if opt.GetSnapshotId() != "" {
		return d.createVolumeFromSnapshot(opt)
	}

	var name = getVolumeName(opt.GetId())
	volConfig := d.GetVolumeConfig(name, opt.GetSize())

	storagePool := &storage.Pool{
		Name:               opt.GetPoolName(),
		StorageClasses:     make([]string, 0),
		Attributes:         make(map[string]sa.Offer),
		InternalAttributes: make(map[string]string),
	}

	err = d.sanStorageDriver.Create(volConfig, storagePool, make(map[string]sa.Request))
	if err != nil {
		log.Errorf("create volume (%s) failed: %v", opt.GetId(), err)
		return nil, err
	}

	lunPath := lunPath(name)

	// Get LUN Serial Number
	lunSerialNumber, err := d.getLunSerialNumber(lunPath)
	if err != nil {
		log.Errorf("create volume (%s) failed: %v", opt.GetId(), err)
		return nil, err
	}

	log.Infof("volume (%s) created successfully.", opt.GetId())

	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Size:        opt.GetSize(),
		Description: opt.GetDescription(),
		Identifier:  &model.Identifier{DurableName: lunSerialNumber, DurableNameFormat: KLvIdFormat},
		Metadata: map[string]string{
			KLvPath: lunPath,
		},
	}, nil
}

func (d *Driver) createVolumeFromSnapshot(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {

	var snapName = getSnapshotName(opt.GetSnapshotId())
	var volName = opt.GetMetadata()["volume"]
	var name = getVolumeName(opt.GetId())

	volConfig := d.GetVolumeConfig(name, opt.GetSize())
	volConfig.CloneSourceVolumeInternal = volName
	volConfig.CloneSourceSnapshot = volName
	volConfig.CloneSourceSnapshot = snapName

	err := d.sanStorageDriver.CreateClone(volConfig)
	if err != nil {
		log.Errorf("create volume (%s) from snapshot (%s) failed: %v", opt.GetId(), opt.GetSnapshotId(), err)
		return nil, err
	}

	lunPath := lunPath(name)

	// Get LUN Serial Number
	lunSerialNumber, err := d.getLunSerialNumber(lunPath)
	if err != nil {
		log.Errorf("create volume (%s) from snapshot (%s) failed: %v", opt.GetId(), opt.GetSnapshotId(), err)
		return nil, err
	}

	log.Infof("volume (%s) created from snapshot (%s) successfully.", opt.GetId(), opt.GetSnapshotId())

	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Size:        opt.GetSize(),
		Description: opt.GetDescription(),
		Identifier:  &model.Identifier{DurableName: lunSerialNumber, DurableNameFormat: KLvIdFormat},
		Metadata: map[string]string{
			KLvPath: lunPath,
		},
	}, err
}

func (d *Driver) PullVolume(volId string) (*model.VolumeSpec, error) {

	return nil, &model.NotImplementError{"method PullVolume has not been implemented yet"}
}

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	var name = getVolumeName(opt.GetId())
	err := d.sanStorageDriver.Destroy(name)
	if err != nil {
		msg := fmt.Sprintf("delete volume (%s) failed: %v", opt.GetId(), err)
		log.Error(msg)
		return err
	}
	log.Infof("volume (%s) deleted successfully.", opt.GetId())
	return nil
}

// ExtendVolume ...
func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	var name = getVolumeName(opt.GetId())
	volConfig := d.GetVolumeConfig(name, opt.GetSize())

	newSize := uint64(opt.GetSize() * bytesGB)
	if err := d.sanStorageDriver.Resize(volConfig, newSize); err != nil {
		log.Errorf("extend volume (%s) failed, error: %v", name, err)
		return nil, err
	}

	log.Infof("volume (%s) extended successfully.", opt.GetId())
	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Size:        opt.GetSize(),
		Description: opt.GetDescription(),
		Metadata:    opt.GetMetadata(),
	}, nil
}

func (d *Driver) InitializeConnection(opt *pb.CreateVolumeAttachmentOpts) (*model.ConnectionInfo, error) {

	var name = getVolumeName(opt.GetVolumeId())
	hostInfo := opt.GetHostInfo()
	initiator := odu.GetInitiatorName(hostInfo.GetInitiators(), opt.GetAccessProtocol())
	hostName := hostInfo.GetHost()

	publishInfo := &utils.VolumePublishInfo{
		HostIQN:  []string{initiator},
		HostIP:   []string{hostInfo.GetIp()},
		HostName: hostName,
	}

	err := d.sanStorageDriver.Publish(name, publishInfo)
	if err != nil {
		msg := fmt.Sprintf("volume (%s) attachment is failed: %v", opt.GetVolumeId(), err)
		log.Errorf(msg)
		return nil, err
	}

	log.Infof("volume (%s) attachment is created successfully", opt.GetVolumeId())

	connInfo := &model.ConnectionInfo{
		DriverVolumeType: opt.GetAccessProtocol(),
		ConnectionData: map[string]interface{}{
			"target_discovered": true,
			"volume_id":         opt.GetVolumeId(),
			"volume":            name,
			"description":       "NetApp ONTAP Attachment",
			"host":              hostName,
			"initiator":         initiator,
			"targetIQN":         []string{publishInfo.IscsiTargetIQN},
			"targetPortal":      []string{hostInfo.GetIp() + ":3260"},
			"targetLun":         publishInfo.IscsiLunNumber,
			"igroup":            publishInfo.IscsiIgroup,
		},
	}

	log.Infof("initialize connection successfully: %v", connInfo)
	return connInfo, nil
}

func (d *Driver) TerminateConnection(opt *pb.DeleteVolumeAttachmentOpts) error {
	var name = getVolumeName(opt.GetVolumeId())

	// Validate Flexvol exists before trying to Unmount
	volExists, err := d.sanStorageDriver.API.VolumeExists(name)
	if err != nil {
		return fmt.Errorf("error checking for existing volume (%s), error: %v", name, err)
	}
	if !volExists {
		log.Infof("volume %s already deleted, skipping destroy.", name)
		return nil
	}

	// Unmount the FlexVolume
	volUnmountResponse, err := d.sanStorageDriver.API.VolumeUnmount(name, true)
	if err != nil {
		return fmt.Errorf("error destroying volume %v: %v", name, err)
	}
	if zerr := api.NewZapiError(volUnmountResponse); !zerr.IsPassed() {
		return fmt.Errorf("error destroying volume %v: %v", name, zerr.Error())
	}

	log.Infof("termination connection successfully")
	return nil
}

func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (snap *model.VolumeSnapshotSpec, err error) {
	var snapName = getSnapshotName(opt.GetId())
	var volName = getVolumeName(opt.GetVolumeId())

	snapConfig := d.GetSnapshotConfig(snapName, volName)

	snapshot, err := d.sanStorageDriver.CreateSnapshot(snapConfig)

	if err != nil {
		msg := fmt.Sprintf("create snapshot %s (%s) failed: %s", opt.GetName(), opt.GetId(), err)
		log.Error(msg)
		return nil, err
	}

	log.Infof("snapshot %s (%s) created successfully.", opt.GetName(), opt.GetId())

	return &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Description: opt.GetDescription(),
		VolumeId:    opt.GetVolumeId(),
		Size:        opt.GetSize(),
		Metadata: map[string]string{
			"name":         snapName,
			"volume":       volName,
			"creationTime": snapshot.Created,
			"size":         strconv.FormatInt(snapshot.SizeBytes/bytesGB, 10) + "GB",
		},
	}, nil
}

func (d *Driver) PullSnapshot(snapIdentifier string) (*model.VolumeSnapshotSpec, error) {
	// not used, do nothing
	return nil, &model.NotImplementError{"method PullSnapshot has not been implemented yet"}
}

func (d *Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {

	var snapName = getSnapshotName(opt.GetId())
	var volName = getVolumeName(opt.GetVolumeId())

	snapConfig := d.GetSnapshotConfig(snapName, volName)

	err := d.sanStorageDriver.DeleteSnapshot(snapConfig)

	if err != nil {
		msg := fmt.Sprintf("delete volume snapshot (%s) failed: %v", opt.GetId(), err)
		log.Error(msg)
		return err
	}
	log.Infof("volume snapshot (%s) deleted successfully", opt.GetId())
	return nil
}

func (d *Driver) ListPools() ([]*model.StoragePoolSpec, error) {

	var pools []*model.StoragePoolSpec

	aggregates, err := d.sanStorageDriver.API.VserverGetAggregateNames()

	if err != nil {
		msg := fmt.Sprintf("list pools failed: %v", err)
		log.Error(msg)
		return nil, err
	}

	c := d.conf
	for _, aggr := range aggregates {
		if _, ok := c.Pool[aggr]; !ok {
			continue
		}
		aggregate, _ := d.sanStorageDriver.API.AggregateCommitment(aggr)
		aggregateCapacity := aggregate.AggregateSize / bytesGB
		aggregateAllocatedCapacity := aggregate.TotalAllocated / bytesGB

		pool := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, aggr).String(),
			},
			Name:             aggr,
			TotalCapacity:    int64(aggregateCapacity),
			FreeCapacity:     int64(aggregateCapacity) - int64(aggregateAllocatedCapacity),
			ConsumedCapacity: int64(aggregateAllocatedCapacity),
			StorageType:      c.Pool[aggr].StorageType,
			Extras:           c.Pool[aggr].Extras,
			AvailabilityZone: c.Pool[aggr].AvailabilityZone,
		}
		if pool.AvailabilityZone == "" {
			pool.AvailabilityZone = DefaultAZ
		}
		pools = append(pools, pool)
	}

	log.Info("list pools successfully")
	return pools, nil
}

func (d *Driver) InitializeSnapshotConnection(opt *pb.CreateSnapshotAttachmentOpts) (*model.ConnectionInfo, error) {

	return nil, &model.NotImplementError{S: "method InitializeSnapshotConnection has not been implemented yet."}
}

func (d *Driver) TerminateSnapshotConnection(opt *pb.DeleteSnapshotAttachmentOpts) error {

	return &model.NotImplementError{S: "method TerminateSnapshotConnection has not been implemented yet."}

}

func (d *Driver) CreateVolumeGroup(opt *pb.CreateVolumeGroupOpts) (*model.VolumeGroupSpec, error) {
	return nil, &model.NotImplementError{"method CreateVolumeGroup has not been implemented yet"}
}

func (d *Driver) UpdateVolumeGroup(opt *pb.UpdateVolumeGroupOpts) (*model.VolumeGroupSpec, error) {
	return nil, &model.NotImplementError{"method UpdateVolumeGroup has not been implemented yet"}
}

func (d *Driver) DeleteVolumeGroup(opt *pb.DeleteVolumeGroupOpts) error {
	return &model.NotImplementError{"method DeleteVolumeGroup has not been implemented yet"}
}
