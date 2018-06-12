// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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
	"flag"
	"fmt"
	"os"

	"github.com/opensds/opensds/contrib/connector"
	_ "github.com/opensds/opensds/contrib/connector/iscsi"
)

const (
	iscsiProtocol = "iscsi"
)

var (
	connInput string
)

func init() {
	flag.StringVar(&connInput, "connection", "", "Connectoin data for attaching the volume to host")
	flag.Parse()
}

func main() {
	if connInput == "" {
		fmt.Println("The number of args is not correct!")
		os.Exit(-1)
	}

	connData := make(map[string]interface{})
	if err := json.Unmarshal([]byte(connInput), &connData); err != nil {
		fmt.Println("The format of connection data is not correct!")
		os.Exit(-1)
	}
	dev, err := connector.NewConnector(iscsiProtocol).Attach(connData)
	if err != nil {
		fmt.Println("Failed to attach volume to the host:", err)
		os.Exit(-1)
	}

	fmt.Println("Got device:", dev)
}
