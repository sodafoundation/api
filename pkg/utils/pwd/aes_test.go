// Copyright 2019 The OpenSDS Authors.
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
	"testing"
)

var (
	password = "123456"
	//ciphertext = "9a3b4d577c86d5c9468123b93382c194399070c3f5e3"
)

func TestDecrypter(t *testing.T) {
	var expected = password
	var aes = &AES{}
	ciphertext, _ := aes.Encrypter(password)
	pwdText, _ := aes.Decrypter(ciphertext)
	if pwdText != expected {
		t.Errorf("Expected %s, got %s\n", expected, pwdText)
	}
}
