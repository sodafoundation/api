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
This module implements a entry into the OpenSDS service.

*/

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/opensds/opensds/pkg/api"

	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "manage OpenSDS versions",
	Run:   versionAction,
}

var versionListCommand = &cobra.Command{
	Use:   "list",
	Short: "get all available versions",
	Run:   versionListAction,
}

var versionShowCommand = &cobra.Command{
	Use:   "show <version id>",
	Short: "show information of specified version",
	Run:   versionShowAction,
}

var fakeVersion api.VersionInfo
var fakeVersions api.AvailableVersions

func init() {
	versionCommand.AddCommand(versionListCommand)
	versionCommand.AddCommand(versionShowCommand)
}

func versionAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func versionListAction(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	versions, err := api.ListVersions()
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(versions, fakeVersions) {
			fmt.Println("List versions failed!")
		} else {
			rbody, _ := json.Marshal(versions)
			fmt.Printf("%s\n", string(rbody))
		}
	}
}

func versionShowAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	switch args[0] {
	case "v1":
		version, err := api.GetVersionv1()
		if err != nil {
			fmt.Println(err)
		} else {
			if reflect.DeepEqual(version, fakeVersion) {
				fmt.Println("Get version failed!")
			} else {
				rbody, _ := json.Marshal(version)
				fmt.Printf("%s\n", string(rbody))
			}
		}
	default:
		fmt.Println("Version not supported!")
	}

}
