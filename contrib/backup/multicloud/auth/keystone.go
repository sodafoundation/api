// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Keystone authentication middleware, only support keystone v3.
package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/emicklei/go-restful"
	log "github.com/golang/glog"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/constants"
)

const (
	ConfFile             = "/etc/opensds/driver/multi-cloud.yaml"
	DefaultUploadTimeout = 30 // in Seconds
)

type Keystone struct {
	identity *gophercloud.ServiceClient
	conf     *MultiCloudConf
}

type MultiCloudConf struct {
	Endpoint      string `yaml:"Endpoint,omitempty"`
	UploadTimeout int64  `yaml:"UploadTimeout,omitempty"`
	AuthOptions   `yaml:"AuthOptions,omitempty"`
}

type AuthOptions struct {
	Strategy        string `yaml:"Strategy"`
	AuthUrl         string `yaml:"AuthUrl,omitempty"`
	DomainName      string `yaml:"DomainName,omitempty"`
	UserName        string `yaml:"UserName,omitempty"`
	Password        string `yaml:"Password,omitempty"`
	PwdEncrypter    string `yaml:"PwdEncrypter,omitempty"`
	EnableEncrypted bool   `yaml:"EnableEncrypted,omitempty"`
	TenantName      string `yaml:"TenantName,omitempty"`
}

func GetIdentity(k *Keystone) *gophercloud.ServiceClient {
	return k.identity
}

func (k *Keystone) loadConf(p string) (*MultiCloudConf, error) {
	conf := &MultiCloudConf{
		Endpoint:      "http://127.0.0.1:8088",
		UploadTimeout: DefaultUploadTimeout,
	}
	confYaml, err := ioutil.ReadFile(p)
	if err != nil {
		log.Errorf("Read config yaml file (%s) failed, reason:(%v)", p, err)
		return nil, err
	}
	if err = yaml.Unmarshal(confYaml, conf); err != nil {
		log.Errorf("Parse error: %v", err)
		return nil, err
	}
	return conf, nil
}

func (k *Keystone) SetUp() error {
	var err error
	if k.conf, err = k.loadConf(ConfFile); err != nil {
		return err
	}

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: k.conf.AuthUrl,
		DomainName:       k.conf.DomainName,
		Username:         k.conf.UserName,
		Password:         k.conf.Password,
		TenantName:       k.conf.TenantName,
	}
	bytes, _ := json.Marshal(opts)
	log.Infof("bytes:%v", string(bytes))
	log.Infof("opts:%+v", opts)
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
	log.V(4).Infof("Service Token Info: %s", provider.TokenID)
	return nil
}

func (k *Keystone) Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// Strip the spaces around the token  ctx.Input.Header(constants.AuthTokenHeader)
	token := strings.TrimSpace(req.HeaderParameter(constants.AuthTokenHeader))
	if err := k.validateToken(req, resp, token); err != nil {
		return
	}
	chain.ProcessFilter(req, resp)
}

func (k *Keystone) validateToken(req *restful.Request, res *restful.Response, token string) error {
	if token == "" {
		return model.HttpError(res, http.StatusUnauthorized, "token not found in header")
	}
	var r tokens.GetResult
	// The service token may be expired or revoked, so retry to get new token.
	err := utils.Retry(2, "verify token", false, func(retryIdx int, lastErr error) error {
		if retryIdx > 0 {
			// Fixme: Is there any better method ?
			if lastErr.Error() == "Authentication failed" {
				k.SetUp()
			} else {
				return lastErr
			}
		}
		log.Info("k.identity:", k.identity)
		r = tokens.Get(k.identity, token)
		log.Info("r:", r)
		log.Info("r.err:", r.Err)
		return r.Err
	})
	if err != nil {
		return model.HttpError(res, http.StatusUnauthorized, "get token failed,%v", r.Err)
	}

	t, err := r.ExtractToken()
	if err != nil {
		return model.HttpError(res, http.StatusUnauthorized, "extract token failed,%v", err)

	}
	log.V(8).Infof("token: %v", t)

	if time.Now().After(t.ExpiresAt) {
		return model.HttpError(res, http.StatusUnauthorized,
			"token has expired, expire time %v", t.ExpiresAt)
	}
	return k.setPolicyContext(req, res, r)
}

func (k *Keystone) setPolicyContext(req *restful.Request, res *restful.Response, r tokens.GetResult) error {
	roles, err := r.ExtractRoles()
	if err != nil {
		return model.HttpError(res, http.StatusUnauthorized, "extract roles failed,%v", err)
	}

	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	project, err := r.ExtractProject()
	if err != nil {
		return model.HttpError(res, http.StatusUnauthorized, "extract project failed,%v", err)
	}

	user, err := r.ExtractUser()
	if err != nil {
		return model.HttpError(res, http.StatusUnauthorized, "extract user failed,%v", err)
	}
	req.SetAttribute("Roles", project.ID)
	req.SetAttribute("TenantId", roleNames)
	req.SetAttribute("UserId", user.ID)
	req.SetAttribute("IsAdminProject", strings.ToLower(project.Name) == "admin")
	return nil
}
