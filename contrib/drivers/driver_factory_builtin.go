// Copyright 2019 The OpenSDS Authors.
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
// See the License for the specific

// +build driver_builtin

package drivers

import (
	"fmt"

	"github.com/opensds/opensds/contrib/drivers/factory"
	oceanstorfs "github.com/opensds/opensds/contrib/drivers/fileshare/huawei/oceanstor"
	"github.com/opensds/opensds/contrib/drivers/fileshare/manila"
	"github.com/opensds/opensds/contrib/drivers/fileshare/nfs"
	"github.com/opensds/opensds/contrib/drivers/utils/constants"
	. "github.com/opensds/opensds/contrib/drivers/utils/constants"
	"github.com/opensds/opensds/contrib/drivers/volume/ceph"
	"github.com/opensds/opensds/contrib/drivers/volume/drbd"
	"github.com/opensds/opensds/contrib/drivers/volume/fujitsu/eternus"
	"github.com/opensds/opensds/contrib/drivers/volume/hpe/nimble"
	"github.com/opensds/opensds/contrib/drivers/volume/huawei/fusionstorage"
	"github.com/opensds/opensds/contrib/drivers/volume/huawei/oceanstor"
	"github.com/opensds/opensds/contrib/drivers/volume/lvm"
	"github.com/opensds/opensds/contrib/drivers/volume/openstack/cinder"
	"github.com/opensds/opensds/contrib/drivers/volume/scutech/cms"
	"github.com/opensds/opensds/pkg/utils/config"
)

func NewDriverFactory() factory.DriverFactory {
	return &DriverFactory{}
}

type DriverFactory struct {
}

type NewDriver func(properties config.BackendProperties) factory.Driver

var VolumeDriverMap = map[string]NewDriver{
	DriverNameCeph:                ceph.NewDriver,
	DriverNameFujitsuEternus:      eternus.NewDriver,
	DriverNameHPENimble:           nimble.NewDriver,
	DriverNameHuaweiOceanStor:     oceanstor.NewDriver,
	DriverNameHuaweiFusionStorage: fusionstorage.NewDriver,
	DriverNameLVM:                 lvm.NewDriver,
	DriverNameCinder:              cinder.NewDriver,
}

var VolumeReplicationDriverMap = map[string]NewDriver{
	DriverNameDRBD:            drbd.NewReplicationDriver,
	DriverNameHuaweiOceanStor: oceanstor.NewReplicationDriver,
	DriverNameScutechCMS:      scms.NewReplicationDriver,
}

var VolumeMetricsDriverMap = map[string]NewDriver{
	DriverNameCeph:            ceph.NewMetricDriver,
	DriverNameHuaweiOceanStor: oceanstor.NewMetricDriver,
	DriverNameLVM:             lvm.NewMetricDriver,
}

var FileShareDriverMap = map[string]NewDriver{
	DriverNameNFS:             nfs.NewDriver,
	DriverNameHuaweiOceanStor: oceanstorfs.NewDriver,
	DriverNameManila:          manila.NewDriver,
}

var drivers = map[string]map[string]NewDriver{}

func init() {
	drivers[string(StorageTypeBlock)+DriverTypeProvision.String()] = VolumeDriverMap
	drivers[string(StorageTypeBlock)+DriverTypeReplication.String()] = VolumeReplicationDriverMap
	drivers[string(StorageTypeBlock)+DriverTypeMetric.String()] = VolumeMetricsDriverMap
	drivers[string(StorageTypeFile)+DriverTypeProvision.String()] = FileShareDriverMap
}

func (d *DriverFactory) GetDriver(st constants.StorageType, dt constants.DriverType,
	bp config.BackendProperties) (factory.Driver, error) {
	driverMap, ok := drivers[string(st)+dt.String()]
	if !ok {
		return nil, fmt.Errorf("cann't find specified driver map, storageType:%s, driverType:%s", st, dt)
	}
	fun, ok := driverMap[bp.DriverName]
	if !ok {
		return nil, fmt.Errorf("cann't find specified driver: %s", bp.DriverName)
	}

	return fun(bp), nil
}
