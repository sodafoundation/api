// Copyright 2019 The OpenSDS Authors.
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
	rootCommand.AddCommand(fileShareCommand)
	flags := rootCommand.PersistentFlags()
	flags.BoolVar(&Debug, "debug", false, "shows debugging output.")
}

type DummyWriter struct{}

// do nothing
func (writer DummyWriter) Write(data []byte) (n int, err error) {
	return len(data), nil
}

type DebugWriter struct{}

// do nothing
func (writer DebugWriter) Write(data []byte) (n int, err error) {
	Debugf("%s", string(data))
	return len(data), nil
}

// Run method indicates how to start a cli tool through cobra.
func Run() error {
	if !utils.Contained("--debug", os.Args) {
		log.SetOutput(DummyWriter{})
	} else {
		log.SetOutput(DebugWriter{})
	}

	ep, ok := os.LookupEnv(c.OpensdsEndpoint)
	if !ok {
		ep = constants.DefaultOpensdsEndpoint
		Warnf("OPENSDS_ENDPOINT is not specified, use default(%s)\n", ep)
	}

	cfg := &c.Config{Endpoint: ep}

	authStrategy, ok := os.LookupEnv(c.OpensdsAuthStrategy)
	if !ok {
		authStrategy = c.Noauth
		Warnf("Not found Env OPENSDS_AUTH_STRATEGY, use default(noauth)\n")
	}

	var authOptions c.AuthOptions
	var err error

	switch authStrategy {
	case c.Keystone:
		authOptions, err = c.LoadKeystoneAuthOptionsFromEnv()
		if err != nil {
			return err
		}
	case c.Noauth:
		authOptions = c.LoadNoAuthOptionsFromEnv()
	default:
		authOptions = c.NewNoauthOptions(constants.DefaultTenantId)
	}

	cfg.AuthOptions = authOptions

	client, err = c.NewClient(cfg)
	if client == nil || err != nil {
		return err
	}

	return rootCommand.Execute()
}
