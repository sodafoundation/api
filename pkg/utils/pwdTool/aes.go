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

package pwdTool

import (
	"fmt"
	"os"

	"github.com/howeyc/gopass"
	"github.com/opensds/opensds/pkg/utils/pwd"
)

func main() {
	fmt.Print("Enter password: ")
	plainText, err := gopass.GetPasswdMasked()
	if err != nil {
		fmt.Printf("Input password error: %v\n", err)
		os.Exit(1)
	}

	// Encrypt the password
	pwdTool := pwd.NewPwdTool("aes")
	pwd, err := pwdTool.Encrypter(string(plainText))

	if err != nil {
		fmt.Printf("Encrypt password error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Encrypted password:", pwd)
}
