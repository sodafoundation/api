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
	"fmt"
	"os"

	"github.com/opensds/opensds/pkg/utils/constants"
	"github.com/opensds/opensds/pkg/utils/pwd"
)

const (
	//Opensds Auth ENVs
	OpensdsAuthStrategy = "OPENSDS_AUTH_STRATEGY"
	OpensdsTenantId     = "OPENSDS_TENANT_ID"

	// Keystone Auth ENVs
	OsAuthUrl       = "OS_AUTH_URL"
	OsUsername      = "OS_USERNAME"
	OsPassword      = "OS_PASSWORD"
	OsTenantName    = "OS_TENANT_NAME"
	OsProjectName   = "OS_PROJECT_NAME"
	OsUserDomainId  = "OS_USER_DOMAIN_ID"
	PwdEncrypter    = "PASSWORD_ENCRYPTER"
	EnableEncrypted = "ENABLE_ENCRYPTED"
	Keystone        = "keystone"
	Noauth          = "noauth"
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
	PwdEncrypter     string
	EnableEncrypted  bool
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

func LoadKeystoneAuthOptionsFromEnv() (*KeystoneAuthOptions, error) {
	// pre-check
	envs := []string{OsAuthUrl, OsUsername, OsPassword, OsTenantName, OsProjectName, OsUserDomainId}
	for _, env := range envs {
		if _, ok := os.LookupEnv(env); !ok {
			return nil, fmt.Errorf("can not get keystone ENV: %s", env)
		}
	}

	opt := NewKeystoneAuthOptions()
	opt.IdentityEndpoint = os.Getenv(OsAuthUrl)
	opt.Username = os.Getenv(OsUsername)
	var pwdCiphertext = os.Getenv(OsPassword)
	if os.Getenv(EnableEncrypted) == "T" {
		// Decrypte the password
		pwdTool := os.Getenv(PwdEncrypter)
		if pwdTool == "" {
			return nil, fmt.Errorf("The password encrypter can not be empty if password encrypted is enabled.")
		}

		password, err := pwd.NewPwdEncrypter(pwdTool).Decrypter(pwdCiphertext)
		if err != nil {
			return nil, fmt.Errorf("Decryption failed, %v", err)
		}
		pwdCiphertext = password
	}

	opt.Password = pwdCiphertext
	opt.TenantName = os.Getenv(OsTenantName)
	projectName := os.Getenv(OsProjectName)
	opt.DomainID = os.Getenv(OsUserDomainId)
	if opt.TenantName == "" {
		opt.TenantName = projectName
	}

	return opt, nil
}

func LoadNoAuthOptionsFromEnv() *NoAuthOptions {
	tenantId, ok := os.LookupEnv(OpensdsTenantId)
	if !ok {
		tenantId = constants.DefaultTenantId
	}
	return NewNoauthOptions(tenantId)
}
