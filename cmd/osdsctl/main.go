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

/*
This module implements a entry into the OpenSDS CLI service.

*/

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/opensds/opensds/pkg/cli"
)

func main() {
	// Open OpenSDS CLI service log file
	f, err := os.OpenFile("/var/log/opensds/osdsctl.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Errorf("Error opening file:%v", err)
		os.Exit(1)
	}
	defer f.Close()

	// assign it to the standard logger
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := cli.Run(); err != nil {
		panic(err)
	}
}
