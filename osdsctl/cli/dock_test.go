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

func TestDockAction(t *testing.T) {
	beBrasher := os.Getenv("BE_CRASHER")

	if beBrasher == "1" {
		var args []string
		dockAction(dockCommand, args)

		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestDockAction")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	e, ok := err.(*exec.ExitError)

	if ok && ("exit status 1" == e.Error()) {
		return
	}

	t.Fatalf("process ran with %s, want exit status 1", e.Error())
}

func TestDockShowAction(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.DockMgr), "GetDock",
		func(_ *c.DockMgr, _ string) (*model.DockSpec, error) {
			var res model.DockSpec

			if err := json.Unmarshal([]byte(ByteDock), &res); err != nil {
				return &res, err
			}

			return &res, nil
		})

	var args []string
	args = append(args, "b7602e18-771e-11e7-8f38-dbd6d291f4e0")
	dockShowAction(dockShowCommand, args)
}

func TestDockListAction(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.DockMgr), "ListDocks",
		func(_ *c.DockMgr) ([]*model.DockSpec, error) {
			var res []*model.DockSpec

			if err := json.Unmarshal([]byte(ByteDocks), &res); err != nil {
				return res, err
			}

			return res, nil
		})

	var args []string
	dockListAction(dockListCommand, args)
}
