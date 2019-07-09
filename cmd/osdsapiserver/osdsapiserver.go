// Copyright 2018 The OpenSDS Authors.
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

	"github.com/opensds/opensds/pkg/api"
	"github.com/opensds/opensds/pkg/db"
	. "github.com/opensds/opensds/pkg/utils/config"
	"github.com/opensds/opensds/pkg/utils/daemon"
	"github.com/opensds/opensds/pkg/utils/logs"
)

func init() {
	// Load global configuration from specified config file.
	CONF.Load()

	// Parse some configuration fields from command line. and it will override the value which is got from config file.
	flag.StringVar(&CONF.OsdsApiServer.ApiEndpoint, "api-endpoint", CONF.OsdsApiServer.ApiEndpoint, "Listen endpoint of api-server service")
	flag.DurationVar(&CONF.OsdsApiServer.LogFlushFrequency, "log-flush-frequency", CONF.OsdsApiServer.LogFlushFrequency, "Maximum number of seconds between log flushes")
	flag.BoolVar(&CONF.OsdsApiServer.Daemon, "daemon", CONF.OsdsApiServer.Daemon, "Run app as a daemon with -daemon=true")
	// prometheus related
	flag.StringVar(&CONF.OsdsApiServer.PrometheusConfHome, "prometheus-conf-home", CONF.OsdsApiServer.PrometheusConfHome, "Prometheus conf. path")
	flag.StringVar(&CONF.OsdsApiServer.PrometheusUrl, "prometheus-url", CONF.OsdsApiServer.PrometheusUrl, "Prometheus URL")
	flag.StringVar(&CONF.OsdsApiServer.PrometheusConfFile, "prometheus-conf-file", CONF.OsdsApiServer.PrometheusConfFile, "Prometheus conf. file")
	// alert manager related
	flag.StringVar(&CONF.OsdsApiServer.AlertmgrConfHome, "alertmgr-conf-home", CONF.OsdsApiServer.AlertmgrConfHome, "Alert manager conf. home")
	flag.StringVar(&CONF.OsdsApiServer.AlertMgrUrl, "alertmgr-url", CONF.OsdsApiServer.AlertMgrUrl, "Alert manager listen endpoint")
	flag.StringVar(&CONF.OsdsApiServer.AlertmgrConfFile, "alertmgr-conf-file", CONF.OsdsApiServer.AlertmgrConfFile, "Alert manager conf. file")
	// grafana related
	flag.StringVar(&CONF.OsdsApiServer.GrafanaConfHome, "grafana-conf-home", CONF.OsdsApiServer.GrafanaConfHome, "Grafana conf. home")
	flag.StringVar(&CONF.OsdsApiServer.GrafanaRestartCmd, "grafana-restart-cmd", CONF.OsdsApiServer.GrafanaRestartCmd, "Grafana restart command")
	flag.StringVar(&CONF.OsdsApiServer.GrafanaConfFile, "grafana-conf-file", CONF.OsdsApiServer.GrafanaConfFile, "Grafana conf file")
	flag.StringVar(&CONF.OsdsApiServer.GrafanaUrl, "grafana-url", CONF.OsdsApiServer.GrafanaUrl, "Grafana listen endpoint")
	// prometheus and alert manager configuration reload url
	flag.StringVar(&CONF.OsdsApiServer.ConfReloadUrl, "conf-reload-url", CONF.OsdsApiServer.ConfReloadUrl, "Prometheus and Alert manager conf. reload URL")
	flag.Parse()

	daemon.CheckAndRunDaemon(CONF.OsdsApiServer.Daemon)
}

func main() {
	// Open OpenSDS orchestrator service log file.
	logs.InitLogs(CONF.OsdsApiServer.LogFlushFrequency)
	defer logs.FlushLogs()

	// Set up database session.
	db.Init(&CONF.Database)

	// Start OpenSDS northbound REST service.
	api.Run(CONF.OsdsApiServer)
}
