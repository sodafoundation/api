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
	"encoding/json"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	c "github.com/opensds/opensds/client"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/testutils/collection"
)

func init() {
	if nil == client {
		ep, ok := os.LookupEnv("OPENSDS_ENDPOINT")

		if !ok {
			ep = "TestEndPoint"
			os.Setenv("OPENSDS_ENDPOINT", ep)
		}

		client = c.NewClient(&c.Config{Endpoint: ep})
	}
}

func TestVolumeSnapshotAction(t *testing.T) {
	beBrasher := os.Getenv("BE_CRASHER")

	if beBrasher == "1" {
		var args []string
		volumeSnapshotAction(volumeSnapshotCommand, args)

		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestVolumeSnapshotAction")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	e, ok := err.(*exec.ExitError)

	if ok && ("exit status 1" == e.Error()) {
		return
	}

	t.Fatalf("process ran with %s, want exit status 1", e.Error())
}

func TestVolumeSnapshotCreateAction(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.VolumeMgr), "CreateVolumeSnapshot",
		func(_ *c.VolumeMgr, body c.VolumeSnapshotBuilder) (*model.VolumeSnapshotSpec, error) {
			var res model.VolumeSnapshotSpec
			if err := json.Unmarshal([]byte(ByteSnapshot), &res); err != nil {
				return nil, err
			}

			return &res, nil
		})

	var args []string
	args = append(args, "bd5b12a8-a101-11e7-941e-d77981b584d8")
	volumeSnapshotCreateAction(volumeSnapshotCreateCommand, args)
}

func TestVolumeSnapshotShowAction(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.VolumeMgr), "GetVolumeSnapshot",
		func(_ *c.VolumeMgr, snpID string) (*model.VolumeSnapshotSpec, error) {
			var res model.VolumeSnapshotSpec

			if err := json.Unmarshal([]byte(ByteSnapshot), &res); err != nil {
				return nil, err
			}

			return &res, nil
		})

	var args []string
	args = append(args, "3769855c-a102-11e7-b772-17b880d2f537")
	volumeSnapshotShowAction(volumeSnapshotShowCommand, args)
}

func TestVolumeSnapshotListAction(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.VolumeMgr), "ListVolumeSnapshots",
		func(_ *c.VolumeMgr) ([]*model.VolumeSnapshotSpec, error) {
			var res []*model.VolumeSnapshotSpec
			if err := json.Unmarshal([]byte(ByteSnapshots), &res); err != nil {
				return nil, err
			}

			return res, nil
		})

	var args []string
	volumeSnapshotListAction(volumeSnapshotListCommand, args)
}

func TestVolumeSnapshotDeleteAction(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.VolumeMgr), "DeleteVolumeSnapshot",
		func(_ *c.VolumeMgr, snpID string, body c.VolumeSnapshotBuilder) error {
			return nil
		})

	var args []string
	args = append(args, "bd5b12a8-a101-11e7-941e-d77981b584d8")
	args = append(args, "3769855c-a102-11e7-b772-17b880d2f537")
	volumeSnapshotDeleteAction(volumeSnapshotDeleteCommand, args)
}
