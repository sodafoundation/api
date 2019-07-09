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

package routers

import (
	"github.com/astaxie/beego"
	"github.com/opensds/opensds/pkg/api/controllers"
	"github.com/opensds/opensds/pkg/utils/constants"
)

func init() {

	// add router for block storage api
	blockns :=
		beego.NewNamespace("/"+constants.APIVersion+"/:tenantId/block",

			// Volume is the logical description of a piece of storage, which can be directly used by users.
			// All operations of volume can be used for both admin and users.
			beego.NSRouter("/volumes", controllers.NewVolumePortal(), "post:CreateVolume;get:ListVolumes"),
			beego.NSRouter("/volumes/:volumeId", controllers.NewVolumePortal(), "get:GetVolume;put:UpdateVolume;delete:DeleteVolume"),
			// Extend Volume
			beego.NSRouter("/volumes/:volumeId/resize", controllers.NewVolumePortal(), "post:ExtendVolume"),

			// Creates, shows, lists, unpdates and deletes attachment.
			beego.NSRouter("/attachments", controllers.NewVolumeAttachmentPortal(), "post:CreateVolumeAttachment;get:ListVolumeAttachments"),
			beego.NSRouter("/attachments/:attachmentId", controllers.NewVolumeAttachmentPortal(), "get:GetVolumeAttachment;put:UpdateVolumeAttachment;delete:DeleteVolumeAttachment"),

			// Snapshot is a point-in-time copy of the data that a volume contains.
			// Creates, shows, lists, unpdates and deletes snapshot.
			beego.NSRouter("/snapshots", controllers.NewVolumeSnapshotPortal(), "post:CreateVolumeSnapshot;get:ListVolumeSnapshots"),
			beego.NSRouter("/snapshots/:snapshotId", controllers.NewVolumeSnapshotPortal(), "get:GetVolumeSnapshot;put:UpdateVolumeSnapshot;delete:DeleteVolumeSnapshot"),

			// Creates, shows, lists, unpdates and deletes replication.
			beego.NSRouter("/replications", controllers.NewReplicationPortal(), "post:CreateReplication;get:ListReplications"),
			beego.NSRouter("/replications/detail", controllers.NewReplicationPortal(), "get:ListReplicationsDetail"),
			beego.NSRouter("/replications/:replicationId", controllers.NewReplicationPortal(), "get:GetReplication;put:UpdateReplication;delete:DeleteReplication"),
			beego.NSRouter("/replications/:replicationId/enable", controllers.NewReplicationPortal(), "post:EnableReplication"),
			beego.NSRouter("/replications/:replicationId/disable", controllers.NewReplicationPortal(), "post:DisableReplication"),
			beego.NSRouter("/replications/:replicationId/failover", controllers.NewReplicationPortal(), "post:FailoverReplication"),
			// Volume group contains a list of volumes that are used in the same application.
			beego.NSRouter("/volumeGroups", controllers.NewVolumeGroupPortal(), "post:CreateVolumeGroup;get:ListVolumeGroups"),
			beego.NSRouter("/volumeGroups/:groupId", controllers.NewVolumeGroupPortal(), "put:UpdateVolumeGroup;get:GetVolumeGroup;delete:DeleteVolumeGroup"),
		)
	beego.AddNamespace(blockns)
}
