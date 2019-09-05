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

// +build integration

package integration

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opensds/opensds/test/integration/utils"
)

func TestFileShare(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FileShare Suite")
}

var (
	OPERATION_FAILED = "OPERATION_FAILED"
)

var _ = Describe("FileShare Testing", func() {
	Context("create FileShare ", func() {
		It("TC_FS_IT_01: Create fileshare with name input ", func() {
			var jsonStr = map[string]interface{}{"name": "share2223", "description": "This is just for test222", "size": 2, "profileId": "df40af1a-17b5-48e5-899f-fa098b0bd5da"}
			url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
			methodName := "POST"
			resp, err := utils.ConnectToHTTP(methodName, url, jsonStr)
			Expect(resp.StatusCode).Should(Equal(202))
			Expect(err).NotTo(HaveOccurred())
		})
		It("TC_FS_IT_02: Create fileshare with empty file share name ", func() {
			var jsonStr2 = map[string]interface{}{"name": "", "description": "This is just for testxxx", "size": 2, "profileId": "df40af1a-17b5-48e5-899f-fa098b0bd5da"}
			url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
			methodName := "POST"
			resp, _ := utils.ConnectToHTTP(methodName, url, jsonStr2)
			Expect(resp.StatusCode).Should(Equal(400))
		})
		It("TC_FS_IT_03: Create file share name with other encoding characters(except utf-8) ", func() {
			var jsonStr2 = map[string]interface{}{"name": "İnanç Esasları", "description": "This is just for testxxx", "size": 2, "profileId": "df40af1a-17b5-48e5-899f-fa098b0bd5da"}
			url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
			methodName := "POST"
			resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
			Expect(resp.StatusCode).Should(Equal(202))
			Expect(err).NotTo(HaveOccurred())
		})
		It("TC_FS_IT_04: Create file share name having special characters ", func() {
			var jsonStr2 = map[string]interface{}{"name": "#FileShare Code!$!test", "description": "This is just for testxxx", "size": 2, "profileId": "df40af1a-17b5-48e5-899f-fa098b0bd5da"}
			url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
			methodName := "POST"
			resp, _ := utils.ConnectToHTTP(methodName, url, jsonStr2)
			Expect(resp.StatusCode).Should(Equal(202))
		})
		It("TC_FS_IT_05: Create file share name starts with numbers ", func() {
			var jsonStr2 = map[string]interface{}{"name": "123test", "description": "This is just for testxxx", "size": 2, "profileId": "df40af1a-17b5-48e5-899f-fa098b0bd5da"}
			url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
			methodName := "POST"
			resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
			Expect(resp.StatusCode).Should(Equal(202))
			Expect(err).NotTo(HaveOccurred())
		})
		It("TC_FS_IT_06: Create file share name length more than 255 characters ", func() {
			var jsonStr2 = map[string]interface{}{"name": "abqwqwqwggg012345678910gggggggggggggghhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg", "description": "This is just for testxxx", "size": 2, "profileId": "df40af1a-17b5-48e5-899f-fa098b0bd5da"}
			url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
			methodName := "POST"
			resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
			Expect(resp.StatusCode).Should(Equal(400))
			Expect(err).NotTo(HaveOccurred())
		})
		It("TC_FS_IT_08: Create file share description with empty string ", func() {
			var jsonStr2 = map[string]interface{}{"name": "abcd123", "description": "#FileShare Code!$!test", "size": 2, "profileId": "df40af1a-17b5-48e5-899f-fa098b0bd5da"}
			url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
			methodName := "POST"
			resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
			Expect(resp.StatusCode).Should(Equal(400))
			Expect(err).NotTo(HaveOccurred())
		})

	})
})