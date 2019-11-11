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
)

func init() {
	client = c.NewFakeClient(&c.Config{Endpoint: c.TestEp})
}

func TestHostAction(t *testing.T) {
	beCrasher := os.Getenv("BE_CRASHER")

	if beCrasher == "1" {
		var args []string
		hostAction(hostCommand, args)

		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestHostAction")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	e, ok := err.(*exec.ExitError)

	if ok && ("exit status 1" == e.Error()) {
		return
	}

	t.Fatalf("process ran with %s, want exit status 1", e.Error())
}

func TestHostCreateAction(t *testing.T) {
	var args []string
	args = append(args, "sap1")
	hostCreateAction(hostCreateCommand, args)
}

func TestHostShowAction(t *testing.T) {
	var args []string
	args = append(args, "202964b5-8e73-46fd-b41b-a8e403f3c30b")
	hostShowAction(hostShowCommand, args)
}

func TestHostListAction(t *testing.T) {
	var args []string
	hostListAction(hostListCommand, args)
}

func TestHostDeleteAction(t *testing.T) {
	var args []string
	args = append(args, "202964b5-8e73-46fd-b41b-a8e403f3c30b")
	hostDeleteAction(hostDeleteCommand, args)
}

func TestHostUpdateAction(t *testing.T) {
	var args []string
	args = append(args, "202964b5-8e73-46fd-b41b-a8e403f3c30b")
	hostUpdateAction(hostDeleteCommand, args)
}

func TestHostAddInitiatorAction(t *testing.T) {
	var args []string
	args = append(args, "202964b5-8e73-46fd-b41b-a8e403f3c30b", "20000024ff5bb888", "iSCSI")
	hostAddInitiatorAction(hostAddInitiatorCommand, args)
}

func TestHostRemoveInitiatorAction(t *testing.T) {
	var args []string
	args = append(args, "202964b5-8e73-46fd-b41b-a8e403f3c30b", "20000024ff5bb888")
	hostRemoveInitiatorAction(hostAddInitiatorCommand, args)
}
