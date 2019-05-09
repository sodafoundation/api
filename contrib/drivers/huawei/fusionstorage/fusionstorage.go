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
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	. "github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/satori/go.uuid"
)

func (d *Driver) Setup() error {
	conf := &Config{}

	d.conf = conf

	path := "./testdata/fusionstorage.yaml"
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

	d.cli = client

	log.Info("get new client success")
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
	err := d.cli.deleteVolume(name)
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
	pools, err := d.cli.queryPoolInfo()
	if err != nil {
		msg := fmt.Sprintf("list pools failed: %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
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
		msg := "host name or initiator is empty."
		log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Create port if not exist.
	err := d.cli.queryPortInfo(initiator)
	if err != nil {
		if strings.Contains(err.Error(), InitiatorNotExistErrorCodeVersion6) || strings.Contains(err.Error(), InitiatorNotExistErrorCodeVersion8) {
			err := d.cli.createPort(initiator)
			if err != nil {
				msg := fmt.Sprintf("create port failed: %v", err)
				log.Error(msg)
				return nil, errors.New(msg)
			}
		} else {
			msg := fmt.Sprintf("query port info failed: %v", err)
			log.Error(msg)
			return nil, errors.New(msg)
		}
	}

	// Create host if not exist.
	isFind, err := d.cli.queryHostInfo(hostName)
	if err != nil {
		msg := fmt.Sprintf("query host info failed: %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	if !isFind {
		err = d.cli.createHost(hostInfo)
		if err != nil {
			msg := fmt.Sprintf("create host failed: %v", err)
			log.Error(msg)
			return nil, errors.New(msg)
		}
	}

	// Add port to host if port not add to the host
	hostPortMap, err := d.cli.queryHostByPort(initiator)
	if err != nil {
		msg := fmt.Sprintf("query host by port failed: %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	h, ok := hostPortMap.PortHostMap[initiator]
	if ok && h[0] != hostName {
		msg := fmt.Sprintf("initiator is already added to another host, host name =%s", h[0])
		log.Error(msg)
		return nil, errors.New(msg)
	}

	if !ok {
		err = d.cli.addPortToHost(hostName, initiator)
		if err != nil {
			msg := fmt.Sprintf("add port to host failed: %v", err)
			log.Error(msg)
			return nil, errors.New(msg)
		}
	}

	// Map volume to host
	err = d.cli.addLunsToHost(hostName, lunId)
	if err != nil && !strings.Contains(err.Error(), VolumeAlreadyInHostErrorCode) {
		msg := fmt.Sprintf("add luns to host failed: %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	// Get target lun id
	hostLunList, err := d.cli.queryHostLunInfo(hostName)
	if err != nil {
		msg := fmt.Sprintf("query host lun info failed: %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	var targetLunId int
	for _, v := range hostLunList.LunList {
		if v.Name == lunId {
			targetLunId = v.Id
		}
	}

	targetIQN, targetPortal, err := d.GetTargetPortal(initiator)
	if err != nil {
		msg := fmt.Sprintf("get target portals and iqns failed: %v", err)
		log.Errorf(msg)
		return nil, errors.New(msg)
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

	log.Infof("initialize connection success: %v", connInfo)
	return connInfo, nil
}

func (d *Driver) GetTargetPortal(initiator string) ([]string, []string, error) {
	version, err := d.cli.getDeviceVersion()
	if err != nil {
		msg := fmt.Sprintf("get device version failed %v", err)
		log.Error(msg)
		return nil, nil, errors.New(msg)
	}

	regVersion6, _ := regexp.Compile("^V100R006C")
	regVersion8, _ := regexp.Compile("^8")

	if regVersion6.MatchString(version.Version) {
		return d.GeTgtPortalAndIQNVersion6(initiator)
	}

	if regVersion8.MatchString(version.Version) {
		return d.GeTgtPortalAndIQNVersion8()
	}

	return nil, nil, errors.New("cannot find any target portal and iqn")
}

func (d *Driver) GeTgtPortalAndIQNVersion6(initiator string) ([]string, []string, error) {
	err := d.cli.StartServer()
	if err != nil {
		msg := fmt.Sprintf("get new client failed, %v", err)
		log.Errorf(msg)
		return nil, nil, errors.New(msg)
	}
	targetPortalInfo, err := d.cli.queryIscsiPortalVersion6(initiator)
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

func (d *Driver) GeTgtPortalAndIQNVersion8() ([]string, []string, error) {
	targetPortalInfo, err := d.cli.queryIscsiPortalVersion8()
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
	hostIsFind, err := d.cli.queryHostInfo(hostName)
	if err != nil {
		msg := fmt.Sprintf("query host failed, host name =%s, error: %v", hostName, err)
		log.Error(msg)
		return errors.New(msg)
	}

	if !hostIsFind {
		msg := fmt.Sprintf("host can not be found, host name =%s", hostName)
		log.Error(msg)
		return errors.New(msg)
	}

	// Check whether the volume attach to the host
	hostLunList, err := d.cli.queryHostLunInfo(hostName)
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

	// Remove lun from host
	err = d.cli.deleteLunFromHost(hostName, lunId)
	if err != nil {
		msg := fmt.Sprintf("delete lun from host failed, %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	// Remove initiator and host if there is no lun belong to the host
	hostLunList, err = d.cli.queryHostLunInfo(hostName)
	if err != nil {
		msg := fmt.Sprintf("query host lun info failed, %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	if len(hostLunList.LunList) == 0 {
		d.cli.deletePortFromHost(hostName, initiator)
		d.cli.deleteHost(hostName)
		d.cli.deletePort(initiator)
	}

	log.Info("terminate Connection success.")
	return nil
}

func (d *Driver) PullVolume(volIdentifier string) (*VolumeSpec, error) {
	// Not used , do nothing
	return nil, nil
}

func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*VolumeSpec, error) {
	err := d.cli.extendVolume(EncodeName(opt.GetId()), opt.GetSize()<<UnitGiShiftBit)
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

	if err := d.cli.createSnapshot(snapName, volName); err != nil {
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
	err := d.cli.deleteSnapshot(EncodeName(opt.GetId()))
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
