// Copyright 2017 The OpenSDS Authors.
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

package selector

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
)

// simplifyPoolCapabilityMap ...
func simplifyPoolCapabilityMap(input map[string]interface{}) (map[string]interface{}, map[string]interface{}) {
	simpleMap := make(map[string]interface{})
	unSimpleMap := make(map[string]interface{})

	for key, value := range input {
		valueMap, ok := value.(map[string]interface{})
		if ok {
			for k, v := range valueMap {
				_, ok = v.(map[string]interface{})
				if ok {
					unSimpleMap[key+"."+k] = v
				} else {
					simpleMap[key+"."+k] = v
				}
			}

		} else {
			simpleMap[key] = value
		}
	}

	return simpleMap, unSimpleMap
}

// GetPoolCapabilityMap ...
func GetPoolCapabilityMap(pool *model.StoragePoolSpec) (map[string]interface{}, error) {
	temMap, err := utils.StructToMap(pool)
	if nil != err {
		return nil, err
	}

	result := make(map[string]interface{})

	// There is no infinite loop here, so set the maximum number of loops to 10.
	for i := 0; i < 10; i++ {
		simpleMap, unSimpleMap := simplifyPoolCapabilityMap(temMap)

		if 0 != len(result) {
			for key, value := range simpleMap {
				result[key] = value
			}
		} else {
			result = simpleMap
		}

		if 0 == len(unSimpleMap) {
			return result, nil
		}

		temMap = unSimpleMap
	}

	return result, nil
}

// IsAvailablePool ...
func IsAvailablePool(filterReq map[string]interface{}, pool *model.StoragePoolSpec) (bool, error) {
	poolMap, err := GetPoolCapabilityMap(pool)
	if nil != err {
		return false, err
	}

	for key, reqValue := range filterReq {
		if strings.HasPrefix(key, ":") {
			log.Info("Because " + key + " is prefixed with a colon, it is not used to filter the pool")
			continue
		}

		poolValue, ok := poolMap[key]
		if !ok {
			log.Error("pool: " + pool.Name + " doesn't provide capability: " + key)
			return false, nil
		}
		ismatch, err := match(key, poolValue, reqValue)
		if nil != err {
			log.Errorf("[%v]: The request value %v is not match the pool value %v.", key, reqValue, poolValue)
			return false, err
		}

		if !ismatch {
			return false, nil
		}
	}

	return true, nil
}

// match ...
func match(key string, value interface{}, reqValue interface{}) (bool, error) {
	reqValueStr, ok := reqValue.(string)
	if !ok {
		return utils.IsEqual(key, value, reqValue)
	}

	words := strings.Split(reqValueStr, " ")
	wordsLen := len(words)

	switch words[0] {
	case "<or>":
		return OrOperator(key, words, value)
	case "=", "==", "!=", ">=", "<=":
		if 2 == wordsLen {
			return ParseFloat64AndCompare(words[0], key, value, words[1])
		}

		return false, errors.New("the format of " + key + ": " + reqValueStr + " is incorrect")
	case "<in>":
		if 2 == wordsLen {
			return InOperator(key, words[1], value)
		}

		return false, errors.New("the format of " + key + ": " + reqValueStr + " is incorrect")
	case "<is>":
		if 2 == wordsLen {
			return ParseBoolAndCompare(key, value, words[1])
		}

		return false, errors.New("the format of " + key + ": " + reqValueStr + " is incorrect")
	case "s==", "s!=", "s<", "s<=", "s>", "s>=":
		if 2 == wordsLen {
			valueStr, ok := value.(string)
			if ok {
				return StringCompare(words[0], key, valueStr, words[1])
			}

			return false, errors.New(key + "is not a string")
		}

		return false, errors.New("the format of " + key + ": " + reqValueStr + " is incorrect")
	default:
		return CompareOperator("", key, words[0], value)
	}
}

// InOperator ...
func InOperator(key string, reqValue string, value interface{}) (bool, error) {
	valueStr, ok := value.(string)
	if !ok {
		return false, errors.New(key + " is not a string, so <in> can not be used")
	}

	isIn, err := regexp.MatchString(reqValue, valueStr)

	if nil != err {
		return false, err
	}

	if isIn {
		return true, nil
	}

	return false, nil
}

// CompareOperator ...
func CompareOperator(op string, key string, reqValue string, value interface{}) (bool, error) {
	switch value.(type) {
	case bool:
		return ParseBoolAndCompare(key, value, reqValue)
	case float64:
		return ParseFloat64AndCompare(op, key, value, reqValue)
	case string:
		valueStr, ok := value.(string)
		if ok {
			return StringCompare(op, key, valueStr, reqValue)
		}

		log.Error("unknown .(string) error")
		return false, nil
	default:
		return false, errors.New("The type of " + key + " must be bool or float64 or string")
	}
}

// StringCompare ...
func StringCompare(op string, key string, a string, b string) (bool, error) {
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	switch op {
	case "s==", "":
		return a == b, nil
	case "s!=":
		return a != b, nil
	case "s>=":
		return a >= b, nil
	case "s>":
		return a > b, nil
	case "s<=":
		return a <= b, nil
	case "s<":
		return a < b, nil
	default:
		return false, errors.New("the operator of string can not be " + op)
	}
}

// ParseBoolAndCompare ...
func ParseBoolAndCompare(key string, a interface{}, b string) (bool, error) {
	B, err := strconv.ParseBool(b)
	if err != nil {
		return false, errors.New("capability is: " + key + ", " + b + " is not bool")
	}

	A, ok := a.(bool)
	if ok {
		if A == B {
			return true, nil
		}

		return false, nil
	}

	return false, errors.New("the value of " + key + " is not bool")
}

// ParseFloat64AndCompare ...
func ParseFloat64AndCompare(op string, key string, a interface{}, b string) (bool, error) {
	B, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return false, errors.New("capability is: " + key + ", " + b + " is not float64")
	}

	A, ok := a.(float64)
	if ok {
		switch op {
		case "<=":
			return ((A < B) || utils.IsFloatEqual(A, B)), nil
		case ">=":
			return ((A > B) || utils.IsFloatEqual(A, B)), nil
		case "==", "":
			return utils.IsFloatEqual(A, B), nil
		case "!=":
			return (!utils.IsFloatEqual(A, B)), nil
		default:
			return false, errors.New("the operator of float64 can not be " + op)
		}
	} else {
		return false, errors.New("the value of " + key + " is not float64")
	}
}

// OrOperator ...
func OrOperator(key string, words []string, value interface{}) (bool, error) {
	wordsLen := len(words)

	if 0 == (wordsLen%2) && wordsLen >= 2 {
		for i := 0; i < wordsLen; i = i + 2 {
			if "<or>" != words[i] {
				return false, errors.New("the first operator is <or>, the following operators must be <or>")
			}
		}

		for i := 1; i < wordsLen; i = i + 2 {
			isOk, err := CompareOperator("", key, words[i], value)

			if nil != err {
				return false, err
			}

			if isOk {
				return true, nil
			}
		}
	} else {
		return false, errors.New("when using <or> as an operator, the <or> and value must appear in pairs")
	}

	return false, nil
}
