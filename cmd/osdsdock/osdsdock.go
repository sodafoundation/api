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
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/opensds/opensds/cmd/osdsdock/app"
	"github.com/opensds/opensds/pkg/db"
	dockServer "github.com/opensds/opensds/pkg/grpc/dock/server"
)

var (
	// The endpoint of dock server, format: "host_ip:port"
	apiEdp string

	// The driver name of database, for example etcd, mysql etc
	dbDriver string
	// Connecting endpoint of database client, only be used in etcd
	dbEdp string
	// Connecting credentials of database, only be used in mysql
	// format: "username:password@tcp(ip:port)/dbname"
	dbCredential string

	// Path of file for logging, default is "/var/log/opensds/osdsdock.log"
	osdsdockLogFile string
)

func init() {
	flag.StringVar(&apiEdp, "api-endpoint", "localhost:50050", "Listen endpoint of controller service")
	flag.StringVar(&dbEdp, "db-endpoint", "localhost:2379,localhost:2380", "Connection endpoint of database service")
	flag.StringVar(&dbDriver, "db-driver", "etcd", "Driver name of database service")
	flag.StringVar(&dbCredential, "db-credential", "username:password@tcp(ip:port)/dbname", "Connection credential of database service")
	flag.StringVar(&osdsdockLogFile, "osdsdocklog-file", "/var/log/opensds/osdsdock.log", "Location of osdsdock log file")

	flag.Parse()
}

func main() {
	// Open OpenSDS dock service log file.
	f, err := os.OpenFile(osdsdockLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer f.Close()

	// assign it to the standard logger.
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Set up database session.
	db.Init(&db.DBConfig{
		DriverName: dbDriver,
		Endpoints:  strings.Split(dbEdp, ","),
		Credential: dbCredential,
	})

	// Automatically discover dock and pool resources from backends.
	if err = app.ResourceDiscovery(); err != nil {
		panic(err)
	}

	// Construct dock module grpc server struct and do some initialization.
	ds := dockServer.NewDockServer(apiEdp)
	// Start the listen mechanism of dock module.
	ds.ListenAndServe()
}
