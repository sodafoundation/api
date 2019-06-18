// Copyright 2017 The OpenSDS Authors.
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
	"os"

	"github.com/opensds/opensds/pkg/model"
	"github.com/spf13/cobra"
)

var volumeAttachmentCommand = &cobra.Command{
	Use:   "attachment",
	Short: "manage volume attachments in the cluster",
	Run:   volumeAttachmentAction,
}

var volumeAttachmentCreateCommand = &cobra.Command{
	Use:   "create <attachment info>",
	Short: "create an attachment of specified volume in the cluster",
	Run:   volumeAttachmentCreateAction,
}

var volumeAttachmentShowCommand = &cobra.Command{
	Use:   "show <attachment id>",
	Short: "show a volume attachment in the cluster",
	Run:   volumeAttachmentShowAction,
}

var volumeAttachmentListCommand = &cobra.Command{
	Use:   "list",
	Short: "list all volume attachments in the cluster",
	Run:   volumeAttachmentListAction,
}

var volumeAttachmentDeleteCommand = &cobra.Command{
	Use:   "delete <attachment id>",
	Short: "delete a volume attachment of specified volume in the cluster",
	Run:   volumeAttachmentDeleteAction,
}

var volumeAttachmentUpdateCommand = &cobra.Command{
	Use:   "update <attachment id> <attachment info>",
	Short: "update a volume attachment in the cluster",
	Run:   volumeAttachmentUpdateAction,
}

var (
	volAtmLimit      string
	volAtmOffset     string
	volAtmSortDir    string
	volAtmSortKey    string
	volAtmId         string
	volAtmUserId     string
	volAtmVolumeId   string
	volAtmMountpoint string
	volAtmStatus     string
)

func init() {
	volumeAttachmentListCommand.Flags().StringVarP(&volAtmLimit, "limit", "", "50", "the number of ertries displayed per page")
	volumeAttachmentListCommand.Flags().StringVarP(&volAtmOffset, "offset", "", "0", "all requested data offsets")
	volumeAttachmentListCommand.Flags().StringVarP(&volAtmSortDir, "sortDir", "", "desc", "the sort direction of all requested data. supports asc or desc(default)")
	volumeAttachmentListCommand.Flags().StringVarP(&volAtmSortKey, "sortKey", "", "id",
		"the sort key of all requested data. supports id(default), volumeid, status, userid, tenantid")
	volumeAttachmentListCommand.Flags().StringVarP(&volAtmId, "id", "", "", "list volume attachment by id")
	volumeAttachmentListCommand.Flags().StringVarP(&volAtmUserId, "userId", "", "", "list volume attachment by storage userId")
	volumeAttachmentListCommand.Flags().StringVarP(&volAtmVolumeId, "volumeId", "", "", "list volume attachment by volumeId")
	volumeAttachmentListCommand.Flags().StringVarP(&volAtmStatus, "status", "", "", "list volume attachment by status")
	volumeAttachmentListCommand.Flags().StringVarP(&volAtmMountpoint, "mountpoint", "", "", "list volume attachment by mountpoint")

	volumeAttachmentCommand.AddCommand(volumeAttachmentCreateCommand)
	volumeAttachmentCommand.AddCommand(volumeAttachmentShowCommand)
	volumeAttachmentCommand.AddCommand(volumeAttachmentListCommand)
	volumeAttachmentCommand.AddCommand(volumeAttachmentDeleteCommand)
	volumeAttachmentCommand.AddCommand(volumeAttachmentUpdateCommand)
}

func volumeAttachmentAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

var attachmentFormatters = FormatterList{"HostInfo": JsonFormatter, "ConnectionInfo": JsonFormatter,
	"Metadata": JsonFormatter}

func volumeAttachmentCreateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	attachment := &model.VolumeAttachmentSpec{}
	if err := json.Unmarshal([]byte(args[0]), attachment); err != nil {
		Errorln(err)
		cmd.Usage()
		os.Exit(1)
	}
	resp, err := client.CreateVolumeAttachment(attachment)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "HostInfo", "ConnectionInfo", "Mountpoint",
		"Status", "VolumeId", "AttachMode", "Metadata"}
	PrintDict(resp, keys, attachmentFormatters)
}

func volumeAttachmentShowAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	resp, err := client.GetVolumeAttachment(args[0])
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "CreatedAt", "UpdatedAt", "TenantId", "UserId", "HostInfo", "ConnectionInfo",
		"Mountpoint", "Status", "VolumeId", "AccessProtocol", "AttachMode", "Metadata"}
	PrintDict(resp, keys, attachmentFormatters)
}

func volumeAttachmentListAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 0)

	var opts = map[string]string{"limit": volAtmLimit, "offset": volAtmOffset,
		"sortDir": volAtmSortDir, "sortKey": volAtmSortKey, "Id": volAtmId,
		"UserId": volAtmUserId, "VolumeId": volAtmVolumeId,
		"Status": volAtmStatus, "Mountpoint": volAtmMountpoint}

	resp, err := client.ListVolumeAttachments(opts)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "Mountpoint", "Status", "VolumeId", "AccessProtocol"}
	PrintList(resp, keys, attachmentFormatters)
}

func volumeAttachmentDeleteAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 1)
	attachment := &model.VolumeAttachmentSpec{}
	err := client.DeleteVolumeAttachment(args[0], attachment)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
}

func volumeAttachmentUpdateAction(cmd *cobra.Command, args []string) {
	ArgsNumCheck(cmd, args, 2)
	attachment := &model.VolumeAttachmentSpec{}
	if err := json.Unmarshal([]byte(args[1]), attachment); err != nil {
		Errorln(err)
		cmd.Usage()
		os.Exit(1)
	}

	resp, err := client.UpdateVolumeAttachment(args[0], attachment)
	if err != nil {
		Fatalln(HttpErrStrip(err))
	}
	keys := KeyList{"Id", "UpdatedAt", "HostInfo", "ConnectionInfo", "Mountpoint",
		"Status", "VolumeId", "AttachMode", "Metadata"}
	PrintDict(resp, keys, attachmentFormatters)
}
