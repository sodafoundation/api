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

func TestFileShareAclAction(t *testing.T) {
	beCrasher := os.Getenv("BE_CRASHER")

	if beCrasher == "1" {
		var args []string
		fileShareAclAction(fileShareAclCommand, args)

		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFileShareAclAction")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	e, ok := err.(*exec.ExitError)

	if ok && ("exit status 1" == e.Error()) {
		return
	}

	t.Fatalf("process ran with %s, want exit status 1", e.Error())
}

func TestFileShareAclCreateAction(t *testing.T) {
	var args []string
	args = append(args, "bd5b12a8-a101-11e7-941e-d77981b584d8")
	fileShareAclCreateAction(fileShareAclCreateCommand, args)
}

func TestFileShareAclDeleteAction(t *testing.T) {
	var args []string
	args = append(args, "d2975ebe-d82c-430f-b28e-f373746a71ca")
	fileShareAclDeleteAction(fileShareAclDeleteCommand, args)
}

func TestFileShareAclShowAction(t *testing.T) {
	var args []string
	args = append(args, "d2975ebe-d82c-430f-b28e-f373746a71ca")
	fileShareAclShowAction(fileShareAclShowCommand, args)
}

func TestFileShareAclListAction(t *testing.T) {
	var args []string
	fileSharesAclListAction(fileShareAclListCommand, args)
}
