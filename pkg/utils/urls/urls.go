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

package urls

import (
	"strings"
)

func GenerateDockURL(in ...string) string {
	return generateURL("docks", in...)
}

func GeneratePoolURL(in ...string) string {
	return generateURL("pools", in...)
}

func GenerateProfileURL(in ...string) string {
	return generateURL("profiles", in...)
}

func GenerateVolumeURL(in ...string) string {
	return generateURL("block/volumes", in...)
}

func GenerateAttachmentURL(in ...string) string {
	return generateURL("block/attachments", in...)
}

func GenerateSnapshotURL(in ...string) string {
	return generateURL("block/snapshots", in...)
}

func generateURL(resource string, in ...string) string {
	value := []string{CurrentVersion(), resource}
	value = append(value, in...)

	return strings.Join(value, "/")
}

func CurrentVersion() string {
	return "v1beta"
}
