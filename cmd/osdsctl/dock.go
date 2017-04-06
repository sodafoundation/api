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
	"fmt"
	"os"

	"github.com/opensds/opensds/pkg/controller/api/v1/docks"

	"github.com/spf13/cobra"
)

var dockCommand = &cobra.Command{
	Use:   "dock",
	Short: "manage OpenSDS dock resources",
	Run:   dockAction,
}

var dockRegisterCommand = &cobra.Command{
	Use:   "register <dock ip>",
	Short: "register a new dock with ip to controller",
	Run:   dockRegisterAction,
}

var dockDeregisterCommand = &cobra.Command{
	Use:   "deregister <dock ip>",
	Short: "deregister an existing dock with ip from controller",
	Run:   dockDeregisterAction,
}

var dockShowCommand = &cobra.Command{
	Use:   "show <dock id>",
	Short: "show information of specified dock",
	Run:   dockShowAction,
}

var dockListCommand = &cobra.Command{
	Use:   "list",
	Short: "get all dock resources",
	Run:   dockListAction,
}

func init() {
	dockCommand.AddCommand(dockRegisterCommand)
	dockCommand.AddCommand(dockDeregisterCommand)
	dockCommand.AddCommand(dockShowCommand)
	dockCommand.AddCommand(dockListCommand)
}

func dockAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func dockRegisterAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	dockRequest := docks.DockRequest{}

	dock, err := docks.RegisterDock(dockRequest, args[0])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", dock)
	}
}

func dockDeregisterAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	dockRequest := docks.DockRequest{}

	res, err := docks.DeregisterDock(dockRequest, args[0])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s\n", res)
	}
}

func dockShowAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	dockRequest := docks.DockRequest{}

	dock, err := docks.GetDock(dockRequest, args[0])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", dock)
	}
}

func dockListAction(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	dockRequest := docks.DockRequest{}

	docks, err := docks.ListDocks(dockRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", docks)
	}
}
