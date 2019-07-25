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

package drbd

import (
	"path/filepath"

	"github.com/LINBIT/godrbdutils"
)

// Probably make some of these configurable later
const (
	portPostfix  = "-drbd-port"
	minorPostfix = "-drbd-minor"

	defaultPortMin = 7000
	defaultPortMax = 8000

	minorMin = 1
	minorMax = 1000

	// for the time being opensds only has one primary and one secondary,
	// but reserve slots for 7 peers anyways
	maxPeers = 7

	resDir          = "/etc/drbd.d"
	defaultConfPath = "/etc/opensds/driver/drbd.yaml"
)

type drbdConf struct {
	Hosts   []godrbdutils.Host `yaml:"Hosts,omitempty"`
	PortMin int                `yaml:"PortMin,omitempty"`
	PortMax int                `yaml:"PortMax,omitempty"`
}

func portKey(s string) string  { return s + portPostfix }
func minorKey(s string) string { return s + minorPostfix }

func resFilePath(resName string) string {
	return filepath.Join(resDir, resName) + ".res"
}

func cfgOrDefault(c, d int) int {
	if c > 0 {
		return c
	}
	return d
}
