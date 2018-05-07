// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

var volumeGroupCommand = &cobra.Command{
	Use:   "group",
	Short: "manage volume group in the cluster",
	Run:   volumeGroupAction,
}

var volumeGroupCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "create a volume group in the cluster",
	Run:   volumeGroupCreateAction,
}

var volumeGroupShowCommand = &cobra.Command{
	Use:   "show <id>",
	Short: "show a volume group in the cluster",
	Run:   volumeGroupShowAction,
}

var volumeGroupListCommand = &cobra.Command{
	Use:   "list",
	Short: "list all volume groups in the cluster",
	Run:   volumeGroupListAction,
}

var volumeGroupDeleteCommand = &cobra.Command{
	Use:   "delete <id>",
	Short: "delete a volume group in the cluster",
	Run:   volumeGroupDeleteAction,
}

var volumeGroupUpdateCommand = &cobra.Command{
	Use:   "update <id>",
	Short: "update a volume group in the cluster",
	Run:   volumeGroupUpdateAction,
}

var (
	vgName        string
	vgDesp        string
	vgAZ          string
	addVolumes    *[]string
	removeVolumes *[]string
)

func init() {
	volumeGroupCommand.AddCommand(volumeGroupCreateCommand)
	volumeGroupCreateCommand.Flags().StringVarP(&vgName, "name", "n", "", "the name of created volume group")
	volumeGroupCreateCommand.Flags().StringVarP(&vgDesp, "description", "d", "", "the description of created volume group")
	volumeGroupCreateCommand.Flags().StringVarP(&vgAZ, "availabilityZone", "a", "", "the availabilityZone of created volume group")
	volumeGroupCommand.AddCommand(volumeGroupShowCommand)
	volumeGroupCommand.AddCommand(volumeGroupListCommand)
	volumeGroupCommand.AddCommand(volumeGroupDeleteCommand)
	volumeGroupCommand.AddCommand(volumeGroupUpdateCommand)
	volumeGroupUpdateCommand.Flags().StringVarP(&vgName, "name", "n", "", "the name of updated volume group")
	volumeGroupUpdateCommand.Flags().StringVarP(&vgDesp, "description", "d", "", "the description of updated volume group")
	addVolumes = volumeGroupUpdateCommand.Flags().StringSliceP("addVolumes", "a", nil, "the addVolumes of updated volume group")
	removeVolumes = volumeGroupUpdateCommand.Flags().StringSliceP("removeVolumes", "r", nil, "the removeVolumes of updated volume group")
}

func volumeGroupAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func volumeGroupCreateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)
	vg := &model.VolumeGroupSpec{
		Name:             vgName,
		Description:      vgDesp,
		AvailabilityZone: vgAZ,
	}

	resp, err := client.CreateVolumeGroup(vg)
	PrintResponse(resp)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Status", "AvailabilityZone", "PoolId"}
	PrintDict(resp, keys, FormatterList{})
}

func volumeGroupShowAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	resp, err := client.GetVolumeGroup(args[0])
	PrintResponse(resp)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Status", "AvailabilityZone", "PoolId"}
	PrintDict(resp, keys, FormatterList{})
}

func volumeGroupListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)
	resp, err := client.ListVolumeGroups()
	PrintResponse(resp)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "Name", "Description", "Status", "AvailabilityZone", "PoolId"}
	PrintList(resp, keys, FormatterList{})
}

func volumeGroupDeleteAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	err := client.DeleteVolumeGroup(args[0], nil)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	fmt.Printf("Delete group(%s) success.\n", args[0])
}

func volumeGroupUpdateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	snp := &model.VolumeGroupSpec{
		Name:          vgName,
		Description:   vgDesp,
		AddVolumes:    *addVolumes,
		RemoveVolumes: *removeVolumes,
	}

	resp, err := client.UpdateVolumeGroup(args[0], snp)
	PrintResponse(resp)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Status", "AvailabilityZone", "PoolId"}
	PrintDict(resp, keys, FormatterList{})
}
