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

func TestChangeURL(t *testing.T) {
	var org = "v1beta/docks/d3c7e0c7-6e92-406c-9767-3ab73b39b64f"
	var expected = "v1beta/docks/0105e3e4d44d40b59472688a3f28d469"
	if url := ChangeURL(org, "d3c7e0c7-6e92-406c-9767-3ab73b39b64f", "0105e3e4d44d40b59472688a3f28d469"); url != expected {
		t.Errorf("Expected %v, got %v\n", expected, url)
	}

	org = "v1beta/pools/d3c7e0c7-6e92-406c-9767-3ab73b39b64f/8e5e92ca-d673-11e7-bca8-2ba95b86eb06"
	expected = "v1beta/pools/0105e3e4d44d40b59472688a3f28d469/8e5e92ca-d673-11e7-bca8-2ba95b86eb06"
	if url := ChangeURL(org, "d3c7e0c7-6e92-406c-9767-3ab73b39b64f", "0105e3e4d44d40b59472688a3f28d469"); url != expected {
		t.Errorf("Expected %v, got %v\n", expected, url)
	}

	org = "v1beta/d3c7e0c7-6e92-406c-9767-3ab73b39b64f/docks"
	expected = "v1beta/0105e3e4d44d40b59472688a3f28d469/docks"
	if url := ChangeURL(org, "d3c7e0c7-6e92-406c-9767-3ab73b39b64f", "0105e3e4d44d40b59472688a3f28d469"); url != expected {
		t.Errorf("Expected %v, got %v\n", expected, url)
	}

	org = "v1beta/d3c7e0c7-6e92-406c-9767-3ab73b39b64f/pools/8e5e92ca-d673-11e7-bca8-2ba95b86eb06"
	expected = "v1beta/0105e3e4d44d40b59472688a3f28d469/pools/8e5e92ca-d673-11e7-bca8-2ba95b86eb06"
	if url := ChangeURL(org, "d3c7e0c7-6e92-406c-9767-3ab73b39b64f", "0105e3e4d44d40b59472688a3f28d469"); url != expected {
		t.Errorf("Expected %v, got %v\n", expected, url)
	}
}
