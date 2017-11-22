// Copyright 2017 The OpenSDS Authors.
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

var dockCommand = &cobra.Command{
	Use:   "dock",
	Short: "manage OpenSDS dock resources",
	Run:   dockAction,
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
	dockCommand.AddCommand(dockShowCommand)
	dockCommand.AddCommand(dockListCommand)
}

func dockAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func dockShowAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}
	resp, err := client.GetDock(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Endpoint", "DriverName", "Parameters"}
	PrintDict(resp, keys, FormatterList{})
}

func dockListAction(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	resp, err := client.ListDocks()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	keys := KeyList{"Id", "Name", "Description", "Endpoint", "DriverName", "Parameters"}
	PrintList(resp, keys, FormatterList{})
}

