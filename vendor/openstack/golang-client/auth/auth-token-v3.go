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

package auth

import (
	"errors"
	"time"
)

type AuthV3Token struct {
	Token struct {
		ID      string    `json:"-"`
		Methods []string  `json:"methods"`
		Expires time.Time `json:"expires"`
		Project struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
		Catalog []ServiceCatalog `json:"catalog"`
	} `json:"token"`
}

func (s AuthV3Token) GetToken() string {
	return s.Token.ID
}

func (s AuthV3Token) GetExpiration() time.Time {
	// The expiration functuion is not supported right now.
	return time.Now()
}

func (s AuthV3Token) GetEndpoint(serviceType, regionName string) (string, error) {
	for _, catalog := range s.Token.Catalog {
		if catalog.Type == serviceType {
			ept, err := catalog.GetEndpoint("public", regionName)
			if err != nil {
				return "", err
			} else {
				return ept, nil
			}
		}
	}
	err := errors.New("Endpoint not found for service type: " + serviceType)
	return "", err
}

func (s AuthV3Token) GetProject() string {
	return s.Token.Project.Name
}
