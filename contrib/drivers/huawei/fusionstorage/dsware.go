// Copyright (c) 2019 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package main

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	. "github.com/opensds/opensds/pkg/model"
	"github.com/satori/go.uuid"
)

type Driver struct {
	cli  *RestCommon
	conf *Config
}

type AuthOptions struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Url      string `yaml:"url"`
}

type Config struct {
	AuthOptions `yaml:"authOptions"`
	Pool        map[string]PoolProperties `yaml:"pool,flow"`
}

func (d *Driver) Setup() error {
	conf := &Config{}

	d.conf = conf

	path := "./testdata/fusionstorage.yaml"
	if path == "" {
		path = ""
	}

	Parse(conf, path)

	client := newRestCommon(conf.Username, conf.Password, conf.Url)

	err := client.login()
	if err != nil {
		fmt.Printf("Get new client failed, %v", err)
		return err
	}

	d.cli = client

	return nil
}

func (d *Driver) Unset() error {
	return nil
}

func EncodeName(id string) string {
	return NamePrefix + "-" + id
}

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*VolumeSpec, error) {

	name := EncodeName(opt.GetId())
	err := d.cli.createVolume(name, opt.GetPoolName(), opt.GetSize()<<UnitGiShiftBit)
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

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	name := EncodeName(opt.GetId())
	err := d.cli.deleteVolume(name)
	if err != nil {
		log.Errorf("Delete volume (%s) failed: %v", opt.GetId(), err)
		return err
	}
	log.Infof("Delete volume (%s) success.", opt.GetId())
	return nil
}

func (d *Driver) ListPools() ([]*StoragePoolSpec, error) {
	var pols []*StoragePoolSpec
	pools, err := d.cli.queryPoolInfo()
	if err != nil {
		return nil, err
	}

	c := d.conf
	for _, p := range pools.Pools {
		poolId := strconv.Itoa(p.PoolId)
		if _, ok := c.Pool[poolId]; !ok {
			continue
		}
		host, _ := os.Hostname()
		name := fmt.Sprintf("%s:%s:%s", host, d.conf.Url, poolId)
		pol := &StoragePoolSpec{
			BaseModel: &BaseModel{
				// Make sure uuid is unique
				Id: uuid.NewV5(uuid.NamespaceOID, name).String(),
			},
			Name:             poolId,
			TotalCapacity:    p.TotalCapacity >> UnitGiShiftBit,
			FreeCapacity:     (p.TotalCapacity - p.UsedCapacity) >> UnitGiShiftBit,
			StorageType:      c.Pool[poolId].StorageType,
			Extras:           c.Pool[poolId].Extras,
			AvailabilityZone: c.Pool[poolId].AvailabilityZone,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = DefaultAZ
		}
		pols = append(pols, pol)
	}
	return pols, nil
}
