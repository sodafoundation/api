// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
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

/*
This module implements a entry into the OpenSDS REST service.

*/

package main

import (
	"github.com/opensds/opensds/pkg/api"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/logs"
	"log"
	"strings"
)

func init() {
	conf := utils.CONF
	flag := utils.CONF.Flag
	flag.StringVar(&conf.OsdsLet.ApiEndpoint, "api-endpoint", conf.OsdsLet.ApiEndpoint, "Listen endpoint of controller service")
	flag.StringVar(&conf.Database.Endpoint, "db-endpoint", conf.Database.Endpoint, "Connection endpoint of database service")
	flag.StringVar(&conf.Database.Driver, "db-driver", conf.Database.Driver, "Driver name of database service")
	flag.StringVar(&conf.Database.Credential, "db-credential", conf.Database.Credential, "Connection credential of database service")
	conf.Load("/etc/opensds/opensds.conf")
}

func main() {
	// Open OpenSDS orchestrator service log file.
	conf := utils.CONF
	logs.InitLogs()
	defer logs.FlushLogs()

	// Set up database session.
	db.Init(&db.DBConfig{
		DriverName: conf.Database.Driver,
		Endpoints:  strings.Split(conf.Database.Endpoint, ","),
		Credential: conf.Database.Credential,
	})

	// Start OpenSDS northbound REST service.
	api.Run(conf.OsdsLet.ApiEndpoint)
}
