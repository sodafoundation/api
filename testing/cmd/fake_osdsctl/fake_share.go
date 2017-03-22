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
	shares "github.com/opensds/opensds/testing/pkg/controller/api/v1/fake_shares"

	"github.com/spf13/cobra"
)

var shareCommand = &cobra.Command{
	Use:   "share",
	Short: "manage shares in the cluster",
	Run:   shareAction,
}

var shareCreateCommand = &cobra.Command{
	Use:   "create <share_proto> <size>",
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

var shareDeleteCommand = &cobra.Command{
	Use:   "delete <id>",
	Short: "delete a share in the cluster",
	Run:   shareDeleteAction,
}

var falseShareResponse api.ShareResponse
var falseShareDetailResponse api.ShareDetailResponse
var falseAllSharesResponse []api.ShareResponse
var falseAllSharesDetailResponse []api.ShareDetailResponse

var (
	shrResourceType string
	shrName         string
	shrType         string
	shrAllowDetails bool
)

func init() {
	shareCommand.PersistentFlags().StringVarP(&shrResourceType, "backend", "b", "manila", "backend resource type")
	shareCommand.AddCommand(shareCreateCommand)
	shareCommand.AddCommand(shareShowCommand)
	shareCommand.AddCommand(shareListCommand)
	shareCommand.AddCommand(shareDeleteCommand)
	shareCreateCommand.Flags().StringVarP(&shrName, "name", "n", "null", "the name of created share")
	shareCreateCommand.Flags().StringVarP(&shrType, "type", "t", "default", "the type of created share")
	shareListCommand.Flags().BoolVarP(&shrAllowDetails, "detail", "d", false, "list shares in details")
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
		die("error parsing size %s: %v", args[0], err)
	}

	shareRequest := shares.ShareRequest{
		ResourceType: shrResourceType,
		Name:         shrName,
		ShareType:    shrType,
		ShareProto:   shrProto,
		Size:         int32(size),
	}
	result, err := shares.CreateShare(shareRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseShareResponse) {
			fmt.Println("Create share failed!")
		} else {
			rbody, _ := json.Marshal(result)
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
		ResourceType: shrResourceType,
		Id:           args[0],
	}
	result, err := shares.GetShare(shareRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseShareDetailResponse) {
			fmt.Println("Show share failed!")
		} else {
			rbody, _ := json.Marshal(result)
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
		ResourceType: shrResourceType,
		AllowDetails: shrAllowDetails,
	}
	result, err := shares.ListShares(shareRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseAllSharesResponse) {
			fmt.Println("List shares failed!")
		} else {
			rbody, _ := json.Marshal(result)
			fmt.Printf("%s\n", string(rbody))
		}
	}
}

func shareDeleteAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(1)
	}

	shareRequest := shares.ShareRequest{
		ResourceType: shrResourceType,
		Id:           args[0],
	}

	result := shares.DeleteShare(shareRequest)
	fmt.Printf("%v\n", result)
}
