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

package config

import (
	gflag "flag"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
	log "github.com/golang/glog"
)

const (
	ConfKeyName = iota
	ConfDefaultValue
)

func setSlice(v reflect.Value, str string) {
	sList := strings.Split(str, ",")
	s := reflect.MakeSlice(v.Type(), 0, 5)
	switch v.Type().Elem().Kind() {
	case reflect.Bool:
		for _, elm := range sList {
			val, _ := strconv.ParseBool(elm)
			s = reflect.Append(s, reflect.ValueOf(val))
		}
	case reflect.Int:
		for _, elm := range sList {
			val, _ := strconv.Atoi(elm)
			s = reflect.Append(s, reflect.ValueOf(val))
		}
	case reflect.Int8:
		for _, elm := range sList {
			val, _ := strconv.ParseInt(elm, 10, 64)
			s = reflect.Append(s, reflect.ValueOf(int8(val)))
		}
	case reflect.Int16:
		for _, elm := range sList {
			val, _ := strconv.ParseInt(elm, 10, 64)
			s = reflect.Append(s, reflect.ValueOf(int16(val)))
		}
	case reflect.Int32:
		for _, elm := range sList {
			val, _ := strconv.ParseInt(elm, 10, 64)
			s = reflect.Append(s, reflect.ValueOf(int32(val)))
		}
	case reflect.Int64:
		for _, elm := range sList {
			val, _ := strconv.ParseInt(elm, 10, 64)
			s = reflect.Append(s, reflect.ValueOf(int64(val)))
		}
	case reflect.Uint:
		for _, elm := range sList {
			val, _ := strconv.ParseUint(elm, 10, 64)
			s = reflect.Append(s, reflect.ValueOf(uint(val)))
		}
	case reflect.Uint8:
		for _, elm := range sList {
			val, _ := strconv.ParseUint(elm, 10, 64)
			s = reflect.Append(s, reflect.ValueOf(uint8(val)))
		}
	case reflect.Uint16:
		for _, elm := range sList {
			val, _ := strconv.ParseUint(elm, 10, 64)
			s = reflect.Append(s, reflect.ValueOf(uint16(val)))
		}
	case reflect.Uint32:
		for _, elm := range sList {
			val, _ := strconv.ParseUint(elm, 10, 64)
			s = reflect.Append(s, reflect.ValueOf(uint32(val)))
		}
	case reflect.Uint64:
		for _, elm := range sList {
			val, _ := strconv.ParseUint(elm, 10, 64)
			s = reflect.Append(s, reflect.ValueOf(uint64(val)))
		}
	case reflect.Float32:
		for _, elm := range sList {
			val, _ := strconv.ParseFloat(elm, 64)
			s = reflect.Append(s, reflect.ValueOf(float32(val)))
		}
	case reflect.Float64:
		for _, elm := range sList {
			val, _ := strconv.ParseFloat(elm, 64)
			s = reflect.Append(s, reflect.ValueOf(val))
		}
	case reflect.String:
		for _, elm := range sList {
			s = reflect.Append(s, reflect.ValueOf(elm))
		}
	default:
		log.Error("Not support this type of slice.")
	}
	v.Set(s)
}

func parseItems(section string, v reflect.Value, cfg *ini.File) {
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
			val, _ := strconv.ParseBool(strVal)
			field.SetBool(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, _ := strconv.ParseInt(strVal, 10, 64)
			field.SetInt(val)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, _ := strconv.ParseUint(strVal, 10, 64)
			field.SetUint(val)
		case reflect.Float32, reflect.Float64:
			val, _ := strconv.ParseFloat(strVal, 64)
			field.SetFloat(val)
		case reflect.String:
			field.SetString(strVal)
		case reflect.Slice:
			setSlice(field, strVal)
		default:
		}
	}
}

func parseSections(cfg *ini.File, t reflect.Type, v reflect.Value) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		section := t.Field(i).Tag.Get("conf")
		if "FlagSet" == field.Type().Name() {
			continue
		}
		if "" == section {
			parseSections(cfg, field.Type(), field)
		}
		parseItems(section, field, cfg)
	}
}

func initConf(confFile string, conf interface{}) {
	cfg, err := ini.Load(confFile)
	if err != nil && confFile != "" {
		log.Info("Read configuration failed, use default value")
	}
	t := reflect.TypeOf(conf)
	v := reflect.ValueOf(conf)
	parseSections(cfg, t, v)

}

// Global Configuration Variable
var CONF *Config = GetDefaultConfig()

//Create a Config and init default value.
func GetDefaultConfig() *Config {
	var conf *Config = new(Config)
	initConf("", conf)
	return conf
}

func (c *Config) Load(confFile string) {
	gflag.StringVar(&confFile, "config-file", confFile, "The configuration file of OpenSDS")
	c.Flag.Parse()
	initConf(confFile, CONF)
	c.Flag.AssignValue()
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
