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

package cli

import (
	"os"
	"os/exec"
	"testing"

	c "github.com/opensds/opensds/client"
	. "github.com/opensds/opensds/testutils/collection"
)

func init() {
	client = c.NewFakeClient(&c.Config{Endpoint: c.TestEp})
}

func TestZoneAction(t *testing.T) {
	beCrasher := os.Getenv("BE_CRASHER")

	if beCrasher == "1" {
		var args []string
		zoneAction(zoneCommand, args)

		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestZoneAction")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	e, ok := err.(*exec.ExitError)

	if ok && ("exit status 1" == e.Error()) {
		return
	}

	t.Fatalf("process ran with %s, want exit status 1", e.Error())
}

func TestZoneCreateAction(t *testing.T) {
	var args []string
	args = append(args, ByteAvailabilityZone)
	zoneCreateAction(zoneCreateCommand, args)
}

func TestZoneShowAction(t *testing.T) {
	var args []string
	args = append(args, "1106b972-66ef-11e7-b172-db03f3689c9c")
	zoneShowAction(zoneShowCommand, args)
}

func TestZoneListAction(t *testing.T) {
	var args []string
	zoneListAction(zoneListCommand, args)
}

func TestZoneDeleteAction(t *testing.T) {
	var args []string
	args = append(args, "1106b972-66ef-11e7-b172-db03f3689c9c")
	zoneDeleteAction(zoneDeleteCommand, args)
}
