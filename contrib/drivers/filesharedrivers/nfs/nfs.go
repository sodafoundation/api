// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package nfs

import (
	"errors"
	"path"
	"strings"

	log "github.com/golang/glog"
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
	defaultConfPath   = "/etc/opensds/driver/nfs.yaml"
	FileSharePrefix   = "fileshare-"
	snapshotPrefix    = "_snapshot-"
	blocksize         = 4096
	sizeShiftBit      = 30
	opensdsnvmepool   = "opensds-nvmegroup"
	nvmeofAccess      = "nvmeof"
	iscsiAccess       = "iscsi"
)

const (
	KLvPath            = "lvPath"
	KLvsPath           = "lvsPath"
	KFileshareName     = "nfsFileshareName"
	KFileshareID       = "nfsFileshareID"
	KFileshareSnapName = "snapshotName"
	KFileshareSnapID   = "snapshotID"
	AccessLevelRo      = "ro"
	AccessLevelRw      = "rw"
)

type NFSConfig struct {
	TgtBindIp      string                    `yaml:"tgtBindIp"`
	TgtConfDir     string                    `yaml:"tgtConfDir"`
	EnableChapAuth bool                      `yaml:"enableChapAuth"`
	Pool           map[string]PoolProperties `yaml:"pool,flow"`
}

type Driver struct {
	conf *NFSConfig
	cli  *Cli
}

func (d *Driver) Setup() error {
	// Read nfs config file
	d.conf = &NFSConfig{TgtBindIp: defaultTgtBindIp, TgtConfDir: defaultTgtConfDir}
	p := config.CONF.OsdsDock.Backends.NFS.ConfigPath
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

func (d *Driver) CreateFileShareAcl(opt *pb.CreateFileShareAclOpts) (*model.FileShareAclSpec, error) {
	var access string
	// Get accessto list
	accessTo := opt.GetAccessTo()
	// get accessCapability list
	accessCapability := opt.GetAccessCapability()
	// get fileshare name
	fname := opt.Name

	permissions := []string{"write"}
	WriteAccess := false

	for _, value := range accessCapability {
		value = strings.ToLower(value)
		if utils.Contains(permissions, value) {
			WriteAccess = true
		}
		if value == "Execute" {
			log.Error("invalid permission:", value)
			return nil, nil
		}
	}
	if WriteAccess {
		access = AccessLevelRw
	} else {
		access = AccessLevelRo
	}

	if err := d.cli.CreateAccess(accessTo, access, fname); err != nil {
		log.Errorf("grant access %s to %s failed %v", accessTo, fname, err)
		return nil, err
	}

	shareAccess := &model.FileShareAclSpec{
		BaseModel: &model.BaseModel{
			Id: opt.Id,
		},
		FileShareId:      opt.FileshareId,
		Type:             opt.Type,
		AccessCapability: accessCapability,
		AccessTo:         accessTo,
		Metadata:         map[string]string{},
	}
	return shareAccess, nil
}

func (d *Driver) DeleteFileShareAcl(opt *pb.DeleteFileShareAclOpts) error {
	// Get accessto list
	accessTo := opt.GetAccessTo()
	// get fileshare name
	fname := opt.Name

	if err := d.cli.DeleteAccess(accessTo, fname); err != nil {
		log.Error("cannot revoke access:", err)
		return err
	}

	return nil
}

func (d *Driver) CreateFileShare(opt *pb.CreateFileShareOpts) (*model.FileShareSpec, error) {
	var fshare *model.FileShareSpec
	//get the server ip for configuration
	var server = d.conf.TgtBindIp
	//get fileshare name
	var name = opt.GetName()
	//get volume group
	var vg = opt.GetPoolName()
	// Crete a directory to mount
	var dirName = path.Join("/var/", name)
	// create a fileshare path
	var lvPath = path.Join("/dev", vg, name)

	if err := d.cli.CreateDirectory(dirName); err != nil {
		log.Error("failed to create a directory:", err)
		return nil, err
	}

	if err := d.cli.CreateVolume(name, vg, opt.GetSize()); err != nil {
		return nil, err
	}
	// remove created volume if got error
	defer func() {
		// using return value as the error flag
		if fshare == nil {
			if err := d.cli.Delete(name, vg); err != nil {
				log.Error("failed to remove volume fileshare:", err)
			}
		}
	}()

	// Crete fileshare on this path
	if err := d.cli.CreateFileShare(lvPath); err != nil {
		log.Error("failed to create filesystem logic volume:", err)
		return nil, err
	}
	// mount the volume to directory
	if err := d.cli.Mount(lvPath, dirName); err != nil {
		log.Error("failed to mount a directory:", err)
		return nil, err
	}
	// Set permission to directory
	if err := d.cli.SetPermission(dirName); err != nil {
		log.Error("failed to set permission:", err)
		return nil, err
	}
	// get export location of fileshare
	var location []string
	location = []string{d.cli.GetExportLocation(name, server)}
	if len(location) == 0 {
		errMsg := errors.New("failed to get exportlocation: export location is empty!")
		log.Error(errMsg)
		return nil, errMsg
	}

	fshare = &model.FileShareSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.GetSize(),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
		PoolId:           opt.GetPoolId(),
		Protocols:        []string{NFSProtocol},
		ExportLocations:  location,
		Metadata: map[string]string{
			KFileshareName: name,
			KFileshareID:   opt.GetId(),
			KLvPath:        lvPath,
		},
	}
	return fshare, nil
}

