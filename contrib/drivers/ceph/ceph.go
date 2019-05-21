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

/*
This module implements ceph driver for OpenSDS. Ceph driver will pass these
operation requests about volume to go-ceph module.
*/

package ceph

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ceph/go-ceph/rados"
	"github.com/ceph/go-ceph/rbd"
	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils/config"
	uuid "github.com/satori/go.uuid"
)

const (
	opensdsPrefix   = "opensds-"
	sizeShiftBit    = 30
	defaultConfPath = "/etc/opensds/driver/ceph.yaml"
	defaultAZ       = "default"
)

const (
	KPoolName  = "CephPoolName"
	KImageName = "CephImageName"
)

type CephConfig struct {
	ConfigFile string                    `yaml:"configFile,omitempty"`
	Pool       map[string]PoolProperties `yaml:"pool,flow"`
}

func EncodeName(id string) string {
	return opensdsPrefix + id
}

func NewSrcMgr(conf *CephConfig) *SrcMgr {
	return &SrcMgr{conf: conf}
}

type SrcMgr struct {
	conn  *rados.Conn
	ioctx *rados.IOContext
	img   *rbd.Image
	conf  *CephConfig
}

func (s *SrcMgr) GetConn() (*rados.Conn, error) {
	if s.conn != nil {
		return s.conn, nil
	}
	conn, err := rados.NewConn()
	if err != nil {
		log.Error("New connect failed:", err)
		return nil, err
	}

	if err = conn.ReadConfigFile(s.conf.ConfigFile); err != nil {
		log.Error("Read config file failed:", err)
		return nil, err
	}
	if err = conn.Connect(); err != nil {
		log.Error("Connect failed:", err)
		return nil, err
	}
	s.conn = conn
	return s.conn, nil
}

func (s *SrcMgr) GetIoctx(poolName string) (*rados.IOContext, error) {
	if s.ioctx != nil {
		return s.ioctx, nil
	}

	conn, err := s.GetConn()
	if err != nil {
		return nil, err
	}
	ioctx, err := conn.OpenIOContext(poolName)
	if err != nil {
		log.Error("Open IO context failed, poolName:", poolName, err)
		return nil, err
	}
	s.ioctx = ioctx
	return s.ioctx, err
}

func (s *SrcMgr) GetImage(poolName string, imgName string, args ...interface{}) (*rbd.Image, error) {
	if s.img != nil {
		return s.img, nil
	}
	ioctx, err := s.GetIoctx(poolName)
	if err != nil {
		return nil, err
	}
	img := rbd.GetImage(ioctx, imgName)
	if err := img.Open(args...); err != nil {
		log.Error("When open image:", err)
		return nil, err
	}
	s.img = img
	return s.img, nil
}

func (s *SrcMgr) destroy() {
	if s.img != nil {
		s.img.Close()
		s.img = nil
	}
	if s.ioctx != nil {
		s.ioctx.Destroy()
		s.ioctx = nil
	}
	if s.conn != nil {
		s.conn.Shutdown()
		s.conn = nil
	}
}

type Driver struct {
	conf *CephConfig
}

func (d *Driver) Setup() error {
	d.conf = &CephConfig{ConfigFile: "/etc/ceph/ceph.conf"}
	p := config.CONF.OsdsDock.Backends.Ceph.ConfigPath
	if "" == p {
		p = defaultConfPath
	}
	_, err := Parse(d.conf, p)
	return err
}

func (d *Driver) Unset() error { return nil }

func (d *Driver) createVolumeFromSnapshot(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	poolName := opt.GetPoolName()
	srcSnapName := EncodeName(opt.GetSnapshotId())
	srcImgName := opt.GetMetadata()[KImageName]
	destImgName := EncodeName(opt.GetId())

	mgr := NewSrcMgr(d.conf)
	defer mgr.destroy()

	ioctx, err := mgr.GetIoctx(poolName)
	if err != nil {
		return nil, err
	}

	img, err := mgr.GetImage(poolName, srcImgName, srcSnapName)
	if err != nil {
		return nil, err
	}
	snap := img.GetSnapshot(srcSnapName)
	if ok, _ := snap.IsProtected(); !ok {
		if err := snap.Protect(); err != nil {
			log.Errorf("protect failed, %v", err)
			return nil, err
		}
		defer snap.Unprotect()
	}

	_, err = img.Clone(srcSnapName, ioctx, destImgName, rbd.RbdFeatureLayering, 20)
	if err != nil {
		log.Errorf("create volume (%s) from snapshot (%s) failed, %v",
			opt.GetId(), opt.GetSnapshotId(), err)
		return nil, err
	}
	log.Infof("create volume (%s) from snapshot (%s) success",
		opt.GetId(), opt.GetSnapshotId())
	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.GetSize(),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
		Metadata: map[string]string{
			KPoolName: opt.GetPoolName(),
		},
	}, nil
}

