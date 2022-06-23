// Copyright 2021 The SODA Foundation Authors.
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
	"github.com/sodafoundation/api/pkg/api/controllers"
	"github.com/sodafoundation/api/pkg/utils/constants"
)

func init() {

	// add router for aksk api
	akskns :=
		beego.NewNamespace("/"+constants.APIVersion+"/:tenantId/aksk",
			beego.NSRouter("/aksks", controllers.NewAkSkPortal(), "post:CreateAkSk;get:ListAkSks"),
			beego.NSRouter("/aksks/:UserId", controllers.NewAkSkPortal(), "get:GetAkSk;delete:DeleteAkSk"),
		)
	beego.AddNamespace(akskns)
}
