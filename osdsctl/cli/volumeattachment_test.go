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
	. "github.com/opensds/opensds/testutils/collection"
)

func init() {
	client = c.NewFakeClient(&c.Config{Endpoint: c.TestEp})
}

func TestVolumeAttachmentAction(t *testing.T) {
	beCrasher := os.Getenv("BE_CRASHER")

	if beCrasher == "1" {
		var args []string
		volumeAttachmentAction(volumeAttachmentCommand, args)

		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestVolumeAttachmentAction")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	e, ok := err.(*exec.ExitError)

	if ok && ("exit status 1" == e.Error()) {
		return
	}

	t.Fatalf("process ran with %s, want exit status 1", e.Error())
}

func TestVolumeAttachmentCreateAction(t *testing.T) {
	var args []string
	args = append(args, ByteAttachment)
	volumeAttachmentCreateAction(volumeAttachmentCreateCommand, args)
}

func TestVolumeAttachmentShowAction(t *testing.T) {
	var args []string
	args = append(args, "f2dda3d2-bf79-11e7-8665-f750b088f63e")
	volumeAttachmentShowAction(volumeAttachmentShowCommand, args)
}

func TestVolumeAttachmentListAction(t *testing.T) {
	var args []string
	volumeAttachmentListAction(volumeAttachmentListCommand, args)
}

func TestVolumeAttachmentDeleteAction(t *testing.T) {
	var args []string
	args = append(args, "f2dda3d2-bf79-11e7-8665-f750b088f63e")
	volumeAttachmentDeleteAction(volumeAttachmentDeleteCommand, args)
}

func TestVolumeAttachmentUpdateAction(t *testing.T) {
	var args []string
	args = append(args, "f2dda3d2-bf79-11e7-8665-f750b088f63e")
	args = append(args, ByteAttachment)
	volumeAttachmentUpdateAction(volumeAttachmentDeleteCommand, args)
}
