// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package dorado

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/config"
)

const (
	defaultConfPath = "/etc/opensds/driver/huawei_dorado.yaml"
	defaultAZ       = "default"
)

const (
	KLunId  = "huaweiLunId"
	KSnapId = "huaweiSnapId"
)

type AuthOptions struct {
	Username  string `yaml:"username,omitempty"`
	Password  string `yaml:"password,omitempty"`
	Endpoints string `yaml:"endpoints,omitempty"`
	Insecure  bool   `yaml:"insecure,omitempty"`
}

type DoradoConfig struct {
	AuthOptions `yaml:"authOptions"`
	Pool        map[string]PoolProperties `yaml:"pool,flow"`
	TargetIp    string                    `yaml:"targetIp,omitempty"`
}

type Driver struct {
	conf   *DoradoConfig
	client *DoradoClient
}

func (d *Driver) Setup() error {
	// Read huawei dorado config file
	conf := &DoradoConfig{}
	d.conf = conf
	path := config.CONF.OsdsDock.Backends.HuaweiDorado.ConfigPath

	if "" == path {
		path = defaultConfPath
	}
	Parse(conf, path)
	dp := strings.Split(conf.Endpoints, ",")
	client, err := NewClient(conf.Username, conf.Password, dp, conf.Insecure)
	d.client = client
	if err != nil {
		log.Errorf("Get new client failed, %v", err)
		return err
	}
	return nil
}

func (d *Driver) Unset() error {
	d.client.logout()
	return nil
}

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	name := EncodeName(opt.GetId())
	desc := TruncateDescription(opt.GetDescription())
	lun, err := d.client.CreateVolume(name, opt.GetSize(), desc)
	if err != nil {
		log.Error("Create Volume Failed:", err)
		return nil, err
	}
	log.Infof("Create volume %s (%s) success.", opt.GetName(), lun.Id)
	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             Sector2Gb(lun.Capacity),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
		Metadata: map[string]string{
			KLunId: lun.Id,
		},
	}, nil
}

func (d *Driver) PullVolume(volID string) (*model.VolumeSpec, error) {
	name := EncodeName(volID)
	lun, err := d.client.GetVolumeByName(name)
	if err != nil {
		return nil, err
	}

	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: lun.Id,
		},
		Size:             Sector2Gb(lun.Capacity),
		Description:      lun.Description,
		AvailabilityZone: lun.ParentName,
	}, nil
}

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	lunId := opt.GetMetadata()[KLunId]
	err := d.client.DeleteVolume(lunId)
	if err != nil {
		log.Errorf("Delete volume failed, volume id =%s , Error:%s", opt.GetId(), err)
		return err
	}
	log.Info("Remove volume success, volume id =", opt.GetId())
	return nil
}

// ExtendVolume ...
func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	lunId := opt.GetMetadata()[KLunId]
	err := d.client.ExtendVolume(opt.GetSize(), lunId)
	if err != nil {
		log.Error("Extend Volume Failed:", err)
		return nil, err
	}

	log.Infof("Extend volume %s (%s) success.", opt.GetName(), opt.GetId())
	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.GetSize(),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
	}, nil
}

func (d *Driver) getTargetInfo() (string, string, error) {
	tgtIp := d.conf.TargetIp
	resp, err := d.client.ListTgtPort()
	if err != nil {
		return "", "", err
	}
	for _, itp := range resp.Data {
		items := strings.Split(itp.Id, ",")
		iqn := strings.Split(items[0], "+")[1]
		items = strings.Split(iqn, ":")
		ip := items[len(items)-1]
		if tgtIp == ip {
			return iqn, ip, nil
		}
	}
	msg := fmt.Sprintf("Not find configuration targetIp: %v in device", tgtIp)
	return "", "", errors.New(msg)
}

