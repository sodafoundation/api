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
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/opensds/opensds/pkg/api"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/db"
)

func init() {
	conf := utils.CONF
	flag := utils.CONF.Flag
	flag.StringVar(&conf.OsdsLet.ApiEndpoint, "api-endpoint", conf.OsdsLet.ApiEndpoint, "Listen endpoint of controller service")
	flag.StringVar(&conf.Database.Endpoint, "db-endpoint", conf.Database.Endpoint, "Connection endpoint of database service")
	flag.StringVar(&conf.Database.Driver, "db-driver", conf.Database.Driver, "Driver name of database service")
	flag.StringVar(&conf.Database.Credential, "db-credential", conf.Database.Credential, "Connection credential of database service")
	flag.StringVar(&conf.OsdsLet.LogFile, "osdsletlog-file", conf.OsdsLet.LogFile, "Location of osdslet log file")
	conf.Load("/etc/opensds/opensds.conf")
}

func main() {
	// Open OpenSDS orchestrator service log file.
	conf := utils.CONF
	f, err := os.OpenFile(conf.OsdsLet.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Assign it to the standard logger.
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Set up database session.
	db.Init(&db.DBConfig{
		DriverName: conf.Database.Driver,
		Endpoints:  strings.Split(conf.Database.Endpoint, ","),
		Credential: conf.Database.Credential,
	})

	// Start OpenSDS northbound REST service.
	api.Run(conf.OsdsLet.ApiEndpoint)
}

