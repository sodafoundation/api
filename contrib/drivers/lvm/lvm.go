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

package lvm

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/backup"
	"github.com/opensds/opensds/contrib/connector"
	"github.com/opensds/opensds/contrib/drivers/lvm/targets"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/config"
	uuid "github.com/satori/go.uuid"
)

const (
	defaultTgtConfDir = "/etc/tgt/conf.d"
	defaultTgtBindIp  = "127.0.0.1"
	defaultConfPath   = "/etc/opensds/driver/lvm.yaml"
	volumePrefix      = "volume-"
	snapshotPrefix    = "_snapshot-"
	blocksize         = 4096
	sizeShiftBit      = 30
	opensdsnvmepool   = "opensds-nvmegroup"
	nvmeofAccess      = "nvmeof"
	iscsiAccess       = "iscsi"
)

const (
	KLvPath  = "lvPath"
	KLvsPath = "lvsPath"
)

type LVMConfig struct {
	TgtBindIp      string                    `yaml:"tgtBindIp"`
	TgtConfDir     string                    `yaml:"tgtConfDir"`
	EnableChapAuth bool                      `yaml:"enableChapAuth"`
	Pool           map[string]PoolProperties `yaml:"pool,flow"`
}

type Driver struct {
	conf *LVMConfig
	cli  *Cli
}

func (d *Driver) Setup() error {
	// Read lvm config file
	d.conf = &LVMConfig{TgtBindIp: defaultTgtBindIp, TgtConfDir: defaultTgtConfDir}
	p := config.CONF.OsdsDock.Backends.LVM.ConfigPath
	if "" == p {
		p = defaultConfPath
	}
	if _, err := Parse(d.conf, p); err != nil {
		return err
	}
	cli, err := NewCli()
	if err != nil {
		return err
	}
	d.cli = cli

	return nil
}

func (*Driver) Unset() error { return nil }

func (d *Driver) downloadSnapshot(bucket, backupId, dest string) error {
	mc, err := backup.NewBackup("multi-cloud")
	if err != nil {
		log.Errorf("get backup driver, err: %v", err)
		return err
	}

	if err := mc.SetUp(); err != nil {
		return err
	}
	defer mc.CleanUp()

	file, err := os.OpenFile(dest, os.O_RDWR, 0666)
	if err != nil {
		log.Errorf("open lvm snapshot file, err: %v", err)
		return err
	}
	defer file.Close()

	metadata := map[string]string{
		"bucket": bucket,
	}
	b := &backup.BackupSpec{
		Metadata: metadata,
	}

	if err := mc.Restore(b, backupId, file); err != nil {
		log.Errorf("upload snapshot to multi-cloud failed, err: %v", err)
		return err
	}
	return nil
}

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (vol *model.VolumeSpec, err error) {
	var name = volumePrefix + opt.GetId()
	var vg = opt.GetPoolName()
	if err = d.cli.CreateVolume(name, vg, opt.GetSize()); err != nil {
		return
	}

	// remove created volume if got error
	defer func() {
		// using return value as the error flag
		if vol == nil {
			if err := d.cli.Delete(name, vg); err != nil {
				log.Error("Failed to remove logic volume:", err)
			}
		}
	}()

	var lvPath = path.Join("/dev", vg, name)
	// Create volume from snapshot
	if opt.GetSnapshotId() != "" {
		if opt.SnapshotFromCloud {
			// download cloud snapshot to volume
			data := opt.GetMetadata()
			backupId, ok := data["backupId"]
			if !ok {
				return nil, errors.New("can't find backupId in metadata")
			}
			bucket, ok := data["bucket"]
			if !ok {
				return nil, errors.New("can't find bucket name in metadata")
			}
			err := d.downloadSnapshot(bucket, backupId, lvPath)
			if err != nil {
				log.Errorf("Download snapshot failed, %v", err)
				return nil, err
			}
		} else {
			// copy local snapshot to volume
			var lvsPath = path.Join("/dev", vg, snapshotPrefix+opt.GetSnapshotId())
			if err := d.cli.CopyVolume(lvsPath, lvPath, opt.GetSize()); err != nil {
				log.Error("Failed to create logic volume:", err)
				return nil, err
			}
		}
	}

	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Size:        opt.GetSize(),
		Description: opt.GetDescription(),
		Metadata: map[string]string{
			KLvPath: lvPath,
		},
	}, nil
}

