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

func TestProfileAction(t *testing.T) {
	beBrasher := os.Getenv("BE_CRASHER")

	if beBrasher == "1" {
		var args []string
		profileAction(profileCommand, args)

		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestProfileAction")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	e, ok := err.(*exec.ExitError)

	if ok && ("exit status 1" == e.Error()) {
		return
	}

	t.Fatalf("process ran with %s, want exit status 1", e.Error())
}

func TestProfileCreateAction(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.ProfileMgr), "CreateProfile",
		func(_ *c.ProfileMgr, body c.ProfileBuilder) (*model.ProfileSpec, error) {
			return body, nil
		})

	var args []string
	args = append(args, ByteProfile)
	profileCreateAction(profileCreateCommand, args)
}

func TestProfileShowAction(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.ProfileMgr), "GetProfile",
		func(_ *c.ProfileMgr, prfID string) (*model.ProfileSpec, error) {
			var res model.ProfileSpec
			if err := json.Unmarshal([]byte(ByteProfile), &res); err != nil {
				return nil, err
			}
			res.Id = prfID

			return &res, nil
		})

	var args []string
	args = append(args, "1106b972-66ef-11e7-b172-db03f3689c9c")
	profileShowAction(profileShowCommand, args)
}

func TestProfileListAction(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.ProfileMgr), "ListProfiles",
		func(_ *c.ProfileMgr) ([]*model.ProfileSpec, error) {
			var res []*model.ProfileSpec
			if err := json.Unmarshal([]byte(ByteProfiles), &res); err != nil {
				return nil, err
			}

			return res, nil
		})

	var args []string
	profileListAction(profileListCommand, args)
}

func TestProfileDeleteAction(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.ProfileMgr), "DeleteProfile",
		func(_ *c.ProfileMgr, prfID string) error {
			return nil
		})

	var args []string
	args = append(args, "1106b972-66ef-11e7-b172-db03f3689c9c")
	profileDeleteAction(profileDeleteCommand, args)
}
