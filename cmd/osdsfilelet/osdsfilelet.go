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
// See the License for the specific language governing permissions and
// limitations under the License.

/*
This module implements a entry into the OpenSDS REST service.

*/

package main

import (
	"flag"

	"github.com/opensds/opensds/pkg/db"
	c "github.com/opensds/opensds/pkg/filesharecontroller"
	. "github.com/opensds/opensds/pkg/utils/config"
	"github.com/opensds/opensds/pkg/utils/daemon"
	"github.com/opensds/opensds/pkg/utils/logs"
)

func init() {
	// Load global configuration from specified config file.
	CONF.Load()

	// Parse some configuration fields from command line. and it will override the value which is got from config file.
	flag.StringVar(&CONF.OsdsfileLet.ApiEndpoint, "api-endpoint", CONF.OsdsfileLet.ApiEndpoint, "Listen endpoint of controller service")
	flag.BoolVar(&CONF.OsdsfileLet.Daemon, "daemon", CONF.OsdsfileLet.Daemon, "Run app as a daemon with -daemon=true")
	flag.DurationVar(&CONF.OsdsfileLet.LogFlushFrequency, "log-flush-frequency", CONF.OsdsfileLet.LogFlushFrequency, "Maximum number of seconds between log flushes")
	flag.Parse()

	daemon.CheckAndRunDaemon(CONF.OsdsfileLet.Daemon)
}

func main() {
	// Open OpenSDS orchestrator service log file.
	logs.InitLogs(CONF.OsdsfileLet.LogFlushFrequency)
	defer logs.FlushLogs()

	// Set up database session.
	db.Init(&CONF.Database)

	// Construct controller module grpc server struct and run controller server process.
	if err := c.FileShareNewController(CONF.OsdsfileLet.ApiEndpoint).Run(); err != nil {
		panic(err)
	}
}
