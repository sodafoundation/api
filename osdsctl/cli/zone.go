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
This module implements a entry into the OpenSDS service.
*/

package cli

import (
	"encoding/json"
	"os"

	"github.com/opensds/opensds/pkg/model"
	"github.com/spf13/cobra"
)

var zoneCommand = &cobra.Command{
	Use:   "zone",
	Short: "manage OpenSDS Availability Zone resources",
	Run:   zoneAction,
}

var zoneCreateCommand = &cobra.Command{
	Use:   "create <availability zone info>",
	Short: "create a new availability zone resource",
	Run:   zoneCreateAction,
}

var zoneShowCommand = &cobra.Command{
	Use:   "show <availability zone id>",
	Short: "show information of specified availability zone",
	Run:   zoneShowAction,
}

var zoneListCommand = &cobra.Command{
	Use:   "list",
	Short: "get all availability zone resources",
	Run:   zoneListAction,
}

var zoneDeleteCommand = &cobra.Command{
	Use:   "delete <availability zone id>",
	Short: "delete a specified availability zone resource",
	Run:   zoneDeleteAction,
}

var zoneUpdateCommand = &cobra.Command{
	Use:   "update <availability zone id> <availability zone info>",
	Short: "update a specified zone resource",
	Run:   zoneUpdateAction,
}

var (
	zoneLimit       string
	zoneOffset      string
	zoneSortDir     string
	zoneSortKey     string
	zoneId          string
	zoneName        string
	zoneDescription string
)

func init() {
	zoneListCommand.Flags().StringVarP(&zoneLimit, "limit", "", "50", "the number of ertries displayed per page")
	zoneListCommand.Flags().StringVarP(&zoneOffset, "offset", "", "0", "all requested data offsets")
	zoneListCommand.Flags().StringVarP(&zoneSortDir, "sortDir", "", "desc", "the sort direction of all requested data. supports asc or desc(default)")
	zoneListCommand.Flags().StringVarP(&zoneSortKey, "sortKey", "", "id", "the sort key of all requested data. supports id(default), name, description")
	zoneListCommand.Flags().StringVarP(&zoneId, "id", "", "", "list availability zone by id")
	zoneListCommand.Flags().StringVarP(&zoneName, "name", "", "", "list availability zone by name")
	zoneListCommand.Flags().StringVarP(&zoneDescription, "description", "", "", "list availability zone by description")

	zoneCommand.AddCommand(zoneCreateCommand)
	zoneCommand.AddCommand(zoneShowCommand)
	zoneCommand.AddCommand(zoneListCommand)
	zoneCommand.AddCommand(zoneDeleteCommand)
	zoneCommand.AddCommand(zoneUpdateCommand)
}

func zoneAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

var zoneFormatters = FormatterList{}

func zoneCreateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	az := &model.AvailabilityZoneSpec{}
	if err := json.Unmarshal([]byte(args[0]), az); err != nil {
		Errorln(err)
		cmd.Usage()
		os.Exit(1)
	}

	resp, err := client.CreateAvailabilityZone(az)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description"}
	PrintDict(resp, keys, zoneFormatters)
}

func zoneShowAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	resp, err := client.GetAvailabilityZone(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description"}
	PrintDict(resp, keys, zoneFormatters)
}

func zoneListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)
	var opts = map[string]string{"limit": zoneLimit, "offset": zoneOffset, "sortDir": zoneSortDir,
		"sortKey": zoneSortKey, "Id": zoneId,
		"Name": zoneName, "Description": zoneDescription}

	resp, err := client.ListAvailabilityZones(opts)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description"}
	PrintList(resp, keys, FormatterList{})
}

func zoneDeleteAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	err := client.DeleteAvailabilityZone(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
}

func zoneUpdateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 2)
	az := &model.AvailabilityZoneSpec{}

	if err := json.Unmarshal([]byte(args[1]), az); err != nil {
		Errorln(err)
		cmd.Usage()
		os.Exit(1)
	}

	resp, err := client.UpdateAvailabilityZone(args[0], az)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "UpdatedAt", "Name", "Description"}
	PrintDict(resp, keys, FormatterList{})
}