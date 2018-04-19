// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package api

import (
	"reflect"
	"strconv"
	"strings"

	log "github.com/golang/glog"
)

const (
	defaultLimit   = -1
	defaultOffset  = 0
	defaultSortDir = "desc"
	defaultSortKey = "ID"
)

//Parameter
type Parameter struct {
	beginIdx, endIdx int
	sortDir, sortKey string
}

//IsInArray
func IsInArray(e string, s []string) bool {
	for _, v := range s {
		if strings.EqualFold(e, v) {
			return true
		}
	}
	return false
}

func SelectOrNot(m map[string][]string) bool {
	for key := range m {
		if key != "limit" && key != "offset" && key != "sortDir" && key != "sortKey" {
			return true
		}
	}
	return false
}

//Get parameter limit
func GetLimit(m map[string][]string) int {
	var limit int
	var err error
	v, ok := m["limit"]
	if ok {
		limit, err = strconv.Atoi(v[0])
		if err != nil || limit < 0 {
			log.Warning("Invalid input limit:", limit, ",use default value instead:50")
			return defaultLimit
		}
	} else {
		log.Warning("The parameter limit is not present,use default value instead:50")
		return defaultLimit
	}
	return limit
}

//Get parameter offset
func GetOffset(m map[string][]string, size int) int {
	var offset int
	var err error
	v, ok := m["offset"]
	if ok {
		offset, err = strconv.Atoi(v[0])

		if err != nil || offset < 0 || offset > size {
			log.Warning("Invalid input offset or input offset is out of bounds:", offset, ",use default value instead:0")

			return defaultOffset
		}

	} else {
		log.Warning("The parameter offset is not present,use default value instead:0")
		return defaultOffset
	}
	return offset
}

//Get parameter sortDir
func GetSortDir(m map[string][]string) string {
	var sortDir string
	v, ok := m["sortDir"]
	if ok {
		sortDir = v[0]
		if !strings.EqualFold(sortDir, "desc") && !strings.EqualFold(sortDir, "asc") {
			log.Warning("Invalid input sortDir:", sortDir, ",use default value instead:desc")
			return defaultSortDir
		}
	} else {
		log.Warning("The parameter sortDir is not present,use default value instead:desc")
		return defaultSortDir
	}
	return sortDir
}

//Get parameter sortKey
func GetSortKey(m map[string][]string, sortKeys []string) string {
	var sortKey string
	v, ok := m["sortKey"]
	if ok {
		sortKey = strings.ToUpper(v[0])
		if !IsInArray(sortKey, sortKeys) {
			log.Warning("Invalid input sortKey:", sortKey, ",use default value instead:ID")
			return defaultSortKey
		}

	} else {
		log.Warning("The parameter sortKey is not present,use default value instead:ID")
		return defaultSortKey
	}
	return sortKey
}

//ParameterFilter
func ParameterFilter(m map[string][]string, size int, sortKeys []string) *Parameter {
	limit := GetLimit(m)
	offset := GetOffset(m, size)
	beginIdx := offset
	endIdx := limit + offset

	// If use not specified the limit return all the items.
	if limit == defaultLimit || endIdx > size {
		endIdx = size
	}

	sortDir := GetSortDir(m)
	sortKey := GetSortKey(m, sortKeys)

	return &Parameter{beginIdx, endIdx, sortDir, sortKey}
}

func Select(filter map[string][]string, slice interface{}) interface{} {
	if !SelectOrNot(filter) {
		return slice
	}

	var sliceNew []interface{}
	var flag bool
	for _, u := range sliceconv(slice) {

		v := reflect.ValueOf(u)

		method := v.MethodByName("FindValue")

		flag = true
		for key := range filter {

			result := method.Call([]reflect.Value{reflect.ValueOf(key), v})
			k := result[0].Interface().(string)
			if k != "" && !strings.EqualFold(filter[key][0], k) {
				flag = false
				break
			}
		}
		if flag {
			sliceNew = append(sliceNew, u)
		}
	}

	return sliceNew
}

func sliceconv(slice interface{}) []interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		panic("param \"slice\" should be on slice value")
	}

	l := v.Len()
	r := make([]interface{}, l)
	for i := 0; i < l; i++ {
		r[i] = v.Index(i).Interface()
	}
	return r
}
