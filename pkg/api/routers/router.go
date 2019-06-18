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
This module implements a entry into the OpenSDS northbound REST service.

*/

package routers

import (
	"github.com/astaxie/beego"
	bctx "github.com/astaxie/beego/context"
	"github.com/opensds/opensds/pkg/api/controllers"
	"github.com/opensds/opensds/pkg/utils/constants"
)

func init() {

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
			beego.NSRouter("/:tenantId/docks", &controllers.DockPortal{}, "get:ListDocks"),
			// Show one dock service, including endpoint, driverName and so on
			beego.NSRouter("/:tenantId/docks/:dockId", &controllers.DockPortal{}, "get:GetDock"),

			// Profile is a set of policies configured by admin and provided for users
			// CreateProfile, UpdateProfile and DeleteProfile are used for admin only
			// ListProfiles and GetProfile are used for both admin and users
			beego.NSRouter("/:tenantId/profiles", &controllers.ProfilePortal{}, "post:CreateProfile;get:ListProfiles"),
			beego.NSRouter("/:tenantId/profiles/:profileId", &controllers.ProfilePortal{}, "get:GetProfile;put:UpdateProfile;delete:DeleteProfile"),

			// All operations of customProperties are used for admin only
			beego.NSRouter("/:tenantId/profiles/:profileId/customProperties", &controllers.ProfilePortal{}, "post:AddCustomProperty;get:ListCustomProperties"),
			beego.NSRouter("/:tenantId/profiles/:profileId/customProperties/:customKey", &controllers.ProfilePortal{}, "delete:RemoveCustomProperty"),

			// Pool is the virtual description of backend storage, usually divided into block, file and object,
			// and every pool is atomic, which means every pool contains a specific set of features.
			// ListPools and GetPool are used for checking the status of backend pool, admin only
			beego.NSRouter("/:tenantId/pools", &controllers.PoolPortal{}, "get:ListPools"),
			beego.NSRouter("/:tenantId/pools/:poolId", &controllers.PoolPortal{}, "get:GetPool"),
			beego.NSRouter("/:tenantId/availabilityZones", &controllers.PoolPortal{}, "get:ListAvailabilityZones"),
		)
	beego.AddNamespace(ns)

	// add router for api version
	beego.Router("/", &controllers.VersionPortal{}, "get:ListVersions")
	beego.Router("/:apiVersion", &controllers.VersionPortal{}, "get:GetVersion")
}
