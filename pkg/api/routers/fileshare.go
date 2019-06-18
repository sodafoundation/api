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

package routers

import (
	"github.com/astaxie/beego"
	"github.com/opensds/opensds/pkg/api/controllers"
	"github.com/opensds/opensds/pkg/utils/constants"
)

func init() {

	// add router for file storage api
	filens :=
		// Share is a part of files. At the same time multiple users can access the the same shares.
		beego.NewNamespace("/"+constants.APIVersion+"/:tenantId/file",
			beego.NSRouter("/shares", controllers.NewFileSharePortal(), "post:CreateFileShare;get:ListFileShares"),
			beego.NSRouter("/shares/:fileshareId", controllers.NewFileSharePortal(), "get:GetFileShare;put:UpdateFileShare;delete:DeleteFileShare"),
			// Snapshot is a point-in-time copy of the data that a FileShare contains.
			// Creates, shows, lists, unpdates and deletes snapshot.
			beego.NSRouter("/snapshots", controllers.NewFileShareSnapshotPortal(), "post:CreateFileShareSnapshot;get:ListFileShareSnapshots"),
			beego.NSRouter("/snapshots/:snapshotId", controllers.NewFileShareSnapshotPortal(), "get:GetFileShareSnapshot;put:UpdateFileShareSnapshot;delete:DeleteFileShareSnapshot"),
			// Access is to set acl's for fileshare
			beego.NSRouter("/acls", controllers.NewFileSharePortal(), "post:CreateFileShareAcl;get:ListFileSharesAcl"),
			beego.NSRouter("/acls/:aclId", controllers.NewFileSharePortal(), "get:GetFileShareAcl;delete:DeleteFileShareAcl"),
		)
	beego.AddNamespace(filens)
}
