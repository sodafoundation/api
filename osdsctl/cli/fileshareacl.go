// Copyright 2019 The OpenSDS Authors.
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

	"github.com/sodafoundation/api/pkg/model"
	"github.com/spf13/cobra"
)

var fileShareAclCommand = &cobra.Command{
	Use:   "acl",
	Short: "manage fileshare acls in the cluster",
	Run:   fileShareAclAction,
}

var fileShareAclCreateCommand = &cobra.Command{
	Use:   "create <fileshare id>",
	Short: "create a acl of specified fileshare in the cluster",
	Example: "osdsctl fileshare acl create -a 10.0.0.10 -c \"Write\" -t ip 87be9ce5-6ecc-4ac3-8d6c-f5a58c9110e4",
	Run:   fileShareAclCreateAction,
}

var fileShareAclDeleteCommand = &cobra.Command{
	Use:   "delete <fileshare acl id>",
	Short: "delete a fileshare acl of specified fileshare in the cluster",
	Run:   fileShareAclDeleteAction,
}

var fileShareAclShowCommand = &cobra.Command{
	Use:   "show <fileshare acl id>",
	Short: "show a fileshare acl in the cluster",
	Run:   fileShareAclShowAction,
}

var fileShareAclListCommand = &cobra.Command{
	Use:   "list",
	Short: "list all fileshare acls in the cluster",
	Run:   fileSharesAclListAction,
}

var (
	shareAclType             string
	shareAclAccessCapability []string
	shareAclAccessTo         string
	shareAclDesp             string

	shareAclFormatters = FormatterList{"Metadata": JsonFormatter}
)

func init() {
	fileShareAclCommand.AddCommand(fileShareAclCreateCommand)
	fileShareAclCommand.AddCommand(fileShareAclDeleteCommand)
	fileShareAclCommand.AddCommand(fileShareAclShowCommand)
	fileShareAclCommand.AddCommand(fileShareAclListCommand)

	fileShareAclCreateCommand.Flags().StringVarP(&shareAclType, "type", "t", "", "the type of access. The Only current supported type is: ip")
	fileShareAclCreateCommand.Flags().StringSliceVarP(&shareAclAccessCapability, "capability", "c", shareAclAccessCapability, "the accessCapability \"Read\" or \"Write\" for fileshare")
	fileShareAclCreateCommand.Flags().StringVarP(&shareAclAccessTo, "accessTo", "a", "", "accessTo of the fileshare. A valid IPv4 format is supported")
	fileShareAclCreateCommand.Flags().StringVarP(&shareAclDesp, "description", "d", "", "the description of of the fileshare acl")
}

func fileShareAclAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func fileShareAclCreateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	acl := &model.FileShareAclSpec{
		FileShareId:      args[0],
		Type:             shareAclType,
		AccessCapability: shareAclAccessCapability,
		AccessTo:         shareAclAccessTo,
		Description:      shareAclDesp,
	}

	resp, err := client.CreateFileShareAcl(acl)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}

	keys := KeyList{"Id", "CreatedAt", "TenantId", "FileShareId", "Metadata",
		"Type", "AccessCapability", "AccessTo", "Description"}
	PrintDict(resp, keys, shareAclFormatters)
}

func fileShareAclDeleteAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)

	err := client.DeleteFileShareAcl(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
}

func fileShareAclShowAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	resp, err := client.GetFileShareAcl(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}

	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "TenantId", "FileShareId",
		"Type", "AccessCapability", "AccessTo", "Description", "Metadata"}
	PrintDict(resp, keys, shareAclFormatters)
}

func fileSharesAclListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)
	resp, err := client.ListFileSharesAcl()
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}

	keys := KeyList{"Id", "FileShareId",
		"Type", "AccessCapability", "AccessTo", "Description"}
	PrintList(resp, keys, shareAclFormatters)
}
