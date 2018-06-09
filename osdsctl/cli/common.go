// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

	c "github.com/opensds/opensds/client"
	"github.com/spf13/cobra"
)

const (
	errorPrefix = "ERROR:"
	debugPrefix = "DEBUG:"
	warnPrefix  = "WARNING:"
)

func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(os.Stdout, format, a...)
}

func Debugf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		return fmt.Fprintf(os.Stdout, debugPrefix+" "+format, a...)
	}
	return 0, nil
}

func Warnf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(os.Stdout, warnPrefix+" "+format, a...)
}

func Errorf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(os.Stderr, errorPrefix+" "+format, a...)
}

func Fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, errorPrefix+" "+format, a...)
	os.Exit(-1)
}

func Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stdout, a...)
}

func Debugln(a ...interface{}) (n int, err error) {
	if Debug {
		a = append([]interface{}{debugPrefix}, a...)
		return fmt.Fprintln(os.Stdout, a...)
	}
	return 0, nil
}

func Warnln(a ...interface{}) (n int, err error) {
	a = append([]interface{}{warnPrefix}, a...)
	return fmt.Fprintln(os.Stdout, a...)
}

func Errorln(a ...interface{}) (n int, err error) {
	a = append([]interface{}{errorPrefix}, a...)
	return fmt.Fprintln(os.Stderr, a...)
}

func Fatalln(a ...interface{}) {
	a = append([]interface{}{errorPrefix}, a...)
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(-1)
}

// Strip some redundant message from client http error.
func HttpErrStrip(err error) error {
	if httpErr, ok := err.(*c.HttpError); ok {
		httpErr.Decode()
		return fmt.Errorf(httpErr.Msg)
	}
	return err
}

func ArgsNumCheck(cmd *cobra.Command, args []string, invalidNum int) {
	if len(args) != invalidNum {
		Errorln("The number of args is not correct!")
		cmd.Usage()
		os.Exit(1)
	}
}

func PrintResponse(v interface{}) {
	if Debug {
		b, _ := json.Marshal(v)
		Debugln(string(b))
	}
}
