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

package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
)

var execCmdHandler = execCmd

func execCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "%s [PID] %d runn in daemon failed.\n,[Error]%s\n",
			os.Args[0], cmd.Process.Pid, err)
		os.Exit(1)
	}
	fmt.Printf("%s [PID] %d running...\n", os.Args[0], cmd.Process.Pid)
	os.Exit(0)
}

func CheckAndRunDaemon(isDaemon bool) {
	if !isDaemon {
		return
	}
	if runtime.GOOS == "windows" {
		fmt.Fprintf(os.Stderr, "Windows does not support daemon mode.")
		return
	}
	var args []string
	for _, arg := range os.Args[1:] {
		if m, _ := regexp.MatchString(`^-{1,2}daemon(=true)?$`, arg); m {
			continue
		}
		args = append(args, arg)
	}
	execCmdHandler(os.Args[0], args...)
}
