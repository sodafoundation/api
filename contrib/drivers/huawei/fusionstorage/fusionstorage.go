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

package fusionstorage

import (
	"errors"
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

func (d *Driver) Setup() error {
	conf := &Config{}

	d.Conf = conf

	path := config.CONF.OsdsDock.Backends.HuaweiFusionStorage.ConfigPath
	if path == "" {
		path = DefaultConfPath
	}

	Parse(conf, path)

	client, err := newRestCommon(conf)
	if err != nil {
		msg := fmt.Sprintf("get new client failed, %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	d.Client = client

	log.Info("get new client success")
	return nil
}

func (d *Driver) Unset() error {
	return d.Client.logout()
}

func EncodeName(id string) string {
	return NamePrefix + "-" + id
}

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*VolumeSpec, error) {
	name := EncodeName(opt.GetId())
	err := d.Client.createVolume(name, opt.GetPoolName(), opt.GetSize()<<UnitGiShiftBit)
	if err != nil {
		msg := fmt.Sprintf("create volume %s (%s) failed: %s", opt.GetName(), opt.GetId(), err)
		log.Error(msg)
		return nil, errors.New(msg)
	}
	log.V(8).Infof("create volume %s (%s) success.", opt.GetName(), opt.GetId())
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
	err := d.Client.deleteVolume(name)
	if err != nil {
		msg := fmt.Sprintf("delete volume (%s) failed: %v", opt.GetId(), err)
		log.Error(msg)
		return errors.New(msg)
	}
	log.Infof("delete volume (%s) success.", opt.GetId())
	return nil
}

func (d *Driver) ListPools() ([]*StoragePoolSpec, error) {
	var pols []*StoragePoolSpec
	pools, err := d.Client.queryPoolInfo()
	if err != nil {
		msg := fmt.Sprintf("list pools failed: %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	c := d.Conf
	for _, p := range pools.Pools {
		poolId := strconv.Itoa(p.PoolId)
		if _, ok := c.Pool[poolId]; !ok {
			continue
		}
		host, _ := os.Hostname()
		name := fmt.Sprintf("%s:%s:%s", host, d.Conf.Url, poolId)
		pol := &StoragePoolSpec{
			BaseModel: &BaseModel{
				// Make sure uuid is unique
				Id: uuid.NewV5(uuid.NamespaceOID, name).String(),
			},
			Name:             poolId,
			TotalCapacity:    p.TotalCapacity >> UnitGiShiftBit,
			FreeCapacity:     (p.TotalCapacity - p.UsedCapacity) >> UnitGiShiftBit,
			StorageType:      c.Pool[poolId].StorageType,
			Extras:           c.Pool[poolId].Extras,
			AvailabilityZone: c.Pool[poolId].AvailabilityZone,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = DefaultAZ
		}
		pols = append(pols, pol)
	}

	log.Info("list pools success")
	return pols, nil
}

func (d *Driver) InitializeConnection(opt *pb.CreateVolumeAttachmentOpts) (*ConnectionInfo, error) {
	lunId := opt.GetMetadata()[LunId]
	if lunId == "" {
		msg := "lun id is empty"
		log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	hostInfo := opt.GetHostInfo()

	initiator := hostInfo.GetInitiator()
	hostName := hostInfo.GetHost()

	if initiator == "" || hostName == "" {
		msg := "host name or initiator cannot be empty"
		log.Error(msg)
		return nil, errors.New(msg)
	}

	// Create port if not exist.
	if err := d.CreatePortIfNotExist(initiator); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Create host if not exist.
	if err := d.CreateHostIfNotExist(hostInfo); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Add port to host if port not add to the host
	if err := d.AddPortToHost(initiator, hostName); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Map volume to host
	err := d.Client.addLunsToHost(hostName, lunId)
	if err != nil && !strings.Contains(err.Error(), VolumeAlreadyInHostErrorCode) {
		msg := fmt.Sprintf("add luns to host failed: %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	// Get target lun id
	targetLunId, err := d.GetTgtLunID(hostName, lunId)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	targetIQN, targetPortal, err := d.GetTargetPortal(initiator)
	if err != nil {
		msg := fmt.Sprintf("get target portals and iqns failed: %v", err)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}

	connInfo := &ConnectionInfo{
		DriverVolumeType: opt.GetAccessProtocol(),
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

	log.Infof("initialize connection success: %v", connInfo)
	return connInfo, nil
}

func (d *Driver) GetTgtLunID(hostName, sourceLunID string) (int, error) {
	hostLunList, err := d.Client.queryHostLunInfo(hostName)
	if err != nil {
		return -1, fmt.Errorf("query host lun info failed: %v", err)
	}

	var targetLunId int
	for _, v := range hostLunList.LunList {
		if v.Name == sourceLunID {
			targetLunId = v.Id
		}
	}

	return targetLunId, nil
}

func (d *Driver) AddPortToHost(initiator, hostName string) error {
	hostPortMap, err := d.Client.queryHostByPort(initiator)
	if err != nil {
		return fmt.Errorf("query host by port failed: %v", err)
	}

	h, ok := hostPortMap.PortHostMap[initiator]
	if ok && h[0] != hostName {
		return fmt.Errorf("initiator is already added to another host, host name =%s", h[0])
	}

	if !ok {
		err = d.Client.addPortToHost(hostName, initiator)
		if err != nil {
			return fmt.Errorf("add port to host failed: %v", err)
		}
	}

	return nil
}

func (d *Driver) CreateHostIfNotExist(hostInfo *pb.HostInfo) error {
	isFind, err := d.Client.queryHostInfo(hostInfo.GetHost())
	if err != nil {
		return fmt.Errorf("query host info failed: %v", err)
	}

	if !isFind {
		err = d.Client.createHost(hostInfo)
		if err != nil {
			return fmt.Errorf("create host failed: %v", err)
		}
	}

	return nil
}

func (d *Driver) CreatePortIfNotExist(initiator string) error {
	err := d.Client.queryPortInfo(initiator)
	if err != nil {
		if strings.Contains(err.Error(), InitiatorNotExistErrorCodeVersion6) ||
			strings.Contains(err.Error(), InitiatorNotExistErrorCodeVersion8) {
			err := d.Client.createPort(initiator)
			if err != nil {
				return fmt.Errorf("create port failed: %v", err)
			}
		} else {
			return fmt.Errorf("query port info failed: %v", err)
		}
	}
	return nil
}

func (d *Driver) GetTargetPortal(initiator string) ([]string, []string, error) {
	if d.Conf.Version == ClientVersion6_3 {
		return d.GeTgtPortalAndIQNVersion6_3(initiator)
	}
	if d.Conf.Version == ClientVersion8_0 {
		return d.GeTgtPortalAndIQNVersion8_0()
	}

	return nil, nil, errors.New("cannot find any target portal and iqn")
}

func (d *Driver) GeTgtPortalAndIQNVersion6_3(initiator string) ([]string, []string, error) {
	targetPortalInfo, err := d.Client.queryIscsiPortalVersion6(initiator)
	if err != nil {
		msg := fmt.Sprintf("query iscsi portal failed: %v", err)
		log.Error(msg)
		return nil, nil, errors.New(msg)
	}

	var targetIQN []string
	var targetPortal []string

	for _, v := range targetPortalInfo {
		iscsiTarget := strings.Split(v, ",")
		targetIQN = append(targetIQN, iscsiTarget[1])
		targetPortal = append(targetPortal, iscsiTarget[0])
	}

	return targetIQN, targetPortal, nil
}

func (d *Driver) GeTgtPortalAndIQNVersion8_0() ([]string, []string, error) {
	targetPortalInfo, err := d.Client.queryIscsiPortalVersion8()
	if err != nil {
		msg := fmt.Sprintf("query iscsi portal failed: %v", err)
		log.Error(msg)
		return nil, nil, errors.New(msg)
	}

	var targetPortal []string

	for _, v := range targetPortalInfo.NodeResultList {
		for _, p := range v.PortalList {
			if p.Status == "active" {
				targetPortal = append(targetPortal, p.IscsiPortal)
			}
		}
	}

	if len(targetPortal) == 0 {
		return nil, nil, errors.New("the iscsi target portal is empty")
	}

	return nil, targetPortal, nil
}

func (d *Driver) TerminateConnection(opt *pb.DeleteVolumeAttachmentOpts) error {
	lunId := opt.GetMetadata()[LunId]
	if lunId == "" {
		msg := "lun id is empty."
		log.Error(msg)
		return errors.New(msg)
	}

	hostInfo := opt.GetHostInfo()

	initiator := hostInfo.GetInitiator()
	hostName := hostInfo.GetHost()

	if initiator == "" || hostName == "" {
		msg := "host name or initiator is empty."
		log.Error(msg)
		return errors.New(msg)
	}

	// Make sure that host is exist.
	if err := d.CheckHostIsExist(hostName); err != nil {
		return err
	}

	// Check whether the volume attach to the host
	if err := d.CheckVolAttachToHost(hostName, lunId); err != nil {
		return err
	}

	// Remove lun from host
	err := d.Client.deleteLunFromHost(hostName, lunId)
	if err != nil {
		msg := fmt.Sprintf("delete lun from host failed, %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	// Remove initiator and host if there is no lun belong to the host
	hostLunList, err := d.Client.queryHostLunInfo(hostName)
	if err != nil {
		msg := fmt.Sprintf("query host lun info failed, %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	if len(hostLunList.LunList) == 0 {
		d.Client.deletePortFromHost(hostName, initiator)
		d.Client.deleteHost(hostName)
		d.Client.deletePort(initiator)
	}

	log.Info("terminate Connection success.")
	return nil
}

func (d *Driver) CheckVolAttachToHost(hostName, lunId string) error {
	hostLunList, err := d.Client.queryHostLunInfo(hostName)
	if err != nil {
		msg := fmt.Sprintf("query host lun info failed, %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	var lunIsFind = false
	for _, v := range hostLunList.LunList {
		if v.Name == lunId {
			lunIsFind = true
			break
		}
	}

	if !lunIsFind {
		msg := fmt.Sprintf("the lun %s is not attach to the host %s", lunId, hostName)
		log.Error(msg)
		return errors.New(msg)
	}

	return nil
}

func (d *Driver) CheckHostIsExist(hostName string) error {
	hostIsFind, err := d.Client.queryHostInfo(hostName)
	if err != nil {
		msg := fmt.Sprintf("query host failed, host name %s, error: %v", hostName, err)
		log.Error(msg)
		return errors.New(msg)
	}

	if !hostIsFind {
		msg := fmt.Sprintf("host can not be found, host name =%s", hostName)
		log.Error(msg)
		return errors.New(msg)
	}
	return nil
}

func (d *Driver) PullVolume(volIdentifier string) (*VolumeSpec, error) {
	// Not used , do nothing
	return nil, nil
}

func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*VolumeSpec, error) {
	err := d.Client.extendVolume(EncodeName(opt.GetId()), opt.GetSize()<<UnitGiShiftBit)
	if err != nil {
		msg := fmt.Sprintf("extend volume %s (%s) failed: %v", opt.GetName(), opt.GetId(), err)
		log.Error(msg)
		return nil, errors.New(msg)
	}
	log.Infof("extend volume %s (%s) success.", opt.GetName(), opt.GetId())
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

	if err := d.Client.createSnapshot(snapName, volName); err != nil {
		msg := fmt.Sprintf("create snapshot %s (%s) failed: %s", opt.GetName(), opt.GetId(), err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	log.Infof("create snapshot %s (%s) success.", opt.GetName(), opt.GetId())
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
	err := d.Client.deleteSnapshot(EncodeName(opt.GetId()))
	if err != nil {
		msg := fmt.Sprintf("delete volume snapshot (%s) failed: %v", opt.GetId(), err)
		log.Error(msg)
		return errors.New(msg)
	}
	log.Infof("remove volume snapshot (%s) success", opt.GetId())
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