func (d *Driver) PullVolume(volIdentifier string) (*model.VolumeSpec, error) {
	// Not used , do nothing
	return nil, nil
}

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {

	var name = volumePrefix + opt.GetId()
	if !d.cli.Exists(name) {
		log.Warningf("Volume(%s) does not exist, nothing to remove", name)
		return nil
	}

	lvPath, ok := opt.GetMetadata()[KLvPath]
	if !ok {
		err := errors.New("can't find 'lvPath' in volume metadata")
		log.Error(err)
		return err
	}

	field := strings.Split(lvPath, "/")
	vg := field[2]
	if d.cli.LvHasSnapshot(name, vg) {
		err := fmt.Errorf("unable to delete due to existing snapshot for volume: %s", name)
		log.Error(err)
		return err
	}

	if err := d.cli.Delete(name, vg); err != nil {
		log.Error("Failed to remove logic volume:", err)
		return err
	}

	return nil
}

// ExtendVolume ...
func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	var name = volumePrefix + opt.GetId()
	if err := d.cli.ExtendVolume(name, opt.GetPoolName(), opt.GetSize()); err != nil {
		log.Errorf("extend volume(%s) failed, error: %v", name, err)
		return nil, err
	}
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
	log.V(8).Infof("lvm initialize connection information: %v", opt)
	initiator := opt.HostInfo.GetInitiator()
	if initiator == "" {
		initiator = "ALL"
	}

	hostIP := opt.HostInfo.GetIp()
	if hostIP == "" {
		hostIP = "ALL"
	}

	lvPath, ok := opt.GetMetadata()[KLvPath]
	if !ok {
		err := errors.New("can't find 'lvPath' in volume metadata")
		log.Error(err)
		return nil, err
	}
	var chapAuth []string
	if d.conf.EnableChapAuth {
		chapAuth = []string{utils.RandSeqWithAlnum(20), utils.RandSeqWithAlnum(16)}
	}

	// create target according to the pool's access protocol
	accPro := opt.AccessProtocol
	log.Info("accpro:", accPro)
	t := targets.NewTarget(d.conf.TgtBindIp, d.conf.TgtConfDir, accPro)
	expt, err := t.CreateExport(opt.GetVolumeId(), lvPath, hostIP, initiator, chapAuth)
	if err != nil {
		log.Error("Failed to initialize connection of logic volume:", err)
		return nil, err
	}

	log.V(8).Infof("lvm ConnectionData: %v", expt)

	return &model.ConnectionInfo{
		DriverVolumeType: accPro,
		ConnectionData:   expt,
	}, nil
}

func (d *Driver) TerminateConnection(opt *pb.DeleteVolumeAttachmentOpts) error {
	log.V(8).Infof("TerminateConnection: opt info is %v", opt)
	accPro := opt.AccessProtocol
	t := targets.NewTarget(d.conf.TgtBindIp, d.conf.TgtConfDir, accPro)
	if err := t.RemoveExport(opt.GetVolumeId(), opt.GetHostInfo().GetIp()); err != nil {
		log.Error("failed to terminate connection of logic volume:", err)
		return err
	}
	return nil
}

func (d *Driver) AttachSnapshot(snapshotId string, lvsPath string) (string, *model.ConnectionInfo, error) {

	var err error
	createOpt := &pb.CreateSnapshotAttachmentOpts{
		SnapshotId: snapshotId,
		Metadata: map[string]string{
			KLvsPath: lvsPath,
		},
		HostInfo: &pb.HostInfo{
			Platform:  runtime.GOARCH,
			OsType:    runtime.GOOS,
			Host:      d.conf.TgtBindIp,
			Initiator: "",
		},
	}

	info, err := d.InitializeSnapshotConnection(createOpt)
	if err != nil {
		return "", nil, err
	}

	// rollback
	defer func() {
		if err != nil {
			deleteOpt := &pb.DeleteSnapshotAttachmentOpts{}
			d.TerminateSnapshotConnection(deleteOpt)
		}
	}()

	conn := connector.NewConnector(info.DriverVolumeType)
	mountPoint, err := conn.Attach(info.ConnectionData)
	if err != nil {
		return "", nil, err
	}
	log.V(8).Infof("Attach snapshot success, MountPoint:%s", mountPoint)
	return mountPoint, info, nil
}

