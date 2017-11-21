// Copyright (c) 2017 OpenSDS Authors.
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
	"strconv"
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
	UnitGi          = 1024 * 1024 * 1024
)

type AuthOptions struct {
	Username  string `yaml:"userName,omitempty"`
	Password  string `yaml:"password,omitempty"`
	Endpoints string `yaml:"endpoints,omitempty"`
}

type DoradoConfig struct {
	AuthOptions `yaml:"authOptions"`
	Pool        map[string]PoolProperties `yaml:"pool,flow"`
}

type Driver struct {
	conf   *DoradoConfig
	client *DoradoClient
}

func (d *Driver) sector2Gb(sec string) int64 {
	capa, err := strconv.ParseInt(sec, 10, 64)
	if err != nil {
		log.Error("Convert capacity from string to number failed, error:", err)
		return 0
	}
	return capa * 512 / UnitGi
}

func (d *Driver) gb2Sector(gb int64) int64 {
	return gb * UnitGi / 512
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
	client, err := NewClient(conf.Username, conf.Password, dp)
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
	//Convert the storage unit Giga to sector

	lun, err := d.client.CreateVolume(opt.GetName(), d.gb2Sector(opt.GetSize()), opt.GetDescription())
	if err != nil {
		log.Error("Create Volume Failed:", err)
		return nil, err
	}
	log.Infof("Create volume %s (%s) success.", opt.GetName(), lun.Id)
	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: lun.Id,
		},
		Name:             lun.Name,
		Size:             d.sector2Gb(lun.Capacity),
		Description:      lun.Description,
		AvailabilityZone: "dorado",
	}, nil
}

func (d *Driver) PullVolume(volID string) (*model.VolumeSpec, error) {
	lun, err := d.client.GetVolume(volID)
	if err != nil {
		return nil, err
	}

	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: lun.Id,
		},
		Name:             lun.Name,
		Size:             d.sector2Gb(lun.Capacity),
		Description:      lun.Description,
		AvailabilityZone: "dorado",
	}, nil
}

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	err := d.client.DeleteVolume(opt.Id)
	if err != nil {
		log.Errorf("Delete volume failed, volume id =%s , Error:%s", opt.GetId())
	}
	log.Info("Remove volume success, volume id =", opt.GetId())
	return nil
}

func (d *Driver) InitializeConnection(opt *pb.CreateAttachmentOpts) (*model.ConnectionInfo, error) {
	return &model.ConnectionInfo{}, nil
}

func (d *Driver) TerminateConnection(opt *pb.DeleteAttachmentOpts) error { return nil }

func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	snap, err := d.client.CreateSnapshot(opt.GetVolumeId(), opt.GetName(), opt.GetDescription())
	if err != nil {
		return nil, err
	}
	return &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: snap.Id,
		},
		Name:        snap.Name,
		Description: snap.Description,
		VolumeId:    snap.ParentId,
		Size:        0,
	}, nil
}

func (d *Driver) PullSnapshot(id string) (*model.VolumeSnapshotSpec, error) {
	snap, err := d.client.GetSnapshot(id)
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
	err := d.client.DeleteSnapshot(opt.GetId())
	if err != nil {
		log.Errorf("Delete volume snapshot failed, volume id =%s , Error:%s", opt.GetId())
	}
	log.Info("Remove volume snapshot success, volume id =", opt.GetId())
	return nil
}

func (d *Driver) buildPoolParam(proper PoolProperties) map[string]interface{} {
	param := make(map[string]interface{})
	param["diskType"] = proper.DiskType
	param["iops"] = proper.IOPS
	param["bandwidth"] = proper.BandWidth
	return param
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
		param := d.buildPoolParam(c.Pool[p.Name])
		pol := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: p.Id,
			},
			Name:             p.Name,
			TotalCapacity:    d.sector2Gb(p.UserTotalCapacity),
			FreeCapacity:     d.sector2Gb(p.UserFreeCapacity),
			Parameters:       param,
			AvailabilityZone: c.Pool[p.Name].AZ,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = defaultAZ
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

