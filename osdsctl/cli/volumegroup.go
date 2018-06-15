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
	vgLimit    string
	vgOffset   string
	vgSortDir  string
	vgSortKey  string
	vgId       string
	vgTenantId string
	vgUserId   string

	vgName        string
	vgDesp        string
	vgAZ          string
	addVolumes    *[]string
	removeVolumes *[]string
	vgprofiles    *[]string
	vgStatus      string
	vgPoolId      string
)

func init() {
	volumeGroupListCommand.Flags().StringVarP(&vgLimit, "limit", "", "50", "the number of ertries displayed per page")
	volumeGroupListCommand.Flags().StringVarP(&vgOffset, "offset", "", "0", "all requested data offsets")
	volumeGroupListCommand.Flags().StringVarP(&vgSortDir, "sortDir", "", "desc", "the sort direction of all requested data. supports asc or desc(default)")
	volumeGroupListCommand.Flags().StringVarP(&vgSortKey, "sortKey", "", "id",
		"the sort key of all requested data. supports id(default), name, status, availability zone, tenantid, pool id")
	volumeGroupListCommand.Flags().StringVarP(&vgId, "id", "", "", "list volume group by id")
	volumeGroupListCommand.Flags().StringVarP(&vgTenantId, "tenantId", "", "", "list volume group by tenantId")
	volumeGroupListCommand.Flags().StringVarP(&vgUserId, "userId", "", "", "list volume group by storage userId")
	volumeGroupListCommand.Flags().StringVarP(&vgStatus, "status", "", "", "list volume group by status")
	volumeGroupListCommand.Flags().StringVarP(&vgName, "name", "", "", "list volume group by Name")
	volumeGroupListCommand.Flags().StringVarP(&vgDesp, "description", "", "", "list volume group by description")
	volumeGroupListCommand.Flags().StringVarP(&vgAZ, "availabilityZone", "", "", "list volume group by availability zone")
	volumeGroupListCommand.Flags().StringVarP(&vgPoolId, "poolId", "", "", "list volume group by pool id")

	volumeGroupCommand.AddCommand(volumeGroupCreateCommand)
	volumeGroupCreateCommand.Flags().StringVarP(&vgName, "name", "n", "", "the name of created volume group")
	volumeGroupCreateCommand.Flags().StringVarP(&vgDesp, "description", "d", "", "the description of created volume group")
	volumeGroupCreateCommand.Flags().StringVarP(&vgAZ, "availabilityZone", "a", "", "the availabilityZone of created volume group")
	vgprofiles = volumeGroupCreateCommand.Flags().StringSliceP("profiles", "", nil, "the profiles of created volume group")
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
		Profiles:         *vgprofiles,
	}

	resp, err := client.CreateVolumeGroup(vg)
	PrintResponse(resp)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Status", "AvailabilityZone", "PoolId", "Profiles"}
	PrintDict(resp, keys, FormatterList{})
}

func volumeGroupShowAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	resp, err := client.GetVolumeGroup(args[0])
	PrintResponse(resp)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Status", "AvailabilityZone", "PoolId", "Profiles"}
	PrintDict(resp, keys, FormatterList{})
}

func volumeGroupListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)

	var opts = map[string]string{"limit": vgLimit, "offset": vgOffset, "sortDir": vgSortDir,
		"sortKey": vgSortKey, "Id": vgId,
		"Name": vgName, "Description": vgDesp, "UserId": vgUserId, "AvailabilityZone": vgAZ,
		"Status": vgStatus, "PoolId": vgPoolId}

	resp, err := client.ListVolumeGroups(opts)
	PrintResponse(resp)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "Name", "Description", "Status", "AvailabilityZone", "PoolId", "Profiles"}
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
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "Status", "AvailabilityZone", "PoolId", "Profiles"}
	PrintDict(resp, keys, FormatterList{})
}
