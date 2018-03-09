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
	"os"

	c "github.com/opensds/opensds/client"
	"github.com/opensds/opensds/pkg/utils/constants"
	"github.com/spf13/cobra"
)

const (
	// Opensds Auth EVNs
	OpensdsEndpoint     = "OPENSDS_ENDPOINT"
	OpensdsAuthStrategy = "OPENSDS_AUTH_STRATEGY"
	OpensdsTenantId     = "OPENSDS_TENANT_ID"

	// Keystone Auth ENVs
	OsAuthUrl      = "OS_AUTH_URL"
	OsUsername     = "OS_USERNAME"
	OsPassword     = "OS_PASSWORD"
	OsTenantName   = "OS_TENANT_NAME"
	OsProjectName  = "OS_PROJECT_NAME"
	OsUserDomainId = "OS_USER_DOMAIN_ID"
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
	rootCommand.AddCommand(versionCommand)
	rootCommand.AddCommand(volumeCommand)
	rootCommand.AddCommand(dockCommand)
	rootCommand.AddCommand(poolCommand)
	rootCommand.AddCommand(profileCommand)
}

func LoadKeystoneAuthOptionsFromEnv() *c.KeystoneAuthOptions {
	opt := c.NewKeystoneAuthOptions()
	opt.IdentityEndpoint = os.Getenv(OsAuthUrl)
	opt.Username = os.Getenv(OsUsername)
	opt.Password = os.Getenv(OsPassword)
	opt.TenantName = os.Getenv(OsTenantName)
	projectName := os.Getenv(OsProjectName)
	opt.DomainID = os.Getenv(OsUserDomainId)
	if opt.TenantName == "" {
		opt.TenantName = projectName
	}
	return opt
}

func LoadNoAuthOptionsFromEnv() *c.NoAuthOptions {
	tenantId, ok := os.LookupEnv(OpensdsTenantId)
	if !ok {
		return c.NewNoauthOptions(constants.DefaultTenantId)
	}
	return c.NewNoauthOptions(tenantId)
}

// Run method indicates how to start a cli tool through cobra.
func Run() error {
	ep, ok := os.LookupEnv(OpensdsEndpoint)
	if !ok {
		return fmt.Errorf("ERROR: You must provide the endpoint by setting " +
			"the environment variable OPENSDS_ENDPOINT")
	}
	cfg := &c.Config{Endpoint: ep}

	authStrategy, ok := os.LookupEnv(OpensdsAuthStrategy)
	if !ok {
		authStrategy = c.Noauth
		fmt.Println("WARNING: Not found Env OPENSDS_AUTH_STRATEGY, use default(noauth)\n")
	}

	switch authStrategy {
	case c.Keystone:
		cfg.AuthOptions = LoadKeystoneAuthOptionsFromEnv()
	case c.Noauth:
		cfg.AuthOptions = LoadNoAuthOptionsFromEnv()
	default:
		cfg.AuthOptions = c.NewNoauthOptions(constants.DefaultTenantId)
	}

	client = c.NewClient(cfg)

	return rootCommand.Execute()
}
