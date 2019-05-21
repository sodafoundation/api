// Copyright (c) 2019 The OpenSDS Authors.
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
	"log"
	"os"
	"strconv"

	"github.com/opensds/opensds/pkg/model"
	"github.com/spf13/cobra"
)

var fileShareCommand = &cobra.Command{
	Use:   "fileshare",
	Short: "manage fileshares in the cluster",
	Run:   fileShareAction,
}

var fileShareCreateCommand = &cobra.Command{
	Use:   "create <size>",
	Short: "create a fileshare in the cluster",
	Run:   fileShareCreateAction,
}

var fileShareDeleteCommand = &cobra.Command{
	Use:   "delete <id>",
	Short: "delete a fileshare in the cluster",
	Run:   fileShareDeleteAction,
}

var fileShareShowCommand = &cobra.Command{
	Use:   "show <id>",
	Short: "show a fileshare in the cluster",
	Run:   fileShareShowAction,
}

var fileShareListCommand = &cobra.Command{
	Use:   "list",
	Short: "list all fileshares in the cluster",
	Run:   fileShareListAction,
}

var fileShareUpdateCommand = &cobra.Command{
	Use:   "update <id>",
	Short: "update a fileshare in the cluster",
	Run:   fileShareUpdateAction,
}

var (
	shareAZ              string
	shareDescription     string
	shareExportLocations string
	shareID              string
	shareName            string
	sharePoolID          string
	shareProfileID       string
	shareProtocols       string
	shareSnapshotID      string
	shareStatus          string
	shareTenantID        string
	shareUserID          string
	fileShareID          string

	shareLimit   string
	shareOffset  string
	shareSortDir string
	shareSortKey string

	sharekeys = KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Size",
		"AvailabilityZone", "Status", "PoolId", "ProfileId", "Protocols",
		"TenantId", "UserId", "SnapshotId", "ExportLocations"}
)

func init() {
	fileShareCommand.AddCommand(fileShareCreateCommand)
	fileShareCommand.AddCommand(fileShareDeleteCommand)
	fileShareCommand.AddCommand(fileShareShowCommand)
	fileShareCommand.AddCommand(fileShareListCommand)
	fileShareCommand.AddCommand(fileShareUpdateCommand)
	fileShareCommand.AddCommand(fileShareSnapshotCommand)
	fileShareCommand.AddCommand(fileShareAclCommand)

	fileShareCreateCommand.Flags().StringVarP(&shareName, "name", "n", "", "the name of the fileshare")
	fileShareCreateCommand.Flags().StringVarP(&shareDescription, "description", "d", "", "the description of the fileshare")
	fileShareCreateCommand.Flags().StringVarP(&shareAZ, "availabilityZone", "a", "", "the locality that fileshare belongs to")
	fileShareCreateCommand.Flags().StringVarP(&shareSnapshotID, "snapshotId", "s", "", "the uuid of the snapshot which the fileshare is created")
	fileShareCreateCommand.Flags().StringVarP(&shareProfileID, "profileId", "p", "", "the uuid of the profile which the fileshare belongs to")
	fileShareCreateCommand.Flags().StringVarP(&shareExportLocations, "exportLocations", "e", "", "exportLocations of the fileshare")
	fileShareCreateCommand.Flags().StringVarP(&shareUserID, "userId", "u", "", "exportLocations of the fileshare")

	fileShareListCommand.Flags().StringVarP(&shareLimit, "limit", "", "50", "the number of ertries displayed per page")
	fileShareListCommand.Flags().StringVarP(&shareOffset, "offset", "", "0", "all requested data offsets")
	fileShareListCommand.Flags().StringVarP(&shareSortDir, "sortDir", "", "desc", "the sort direction of all requested data. supports asc or desc(default)")
	fileShareListCommand.Flags().StringVarP(&shareSortKey, "sortKey", "", "id",
		"the sort key of all requested data. supports id(default), name, status, availabilityZone, profileId, tenantId, userId, size, poolId, description, protocols, snapshotId, exportLocations")
	fileShareListCommand.Flags().StringVarP(&shareID, "id", "", "", "list share by id")
	fileShareListCommand.Flags().StringVarP(&shareName, "name", "", "", "list share by name")
	fileShareListCommand.Flags().StringVarP(&shareDescription, "description", "", "", "list share by description")
	fileShareListCommand.Flags().StringVarP(&shareTenantID, "tenantId", "", "", "list share by tenantId")
	fileShareListCommand.Flags().StringVarP(&shareUserID, "userId", "", "", "list share by userId")
	fileShareListCommand.Flags().StringVarP(&shareStatus, "status", "", "", "list share by status")
	fileShareListCommand.Flags().StringVarP(&sharePoolID, "poolId", "", "", "list share by poolId")
	fileShareListCommand.Flags().StringVarP(&shareAZ, "availabilityZone", "", "", "list share by availabilityZone")
	fileShareListCommand.Flags().StringVarP(&shareProfileID, "profileId", "", "", "list share by profileId")
	fileShareListCommand.Flags().StringVarP(&shareProtocols, "protocols", "", "", "list share by protocols")
	fileShareListCommand.Flags().StringVarP(&shareSnapshotID, "snapshotId", "", "", "list share by snapshotId")
	fileShareListCommand.Flags().StringVarP(&shareSize, "size", "", "", "list share by size")
	fileShareListCommand.Flags().StringVarP(&shareExportLocations, "exportLocations", "", "", "list share by exportLocations")

	fileShareUpdateCommand.Flags().StringVarP(&shareName, "name", "n", "", "the name of the fileshare")
	fileShareUpdateCommand.Flags().StringVarP(&shareDescription, "description", "d", "", "the description of the fileshare")
}

func fileShareAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func fileShareCreateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	size, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("error parsing size %s: %+v", args[0], err)
	}

	var exportLocations []string
	if "" != shareExportLocations {
		err = json.Unmarshal([]byte(shareExportLocations), &exportLocations)
		if err != nil {
			log.Fatalf("error parsing exportLocations %s: %+v", shareExportLocations, err)
		}
	}

	share := &model.FileShareSpec{
		Description:      shareDescription,
		Name:             shareName,
		Size:             int64(size),
		UserId:           shareUserID,
		AvailabilityZone: shareAZ,
		ExportLocations:  exportLocations,
		ProfileId:        shareProfileID,
		SnapshotId:       shareSnapshotID,
	}

	resp, err := client.CreateFileShare(share)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}

	PrintDict(resp, sharekeys, FormatterList{})
}

func fileShareDeleteAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	err := client.DeleteFileShare(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
}

func fileShareShowAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	resp, err := client.GetFileShare(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}

	PrintDict(resp, sharekeys, FormatterList{})
}

func fileShareListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)

	var opts = map[string]string{"limit": shareLimit, "offset": shareOffset, "sortDir": shareSortDir,
		"sortKey": shareSortKey, "Id": shareID, "Name": shareName, "Description": shareDescription,
		"TenantId": shareTenantID, "UserId": shareUserID, "AvailabilityZone": shareAZ, "Status": shareStatus,
		"PoolId": sharePoolID, "ProfileId": shareProfileID, "Protocols": shareProtocols, "snapshotId": shareSnapshotID,
		"size": shareSize, "ExportLocations": shareExportLocations}

	resp, err := client.ListFileShares(opts)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}

	PrintList(resp, sharekeys, FormatterList{})
}

func fileShareUpdateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	share := &model.FileShareSpec{
		Name:        shareName,
		Description: shareDescription,
	}

	resp, err := client.UpdateFileShare(args[0], share)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}

	PrintDict(resp, sharekeys, FormatterList{})
}
