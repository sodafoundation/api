// Copyright 2017 The OpenSDS Authors.
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

/*
This module implements the common data structure.

*/

package model

// An OpenSDS profile is identified by a unique name and ID. With adding
// extra properties, each profile can contains a set of tags of storage
// capabilities which are desirable features for a class of applications.
type ProfileSpec struct {
	*BaseModel
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	StorageType string    `json:"storageType,omitempty"`
	Extras      ExtraSpec `json:"extras,omitempty"`
}
