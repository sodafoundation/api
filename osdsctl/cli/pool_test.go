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

func TestPoolAction(t *testing.T) {
	beBrasher := os.Getenv("BE_CRASHER")

	if beBrasher == "1" {
		var args []string
		poolAction(dockCommand, args)

		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestPoolAction")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	e, ok := err.(*exec.ExitError)

	if ok && ("exit status 1" == e.Error()) {
		return
	}

	t.Fatalf("process ran with %s, want exit status 1", e.Error())
}

func TestPoolShowAction(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.PoolMgr), "GetPool",
		func(_ *c.PoolMgr, _ string) (*model.StoragePoolSpec, error) {
			var res model.StoragePoolSpec
			if err := json.Unmarshal([]byte(BytePool), &res); err != nil {
				return &res, err
			}
			return &res, nil
		})

	var args []string
	args = append(args, "084bf71e-a102-11e7-88a8-e31fe6d52248")
	poolShowAction(poolShowCommand, args)
}

func TestPoolListAction(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(client.PoolMgr), "ListPools",
		func(_ *c.PoolMgr) ([]*model.StoragePoolSpec, error) {
			var res []*model.StoragePoolSpec

			if err := json.Unmarshal([]byte(BytePools), &res); err != nil {
				fmt.Println(err)
				return res, err
			}

			return res, nil
		})

	var args []string
	poolListAction(dockListCommand, args)
}
