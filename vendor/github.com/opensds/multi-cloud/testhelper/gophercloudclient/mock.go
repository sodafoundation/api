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

package gophercloudclient

import (
	"github.com/gophercloud/gophercloud"
	"github.com/opensds/multi-cloud/testhelper"
)

// Fake token to use.
const TokenID = "cbc36478b0bd8e67e89469c7749d4127"

// ServiceClient returns a generic service client for use in tests.
func ServiceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{TokenID: TokenID},
		Endpoint:       testhelper.Endpoint(),
	}
}
