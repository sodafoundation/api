// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	bctx "github.com/astaxie/beego/context"
	"github.com/opensds/opensds/pkg/api/filter/accesslog"
	"github.com/opensds/opensds/pkg/api/filter/auth"
	"github.com/opensds/opensds/pkg/api/filter/context"
	cfg "github.com/opensds/opensds/pkg/utils/config"
	"github.com/opensds/opensds/pkg/utils/constants"
)

const (
	AddressIdx = iota
	PortIdx
)

const (
	StatusOK       = http.StatusOK
	StatusAccepted = http.StatusAccepted
)

func Run(apiServerCfg cfg.OsdsApiServer) {

	// add router for v1beta api
	ns :=
		beego.NewNamespace("/"+constants.APIVersion,
			beego.NSCond(func(ctx *bctx.Context) bool {
				// To judge whether the scheme is legal or not.
				if ctx.Input.Scheme() != "http" && ctx.Input.Scheme() != "https" {
					return false
				}

				return true
			}),

			// List all dock services, including a list of dock object
			beego.NSRouter("/:tenantId/metrics", NewMetricsPortal(), "post:CollectMetrics"),

			// List all dock services, including a list of dock object
			beego.NSRouter("/:tenantId/docks", &DockPortal{}, "get:ListDocks"),
			// Show one dock service, including endpoint, driverName and so on
			beego.NSRouter("/:tenantId/docks/:dockId", &DockPortal{}, "get:GetDock"),

			// Profile is a set of policies configured by admin and provided for users
			// CreateProfile, UpdateProfile and DeleteProfile are used for admin only
			// ListProfiles and GetProfile are used for both admin and users
			beego.NSRouter("/:tenantId/profiles", &ProfilePortal{}, "post:CreateProfile;get:ListProfiles"),
			beego.NSRouter("/:tenantId/profiles/:profileId", &ProfilePortal{}, "get:GetProfile;put:UpdateProfile;delete:DeleteProfile"),

			// All operations of customProperties are used for admin only
			beego.NSRouter("/:tenantId/profiles/:profileId/customProperties", &ProfilePortal{}, "post:AddCustomProperty;get:ListCustomProperties"),
			beego.NSRouter("/:tenantId/profiles/:profileId/customProperties/:customKey", &ProfilePortal{}, "delete:RemoveCustomProperty"),

			// Pool is the virtual description of backend storage, usually divided into block, file and object,
			// and every pool is atomic, which means every pool contains a specific set of features.
			// ListPools and GetPool are used for checking the status of backend pool, admin only
			beego.NSRouter("/:tenantId/pools", &PoolPortal{}, "get:ListPools"),
			beego.NSRouter("/:tenantId/pools/:poolId", &PoolPortal{}, "get:GetPool"),
			beego.NSRouter("/:tenantId/availabilityZones", &PoolPortal{}, "get:ListAvailabilityZones"),

			beego.NSNamespace("/:tenantId/block",

				// Volume is the logical description of a piece of storage, which can be directly used by users.
				// All operations of volume can be used for both admin and users.
				beego.NSRouter("/volumes", NewVolumePortal(), "post:CreateVolume;get:ListVolumes"),
				beego.NSRouter("/volumes/:volumeId", NewVolumePortal(), "get:GetVolume;put:UpdateVolume;delete:DeleteVolume"),
				// Extend Volume
				beego.NSRouter("/volumes/:volumeId/resize", NewVolumePortal(), "post:ExtendVolume"),

				// Creates, shows, lists, unpdates and deletes attachment.
				beego.NSRouter("/attachments", NewVolumeAttachmentPortal(), "post:CreateVolumeAttachment;get:ListVolumeAttachments"),
				beego.NSRouter("/attachments/:attachmentId", NewVolumeAttachmentPortal(), "get:GetVolumeAttachment;put:UpdateVolumeAttachment;delete:DeleteVolumeAttachment"),

				// Snapshot is a point-in-time copy of the data that a volume contains.
				// Creates, shows, lists, unpdates and deletes snapshot.
				beego.NSRouter("/snapshots", NewVolumeSnapshotPortal(), "post:CreateVolumeSnapshot;get:ListVolumeSnapshots"),
				beego.NSRouter("/snapshots/:snapshotId", NewVolumeSnapshotPortal(), "get:GetVolumeSnapshot;put:UpdateVolumeSnapshot;delete:DeleteVolumeSnapshot"),

				// Creates, shows, lists, unpdates and deletes replication.
				beego.NSRouter("/replications", NewReplicationPortal(), "post:CreateReplication;get:ListReplications"),
				beego.NSRouter("/replications/detail", NewReplicationPortal(), "get:ListReplicationsDetail"),
				beego.NSRouter("/replications/:replicationId", NewReplicationPortal(), "get:GetReplication;put:UpdateReplication;delete:DeleteReplication"),
				beego.NSRouter("/replications/:replicationId/enable", NewReplicationPortal(), "post:EnableReplication"),
				beego.NSRouter("/replications/:replicationId/disable", NewReplicationPortal(), "post:DisableReplication"),
				beego.NSRouter("/replications/:replicationId/failover", NewReplicationPortal(), "post:FailoverReplication"),
				// Volume group contains a list of volumes that are used in the same application.
				beego.NSRouter("/volumeGroups", NewVolumeGroupPortal(), "post:CreateVolumeGroup;get:ListVolumeGroups"),
				beego.NSRouter("/volumeGroups/:groupId", NewVolumeGroupPortal(), "put:UpdateVolumeGroup;get:GetVolumeGroup;delete:DeleteVolumeGroup"),
			),
		)
	pattern := fmt.Sprintf("/%s/*", constants.APIVersion)
	beego.InsertFilter(pattern, beego.BeforeExec, context.Factory())
	beego.InsertFilter(pattern, beego.BeforeExec, auth.Factory())
	beego.InsertFilter("*", beego.BeforeExec, accesslog.Factory())
	beego.AddNamespace(ns)

	// add router for api version
	beego.Router("/", &VersionPortal{}, "get:ListVersions")
	beego.Router("/:apiVersion", &VersionPortal{}, "get:GetVersion")

	if apiServerCfg.HTTPSEnabled {
		if apiServerCfg.BeegoHTTPSCertFile == "" || apiServerCfg.BeegoHTTPSKeyFile == "" {
			fmt.Println("If https is enabled in hotpot, please ensure key file and cert file of the hotpot are not empty.")
			return
		}

		// beego https config
		beego.BConfig.Listen.EnableHTTP = false
		beego.BConfig.Listen.EnableHTTPS = true
		strs := strings.Split(apiServerCfg.ApiEndpoint, ":")
		beego.BConfig.Listen.HTTPSAddr = strs[AddressIdx]
		beego.BConfig.Listen.HTTPSPort, _ = strconv.Atoi(strs[PortIdx])
		beego.BConfig.Listen.HTTPSCertFile = apiServerCfg.BeegoHTTPSCertFile
		beego.BConfig.Listen.HTTPSKeyFile = apiServerCfg.BeegoHTTPSKeyFile
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			},
		}

		beego.BeeApp.Server.TLSConfig = tlsConfig
	}

	beego.BConfig.Listen.ServerTimeOut = constants.BeegoServerTimeOut
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.EnableErrorsShow = false
	beego.BConfig.EnableErrorsRender = false
	beego.BConfig.WebConfig.AutoRender = false

	// start service
	beego.Run(apiServerCfg.ApiEndpoint)
}
