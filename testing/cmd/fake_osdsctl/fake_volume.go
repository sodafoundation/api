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

	"github.com/opensds/opensds/testing/pkg/controller/api"
	volumes "github.com/opensds/opensds/testing/pkg/controller/api/v1/fake_volumes"

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

var volumeAttachCommand = &cobra.Command{
	Use:   "attach <id>",
	Short: "attach a volume in the cluster",
	Run:   volumeAttachAction,
}

var volumeDetachCommand = &cobra.Command{
	Use:   "detach <id> <attachment id>",
	Short: "detach a volume with attachment_id in the cluster",
	Run:   volumeDetachAction,
}

var volumeMountCommand = &cobra.Command{
	Use:   "mount <target mount dir> <mount device> <volume id>",
	Short: "mount a volume in the cluster",
	Run:   volumeMountAction,
}

var volumeUnmountCommand = &cobra.Command{
	Use:   "unmount <mount dir>",
	Short: "unmount a volume in the cluster",
	Run:   volumeUnmountAction,
}

var falseVolumeResponse api.VolumeResponse
var falseVolumeDetailResponse api.VolumeDetailResponse
var falseAllVolumesResponse []api.VolumeResponse
var falseAllVolumesDetailResponse api.VolumeDetailResponse

var (
	volResourceType string
	volName         string
	volAllowDetails bool
	host            string
	attachDevice    string
	fsType          string
)

func init() {
	defaultHost, err := os.Hostname()
	if err != nil {
		panic("Can't get the host name!")
	}

	volumeCommand.PersistentFlags().StringVarP(&volResourceType, "backend", "b", "cinder", "backend resource type")
	volumeCommand.AddCommand(volumeCreateCommand)
	volumeCommand.AddCommand(volumeShowCommand)
	volumeCommand.AddCommand(volumeListCommand)
	volumeCommand.AddCommand(volumeDeleteCommand)
	volumeCommand.AddCommand(volumeAttachCommand)
	volumeCommand.AddCommand(volumeDetachCommand)
	volumeCreateCommand.Flags().StringVarP(&volName, "name", "n", "null", "the name of created volume")
	volumeListCommand.Flags().BoolVarP(&volAllowDetails, "detail", "d", false, "list volumes in details")
	volumeAttachCommand.Flags().StringVarP(&host, "host", "o", defaultHost, "the name of attaching host")
	volumeAttachCommand.Flags().StringVarP(&attachDevice, "path", "p", "/mnt", "the path of attaching device")
	volumeMountCommand.Flags().StringVarP(&fsType, "type", "t", "ext4", "the file system type")
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
		die("error parsing size %s: %v", args[0], err)
	}

	volumeRequest := volumes.VolumeRequest{
		ResourceType: volResourceType,
		Name:         volName,
		Size:         int32(size),
	}
	result, err := volumes.CreateVolume(volumeRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseVolumeResponse) {
			fmt.Println("Create volume failed!")
		} else {
			rbody, _ := json.Marshal(result)
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

	volumeRequest := volumes.VolumeRequest{
		ResourceType: volResourceType,
		Id:           args[0],
	}
	result, err := volumes.GetVolume(volumeRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseVolumeDetailResponse) {
			fmt.Println("Show volume failed!")
		} else {
			rbody, _ := json.Marshal(result)
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

	volumeRequest := volumes.VolumeRequest{
		ResourceType: volResourceType,
		AllowDetails: volAllowDetails,
	}
	result, err := volumes.ListVolumes(volumeRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseAllVolumesResponse) {
			fmt.Println("List volumes failed!")
		} else {
			rbody, _ := json.Marshal(result)
			fmt.Printf("%s\n", string(rbody))
		}
	}
}

func volumeDeleteAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volumeRequest := volumes.VolumeRequest{
		ResourceType: volResourceType,
		Id:           args[0],
	}

	result := volumes.DeleteVolume(volumeRequest)
	fmt.Printf("%v\n", result)
}

func volumeAttachAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volumeRequest := volumes.VolumeRequest{
		ResourceType: volResourceType,
		Id:           args[0],
		Host:         host,
		Device:       attachDevice,
	}

	result := volumes.AttachVolume(volumeRequest)
	fmt.Printf("%v\n", result)
}

func volumeDetachAction(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volumeRequest := volumes.VolumeRequest{
		ResourceType: volResourceType,
		Id:           args[0],
		Attachment:   args[1],
	}

	result := volumes.DetachVolume(volumeRequest)
	fmt.Printf("%v\n", result)
}

func volumeMountAction(cmd *cobra.Command, args []string) {
	if len(args) != 3 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volumeRequest := volumes.VolumeRequest{
		MountDir: args[0],
		Device:   args[1],
		Id:       args[2],
		FsType:   fsType,
	}

	result := volumes.MountVolume(volumeRequest)
	fmt.Printf("%v\n", result)
}

func volumeUnmountAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volumeRequest := volumes.VolumeRequest{
		MountDir: args[0],
	}

	result := volumes.UnmountVolume(volumeRequest)
	fmt.Printf("%v\n", result)
}
