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
This module implements a entry into the OpenSDS CLI service.

*/

package cli

import (
	"fmt"
	"log"
	"os"

	c "github.com/opensds/opensds/client"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/constants"
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
	Debug bool
)

func init() {
	rootCommand.AddCommand(versionCommand)
	rootCommand.AddCommand(volumeCommand)
	rootCommand.AddCommand(dockCommand)
	rootCommand.AddCommand(poolCommand)
	rootCommand.AddCommand(profileCommand)
	rootCommand.AddCommand(replicationCommand)
	flags := rootCommand.PersistentFlags()
	flags.BoolVar(&Debug, "debug", false, "shows debugging output.")
}

type Writer struct{}

// do nothing
func (writer Writer) Write(data []byte) (n int, err error) {
	return len(data), nil
}

// Run method indicates how to start a cli tool through cobra.
func Run() error {
	if !utils.Contained("--debug", os.Args) {
		log.SetOutput(Writer{})
	}

	ep, ok := os.LookupEnv(c.OpensdsEndpoint)
	if !ok {
		return fmt.Errorf("ERROR: You must provide the endpoint by setting " +
			"the environment variable OPENSDS_ENDPOINT")
	}
	cfg := &c.Config{Endpoint: ep}

	authStrategy, ok := os.LookupEnv(c.OpensdsAuthStrategy)
	if !ok {
		authStrategy = c.Noauth
		fmt.Println("WARNING: Not found Env OPENSDS_AUTH_STRATEGY, use default(noauth)")
	}

	switch authStrategy {
	case c.Keystone:
		cfg.AuthOptions = c.LoadKeystoneAuthOptionsFromEnv()
	case c.Noauth:
		cfg.AuthOptions = c.LoadNoAuthOptionsFromEnv()
	default:
		cfg.AuthOptions = c.NewNoauthOptions(constants.DefaultTenantId)
	}

	client = c.NewClient(cfg)

	return rootCommand.Execute()
}
