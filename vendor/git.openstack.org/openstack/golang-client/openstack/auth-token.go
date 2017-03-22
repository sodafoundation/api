// auth-token - Token Authentication
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
	"errors"
	"time"
)

// Identity Response Types

type AccessType struct {
	Token          Token                 `json:"token"`
	User           interface{}           `json:"user"`
	ServiceCatalog []ServiceCatalogEntry `json:"servicecatalog"`
}

type AuthToken struct {
	Access AccessType `json:"access"`
}

type Token struct {
	ID      string    `json:"id"`
	Expires time.Time `json:"expires"`
	Project struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"tenant"`
}

func (s AuthToken) GetToken() string {
	return s.Access.Token.ID
}

func (s AuthToken) GetExpiration() time.Time {
	return s.Access.Token.Expires
}

func (s AuthToken) GetEndpoint(serviceType string, regionName string) (string, error) {

	// Parse service catalog
	for _, v := range s.Access.ServiceCatalog {
		ep, err := v.GetEndpoint(serviceType, "public", regionName)
		if err == nil {
			return ep, nil
		}
	}
	return "", errors.New("err: endpoint not found")
}

func (s AuthToken) GetProject() string {
	return s.Access.Token.Project.Name
}
