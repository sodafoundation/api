// Copyright 2017 The OpenSDS Authors.
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

package config

type Default struct{}

type OsdsLet struct {
	ApiEndpoint string `conf:"api_endpoint,localhost:50040"`
	Graceful    bool   `conf:"graceful,true"`
	SocketOrder string `conf:"socket_order"`
	Daemon      bool   `conf:"daemon,false"`
}

type OsdsDock struct {
	ApiEndpoint     string   `conf:"api_endpoint,localhost:50050"`
	EnabledBackends []string `conf:"enabled_backends,ceph"`
	Daemon          bool     `conf:"daemon,false"`
	Backends
}

type Database struct {
	Credential string `conf:"credential,username:password@tcp(ip:port)/dbname"`
	Driver     string `conf:"driver,etcd"`
	Endpoint   string `conf:"endpoint,localhost:2379,localhost:2380"`
}

type BackendProperties struct {
	Name        string `conf:"name"`
	Description string `conf:"description"`
	DriverName  string `conf:"driver_name"`
	ConfigPath  string `conf:"config_path"`
}

type Backends struct {
	Ceph         BackendProperties `conf:"ceph"`
	Cinder       BackendProperties `conf:"cinder"`
	Sample       BackendProperties `conf:"sample"`
	LVM          BackendProperties `conf:"lvm"`
	HuaweiDorado BackendProperties `conf:"huawei_dorado"`
}

type Config struct {
	Default  `conf:"default"`
	OsdsLet  `conf:"osdslet"`
	OsdsDock `conf:"osdsdock"`
	Database `conf:"database"`
	Flag     FlagSet
}
