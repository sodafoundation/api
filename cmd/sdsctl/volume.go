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
	"fmt"
	"os"
	"strconv"

	"api/volumes"

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

var volumeUpdateCommand = &cobra.Command{
	Use:   "update <id>",
	Short: "update a volume in the cluster",
	Run:   volumeUpdateAction,
}

var volumeDeleteCommand = &cobra.Command{
	Use:   "delete <id>",
	Short: "delete a volume in the cluster",
	Run:   volumeDeleteAction,
}

var volumeMountCommand = &cobra.Command{
	Use:   "mount <id>",
	Short: "mount a volume in the cluster",
	Run:   volumeMountAction,
}

var volumeUnmountCommand = &cobra.Command{
	Use:   "unmount <id> <attachment_id>",
	Short: "unmount a volume with attachment_id in the cluster",
	Run:   volumeUnmountAction,
}

var (
	resourceType string
	name         string
	allowDetails bool
	host         string
	mountpoint   string
)

func init() {
	volumeCommand.PersistentFlags().StringVarP(&resourceType, "backend", "b", "cinder", "backend resource type")
	volumeCommand.AddCommand(volumeCreateCommand)
	volumeCommand.AddCommand(volumeShowCommand)
	volumeCommand.AddCommand(volumeListCommand)
	volumeCommand.AddCommand(volumeUpdateCommand)
	volumeCommand.AddCommand(volumeDeleteCommand)
	volumeCommand.AddCommand(volumeMountCommand)
	volumeCommand.AddCommand(volumeUnmountCommand)
	volumeCreateCommand.Flags().StringVarP(&name, "name", "n", "null", "the name of created volume")
	volumeListCommand.Flags().BoolVarP(&allowDetails, "detail", "d", false, "list volumes in details")
	volumeMountCommand.Flags().StringVarP(&host, "host", "h", "localhost", "the hostname mounting volume")
	volumeMountCommand.Flags().StringVarP(&mountpoint, "mountpoint", "m", "/dev/vdc", "mountpoint of volume")
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

	result, err := volumes.Create(resourceType, name, size)
	if err != nil {
		fmt.Println(err)
	} else {
		if result == "" {
			fmt.Println("Create volume failed!")
		} else {
			fmt.Printf("%v\n", result)
		}
	}
}

func volumeShowAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volID := args[0]

	result, err := volumes.Show(resourceType, volID)
	if err != nil {
		fmt.Println(err)
	} else {
		if result == "" {
			fmt.Println("Show volume failed!")
		} else {
			fmt.Printf("%v\n", result)
		}
	}
}

func volumeListAction(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	result, err := volumes.List(resourceType, allowDetails)
	if err != nil {
		fmt.Println(err)
	} else {
		if result == "" {
			fmt.Println("List volumes failed!")
		} else {
			fmt.Printf("%v\n", result)
		}
	}
}

func volumeUpdateAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volID := args[0]

	result, err := volumes.Update(resourceType, volID, name)
	if err != nil {
		fmt.Println(err)
	} else {
		if result == "" {
			fmt.Println("Update volume failed!")
		} else {
			fmt.Printf("%v\n", result)
		}
	}
}

func volumeDeleteAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volID := args[0]

	result, err := volumes.Delete(resourceType, volID)
	if err != nil {
		fmt.Println(err)
	} else {
		if result == "" {
			fmt.Println("Delete volume failed!")
		} else {
			fmt.Printf("%v\n", result)
		}
	}
}

func volumeMountAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volID := args[0]

	result, err := volumes.Mount(resourceType, volID, host, mountpoint)
	if err != nil {
		fmt.Println(err)
	} else {
		if result == "" {
			fmt.Println("Mount volume failed!")
		} else {
			fmt.Printf("%v\n", result)
		}
	}
}

func volumeUnmountAction(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volID := args[0]
	attachment := args[1]

	result, err := volumes.Unmount(resourceType, volID, attachment)
	if err != nil {
		fmt.Println(err)
	} else {
		if result == "" {
			fmt.Println("Unmount volume failed!")
		} else {
			fmt.Printf("%v\n", result)
		}
	}
}
