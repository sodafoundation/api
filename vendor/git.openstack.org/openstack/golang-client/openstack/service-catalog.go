// service-catalog - Service Catalog structs
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
	"fmt"
	"strings"
)

type ServiceCatalogEntry struct {
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Endpoints []ServiceEndpoint `json:"endpoints"`
	// Endpoints []map[string]string `json:"endpoints"`
}

type ServiceEndpoint struct {
	Type        string `json:"type"`
	Region      string `json:"region"`
	PublicURL   string `json:"publicurl"`
	AdminURL    string `json:"adminurl"`
	InternalURL string `json:"internalurl"`
	VersionID   string `json:"versionid"`
}

// Valid interfaceType values: 'public', 'publicURL', 'admin', 'admin URL', 'internal', 'internalURL'
func (sce ServiceCatalogEntry) GetEndpoint(
	serviceType string,
	interfaceType string,
	regionName string,
) (string, error) {
	if interfaceType == "" {
		// Set the default value
		interfaceType = "public"
	}
	if sce.Type == serviceType {
		for _, r := range sce.Endpoints {
			if regionName == "" || r.Region == regionName {
				// Translate passed interface types
				sc_int := strings.ToLower(interfaceType)
				if sc_int == "public" || sc_int == "publicurl" {
					return r.PublicURL, nil
				}
				if sc_int == "admin" || sc_int == "adminurl" {
					return r.AdminURL, nil
				}
				if sc_int == "internal" || sc_int == "internalURL" {
					return r.InternalURL, nil
				}
			}
		}
	}

	var msg string
	if regionName != "" {
		msg = fmt.Sprintf("%s endpoint for %s service in %s region not found",
			interfaceType, serviceType, regionName)
	} else {
		msg = fmt.Sprintf("%s endpoint for %s service not found",
			interfaceType, serviceType)
	}
	return "", errors.New(msg)
}
