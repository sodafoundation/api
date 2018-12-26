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
	c "github.com/opensds/opensds/pkg/controller"
	"github.com/opensds/opensds/pkg/db"
	. "github.com/opensds/opensds/pkg/utils/config"
	"github.com/opensds/opensds/pkg/utils/daemon"
	"github.com/opensds/opensds/pkg/utils/logs"
)

func init() {
	// Load global configuration from specified config file.
	CONF.Load("/etc/opensds/opensds.conf")

	flag := &CONF.Flag
	flag.StringVar(&CONF.OsdsLet.ApiEndpoint, "api-endpoint", CONF.OsdsLet.ApiEndpoint, "Listen endpoint of controller service")
	flag.StringVar(&CONF.Database.Endpoint, "db-endpoint", CONF.Database.Endpoint, "Connection endpoint of database service")
	flag.StringVar(&CONF.Database.Driver, "db-driver", CONF.Database.Driver, "Driver name of database service")
	flag.DurationVar(&CONF.OsdsLet.LogFlushFrequency, "log-flush-frequency", CONF.OsdsLet.LogFlushFrequency, "Maximum number of seconds between log flushes")

	daemon.SetDaemonFlag(&CONF.OsdsLet.Daemon, CONF.OsdsLet.Daemon)
	daemon.CheckAndRunDaemon(CONF.OsdsLet.Daemon)
}

func main() {
	// Open OpenSDS orchestrator service log file.
	logs.InitLogs(CONF.OsdsLet.LogFlushFrequency)
	defer logs.FlushLogs()

	// Set up database session.
	db.Init(&CONF.Database)

	// Initialize Controller object.
	c.Brain = c.NewController()

	// Start OpenSDS northbound REST service.
	api.Run(CONF.OsdsLet.ApiEndpoint)
}
