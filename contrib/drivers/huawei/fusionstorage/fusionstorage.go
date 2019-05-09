// Copyright (c) 2019 Huawei Technologies Co., Ltd. All Rights Reserved.
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
	"os"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	. "github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils/config"
	uuid "github.com/satori/go.uuid"
)

type Driver struct {
	cli  *FsCli
	conf *Config
}

type AuthOptions struct {
	Username        string   `yaml:"username"`
	Password        string   `yaml:"password"`
	Url             string   `yaml:"url"`
	FmIp            string   `yaml:"fmIp,omitempty"`
	FsaIp           []string `yaml:"fsaIp,flow"`
	PwdEncrypter    string   `yaml:"PwdEncrypter,omitempty"`
	EnableEncrypted bool     `yaml:"EnableEncrypted,omitempty"`
}

type Config struct {
	AuthOptions `yaml:"authOptions"`
	Pool        map[string]PoolProperties `yaml:"pool,flow"`
}

func (d *Driver) Setup() error {
	conf := &Config{}

	d.conf = conf

	path := config.CONF.OsdsDock.Backends.HuaweiFusionStorage.ConfigPath
	if path == "" {
		path = DefaultConfPath
	}

	Parse(conf, path)

	client, err := newRestCommon(conf)
	if err != nil {
		log.Errorf("Get new client failed, %v", err)
		return err
	}

	err = client.StartServer()
	if err != nil {
		log.Errorf("Get new client failed, %v", err)
		return err
	}

	d.cli = client

	log.Info("Get new client success")
	return nil
}

func (d *Driver) Unset() error {
	return nil
}

func EncodeName(id string) string {
	return NamePrefix + "-" + id
}

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*VolumeSpec, error) {
	name := EncodeName(opt.GetId())
	err := d.cli.createVolume(name, opt.GetPoolName(), opt.GetSize()<<UnitGiShiftBit)
	if err != nil {
		log.Errorf("Create volume %s (%s) failed: %s", opt.GetName(), opt.GetId(), err)
		return nil, err
	}
	log.Infof("Create volume %s (%s) success.", opt.GetName(), opt.GetId())
	return &VolumeSpec{
		BaseModel: &BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.Size,
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
		PoolId:           opt.GetPoolId(),
		Metadata: map[string]string{
			LunId: name,
		},
	}, nil
}

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	name := EncodeName(opt.GetId())
	err := d.cli.deleteVolume(name)
	if err != nil {
		log.Errorf("Delete volume (%s) failed: %v", opt.GetId(), err)
		return err
	}
	log.Infof("Delete volume (%s) success.", opt.GetId())
	return nil
}

func (d *Driver) ListPools() ([]*StoragePoolSpec, error) {
	var pols []*StoragePoolSpec
	pools, err := d.cli.queryPoolInfo()
	if err != nil {
		log.Errorf("List pools failed: %v", err)
		return nil, err
	}

	c := d.conf
	for _, p := range pools.Pools {
		poolId := strconv.Itoa(p.PoolId)
		if _, ok := c.Pool[poolId]; !ok {
			continue
		}
		host, _ := os.Hostname()
		name := fmt.Sprintf("%s:%s:%s", host, d.conf.Url, poolId)
		pol := &StoragePoolSpec{
			BaseModel: &BaseModel{
				// Make sure uuid is unique
				Id: uuid.NewV5(uuid.NamespaceOID, name).String(),
			},
			Name:             "fusionstorage_pool" + poolId,
			TotalCapacity:    p.TotalCapacity >> UnitGiShiftBit,
			FreeCapacity:     (p.TotalCapacity - p.UsedCapacity) >> UnitGiShiftBit,
			StorageType:      c.Pool[poolId].StorageType,
			Extras:           c.Pool[poolId].Extras,
			AvailabilityZone: c.Pool[poolId].AvailabilityZone,
			MultiAttach:      c.Pool[poolId].MultiAttach,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = DefaultAZ
		}
		pols = append(pols, pol)
	}

	log.Info("List pools success")
	return pols, nil
}

