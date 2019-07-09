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

var (
	poolLimit            string
	poolOffset           string
	poolSortDir          string
	poolSortKey          string
	poolId               string
	poolName             string
	poolDescription      string
	poolStatus           string
	poolDockId           string
	poolAvailabilityZone string
	poolStorageType      string
)

func init() {
	poolListCommand.Flags().StringVarP(&poolLimit, "limit", "", "50", "the number of ertries displayed per page")
	poolListCommand.Flags().StringVarP(&poolOffset, "offset", "", "0", "all requested data offsets")
	poolListCommand.Flags().StringVarP(&poolSortDir, "sortDir", "", "desc", "the sort direction of all requested data. supports asc or desc(default)")
	poolListCommand.Flags().StringVarP(&poolSortKey, "sortKey", "", "id", "the sort key of all requested data. supports id(default), name, status, availabilityzone, dock id, description")
	poolListCommand.Flags().StringVarP(&poolId, "id", "", "", "list pools by id")
	poolListCommand.Flags().StringVarP(&poolName, "name", "", "", "list pools by name")
	poolListCommand.Flags().StringVarP(&poolDescription, "description", "", "", "list pools by description")
	poolListCommand.Flags().StringVarP(&poolStatus, "status", "", "", "list pools by status")
	poolListCommand.Flags().StringVarP(&poolStorageType, "storageType", "", "", "list pools by storage type")
	poolListCommand.Flags().StringVarP(&poolDockId, "dockId", "", "", "list pools by dock id")
	poolListCommand.Flags().StringVarP(&poolAvailabilityZone, "availabilityZone", "", "", "list pools by availability zone")

	poolCommand.AddCommand(poolShowCommand)
	poolCommand.AddCommand(poolListCommand)
}

func poolAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

var poolFormatters = FormatterList{"Extras": JsonFormatter}

func poolShowAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	pols, err := client.GetPool(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Status", "DockId",
		"AvailabilityZone", "TotalCapacity", "FreeCapacity", "StorageType", "MultiAttach", "Extras"}
	PrintDict(pols, keys, poolFormatters)
}

func poolListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)

	var opts = map[string]string{"limit": poolLimit, "offset": poolOffset, "sortDir": poolSortDir,
		"sortKey": poolSortKey, "Id": poolId,
		"Name": poolName, "Description": poolDescription, "AvailabilityZone": poolAvailabilityZone,
		"Status": poolStatus,
		"DockId": poolDockId, "StorageType": poolStorageType}

	pols, err := client.ListPools(opts)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "Name", "Description", "Status", "TotalCapacity", "FreeCapacity"}
	PrintList(pols, keys, poolFormatters)
}