func (d *Driver) createVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	mgr := NewSrcMgr(d.conf)
	defer mgr.destroy()

	ioctx, err := mgr.GetIoctx(opt.GetPoolName())
	if err != nil {
		return nil, err
	}

	name := EncodeName(opt.GetId())
	_, err = rbd.Create(ioctx, name, uint64(opt.GetSize())<<sizeShiftBit, 20)
	if err != nil {
		log.Errorf("Create rbd image (%s) failed, (%v)", name, err)
		return nil, err
	}

	log.Infof("Create volume %s (%s) success.", opt.GetName(), opt.GetId())
	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.GetSize(),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
		Metadata: map[string]string{
			KPoolName: opt.GetPoolName(),
		},
	}, nil
}

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	// create a volume from snapshot
	if opt.GetSnapshotId() != "" {
		return d.createVolumeFromSnapshot(opt)
	}
	return d.createVolume(opt)
}

// ExtendVolume ...
func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	mgr := NewSrcMgr(d.conf)
	defer mgr.destroy()

	img, err := mgr.GetImage(opt.GetPoolName(), EncodeName(opt.GetId()))
	if err != nil {
		return nil, err
	}

	if err := img.Resize(uint64(opt.GetSize()) << sizeShiftBit); err != nil {
		log.Error("When resize image:", err)
		return nil, err
	}
	log.Info("Resize image success, volume id =", opt.GetId())

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

func (d *Driver) PullVolume(volID string) (*model.VolumeSpec, error) {
	// Not used, do nothing.
	return nil, nil
}

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	mgr := NewSrcMgr(d.conf)
	defer mgr.destroy()

	log.Info(opt.GetMetadata()[KPoolName], EncodeName(opt.GetId()))
	ioctx, err := mgr.GetIoctx(opt.GetMetadata()[KPoolName])
	if err != nil {
		return err
	}

	err = rbd.GetImage(ioctx, EncodeName(opt.GetId())).Remove()
	if err != nil && err != rbd.RbdErrorNotFound {
		log.Errorf("Remove volume(%s) filed, %v", opt.GetId(), err)
		return err
	}

	log.Infof("Remove volume (%s) success", opt.GetId())
	return nil
}

func (d *Driver) InitializeConnection(opt *pb.CreateVolumeAttachmentOpts) (*model.ConnectionInfo, error) {
	poolName, ok := opt.GetMetadata()[KPoolName]
	if !ok {
		err := errors.New("Failed to find poolName in volume metadata!")
		log.Error(err)
		return nil, err
	}
	return &model.ConnectionInfo{
		DriverVolumeType: RBDProtocol,
		ConnectionData: map[string]interface{}{
			"secret_type":  "ceph",
			"name":         poolName + "/" + opensdsPrefix + opt.GetVolumeId(),
			"cluster_name": "ceph",
			"hosts":        []string{opt.GetHostInfo().Host},
			"volume_id":    opt.GetVolumeId(),
			"access_mode":  "rw",
			"ports":        []string{"6789"},
		},
	}, nil
}

func (d *Driver) TerminateConnection(opt *pb.DeleteVolumeAttachmentOpts) error { return nil }

func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	mgr := NewSrcMgr(d.conf)
	defer mgr.destroy()

	poolName := opt.GetMetadata()[KPoolName]
	img, err := mgr.GetImage(poolName, EncodeName(opt.GetVolumeId()))
	if err != nil {
		return nil, err
	}

	if _, err := img.CreateSnapshot(EncodeName(opt.GetId())); err != nil {
		log.Error("When create snapshot:", err)
		return nil, err
	}

	log.Infof("Create snapshot (name:%s, id:%s, volID:%s) success",
		opt.GetName(), opt.GetId(), opt.GetVolumeId())

	return &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Description: opt.GetDescription(),
		VolumeId:    opt.GetVolumeId(),
		Size:        opt.GetSize(),
		Metadata: map[string]string{
			KPoolName:  poolName,
			KImageName: EncodeName(opt.GetVolumeId()),
		},
	}, nil

}

