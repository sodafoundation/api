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

package config

import (
	gflag "flag"
	"fmt"
	"os"
	"testing"
)

func ResetForTesting(usage func()) {
	gflag.CommandLine = gflag.NewFlagSet(os.Args[0], gflag.ContinueOnError)
	gflag.CommandLine.Usage = func() {}
	gflag.Usage = usage
}

func TestFlag(t *testing.T) {
	f := FlagSet{}
	var stringval string
	var boolval bool
	var intval int
	var int64val int64
	var uintval uint
	var uint64val uint64
	var flat64val float64
	f.StringVar(&stringval, "StringVar", "DefualtString", "test")
	f.BoolVar(&boolval, "BoolVar", false, "test")
	f.IntVar(&intval, "IntVar", -321, "test")
	f.Int64Var(&int64val, "Int64Var", -321, "test")
	f.UintVar(&uintval, "UintVar", 321, "test")
	f.Uint64Var(&uint64val, "Uint64Var", 321, "test")
	f.Float64Var(&flat64val, "Float64Var", 0.321, "test")
	cmd := []string{
		"-StringVar", "HelloWorld",
		"-BoolVar",
		"-IntVar", "-123",
		"-Int64Var", "-123",
		"-UintVar", "123",
		"-Uint64Var", "123",
		"-Float64Var", "0.123",
	}
	err := gflag.CommandLine.Parse(cmd)
	if err != nil {
		t.Error(err)
	}
	f.AssignValue()
	if stringval != "HelloWorld" {
		t.Error("Test StringVar Failed!")
	}
	if boolval != true {
		t.Error("Test BoolVar Failed!")
	}
	if intval != -123 {
		t.Error("Test IntVar Failed!")
	}
	if int64val != -123 {
		t.Error("Test Int64Var Failed!")
	}
	if uintval != 123 {
		t.Error("Test UintVar Failed!")
	}
	if uint64val != 123 {
		t.Error("Test Uint64Var Failed!")
	}
	if flat64val != 0.123 {
		t.Error("Test Float64Var Failed!")
	}
}

func TestDefaultValue(t *testing.T) {
	ResetForTesting(func() { t.Fatal("bad parse") })
	f := FlagSet{}
	var stringval string
	f.StringVar(&stringval, "StringVar", "DefualtString", "test")
	err := gflag.CommandLine.Parse([]string{})
	if err != nil {
		t.Error(err)
	}
	f.AssignValue()
	fmt.Println(stringval)
	if stringval != "" {
		t.Error("Test StringVar Failed!")
	}
}
