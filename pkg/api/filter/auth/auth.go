// Copyright 2018 The OpenSDS Authors.
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

package auth

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	log "github.com/golang/glog"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/utils/config"
	"github.com/opensds/opensds/pkg/utils/constants"
)

type AuthBase interface {
	Filter(ctx *context.Context)
}

func NewNoAuth() AuthBase {
	return &NoAuth{}
}

type NoAuth struct{}

func (auth *NoAuth) Filter(httpCtx *context.Context) {
	ctx := c.GetContext(httpCtx)
	ctx.TenantId = httpCtx.Input.Param(":tenantId")
	// In noauth case, only the default id is treated as admin role.
	if ctx.TenantId == constants.DefaultTenantId {
		ctx.IsAdmin = true
	}
	httpCtx.Input.SetData("context", ctx)
}

func Factory() beego.FilterFunc {
	var auth AuthBase
	log.Infof(config.CONF.AuthStrategy)
	switch config.CONF.AuthStrategy {
	case "keystone":
		auth = NewKeystone()
	case "noauth":
		auth = NewNoAuth()
	default:
		auth = NewNoAuth()
	}
	log.Info(auth)
	return auth.Filter
}
