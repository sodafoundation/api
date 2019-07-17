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

package utils

import (
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"reflect"

	log "github.com/golang/glog"
)

//remove redundant elements
func RvRepElement(arr []string) []string {
	result := []string{}
	for i := 0; i < len(arr); i++ {
		flag := true
		for j := range result {
			if arr[i] == result[j] {
				flag = false
				break
			}
		}
		if flag == true {
			result = append(result, arr[i])
		}
	}
	return result
}

func Contained(obj, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	default:
		return false
	}
	return false
}

func MergeGeneralMaps(maps ...map[string]interface{}) map[string]interface{} {
	var out = make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

func MergeStringMaps(maps ...map[string]string) map[string]string {
	var out = make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Retry(retryNum int, desc string, silent bool, fn func(retryIdx int, lastErr error) error) error {
	var err error
	for i := 0; i < retryNum; i++ {
		if err = fn(i, err); err != nil {
			if !silent {
				log.Errorf("%s:%s, retry %d time(s)", desc, err, i+1)
			}
		} else {
			return nil
		}
	}
	if !silent {
		log.Errorf("%s retry exceed the max retry times(%d).", desc, retryNum)
	}
	return err
}

// StructToMap ...
func StructToMap(structObj interface{}) (map[string]interface{}, error) {
	jsonStr, err := json.Marshal(structObj)
	if nil != err {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonStr, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Epsilon ...
const Epsilon float64 = 0.00000001

// IsFloatEqual ...
func IsFloatEqual(a, b float64) bool {
	if (a-b) < Epsilon && (b-a) < Epsilon {
		return true
	}

	return false
}

// IsEqual ...
func IsEqual(key string, value interface{}, reqValue interface{}) (bool, error) {
	switch value.(type) {
	case bool:
		v, ok1 := value.(bool)
		r, ok2 := reqValue.(bool)

		if ok1 && ok2 {
			if v == r {
				return true, nil
			}

			return false, nil
		}

		return false, errors.New("the type of " + key + " must be bool")
	case float64:
		v, ok1 := value.(float64)
		r, ok2 := reqValue.(float64)

		if ok1 && ok2 {
			if IsFloatEqual(v, r) {
				return true, nil
			}

			return false, nil
		}

		return false, errors.New("the type of " + key + " must be float64")
	case string:
		v, ok1 := value.(string)
		r, ok2 := reqValue.(string)
		if ok1 && ok2 {
			if v == r {
				return true, nil
			}

			return false, nil
		}

		return false, errors.New("the type of " + key + " must be string")
	default:
		return false, errors.New("the type of " + key + " must be bool or float64 or string")
	}
}

func RandSeqWithAlnum(n int) string {
	alnum := []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	return RandSeq(n, alnum)
}

func RandSeq(n int, chs []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = chs[rand.Intn(len(chs))]
	}
	return string(b)
}
