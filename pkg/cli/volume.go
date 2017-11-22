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
	"log"
	"os"
	"strconv"

	"github.com/opensds/opensds/pkg/model"
	"github.com/spf13/cobra"
)

var volumeCommand = &cobra.Command{
	Use:   "volume",
	Short: "manage volumes in the cluster",
	Run:   volumeAction,
}

var volumeCreateCommand = &cobra.Command{
	Use:   "create <size>",
	Short: "create a volume in the cluster",
	Run:   volumeCreateAction,
}

var volumeShowCommand = &cobra.Command{
	Use:   "show <id>",
	Short: "show a volume in the cluster",
	Run:   volumeShowAction,
}

var volumeListCommand = &cobra.Command{
	Use:   "list",
	Short: "list all volumes in the cluster",
	Run:   volumeListAction,
}

var volumeDeleteCommand = &cobra.Command{
	Use:   "delete <id>",
	Short: "delete a volume in the cluster",
	Run:   volumeDeleteAction,
}

var (
	profileId string
	volName   string
	volDesp   string
	volAz     string
)

func init() {
	volumeCommand.PersistentFlags().StringVarP(&profileId, "profile", "p", "", "the name of profile configured by admin")

	volumeCommand.AddCommand(volumeCreateCommand)
	volumeCreateCommand.Flags().StringVarP(&volName, "name", "n", "null", "the name of created volume")
	volumeCreateCommand.Flags().StringVarP(&volDesp, "description", "d", "", "the description of created volume")
	volumeCreateCommand.Flags().StringVarP(&volAz, "az", "a", "", "the availabilty zone of created volume")
	volumeCommand.AddCommand(volumeShowCommand)
	volumeCommand.AddCommand(volumeListCommand)
	volumeCommand.AddCommand(volumeDeleteCommand)

	volumeCommand.AddCommand(volumeSnapshotCommand)
}

func volumeAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func volumeCreateAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}
	size, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("error parsing size %s: %+v", args[0], err)
	}

	vol := &model.VolumeSpec{
		Name:             volName,
		Description:      volDesp,
		AvailabilityZone: volAz,
		Size:             int64(size),
		ProfileId:        profileId,
	}

	resp, err := client.CreateVolume(vol)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Size",
		"AvailabilityZone", "Status", "PoolId", "ProfileId", "Metadata"}
	PrintDict(resp, keys, FormatterList{})
}

func volumeShowAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	resp, err := client.GetVolume(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Size",
		"AvailabilityZone", "Status", "PoolId", "ProfileId", "Metadata"}
	PrintDict(resp, keys, FormatterList{})
}

func volumeListAction(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	resp, err := client.ListVolumes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	keys := KeyList{"Id", "Name", "Description", "Size",
		"AvailabilityZone", "Status", "PoolId", "ProfileId"}
	PrintList(resp, keys, FormatterList{})
}

func volumeDeleteAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}
	vol := &model.VolumeSpec{
		ProfileId: profileId,
	}
	err := client.DeleteVolume(args[0], vol)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Delete volume(%s) sucess.\n", args[0])
}

