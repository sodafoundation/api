// auth-password - Username/Password Authentication
// Copyright 2015 Dean Troyer
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openstack

import (
	"encoding/json"
	"errors"
	// "strings"
)

// The token request structure for Identity v2

type PasswordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type OSAuth struct {
	PasswordCredentials `json:"passwordCredentials"`
	ProjectName         string `json:"tenantName"`
	ProjectId           string `json:"tenantId"`
}

type UserPassV2 struct {
	OSAuth  `json:"auth"`
	AuthUrl string `json:"-"`
}

func NewUserPassV2(ao AuthOpts) (upv2 *UserPassV2, err error) {
	// Validate incoming values
	if ao.AuthUrl == "" {
		err = errors.New("AuthUrl required")
		return nil, err
	}
	if ao.Username == "" {
		err = errors.New("Username required")
		return nil, err
	}
	if ao.Password == "" {
		err = errors.New("Password required")
		return nil, err
	}
	upv2 = &UserPassV2{
		AuthUrl: ao.AuthUrl,
		OSAuth: OSAuth{
			PasswordCredentials: PasswordCredentials{
				Username: ao.Username,
				Password: ao.Password,
			},
			ProjectName: ao.ProjectName,
			ProjectId:   ao.ProjectId,
		},
	}
	return upv2, nil
}

// Produce JSON output
func (s *UserPassV2) JSON() []byte {
	reqAuth, err := json.Marshal(s)
	if err != nil {
		// Return an empty structure
		reqAuth = []byte{'{', '}'}
	}
	return reqAuth
}

// func (self *UserPassV2) AuthUserPassV2(opts interface{}) (AuthRef, error) {
//     auth, err := self.GetAuthRef()
//     return AuthRef(auth), err
// }
