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
	"os"
	"strings"

	"github.com/astaxie/beego"
	bctx "github.com/astaxie/beego/context"
	c "github.com/opensds/opensds/client"
	"github.com/opensds/opensds/contrib/cindercompatibleapi/converter"
	"github.com/opensds/opensds/pkg/api/filter/auth"
	"github.com/opensds/opensds/pkg/api/filter/context"
	"github.com/opensds/opensds/pkg/utils/constants"
)

var (
	client *c.Client
)

// Run ...
func Run(CinderEndPoint string) {
	ep, ok := os.LookupEnv(c.OpensdsEndpoint)
	if !ok {
		fmt.Println("ERROR: You must provide the endpoint by setting " +
			"the environment variable " + c.OpensdsEndpoint)
		return
	}
	cfg := &c.Config{Endpoint: ep}

	authStrategy, ok := os.LookupEnv(c.OpensdsAuthStrategy)
	if !ok {
		authStrategy = c.Noauth
		fmt.Println("WARNING: Not found Env " + c.OpensdsAuthStrategy + ", use default(noauth)")
	}

	switch authStrategy {
	case c.Keystone:
		cfg.AuthOptions = c.LoadKeystoneAuthOptionsFromEnv()
	case c.Noauth:
		cfg.AuthOptions = c.LoadNoAuthOptionsFromEnv()
	default:
		cfg.AuthOptions = c.NewNoauthOptions(constants.DefaultTenantId)
	}

	client = c.NewClient(cfg)
	// CinderEndPoint: http://127.0.0.1:8777/v3 http://10.10.3.173/volume/v3
	words := strings.Split(CinderEndPoint, "/")
	versionPosition := 3
	isHaveV3 := false

	for i := 0; i < len(words); i++ {
		if words[i] == converter.APIVersion {
			versionPosition = i
			isHaveV3 = true
		}
	}

	fmt.Println(versionPosition)
	if (versionPosition < 3) || (false == isHaveV3) {
		fmt.Println("The environment variable CINDER_ENDPOINT is set incorrectly")
		return
	}

	prefix := ""
	for j := 3; j <= versionPosition; j++ {
		prefix = prefix + words[j]
		if j != versionPosition {
			prefix = prefix + "/"
		}
	}

	fmt.Println(prefix)

	ns :=
		beego.NewNamespace(prefix,
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

	beego.InsertFilter("*", beego.BeforeExec, context.Factory())
	beego.InsertFilter("*", beego.BeforeExec, auth.Factory())
	beego.AddNamespace(ns)

	beego.Router("/", &VersionPortal{}, "get:ListAllAPIVersions")

	// start service
	beego.Run(words[2])
}
