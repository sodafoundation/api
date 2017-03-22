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

	"github.com/opensds/opensds/cmd/osdslet/northbound"
	"github.com/opensds/opensds/cmd/utils"
	orchServer "github.com/opensds/opensds/pkg/grpc/controller/orchestration/server"
)

const (
	NORTHBOUND_PORT    = ":50048"
	ORCHESTRATION_PORT = ":50049"
)

func main() {
	// Open OpenSDS log file
	f, err := os.OpenFile("/var/log/opensds/osdslet.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer f.Close()
	// assign it to the standard logger
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Get OpenSDS host IP.
	host, err := utils.GetHostIP()
	if err != nil {
		panic(err)
	}

	// Start OpenSDS northbound REST service.
	go northbound.Run(host + NORTHBOUND_PORT)

	// Construct orchestration module grpc server struct and do some initialization.
	os := orchServer.NewOrchServer(host + ORCHESTRATION_PORT)

	// Start the listen mechanism of controller orchestration module.
	os.ListenAndServe()
}
