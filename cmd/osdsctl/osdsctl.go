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
This module implements a entry into the OpenSDS CLI service.

*/

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "osdsctl",
	Short: "Administer the opensds storage cluster",
	Long:  `Admin utility for the opensds unified storage cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(1)
	},
}

func init() {
	rootCommand.AddCommand(versionCommand)
	rootCommand.AddCommand(shareCommand)
	rootCommand.AddCommand(volumeCommand)
	rootCommand.AddCommand(dockCommand)
	rootCommand.AddCommand(profileCommand)
	rootCommand.AddCommand(poolCommand)
}

func main() {
	// Open OpenSDS CLI service log file
	f, err := os.OpenFile("/var/log/opensds/osdsctl.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer f.Close()
	// assign it to the standard logger
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := rootCommand.Execute(); err != nil {
		die("%+v", err)
	}
}

func die(why string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, why+"\n", args...)
	os.Exit(1)
}
