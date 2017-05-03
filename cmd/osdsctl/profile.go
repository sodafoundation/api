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

	profiles "github.com/opensds/opensds/pkg/controller/api"

	"github.com/spf13/cobra"
)

var profileCommand = &cobra.Command{
	Use:   "profile",
	Short: "manage OpenSDS profile resources",
	Run:   profileAction,
}

var profileShowCommand = &cobra.Command{
	Use:   "show <profile name>",
	Short: "show information of specified profile",
	Run:   profileShowAction,
}

var profileListCommand = &cobra.Command{
	Use:   "list",
	Short: "get all profile resources",
	Run:   profileListAction,
}

func init() {
	profileCommand.AddCommand(profileShowCommand)
	profileCommand.AddCommand(profileListCommand)
}

func profileAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func profileShowAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	profileRequest := profiles.ProfileRequest{}

	result, err := profiles.GetProfile(profileRequest, args[0])
	if err != nil {
		fmt.Println("Get profile resource failed: ", err)
	} else {
		rbody, _ := json.MarshalIndent(result, "", "  ")
		fmt.Printf("%s\n", string(rbody))
	}
}

func profileListAction(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	profileRequest := profiles.ProfileRequest{}

	result, err := profiles.ListProfiles(profileRequest)
	if err != nil {
		fmt.Println("List profile resources failed: ", err)
	} else {
		rbody, _ := json.MarshalIndent(result, "", "  ")
		fmt.Printf("%s\n", string(rbody))
	}
}