func (d *Driver) DetachSnapshot(snapshotId string, info *model.ConnectionInfo) error {

	con := connector.NewConnector(info.DriverVolumeType)
	if con == nil {
		return fmt.Errorf("Can not find connector (%s)!", info.DriverVolumeType)
	}

	con.Detach(info.ConnectionData)
	attach := &pb.DeleteSnapshotAttachmentOpts{
		SnapshotId:     snapshotId,
		AccessProtocol: info.DriverVolumeType,
	}
	return d.TerminateSnapshotConnection(attach)
}

func (d *Driver) uploadSnapshot(lvsPath string, bucket string) (string, error) {
	mc, err := backup.NewBackup("multi-cloud")
	if err != nil {
		log.Errorf("get backup driver, err: %v", err)
		return "", err
	}

	if err := mc.SetUp(); err != nil {
		return "", err
	}
	defer mc.CleanUp()

	file, err := os.Open(lvsPath)
	if err != nil {
		log.Errorf("open lvm snapshot file, err: %v", err)
		return "", err
	}
	defer file.Close()

	metadata := map[string]string{
		"bucket": bucket,
	}
	b := &backup.BackupSpec{
		Id:       uuid.NewV4().String(),
		Metadata: metadata,
	}

	if err := mc.Backup(b, file); err != nil {
		log.Errorf("upload snapshot to multi-cloud failed, err: %v", err)
		return "", err
	}
	return b.Id, nil
}

func (d *Driver) deleteUploadedSnapshot(backupId string, bucket string) error {
	mc, err := backup.NewBackup("multi-cloud")
	if err != nil {
		log.Errorf("get backup driver failed, err: %v", err)
		return err
	}

	if err := mc.SetUp(); err != nil {
		return err
	}
	defer mc.CleanUp()

	metadata := map[string]string{
		"bucket": bucket,
	}
	b := &backup.BackupSpec{
		Id:       backupId,
		Metadata: metadata,
	}
	if err := mc.Delete(b); err != nil {
		log.Errorf("delete backup snapshot  failed, err: %v", err)
		return err
	}
	return nil
}

func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (snap *model.VolumeSnapshotSpec, err error) {
	var snapName = snapshotPrefix + opt.GetId()

	lvPath, ok := opt.GetMetadata()[KLvPath]
	if !ok {
		err := errors.New("can't find 'lvPath' in snapshot metadata")
		log.Error(err)
		return nil, err
	}

	fields := strings.Split(lvPath, "/")
	vg, sourceLvName := fields[2], fields[3]
	if err := d.cli.CreateLvSnapshot(snapName, sourceLvName, vg, opt.GetSize()); err != nil {
		log.Error("Failed to create logic volume snapshot:", err)
		return nil, err
	}

	lvsPath := path.Join("/dev", vg, snapName)
	metadata := map[string]string{KLvsPath: lvsPath}

	if bucket, ok := opt.Metadata["bucket"]; ok {
		//nvmet right now can not support snap volume serve as nvme target
		if vg == opensdsnvmepool {
			log.Infof("nvmet right now can not support snap volume serve as nvme target")
			log.Infof("still store in nvme pool but initialize connection by iscsi protocol")
		}
		mountPoint, info, err := d.AttachSnapshot(opt.GetId(), lvsPath)
		if err != nil {
			d.cli.Delete(snapName, vg)
			return nil, err
		}
		defer d.DetachSnapshot(opt.GetId(), info)

		log.Info("update load snapshot to :", bucket)
		backupId, err := d.uploadSnapshot(mountPoint, bucket)
		if err != nil {
			d.cli.Delete(snapName, vg)
			return nil, err
		}
		metadata["backupId"] = backupId
		metadata["bucket"] = bucket
	}

	return &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Size:        opt.GetSize(),
		Description: opt.GetDescription(),
		VolumeId:    opt.GetVolumeId(),
		Metadata:    metadata,
	}, nil
}

