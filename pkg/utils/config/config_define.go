// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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
	ApiEndpoint  string `conf:"api_endpoint,localhost:50040"`
	Graceful     bool   `conf:"graceful,true"`
	SocketOrder  string `conf:"socket_order"`
	AuthStrategy string `conf:"auth_strategy,noauth"`
	Daemon       bool   `conf:"daemon,false"`
	PolicyPath   string `conf:"policy_path,/etc/opensds/policy.json"`
}

type OsdsDock struct {
	ApiEndpoint                string   `conf:"api_endpoint,localhost:50050"`
	DockType                   string   `conf:"dock_type,provisioner"`
	EnabledBackends            []string `conf:"enabled_backends,lvm"`
	Daemon                     bool     `conf:"daemon,false"`
	BindIp                     string   `conf:"bind_ip"` // Just used for attacher dock
	HostBasedReplicationDriver string   `conf:"host_based_replication_driver,drbd"`
	Backends
}

type Database struct {
	Credential string `conf:"credential,username:password@tcp(ip:port)/dbname"`
	Driver     string `conf:"driver,etcd"`
	Endpoint   string `conf:"endpoint,localhost:2379,localhost:2380"`
}

type BackendProperties struct {
	Name               string `conf:"name"`
	Description        string `conf:"description"`
	DriverName         string `conf:"driver_name"`
	ConfigPath         string `conf:"config_path"`
	SupportReplication bool   `conf:"support_replication,false"`
}

type Backends struct {
	Ceph         BackendProperties `conf:"ceph"`
	Cinder       BackendProperties `conf:"cinder"`
	Sample       BackendProperties `conf:"sample"`
	LVM          BackendProperties `conf:"lvm"`
	HuaweiDorado BackendProperties `conf:"huawei_dorado"`
}

type KeystoneAuthToken struct {
	MemcachedServers  string `conf:"memcached_servers"`
	SigningDir        string `conf:"signing_dir"`
	Cafile            string `conf:"cafile"`
	AuthUri           string `conf:"auth_uri"`
	ProjectDomainName string `conf:"project_domain_name"`
	ProjectName       string `conf:"project_name"`
	UserDomainName    string `conf:"user_domain_name"`
	Password          string `conf:"password"`
	Username          string `conf:"username"`
	AuthUrl           string `conf:"auth_url"`
	AuthType          string `conf:"auth_type"`
}

type Config struct {
	Default           `conf:"default"`
	OsdsLet           `conf:"osdslet"`
	OsdsDock          `conf:"osdsdock"`
	Database          `conf:"database"`
	KeystoneAuthToken `conf:"keystone_authtoken"`
	Flag              FlagSet
}
