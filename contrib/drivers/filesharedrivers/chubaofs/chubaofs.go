// Copyright (c) 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package chubaofs

import (
	"errors"
	"fmt"
	"os"
	"path"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	. "github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils/config"
	uuid "github.com/satori/go.uuid"
)

const (
	DefaultConfPath = "/etc/opensds/driver/chubaofs.yaml"
	NamePrefix      = "chubaofs"
)

const (
	KMountPoint = "mountPoint"
	KVolumeName = "volName"
	KMasterAddr = "masterAddr"
	KLogDir     = "logDir"
	KWarnLogDir = "warnLogDir"
	KLogLevel   = "logLevel"
	KOwner      = "owner"
	KProfPort   = "profPort"
)

const (
	KClientPath = "clientPath"
)

const (
	defaultLogLevel = "error"
	defaultOwner    = "chubaofs"
	defaultProfPort = "10094"

	defaultVolumeCapLimit int64 = 1000000
)

const (
	clientConfigFileName = "client.json"
	clientCmdName        = "cfs-client"
)

type ClusterInfo struct {
	Name           string   `yaml:"name"`
	MasterAddr     []string `yaml:"masterAddr"`
	VolumeCapLimit int64    `yaml:"volumeCapLimit"`
}

type RuntimeEnv struct {
	MntPoint   string `yaml:"mntPoint"`
	ClientPath string `yaml:"clientPath"`
	LogLevel   string `yaml:"logLevel"`
	Owner      string `yaml:"owner"`
	ProfPort   string `yaml:"profPort"`
}

type Config struct {
	ClusterInfo `yaml:"clusterInfo"`
	RuntimeEnv  `yaml:"runtimeEnv"`

	Pool map[string]PoolProperties `yaml:"pool,flow"`
}

type Driver struct {
	conf *Config
}

func (d *Driver) Setup() error {
	conf := &Config{}
	path := config.CONF.OsdsDock.Backends.Chubaofs.ConfigPath
	if "" == path {
		path = DefaultConfPath
	}

	if _, err := Parse(conf, path); err != nil {
		return err
	}

	if conf.MntPoint == "" || conf.ClientPath == "" {
		return errors.New(fmt.Sprintf("chubaofs: lack of necessary config, mntPoint(%v) clientPath(%v)", conf.MntPoint, conf.ClientPath))
	}

	if conf.VolumeCapLimit <= 0 {
		conf.VolumeCapLimit = defaultVolumeCapLimit
	}

	d.conf = conf
	return nil
}

func (d *Driver) Unset() error {
	return nil
}

func (d *Driver) CreateFileShare(opt *pb.CreateFileShareOpts) (fshare *FileShareSpec, err error) {
	log.Info("CreateFileShare ...")

	volName := opt.GetId()
	volSize := opt.GetSize()

	configFiles, fsMntPoints, owner, err := prepareConfigFiles(d, opt)
	if err != nil {
		return nil, err
	}

	/*
	 * Only the master raft leader can repsonse to create volume requests.
	 */
	leader, err := getClusterInfo(d.conf.MasterAddr[0])
	if err != nil {
		return nil, err
	}

	err = createOrDeleteVolume(createVolumeRequest, leader, volName, owner, volSize)
	if err != nil {
		return nil, err
	}

	err = doMount(clientCmdName, configFiles)
	if err != nil {
		doUmount(fsMntPoints)
		createOrDeleteVolume(deleteVolumeRequest, leader, volName, owner, 0)
		return nil, err
	}

	log.Infof("Start client daemon successful: volume name: %v", volName)

	fshare = &FileShareSpec{
		BaseModel: &BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.GetSize(),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
		PoolId:           opt.GetPoolId(),
		ExportLocations:  fsMntPoints,
		Metadata: map[string]string{
			KVolumeName: volName,
			KClientPath: d.conf.ClientPath,
			KOwner:      owner,
		},
	}

	return fshare, nil
}

func (d *Driver) DeleteFileShare(opts *pb.DeleteFileShareOpts) error {
	volName := opts.GetMetadata()[KVolumeName]
	clientPath := opts.GetMetadata()[KClientPath]
	owner := opts.GetMetadata()[KOwner]
	fsMntPoints := make([]string, 0)
	fsMntPoints = append(fsMntPoints, opts.ExportLocations...)

	/*
	 * Umount export locations
	 */
	err := doUmount(fsMntPoints)
	if err != nil {
		return err
	}

	/*
	 * Remove generated mount points dir
	 */
	for _, mnt := range fsMntPoints {
		err = os.RemoveAll(mnt)
		if err != nil {
			return errors.New(fmt.Sprintf("chubaofs: failed to remove export locations, err: %v", err))
		}
	}

	/*
	 * Remove generated client runtime path
	 */
	err = os.RemoveAll(path.Join(clientPath, volName))
	if err != nil {
		return errors.New(fmt.Sprintf("chubaofs: failed to remove client path %v , volume name: %v , err: %v", clientPath, volName, err))
	}

	/*
	 * Only the master raft leader can repsonse to delete volume requests.
	 */
	leader, err := getClusterInfo(d.conf.MasterAddr[0])
	if err != nil {
		return err
	}
	err = createOrDeleteVolume(deleteVolumeRequest, leader, volName, owner, 0)
	return err
}

func (d *Driver) CreateFileShareSnapshot(opts *pb.CreateFileShareSnapshotOpts) (*FileShareSnapshotSpec, error) {
	return nil, &NotImplementError{"CreateFileShareSnapshot not implemented yet"}
}

func (d *Driver) DeleteFileShareSnapshot(opts *pb.DeleteFileShareSnapshotOpts) error {
	return &NotImplementError{"DeleteFileShareSnapshot not implemented yet"}
}

func (d *Driver) CreateFileShareAcl(opts *pb.CreateFileShareAclOpts) (*FileShareAclSpec, error) {
	return nil, &NotImplementError{"CreateFileShareAcl not implemented yet"}
}

func (d *Driver) DeleteFileShareAcl(opts *pb.DeleteFileShareAclOpts) error {
	return &NotImplementError{"DeleteFileShareAcl not implemented yet"}
}

func (d *Driver) ListPools() ([]*StoragePoolSpec, error) {
	pools := make([]*StoragePoolSpec, 0)
	for name, prop := range d.conf.Pool {
		pool := &StoragePoolSpec{
			BaseModel: &BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, name).String(),
			},
			Name:             name,
			TotalCapacity:    d.conf.VolumeCapLimit,
			FreeCapacity:     d.conf.VolumeCapLimit,
			StorageType:      prop.StorageType,
			Extras:           prop.Extras,
			AvailabilityZone: prop.AvailabilityZone,
		}
		if pool.AvailabilityZone == "" {
			pool.AvailabilityZone = "default"
		}
		pools = append(pools, pool)
	}
	return pools, nil
}
