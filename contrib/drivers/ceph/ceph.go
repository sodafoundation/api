// Copyright 2017 The OpenSDS Authors.
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

/*
This module implements ceph driver for OpenSDS. Ceph driver will pass these
operation requests about volume to go-ceph module.
*/

package ceph

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/ceph/go-ceph/rados"
	"github.com/ceph/go-ceph/rbd"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/backup"
	"github.com/opensds/opensds/contrib/connector"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils"
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

// this function only used for open origin image rather than clone image or copy image
func (s *SrcMgr) GetOriginImage(poolName string, imgName string, args ...interface{}) (*rbd.Image, error) {
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

func (d *Driver) createVolumeFromSnapshot(opt *pb.CreateVolumeOpts) error {
	mgr := NewSrcMgr(d.conf)
	defer mgr.destroy()

	poolName := opt.GetPoolName()
	srcSnapName := EncodeName(opt.GetSnapshotId())
	srcImgName := opt.GetMetadata()[KImageName]
	destImgName := EncodeName(opt.GetId())

	img, err := mgr.GetOriginImage(poolName, srcImgName, srcSnapName)
	if err != nil {
		return err
	}
	snap := img.GetSnapshot(srcSnapName)
	if ok, _ := snap.IsProtected(); !ok {
		if err := snap.Protect(); err != nil {
			log.Errorf("protect snapshot failed, %v", err)
			return err
		}
		defer snap.Unprotect()
	}

	ioctx, err := mgr.GetIoctx(poolName)
	if err != nil {
		return err
	}

	destImg, err := img.Clone(srcSnapName, ioctx, destImgName, rbd.RbdFeatureLayering, 20)
	if err != nil {
		log.Errorf("snapshot clone failed:%v", err)
		return err
	}

	// flatten dest image
	if err := destImg.Open(); err != nil {
		log.Error("new image open failed:", err)
		return err
	}
	defer destImg.Close()
	if err := destImg.Flatten(); err != nil {
		log.Errorf("new image flatten failed, %v", err)
		return err
	}

	log.Infof("create volume (%s) from snapshot (%s) success", srcImgName, srcSnapName)
	return nil
}

func (d *Driver) createVolume(opt *pb.CreateVolumeOpts) error {
	mgr := NewSrcMgr(d.conf)
	defer mgr.destroy()

	ioctx, err := mgr.GetIoctx(opt.GetPoolName())
	if err != nil {
		return err
	}

	name := EncodeName(opt.GetId())
	_, err = rbd.Create(ioctx, name, uint64(opt.GetSize())<<sizeShiftBit, 20)
	if err != nil {
		log.Errorf("Create rbd image (%s) failed, (%v)", name, err)
		return err
	}

	log.Infof("Create volume %s (%s) success.", opt.GetName(), opt.GetId())
	return nil
}

func (d *Driver) createVolumeFromCloud(opt *pb.CreateVolumeOpts) error {

	if err := d.createVolume(opt); err != nil {
		log.Errorf("create image failed, %s", err)
		return err
	}
	if err := d.downloadSnapshotFromCloud(opt); err != nil {
		log.Errorf("create image failed, %s", err)
		// roll back
		d.deleteVolume(opt.GetPoolName(), opt.GetId(), nil)
		return err
	}
	return nil
}

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	var err error
	// create a volume from snapshot
	if opt.GetSnapshotId() != "" {
		if opt.SnapshotFromCloud {
			err = d.createVolumeFromCloud(opt)
		} else {
			err = d.createVolumeFromSnapshot(opt)
		}
	} else {
		err = d.createVolume(opt)
	}

	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.GetSize(),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
		Metadata:         map[string]string{KPoolName: opt.GetPoolName()},
	}, err

}

// ExtendVolume ...
func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	mgr := NewSrcMgr(d.conf)
	defer mgr.destroy()

	img, err := mgr.GetOriginImage(opt.GetPoolName(), EncodeName(opt.GetId()))
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

