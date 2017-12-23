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
	"fmt"
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

func TestVolumeAction(t *testing.T) {
	beBrasher := os.Getenv("BE_CRASHER")

	if beBrasher == "1" {
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
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.VolumeMgr), "CreateVolume",
		func(_ *c.VolumeMgr, body c.VolumeBuilder) (*model.VolumeSpec, error) {
			var res model.VolumeSpec
			if err := json.Unmarshal([]byte(ByteVolume), &res); err != nil {
				return nil, err
			}

			return &res, nil
		})

	var args []string
	args = append(args, "1")
	volumeCreateAction(volumeCreateCommand, args)
}

func TestVolumeShowAction(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.VolumeMgr), "GetVolume",
		func(_ *c.VolumeMgr, volID string) (*model.VolumeSpec, error) {
			var res model.VolumeSpec

			if err := json.Unmarshal([]byte(ByteVolume), &res); err != nil {
				fmt.Println(err)
				return nil, err
			}

			res.Id = volID

			return &res, nil
		})

	var args []string
	args = append(args, "bd5b12a8-a101-11e7-941e-d77981b584d8")
	volumeShowAction(volumeShowCommand, args)
}

func TestVolumeListAction(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.VolumeMgr), "ListVolumes",
		func(_ *c.VolumeMgr) ([]*model.VolumeSpec, error) {
			var res []*model.VolumeSpec
			if err := json.Unmarshal([]byte(ByteVolumes), &res); err != nil {
				return nil, err
			}

			return res, nil
		})

	var args []string
	volumeListAction(volumeListCommand, args)
}

func TestVolumeDeleteAction(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.VolumeMgr), "DeleteVolume",
		func(_ *c.VolumeMgr, volID string, body c.VolumeBuilder) error {
			return nil
		})

	var args []string
	args = append(args, "bd5b12a8-a101-11e7-941e-d77981b584d8")
	volumeDeleteAction(volumeDeleteCommand, args)
}
