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

package v3

import (
	"errors"
	"fmt"
	"strings"
)

type ServiceEndpoint struct {
	URL       string `json:"url"`
	Interface string `json:"interface"`
	Region    string `json:"region"`
	RegionID  string `json:"region_id"`
	ID        string `json:"id"`
}

type ServiceCatalog struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Endpoints []ServiceEndpoint `json:"endpoints"`
}

// Valid interfaceType values: 'public', 'publicURL', 'admin', 'admin URL', 'internal', 'internalURL'
func (sc ServiceCatalog) GetEndpoint(interfaceType, regionName string) (string, error) {
	if interfaceType == "" {
		// Set the default value
		interfaceType = "public"
	}

	for _, epoint := range sc.Endpoints {
		if regionName == "" || epoint.Region == regionName {
			// Translate passed interface types
			sc_int := strings.ToLower(interfaceType)

			if epoint.Interface == sc_int {
				return epoint.URL, nil
			}
		}
	}

	var msg string
	if regionName != "" {
		msg = fmt.Sprintf("%s endpoint in %s region not found",
			interfaceType, regionName)
	} else {
		msg = fmt.Sprintf("%s endpoint not found", interfaceType)
	}
	return "", errors.New(msg)
}
