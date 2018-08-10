// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package constants

const (
	// It's RFC 8601 format that decodes and encodes with
	// exactly precision to seconds.
	TimeFormat = `2006-01-02T15:04:05`

	DefaultOpensdsEndpoint = "http://localhost:50040"

	// This is set for None Auth
	DefaultTenantId = "e93b4c0934da416eb9c8d120c5d04d96"

	// Token parameter name
	AuthTokenHeader    = "X-Auth-Token"
	SubjectTokenHeader = "X-Subject-Token"

	// OpenSDS current api version
	APIVersion = "v1beta"
)
