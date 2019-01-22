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

package client

import (
	"os"

	"github.com/opensds/opensds/pkg/utils/constants"
)

const (
	//Opensds Auth ENVs
	OpensdsAuthStrategy = "OPENSDS_AUTH_STRATEGY"
	OpensdsTenantId     = "OPENSDS_TENANT_ID"

	// Keystone Auth ENVs
	OsAuthUrl      = "OS_AUTH_URL"
	OsUsername     = "OS_USERNAME"
	OsPassword     = "OS_PASSWORD"
	OsTenantName   = "OS_TENANT_NAME"
	OsProjectName  = "OS_PROJECT_NAME"
	OsUserDomainId = "OS_USER_DOMAIN_ID"

	Keystone = "keystone"
	Noauth   = "noauth"
)

type AuthOptions interface {
	GetTenantId() string
}

func NewKeystoneAuthOptions() *KeystoneAuthOptions {
	return &KeystoneAuthOptions{}
}

type KeystoneAuthOptions struct {
	IdentityEndpoint string
	Username         string
	UserID           string
	Password         string
	DomainID         string
	DomainName       string
	TenantID         string
	TenantName       string
	AllowReauth      bool
	TokenID          string
}

func (k *KeystoneAuthOptions) GetTenantId() string {
	return k.TenantID
}

func NewNoauthOptions(tenantId string) *NoAuthOptions {
	return &NoAuthOptions{TenantID: tenantId}
}

type NoAuthOptions struct {
	TenantID string
}

func (n *NoAuthOptions) GetTenantId() string {
	return n.TenantID
}

func LoadKeystoneAuthOptionsFromEnv() *KeystoneAuthOptions {
	opt := NewKeystoneAuthOptions()
	opt.IdentityEndpoint = os.Getenv(OsAuthUrl)
	opt.Username = os.Getenv(OsUsername)
	opt.Password = os.Getenv(OsPassword)
	opt.TenantName = os.Getenv(OsTenantName)
	projectName := os.Getenv(OsProjectName)
	opt.DomainID = os.Getenv(OsUserDomainId)
	if opt.TenantName == "" {
		opt.TenantName = projectName
	}
	return opt
}

func LoadNoAuthOptionsFromEnv() *NoAuthOptions {
	tenantId, ok := os.LookupEnv(OpensdsTenantId)
	if !ok {
		tenantId = constants.DefaultTenantId
	}
	return NewNoauthOptions(tenantId)
}
