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

// +build e2e

package e2e

import (
	"fmt"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

//Function to run the Ginkgo Test
func TestFileShareIntegration(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)

	var _ = ginkgo.BeforeSuite(func() {
		fmt.Println("Before Suite Execution")

	})
	ginkgo.AfterSuite(func() {
		ginkgo.By("After Suite Execution....!")
	})

	ginkgo.RunSpecs(t, "File Share E2E Test Suite")
}
