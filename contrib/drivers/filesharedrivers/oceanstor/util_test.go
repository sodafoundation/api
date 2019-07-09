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

package oceanstor

import (
	"reflect"
	"testing"
)

var assertTestResult = func(t *testing.T, got, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %v, got: %v\n", expected, got)
	}
}

func TestHandleResponse(t *testing.T) {
	t.Run("error response", func(t *testing.T) {
		sample :=
			`{
		        "data": [
                    {
			            "ID":"12",
			            "IPV4ADDR":"1.2.3.5"
			        },
					{
						"ID":"34",
			            "IPV4ADDR":"3.4.5.6"
					}
		        ],
                "error": {
			        "code":3,
			        "description":"other error"
			    }
		    }`

		var logicalPortList LogicalPortList

		err := handleReponse([]byte(sample), &logicalPortList)
		assertTestResult(t, err.Error(), "other error")
	})

	t.Run("normal response", func(t *testing.T) {
		sample :=
			`{
		        "data": [
                    {
			            "ID":"12",
			            "IPV4ADDR":"1.2.3.5"
			        },
					{
						"ID":"34",
			            "IPV4ADDR":"3.4.5.6"
					}
		        ],
                "error": {
			        "code":0,
			        "description":"0"
			    }
		    }`

		var logicalPortList LogicalPortList

		err := handleReponse([]byte(sample), &logicalPortList)
		assertTestResult(t, err, nil)
	})

	t.Run("no error in response", func(t *testing.T) {
		sample :=
			`{
			        "data": [
                    {
			            "ID":"12",
			            "IPV4ADDR":"1.2.3.5"
			        },
					{
						"ID":"34",
			            "IPV4ADDR":"3.4.5.6"
					}
		        ]
		    }`

		var logicalPortList LogicalPortList

		err := handleReponse([]byte(sample), &logicalPortList)
		assertTestResult(t, err.Error(), "unable to get execution result from response content")
	})
}

func TestFindSpecifiedStruct(t *testing.T) {
	type Sample4 struct {
		Error
		Filed1 bool
	}

	type Sample3 struct {
		Sample4
		Filed1 string
		Filed2 int
	}
	type Sample2 struct {
		Filed1 string
		Sample3
		Filed2 int
	}

	type Sample1 struct {
		Filed1 string
		Filed2 int
		Sample2
	}

	errStruct := Error{
		Code:        1,
		Description: "test error",
	}

	sample4 := Sample4{
		Error:  errStruct,
		Filed1: false,
	}

	sample3 := Sample3{
		Sample4: sample4,
		Filed1:  "test3",
		Filed2:  3,
	}

	sample2 := Sample2{
		Filed1:  "test2",
		Sample3: sample3,
		Filed2:  2,
	}

	sample1 := Sample1{
		Filed1:  "test1",
		Filed2:  1,
		Sample2: sample2,
	}

	t.Run("search substructure named Error from nested structure", func(t *testing.T) {
		result, _ := findSpecifiedStruct("Error", sample1)
		errResult := result.(Error)
		assertTestResult(t, errResult.Description, "test error")
	})

	t.Run("search substructure named Sample3 from nested structure", func(t *testing.T) {
		result, _ := findSpecifiedStruct("Sample3", sample1)
		resultStruct := result.(Sample3)
		assertTestResult(t, resultStruct, sample3)
	})

	t.Run("search substructure named Sample5 from nested structure", func(t *testing.T) {
		result, _ := findSpecifiedStruct("Sample5", sample1)
		assertTestResult(t, result, nil)
	})

	t.Run("search substructure named Sample5 from ptr", func(t *testing.T) {
		result, _ := findSpecifiedStruct("Sample3", &sample1)
		resultStruct := result.(Sample3)
		assertTestResult(t, resultStruct, sample3)
	})
}

func TestCheckProtocol(t *testing.T) {
	t.Run("protocol is nfs", func(t *testing.T) {
		result := checkProtocol(NFSProto)
		assertTestResult(t, result, true)
	})

	t.Run("protocol is cifs", func(t *testing.T) {
		result := checkProtocol(CIFSProto)
		assertTestResult(t, result, true)
	})

	t.Run("protocol is test", func(t *testing.T) {
		result := checkProtocol("test")
		assertTestResult(t, result, false)
	})
}

func TestGetSharePath(t *testing.T) {
	t.Run("get share path", func(t *testing.T) {
		result := getSharePath("fileshare")
		assertTestResult(t, result, "/fileshare/")
	})
}

func TestCheckAccessLevel(t *testing.T) {
	t.Run("check access level", func(t *testing.T) {
		result := checkAccessLevel("rr")
		assertTestResult(t, result, false)
	})

	t.Run("check access level", func(t *testing.T) {
		result := checkAccessLevel("rw")
		assertTestResult(t, result, true)
	})

	t.Run("check access level", func(t *testing.T) {
		result := checkAccessLevel("ro")
		assertTestResult(t, result, true)
	})
}

func TestCheckAccessType(t *testing.T) {
	t.Run("check access level", func(t *testing.T) {
		result := checkAccessType("block")
		assertTestResult(t, result, false)
	})

	t.Run("check access level", func(t *testing.T) {
		result := checkAccessType("ip")
		assertTestResult(t, result, true)
	})

	t.Run("check access level", func(t *testing.T) {
		result := checkAccessType("user")
		assertTestResult(t, result, true)
	})
}

func TestTryTimes(t *testing.T) {
	t.Run("test try times", func(t *testing.T) {
		f := func() error {
			return nil
		}
		result := tryTimes(f)
		assertTestResult(t, result, nil)
	})
}

func TestSector2Gb(t *testing.T) {
	t.Run("test sector to Gb", func(t *testing.T) {
		result := Sector2Gb(12582912)
		assertTestResult(t, result, int64(6))
	})
}

func TestGb2Sector(t *testing.T) {
	t.Run("test Gb to sector", func(t *testing.T) {
		result := Gb2Sector(6)
		assertTestResult(t, result, int64(12582912))
	})
}
