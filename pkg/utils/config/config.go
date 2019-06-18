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

package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-ini/ini"
	"github.com/opensds/opensds/pkg/utils/constants"
)

const (
	ConfKeyName = iota
	ConfDefaultValue
)

func setSlice(v reflect.Value, str string) error {
	sList := strings.Split(str, ",")
	s := reflect.MakeSlice(v.Type(), 0, 5)
	switch v.Type().Elem().Kind() {
	case reflect.Bool:
		for _, elm := range sList {
			val, err := strconv.ParseBool(elm)
			if err != nil {
				return fmt.Errorf("cann't convert slice item %s to Bool, %v", elm, err)
			}
			s = reflect.Append(s, reflect.ValueOf(val))
		}
	case reflect.Int:
		for _, elm := range sList {
			val, err := strconv.Atoi(elm)
			if err != nil {
				return fmt.Errorf("cann't convert slice item %s to Int, %v", elm, err)
			}
			s = reflect.Append(s, reflect.ValueOf(val))
		}
	case reflect.Int8:
		for _, elm := range sList {
			val, err := strconv.ParseInt(elm, 10, 64)
			if err != nil {
				return fmt.Errorf("cann't convert slice item %s to Int8, %v", elm, err)
			}
			s = reflect.Append(s, reflect.ValueOf(int8(val)))
		}
	case reflect.Int16:
		for _, elm := range sList {
			val, err := strconv.ParseInt(elm, 10, 64)
			if err != nil {
				return fmt.Errorf("cann't convert slice item %s to Int16, %v", elm, err)
			}
			s = reflect.Append(s, reflect.ValueOf(int16(val)))
		}
	case reflect.Int32:
		for _, elm := range sList {
			val, err := strconv.ParseInt(elm, 10, 64)
			if err != nil {
				return fmt.Errorf("cann't convert slice item %s to Int32, %v", elm, err)
			}
			s = reflect.Append(s, reflect.ValueOf(int32(val)))
		}
	case reflect.Int64:
		for _, elm := range sList {
			val, err := strconv.ParseInt(elm, 10, 64)
			if err != nil {
				return fmt.Errorf("cann't convert slice item %s to Int64, %v", elm, err)
			}
			s = reflect.Append(s, reflect.ValueOf(int64(val)))
		}
	case reflect.Uint:
		for _, elm := range sList {
			val, err := strconv.ParseUint(elm, 10, 64)
			if err != nil {
				return fmt.Errorf("cann't convert slice item %s to Uint, %v", elm, err)
			}
			s = reflect.Append(s, reflect.ValueOf(uint(val)))
		}
	case reflect.Uint8:
		for _, elm := range sList {
			val, err := strconv.ParseUint(elm, 10, 64)
			if err != nil {
				return fmt.Errorf("cann't convert slice item %s to Uint8, %v", elm, err)
			}
			s = reflect.Append(s, reflect.ValueOf(uint8(val)))
		}
	case reflect.Uint16:
		for _, elm := range sList {
			val, err := strconv.ParseUint(elm, 10, 64)
			if err != nil {
				return fmt.Errorf("cann't convert slice item %s to Uint16, %v", elm, err)
			}
			s = reflect.Append(s, reflect.ValueOf(uint16(val)))
		}
	case reflect.Uint32:
		for _, elm := range sList {
			val, err := strconv.ParseUint(elm, 10, 64)
			if err != nil {
				return fmt.Errorf("cann't convert slice item %s to Uint32, %v", elm, err)
			}
			s = reflect.Append(s, reflect.ValueOf(uint32(val)))
		}
	case reflect.Uint64:
		for _, elm := range sList {
			val, err := strconv.ParseUint(elm, 10, 64)
			if err != nil {
				return fmt.Errorf("cann't convert slice item %s to Uint64, %v", elm, err)
			}
			s = reflect.Append(s, reflect.ValueOf(uint64(val)))
		}
	case reflect.Float32:
		for _, elm := range sList {
			val, err := strconv.ParseFloat(elm, 64)
			if err != nil {
				return fmt.Errorf("cann't convert slice item %s to Float32, %v", elm, err)
			}
			s = reflect.Append(s, reflect.ValueOf(float32(val)))
		}
	case reflect.Float64:
		for _, elm := range sList {
			val, err := strconv.ParseFloat(elm, 64)
			if err != nil {
				return fmt.Errorf("cann't convert slice item %s to Float54, %v", elm, err)
			}
			s = reflect.Append(s, reflect.ValueOf(val))
		}
	case reflect.String:
		for _, elm := range sList {
			s = reflect.Append(s, reflect.ValueOf(elm))
		}
	default:
		log.Printf("[ERROR] Does not support this type of slice.")
	}
	v.Set(s)
	return nil
}

