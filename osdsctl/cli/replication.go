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
	"os"

	"encoding/json"
	"strings"

	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/spf13/cobra"
)

var replicationCommand = &cobra.Command{
	Use:   "replication",
	Short: "manage replications in the cluster",
	Run:   replicationAction,
}

var replicationCreateCommand = &cobra.Command{
	Use:   "create <primary volume id> <secondary volume id>",
	Short: "create a replication of specified volumes in the cluster",
	Run:   replicationCreateAction,
}

var replicationShowCommand = &cobra.Command{
	Use:   "show <replication id>",
	Short: "show a replication in the cluster",
	Run:   replicationShowAction,
}

var replicationListCommand = &cobra.Command{
	Use:   "list",
	Short: "list all replications in the cluster",
	Run:   replicationListAction,
}

var replicationDeleteCommand = &cobra.Command{
	Use:   "delete <replication id>",
	Short: "delete a replication in the cluster",
	Run:   replicationDeleteAction,
}

var replicationUpdateCommand = &cobra.Command{
	Use:   "update <replication id>",
	Short: "update a replication in the cluster",
	Run:   replicationUpdateAction,
}
var replicationEnableCommand = &cobra.Command{
	Use:   "enable <replication id>",
	Short: "enable a replication in the cluster",
	Run:   replicationEnableAction,
}

var replicationDisableCommand = &cobra.Command{
	Use:   "disable <replication id>",
	Short: "disable a replication in the cluster",
	Run:   replicationDisableAction,
}

var replicationFailoverCommand = &cobra.Command{
	Use:   "failover <replication id>",
	Short: "failover a replication in the cluster",
	Run:   replicationFailoverAction,
}

var (
	replicationName                string
	replicationDesp                string
	primaryReplicationDriverData   string
	secondaryReplicationDriverData string
	replicationMode                string
	replicationPeriod              int64
	allowAttachedVolume            bool
	secondaryBackendId             string
)

var (
	repLimit             string
	repOffset            string
	repSortDir           string
	repSortKey           string
	repId                string
	repName              string
	repDesp              string
	repPrimaryVolumeId   string
	repSecondaryVolumeId string
)

func init() {
	replicationListCommand.Flags().StringVarP(&repLimit, "limit", "", "50", "the number of ertries displayed per page")
	replicationListCommand.Flags().StringVarP(&repOffset, "offset", "", "0", "all requested data offsets")
	replicationListCommand.Flags().StringVarP(&repSortDir, "sortDir", "", "desc", "the sort direction of all requested data. supports asc or desc(default)")
	replicationListCommand.Flags().StringVarP(&repSortKey, "sortKey", "", "id",
		"the sort key of all requested data. supports id(default), name, primaryVolumeId, secondaryVolumeId,  description, create time, updatetime")
	replicationListCommand.Flags().StringVarP(&repId, "id", "", "", "list replication by id")
	replicationListCommand.Flags().StringVarP(&repName, "name", "", "", "list replication by name")
	replicationListCommand.Flags().StringVarP(&repDesp, "description", "", "", "list replication by description")
	replicationListCommand.Flags().StringVarP(&repPrimaryVolumeId, "primaryVolumeId", "", "", "list replication by PrimaryVolumeId")
	replicationListCommand.Flags().StringVarP(&repSecondaryVolumeId, "secondaryVolumeId", "", "", "list replication by storage userId")

	replicationCommand.AddCommand(replicationCreateCommand)
	flags := replicationCreateCommand.Flags()
	flags.StringVarP(&replicationName, "name", "n", "", "the name of created replication")
	flags.StringVarP(&replicationDesp, "description", "d", "", "the description of created replication")
	flags.StringVarP(&primaryReplicationDriverData, "primary_driver_data", "p", "", "the primary replication driver data of created replication")
	flags.StringVarP(&secondaryReplicationDriverData, "secondary_driver_data", "s", "", "the secondary replication driver data of created replication")
	flags.StringVarP(&replicationMode, "replication_mode", "m", model.ReplicationModeSync, "the replication mode of created replication, value can be sync/async")
	flags.Int64VarP(&replicationPeriod, "replication_period", "t", 0, "the replication period(minute) of created replication, the value must greater than 0, only in sync replication mode should set this value (default 60)")
	replicationUpdateCommand.Flags().StringVarP(&replicationName, "name", "n", "", "the name of updated replication")
	replicationUpdateCommand.Flags().StringVarP(&replicationDesp, "description", "d", "", "the description of updated replication")
	// TODO: Add some other update items, such as status, replicatoin_period ... etc.
	replicationFailoverCommand.Flags().BoolVarP(&allowAttachedVolume, "allow_attached_volume", "a", false, "whether allow attached volume when failing over replication")
	replicationFailoverCommand.Flags().StringVarP(&secondaryBackendId, "secondary_backend_id", "s", model.ReplicationDefaultBackendId, "the secondary backend id of failoverr replication")
	replicationCommand.AddCommand(replicationShowCommand)
	replicationCommand.AddCommand(replicationListCommand)
	replicationCommand.AddCommand(replicationDeleteCommand)
	replicationCommand.AddCommand(replicationUpdateCommand)
	replicationCommand.AddCommand(replicationEnableCommand)
	replicationCommand.AddCommand(replicationDisableCommand)
	replicationCommand.AddCommand(replicationFailoverCommand)
}

func replicationAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

var replicationFormatters = FormatterList{"PrimaryReplicationDriverData": JsonFormatter,
	"SecondaryReplicationDriverData": JsonFormatter}

func replicationCreateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 2)
	validMode := []string{model.ReplicationModeSync, model.ReplicationModeAsync}
	var mode = strings.ToLower(replicationMode)
	if !utils.Contained(mode, validMode) {
		Fatalf("invalid replication mode '%s'\n", replicationMode)
	}

	prdd := map[string]string{}
	if len(primaryReplicationDriverData) != 0 {
		if err := json.Unmarshal([]byte(primaryReplicationDriverData), &prdd); err != nil {
			Debugln(err)
			Fatalln("invalid replication primary driver data")
		}
	}

	srdd := map[string]string{}
	if len(secondaryReplicationDriverData) != 0 {
		if err := json.Unmarshal([]byte(secondaryReplicationDriverData), &prdd); err != nil {
			Debugln(err)
			Fatalln("invalid replication secondary driver data")
		}
	}

	switch {
	case replicationPeriod < 0:
		Fatalf("invalid replication period '%d'\n", replicationPeriod)
	case replicationPeriod != 0 && replicationMode == model.ReplicationModeSync:
		Fatalf("no need to set replication_period when the replication mode is 'sync'\n")
	case replicationPeriod != 0:
		break
	case replicationPeriod == 0 && replicationMode == model.ReplicationModeAsync:
		replicationPeriod = model.ReplicationDefaultPeriod
	}

	replica := &model.ReplicationSpec{
		Name:                           replicationName,
		Description:                    replicationDesp,
		PrimaryVolumeId:                args[0],
		SecondaryVolumeId:              args[1],
		PrimaryReplicationDriverData:   prdd,
		SecondaryReplicationDriverData: srdd,
		ReplicationMode:                mode,
		ReplicationPeriod:              replicationPeriod,
	}

	resp, err := client.CreateReplication(replica)
	PrintResponse(resp)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "AvailabilityZone",
		"PrimaryVolumeId", "SecondaryVolumeId", "PrimaryReplicationDriverData", "SecondaryReplicationDriverData",
		"ReplicationStatus", "ReplicationMode", "ReplicationPeriod", "ProfileId"}
	PrintDict(resp, keys, replicationFormatters)
}

func replicationShowAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	resp, err := client.GetReplication(args[0])
	PrintResponse(resp)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "AvailabilityZone",
		"PrimaryVolumeId", "SecondaryVolumeId", "PrimaryReplicationDriverData", "SecondaryReplicationDriverData",
		"ReplicationStatus", "ReplicationMode", "ReplicationPeriod", "ProfileId"}
	PrintDict(resp, keys, replicationFormatters)
}

func replicationListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)

	var opts = map[string]string{"limit": repLimit, "offset": repOffset, "sortDir": repSortDir,
		"sortKey": repSortKey, "Id": repId,
		"Name": repName, "Description": repDesp, "PrimaryVolumeId": repPrimaryVolumeId,
		"SecondaryVolumeId": repSecondaryVolumeId}

	resp, err := client.ListReplications(opts)
	PrintResponse(resp)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "Name", "Description", "AvailabilityZone",
		"PrimaryVolumeId", "SecondaryVolumeId", "ReplicationStatus", "ReplicationMode"}
	PrintList(resp, keys, FormatterList{})
}

func replicationUpdateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	replica := &model.ReplicationSpec{
		Name:        replicationName,
		Description: replicationDesp,
	}

	resp, err := client.UpdateReplication(args[0], replica)
	PrintResponse(resp)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "Name", "Description", "AvailabilityZone",
		"PrimaryVolumeId", "SecondaryVolumeId", "PrimaryReplicationDriverData", "SecondaryReplicationDriverData",
		"ReplicationStatus", "ReplicationMode", "ReplicationPeriod", "ProfileId"}
	PrintDict(resp, keys, replicationFormatters)
}

func replicationDeleteAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	replicaId := args[0]
	err := client.DeleteReplication(replicaId, nil)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
}

func replicationEnableAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	replicaId := args[0]
	err := client.EnableReplication(replicaId)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
}

func replicationDisableAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	replicaId := args[0]
	err := client.DisableReplication(replicaId)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
}

func replicationFailoverAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	replicaId := args[0]
	failoverReplication := &model.FailoverReplicationSpec{
		AllowAttachedVolume: allowAttachedVolume,
		SecondaryBackendId:  secondaryBackendId,
	}
	err := client.FailoverReplication(replicaId, failoverReplication)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
}
