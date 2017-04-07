// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

/*
This module implements a entry into the OpenSDS northbound REST service.

*/

package northbound

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func Run(host string) {
	ns :=
		beego.NewNamespace("/api",
			beego.NSCond(func(ctx *context.Context) bool {
				if ctx.Input.Scheme() == "http" {
					return true
				}
				return false
			}),
			beego.NSGet("/", GetAllVersions),
			beego.NSNamespace("/v1",
				beego.NSGet("/", GetVersionv1),
				beego.NSNamespace("/volumes",
					beego.NSRouter("/:resource", &VolumeController{}),
					beego.NSGet("/:resource/:id", GetVolume),
					beego.NSPut("/attach", AttachVolume),
					beego.NSDelete("/attach", DetachVolume),
					beego.NSPut("/mount", MountVolume),
					beego.NSDelete("/mount", UnmountVolume),
				),
				beego.NSNamespace("/shares",
					beego.NSRouter("/:resource", &ShareController{}),
					beego.NSGet("/:resource/:id", GetShare),
					beego.NSPut("/attach", AttachShare),
					beego.NSDelete("/attach", DetachShare),
					beego.NSPut("/mount", MountShare),
					beego.NSDelete("/mount", UnmountShare),
				),
			),
		)

	beego.AddNamespace(ns)
	beego.Run(host)
}
