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

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/opensds/opensds/contrib/connector"
	_ "github.com/opensds/opensds/contrib/connector/iscsi"
	_ "github.com/opensds/opensds/contrib/connector/nvmeof"
)

const (
	iscsiProtocol  = "iscsi"
	nvmeofProtocol = "nvmeof"

	attachCommand = "attach"
	detachCommand = "detach"
)

var (
	connInput string
)

func main() {
	connData := make(map[string]interface{})
	if err := json.Unmarshal([]byte(os.Args[2]), &connData); err != nil {
		fmt.Printf("The format of connection data(%v) is not correct!\n", connData)
		os.Exit(-1)
	}

	accPro := os.Args[3]
	switch os.Args[1] {
	case attachCommand:
		dev, err := connector.NewConnector(accPro).Attach(connData)
		if err != nil {
			fmt.Println("Failed to attach volume to the host:", err)
			os.Exit(-1)
		}
		fmt.Println("Got device:", dev)
		break

	case detachCommand:
		if err := connector.NewConnector(accPro).Detach(connData); err != nil {
			fmt.Println("Failed to detach volume to the host:", err)
			os.Exit(-1)
		}
		fmt.Println("Detach volume success!")
		break
	default:
		fmt.Println(os.Args)
		os.Exit(-1)
	}
}
