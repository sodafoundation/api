// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

/*
This module implements a entry into the OpenSDS service.

*/

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	api "github.com/opensds/opensds/pkg/api/v1"
	volumes "github.com/opensds/opensds/pkg/controller/api"

	"github.com/spf13/cobra"
)

var volumeSnapshotCommand = &cobra.Command{
	Use:   "snapshot",
	Short: "manage volume snapshots in the cluster",
	Run:   volumeSnapshotAction,
}

var volumeSnapshotCreateCommand = &cobra.Command{
	Use:   "create <volume id>",
	Short: "create a snapshot of specified volume in the specified backend of OpenSDS cluster",
	Run:   volumeSnapshotCreateAction,
}

var volumeSnapshotShowCommand = &cobra.Command{
	Use:   "show <snapshot id>",
	Short: "show a volume snapshot in the specified backend of OpenSDS cluster",
	Run:   volumeSnapshotShowAction,
}

var volumeSnapshotListCommand = &cobra.Command{
	Use:   "list",
	Short: "list all volume snapshots in the specified backend of OpenSDS cluster",
	Run:   volumeSnapshotListAction,
}

var volumeSnapshotDeleteCommand = &cobra.Command{
	Use:   "delete <snapshot id>",
	Short: "delete a volume snapshot in the specified backend of OpenSDS cluster",
	Run:   volumeSnapshotDeleteAction,
}

var falseVolumeSnapshotResponse api.VolumeSnapshotResponse
var falseVolumeSnapshotsResponse []api.VolumeSnapshotResponse

var (
	volSnapshotName        string
	volSnapshotDescription string
	volForceSnapshoted     bool
)

func init() {
	volumeSnapshotCommand.AddCommand(volumeSnapshotCreateCommand)
	volumeSnapshotCreateCommand.Flags().StringVarP(&volSnapshotName, "name", "n", "null", "the name of created volume snapshot")
	volumeSnapshotCreateCommand.Flags().StringVarP(&volSnapshotDescription, "description", "d", "", "description of created volume snapshot")
	volumeSnapshotCreateCommand.Flags().BoolVarP(&volForceSnapshoted, "force", "f", true, "create a snapshot by force")
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

	volumeRequest := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			SnapshotName:    volSnapshotName,
			Id:              args[0],
			Description:     volSnapshotDescription,
			ForceSnapshoted: volForceSnapshoted,
		},
		Profile: &api.StorageProfile{
			BackendDriver: volBackendDriver,
		},
	}
	result, err := volumes.CreateVolumeSnapshot(volumeRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseVolumeSnapshotResponse) {
			fmt.Println("Create volume snapshot failed!")
		} else {
			rbody, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("%s\n", string(rbody))
		}
	}
}

func volumeSnapshotShowAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volumeRequest := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			SnapshotId: args[0],
		},
		Profile: &api.StorageProfile{
			BackendDriver: volBackendDriver,
		},
	}
	result, err := volumes.GetVolumeSnapshot(volumeRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseVolumeSnapshotResponse) {
			fmt.Println("Show volume snapshot failed!")
		} else {
			rbody, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("%s\n", string(rbody))
		}
	}
}

func volumeSnapshotListAction(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volumeRequest := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{},
		Profile: &api.StorageProfile{
			BackendDriver: volBackendDriver,
		},
	}
	result, err := volumes.ListVolumeSnapshots(volumeRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseVolumeSnapshotsResponse) {
			fmt.Println("List volume snapshots failed!")
		} else {
			rbody, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("%s\n", string(rbody))
		}
	}
}

func volumeSnapshotDeleteAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volumeRequest := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			SnapshotId: args[0],
		},
		Profile: &api.StorageProfile{
			BackendDriver: volBackendDriver,
		},
	}

	result := volumes.DeleteVolumeSnapshot(volumeRequest)
	rbody, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("%s\n", string(rbody))
}
