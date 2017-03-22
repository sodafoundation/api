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

	northbound "github.com/opensds/opensds/testing/cmd/fake_osdslet/fake_northbound"
	orchServer "github.com/opensds/opensds/testing/pkg/grpc/controller/orchestration/fake_server"
	dockServer "github.com/opensds/opensds/testing/pkg/grpc/dock/fake_server"
)

const (
	defaultOrchestrationPort = "localhost:50049"
	defaultDockPort          = "localhost:50050"
)

func main() {
	// Open OpenSDS log file
	f, err := os.OpenFile("/var/log/fake_opensds.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error opening file:", err)
		os.Exit(1)
	}
	defer f.Close()
	// assign it to the standard logger
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Start OpenSDS northbound REST service.
	go northbound.Run()

	// Construct orchestration module grpc server struct and do some initialization.
	os := orchServer.NewOrchServer(defaultOrchestrationPort)
	// Construct dock module grpc server struct and do some initialization.
	ds := dockServer.NewDockServer(defaultDockPort)

	// Start the listen mechanism of orchestration and dock module.
	go os.ListenAndServe()
	ds.ListenAndServe()
}
