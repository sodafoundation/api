// Copyright (c) 2019 Huawei Technologies Co., Ltd. All Rights Reserved.
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
	"fmt"
	"os"

	"github.com/opensds/opensds/pkg/utils/pwd"
	"github.com/spf13/cobra"
)

var aesCommand = &cobra.Command{
	Use:   "aes <password>",
	Short: "encryption tool",
	Run:   encrypter,
}

func encrypter(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Usage()
		os.Exit(0)
	}
	ArgsNumCheck(cmd, args, 1)
	// Encrypt the password
	pwdTool := pwd.NewPwdTool("aes")
	pwd, err := pwdTool.Encrypter(args[0])
	fmt.Print(pwd)
	if err != nil {
		Fatalln("Encrypt password error: %v", err)
	}
}
