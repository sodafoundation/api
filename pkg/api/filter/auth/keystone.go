// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

// Keystone authentication middleware, only support keystone v3.

package auth

import (
	"net/http"
	"strings"
	"time"

	bctx "github.com/astaxie/beego/context"
	log "github.com/golang/glog"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/config"
)

const (
	authTokenHeader    = "X-Auth-Token"
	subjectTokenHeader = "X-Subject-Token"
)

func NewKeystone() AuthBase {
	k := &Keystone{}
	if err := k.SetUp(); err != nil {
		// If auth set up failed, raise panic.
		panic(err)
	}
	return k
}

type Keystone struct {
	identity *gophercloud.ServiceClient
	ctx      *bctx.Context
}

func (k *Keystone) SetUp() error {
	c := config.CONF.KeystoneAuthToken
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: c.AuthUrl,
		DomainName:       c.UserDomainName,
		Username:         c.Username,
		Password:         c.Password,
		TenantName:       c.ProjectName,
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Error("When get auth client:", err)
		return err
	}
	// Only support keystone v3
	k.identity, err = openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		log.Error("When get identity session:", err)
		return err
	}
	return nil
}

func (k *Keystone) setPolicyContext(r tokens.GetResult) error {
	roles, err := r.ExtractRoles()
	if err != nil {
		return model.HttpError(k.ctx, http.StatusUnauthorized, "Extract roles failed,%v", err)
	}

	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	project, err := r.ExtractProject()
	if err != nil {
		return model.HttpError(k.ctx, http.StatusUnauthorized, "Extract project failed,%v", err)
	}

	user, err := r.ExtractUser()
	if err != nil {
		return model.HttpError(k.ctx, http.StatusUnauthorized, "Extract user failed,%v", err)
	}

	policyCtx := &context.Context{
		ProjectId:      project.ID,
		Roles:          roleNames,
		UserId:         user.ID,
		IsAdminProject: project.Name == "admin",
	}
	k.ctx.Input.SetData("context", policyCtx)
	return nil
}

func (k *Keystone) validateToken(token string) error {
	r := tokens.Get(k.identity, token)
	if r.Err != nil {
		return model.HttpError(k.ctx, http.StatusUnauthorized, "Get token failed,%v", r.Err)
	}

	t, err := r.ExtractToken()
	if err != nil {
		return model.HttpError(k.ctx, http.StatusUnauthorized, "Extract token failed,%v", err)

	}
	log.V(8).Infof("token: %v", t)

	if time.Now().After(t.ExpiresAt) {
		return model.HttpError(k.ctx, http.StatusUnauthorized,
			"Token has expired, expire time %v", t.ExpiresAt)
	}
	return k.setPolicyContext(r)
}

func (k *Keystone) Filter(ctx *bctx.Context) {
	// Strip the spaces around the token
	token := strings.TrimSpace(ctx.Input.Header(authTokenHeader))
	k.ctx = ctx
	k.validateToken(token)
}
