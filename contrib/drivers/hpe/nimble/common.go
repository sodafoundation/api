// Copyright 2019 The OpenSDS Authors.
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

package nimble

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

const UnitGi = 1024 * 1024 * 1024
const UnitMebi = 1024

func EncodeName(id string) string {
	h := md5.New()
	h.Write([]byte(id))
	encodedName := hex.EncodeToString(h.Sum(nil))
	prefix := strings.Split(id, "-")[0] + "-"
	postfix := encodedName[:MaxNameLength-len(prefix)]
	return prefix + postfix
}

func TruncateDescription(desc string) string {
	if len(desc) > MaxDescriptionLength {
		desc = desc[:MaxDescriptionLength]
	}
	return desc
}

func Byte2Gib(byteSize int64) int64 {
	return byteSize / UnitGi
}

func Gib2Mebi(gigabyte int64) int64 {
	return gigabyte * UnitMebi
}
