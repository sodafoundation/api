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
	"strconv"

	api "github.com/opensds/opensds/pkg/api/v1"
	volumes "github.com/opensds/opensds/pkg/apiserver"

	"github.com/spf13/cobra"
)

var volumeCommand = &cobra.Command{
	Use:   "volume",
	Short: "manage volumes in the cluster",
	Run:   volumeAction,
}

var volumeCreateCommand = &cobra.Command{
	Use:   "create <size>",
	Short: "create a volume in the specified backend of OpenSDS cluster",
	Run:   volumeCreateAction,
}

var volumeShowCommand = &cobra.Command{
	Use:   "show <id>",
	Short: "show a volume in the specified backend of OpenSDS cluster",
	Run:   volumeShowAction,
}

var volumeListCommand = &cobra.Command{
	Use:   "list",
	Short: "list all volumes in the specified backend of OpenSDS cluster",
	Run:   volumeListAction,
}

var volumeDeleteCommand = &cobra.Command{
	Use:   "delete <id>",
	Short: "delete a volume in the specified backend of OpenSDS cluster",
	Run:   volumeDeleteAction,
}

var (
	falseVolumeResponse       api.VolumeResponse
	falseVolumeDetailResponse api.VolumeDetailResponse
	profileName               string
	volName                   string
)

func init() {
	volumeCommand.PersistentFlags().StringVarP(&profileName, "profile", "p", "", "the name of profile configured by admin")

	volumeCommand.AddCommand(volumeCreateCommand)
	volumeCreateCommand.Flags().StringVarP(&volName, "name", "n", "null", "the name of created volume")
	volumeCommand.AddCommand(volumeShowCommand)
	volumeCommand.AddCommand(volumeListCommand)
	volumeCommand.AddCommand(volumeDeleteCommand)
	volumeCommand.AddCommand(volumeAttachmentCommand)
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
		die("error parsing size %s: %+v", args[0], err)
	}

	volumeRequest := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			Name: volName,
			Size: int32(size),
		},
		Profile: &api.StorageProfile{
			Name: profileName,
		},
	}
	result, err := volumes.CreateVolume(volumeRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseVolumeResponse) {
			fmt.Println("Create volume failed!")
		} else {
			rbody, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("%s\n", string(rbody))
		}
	}
}

func volumeShowAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volumeRequest := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			Id: args[0],
		},
	}
	result, err := volumes.GetVolume(volumeRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseVolumeDetailResponse) {
			fmt.Printf("The volume id %s not exists!\n", args[0])
		} else {
			rbody, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("%s\n", string(rbody))
		}
	}
}

func volumeListAction(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volumeRequest := &volumes.VolumeRequest{}

	result, err := volumes.ListVolumes(volumeRequest)
	if err != nil {
		fmt.Println(err)
	}
	rbody, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("%s\n", string(rbody))
}

func volumeDeleteAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volumeRequest := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			Id: args[0],
		},
		Profile: &api.StorageProfile{
			Name: profileName,
		},
	}

	result := volumes.DeleteVolume(volumeRequest)
	rbody, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("%s\n", string(rbody))
}
