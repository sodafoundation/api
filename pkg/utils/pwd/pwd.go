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

package pwd

import (
	"fmt"
)

type PwdEncrypter interface {
	Encrypter(password string) (string, error)
	Decrypter(code string) (string, error)
}

func NewPwdEncrypter(encrypter string) PwdEncrypter {
	switch encrypter {
	case "aes":
		return NewAES()
	default:
		fmt.Println("Use default encryption tool: aes.")
		return NewAES()
	}
}
