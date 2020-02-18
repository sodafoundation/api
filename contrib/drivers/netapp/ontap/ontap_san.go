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
	"github.com/netapp/trident/storage_drivers/ontap/api/azgo"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	log "github.com/golang/glog"
	uuid "github.com/satori/go.uuid"

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

var (
	d                        = SANDriver{}
	initialize               = d.sanStorageDriver.Initialize
	create                   = d.sanStorageDriver.Create
	createClone              = d.sanStorageDriver.CreateClone
	destroy                  = d.sanStorageDriver.Destroy
	resize                   = d.sanStorageDriver.Resize
	publish                  = d.sanStorageDriver.Publish
	createSnapshot           = d.sanStorageDriver.CreateSnapshot
	deleteSnapshot           = d.sanStorageDriver.DeleteSnapshot
	LunGetSerialNumber       = func(lunPath string) (*azgo.LunGetSerialNumberResponse, error) { return nil, nil }
	VserverGetAggregateNames = func() ([]string, error) { return nil, nil }
	AggregateCommitment      = func(aggregate string) (*api.AggregateCommitment, error) { return nil, nil }
	Terminate                = d.sanStorageDriver.Terminate
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
func (d *SANDriver) getLunSerialNumber(lunPath string) (string, error) {

	lunSrNumber, err := LunGetSerialNumber(lunPath)
	if err != nil {
		return "", fmt.Errorf("problem reading maps for LUN %s: %v", lunPath, err)
	}

	return naaPrefix + hex.EncodeToString([]byte(lunSrNumber.Result.SerialNumber())), nil
}

func (d *SANDriver) GetVolumeConfig(name string, size int64) (volConfig *storage.VolumeConfig) {
	volConfig = &storage.VolumeConfig{
		Version:      VolumeVersion,
		Name:         name,
		InternalName: name,
		Size:         strconv.FormatInt(size*bytesGiB, 10),
		Protocol:     d.sanStorageDriver.GetProtocol(),
		AccessMode:   accessMode,
		VolumeMode:   volumeMode,
		AccessInfo:   utils.VolumeAccessInfo{},
	}
	return volConfig
}

func (d *SANDriver) GetSnapshotConfig(snapName string, volName string) (snapConfig *storage.SnapshotConfig) {
	snapConfig = &storage.SnapshotConfig{
		Version:            SnapshotVersion,
		Name:               snapName,
		InternalName:       snapName,
		VolumeName:         volName,
		VolumeInternalName: volName,
	}
	return snapConfig
}

func (d *SANDriver) Setup() error {
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
	if err = initialize(driverContext, configJSON, commonConfig); err != nil {
		log.Errorf("could not initialize storage driver (%s). failed: %v", commonConfig.StorageDriverName, err)
		return err
	}
	log.Infof("storage driver (%s) initialized successfully.", commonConfig.StorageDriverName)
	LunGetSerialNumber = d.sanStorageDriver.API.LunGetSerialNumber
	VserverGetAggregateNames = d.sanStorageDriver.API.VserverGetAggregateNames
	AggregateCommitment = d.sanStorageDriver.API.AggregateCommitment
	return nil
}

func (d *SANDriver) Unset() error {
	//driver to clean up and stop any ongoing operations.
	Terminate()
	return nil
}

func (d *SANDriver) CreateVolume(opt *pb.CreateVolumeOpts) (vol *model.VolumeSpec, err error) {

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

	err = create(volConfig, storagePool, make(map[string]sa.Request))
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

func (d *SANDriver) createVolumeFromSnapshot(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {

	var snapName = getSnapshotName(opt.GetSnapshotId())
	var volName = opt.GetMetadata()["volume"]
	var name = getVolumeName(opt.GetId())

	volConfig := d.GetVolumeConfig(name, opt.GetSize())
	volConfig.CloneSourceVolumeInternal = volName
	volConfig.CloneSourceVolume = volName
	volConfig.CloneSourceSnapshot = snapName

	storagePool := &storage.Pool{
		Name:               opt.GetPoolName(),
		StorageClasses:     make([]string, 0),
		Attributes:         make(map[string]sa.Offer),
		InternalAttributes: make(map[string]string),
	}

	err := createClone(volConfig, storagePool)
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
		SnapshotId:  opt.GetSnapshotId(),
		Identifier:  &model.Identifier{DurableName: lunSerialNumber, DurableNameFormat: KLvIdFormat},
		Metadata: map[string]string{
			KLvPath: lunPath,
		},
	}, err
}

func (d *SANDriver) PullVolume(volId string) (*model.VolumeSpec, error) {

	return nil, &model.NotImplementError{"method PullVolume has not been implemented yet"}
}

func (d *SANDriver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	var name = getVolumeName(opt.GetId())
	err := destroy(name)
	if err != nil {
		msg := fmt.Sprintf("delete volume (%s) failed: %v", opt.GetId(), err)
		log.Error(msg)
		return err
	}
	log.Infof("volume (%s) deleted successfully.", opt.GetId())
	return nil
}

// ExtendVolume ...
func (d *SANDriver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	var name = getVolumeName(opt.GetId())
	volConfig := d.GetVolumeConfig(name, opt.GetSize())

	newSize := uint64(opt.GetSize() * bytesGiB)
	if err := resize(volConfig, newSize); err != nil {
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

func (d *SANDriver) InitializeConnection(opt *pb.CreateVolumeAttachmentOpts) (*model.ConnectionInfo, error) {

	var name = getVolumeName(opt.GetVolumeId())
	hostInfo := opt.GetHostInfo()
	initiator := odu.GetInitiatorName(hostInfo.GetInitiators(), opt.GetAccessProtocol())
	hostName := hostInfo.GetHost()

	publishInfo := &utils.VolumePublishInfo{
		HostIQN:  []string{initiator},
		HostIP:   []string{hostInfo.GetIp()},
		HostName: hostName,
	}

	err := publish(name, publishInfo)
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
			"volumeId":          opt.GetVolumeId(),
			"volume":            name,
			"description":       "NetApp ONTAP Attachment",
			"hostName":          hostName,
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

func (d *SANDriver) TerminateConnection(opt *pb.DeleteVolumeAttachmentOpts) error {
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

func (d *SANDriver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (snap *model.VolumeSnapshotSpec, err error) {
	var snapName = getSnapshotName(opt.GetId())
	var volName = getVolumeName(opt.GetVolumeId())

	snapConfig := d.GetSnapshotConfig(snapName, volName)

	snapshot, err := createSnapshot(snapConfig)

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
			"size":         strconv.FormatInt(snapshot.SizeBytes/bytesKiB, 10) + "K",
		},
	}, nil
}

func (d *SANDriver) PullSnapshot(snapIdentifier string) (*model.VolumeSnapshotSpec, error) {
	// not used, do nothing
	return nil, &model.NotImplementError{"method PullSnapshot has not been implemented yet"}
}

func (d *SANDriver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {

	var snapName = getSnapshotName(opt.GetId())
	var volName = getVolumeName(opt.GetVolumeId())

	snapConfig := d.GetSnapshotConfig(snapName, volName)

	err := deleteSnapshot(snapConfig)

	if err != nil {
		msg := fmt.Sprintf("delete volume snapshot (%s) failed: %v", opt.GetId(), err)
		log.Error(msg)
		return err
	}
	log.Infof("volume snapshot (%s) deleted successfully", opt.GetId())
	return nil
}

func (d *SANDriver) ListPools() ([]*model.StoragePoolSpec, error) {

	var pools []*model.StoragePoolSpec

	aggregates, err := VserverGetAggregateNames() //d.sanStorageDriver.API.VserverGetAggregateNames()

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
		aggregate, _ := AggregateCommitment(aggr) //d.sanStorageDriver.API.AggregateCommitment(aggr)
		aggregateCapacity := aggregate.AggregateSize / bytesGiB
		aggregateAllocatedCapacity := aggregate.TotalAllocated / bytesGiB

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
			MultiAttach:      c.Pool[aggr].MultiAttach,
		}
		if pool.AvailabilityZone == "" {
			pool.AvailabilityZone = DefaultAZ
		}
		pools = append(pools, pool)
	}

	log.Info("list pools successfully")
	return pools, nil
}

func (d *SANDriver) InitializeSnapshotConnection(opt *pb.CreateSnapshotAttachmentOpts) (*model.ConnectionInfo, error) {

	return nil, &model.NotImplementError{S: "method InitializeSnapshotConnection has not been implemented yet."}
}

func (d *SANDriver) TerminateSnapshotConnection(opt *pb.DeleteSnapshotAttachmentOpts) error {

	return &model.NotImplementError{S: "method TerminateSnapshotConnection has not been implemented yet."}

}

func (d *SANDriver) CreateVolumeGroup(opt *pb.CreateVolumeGroupOpts) (*model.VolumeGroupSpec, error) {
	return nil, &model.NotImplementError{"method CreateVolumeGroup has not been implemented yet"}
}

func (d *SANDriver) UpdateVolumeGroup(opt *pb.UpdateVolumeGroupOpts) (*model.VolumeGroupSpec, error) {
	return nil, &model.NotImplementError{"method UpdateVolumeGroup has not been implemented yet"}
}

func (d *SANDriver) DeleteVolumeGroup(opt *pb.DeleteVolumeGroupOpts) error {
	return &model.NotImplementError{"method DeleteVolumeGroup has not been implemented yet"}
}