func (d *Driver) deleteVolume(poolName, volumeId string, mgr *SrcMgr) error {
	if mgr == nil {
		mgr = NewSrcMgr(d.conf)
		defer mgr.destroy()
	}

	log.Info(poolName, EncodeName(volumeId))
	ioctx, err := mgr.GetIoctx(poolName)
	if err != nil {
		return err
	}

	err = rbd.GetImage(ioctx, EncodeName(volumeId)).Remove()
	if err != nil && err != rbd.RbdErrorNotFound {
		log.Errorf("Remove volume(%s) filed, %v", volumeId, err)
		return err
	}

	log.Infof("Remove volume (%s) success", volumeId)
	return nil
}

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	return d.deleteVolume(opt.GetMetadata()[KPoolName], opt.GetId(), nil)
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
			"name":         poolName + "/" + EncodeName(opt.GetVolumeId()),
			"cluster_name": "ceph",
			"hosts":        []string{opt.GetHostInfo().Host},
			"volume_id":    opt.GetVolumeId(),
			"access_mode":  "rw",
			"ports":        []string{"6789"},
		},
	}, nil
}

func (d *Driver) TerminateConnection(opt *pb.DeleteVolumeAttachmentOpts) error { return nil }

func (d *Driver) uploadSnapshotToCloud(opt *pb.CreateVolumeSnapshotOpts, bucket string, mgr *SrcMgr) (map[string]string, error) {

	hostname, err := os.Hostname()
	if err != nil {
		log.Errorf("get host name filed, %v", err)
	}
	createOpt := &pb.CreateSnapshotAttachmentOpts{
		SnapshotId: opt.GetId(),
		Metadata:   opt.GetMetadata(),
		HostInfo: &pb.HostInfo{
			Platform:  runtime.GOARCH,
			OsType:    runtime.GOOS,
			Host:      hostname,
			Initiator: "",
		},
	}

	createOpt.Metadata[KImageName] = EncodeName(opt.GetVolumeId())
	info, err := d.InitializeSnapshotConnection(createOpt)
	if err != nil {
		return nil, err
	}
	defer d.TerminateConnection(&pb.DeleteVolumeAttachmentOpts{})

	log.Errorf("%v", info)
	conn := connector.NewConnector(info.DriverVolumeType)
	mountPoint, err := conn.Attach(info.ConnectionData)
	if err != nil {
		log.Errorf("attach image failed, %v", err)
		return nil, err
	}
	defer conn.Detach(info.ConnectionData)

	bk, err := backup.NewBackup("multi-cloud")
	if err != nil {
		log.Errorf("get backup driver, err: %v", err)
		return nil, err
	}

	if err := bk.SetUp(); err != nil {
		log.Errorf("backup driver setup failed:%v", err)
		return nil, err
	}
	defer bk.CleanUp()

	file, err := os.Open(mountPoint)
	if err != nil {
		log.Errorf("open lvm snapshot file, err: %v", err)
		return nil, err
	}
	defer file.Close()

	b := &backup.BackupSpec{
		Id:       uuid.NewV4().String(),
		Metadata: map[string]string{"bucket": bucket},
	}
	if err := bk.Backup(b, file); err != nil {
		log.Errorf("upload snapshot to multi-cloud failed, err: %v", err)
		return nil, err
	}

	return map[string]string{"backupId": b.Id, "bucket": bucket}, nil
}

func (d *Driver) downloadSnapshotFromCloud(opt *pb.CreateVolumeOpts) error {
	data := opt.GetMetadata()
	backupId, ok := data["backupId"]
	if !ok {
		return errors.New("can't find backupId in metadata")
	}
	bucket, ok := data["bucket"]
	if !ok {
		return errors.New("can't find bucket name in metadata")
	}

	createOpt := &pb.CreateVolumeAttachmentOpts{
		VolumeId: opt.GetId(),
		Metadata: opt.GetMetadata(),
		HostInfo: &pb.HostInfo{
			Platform:  runtime.GOARCH,
			OsType:    runtime.GOOS,
			Initiator: "",
		},
	}

	info, err := d.InitializeConnection(createOpt)
	if err != nil {
		return err
	}
	defer d.TerminateSnapshotConnection(&pb.DeleteSnapshotAttachmentOpts{})

	conn := connector.NewConnector(info.DriverVolumeType)
	mountPoint, err := conn.Attach(info.ConnectionData)
	if err != nil {
		return err
	}
	defer conn.Detach(info.ConnectionData)

	file, err := os.OpenFile(mountPoint, os.O_RDWR, 0666)
	if err != nil {
		log.Errorf("open lvm snapshot file, err: %v", err)
		return err
	}
	defer file.Close()

	bk, err := backup.NewBackup("multi-cloud")
	if err != nil {
		log.Errorf("get backup driver, err: %v", err)
		return err
	}
	if err := bk.SetUp(); err != nil {
		return err
	}
	defer bk.CleanUp()

	b := &backup.BackupSpec{
		Metadata: map[string]string{"bucket": bucket},
	}
	if err := bk.Restore(b, backupId, file); err != nil {
		log.Errorf("upload snapshot to multi-cloud failed, err: %v", err)
		return err
	}
	log.Infof("download snapshot(%s) from cloud bucket (%s) success", opt.SnapshotId, bucket)
	return nil
}