func (d *Driver) PullSnapshot(snapIdentifier string) (*model.VolumeSnapshotSpec, error) {
	// not used, do nothing
	return nil, nil
}

func (d *Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {

	if bucket, ok := opt.Metadata["bucket"]; ok {
		log.Info("remove snapshot in multi-cloud :", bucket)
		if err := d.deleteUploadedSnapshot(opt.Metadata["backupId"], bucket); err != nil {
			return err
		}
	}

	lvsPath, ok := opt.GetMetadata()[KLvsPath]
	if !ok {
		err := errors.New("can't find 'lvsPath' in snapshot metadata, ingnore it!")
		log.Error(err)
		return nil
	}
	fields := strings.Split(lvsPath, "/")
	vg, snapName := fields[2], fields[3]
	if !d.cli.Exists(snapName) {
		log.Warningf("Snapshot(%s) does not exist, nothing to remove", snapName)
		return nil
	}

	if err := d.cli.Delete(snapName, vg); err != nil {
		log.Error("Failed to remove logic volume:", err)
		return err
	}

	return nil
}

func (d *Driver) ListPools() ([]*model.StoragePoolSpec, error) {

	vgs, err := d.cli.ListVgs()
	if err != nil {
		return nil, err
	}

	var pols []*model.StoragePoolSpec
	for _, vg := range *vgs {
		if _, ok := d.conf.Pool[vg.Name]; !ok {
			continue
		}

		pol := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, vg.UUID).String(),
			},
			Name:             vg.Name,
			TotalCapacity:    vg.TotalCapacity,
			FreeCapacity:     vg.FreeCapacity,
			StorageType:      d.conf.Pool[vg.Name].StorageType,
			Extras:           d.conf.Pool[vg.Name].Extras,
			AvailabilityZone: d.conf.Pool[vg.Name].AvailabilityZone,
			MultiAttach:      d.conf.Pool[vg.Name].MultiAttach,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = "default"
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

func (d *Driver) InitializeSnapshotConnection(opt *pb.CreateSnapshotAttachmentOpts) (*model.ConnectionInfo, error) {
	initiator := opt.HostInfo.GetInitiator()
	if initiator == "" {
		initiator = "ALL"
	}

	hostIP := opt.HostInfo.GetIp()
	if hostIP == "" {
		hostIP = "ALL"
	}

	lvsPath, ok := opt.GetMetadata()[KLvsPath]
	if !ok {
		err := errors.New("Failed to find logic volume path in volume attachment metadata!")
		log.Error(err)
		return nil, err
	}
	var chapAuth []string
	if d.conf.EnableChapAuth {
		chapAuth = []string{utils.RandSeqWithAlnum(20), utils.RandSeqWithAlnum(16)}
	}

	accPro := opt.AccessProtocol
	if accPro == nvmeofAccess {
		log.Infof("nvmet right now can not support snap volume serve as nvme target")
		log.Infof("still create snapshot connection by iscsi")
		accPro = iscsiAccess
	}
	t := targets.NewTarget(d.conf.TgtBindIp, d.conf.TgtConfDir, accPro)
	data, err := t.CreateExport(opt.GetSnapshotId(), lvsPath, hostIP, initiator, chapAuth)
	if err != nil {
		log.Error("Failed to initialize snapshot connection of logic volume:", err)
		return nil, err
	}

	return &model.ConnectionInfo{
		DriverVolumeType: accPro,
		ConnectionData:   data,
	}, nil
}

func (d *Driver) TerminateSnapshotConnection(opt *pb.DeleteSnapshotAttachmentOpts) error {
	accPro := opt.AccessProtocol
	if accPro == nvmeofAccess {
		log.Infof("nvmet right now can not support snap volume serve as nvme target")
		log.Infof("still create snapshot connection by iscsi")
		accPro = iscsiAccess
	}
	log.Info("terminate snapshot conn")
	t := targets.NewTarget(d.conf.TgtBindIp, d.conf.TgtConfDir, accPro)
	if err := t.RemoveExport(opt.GetSnapshotId(), opt.GetHostInfo().GetIp()); err != nil {
		log.Error("Failed to terminate snapshot connection of logic volume:", err)
		return err
	}
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
