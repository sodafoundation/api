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
	cryptorand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"math/big"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/sodafoundation/api/pkg/utils/constants"

	log "github.com/golang/glog"
)

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

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

func CompareArray(key string, value interface{}, reqValue interface{}) (bool, error) {
	aInterface := value.([]interface{})
	astring := make([]string, len(aInterface))
	for i, v := range aInterface {
		astring[i] = v.(string)
	}
	switch reflect.TypeOf(reqValue).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(reqValue)
		for i := 0; i < s.Len(); i++ {
			val := s.Index(i).String()
			if !Contains(astring, val) {
				return false, nil
			}
		}
	}
	return true, nil
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

	case []interface{}:
		return CompareArray(key, value, reqValue)

	default:
		return false, errors.New("the type of " + key + " must be bool or float64 or string")
	}
}

// Filter filters slice by struct field
func Filter(arr interface{}, params map[string][]string) interface{} {
	in := reflect.ValueOf(arr)
	out := make([]interface{}, 0, in.Len())

	for i := 0; i < in.Len(); i++ {
		element := in.Index(i).Interface()
		ve := reflect.ValueOf(element).Elem()
		satisfied := true
		for key := range params {
			fieldValue := ve.FieldByName(strings.Title(key))
			if !fieldValue.IsValid() {
				continue
			}
			// Considering multiple input for one parameter, like: ?id=xx1&id=xx2
			for j, value := range params[key] {
				if reflect.DeepEqual(fmt.Sprintf("%v", fieldValue), value) {
					break
				}
				if j == len(params[key])-1 {
					satisfied = false
				}
			}

			if !satisfied {
				break
			}
		}

		if satisfied {
			out = append(out, element)
		}
	}
	return out
}

// Sorting sorts slice with struct field, note: only string and int are supported now
func Sort(arr interface{}, sortKey, sortDir string) interface{} {
	// Sorting
	in := reflect.ValueOf(arr)
	sortKey = strings.Title(sortKey)
	sort.SliceStable(arr, func(i int, j int) bool {
		x := reflect.ValueOf(in.Index(i).Interface()).Elem().FieldByName(sortKey)
		y := reflect.ValueOf(in.Index(j).Interface()).Elem().FieldByName(sortKey)
		if sortDir != constants.DefaultSortDir {
			switch x.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return x.Int() < y.Int()
			}
			return fmt.Sprintf("%v", x) < fmt.Sprintf("%v", y)
		}

		switch x.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return x.Int() > y.Int()
		}
		return fmt.Sprintf("%v", x) > fmt.Sprintf("%v", y)
	})
	return arr
}

// Slicing implements pagination
func Slice(arr interface{}, offset, limit int) interface{} {
	// Slicing
	in := reflect.ValueOf(arr)
	out := make([]interface{}, 0, limit)
	for i := offset; i < in.Len() && i < offset+limit; i++ {
		element := in.Index(i).Interface()
		out = append(out, element)
	}
	return out
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

func WaitForCondition(f func() (bool, error), interval, timeout time.Duration) error {
	endAt := time.Now().Add(timeout)
	time.Sleep(time.Duration(interval))
	for {
		startTime := time.Now()
		ok, err := f()
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		if time.Now().After(endAt) {
			break
		}
		elapsed := time.Now().Sub(startTime)
		time.Sleep(interval - elapsed)
	}
	return fmt.Errorf("wait for condition timeout")
}

func ContainsIgnoreCase(a []string, x string) bool {
	for _, n := range a {
		if strings.EqualFold(x, n) {
			return true
		}
	}
	return false
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret = append(ret, letters[num.Int64()])
	}

	return string(ret), nil
}
