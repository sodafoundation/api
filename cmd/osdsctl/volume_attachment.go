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
	volumes "github.com/opensds/opensds/pkg/apiserver"

	"github.com/spf13/cobra"
)

var volumeAttachmentCommand = &cobra.Command{
	Use:   "attachment",
	Short: "manage volume attachments in the cluster",
	Run:   volumeAttachmentAction,
}

var volumeAttachmentCreateCommand = &cobra.Command{
	Use:   "create <volume id>",
	Short: "create an attachment of specified volume in the specified backend of OpenSDS cluster",
	Run:   volumeAttachmentCreateAction,
}

var volumeAttachmentShowCommand = &cobra.Command{
	Use:   "show <volume id> <attachment id>",
	Short: "show an attachment of specified volume in the specified backend of OpenSDS cluster",
	Run:   volumeAttachmentShowAction,
}

var volumeAttachmentListCommand = &cobra.Command{
	Use:   "list <volume id>",
	Short: "list all attachments of specified volume in the specified backend of OpenSDS cluster",
	Run:   volumeAttachmentListAction,
}

var volumeAttachmentUpdateCommand = &cobra.Command{
	Use:   "update <volume id> <attachment id> <device path>",
	Short: "Update an attachment of specified volume in the specified backend of OpenSDS cluster",
	Run:   volumeAttachmentUpdateAction,
}

var volumeAttachmentDeleteCommand = &cobra.Command{
	Use:   "delete <volume id> <attachment id>",
	Short: "delete an attachment of specified volume in the specified backend of OpenSDS cluster",
	Run:   volumeAttachmentDeleteAction,
}

var (
	falseVolumeAttachment api.VolumeAttachment
	doLocalAttach         bool
	multiPath             bool
)

func init() {
	volumeAttachmentCommand.AddCommand(volumeAttachmentCreateCommand)
	volumeAttachmentCreateCommand.Flags().BoolVarP(&doLocalAttach, "local", "l", false, "specify if it is a local volume")
	volumeAttachmentCreateCommand.Flags().BoolVarP(&multiPath, "multipath", "m", false, "specify if the volume can be attached to multiple host")
	volumeAttachmentCommand.AddCommand(volumeAttachmentShowCommand)
	volumeAttachmentCommand.AddCommand(volumeAttachmentListCommand)
	volumeAttachmentCommand.AddCommand(volumeAttachmentUpdateCommand)
	volumeAttachmentCommand.AddCommand(volumeAttachmentDeleteCommand)
}

func volumeAttachmentAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func volumeAttachmentCreateAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volumeRequest := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			Id:            args[0],
			DoLocalAttach: doLocalAttach,
			MultiPath:     multiPath,
			HostInfo:      api.HostInfo{},
		},
	}
	result, err := volumes.CreateVolumeAttachment(volumeRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseVolumeAttachment) {
			fmt.Println("Create volume attachment failed!")
		} else {
			rbody, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("%s\n", string(rbody))
		}
	}
}

func volumeAttachmentShowAction(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volumeRequest := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			Id:           args[0],
			AttachmentId: args[1],
		},
	}
	result, err := volumes.GetVolumeAttachment(volumeRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseVolumeAttachment) {
			fmt.Printf("The volume id %s or attachment id %s not exists!\n", args[0], args[1])
		} else {
			rbody, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("%s\n", string(rbody))
		}
	}
}

func volumeAttachmentListAction(cmd *cobra.Command, args []string) {
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
	result, err := volumes.ListVolumeAttachments(volumeRequest)
	if err != nil {
		fmt.Println(err)
	}
	rbody, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("%s\n", string(rbody))
}

func volumeAttachmentUpdateAction(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volumeRequest := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			Id:           args[0],
			AttachmentId: args[1],
			Mountpoint:   args[2],
			HostInfo:     api.HostInfo{},
		},
	}
	result, err := volumes.UpdateVolumeAttachment(volumeRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		if reflect.DeepEqual(result, falseVolumeAttachment) {
			fmt.Println("Update volume attachment failed!")
		} else {
			rbody, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("%s\n", string(rbody))
		}
	}
}

func volumeAttachmentDeleteAction(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Println("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}

	volumeRequest := &volumes.VolumeRequest{
		Schema: &api.VolumeOperationSchema{
			Id:           args[0],
			AttachmentId: args[1],
		},
	}

	result := volumes.DeleteVolumeAttachment(volumeRequest)
	rbody, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("%s\n", string(rbody))
}
