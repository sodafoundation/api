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

func TestFileShareSnapshotAction(t *testing.T) {
	beCrasher := os.Getenv("BE_CRASHER")

	if beCrasher == "1" {
		var args []string
		fileShareSnapshotAction(fileShareSnapshotCommand, args)

		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFileShareSnapshotAction")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	e, ok := err.(*exec.ExitError)

	if ok && ("exit status 1" == e.Error()) {
		return
	}

	t.Fatalf("process ran with %s, want exit status 1", e.Error())
}

func TestFileShareSnapshotCreateAction(t *testing.T) {
	var args []string
	args = append(args, "bd5b12a8-a101-11e7-941e-d77981b584d8")
	fileShareSnapshotCreateAction(fileShareSnapshotCreateCommand, args)
}

func TestFileShareSnapshotShowAction(t *testing.T) {
	var args []string
	args = append(args, "3769855c-a102-11e7-b772-17b880d2f537")
	fileShareSnapshotShowAction(fileShareSnapshotShowCommand, args)
}

func TestFileShareSnapshotListAction(t *testing.T) {
	var args []string
	fileShareSnapshotListAction(fileShareSnapshotListCommand, args)
}

func TestFileShareSnapshotDeleteAction(t *testing.T) {
	var args []string
	args = append(args, "3769855c-a102-11e7-b772-17b880d2f537")
	fileShareSnapshotDeleteAction(fileShareSnapshotDeleteCommand, args)
}

func TestFileShareSnapshotUpdateAction(t *testing.T) {
	var args []string
	args = append(args, "3769855c-a102-11e7-b772-17b880d2f537")
	fileShareSnapshotUpdateAction(fileShareSnapshotDeleteCommand, args)
}
