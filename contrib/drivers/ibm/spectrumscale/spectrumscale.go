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

package spectrumscale

import (
	"strconv"
	"strings"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils/config"
	uuid "github.com/satori/go.uuid"
)

const (
	defaultTgtConfDir = "/etc/tgt/conf.d"
	defaultTgtBindIp  = "127.0.0.1"
	username          = "root"
	password          = "ibm"
	port              = "2022"
	defaultConfPath   = "/etc/opensds/driver/ibm.yaml"
	volumePrefix      = "volume-"
	snapshotPrefix    = "_snapshot-"
	blocksize         = 4096
	sizeShiftBit      = 30
	opensdsnvmepool   = "opensds-nvmegroup"
	nvmeofAccess      = "nvmeof"
	iscsiAccess       = "iscsi"
	storageType       = "block"
  timeoutForssh        = 60
)

const (
	KLvIdFormat = "NAA"
	FileSetPath = "FilesetPath"
	SnapshotName = "SnapshotName"
)

type IBMConfig struct {
	TgtBindIp      string                    `yaml:"tgtBindIp"`
	UserName       string                    `yaml:"username"`
	Password       string                    `yaml:"password"`
	Port           string                    `yaml:"port"`
	TgtConfDir     string                    `yaml:"tgtConfDir"`
	EnableChapAuth bool                      `yaml:"enableChapAuth"`
	Pool           map[string]PoolProperties `yaml:"pool,flow"`
}

type Driver struct {
	conf *IBMConfig
	cli  *Cli
}

func (d *Driver) Setup() error {
	// Read ibm config file
	d.conf = &IBMConfig{
	TgtBindIp: defaultTgtBindIp,
	TgtConfDir: defaultTgtConfDir,
	UserName: username,
	Port: port,
	Password: password,
	}
	p := config.CONF.OsdsDock.Backends.IBMSpectrumScale.ConfigPath
	if "" == p {
		p = defaultConfPath
	}
	if _, err := Parse(d.conf, p); err != nil {
		return err
	}
	err:= login()
	if err != nil {
		return err
	}

	return nil
}

func (*Driver) Unset() error { return nil }

// first get the status of spectrumstate. If it is not active just return
func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (vol *model.VolumeSpec, err error) {
		 err = d.cli.GetSpectrumScaleStatus()
		 if err != nil{
			 log.Error("the GPFS cluster is not active")
			 return &model.VolumeSpec{}, err
		 }

		 // if spectrumscale service is active, get the mountPoint and filesystem
     var mountPoint, filesystem string
 		 mountPoint, filesystem, err = d.cli.GetSpectrumScaleMountPoint()
 		 if err != nil{
 			log.Error("not able to find spectrumscale mount point")
 			return &model.VolumeSpec{}, err
 		 }
     log.Infof("the cluster filesystem name:%v and mounpoint is:%v", filesystem, mountPoint)


		 log.Info("IBM driver receive create volume request, vr =", opt)
		 var volName = volumePrefix + opt.GetId()
		 var volSize = opt.GetSize()
		 size := strconv.FormatInt(int64(volSize), 10)
		 if err = d.cli.CreateVolume(volName, size); err != nil {
			return &model.VolumeSpec{}, err
		}

		return &model.VolumeSpec{
			BaseModel: &model.BaseModel{
				Id: opt.GetId(),
			},
			Name:        opt.GetName(),
			Size:        opt.GetSize(),
			Description: opt.GetDescription(),
			Identifier:  &model.Identifier{DurableName:opt.GetId(),
				DurableNameFormat: KLvIdFormat,
			},
			Metadata: map[string]string{
				FileSetPath: mountPoint + "/" + volName,
			},
		}, nil
}

