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
	utils "github.com/opensds/opensds/test/integration/utils"
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

		BeforeEach(func() {

		})
		AfterEach(func() {
		})
		It("TC_FS_IT_01: Create fileshare with name input ", func() {
			var jsonStr = map[string]interface{}{"name": "share2223", "description": "This is just for test222", "size": 2}

			// var jsonStr1 = []byte(`{"name":"share22234", "description":"This is just for test222", "size": 2}`)

			url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares" //curl -X POST -H "Content-Type: application/json" -d '{"name":"share1", "description":"This is just for test", "size": 1}' -url "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
			methodName := "POST"
			resp, err := utils.ConnectToHTTP(methodName, url, jsonStr)

			Expect(resp.StatusCode).Should(Equal(202))
			Expect(err).NotTo(HaveOccurred())

		})
		It("TC_FS_IT_02: Create fileshare with duplicate name input ", func() {
			var jsonStr2 = map[string]interface{}{"name": "sharexxx", "description": "This is just for testxxx", "size": 2}
			url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares" //curl -X POST -H "Content-Type: application/json" -d '{"name":"share1", "description":"This is just for test", "size": 1}' -url "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
			methodName := "POST"
			utils.ConnectToHTTP(methodName, url, jsonStr2)
		})
		// It("has 0 units", func() {})
		// Specify("the total amount is 0.00", func() {})
	})
	Context("Get FileShare ", func() {
		// var jsonStr1 = []byte(`{"name":"share2223", "description":"This is just for test222", "size": 2}`)
		// var jsonStr = map[string]interface{}{"name": "share2223", "description": "This is just for test222", "size": 2}
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("TC_FS_IT_03: fileshare GET all ", func() {
			url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares" //curl -X POST -H "Content-Type: application/json" -d '{"name":"share1", "description":"This is just for test", "size": 1}' -url "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
			methodName := "GET"
			utils.ConnectToHTTP(methodName, url, nil)
		})
		It("TC_FS_IT_04: fileshare GET of specific Id", func() {
			fId := "v1beta/file/shares/e93b4c0934da416eb9c8d120c5d04d96/f2ab9308-f208-40c6-bb1f-6fbfa8bf14b5"
			url := "http://127.0.0.1.123:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares/f2ab9308-f208-40c6-bb1f-6fbfa8bf14b5" //curl -X POST -H "Content-Type: application/json" -d '{"name":"share1", "description":"This is just for test", "size": 1}' -url "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
			methodName := "GET"
			utils.ConnectToHTTP(methodName, url, nil)
			// ctx, kv := utils.ConnectToDB()
			ret := utils.GetValueByKeyFromDB(fId)
			Expect(ret).ShouldNot(Equal(OPERATION_FAILED))
			textFound := utils.ReadAndFindTextInFile("C:/go/src/opensds/opensds/test/integration/utils/output.json", "17c60641-63c9-4f7f-992a-c0dcd9abd502")
			Expect(textFound).To(BeTrue(), "Text found in the log file")

		})
		// It("has 0 units", func() {})
		// Specify("the total amount is 0.00", func() {})
	})
})
