// Copyright 2018 The OpenSDS Authors.
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

func TestReplicationAction(t *testing.T) {
	beCrasher := os.Getenv("BE_CRASHER")

	if beCrasher == "1" {
		volumeAttachmentAction(volumeAttachmentCommand, []string{})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestReplicationAction")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	e, ok := err.(*exec.ExitError)

	if ok && ("exit status 1" == e.Error()) {
		return
	}

	t.Fatalf("process ran with %s, want exit status 1", e.Error())
}

func TestReplicationCreateAction(t *testing.T) {
	var args = []string{
		"3fc90eda-4ef6-410d-b1b9-f39c6476683d",
		"e0bfd484-0a95-429a-9065-2a797f673d0d",
	}
	replicationCreateAction(replicationCreateCommand, args)
}

func TestReplicationShowAction(t *testing.T) {
	var args = []string{"f2dda3d2-bf79-11e7-8665-f750b088f63e"}
	replicationShowAction(replicationShowCommand, args)
}

func TestReplicationListAction(t *testing.T) {
	var args []string
	replicationListAction(replicationListCommand, args)
}

func TestReplicationDeleteAction(t *testing.T) {
	var args = []string{"f2dda3d2-bf79-11e7-8665-f750b088f63e"}
	replicationDeleteAction(replicationDeleteCommand, args)
}

func TestReplicationUpdateAction(t *testing.T) {
	var args = []string{"f2dda3d2-bf79-11e7-8665-f750b088f63e"}
	replicationUpdateAction(replicationUpdateCommand, args)
}

func TestReplicationEnableAction(t *testing.T) {
	var args = []string{"f2dda3d2-bf79-11e7-8665-f750b088f63e"}
	replicationEnableAction(replicationEnableCommand, args)
}
func TestReplicationDisableAction(t *testing.T) {
	var args = []string{"f2dda3d2-bf79-11e7-8665-f750b088f63e"}
	replicationDisableAction(replicationDisableCommand, args)
}
func TestReplicationFailoverAction(t *testing.T) {
	var args = []string{"f2dda3d2-bf79-11e7-8665-f750b088f63e"}
	replicationFailoverAction(replicationFailoverCommand, args)
}
