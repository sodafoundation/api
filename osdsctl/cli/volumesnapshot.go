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

	"github.com/opensds/opensds/pkg/model"
	"github.com/spf13/cobra"
)

var volumeSnapshotCommand = &cobra.Command{
	Use:   "snapshot",
	Short: "manage volume snapshots in the cluster",
	Run:   volumeSnapshotAction,
}

var volumeSnapshotCreateCommand = &cobra.Command{
	Use:   "create <volume id>",
	Short: "create a snapshot of specified volume in the cluster",
	Run:   volumeSnapshotCreateAction,
}

var volumeSnapshotShowCommand = &cobra.Command{
	Use:   "show <snapshot id>",
	Short: "show a volume snapshot in the cluster",
	Run:   volumeSnapshotShowAction,
}

var volumeSnapshotListCommand = &cobra.Command{
	Use:   "list",
	Short: "list all volume snapshots in the cluster",
	Run:   volumeSnapshotListAction,
}

var volumeSnapshotDeleteCommand = &cobra.Command{
	Use:   "delete <snapshot id>",
	Short: "delete a volume snapshot of specified volume in the cluster",
	Run:   volumeSnapshotDeleteAction,
}

var volumeSnapshotUpdateCommand = &cobra.Command{
	Use:   "update <snapshot id>",
	Short: "update a volume snapshot in the cluster",
	Run:   volumeSnapshotUpdateAction,
}

var (
	volSnapshotName string
	volSnapshotDesp string
)

var (
	volSnapLimit       string
	volSnapOffset      string
	volSnapSortDir     string
	volSnapSortKey     string
	volSnapId          string
	volSnapUserId      string
	volSnapName        string
	volSnapDescription string
	volSnapStatus      string
	volSnapVolumeId    string
)

func init() {
	volumeSnapshotListCommand.Flags().StringVarP(&volSnapLimit, "limit", "", "50", "the number of ertries displayed per page")
	volumeSnapshotListCommand.Flags().StringVarP(&volSnapOffset, "offset", "", "0", "all requested data offsets")
	volumeSnapshotListCommand.Flags().StringVarP(&volSnapSortDir, "sortDir", "", "desc", "the sort direction of all requested data. supports asc or desc(default)")
	volumeSnapshotListCommand.Flags().StringVarP(&volSnapSortKey, "sortKey", "", "id",
		"the sort key of all requested data. supports id(default), volumeid, status, userid, tenantid, size")
	volumeSnapshotListCommand.Flags().StringVarP(&volSnapId, "id", "", "", "list volume snapshot by id")
	volumeSnapshotListCommand.Flags().StringVarP(&volSnapUserId, "userId", "", "", "list volume snapshot by storage userId")
	volumeSnapshotListCommand.Flags().StringVarP(&volSnapVolumeId, "volumeId", "", "", "list volume snapshot by volume id")
	volumeSnapshotListCommand.Flags().StringVarP(&volSnapStatus, "status", "", "", "list volume snapshot by status")
	volumeSnapshotListCommand.Flags().StringVarP(&volSnapName, "name", "", "", "list volume snapshot by Name")
	volumeSnapshotListCommand.Flags().StringVarP(&volSnapDescription, "description", "", "", "list volume snapshot by description")

	volumeSnapshotCommand.AddCommand(volumeSnapshotCreateCommand)
	volumeSnapshotCreateCommand.Flags().StringVarP(&volSnapshotName, "name", "n", "", "the name of created volume snapshot")
	volumeSnapshotCreateCommand.Flags().StringVarP(&volSnapshotDesp, "description", "d", "", "the description of created volume snapshot")
	volumeSnapshotCommand.AddCommand(volumeSnapshotShowCommand)
	volumeSnapshotCommand.AddCommand(volumeSnapshotListCommand)
	volumeSnapshotCommand.AddCommand(volumeSnapshotDeleteCommand)
	volumeSnapshotCommand.AddCommand(volumeSnapshotUpdateCommand)
	volumeSnapshotUpdateCommand.Flags().StringVarP(&volSnapshotName, "name", "n", "", "the name of updated volume snapshot")
	volumeSnapshotUpdateCommand.Flags().StringVarP(&volSnapshotDesp, "description", "d", "", "the description of updated volume snapshot")
}

func volumeSnapshotAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

var volSnapshotFormatters = FormatterList{"Metadata": JsonFormatter}

func volumeSnapshotCreateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	snp := &model.VolumeSnapshotSpec{
		Name:        volSnapshotName,
		Description: volSnapshotDesp,
		VolumeId:    args[0],
		ProfileId:   profileId,
	}

	resp, err := client.CreateVolumeSnapshot(snp)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "Name", "Description", "Size", "Status",
		"ProfileId", "VolumeId", "Metadata"}
	PrintDict(resp, keys, volSnapshotFormatters)
}

func volumeSnapshotShowAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	resp, err := client.GetVolumeSnapshot(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Size", "Status",
		"ProfileId", "VolumeId", "Metadata"}
	PrintDict(resp, keys, volSnapshotFormatters)
}

func volumeSnapshotListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)

	var opts = map[string]string{"limit": volSnapLimit, "offset": volSnapOffset, "sortDir": volSnapSortDir,
		"sortKey": volSnapSortKey, "Id": volSnapId,
		"Name": volSnapName, "Description": volSnapDescription, "UserId": volSnapUserId,
		"Status": volSnapStatus, "VolumeId": volSnapVolumeId}

	resp, err := client.ListVolumeSnapshots(opts)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "Name", "Description", "Size", "Status", "ProfileId", "VolumeId"}
	PrintList(resp, keys, volSnapshotFormatters)
}

func volumeSnapshotDeleteAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	snapID := args[0]
	err := client.DeleteVolumeSnapshot(snapID, nil)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
}

func volumeSnapshotUpdateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	snp := &model.VolumeSnapshotSpec{
		Name:        volSnapshotName,
		Description: volSnapshotDesp,
	}

	resp, err := client.UpdateVolumeSnapshot(args[0], snp)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "UpdatedAt", "Name", "Description", "Size", "Status",
		"ProfileId", "VolumeId", "Metadata"}
	PrintDict(resp, keys, volSnapshotFormatters)
}
