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

// +build !driver_builtin driver_plugin

package drivers

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
	"path"
	"plugin"
	"strings"

	"github.com/opensds/opensds/contrib/drivers/factory"
	. "github.com/opensds/opensds/contrib/drivers/utils/constants"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/config"
)

func NewDriverFactory() factory.DriverFactory {
	return &DriverFactory{}
}

type DriverFactory struct {
	DriverLibPath string
}

func (d *DriverFactory) FindDriverSoPath(storageType StorageType, driverType DriverType,
	driverName string) (string, error) {
	var libPaths []string
	// find so in build out directory firstly
	if gopath, ok := os.LookupEnv("GOPATH"); ok {
		libPaths = append(libPaths, path.Join(gopath, "src/github.com/opensds/opensds/build/out/lib"))
	}
	libPaths = append(libPaths, "/opt/opensds-hotpot-linux-amd64/lib")
	libPaths = append(libPaths, "/usr/local/opensds/lib")

	for _, libPath := range libPaths {
		p := path.Join(libPath, string(storageType), driverName+".so")
		if ok, _ := utils.PathExists(p); ok {
			logs.Info("find %s driver so in path %s", driverName, p)
			return p, nil
		}
	}
	return "", fmt.Errorf("can't find %s.so, please check in path '%s'",
		driverName, strings.Join(libPaths, ","))
}

func (d *DriverFactory) GetDriver(storageType StorageType, driverType DriverType, bp config.BackendProperties) (factory.Driver, error) {
	var driverType2Sym = map[DriverType]string{
		DriverTypeProvision:   "NewDriver",
		DriverTypeReplication: "NewReplicationDriver",
		DriverTypeMetric:      "NewMetricDriver",
	}

	soPath, err := d.FindDriverSoPath(storageType, driverType, bp.DriverName)
	if err != nil {
		return nil, err
	}
	plug, err := plugin.Open(soPath)
	if err != nil {
		return nil, err
	}
	sym, err := plug.Lookup(driverType2Sym[driverType])
	if err != nil {
		return nil, err
	}
	fun, ok := sym.(func(properties config.BackendProperties) factory.Driver)
	if !ok {
		return nil, fmt.Errorf("unexpected type from module symbol")
	}
	return fun(bp), nil
}
