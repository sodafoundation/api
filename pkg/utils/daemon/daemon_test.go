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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func mockCmdHandler(name string, args ...string) {
	exec.Command(name, args...).Output()
}

const testFile = "test.txt"

func TestDaemon(t *testing.T) {
	var godaemon bool
	flag.BoolVar(&godaemon, "daemon", false, "Run app as a daemon with -daemon=true")

	execCmdHandler = mockCmdHandler

	var bakArgs = make([]string, len(os.Args))
	copy(bakArgs, os.Args)

	t1 := []string{"testcase1", "-daemon"}
	osArgsHelper(t, t1...)
	CheckAndRunDaemon(true)
	check(t, t1[0])

	t2 := []string{"testcase2", "-daemon=true"}
	osArgsHelper(t, t2...)
	CheckAndRunDaemon(true)
	check(t, t2[0])

	t3 := []string{"testcase3", "-daemon=false"}
	osArgsHelper(t, t3...)
	CheckAndRunDaemon(false)
	check(t)

	t4 := []string{"testcase3", "daemon"}
	osArgsHelper(t, t4...)
	CheckAndRunDaemon(true)
	check(t, t4...)

	os.Remove(testFile)
	os.Args = bakArgs
}

func osArgsHelper(t *testing.T, s ...string) {
	cs := []string{os.Args[0], "-test.run=TestHelperProcess", "--"}
	os.Args = append(cs, s...)
}

func check(t *testing.T, s ...string) {
	buf, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Errorf("File Error: %s\n", err)
	}

	if string(buf) != strings.Join(s, " ") {
		t.Errorf("%s error, got string: %s", s[0], buf)
	}

	err = ioutil.WriteFile(testFile, []byte{}, 0644)
	if err != nil {
		t.Errorf("File Error: %s\n", err)
	}
}

func writeToTestFile(t *testing.T, s string) {
	f, err := os.Create(testFile)
	if err != nil {
		fmt.Fprint(os.Stderr, "An error occurred with file opening or creation\n")
		return
	}
	defer f.Close()
	f.WriteString(s)
}

func TestHelperProcess(t *testing.T) {
	defer os.Exit(0)
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprint(os.Stderr, "No command\n")
		os.Exit(0)
	}

	writeToTestFile(t, strings.Join(args, " "))
}
