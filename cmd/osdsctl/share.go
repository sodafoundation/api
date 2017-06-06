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
	shares "github.com/opensds/opensds/pkg/apiserver"

	"github.com/spf13/cobra"
)

var shareCommand = &cobra.Command{
	Use:   "share",
	Short: "manage shares in the specified backend of OpenSDS cluster",
	Run:   shareAction,
}

var shareCreateCommand = &cobra.Command{
	Use:   "create <share_proto> <size>",
	Short: "create a share in the specified backend of OpenSDS cluster",
	Run:   shareCreateAction,
}

var shareShowCommand = &cobra.Command{
	Use:   "show <id>",
	Short: "show a share in the specified backend of OpenSDS cluster",
	Run:   shareShowAction,
}

var shareListCommand = &cobra.Command{
	Use:   "list",
	Short: "list shares in the specified backend of OpenSDS cluster",
	Run:   shareListAction,
}

var shareDeleteCommand = &cobra.Command{
	Use:   "delete <id>",
	Short: "delete a share in the specified backend of OpenSDS cluster",
	Run:   shareDeleteAction,
}

var shareAttachCommand = &cobra.Command{
	Use:   "attach <id>",
	Short: "attach a share in the specified backend of OpenSDS cluster",
	Run:   shareAttachAction,
}

var shareDetachCommand = &cobra.Command{
	Use:   "detach <device path>",
	Short: "detach a share with device path in the specified backend of OpenSDS cluster",
	Run:   shareDetachAction,
}

var (
	falseShareResponse       api.ShareResponse
	falseShareDetailResponse api.ShareDetailResponse
	shrBackendDriver         string
	shrName                  string
)

func init() {
	shareCommand.PersistentFlags().StringVarP(&shrBackendDriver, "backend", "b", "manila", "backend resource type")
	shareCommand.AddCommand(shareCreateCommand)
	shareCommand.AddCommand(shareShowCommand)
	shareCommand.AddCommand(shareListCommand)
	shareCommand.AddCommand(shareDeleteCommand)
	shareCommand.AddCommand(shareAttachCommand)
	shareCommand.AddCommand(shareDetachCommand)
	shareCreateCommand.Flags().StringVarP(&shrName, "name", "n", "null", "the name of created share")
}

func shareAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func shareCreateAction(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	shrProto := args[0]
	size, err := strconv.Atoi(args[1])
	if err != nil {
		die("error parsing size %s: %+v", args[0], err)
	}

	shareRequest := shares.ShareRequest{
		Schema: &api.ShareOperationSchema{
			Name:       shrName,
			ShareProto: shrProto,
			Size:       int32(size),
		},
		Profile: &api.StorageProfile{
			BackendDriver: shrBackendDriver,
		},
	}
	result, err := shares.CreateShare(shareRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseShareResponse) {
			fmt.Println("Create share failed!")
		} else {
			rbody, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("%s\n", string(rbody))
		}
	}
}

func shareShowAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(1)
	}

	shareRequest := shares.ShareRequest{
		Schema: &api.ShareOperationSchema{
			Id: args[0],
		},
		Profile: &api.StorageProfile{
			BackendDriver: shrBackendDriver,
		},
	}
	result, err := shares.GetShare(shareRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseShareDetailResponse) {
			fmt.Printf("The share id %s not exists!\n", args[0])
		} else {
			rbody, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("%s\n", string(rbody))
		}
	}
}

func shareListAction(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	shareRequest := shares.ShareRequest{
		Profile: &api.StorageProfile{
			BackendDriver: shrBackendDriver,
		},
	}
	result, err := shares.ListShares(shareRequest)
	if err != nil {
		fmt.Println(err)
	}
	rbody, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("%s\n", string(rbody))
}

func shareDeleteAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(1)
	}

	shareRequest := shares.ShareRequest{
		Schema: &api.ShareOperationSchema{
			Id: args[0],
		},
		Profile: &api.StorageProfile{
			BackendDriver: shrBackendDriver,
		},
	}

	result := shares.DeleteShare(shareRequest)
	rbody, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("%s\n", string(rbody))
}

func shareAttachAction(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	shareRequest := &shares.ShareRequest{
		Schema: &api.ShareOperationSchema{
			Id: args[0],
		},
		Profile: &api.StorageProfile{
			BackendDriver: shrBackendDriver,
		},
	}

	result := shares.AttachShare(shareRequest)
	rbody, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("%s\n", string(rbody))
}

func shareDetachAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	shareRequest := shares.ShareRequest{
		Schema: &api.ShareOperationSchema{
			Device: args[0],
		},
		Profile: &api.StorageProfile{
			BackendDriver: shrBackendDriver,
		},
	}

	result := shares.DetachShare(shareRequest)
	rbody, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("%s\n", string(rbody))
}
