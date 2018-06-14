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

/*
This module implements a entry into the OpenSDS REST service.

*/

package main

import (
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/dock"
	"github.com/opensds/opensds/pkg/dock/server"
	. "github.com/opensds/opensds/pkg/utils/config"
	"github.com/opensds/opensds/pkg/utils/daemon"
	"github.com/opensds/opensds/pkg/utils/logs"
)

func init() {
	def := GetDefaultConfig()
	flag := CONF.Flag
	flag.StringVar(&CONF.OsdsDock.ApiEndpoint, "api-endpoint", def.OsdsDock.ApiEndpoint, "Listen endpoint of dock service")
	flag.StringVar(&CONF.OsdsDock.DockType, "dock-type", def.OsdsDock.DockType, "Type of dock service")
	flag.StringVar(&CONF.Database.Endpoint, "db-endpoint", def.Database.Endpoint, "Connection endpoint of database service")
	flag.StringVar(&CONF.Database.Driver, "db-driver", def.Database.Driver, "Driver name of database service")
	// flag.StringVar(&CONF.Database.Credential, "db-credential", def.Database.Credential, "Connection credential of database service")
	daemon.SetDaemonFlag(&CONF.OsdsDock.Daemon, def.OsdsDock.Daemon)
	CONF.Load("/etc/opensds/opensds.conf")
	daemon.CheckAndRunDaemon(CONF.OsdsDock.Daemon)
}

func main() {
	// Open OpenSDS dock service log file.
	logs.InitLogs()
	defer logs.FlushLogs()

	// Set up database session.
	db.Init(&CONF.Database)

	// Automatically discover dock and pool resources from backends.
	dock.Brain = dock.NewDockHub(CONF.OsdsDock.DockType)
	if err := dock.Brain.TriggerDiscovery(); err != nil {
		panic(err)
	}

	// Construct dock module grpc server struct and start the listen mechanism
	// of dock module.
	ds := server.NewDockServer(CONF.OsdsDock.ApiEndpoint)
	if err := ds.Run(); err != nil {
		panic(err)
	}
}
