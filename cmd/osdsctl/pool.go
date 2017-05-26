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

	pools "github.com/opensds/opensds/pkg/apiserver"

	"github.com/spf13/cobra"
)

var poolCommand = &cobra.Command{
	Use:   "pool",
	Short: "manage OpenSDS pool resources",
	Run:   poolAction,
}

var poolShowCommand = &cobra.Command{
	Use:   "show <pool id>",
	Short: "show information of specified pool",
	Run:   poolShowAction,
}

var poolListCommand = &cobra.Command{
	Use:   "list",
	Short: "get all pool resources",
	Run:   poolListAction,
}

func init() {
	poolCommand.AddCommand(poolShowCommand)
	poolCommand.AddCommand(poolListCommand)
}

func poolAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func poolShowAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	poolRequest := &pools.PoolRequest{
		Id: args[0],
	}

	result, err := pools.GetPool(poolRequest)
	if err != nil {
		fmt.Println("Get pool resource failed: ", err)
	} else {
		rbody, _ := json.MarshalIndent(result, "", "  ")
		fmt.Printf("%s\n", string(rbody))
	}
}

func poolListAction(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	poolRequest := &pools.PoolRequest{}

	result, err := pools.ListPools(poolRequest)
	if err != nil {
		fmt.Println("List pool resources failed: ", err)
	} else {
		rbody, _ := json.MarshalIndent(result, "", "  ")
		fmt.Printf("%s\n", string(rbody))
	}
}
