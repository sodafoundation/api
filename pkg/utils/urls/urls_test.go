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

package urls

import (
	"testing"
)

func TestCurrentVersion(t *testing.T) {
	var expected = "v1beta"

	if version := CurrentVersion(); version != expected {
		t.Errorf("Expected %v, got %v\n", expected, version)
	}
}

func TestGenerateURL(t *testing.T) {
	var expected = "v1beta/docks/d3c7e0c7-6e92-406c-9767-3ab73b39b64f"
	if url := generateURL("docks", Etcd, "d3c7e0c7-6e92-406c-9767-3ab73b39b64f"); url != expected {
		t.Errorf("Expected %v, got %v\n", expected, url)
	}
	expected = "v1beta/docks"
	if url := generateURL("docks", Etcd, ""); url != expected {
		t.Errorf("Expected %v, got %v\n", expected, url)
	}
	expected = "v1beta/pools/d3c7e0c7-6e92-406c-9767-3ab73b39b64f/8e5e92ca-d673-11e7-bca8-2ba95b86eb06"
	if url := generateURL("pools", Etcd, "d3c7e0c7-6e92-406c-9767-3ab73b39b64f", "8e5e92ca-d673-11e7-bca8-2ba95b86eb06"); url != expected {
		t.Errorf("Expected %v, got %v\n", expected, url)
	}

	expected = "v1beta/d3c7e0c7-6e92-406c-9767-3ab73b39b64f/docks"
	if url := generateURL("docks", Client, "d3c7e0c7-6e92-406c-9767-3ab73b39b64f"); url != expected {
		t.Errorf("Expected %v, got %v\n", expected, url)
	}
	expected = "v1beta/d3c7e0c7-6e92-406c-9767-3ab73b39b64f/pools/8e5e92ca-d673-11e7-bca8-2ba95b86eb06"
	if url := generateURL("pools", Client, "d3c7e0c7-6e92-406c-9767-3ab73b39b64f", "8e5e92ca-d673-11e7-bca8-2ba95b86eb06"); url != expected {
		t.Errorf("Expected %v, got %v\n", expected, url)
	}
}
