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

package urls

import (
	"strings"

	"github.com/opensds/opensds/pkg/utils/constants"
)

const (
	Etcd   = iota // Etcd == 0
	Client        // Client == 1
)

func GenerateDockURL(urlType int, tenantId string, in ...string) string {
	return generateURL("docks", urlType, tenantId, in...)
}

func GeneratePoolURL(urlType int, tenantId string, in ...string) string {
	return generateURL("pools", urlType, tenantId, in...)
}

func GenerateProfileURL(urlType int, tenantId string, in ...string) string {
	return generateURL("profiles", urlType, tenantId, in...)
}

func GenerateVolumeURL(urlType int, tenantId string, in ...string) string {
	return generateURL("block/volumes", urlType, tenantId, in...)
}

// GenerateNewVolumeURL ...
func GenerateNewVolumeURL(urlType int, tenantId string, in ...string) string {
	return generateURL("volumes", urlType, tenantId, in...)
}

func GenerateAttachmentURL(urlType int, tenantId string, in ...string) string {
	return generateURL("block/attachments", urlType, tenantId, in...)
}

func GenerateSnapshotURL(urlType int, tenantId string, in ...string) string {
	return generateURL("block/snapshots", urlType, tenantId, in...)
}

func GenerateReplicationURL(urlType int, tenantId string, in ...string) string {
	return generateURL("block/replications", urlType, tenantId, in...)
}

func GenerateVolumeGroupURL(urlType int, tenantId string, in ...string) string {
	return generateURL("block/volumeGroup", urlType, tenantId, in...)
}
func generateURL(resource string, urlType int, tenantId string, in ...string) string {
	// If project id is not specified, ignore it.
	if tenantId == "" {
		value := []string{CurrentVersion(), resource}
		value = append(value, in...)
		return strings.Join(value, "/")
	}

	// Set project id after resource url just for etcd query performance.
	var value []string
	if urlType == Etcd {
		value = []string{CurrentVersion(), resource, tenantId}
	} else {
		value = []string{CurrentVersion(), tenantId, resource}
	}
	value = append(value, in...)

	return strings.Join(value, "/")
}

func CurrentVersion() string {
	return constants.APIVersion
}
