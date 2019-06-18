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

	flag.StringVar(&CONF.OsdsLet.PrometheusPushMechanism, "prometheus-push-mechanism", CONF.OsdsLet.PrometheusPushMechanism, "Prometheus push mechanism")
	flag.StringVar(&CONF.OsdsLet.PushGatewayUrl, "prometheus-push-gateway-url", CONF.OsdsLet.PushGatewayUrl, "Prometheus push gateway URL")
	flag.StringVar(&CONF.OsdsLet.NodeExporterWatchFolder, "node-exporter-watch-folder", CONF.OsdsLet.NodeExporterWatchFolder, "Node exporter watch folder")
	flag.StringVar(&CONF.OsdsLet.KafkaEndpoint, "kafka-endpoint", CONF.OsdsLet.KafkaEndpoint, "Kafka endpoint")
	flag.StringVar(&CONF.OsdsLet.KafkaTopic, "kafka-topic", CONF.OsdsLet.KafkaTopic, "Kafka topic")
	flag.StringVar(&CONF.OsdsLet.GrafanaUrl, "grafana-url", CONF.OsdsLet.GrafanaUrl, "Grafana listen endpoint")
	flag.StringVar(&CONF.OsdsLet.AlertMgrUrl, "alertmgr-url", CONF.OsdsLet.AlertMgrUrl, "Alert manager listen endpoint")
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
