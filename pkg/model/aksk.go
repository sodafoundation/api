// Copyright 2021 The SODA Foundation Authors.
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

package model

import "encoding/json"

//AkSkSpec is used to hold the data to generate the AccessKey and SecretKey for the User.

type AkSkSpec struct {
	*BaseModel

	// ProjectId or TenantId is the tenant that the user belongs to.
	ProjectId string `json:"project_id,omitempty"`

	// The id of the user for whom the AkSk is being generated.
	UserId string `json:"user_id,omitempty"`

	// The json containing the accesskey and secretkey
	Blob string `json:"blob,omitempty"`

	// TODO: Confirm the usage of the Type field.
	//The type of backend ??
	Type string `json:"type,omitempty"`
}

// Blob to hold the access key and secret key.
type Blob struct {

	// Access key
	Access string `json:"accessKey,omitempty"`

	// Secret key
	Secret string `json:"secretKey ,omitempty"`
}

// TODO : Check if required or not.
// MarshalJSON to remove sensitive data
func (m AkSkSpec) MarshalJSON() ([]byte, error) {
	type akskResp AkSkSpec
	resp := akskResp(m)
	return json.Marshal(resp)
}