// discover the pool from spectrumscale
func (d *Driver) ListPools() ([]*model.StoragePoolSpec, error) {
	var mountPoint,filesystem string
	mountPoint, filesystem, stderr := d.cli.GetSpectrumScaleMountPoint()
	if stderr != nil {
    log.Error("failed to get mountpoint")
		return nil, stderr
	}
  log.Infof("the cluster filesystem name:%v and mounpoint is:%v", filesystem, mountPoint)

	pools, err := d.cli.ListPools(mountPoint, filesystem)
	if err != nil {
		return nil, err
	}
	// retrive the all details from spectrumscale pool
	var pols []*model.StoragePoolSpec
	for _, pool := range *pools {
		pol := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, pool.UUID).String(),
			},
			Name:             pool.Name,
			TotalCapacity:    pool.TotalCapacity,
			FreeCapacity:     pool.FreeCapacity,
			StorageType:      storageType,
			Extras:           d.conf.Pool[pool.Name].Extras,
			AvailabilityZone: d.conf.Pool[pool.Name].AvailabilityZone,
			MultiAttach:      d.conf.Pool[pool.Name].MultiAttach,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = "default"
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

// this function is for deleting the spectrumscale volume(fileset)
func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error{
	fileSetPath:= opt.GetMetadata()[FileSetPath]
	field := strings.Split(fileSetPath, "/")
	name := field[3]
	if err := d.cli.Delete(name); err != nil {
		log.Error("failed to remove logic volume:", err)
		return err
	}
	log.Info("volume is successfully deleted!")
  return nil
}

// this function is for extending the volume(fileset). It sets the quota for block and files
func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	fileSetPath:= opt.GetMetadata()[FileSetPath]
	field := strings.Split(fileSetPath, "/")
	name := field[3]
	var volsize = opt.GetSize()
	size := strconv.FormatInt(int64(volsize), 10)
	if err := d.cli.ExtendVolume(name, size); err != nil {
		log.Error("failed to extend the volume:", err)
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

// this function is for creating the snapshot of spectrumscale volume(fileset)
func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	fileSetPath:= opt.GetMetadata()[FileSetPath]
	field := strings.Split(fileSetPath, "/")
	volName := field[3]
	var snapName = opt.GetName()
	if err := d.cli.CreateSnapshot(snapName, volName); err != nil {
		log.Error("failed to create snapshot for volume:", err)
		return nil, err
	}
	return &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Size:        opt.GetSize(),
		Description: opt.GetDescription(),
		VolumeId:    opt.GetVolumeId(),
		Metadata:    map[string]string{
			FileSetPath: fileSetPath,
			SnapshotName: snapName,
		},
	}, nil
}

// this function is for deleting the snapshot
func (d *Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error{
	fileSetPath:= opt.GetMetadata()[FileSetPath]
	field := strings.Split(fileSetPath, "/")
	volName := field[3]
	snapName := opt.GetMetadata()[SnapshotName]
	if err := d.cli.DeleteSnapshot(volName, snapName); err != nil {
		log.Error("failed to delete the snapshot:", err)
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

func (d *Driver) PullVolume(volIdentifier string) (*model.VolumeSpec, error) {
        return nil, &model.NotImplementError{"method CreateVolumeGroup has not been implemented yet"}
}

func (d *Driver) PullSnapshot(snapIdentifier string) (*model.VolumeSnapshotSpec, error) {
        return nil, &model.NotImplementError{"method CreateVolumeGroup has not been implemented yet"}
}

func (d *Driver) InitializeSnapshotConnection(opt *pb.CreateSnapshotAttachmentOpts) (*model.ConnectionInfo, error) {
        return nil, &model.NotImplementError{"method CreateVolumeGroup has not been implemented yet"}
}

func (d *Driver) TerminateSnapshotConnection(opt *pb.DeleteSnapshotAttachmentOpts) error{
        return nil
}

func (d *Driver) InitializeConnection(opt *pb.CreateVolumeAttachmentOpts) (*model.ConnectionInfo, error) {
        return nil, &model.NotImplementError{"method CreateVolumeGroup has not been implemented yet"}
}

func (d *Driver) TerminateConnection(opt *pb.DeleteVolumeAttachmentOpts) error{
        return nil
}
