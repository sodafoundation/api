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

package cli

import (
	"os"
	"os/exec"
	"testing"

	c "github.com/opensds/opensds/client"
)

func init() {
	client = c.NewFakeClient(&c.Config{Endpoint: c.TestEp})
}

func TestVolumeAction(t *testing.T) {
	beCrasher := os.Getenv("BE_CRASHER")

	if beCrasher == "1" {
		var args []string
		volumeAction(volumeCommand, args)

		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestVolumeAction")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	e, ok := err.(*exec.ExitError)

	if ok && ("exit status 1" == e.Error()) {
		return
	}

	t.Fatalf("process ran with %s, want exit status 1", e.Error())
}

func TestVolumeCreateAction(t *testing.T) {
	var args []string
	args = append(args, "1")
	volumeCreateAction(volumeCreateCommand, args)
}

func TestVolumeShowAction(t *testing.T) {
	var args []string
	args = append(args, "bd5b12a8-a101-11e7-941e-d77981b584d8")
	volumeShowAction(volumeShowCommand, args)
}

func TestVolumeListAction(t *testing.T) {
	var args []string
	volumeListAction(volumeListCommand, args)
}

func TestVolumeDeleteAction(t *testing.T) {
	var args []string
	args = append(args, "bd5b12a8-a101-11e7-941e-d77981b584d8")
	volumeDeleteAction(volumeDeleteCommand, args)
}

func TestVolumeUpdateAction(t *testing.T) {
	var args []string
	args = append(args, "bd5b12a8-a101-11e7-941e-d77981b584d8")
	volumeUpdateAction(volumeDeleteCommand, args)
}
func TestVolumeExtendAction(t *testing.T) {
	var args []string
	args = append(args, "bd5b12a8-a101-11e7-941e-d77981b584d8")
	args = append(args, "5")
	volumeExtendAction(volumeExtendCommand, args)
}
