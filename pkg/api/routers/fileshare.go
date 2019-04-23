package routers

import (
	"github.com/astaxie/beego"
	"github.com/opensds/opensds/pkg/api/controllers"
	"github.com/opensds/opensds/pkg/utils/constants"
)
func init() {

	// add router for fileshare api
	filens :=
	// Share is a part of files. At the same time multiple users can access the the same shares.
		beego.NewNamespace("/"+constants.APIVersion+"/:tenantId/file",
			beego.NSRouter("/shares", controllers.NewFileSharePortal(), "post:CreateFileShare;get:ListFileShares"),
			beego.NSRouter("/shares/:fileshareId", controllers.NewFileSharePortal(), "get:GetFileShare;put:UpdateFileShare;delete:DeleteFileShare"),
			// Snapshot is a point-in-time copy of the data that a FileShare contains.
			// Creates, shows, lists, unpdates and deletes snapshot.
			//beego.NSRouter("/snapshots", controllers.NewFileShareSnapshotPortal(), "post:CreateFileShareSnapshot;get:ListFileShareSnapshots"),
			//beego.NSRouter("/snapshots/:snapshotId", controllers.NewFileShareSnapshotPortal(), "get:GetFileShareSnapshot;put:UpdateFileShareSnapshot;delete:DeleteFileShareSnapshot"),
		)
	beego.AddNamespace(filens)
}

