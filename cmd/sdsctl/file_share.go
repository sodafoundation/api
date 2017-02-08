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

	"api/shares"

	"github.com/spf13/cobra"
)

var shareCommand = &cobra.Command{
	Use:   "share",
	Short: "manage shares in the cluster",
	Run:   shareAction,
}

var shareCreateCommand = &cobra.Command{
	Use:   "create <size>",
	Short: "create a share in the cluster",
	Run:   shareCreateAction,
}

var shareShowCommand = &cobra.Command{
	Use:   "show <id>",
	Short: "show a share in the cluster",
	Run:   shareShowAction,
}

var shareListCommand = &cobra.Command{
	Use:   "list",
	Short: "list shares in the cluster",
	Run:   shareListAction,
}

var shareUpdateCommand = &cobra.Command{
	Use:   "update <id>",
	Short: "update a share in the cluster",
	Run:   shareUpdateAction,
}

var shareDeleteCommand = &cobra.Command{
	Use:   "delete <id>",
	Short: "delete a share in the cluster",
	Run:   shareDeleteAction,
}

func init() {
	shareCommand.PersistentFlags().StringVarP(&resourceType, "backend", "b", "manila", "backend resource type")
	shareCommand.AddCommand(shareCreateCommand)
	shareCommand.AddCommand(shareShowCommand)
	shareCommand.AddCommand(shareListCommand)
	shareCommand.AddCommand(shareUpdateCommand)
	shareCommand.AddCommand(shareDeleteCommand)
	shareCreateCommand.Flags().StringVarP(&name, "name", "n", "null", "list shares in details")
	shareListCommand.Flags().BoolVarP(&allowDetails, "detail", "d", false, "list shares in details")
}

func shareAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func shareCreateAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	size, err := strconv.Atoi(args[0])
	if err != nil {
		die("error parsing size %s: %v", args[0], err)
	}

	result, err := shares.Create(resourceType, name, size)
	if err != nil {
		fmt.Println(err)
	} else {
		if result == "" {
			fmt.Println("Create share failed!")
		} else {
			fmt.Printf("%v\n", result)
		}
	}
}

func shareShowAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(1)
	}

	shrID := args[0]

	result, err := shares.Show(resourceType, shrID)
	if err != nil {
		fmt.Println(err)
	} else {
		if result == "" {
			fmt.Println("Show share failed!")
		} else {
			fmt.Printf("%v\n", result)
		}
	}
}

func shareListAction(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	result, err := shares.List(resourceType, allowDetails)
	if err != nil {
		fmt.Println(err)
	} else {
		if result == "" {
			fmt.Println("List shares failed!")
		} else {
			fmt.Printf("%v\n", result)
		}
	}
}

func shareUpdateAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	shrID := args[0]

	result, err := shares.Update(resourceType, shrID, name)
	if err != nil {
		fmt.Println(err)
	} else {
		if result == "" {
			fmt.Println("Update share failed!")
		} else {
			fmt.Printf("%v\n", result)
		}
	}
}

func shareDeleteAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(1)
	}

	shrID := args[0]

	result, err := shares.Delete(resourceType, shrID)
	if err != nil {
		fmt.Println(err)
	} else {
		if result == "" {
			fmt.Println("Delete share failed!")
		} else {
			fmt.Printf("%v\n", result)
		}
	}
}