// ListPools
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
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = "default"
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

// delete fileshare from device
func (d *Driver) DeleteFileShare(opt *pb.DeleteFileShareOpts) error {
	// get fileshare name to be deleted
	fname := opt.GetMetadata()[KFileshareName]
	if !d.cli.Exists(fname) {
		log.Warningf("fileshare(%s) does not exist, nothing to remove", fname)
		return nil
	}
	// get fileshare path
	lvPath := opt.GetMetadata()[KLvPath]
	// get directory where fileshare mounted
	var dirName = path.Join("/var/", fname)

	// umount the volume to directory
	if err := d.cli.UnMount(dirName); err != nil {
		log.Error("failed to mount a directory:", err)
		return err
	}
	// delete the actual fileshare from device
	if err := d.cli.Delete(fname, lvPath); err != nil {
		log.Error("failed to remove logic volume:", err)
		return err
	}
	// Delete the directory
	if err := d.cli.DeleteDirectory(dirName); err != nil {
		log.Error("failed to delete the directory:", err)
		return err
	}
	return nil
}

// CreateFileShareSnapshot
func (d *Driver) CreateFileShareSnapshot(opt *pb.CreateFileShareSnapshotOpts) (*model.FileShareSnapshotSpec, error) {
	lvPath, ok := opt.GetMetadata()[KLvPath]
	if !ok {
		err := errors.New("can't find 'lvPath' in snapshot metadata")
		log.Error(err)
		return nil, err
	}
	snapName := opt.GetName()
	fields := strings.Split(lvPath, "/")

	vg, sourceLvName := fields[2], fields[3]
	if err := d.cli.CreateLvSnapshot(snapName, sourceLvName, vg, opt.GetSize()); err != nil {
		log.Error("failed to create logic volume snapshot:", err)
		return nil, err
	}

	return &model.FileShareSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:         opt.GetName(),
		SnapshotSize: opt.GetSize(),
		Description:  opt.GetDescription(),
		Metadata: map[string]string{
			KFileshareSnapName: snapName,
			KFileshareSnapID:   opt.GetId(),
			KLvPath:            lvPath,
		},
	}, nil
}

// DeleteFileShareSnapshot
func (d *Driver) DeleteFileShareSnapshot(opt *pb.DeleteFileShareSnapshotOpts) error {
	lvsPath, ok := opt.GetMetadata()[KLvPath]
	snapName := opt.GetMetadata()[KFileshareSnapName]
	if !ok {
		err := errors.New("can't find 'lvsPath' in snapshot metadata, ingnore it!")
		log.Error(err)
		return nil
	}
	fields := strings.Split(lvsPath, "/")
	vg := fields[2]
	if !d.cli.Exists(snapName) {
		log.Warningf("Snapshot(%s) does not exist, nothing to remove", snapName)
		return nil
	}

	if err := d.cli.DeleteFileShareSnapshots(snapName, vg); err != nil {
		log.Error("failed to remove logic volume:", err)
		return err
	}
	return nil
}
