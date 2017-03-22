// Copyright (c) 2016 eBay Inc.
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

package middleware

import (
	"time"

	"git.openstack.org/openstack/golang-client/openstack"
)

type Validator struct {
	// Service account to talk to keystone
	SvcAuthOpts openstack.AuthOpts
	// File path the signing cert would be stored/cached
	CachedSigningKeyPath string
	TokenId              string
	// Token revocation list memory cache duration
	RevCacheDuration time.Duration
}

// Token revocation response structure
type revokeResp struct {
	Signed string `json:"signed"`
}

type revokedListCache struct {
	Revoked []openstack.Token
	// time when this cache was built
	Time time.Time
}

type RevokedList struct {
	Revoked []openstack.Token
}
