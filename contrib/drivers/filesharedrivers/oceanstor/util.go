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
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	log "github.com/golang/glog"
)

func handleReponse(respContent []byte, out interface{}) error {
	if errUnmarshal := json.Unmarshal(respContent, out); errUnmarshal != nil {
		return errUnmarshal
	}

	errStruct, err := findSpecifiedStruct("Error", out)
	if err != nil {
		return err
	}

	errResult := errStruct.(Error)

	if errResult.Description == "" {
		return errors.New("unable to get execution result from response content")
	}

	if errResult.Code != 0 {
		return errors.New(errResult.Description)
	}

	return nil
}

// findSpecifiedStruct  Non-recursive search a specified structure from a nested structure
func findSpecifiedStruct(specifiedStructName string, input interface{}) (interface{}, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}
	if specifiedStructName == "" {
		return nil, errors.New("specified struct name cannot be empty")
	}

	var list []reflect.Value

	list = append(list, reflect.ValueOf(input))

	for len(list) > 0 {
		value := list[0]
		list = append(list[:0], list[1:]...)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		if value.Kind() == reflect.Struct {
			if value.Type().Name() == specifiedStructName {
				return value.Interface(), nil
			}

			for i := 0; i < value.NumField(); i++ {
				list = append(list, value.Field(i))
			}
		}
	}

	return nil, nil
}

func checkProtocol(proto string) bool {
	proList := []string{NFSProto, CIFSProto}
	for _, v := range proList {
		if v == proto {
			return true
		}
	}
	return false
}

func getSharePath(shareName string) string {
	sharePath := "/" + strings.Replace(shareName, "-", "_", -1) + "/"
	return sharePath
}

func checkAccessLevel(accessLevel string) bool {
	accessLevels := []string{AccessLevelRW, AccessLevelRO}
	for _, v := range accessLevels {
		if v == accessLevel {
			return true
		}
	}

	return false
}

func checkAccessType(accessType string) bool {
	accessTypes := []string{AccessTypeUser, AccessTypeIp}
	for _, v := range accessTypes {
		if v == accessType {
			return true
		}
	}

	return false
}

func tryTimes(f func() error) error {
	var err error

	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	for i := 1; i <= MaxRetry; i++ {
		log.Infof("try to exec function %s for %d times", funcName, i)
		err = f()
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		return fmt.Errorf("exec function %s failed: %v", funcName, err)
	}
	return nil
}

func Sector2Gb(sec int64) int64 {
	return sec * 512 / UnitGi
}

func Gb2Sector(gb int64) int64 {
	return gb * UnitGi / 512
}

func EncodeName(id string) string {
	return NamePrefix + "_" + id
}