func (d *Driver) deleteUploadedSnapshot(backupId string, bucket string) error {
	bk, err := backup.NewBackup("multi-cloud")
	if err != nil {
		log.Errorf("get backup driver failed, err: %v", err)
		return err
	}

	if err := bk.SetUp(); err != nil {
		log.Errorf("backup driver setup failed:%v", err)
		return err
	}
	defer bk.CleanUp()

	b := &backup.BackupSpec{
		Id:       backupId,
		Metadata: map[string]string{"bucket": bucket},
	}

	if err := bk.Delete(b); err != nil {
		log.Errorf("delete backup snapshot  failed, err: %v", err)
		return err
	}
	return nil
}

func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	mgr := NewSrcMgr(d.conf)
	defer mgr.destroy()

	poolName := opt.GetMetadata()[KPoolName]
	img, err := mgr.GetOriginImage(poolName, EncodeName(opt.GetVolumeId()))
	if err != nil {
		return nil, err
	}

	if _, err := img.CreateSnapshot(EncodeName(opt.GetId())); err != nil {
		log.Error("When create snapshot:", err)
		return nil, err
	}

	var metadata = map[string]string{
		KPoolName:  poolName,
		KImageName: EncodeName(opt.GetVolumeId()),
	}

	// upload to cloud
	profile := model.NewProfileFromJson(opt.GetProfile())
	bucket := profile.SnapshotProperties.Topology.Bucket
	if len(bucket) != 0 {
		updateMetadata, err := d.uploadSnapshotToCloud(opt, bucket, mgr)
		if err != nil {
			// rollback
			d.deleteSnapshot(poolName, opt.GetVolumeId(), opt.GetId(), mgr)
			return nil, err
		}
		metadata = utils.MergeStringMaps(metadata, updateMetadata)
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
		Metadata:    metadata,
	}, nil

}

func (d *Driver) PullSnapshot(snapID string) (*model.VolumeSnapshotSpec, error) {
	return nil, fmt.Errorf("Ceph PullSnapshot has not implemented yet.")
}

func (d *Driver) deleteSnapshot(poolName, volumeId, snapshotId string, mgr *SrcMgr) error {
	if mgr == nil {
		mgr = NewSrcMgr(d.conf)
		defer mgr.destroy()
	}

	img, err := mgr.GetOriginImage(poolName, EncodeName(volumeId), EncodeName(snapshotId))
	if err == rbd.RbdErrorNotFound {
		log.Warningf("Specified snapshot (pool:%s,volume:%s,snapshot:%s) does not exist, ignore it",
			poolName, volumeId, snapshotId)
		return nil
	}
	if err != nil {
		return err
	}

	snap := img.GetSnapshot(EncodeName(snapshotId))
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

	return nil
}

func (d *Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	if bucket, ok := opt.Metadata["bucket"]; ok {
		log.Info("remove snapshot in cloud :", bucket)
		if err := d.deleteUploadedSnapshot(opt.Metadata["backupId"], bucket); err != nil {
			return err
		}
	}

	poolName := opt.GetMetadata()[KPoolName]
	if err := d.deleteSnapshot(poolName, opt.GetVolumeId(), opt.GetId(), nil); err != nil {
		log.Infof("Delete snapshot (%s) failed", opt.GetId())
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
		err := errors.New("Failed to find poolName in snapshot attachment metadata!")
		log.Error(err)
		return nil, err
	}

	imgName, ok := opt.GetMetadata()[KImageName]
	if !ok {
		err := errors.New("Failed to find imageName in snapshot attachment metadata!")
		log.Error(err)
		return nil, err
	}

	return &model.ConnectionInfo{
		DriverVolumeType: RBDProtocol,
		ConnectionData: map[string]interface{}{
			"secret_type":  "ceph",
			"name":         fmt.Sprintf("%s/%s@%s", poolName, imgName, EncodeName(opt.GetSnapshotId())),
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
