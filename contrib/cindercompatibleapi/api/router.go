// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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
This module implements a entry into the OpenSDS northbound REST service.

*/

package api

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/astaxie/beego"

	bctx "github.com/astaxie/beego/context"
	log "github.com/golang/glog"
	c "github.com/opensds/opensds/client"
	"github.com/opensds/opensds/contrib/cindercompatibleapi/converter"
	"github.com/opensds/opensds/pkg/utils/constants"
)

var (
	OpensdsEndPoint string
	OpensdsClient   *c.Client
)

// ErrorSpec describes Detailed HTTP error response, which consists of a HTTP
// status code, and a custom error message unique for each failure case.
type ErrorSpec struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Run ...
func Run(CinderEndPoint string) {
	OpensdsEndPointFromEnv, ok := os.LookupEnv(c.OpensdsEndpoint)
	if !ok {
		fmt.Println("ERROR: You must provide the endpoint by setting " +
			"the environment variable " + c.OpensdsEndpoint)
		return
	}

	if "" == OpensdsEndPointFromEnv {
		OpensdsEndPoint = constants.DefaultOpensdsEndpoint
		fmt.Println("Warnning: OpenSDS Endpoint is not specified using the default value:" + OpensdsEndPoint)
	} else {
		OpensdsEndPoint = OpensdsEndPointFromEnv
	}

	// CinderEndPoint: http://127.0.0.1:8777/v3
	words := strings.Split(CinderEndPoint, "/")
	if (len(words) < 4) || (words[3] != converter.APIVersion) {
		fmt.Println("The environment variable CINDER_ENDPOINT is set incorrectly")
		return
	}

	ns :=
		beego.NewNamespace("/"+converter.APIVersion,
			beego.NSCond(func(ctx *bctx.Context) bool {
				// To judge whether the scheme is legal or not.
				if ctx.Input.Scheme() != "http" && ctx.Input.Scheme() != "https" {
					return false
				}
				return true
			}),
			beego.NSNamespace("/:projectId",
				beego.NSRouter("/types", &TypePortal{}, "post:CreateType;get:ListTypes"),
				beego.NSRouter("/types/:volumeTypeId", &TypePortal{}, "get:GetType;put:UpdateType;delete:DeleteType"),
				beego.NSRouter("/types/:volumeTypeId/extra_specs", &TypePortal{}, "post:AddExtraProperty;get:ListExtraProperties"),
				beego.NSRouter("/types/:volumeTypeId/extra_specs/:key", &TypePortal{}, "get:ShowExtraProperty;put:UpdateExtraProperty;delete:DeleteExtraProperty"),

				beego.NSRouter("/volumes", &VolumePortal{}, "post:CreateVolume;get:ListVolumes"),
				beego.NSRouter("/volumes/detail", &VolumePortal{}, "get:ListVolumesDetails"),
				beego.NSRouter("/volumes/:volumeId", &VolumePortal{}, "get:GetVolume;delete:DeleteVolume;put:UpdateVolume"),
				beego.NSRouter("/volumes/:volumeId/action", &VolumePortal{}, "post:VolumeAction"),

				beego.NSRouter("/attachments", &AttachmentPortal{}, "post:CreateAttachment;get:ListAttachments"),
				beego.NSRouter("/attachments/detail", &AttachmentPortal{}, "get:ListAttachmentsDetails"),
				beego.NSRouter("/attachments/:attachmentId", &AttachmentPortal{}, "get:GetAttachment;delete:DeleteAttachment;put:UpdateAttachment"),

				beego.NSRouter("/snapshots", &SnapshotPortal{}, "post:CreateSnapshot;get:ListSnapshots"),
				beego.NSRouter("/snapshots/detail", &SnapshotPortal{}, "get:ListSnapshotsDetails"),
				beego.NSRouter("/snapshots/:snapshotId", &SnapshotPortal{}, "get:GetSnapshot;delete:DeleteSnapshot;put:UpdateSnapshot"),
			),
		)

	beego.AddNamespace(ns)
	beego.Router("/", &VersionPortal{}, "get:ListAllAPIVersions")

	// start service
	beego.Run(words[2])
}

// NewClient method creates a new Client.
func NewClient(Ctx *bctx.Context) {
	tenantId := GetProjectId(Ctx)
	tokenID := Ctx.Input.Header(constants.AuthTokenHeader)

	if len(tenantId) > 0 && len(tokenID) > 0 {
		r := &c.KeystoneReciver{Auth: &c.KeystoneAuthOptions{TenantID: tokenID,
			TokenID: tokenID}}

		OpensdsClient = &c.Client{
			ProfileMgr:     c.NewProfileMgr(r, OpensdsEndPoint, tenantId),
			DockMgr:        c.NewDockMgr(r, OpensdsEndPoint, tenantId),
			PoolMgr:        c.NewPoolMgr(r, OpensdsEndPoint, tenantId),
			VolumeMgr:      c.NewVolumeMgr(r, OpensdsEndPoint, tenantId),
			VersionMgr:     c.NewVersionMgr(r, OpensdsEndPoint, tenantId),
			ReplicationMgr: c.NewReplicationMgr(r, OpensdsEndPoint, tenantId),
		}
	} else {
		log.Error("Failed to create a client, TenantId:" + tenantId + ", " +
			"TokenID:" + tokenID + "!!!")
	}
}

// GetProjectId Get the value of project_id
func GetProjectId(Ctx *bctx.Context) string {
	u, err := url.Parse(Ctx.Request.URL.String())
	if err != nil {
		log.Error("url.Parse failed:" + err.Error())
		return ""
	}

	// /v3/{project_id}/
	words := strings.Split(u.Path, "/")

	if (len(words) >= 3) && (len(words[2]) > 0) {
		log.V(5).Info("project_id is" + words[2])
		return words[2]
	} else {
		log.Error("there is no project_id in the path:" + u.Path)
		return ""
	}
}
