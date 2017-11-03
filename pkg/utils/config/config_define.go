// Copyright (c) 2017 OpenSDS Authors.
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

package config

import (
	gflag "flag"
)

type Default struct{}

type OsdsLet struct {
	ApiEndpoint string `conf:"api_endpoint,localhost:50040"`
	Graceful    bool   `conf:"graceful,true"`
	SocketOrder string `conf:"socket_order"`
}

type OsdsDock struct {
	ApiEndpoint    string   `conf:"api_endpoint,localhost:50050"`
	EnableBackends []string `conf:"enabled_backends,ceph"`
	CinderConfig   string   `conf:"cinder_config,/etc/opensds/driver/cinder.yaml"`
	CephConfig     string   `conf:"ceph_config,/etc/opensds/driver/ceph.yaml"`
	LVMConfig      string   `conf:"lvm_config,/etc/opensds/driver/lvm.yaml"`
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
}

type Ceph BackendProperties
type Cinder BackendProperties
type Sample BackendProperties
type LVM BackendProperties

type Config struct {
	Default  `conf:"default"`
	OsdsLet  `conf:"osdslet"`
	OsdsDock `conf:"osdsdock"`
	Database `conf:"database"`
	Ceph     `conf:"ceph"`
	Cinder   `conf:"cinder"`
	Sample   `conf:"sample"`
	LVM      `conf:"lvm"`
	Flag     FlagSet
}

//Create a Config and init default value.
func GetDefaultConfig() *Config {
	var conf *Config = new(Config)
	initConf("", conf)
	return conf
}

func (c *Config) Load(confFile string) {
	gflag.StringVar(&confFile, "config-file", confFile, "The configuration file of OpenSDS")
	c.Flag.Parse()
	initConf(confFile, CONF)
	c.Flag.AssignValue()
}

var CONF *Config = GetDefaultConfig()
