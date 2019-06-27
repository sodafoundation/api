// Copyright 2018 The OpenSDS Authors.
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

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "manage API versions in the cluster",
	Run:   versionAction,
}

var versionShowCommand = &cobra.Command{
	Use:   "show <apiVersion>",
	Short: "show version details by specified API version in the cluster",
	Run:   versionShowAction,
}

var versionListCommand = &cobra.Command{
	Use:   "list",
	Short: "list information for all SDSController API versions in the cluster",
	Run:   versionListAction,
}

func init() {
	versionCommand.AddCommand(versionShowCommand)
	versionCommand.AddCommand(versionListCommand)
}

func versionAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func versionShowAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	resp, err := client.GetVersion(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Name", "Status", "UpdatedAt"}
	PrintDict(resp, keys, FormatterList{})
}

func versionListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)
	resp, err := client.ListVersions()
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Name", "Status", "UpdatedAt"}
	PrintList(resp, keys, FormatterList{})
}
