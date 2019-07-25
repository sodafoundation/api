// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package eternus

import (
	"testing"
)

func TestIsIPv4(t *testing.T) {
	ip := "192.168.1.1"
	expect := true
	result := IsIPv4(ip)
	if result != expect {
		t.Error("Test IsIPv4 failed")
	}
}

func TestGetPortNumber(t *testing.T) {
	caModuleID := "50"
	portNumber := "00"
	expect := "010"
	result := GetPortNumber(caModuleID, portNumber)
	if result != expect {
		t.Error("Test GetPortNumber failed")
	}
}

func TestGetPortNumberV2(t *testing.T) {
	caModuleID := "11"
	portNumber := "00"
	expect := "0110"
	result := GetPortNumberV2(caModuleID, portNumber)
	if result != expect {
		t.Error("Test GetPortNumberV2 failed")
	}
}

func TestParseIPv4(t *testing.T) {
	ip := "C0A80CE3"
	expect := "192.168.12.227"
	result := ParseIPv4(ip)
	if result != expect {
		t.Error("Test ParseIPv4 failed")
	}
}

func TestParseIPv6(t *testing.T) {
	ip := "FE800000000000001111222233334444"
	expect := "FE80:0000:0000:0000:1111:2222:3333:4444"
	result := ParseIPv6(ip)
	if result != expect {
		t.Error("Test ParseIPv6 failed")
	}
}

func TestEqualIP(t *testing.T) {
	ip := "FE80::1111:2222:3333:4444"
	ip2 := "FE80:0000:0000:0000:1111:2222:3333:4444"
	expect := true
	result := EqualIP(ip, ip2)
	if result != expect {
		t.Error("Test EqualIP failed")
	}

	ip = "FE80::1111:2222:3333:4444"
	ip2 = "FE80:0000:0000:0000:1111:2222:3333:5555"
	expect = false
	result = EqualIP(ip, ip2)
	if result != expect {
		t.Error("Test EqualIP failed")
	}
}

func TestGetFnvHash(t *testing.T) {
	str := "teststring"
	expect := "94dd691390fdc1da"
	result := GetFnvHash(str)
	if result != expect {
		t.Error("Test GetFnvHash failed")
	}
}
