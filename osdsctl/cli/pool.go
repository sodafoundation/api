// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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
This module implements a entry into the OpenSDS service.

*/

package cli

import (
	"fmt"
	"os"

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
		fmt.Fprintln(os.Stderr, "The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}
	pols, err := client.GetPool(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Status", "DockId",
		"AvailabilityZone", "TotalCapacity", "FreeCapacity", "StorageType", "Extras"}
	PrintDict(pols, keys, FormatterList{})
}

func poolListAction(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Fprintln(os.Stderr, "The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}
	pols, err := client.ListPools()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	keys := KeyList{"Id", "Name", "Description", "Status", "AvailabilityZone", "TotalCapacity", "FreeCapacity"}
	PrintList(pols, keys, FormatterList{})
}
