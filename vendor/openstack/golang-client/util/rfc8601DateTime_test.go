// Copyright (c) 2014 Hewlett-Packard Development Company, L.P.
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

package util_test

import (
	"encoding/json"
	"testing"
	"time"

	"git.openstack.org/openstack/golang-client.git/testUtil"
	"git.openstack.org/openstack/golang-client.git/util"
)

var testValue = `{"created_at":"2014-09-29T14:44:31"}`
var testTime, _ = time.Parse(`"2006-01-02T15:04:05"`, `"2014-09-29T14:44:31"`)
var timeTestValue = timeTest{CreatedAt: util.RFC8601DateTime{testTime}}

func TestMarshalTimeTest(t *testing.T) {
	bytes, _ := json.Marshal(timeTestValue)

	testUtil.Equals(t, testValue, string(bytes))
}

func TestUnmarshalValidTimeTest(t *testing.T) {
	val := timeTest{}
	err := json.Unmarshal([]byte(testValue), &val)
	testUtil.IsNil(t, err)
	testUtil.Equals(t, timeTestValue.CreatedAt.Time, val.CreatedAt.Time)
}

func TestUnmarshalInvalidDataFormatTimeTest(t *testing.T) {
	val := timeTest{}
	err := json.Unmarshal([]byte("something other than date time"), &val)
	testUtil.Assert(t, err != nil, "expected an error")
}

// Added this test to ensure that its understood that
// only one specific format is supported at this time.
func TestUnmarshalInvalidDateTimeFormatTimeTest(t *testing.T) {
	val := timeTest{}
	err := json.Unmarshal([]byte("2014-09-29T14:44"), &val)
	testUtil.Assert(t, err != nil, "expected an error")
}

type timeTest struct {
	CreatedAt util.RFC8601DateTime `json:"created_at"`
}
