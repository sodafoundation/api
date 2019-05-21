// Copyright (c) 2019 The OpenSDS Authors.
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
	"encoding/json"
	"log"
	"os"

	"github.com/opensds/opensds/pkg/model"
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
	shareAclLimit            string
	shareAclOffset           string
	shareAclSortDir          string
	shareAclSortKey          string
	shareAclID               string
	shareAclTenantID         string
	shareAclFileShareId      string
	shareAclType             string
	shareAclAccessCapability string
	shareAclAccessTo         string
	shareAclDesp             string

	shareAclKeys = KeyList{"Id", "CreatedAt", "UpdatedAt", "TenantId", "FileShareId",
		"Type", "AccessCapability", "AccessTo", "Description"}
)

func init() {
	fileShareAclCommand.AddCommand(fileShareAclCreateCommand)
	fileShareAclCommand.AddCommand(fileShareAclDeleteCommand)
	fileShareAclCommand.AddCommand(fileShareAclShowCommand)
	fileShareAclCommand.AddCommand(fileShareAclListCommand)

	fileShareAclCreateCommand.Flags().StringVarP(&shareAclType, "type", "t", "", "the type of access")
	fileShareAclCreateCommand.Flags().StringVarP(&shareAclAccessCapability, "capability", "c", "", "the accessCapability for fileshare")
	fileShareAclCreateCommand.Flags().StringVarP(&shareAclAccessTo, "aclTo", "a", "", "accessTo of the fileshare")
	fileShareAclCreateCommand.Flags().StringVarP(&shareAclDesp, "description", "d", "", "the description of of the fileshare acl")

	fileShareAclListCommand.Flags().StringVarP(&shareAclLimit, "limit", "", "50", "the number of ertries displayed per page")
	fileShareAclListCommand.Flags().StringVarP(&shareAclOffset, "offset", "", "0", "all requested data offsets")
	fileShareAclListCommand.Flags().StringVarP(&shareAclSortDir, "sortDir", "", "desc", "the sort direction of all requested data. supports asc or desc(default)")
	fileShareAclListCommand.Flags().StringVarP(&shareAclSortKey, "sortKey", "", "id",
		"the sort key of all requested data. supports id(default), tenantId, fileshareId, type, accessCapability, accessTo, description")
	fileShareAclListCommand.Flags().StringVarP(&shareAclID, "id", "", "", "list fileshare acls by id")
	fileShareAclListCommand.Flags().StringVarP(&shareAclTenantID, "tenantId", "", "", "list fileshare acls by tenantId")
	fileShareAclListCommand.Flags().StringVarP(&fileshareID, "fileshareId", "", "", "list fileshare acls by fileshareId")
	fileShareAclListCommand.Flags().StringVarP(&shareAclType, "type", "", "", "list fileshare acls by type")
	fileShareAclListCommand.Flags().StringVarP(&shareAclAccessCapability, "accessCapability", "", "", "list fileshare acls by accessCapability")
	fileShareAclListCommand.Flags().StringVarP(&shareAclAccessTo, "accessTo", "", "", "list fileshare acls by accessTo")
	fileShareAclListCommand.Flags().StringVarP(&shareAclDesp, "description", "", "", "list fileshare acls by description")
}

func fileShareAclAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func fileShareAclCreateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)

	var accessCapability []string
	if "" != shareAclAccessCapability {
		err := json.Unmarshal([]byte(shareAclAccessCapability), &accessCapability)
		if err != nil {
			log.Fatalf("error parsing accessCapability %s: %+v", shareAclAccessCapability, err)
		}
	}

	var accessTo []string
	if "" != shareAclAccessTo {
		err := json.Unmarshal([]byte(shareAclAccessTo), &accessTo)
		if err != nil {
			log.Fatalf("error parsing accessTo %s: %+v", shareAclAccessTo, err)
		}
	}

	acl := &model.FileShareAclSpec{
		FileShareId:      args[0],
		Type:             shareAclType,
		AccessCapability: accessCapability,
		AccessTo:         accessTo,
		Description:      shareAclDesp,
	}

	resp, err := client.CreateFileShareAcl(acl)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}

	PrintDict(resp, shareAclKeys, FormatterList{})
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

	PrintDict(resp, shareAclKeys, FormatterList{})
}

func fileSharesAclListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)

	var opts = map[string]string{"limit": shareAclLimit, "offset": shareAclOffset, "sortDir": shareAclSortDir,
		"sortKey": shareAclSortKey, "Id": shareAclID,
		"TenantId": shareAclTenantID, "FileshareID": fileshareID, "ShareAclType": shareAclType,
		"ShareAclAccessCapability": shareAclAccessCapability, "ShareAclAccessTo": shareAclAccessTo,
		"ShareAclDesp": shareAclDesp}

	resp, err := client.ListFileSharesAcl(opts)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}

	PrintList(resp, shareAclKeys, FormatterList{})
}
