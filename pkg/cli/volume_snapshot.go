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
	"fmt"
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
	Use:   "delete <volume id> <snapshot id>",
	Short: "delete a volume snapshot of specified volume in the cluster",
	Run:   volumeSnapshotDeleteAction,
}

var (
	volSnapshotName string
	volSnapshotDesp string
)

func init() {
	volumeSnapshotCommand.AddCommand(volumeSnapshotCreateCommand)
	volumeSnapshotCreateCommand.Flags().StringVarP(&volSnapshotName, "name", "n", "", "the name of created volume snapshot")
	volumeSnapshotCreateCommand.Flags().StringVarP(&volSnapshotDesp, "description", "d", "", "description of created volume snapshot")
	volumeSnapshotCommand.AddCommand(volumeSnapshotShowCommand)
	volumeSnapshotCommand.AddCommand(volumeSnapshotListCommand)
	volumeSnapshotCommand.AddCommand(volumeSnapshotDeleteCommand)
}

func volumeSnapshotAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func volumeSnapshotCreateAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	snp := &model.VolumeSnapshotSpec{
		Name:        volSnapshotName,
		Description: volSnapshotDesp,
		VolumeId:    args[0],
	}

	resp, err := client.CreateVolumeSnapshot(snp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Size", "Status", "VolumeId"}
	PrintDict(resp, keys, FormatterList{})
}

func volumeSnapshotShowAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	resp, err := client.GetVolumeSnapshot(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Size", "Status", "VolumeId"}
	PrintDict(resp, keys, FormatterList{})
}

func volumeSnapshotListAction(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	resp, err := client.ListVolumeSnapshots()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	keys := KeyList{"Id", "Name", "Description", "Size", "Status", "VolumeId"}
	PrintList(resp, keys, FormatterList{})
}

func volumeSnapshotDeleteAction(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}
	snp := &model.VolumeSnapshotSpec{
		VolumeId: args[0],
	}
	err := client.DeleteVolumeSnapshot(args[1], snp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Delete snapshot(%s) sucess.\n", args[1])
}