func parseItems(section string, v reflect.Value, cfg *ini.File) error {
	for i := 0; i < v.Type().NumField(); i++ {

		field := v.Field(i)
		tag := v.Type().Field(i).Tag.Get("conf")
		if "" == tag {
			parseSections(cfg, field.Type(), field)
		}
		tags := strings.SplitN(tag, ",", 2)
		if !field.CanSet() {
			continue
		}
		var strVal = ""
		if len(tags) > 1 {
			strVal = tags[ConfDefaultValue]
		}
		if cfg != nil {
			key, err := cfg.Section(section).GetKey(tags[ConfKeyName])
			if err == nil {
				strVal = key.Value()
			}
		}
		switch field.Kind() {
		case reflect.Bool:
			val, err := strconv.ParseBool(strVal)
			if err != nil {
				return fmt.Errorf("cann't convert %s:%s to Bool, %v", tags[0], strVal, err)
			}
			field.SetBool(val)
		case reflect.ValueOf(time.Second).Kind():
			if field.Type().String() == "time.Duration" {
				v, err := time.ParseDuration(strVal)
				if err != nil {
					return fmt.Errorf("cann't convert %s:%s to Duration, %v", tags[0], strVal, err)
				}
				field.SetInt(int64(v))
				break
			}
			fallthrough
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				return fmt.Errorf("cann't convert %s:%s to Int, %v", tags[0], strVal, err)
			}
			field.SetInt(val)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := strconv.ParseUint(strVal, 10, 64)
			if err != nil {
				return fmt.Errorf("cann't convert %s:%s to Uint, %v", tags[0], strVal, err)
			}
			field.SetUint(val)
		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				return fmt.Errorf("cann't convert %s:%s to Float, %v", tags[0], strVal, err)
			}
			field.SetFloat(val)
		case reflect.String:
			field.SetString(strVal)
		case reflect.Slice:
			setSlice(field, strVal)
		default:
		}
	}
	return nil
}

func parseSections(cfg *ini.File, t reflect.Type, v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		section := t.Field(i).Tag.Get("conf")
		if "" == section {
			if err := parseSections(cfg, field.Type(), field); err != nil {
				return err
			}
		}
		if err := parseItems(section, field, cfg); err != nil {
			return err
		}
	}
	return nil
}

func initConf(confFile string, conf interface{}) {
	cfg, err := ini.Load(confFile)
	if err != nil && confFile != "" {
		log.Printf("[ERROR] Read configuration failed, use default value")
	}
	t := reflect.TypeOf(conf)
	v := reflect.ValueOf(conf)
	if err := parseSections(cfg, t, v); err != nil {
		log.Fatalf("[ERROR] parse configure file failed: %v", err)
	}
}

// Global Configuration Variable
var CONF *Config = GetDefaultConfig()

//Create a Config and init default value.
func GetDefaultConfig() *Config {
	var conf *Config = new(Config)
	initConf("", conf)
	return conf
}

func GetConfigPath() string {
	path := constants.OpensdsConfigPath
	for i := 1; i < len(os.Args)-1; i++ {
		if m, _ := regexp.MatchString(`^-{1,2}config-file$`, os.Args[i]); m {
			if !strings.HasSuffix(os.Args[i+1], "-") {
				path = os.Args[i+1]
			}
		}
	}
	return path
}

func (c *Config) Load() {
	var dummyConfigPath string
	// Flag 'config-file' is set here for usage show and parameter check, the config path will be parsed by GetConfigPath
	flag.StringVar(&dummyConfigPath, "config-file", constants.OpensdsConfigPath, "OpenSDS config file path")
	initConf(GetConfigPath(), CONF)
}

func GetBackendsMap() map[string]BackendProperties {
	backendsMap := map[string]BackendProperties{}
	v := reflect.ValueOf(CONF.Backends)
	t := reflect.TypeOf(CONF.Backends)

	for i := 0; i < t.NumField(); i++ {
		feild := v.Field(i)
		name := t.Field(i).Tag.Get("conf")
		backendsMap[name] = feild.Interface().(BackendProperties)
	}
	return backendsMap
}
