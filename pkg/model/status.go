// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

// volume status
const (
	VOLUME_CREATING        = "creating"
	VOLUME_AVAILABLE       = "available"
	VOLUME_IN_USE          = "inUse"
	VOLUME_DELETING        = "deleting"
	VOLUME_ERROR           = "error"
	VOLUEM_ERROR_DELETING  = "errorDeleting"
	VOLUME_ERROR_EXTENDING = "errorExtending"
	VOLUME_EXTENDING       = "extending"
)

// volume attach status
const (
	VOLUME_ATTACHING       = "attaching"
	VOLUME_ATTACHED        = "attached"
	VOLUME_DETACHED        = "detached"
	VOLUME_RESERVED        = "reserved"
	VOLUME_ERROR_ATTACHING = "errorAttaching"
	VOLUME_ERROR_DETACHING = "errorDetaching"
	VOLUME_DELETED         = "deleted"
)

// volume snapshot status
const (
	VOLUMESNAP_CREATING       = "creating"
	VOLUMESNAP_AVAILABLE      = "available"
	VOLUMESNAP_DELETING       = "deleting"
	VOLUMESNAP_ERROR          = "error"
	VOLUMESNAP_ERROR_DELETING = "errorDeleting"
	VOLUMESNAP_DELETED        = "deleted"
)

// volume attachment status
const (
	VOLUMEATM_CREATING       = "creating"
	VOLUMEATM_AVAILABLE      = "available"
	VOLUMEATM_ERROR_DELETING = "errorDeleting"
	VOLUMEATM_ERROR          = "error"
)

// volume group status
const (
	VOLUMEGROUP_CREATING       = "creating"
	VOLUMEGROUP_AVAILABLE      = "available"
	VOLUMEGROUP_ERROR_DELETING = "errorDeleting"
	VOLUMEGROUP_ERROR          = "error"
	VOLUMEGROUP_DELETING       = "deleting"
	VOLUMEGROUP_DELETED        = "deleted"
	VOLUMEGROUP_UPDATING       = "updating"
	VOLUMEGROUP_IN_USE         = "inUse"
)