func (d *Driver) PullSnapshot(snapID string) (*model.VolumeSnapshotSpec, error) {
	return nil, fmt.Errorf("Ceph PullSnapshot has not implemented yet.")
}

func (d *Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	mgr := NewSrcMgr(d.conf)
	defer mgr.destroy()

	poolName := opt.GetMetadata()[KPoolName]
	img, err := mgr.GetImage(poolName, EncodeName(opt.GetVolumeId()), EncodeName(opt.GetId()))
	if err == rbd.RbdErrorNotFound {
		log.Warningf("Specified snapshot (%s) does not exist, ignore it", opt.GetId())
		return nil
	}
	if err != nil {
		return err
	}

	snap := img.GetSnapshot(EncodeName(opt.GetId()))
	if ok, _ := snap.IsProtected(); ok {
		if err := snap.Unprotect(); err != nil {
			log.Errorf("unprotect failed, %v", err)
			return err
		}
	}

	err = snap.Remove()
	if err != nil && err != rbd.RbdErrorNotFound {
		log.Error("When remove snapshot:", err)
		return err
	}

	log.Infof("Delete snapshot (%s) success", opt.GetId())
	return nil
}

type TotalStats struct {
	TotalBytes      int64 `json:"total_bytes,omitempty"`
	TotalUsedBytes  int64 `json:"total_used_bytes,omitempty"`
	TotalAvailBytes int64 `json:"total_avail_bytes,omitempty"`
}

type PoolStats struct {
	Name  string `json:"name,omitempty"`
	Id    int64  `json:"id,omitempty"`
	Stats struct {
		KbUsed      int64 `json:"kb_used,omitempty"`
		BytesUsed   int64 `json:"bytes_used,omitempty"`
		PercentUsed int64 `json:"percent_used,omitempty"`
		MaxAvail    int64 `json:"max_avail,omitempty"`
		Objects     int64 `json:"objects,omitempty"`
	} `json:"stats,omitempty"`
}

type DfInfo struct {
	Stats TotalStats  `json:"stats,omitempty"`
	Pools []PoolStats `json:"pools,omitempty"`
}

func (d *Driver) ListPools() ([]*model.StoragePoolSpec, error) {
	mgr := NewSrcMgr(d.conf)
	defer mgr.destroy()

	conn, err := mgr.GetConn()
	if err != nil {
		return nil, err
	}
	buf, info, err := conn.MonCommand([]byte(`{"prefix":"df", "format":"json"}`))
	if err != nil {
		log.Errorf("get mon df info filed, info: %s, err:%v", info, err)
		return nil, err
	}
	dfinfo := DfInfo{}
	json.Unmarshal([]byte(buf), &dfinfo)

	var pols []*model.StoragePoolSpec
	for _, p := range dfinfo.Pools {
		if _, ok := d.conf.Pool[p.Name]; !ok {
			continue
		}

		pol := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, p.Name).String(),
			},
			Name:             p.Name,
			TotalCapacity:    (p.Stats.BytesUsed + p.Stats.MaxAvail) >> sizeShiftBit,
			FreeCapacity:     p.Stats.MaxAvail >> sizeShiftBit,
			StorageType:      d.conf.Pool[p.Name].StorageType,
			Extras:           d.conf.Pool[p.Name].Extras,
			AvailabilityZone: d.conf.Pool[p.Name].AvailabilityZone,
			MultiAttach:      d.conf.Pool[p.Name].MultiAttach,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = defaultAZ
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

func (d *Driver) InitializeSnapshotConnection(opt *pb.CreateSnapshotAttachmentOpts) (*model.ConnectionInfo, error) {
	poolName, ok := opt.GetMetadata()[KPoolName]
	if !ok {
		err := errors.New("Failed to find poolName in snapshot metadata!")
		log.Error(err)
		return nil, err
	}
	return &model.ConnectionInfo{
		DriverVolumeType: RBDProtocol,
		ConnectionData: map[string]interface{}{
			"secret_type":  "ceph",
			"name":         poolName + "/" + opensdsPrefix + opt.GetSnapshotId(),
			"cluster_name": "ceph",
			"hosts":        []string{opt.GetHostInfo().Host},
			"volume_id":    opt.GetSnapshotId(),
			"access_mode":  "rw",
			"ports":        []string{"6789"},
		},
	}, nil
}

func (d *Driver) TerminateSnapshotConnection(opt *pb.DeleteSnapshotAttachmentOpts) error {
	return nil
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
