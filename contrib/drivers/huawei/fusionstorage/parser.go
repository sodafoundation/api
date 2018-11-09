// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package fusionstorage

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ParseResp(data string) map[string]string {
	datamap := map[string]string{}
	items := strings.Split(string(data), ",")
	for _, item := range items {
		kv := strings.SplitN(item, "=", 2)
		sv := ""
		if len(kv) == 2 {
			sv = kv[1]
		}
		datamap[kv[0]] = sv
	}
	return datamap
}

func Value(data map[string]string, v reflect.Value) error {
	for i := 0; i < v.NumField(); i++ {
		iv := v.Field(i)
		it := iv.Type()
		tag := v.Type().Field(i).Tag.Get("fsc")
		if len(tag) == 0 {
			continue
		}
		s, ok := data[tag]
		if !ok {
			continue
		}

		switch iv.Kind() {
		case reflect.Bool:
			val, err := strconv.ParseBool(s)
			if err != nil {
				return fmt.Errorf("can convert %s:%s to integer", tag, s)
			}
			iv.SetBool(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil || reflect.Zero(it).OverflowInt(n) {
				return fmt.Errorf("can convert %s:%s to integer", tag, s)
			}
			iv.SetInt(n)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n, err := strconv.ParseUint(s, 10, 64)
			if err != nil || reflect.Zero(it).OverflowUint(n) {
				return fmt.Errorf("can convert %s:%s to unsigned integer", tag, s)
			}
			iv.SetUint(n)
		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(s, 64)
			if err != nil || reflect.Zero(it).OverflowFloat(f) {
				return fmt.Errorf("can convert %s:%s to float", tag, s)
			}
			iv.SetFloat(f)
		case reflect.String:
			iv.SetString(s)
		default:
			return fmt.Errorf("fsc: Unexpected key type", iv.Kind()) // should never occur
		}
	}
	return nil
}
func unmarshalStruct(data string, v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	datamap := ParseResp(string(data))
	return Value(datamap, v)
}

func unmarshalSlice(data string, v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	elmType := v.Type().Elem()
	if elmType.Kind() == reflect.Ptr {
		elmType = elmType.Elem()
	}
	for _, line := range strings.Split(string(data), "\n") {

		elm := reflect.New(elmType)
		if err := unmarshalStruct(line, elm); err != nil {
			return err
		}
		v.Set(reflect.Append(v, elm))
	}
	return nil
}

func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("invalid type of %s ,reflect.TypeOf(v) ")
	}
	rv = rv.Elem()
	switch rv.Kind() {
	case reflect.Slice:
		return unmarshalSlice(string(data), rv)
	case reflect.Struct:
		return unmarshalStruct(string(data), rv)
	default:
		return fmt.Errorf("unsupported type: %s", rv.Kind())
	}
	return nil
}