func (d *Driver) InitializeConnection(opt *pb.CreateVolumeAttachmentOpts) (*ConnectionInfo, error) {
	lunId := opt.GetMetadata()[LunId]
	if lunId == "" {
		msg := "Lun id is empty"
		log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	hostInfo := opt.GetHostInfo()

	initiator := hostInfo.GetInitiator()
	hostName := hostInfo.GetHost()

	if initiator == "" || hostName == "" {
		msg := "Host name or initiator is empty."
		log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Create port if not exist.
	err := d.cli.queryPortInfo(initiator)
	if err != nil {
		if strings.Contains(err.Error(), InitiatorNotExistErrorCode) {
			err := d.cli.createPort(initiator)
			if err != nil {
				log.Errorf("Create port failed: %v", err)
				return nil, err
			}
		} else {
			log.Errorf("Query port info failed: %v", err)
			return nil, err
		}
	}

	// Create host if not exist.
	isFind, err := d.cli.queryHostInfo(hostName)
	if err != nil {
		log.Errorf("Query host info failed: %v", err)
		return nil, err
	}

	if !isFind {
		err = d.cli.createHost(hostInfo)
		if err != nil {
			log.Errorf("Create host failed: %v", err)
			return nil, err
		}
	}

	// Add port to host if port not add to the host
	hostPortMap, err := d.cli.queryHostByPort(initiator)
	if err != nil {
		log.Errorf("Query host by port failed: %v", err)
		return nil, err
	}

	h, ok := hostPortMap.PortHostMap[initiator]
	if ok && h[0] != hostName {
		msg := fmt.Sprintf("Initiator is already added to another host, host name =%s", h[0])
		log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	if !ok {
		err = d.cli.addPortToHost(hostName, initiator)
		if err != nil {
			log.Errorf("Add port to host failed: %v", err)
			return nil, err
		}
	}

	// Map volume to host
	err = d.cli.addLunsToHost(hostName, lunId)
	if err != nil {
		log.Errorf("Add luns to host failed: %v", err)
		return nil, err
	}

	// Get target lun id
	hostLunList, err := d.cli.queryHostLunInfo(hostName)
	if err != nil {
		log.Errorf("Query host lun info failed: %v", err)
		return nil, err
	}

	var targetLunId int
	for _, v := range hostLunList.LunList {
		if v.Name == lunId {
			targetLunId = v.Id
		}
	}

	// Get target iscsi portal info
	targetPortalInfo, err := d.cli.queryIscsiPortal(initiator)
	if err != nil {
		log.Errorf("Query iscsi portal failed: %v", err)
		return nil, err
	}

	var targetIQN []string
	var targetPortal []string

	for _, v := range targetPortalInfo {
		iscsiTarget := strings.Split(v, ",")
		targetIQN = append(targetIQN, iscsiTarget[1])
		targetPortal = append(targetPortal, iscsiTarget[0])
	}

	connInfo := &ConnectionInfo{
		DriverVolumeType: ISCSIProtocol,
		ConnectionData: map[string]interface{}{
			"target_discovered": true,
			"volume_id":         opt.GetVolumeId(),
			"description":       "huawei",
			"host_name":         hostName,
			"targetLun":         targetLunId,
			"connect_type":      FusionstorageIscsi,
			"host":              hostName,
			"initiator":         initiator,
			"targetIQN":         targetIQN,
			"targetPortal":      targetPortal,
		},
	}

	log.Infof("Initialize connection success: %v", connInfo)
	return connInfo, nil
}

func (d *Driver) TerminateConnection(opt *pb.DeleteVolumeAttachmentOpts) error {
	lunId := opt.GetMetadata()[LunId]
	if lunId == "" {
		msg := "Lun id is empty."
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	hostInfo := opt.GetHostInfo()

	initiator := hostInfo.GetInitiator()
	hostName := hostInfo.GetHost()

	if initiator == "" || hostName == "" {
		msg := "Host name or initiator is empty."
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	// Make sure that host is exist.
	hostIsFind, err := d.cli.queryHostInfo(hostName)
	if err != nil {
		msg := fmt.Sprintf("Query host failed, host name =%s, error: %v", hostName, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	if !hostIsFind {
		msg := fmt.Sprintf("Host can not be found, host name =%s", hostName)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	// Check whether the volume attach to the host
	hostLunList, err := d.cli.queryHostLunInfo(hostName)
	if err != nil {
		log.Errorf("Query host lun info failed, %v", err)
		return err
	}

	var lunIsFind = false
	for _, v := range hostLunList.LunList {
		if v.Name == lunId {
			lunIsFind = true
			break
		}
	}

	if !lunIsFind {
		msg := fmt.Sprintf("The lun %s is not attach to the host %s", lunId, hostName)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	// Remove lun from host
	err = d.cli.deleteLunFromHost(hostName, lunId)
	if err != nil {
		log.Errorf("Delete lun from host failed, %v", err)
		return err
	}

	// Remove initiator and host if there is no lun belong to the host
	hostLunList, err = d.cli.queryHostLunInfo(hostName)
	if err != nil {
		log.Errorf("Query host lun info failed, %v", err)
		return err
	}

	if len(hostLunList.LunList) == 0 {
		d.cli.deletePortFromHost(hostName, initiator)
		d.cli.deleteHost(hostName)
		d.cli.deletePort(initiator)
	}

	log.Info("Terminate Connection success.")
	return nil
}

func (d *Driver) PullVolume(volIdentifier string) (*VolumeSpec, error) {
	// Not used , do nothing
	return nil, nil
}

func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*VolumeSpec, error) {
	err := d.cli.extendVolume(EncodeName(opt.GetId()), opt.GetSize()<<UnitGiShiftBit)
	if err != nil {
		log.Errorf("Extend volume %s (%s) failed: %v", opt.GetName(), opt.GetId(), err)
		return nil, err
	}
	log.Infof("Extend volume %s (%s) success.", opt.GetName(), opt.GetId())
	return &VolumeSpec{
		BaseModel: &BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.GetSize(),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
	}, nil
}

func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*VolumeSnapshotSpec, error) {
	snapName := EncodeName(opt.GetId())
	volName := EncodeName(opt.GetVolumeId())

	if err := d.cli.createSnapshot(snapName, volName); err != nil {
		log.Errorf("Create snapshot %s (%s) failed: %s", opt.GetName(), opt.GetId(), err)
		return nil, err
	}

	log.Errorf("Create snapshot %s (%s) success.", opt.GetName(), opt.GetId())
	return &VolumeSnapshotSpec{
		BaseModel: &BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Description: opt.GetDescription(),
		VolumeId:    opt.GetVolumeId(),
		Size:        opt.GetSize(),
	}, nil
}

func (d *Driver) PullSnapshot(snapIdentifier string) (*VolumeSnapshotSpec, error) {
	return nil, nil
}

func (d *Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	err := d.cli.deleteSnapshot(EncodeName(opt.GetId()))
	if err != nil {
		log.Errorf("Delete volume snapshot (%s) failed: %v", opt.GetId(), err)
		return err
	}
	log.Infof("Remove volume snapshot (%s) success", opt.GetId())
	return nil
}

func (d *Driver) InitializeSnapshotConnection(opt *pb.CreateSnapshotAttachmentOpts) (*ConnectionInfo, error) {
	return nil, &NotImplementError{S: "method InitializeSnapshotConnection has not been implemented yet."}
}

func (d *Driver) TerminateSnapshotConnection(opt *pb.DeleteSnapshotAttachmentOpts) error {
	return &NotImplementError{S: "method TerminateSnapshotConnection has not been implemented yet."}
}

func (d *Driver) CreateVolumeGroup(opt *pb.CreateVolumeGroupOpts) (*VolumeGroupSpec, error) {
	return nil, &NotImplementError{"method CreateVolumeGroup has not been implemented yet"}
}

func (d *Driver) UpdateVolumeGroup(opt *pb.UpdateVolumeGroupOpts) (*VolumeGroupSpec, error) {
	return nil, &NotImplementError{"method UpdateVolumeGroup has not been implemented yet"}
}

func (d *Driver) DeleteVolumeGroup(opt *pb.DeleteVolumeGroupOpts) error {
	return &NotImplementError{"method DeleteVolumeGroup has not been implemented yet"}
}
