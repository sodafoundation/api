// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

/*
This module implements the common data structure.

*/

package api

import (
	"reflect"
	"testing"
)

func TestListVersions(t *testing.T) {
	aVersions, err := ListVersions()
	if err != nil {
		t.Fatal(err)
	}

	expectedResult := VersionInfo{
		Id:     "v1",
		Status: "SUPPORTED",
	}
	if !reflect.DeepEqual(expectedResult, aVersions.Versions[0]) {
		t.Fatalf("Expected\n%#v\ngot\n%#v", expectedResult, aVersions.Versions[1])
	}
}

func TestGetVersionv1(t *testing.T) {
	version, err := GetVersionv1()
	if err != nil {
		t.Fatal(err)
	}

	expectedResult := VersionInfo{
		Id:     "v1",
		Status: "SUPPORTED",
	}
	if !reflect.DeepEqual(expectedResult, version) {
		t.Fatalf("Expected\n%#v\ngot\n%#v", expectedResult, version)
	}
}
