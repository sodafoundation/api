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

package context

import (
	"encoding/json"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-log"
)

const (
	DefaultTenantId     = "tenantId"
	DefaultUserId       = "userId"
	NoAuthAdminTenantId = "adminTenantId"
)

const (
	KContext = "context"
)

type Context struct {
	TenantId string `json:"tenantId"`
	IsAdmin  bool   `json:"isAdmin"`
	UserId   string `json:"userId"`
}

func NewAdminContext() *Context {
	return &Context{
		TenantId: NoAuthAdminTenantId,
		IsAdmin:  true,
		UserId:   "unkown",
	}
}

func NewContext() *Context {
	return &Context{
		TenantId: DefaultTenantId,
		IsAdmin:  true,
		UserId:   DefaultUserId,
	}
}

func NewContextFromJson(s string) *Context {
	ctx := &Context{}
	err := json.Unmarshal([]byte(s), ctx)
	if err != nil {
		log.Logf("Unmarshal json to context failed, reason: %v", err)
	}
	return ctx
}

func (ctx *Context) ToJson() string {
	b, err := json.Marshal(ctx)
	if err != nil {
		log.Logf("Context convert to json failed, reason: %v", err)
	}
	return string(b)
}

func FilterFactory() restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		req.SetAttribute(KContext, NewContext())
		chain.ProcessFilter(req, resp)
	}
}
