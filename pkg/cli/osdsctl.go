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
This module implements a entry into the OpenSDS CLI service.

*/

package cli

import (
	"os"

	c "github.com/opensds/opensds/client"
	"github.com/spf13/cobra"
)

var (
	client      *c.Client
	rootCommand = &cobra.Command{
		Use:   "osdsctl",
		Short: "Administer the opensds storage cluster",
		Long:  `Admin utility for the opensds unified storage cluster.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
			os.Exit(1)
		},
	}
)

func init() {
	client = c.NewClient(&c.Config{})
	//	rootCommand.AddCommand(versionCommand)
	rootCommand.AddCommand(volumeCommand)
	rootCommand.AddCommand(dockCommand)
	rootCommand.AddCommand(poolCommand)
	rootCommand.AddCommand(profileCommand)
}

func Run() error {
	return rootCommand.Execute()
}