func (d *Driver) InitializeConnection(opt *pb.CreateAttachmentOpts) (*model.ConnectionInfo, error) {

	lunId := opt.GetMetadata()[KLunId]
	hostInfo := opt.GetHostInfo()
	// Create host if not exist.
	hostId, err := d.client.AddHostWithCheck(hostInfo)
	if err != nil {
		log.Errorf("Add host failed, host name =%s, error: %v", hostInfo.Host, err)
		return nil, err
	}

	// Add initiator to the host.
	if err = d.client.AddInitiatorToHostWithCheck(hostId, hostInfo.Initiator); err != nil {
		log.Errorf("Add initiator to host failed, host id=%s, initiator=%s, error: %v", hostId, hostInfo.Initiator, err)
		return nil, err
	}

	// Add host to hostgroup.
	hostGrpId, err := d.client.AddHostToHostGroup(hostId)
	if err != nil {
		log.Errorf("Add host to group failed, host id=%s, error: %v", hostId, err)
		return nil, err
	}

	// Mapping lungroup and hostgroup to view.
	if err = d.client.DoMapping(lunId, hostGrpId, hostId); err != nil {
		log.Errorf("Do mapping failed, lun id=%s, hostGrpId=%s, hostId=%s, error: %v",
			lunId, hostGrpId, hostId, err)
		return nil, err
	}

	tgtIqn, tgtIp, err := d.getTargetInfo()
	if err != nil {
		log.Error("Get the target info failed,", err)
		return nil, err
	}
	tgtLun, err := d.client.GetHostLunId(hostId, lunId)
	if err != nil {
		log.Error("Get the get host lun id failed,", err)
		return nil, err
	}
	connInfo := &model.ConnectionInfo{
		DriverVolumeType: "iscsi",
		ConnectionData: map[string]interface{}{
			"targetDiscovered": true,
			"targetIQN":        tgtIqn,
			"targetPortal":     tgtIp + ":3260",
			"discard":          false,
			"targetLun":        tgtLun,
		},
	}
	return connInfo, nil
}

func (d *Driver) TerminateConnection(opt *pb.DeleteAttachmentOpts) error {
	lunId := opt.GetMetadata()[KLunId]
	hostId, err := d.client.GetHostIdByName(opt.GetHostInfo().GetHost())
	if err != nil {
		return err
	}
	lunGrpId, _ := d.client.FindLunGroup(LunGroupPrefix + hostId)
	hostGrpId, _ := d.client.FindHostGroup(HostGroupPrefix + hostId)
	viewId, _ := d.client.FindMappingView(MappingViewPrefix + hostId)
	if viewId != "" {
		d.client.RemoveLunGroupFromMappingView(viewId, lunGrpId)
		d.client.RemoveHostGroupFromMappingView(viewId, hostGrpId)
		d.client.DeleteMappingView(viewId)
	}
	if hostGrpId != "" {
		d.client.RemoveHostFromHostGroup(hostGrpId, hostId)
		d.client.DeleteHostGroup(hostGrpId)
	}
	if lunGrpId != "" {
		d.client.RemoveLunFromLunGroup(lunGrpId, lunId)
		d.client.DeleteLunGroup(lunGrpId)
	}
	d.client.RemoveIscsiFromHost(opt.GetHostInfo().GetInitiator())
	d.client.DeleteHost(hostId)
	return nil
}

func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	lunId := opt.GetMetadata()[KLunId]
	name := EncodeName(opt.GetId())
	desc := TruncateDescription(opt.GetDescription())
	snap, err := d.client.CreateSnapshot(lunId, name, desc)
	if err != nil {
		return nil, err
	}
	return &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Description: opt.GetDescription(),
		VolumeId:    opt.GetVolumeId(),
		Size:        0,
		Metadata: map[string]string{
			KSnapId: snap.Id,
		},
	}, nil
}

func (d *Driver) PullSnapshot(id string) (*model.VolumeSnapshotSpec, error) {
	name := EncodeName(id)
	snap, err := d.client.GetSnapshotByName(name)
	if err != nil {
		return nil, err
	}
	return &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: snap.Id,
		},
		Name:        snap.Name,
		Description: snap.Description,
		Size:        0,
		VolumeId:    snap.ParentId,
	}, nil
}

func (d *Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	id := opt.GetMetadata()[KSnapId]
	err := d.client.DeleteSnapshot(id)
	if err != nil {
		log.Errorf("Delete volume snapshot failed, volume snapshot id = %s , error: %v", opt.GetId(), err)
		return err
	}
	log.Info("Remove volume snapshot success, volume snapshot id =", opt.GetId())
	return nil
}

func (d *Driver) ListPools() ([]*model.StoragePoolSpec, error) {
	var pols []*model.StoragePoolSpec
	sp, err := d.client.ListStoragePools()
	if err != nil {
		return nil, err
	}
	for _, p := range sp {
		c := d.conf
		if _, ok := c.Pool[p.Name]; !ok {
			continue
		}

		pol := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: p.Id,
			},
			Name:             p.Name,
			TotalCapacity:    Sector2Gb(p.UserTotalCapacity),
			FreeCapacity:     Sector2Gb(p.UserFreeCapacity),
			StorageType:      c.Pool[p.Name].StorageType,
			Extras:           c.Pool[p.Name].Extras,
			AvailabilityZone: c.Pool[p.Name].AvailabilityZone,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = defaultAZ
		}
		pols = append(pols, pol)
	}
	return pols, nil
}
