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
	"encoding/json"
	"log"
	"reflect"

	"github.com/opensds/opensds/pkg/api"

	"github.com/astaxie/beego/context"
)

var fakeVersion api.VersionInfo
var fakeVersions api.AvailableVersions

func GetAllVersions(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	versions, err := api.ListVersions()
	if err != nil {
		log.Println(err)
		rbody, _ := json.Marshal("List versions failed!")
		ctx.Output.Body(rbody)
	} else {
		if reflect.DeepEqual(versions, fakeVersions) {
			log.Println("List versions failed!")
			rbody, _ := json.Marshal("List versions failed!")
			ctx.Output.Body(rbody)
		} else {
			rbody, _ := json.Marshal(versions)
			ctx.Output.Body(rbody)
		}
	}
}

func GetVersionv1(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.Output.ContentType("application/json")

	version, err := api.GetVersionv1()
	if err != nil {
		log.Println(err)
		rbody, _ := json.Marshal("Get version v1 failed!")
		ctx.Output.Body(rbody)
	} else {
		if reflect.DeepEqual(version, fakeVersion) {
			log.Println("Get version v1 failed!")
			rbody, _ := json.Marshal("Get version v1 failed!")
			ctx.Output.Body(rbody)
		} else {
			rbody, _ := json.Marshal(version)
			ctx.Output.Body(rbody)
		}
	}
}
