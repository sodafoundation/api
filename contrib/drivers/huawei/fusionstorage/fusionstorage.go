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
	"os"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	. "github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/config"
	"github.com/satori/go.uuid"
)

const (
	DefaultConfPath = "/etc/opensds/driver/fusionstorage.yaml"
	DefaultAZ       = "default"
	NamePrefix      = "opensds"
	UnitGiShiftBit  = 10
)

type AuthOptions struct {
	FmIp  string   `yaml:"fmIp,omitempty"`
	FsaIp []string `yaml:"fsaIp,flow"`
}

type Config struct {
	AuthOptions `yaml:"authOptions"`
	Pool        map[string]PoolProperties `yaml:"pool,flow"`
}

type Driver struct {
	cli  *Cli
	conf *Config
}

func EncodeName(id string) string {
	return NamePrefix + "-" + id
}

func (d *Driver) Setup() error {
	conf := &Config{}
	d.conf = conf

	path := config.CONF.OsdsDock.Backends.HuaweiFusionStorage.ConfigPath
	if "" == path {
		path = DefaultConfPath
	}

	Parse(conf, path)
	cli, err := NewCli(conf.FmIp, conf.FsaIp)
	if err != nil {
		log.Errorf("Get new client failed, %v", err)
		return err
	}
	d.cli = cli
	return nil
}

func (d *Driver) Unset() error {
	return nil
}

func (d *Driver) createVolumeFromSnapshot(opt *pb.CreateVolumeOpts) (*VolumeSpec, error) {
	name := EncodeName(opt.GetId())
	snapName := EncodeName(opt.GetSnapshotId())
	err := d.cli.CreateVolumeFromSnapshot(name, opt.GetSize()<<UnitGiShiftBit, snapName)
	if err != nil {
		return nil, err
	}

	return &VolumeSpec{
		BaseModel: &BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.GetSize(),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
		PoolId:           opt.GetPoolId(),
		SnapshotId:       opt.GetSnapshotId(),
		Metadata:         nil,
	}, nil
}

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*VolumeSpec, error) {

	if opt.GetSnapshotId() != "" {
		return d.createVolumeFromSnapshot(opt)
	}
	name := EncodeName(opt.GetId())
	err := d.cli.CreateVolume(name, opt.GetSize()<<UnitGiShiftBit,
		true, opt.GetPoolName(), nil)
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
		Metadata:         nil,
	}, nil
}

func (d *Driver) PullVolume(volIdentifier string) (*VolumeSpec, error) {
	// Not used , do nothing
	return nil, nil
}

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	name := EncodeName(opt.GetId())
	err := d.cli.DeleteVolume(name)
	if err != nil {
		log.Errorf("Delete volume (%s) failed: %v", opt.GetId(), err)
		return err
	}
	log.Infof("Delete volume (%s) success.", opt.GetId())
	return nil
}

func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*VolumeSpec, error) {
	err := d.cli.ExtendVolume(EncodeName(opt.GetId()), opt.GetSize()<<UnitGiShiftBit)
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

func (d *Driver) InitializeConnection(opt *pb.CreateAttachmentOpts) (*ConnectionInfo, error) {
	connInfo := &ConnectionInfo{

		DriverVolumeType: DSWARE,
		ConnectionData: map[string]interface{}{
			"volumeId":  opt.GetVolumeId(),
			"manage_ip": d.conf.FmIp,
		},
	}
	return connInfo, nil
}

func (d *Driver) TerminateConnection(opt *pb.DeleteAttachmentOpts) error {
	// do nothing
	return nil
}

func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*VolumeSnapshotSpec, error) {
	name := EncodeName(opt.GetId())
	volName := EncodeName(opt.GetVolumeId())

	if err := d.cli.CreateSnapshot(name, volName, false); err != nil {
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
	err := d.cli.DeleteSnapshot(EncodeName(opt.GetId()))
	if err != nil {
		log.Errorf("Delete volume snapshot (%s) failed: %v", opt.GetId(), err)
		return err
	}
	log.Info("Remove volume snapshot (%s) success", opt.GetId())
	return nil
}

func (d *Driver) ListPools() ([]*StoragePoolSpec, error) {
	var pols []*StoragePoolSpec
	pools, err := d.cli.QueryAllPoolInfo()
	if err != nil {
		return nil, err
	}

	c := d.conf
	for _, p := range pools {
		if _, ok := c.Pool[p.PoolId]; !ok {
			continue
		}
		host, _ := os.Hostname()
		name := fmt.Sprintf("%s:%s:%s", host, c.FmIp, p.PoolId)
		pol := &StoragePoolSpec{
			BaseModel: &BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, name).String(),
			},
			Name:             p.PoolId,
			TotalCapacity:    p.TotalCapacity >> UnitGiShiftBit,
			FreeCapacity:     (p.TotalCapacity - p.UsedCapacity) >> UnitGiShiftBit,
			StorageType:      c.Pool[p.PoolId].StorageType,
			Extras:           c.Pool[p.PoolId].Extras,
			AvailabilityZone: c.Pool[p.PoolId].AvailabilityZone,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = DefaultAZ
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

func (d *Driver) InitializeSnapshotConnection(opt *pb.CreateSnapshotAttachmentOpts) (*ConnectionInfo, error) {
	return nil, &NotImplementError{S: "Method InitializeSnapshotConnection has not been implemented yet."}
}

func (d *Driver) TerminateSnapshotConnection(opt *pb.DeleteSnapshotAttachmentOpts) error {
	return &NotImplementError{S: "Method TerminateSnapshotConnection has not been implemented yet."}
}

func (d *Driver) CreateVolumeGroup(
	opt *pb.CreateVolumeGroupOpts,
	vg *VolumeGroupSpec) (*VolumeGroupSpec, error) {
	return nil, &NotImplementError{S: "Method CreateVolumeGroup has not been implemented yet."}
}
func (d *Driver) UpdateVolumeGroup(
	opt *pb.UpdateVolumeGroupOpts,
	vg *VolumeGroupSpec,
	addVolumesRef []*VolumeSpec,
	removeVolumesRef []*VolumeSpec) (*VolumeGroupSpec, []*VolumeSpec, []*VolumeSpec, error) {
	return nil, nil, nil, &NotImplementError{"Method UpdateVolumeGroup has not been implemented yet"}
}
func (d *Driver) DeleteVolumeGroup(
	opt *pb.DeleteVolumeGroupOpts,
	vg *VolumeGroupSpec,
	volumes []*VolumeSpec) (*VolumeGroupSpec, []*VolumeSpec, error) {
	return nil, nil, &NotImplementError{S: "Method DeleteVolumeGroup has not been implemented yet."}
}
