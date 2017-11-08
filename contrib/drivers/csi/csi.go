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
This module implements a sample driver for OpenSDS. This driver will handle all
operations of volume and return a fake value.

*/

package csi

import (
	"errors"
	"fmt"
	"io/ioutil"

	csipb "github.com/container-storage-interface/spec/lib/go/csi"
	log "github.com/golang/glog"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"

	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/config"
)

var conf = CSIConfig{}

type Driver struct {
	Client

	config CSIConfig
}

type CSIConfig struct {
	AuthOptions struct {
		IdentityEndpoint string `yaml:"endpoint,omitempty"`
	} `yaml:"authOptions"`
	Pool map[string]PoolProperties `yaml:"pool,flow"`
}

type PoolProperties struct {
	DiskType  string `yaml:"diskType"`
	IOPS      int64  `yaml:"iops"`
	BandWidth int64  `yaml:"bandwidth"`
}

func (d *Driver) Setup() error {
	d.Client = NewClient(d.config.AuthOptions.IdentityEndpoint)

	// Read csi config file
	confYaml, err := ioutil.ReadFile(config.CONF.CSIConfig)
	if err != nil {
		log.Fatalf("Read csi config yaml file (%s) failed, reason:(%v)", config.CONF.CSIConfig, err)
		return err
	}
	if err = yaml.Unmarshal(confYaml, &conf); err != nil {
		log.Fatal("Parse error: %v", err)
		return err
	}
	d.config = conf

	return nil
}

func (d *Driver) Unset() error {
	d.Close()
	return nil
}

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	in := &csipb.CreateVolumeRequest{
		Name: opt.GetName(),
		CapacityRange: &csipb.CapacityRange{
			RequiredBytes: uint64(opt.Size * 10 << 5),
		},
		VolumeCapabilities: []*csipb.VolumeCapability{
			{
				AccessType: &csipb.VolumeCapability_Block{},
			},
		},
	}

	response, err := d.Client.CreateVolume(context.Background(), in)
	if err != nil {
		log.Error("create volume failed in volume controller:", err)
		return nil, err
	}
	defer d.Client.Close()

	if errorMsg := response.GetError().GetCreateVolumeError(); errorMsg != nil {
		return nil,
			fmt.Errorf("failed to create volume in csi driver, code: %v, message: %v",
				errorMsg.GetErrorCode(), errorMsg.GetErrorDescription())
	}
	vol := response.GetResult().GetVolumeInfo()

	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: vol.GetId(),
		},
		Size:     int64(vol.GetCapacityBytes()),
		Metadata: vol.GetAttributes(),
	}, nil
}

func (d *Driver) PullVolume(volIdentifier string) (*model.VolumeSpec, error) {
	in := &csipb.ListVolumesRequest{}

	response, err := d.Client.ListVolumes(context.Background(), in)
	if err != nil {
		log.Error("List volumes failed in csi driver:", err)
		return nil, err
	}
	defer d.Client.Close()

	if errorMsg := response.GetError().GetGeneralError(); errorMsg != nil {
		return nil,
			fmt.Errorf("failed to list volumes in csi driver, code: %v, message: %v",
				errorMsg.GetErrorCode(), errorMsg.GetErrorDescription())
	}
	ents := response.GetResult().GetEntries()

	for _, ent := range ents {
		vol := ent.GetVolumeInfo()
		if vol.GetId() == volIdentifier {
			return &model.VolumeSpec{
				BaseModel: &model.BaseModel{
					Id: vol.GetId(),
				},
				Size:     int64(vol.GetCapacityBytes()),
				Metadata: vol.GetAttributes(),
			}, nil
		}
	}

	return nil, fmt.Errorf("failed to get volume in csi driver")
}

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	in := &csipb.DeleteVolumeRequest{
		VolumeId: opt.GetId(),
		UserCredentials: &csipb.Credentials{
			Data: opt.GetMetadata(),
		},
	}

	response, err := d.Client.DeleteVolume(context.Background(), in)
	if err != nil {
		log.Error("delete volume failed in csi driver:", err)
		return err
	}
	defer d.Client.Close()

	if errorMsg := response.GetError().GetDeleteVolumeError(); errorMsg != nil {
		return fmt.Errorf("failed to delete volume in csi driver, code: %v, message: %v",
			errorMsg.GetErrorCode(), errorMsg.GetErrorDescription())
	}

	return nil
}

func (*Driver) InitializeConnection(opt *pb.CreateAttachmentOpts) (*model.ConnectionInfo, error) {
	return nil, errors.New("Not implemented!")
}

func (*Driver) TerminateConnection(opt *pb.DeleteAttachmentOpts) error { return nil }

func (*Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (*Driver) PullSnapshot(snapIdentifier string) (*model.VolumeSnapshotSpec, error) {
	return nil, errors.New("Not implemented!")
}

func (*Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	return errors.New("Not implemented!")
}

func (d *Driver) ListPools() ([]*model.StoragePoolSpec, error) {
	var pols []*model.StoragePoolSpec

	for name := range d.config.Pool {
		var param = func(proper PoolProperties) map[string]interface{} {
			var param = make(map[string]interface{})
			param["diskType"] = proper.DiskType
			param["iops"] = proper.IOPS
			param["bandwidth"] = proper.BandWidth
			return param
		}(d.config.Pool[name])
		pol := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, name).String(),
			},
			Name:       name,
			Parameters: param,
		}
		pols = append(pols, pol)
	}

	return pols, nil
}
