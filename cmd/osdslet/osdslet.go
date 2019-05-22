// Copyright (c) 2019 The OpenSDS Authors.
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

	c "github.com/opensds/opensds/pkg/controller"
	"github.com/opensds/opensds/pkg/db"
	. "github.com/opensds/opensds/pkg/utils/config"
	"github.com/opensds/opensds/pkg/utils/constants"
	"github.com/opensds/opensds/pkg/utils/daemon"
	"github.com/opensds/opensds/pkg/utils/logs"
)

func init() {
	// Load global configuration from specified config file.
	CONF.Load()

	// Parse some configuration fields from command line. and it will override the value which is got from config file.
	flag.StringVar(&CONF.OsdsLet.ApiEndpoint, "api-endpoint", CONF.OsdsLet.ApiEndpoint, "Listen endpoint of controller service")
	flag.BoolVar(&CONF.OsdsLet.Daemon, "daemon", CONF.OsdsLet.Daemon, "Run app as a daemon with -daemon=true")
	flag.DurationVar(&CONF.OsdsLet.LogFlushFrequency, "log-flush-frequency", CONF.OsdsLet.LogFlushFrequency, "Maximum number of seconds between log flushes")
	// prometheus related
	flag.StringVar(&CONF.OsdsLet.PrometheusPushMechanism, "prometheus-push-mechanism", CONF.OsdsLet.PrometheusPushMechanism, "Prometheus push mechanism")
	flag.StringVar(&CONF.OsdsLet.PrometheusConfHome, "prometheus-conf-home", CONF.OsdsLet.PrometheusConfHome, "Prometheus conf. path")
	flag.StringVar(&CONF.OsdsLet.PrometheusUrl, "prometheus-url", CONF.OsdsLet.PrometheusUrl, "Prometheus URL")
	flag.StringVar(&CONF.OsdsLet.PrometheusConfFile, "prometheus-conf-file", CONF.OsdsLet.PrometheusConfFile, "Prometheus conf. file")
	// alert manager related
	flag.StringVar(&CONF.OsdsLet.AlertmgrConfHome, "alertmgr-conf-home", CONF.OsdsLet.AlertmgrConfHome, "Alert manager conf. home")
	flag.StringVar(&CONF.OsdsLet.AlertMgrUrl, "alertmgr-url", CONF.OsdsLet.AlertMgrUrl, "Alert manager listen endpoint")
	flag.StringVar(&CONF.OsdsLet.AlertmgrConfFile, "alertmgr-conf-file", CONF.OsdsLet.AlertmgrConfFile, "Alert manager conf. file")
	// grafana related
	flag.StringVar(&CONF.OsdsLet.GrafanaConfHome, "grafana-conf-home", CONF.OsdsLet.GrafanaConfHome, "Grafana conf. home")
	flag.StringVar(&CONF.OsdsLet.GrafanaRestartCmd, "grafana-restart-cmd", CONF.OsdsLet.GrafanaRestartCmd, "Grafana restart command")
	flag.StringVar(&CONF.OsdsLet.GrafanaConfFile, "grafana-conf-file", CONF.OsdsLet.GrafanaConfFile, "Grafana conf file")
	flag.StringVar(&CONF.OsdsLet.GrafanaUrl, "grafana-url", CONF.OsdsLet.GrafanaUrl, "Grafana listen endpoint")
	// prometheus and alert manager configuration reload url
	flag.StringVar(&CONF.OsdsLet.ConfReloadUrl, "conf-reload-url", CONF.OsdsLet.ConfReloadUrl, "Prometheus and Alert manager conf. reload URL")

	flag.Parse()

	daemon.CheckAndRunDaemon(CONF.OsdsLet.Daemon)
}

func main() {
	// Open OpenSDS orchestrator service log file.
	logs.InitLogs(CONF.OsdsLet.LogFlushFrequency)
	defer logs.FlushLogs()

	// Set up database session.
	db.Init(&CONF.Database)

	// Construct controller module grpc server struct and run controller server process.
	if err := c.NewController(constants.OpensdsCtrBindEndpoint).Run(); err != nil {
		panic(err)
	}
}
