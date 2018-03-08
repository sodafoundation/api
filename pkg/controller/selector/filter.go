// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
)

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

	for key, value := range result {
		valueMap, ok := value.(map[string]interface{})
		if ok {
			for k, v := range valueMap {
				result[key+"."+k] = v
			}
			delete(result, key)
		}
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

// IsAvailablePool ...
func IsAvailablePool(filterReq map[string]interface{}, pool *model.StoragePoolSpec) (bool, error) {
	poolMap, err := StructToMap(pool)
	if nil != err {
		return false, err
	}

	for key, reqValue := range filterReq {
		poolValue, ok := poolMap[key]
		if !ok {
			return false, errors.New("pool doesn't provide capability: " + key)
		}

		ismatch, err := match(key, poolValue, reqValue)
		if nil != err {
			return false, err
		}

		if false == ismatch {
			return false, nil
		}
	}

	return true, nil
}

// match ...
func match(key string, value interface{}, reqValue interface{}) (bool, error) {
	reqValueStr, ok := reqValue.(string)
	if !ok {
		return IsEqual(key, value, reqValue)
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
			return ((A < B) || IsFloatEqual(A, B)), nil
		case ">=":
			return ((A > B) || IsFloatEqual(A, B)), nil
		case "==", "":
			return IsFloatEqual(A, B), nil
		case "!=":
			return (!IsFloatEqual(A, B)), nil
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
