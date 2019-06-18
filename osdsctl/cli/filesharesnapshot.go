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
	"os"

	"github.com/opensds/opensds/pkg/model"
	"github.com/spf13/cobra"
)

var fileShareSnapshotCommand = &cobra.Command{
	Use:   "snapshot",
	Short: "manage fileshare snapshots in the cluster",
	Run:   fileShareSnapshotAction,
}

var fileShareSnapshotCreateCommand = &cobra.Command{
	Use:   "create <fileshare id>",
	Short: "create a snapshot of specified fileshare in the cluster",
	Run:   fileShareSnapshotCreateAction,
}

var fileShareSnapshotShowCommand = &cobra.Command{
	Use:   "show <fileshare snapshot id>",
	Short: "show a fileshare snapshot in the cluster",
	Run:   fileShareSnapshotShowAction,
}

var fileShareSnapshotListCommand = &cobra.Command{
	Use:   "list",
	Short: "list all fileshare snapshots in the cluster",
	Run:   fileShareSnapshotListAction,
}

var fileShareSnapshotDeleteCommand = &cobra.Command{
	Use:   "delete <fileshare snapshot id>",
	Short: "delete a fileshare snapshot of specified fileshare in the cluster",
	Run:   fileShareSnapshotDeleteAction,
}

var fileShareSnapshotUpdateCommand = &cobra.Command{
	Use:   "update <fileshare snapshot id>",
	Short: "update a fileshare snapshot in the cluster",
	Run:   fileShareSnapshotUpdateAction,
}

var (
	shareSnapshotName string
	shareSnapshotDesp string
)

var (
	shareSnapLimit     string
	shareSnapOffset    string
	shareSnapSortDir   string
	shareSnapSortKey   string
	shareSnapID        string
	shareSnapUserID    string
	shareSnapName      string
	shareSnapDesp      string
	shareSnapStatus    string
	shareSize          string
	shareSnapSize      string
	shareSnapTenantID  string
	shareSnapProfileID string
	fileshareID        string
	shareSnapCreatedAt string
	shareSnapUpdatedAt string

	shareSnapFormatters = FormatterList{"Metadata": JsonFormatter}
)

func init() {
	fileShareSnapshotCommand.AddCommand(fileShareSnapshotCreateCommand)
	fileShareSnapshotCommand.AddCommand(fileShareSnapshotDeleteCommand)
	fileShareSnapshotCommand.AddCommand(fileShareSnapshotShowCommand)
	fileShareSnapshotCommand.AddCommand(fileShareSnapshotListCommand)
	fileShareSnapshotCommand.AddCommand(fileShareSnapshotUpdateCommand)

	fileShareSnapshotCreateCommand.Flags().StringVarP(&shareSnapName, "name", "n", "", "the name of the fileshare snapshot")
	fileShareSnapshotCreateCommand.Flags().StringVarP(&shareSnapDesp, "description", "d", "", "the description of the fileshare snapshot")
	fileShareSnapshotCreateCommand.Flags().StringVarP(&shareSnapProfileID, "profileId", "P", "", "the description of the fileshare snapshot")

	fileShareSnapshotListCommand.Flags().StringVarP(&shareSnapLimit, "limit", "", "50", "the number of ertries displayed per page")
	fileShareSnapshotListCommand.Flags().StringVarP(&shareSnapOffset, "offset", "", "0", "all requested data offsets")
	fileShareSnapshotListCommand.Flags().StringVarP(&shareSnapSortDir, "sortDir", "", "desc", "the sort direction of all requested data. supports asc or desc(default)")
	fileShareSnapshotListCommand.Flags().StringVarP(&shareSnapSortKey, "sortKey", "", "id",
		"the sort key of all requested data. supports id(default), createdAt, updatedAt, name, description, snapshotSize, status, userId, tenantId, fileshareId, profileId")
	fileShareSnapshotListCommand.Flags().StringVarP(&shareSnapID, "id", "", "", "list fileshare snapshot by id")
	fileShareSnapshotListCommand.Flags().StringVarP(&shareSnapName, "name", "", "", "list fileshare snapshot by name")
	fileShareSnapshotListCommand.Flags().StringVarP(&shareSnapDesp, "description", "", "", "list fileshare snapshot by description")
	fileShareSnapshotListCommand.Flags().StringVarP(&shareSnapSize, "snapshotSize", "", "", "list fileshare snapshot by snapshotSize")
	fileShareSnapshotListCommand.Flags().StringVarP(&shareSnapStatus, "status", "", "", "list fileshare snapshot by status")
	fileShareSnapshotListCommand.Flags().StringVarP(&shareSnapUserID, "userId", "", "", "list fileshare snapshot by userId")
	fileShareSnapshotListCommand.Flags().StringVarP(&shareSnapTenantID, "tenantId", "", "", "list fileshare snapshot by tenantId")
	fileShareSnapshotListCommand.Flags().StringVarP(&fileshareID, "fileshareId", "", "", "list fileshare snapshot by fileshareId")
	fileShareSnapshotListCommand.Flags().StringVarP(&shareSnapProfileID, "profileId", "", "", "list fileshare snapshot by profileId")
	fileShareSnapshotListCommand.Flags().StringVarP(&shareSnapCreatedAt, "createdAt", "", "", "list fileshare snapshot by createdAt")
	fileShareSnapshotListCommand.Flags().StringVarP(&shareSnapUpdatedAt, "updatedAt", "", "", "list fileshare snapshot by updatedAt")

	fileShareSnapshotUpdateCommand.Flags().StringVarP(&shareSnapshotName, "name", "n", "", "the name of the fileshare snapshot")
	fileShareSnapshotUpdateCommand.Flags().StringVarP(&shareSnapshotDesp, "description", "d", "", "the description of the fileshare snapshot")
}

func fileShareSnapshotAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func fileShareSnapshotCreateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	snp := &model.FileShareSnapshotSpec{
		Name:        shareSnapName,
		Description: shareSnapDesp,
		FileShareId: args[0],
		ProfileId:   shareSnapProfileID,
	}

	resp, err := client.CreateFileShareSnapshot(snp)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}

	keys := KeyList{"Id", "CreatedAt", "Name", "Description", "ShareSize",
		"Status", "FileShareId", "SnapshotSize", "TenantId", "UserId",
		"ProfileId", "Metadata"}
	PrintDict(resp, keys, shareSnapFormatters)
}

func fileShareSnapshotShowAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	resp, err := client.GetFileShareSnapshot(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}

	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description",
		"ShareSize", "Status", "FileShareId", "SnapshotSize", "TenantId",
		"UserId", "ProfileId", "Metadata"}
	PrintDict(resp, keys, shareSnapFormatters)
}

func fileShareSnapshotListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)

	var opts = map[string]string{"limit": shareSnapLimit, "offset": shareSnapOffset, "sortDir": shareSnapSortDir,
		"sortKey": shareSnapSortKey, "Id": shareSnapID,
		"Name": shareSnapName, "Description": shareSnapDesp, "UserId": shareSnapUserID,
		"Status": shareSnapStatus, "CreatedAt": shareSnapCreatedAt, "UpdatedAt": shareSnapUpdatedAt,
		"Size": shareSnapSize, "TenantId": shareSnapTenantID, "FileShareId": fileshareID}

	resp, err := client.ListFileShareSnapshots(opts)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}

	keys := KeyList{"Id", "Name", "Description", "Status",
		"FileShareId", "ProfileId", "SnapshotSize"}
	PrintList(resp, keys, shareSnapFormatters)
}

func fileShareSnapshotDeleteAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)

	err := client.DeleteFileShareSnapshot(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
}

func fileShareSnapshotUpdateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	snp := &model.FileShareSnapshotSpec{
		Name:        shareSnapshotName,
		Description: shareSnapshotDesp,
	}

	resp, err := client.UpdateFileShareSnapshot(args[0], snp)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}

	keys := KeyList{"Id", "UpdatedAt", "Name", "Description", "ShareSize",
		"Status", "FileShareId", "SnapshotSize", "TenantId", "UserId",
		"ProfileId", "Metadata"}
	PrintDict(resp, keys, shareSnapFormatters)
}
